#!/bin/python3
import pexpect
from pexpect import pxssh
import os
from getpass import getpass
import time
import sys
import subprocess
import re
import socket
import requests
import logging
import warnings

logging.getLogger("requests").setLevel(logging.FATAL)
logging.getLogger("subprocess").setLevel(logging.FATAL)

########################################################################Function Definitions###########################################################################################################3

def get_password(): # Function to re-prompt for a password if it appears that the esxi host is not accepting the original one given
    global times_password_prompted
    global esxi_password
    times_password_prompted=0
    esxi_password = getpass("Please enter the root password for the esxi host: ")

def write_to_log(path='/tmp/esxi_update_log'):
    global ssh_connection
    print('write_to_log function called')
    logfile_read=open(path, 'a')
    print(ssh_connection.before)
    logfile_read.write(ssh_connection.before)
    logfile_read.flush()
    logfile_read.close()

def accept_fingerprint(spawn): # Conditionally accept new server fingerprint based on user response
    if isinstance(spawn, pexpect.spawn):
        response = input('\033[01;34mThe host @'+ esxi_hostname + ' does not have a recognized fingerprint. Would you like to add the fingerprint to the known hosts [yes/no]? ')
        if re.match('yes|y|Y', response):
            spawn.sendline('yes')
            return
        else:
            print('\033[01;33mHost fingerprint not accepted! Exiting script...\033[00m')
            spawn.close()
            exit(10)
    else:
        print('Provided object is not and instance of pexpect.spawn nor its subclasses!')
        exit(11)

def wait_for_reboot(ip_address): # Wait up to 900 seconds for the esxi host to come back up after a reboot
    std_out = ''
    seconds_elapsed=0
    result=-1
    while seconds_elapsed <= 900:
        try:
            result = subprocess.run(['ping', '-c', '1', ip_address], stdout=subprocess.DEVNULL, timeout=1).returncode
            if result == 0:
                print('\033[01;32mESXi host has finished rebooting, the ESXi host update is now complete\033[00m')
                break
        except:
            if seconds_elapsed == 0:
                print('\033[2J')
            print('\033[0;0H\033[01;34mWaiting '+ str(seconds_elapsed) + ' seconds for esxi host to finish reboot...\033[00m')
            #time.sleep(1)
            seconds_elapsed+=1
            if seconds_elapsed>=900:
                print('\033[01;33mScript timed out waiting for ESXi host to come back up, please investigate further\033[00m')
                exit(7)
    exit(0)

def scp_file(username, hostname, password, source, destination): # copy file to remote server based on local source, remote destination, server, username and hostname
    if not os.path.exists(source):
        print('\033[01;31mEsxi depot file not located on the CSS hdd! Depuyt file must be placed at: \033[00m' + source)
    global times_password_prompted
    times_password_prompted=0
    scp_process = pexpect.spawn('scp '+ source  +' ' + username + '@' + hostname + ':' + destination, encoding='utf-8') # use pexpect to create a child scp process that can be controlled and passing in all the necessary values
    scp_process.logfile_read = open(log_file, 'a')
    while scp_process.isalive():
        scp_index = scp_process.expect (['[Pp]assword','Are you sure you want to continue connecting', pexpect.EOF, pexpect.TIMEOUT]) # Scans the output of the child process listening for the patterns passed to it
        if scp_index == 0: # we have recieved some form of password prompt 
            times_password_prompted+=1
            if times_password_prompted >= 2: # Check to see if we have been prompted twice, if so break process and allow user to re-enter the password
                print('\033[01;33mIt appears the password is incorrect, please verify the password and run this script again!\033[01;33m')
                exit(1)
            else:
                scp_process.logfile_read.close()
                scp_process.logfile_read = None
                print("\033[01;34mSending root password...\033[00m")
                scp_process.sendline(esxi_password)
                scp_process.logfile_read = open(log_file, 'a')
                scp_index = scp_process.expect(['scp: /vmfs/volumes/Datastore01/Updates/VMware-ESXi-depot.zip: No such file or directory', 'VMware-ESXI-depot.zip$', pexpect.EOF, pexpect.TIMEOUT])
                if scp_index == 0:
                    print('\033[01;31mUnable to copy the depot file to the esxi host, please ensure that the directory structure: /vmfs/volumes/Datastore01/Updates exists and is writable!\033[00m')
                    exit (1)
                elif scp_index == 1:
                    print('\033[01;32mScp of the depot file was successful!\033[00m')
        elif scp_index == 1: # add esxi fingerprint to ~/.ssh/known_hosts if it isn't there already for some reason
            accept_fingerprint(scp_process)
        elif scp_index == 2:
            scp_process.close(force=True)
            break
        elif scp_index == 3:
            print('Connection timed out with the scp process. Logs are located at ' +  log_file)
            exit(1)
    scp_process.close()
    
def find_pattern(pattern, search_object, is_file=False):
    if not isinstance(pattern, str):
        print('This function must take a string object as the expression pattern')
        return 1
    else:
        regex = re.compile(pattern)
    if is_file:
        my_file = open(search_object, 'r')
        to_search = my_file.read()
        my_file.close()
    elif isinstance(search_object, str):
        to_search = search_object
    else:
        print('This function only takes files or strings to search')
        return 2
    result = regex.search(to_search)
    if result != None:
        return result.group()
    else:
        return None

def replace_pattern(original_pattern, replacement_pattern,  replacement_object):
    result = re.sub(original_pattern, replacement_pattern, replacement_object)
    if result != replacement_object:
        return result
    else:
        print('No replacements were made!')
        return result

def check_maintenance_mode(ssh_connection):
    print('\033[01;34mEnsuring the host is in maintenance mode...\033[00m')
    if not isinstance(ssh_connection, pexpect.spawn):
        print('You most provide this function with an instance of a pexpect.spawn object!')
        exit(1)
    else:
        ssh_connection.sendline('vim-cmd hostsvc/hostsummary')
        match_prompt(ssh_connection)
        local_log_file = open(log_file, 'r')
        cmd_result = local_log_file.read()
        local_log_file.close()
        result = re.search('inMaintenanceMode = [Tt]rue', cmd_result)
        if result == None:
            print('\033[01;33m\nThe ESXi host is not in maintenance mode!\nPlease shut down all virtual machines and put the host in maintenance mode before continuing!\033[00m')
            exit(1)
        else:
            print('\033[01;32mConfirmed host is in maintenance mode\033[00m')
            return

def get_ip():
    s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    s.settimeout(0)
    try:
        s.connect(('10.254.254.254', 1))
        IP = s.getsockname()[0]
    except Exception:
        IP = '127.0.0.1'
    finally:
        s.close()
    return IP

def get_config_bundle():
    
    config_bundle_url = find_pattern( pattern = 'http://\*/downloads/[a-z0-9]{8}-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{12}/configBundle-.*\.tgz', search_object = log_file, is_file = True)
    if config_bundle_url == None:
        print('\033[01;31mNo configBundle url found!\n Logs are located at ' + log_file +'\033[00m')
        exit(1)
    config_bundle_url = replace_pattern(original_pattern = 'http://\*', replacement_pattern = 'https://' + esxi_hostname, replacement_object = config_bundle_url)
    print('\033[01;36mESXi Config bundle located at: ' + config_bundle_url + '\033[00m')
    with warnings.catch_warnings():
        warnings.simplefilter('ignore')
        config_bundle_response = requests.get(config_bundle_url, verify=False)
        config_bundle_download_location = '/tmp/' + config_bundle_url.split('/')[-1]
    open(config_bundle_download_location , 'wb').write(config_bundle_response.content)
    print('\033[01;36mESXi Config bundle downloaded to: ' + config_bundle_download_location + '\033[00m')

def set_prompt(ssh_connection):
    if not isinstance(ssh_connection, pexpect.spawn):
        print('This function requires an instance of pexpect.spawn!')
        exit(1)
    else:
        ssh_connection.sendline('unset prompt; unset PROMPT_COMMAND')
        ssh_connection.sendline('PS1="[PEXPECT]# "')
        match_prompt(ssh_connection)

def match_prompt(ssh_connection):
    global match_prompt_counter
    if not isinstance(ssh_connection, pexpect.spawn):
        print('This function requires an instance of pexpect.spawn!')
        exit(1)
    else:
        prompt = '\[PEXPECT\]#'
        match_prompt_counter+=1
        index = ssh_connection.expect(prompt)
        if index != 0:
            print('Prompt was not matched!\n Exiting...')
            return False
        else:
            return True

def ssh_connect(ssh_connection):
    global times_password_prompted
    times_password_prompted = 0
    if not isinstance(ssh_connection, pexpect.spawn):
        print('This function requires an instance of pexpect.spawn!')
        exit(1)
    else:
        while True:
            ssh_index = ssh_connection.expect(['[Pp]assword', 'Are you sure you want to continue connecting', pexpect.TIMEOUT, pexpect.EOF])
            if ssh_index == 0:
                times_password_prompted+=1
                if times_password_prompted > 1:
                    print('The password is not correct, please re-run this script with the proper password')
                    exit(1)
                else:
                    ssh_connection.logfile_read.close()
                    ssh_connection.logfile_read = None
                    ssh_connection.sendline(esxi_password)
                    ssh_connection.logfile_read = open(log_file, 'a')
                    time.sleep(5)
                    if ssh_connection.expect(['[Pp]assword:', '\[root', pexpect.TIMEOUT]) == 0:
                        print('\033[01;31mPassword was not correct, exiting...')
                        ssh_connection.send('\003')
                        ssh_connection.close()
                        exit(1)
                    print('\033[01;32mPassword was correct...\033[00m')
                    print('\033[01;34mSetting custom prompt...\033[00m')
                    set_prompt(ssh_connection)
                    match_prompt(ssh_connection)
                    return
            elif ssh_index == 1:
                accept_fingerprint(ssh_connection)
                continue 
            else:
                print('C\033[01;33monnection was not successfull, please check the logs @ ' + log_file + 'for more details\033[00m')
                exit(1)

##########################################################Execution Begins#################################################################################################################################

match_prompt_counter=0
times_password_prompted=0 # This value will be used to check if we immediately get re-prompted for a password for ssh/scp connections, causing the attempt to be terminated after 2 attempts to avoid locking accounts
esxi_password = getpass("\033[01;34mPlease enter the root password for the esxi host:\033[00m ")
ip_address = get_ip()
esxi_hostname = re.sub('\.[0-9]{1,3}$', '.80', ip_address)
esxi_hostname = '192.168.1.230'
esxi_username='root'
depot_file_location = '/run/media/'+ os.getlogin() + '/CSS/css/repository/applications/esxi/updates/VMware-ESXi-depot.zip' 
if not os.path.exists(depot_file_location):
    print('\033[01;31mDepot file not found! the depot file must be located at:\n /run/media/' + os.getlogin() + '/CSS/css/repository/applications/esxi/updates/VMware-ESXi-depot.zip\033[00m')
    exit(1)

depot_file_destination = '/vmfs/volumes/Datastore01/Updates/VMware-ESXi-depot.zip'
depot_file_destination = '/vmfs/volumes/datastore1/Updates/VMware-ESXi-depot.zip'
log_file='/tmp/esxi_update_log'

ssh_connection = pexpect.spawn('ssh ' + esxi_username + '@' + esxi_hostname, encoding='utf-8')
ssh_connection.logfile_read = open(log_file, 'w')
ssh_connect(ssh_connection)
check_maintenance_mode(ssh_connection)
ssh_connection.sendline('mkdir -p /vmfs/volumes/Datastore01/Updates')
match_prompt(ssh_connection)
print('\033[01;34mExecuting scp '+ depot_file_location +' ' + esxi_username + '@' + esxi_hostname + ':' + depot_file_destination + '\033[00m')
scp_file(esxi_username, esxi_hostname, esxi_password, depot_file_location, depot_file_destination)
print('\033[01;34mSyncing host configuration...\033[00m')
ssh_connection.sendline('vim-cmd hostsvc/firmware/sync_config')
match_prompt(ssh_connection)
print('\033[01;34mBacking up host configuration...\033[00m')
ssh_connection.sendline('vim-cmd hostsvc/firmware/backup_config')
match_prompt(ssh_connection)
print('\033[01;34mGetting location of config bundle...\033[00m')
get_config_bundle()
ssh_connection.sendline('esxcli software vib install --depot=' + depot_file_destination) # install the new depot file

while True:
    ssh_index = ssh_connection.expect (['[Pp]assword','Are you sure you want to continue connecting', pexpect.EOF, pexpect.TIMEOUT, 'The update completed successfully','Host is not changed', 'Could not download from depot at zip']) # start scanning stdout again, mostly listening for 'The update completed successfullly'
    if ssh_index == 0: 
        if times_password_prompted >= 2:
            print('\033[01;33mIt appears the password is incorrect, please run this script again with the correct password...\033[00m')
            exit(5)
        else:                                                                                          
            times_password_prompted+=1                                                                 
            ssh_connection.sendline(esxi_password)                                                    
    elif ssh_index == 1:                                                                               
        accept_fingerprint(ssh_connection)
    elif ssh_index == 2:                                                                               
        ssh_connection.close(force=True)                                                               
        break                                                                                          
    elif ssh_index == 3:                                                                               
        print('Connection timed out with the ssh process. Logs are located at ' +  log_file)           
        exit(3)
    elif ssh_index == 4 or ssh_index == 5:
        ssh_index = None
        ssh_index = ssh_connection.expect('PEXPECT', pexpect.EOF, pexpect.TIMEOUT, '.*') #listen to see if the prompt returns after the update completes successfully
        if match_prompt(ssh_connection):
            print('\033[01;32mThe update completed successfully, the ESXi host will need to reboot\033[00m')
            response = input('\033[01;34mPress y to reboot the host now:\033[00m ')
            if  re.match('y|Y|yes|YES', response):
                ssh_connection.sendline('reboot')
                print('\033[01;34mRebooting ESXi host...\033[00m')
                print('\033[01;34mWaiting for host to go down...\033[00m')
                while ssh_connection.isalive():
                    time.sleep(5)
                time.sleep(5)
                wait_for_reboot(esxi_hostname)
            else:
                print('\033[01;33mThe script will not reboot the host now. Ensure that you reboot the host at your earliest convenience\033[00m')
                exit(100)
        else:
            print('The update completed successfully, but the prompt was never returned')
            print('Please ssh into the host manually and ensure that it is rebooted')
            exit(8)
    elif ssh_index == 6:
        print('\033[01;33mSomething went wrong and the depot file could not be found on the esxi host!\nPlease ensure that a depot file dxists at: /vmfs/volumes/Datastore01/Updates/VMware-ESXi-depot.zip\033[00m')
        exit(1)
