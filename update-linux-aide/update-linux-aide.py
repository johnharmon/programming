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
import multiprocessing as mp
import asyncio 


#ssh_connections = [sat01_ssh, wks01_ssh, wks02_ssh, cms01_ssh]

#vault_password = getpass("Please enter the vault password: ")

class Timer():
    """
    This is a simple timer class with some built in functionality to keep track of split times to be used to track how long different 
    steps in the aide update have been running. Honestly surprised python's time module didn't have a built in one. 
    """
    def __init__(self):
        self.start_time = 0
        self.time_elapsed = 0
        self.end_time = 0
        self.isrunning = False
        self.total_time = 0
        self.split_starts = []
        self.split_stops = []
        self.split_times = []
        self.total_split_time = 0

    @classmethod 
    def start(self):
        self.start_time = time.time()
        self.isrunning = True

    @classmethod 
    def get_time_elapsed(self):
        if self.isrunning == True: 
            self.time_elapsed = time.time() - self.start_time
        return f'{self.time_elapsed:.1f}'

    @classmethod 
    def stop(self):
        self.end_time = time.time()
        self.total_time += self.end_time - self.start_time
        self.start_time = 0
        self.end_time = 0
        self.time_elapsed = 0
        self.isrunning = False
        return f'{self.total_time:.1f}'

    @classmethod    
    def get_total_time(self):
        return f'{self.total_time:.1f}'

    @classmethod    
    def reset(self):
        self.total_time = 0
        self.start_time = 0
        self.end_time = 0
        self.time_elapsed = 0
        self.isrunning = False
        self.split_starts.clear()
        self.split_stops.clear()
        self.total_split_time = 0
    
    @classmethod
    def start_split(self):
        self.split_starts.append(time.time())

    @classmethod
    def get_split_time(self, split_number = -1):
        return self.split_stops[split_number] - self.split_starts[split_number]

    @classmethod
    def stop_split(self):
        self.split_stops.append(time.time())
        self.split_times.append(self.get_split_time())
        self.total_split_time += self.split_times[-1]
    
    
class aide_pxssh_wrapper_class():

    """
    This class will serve as a wrapper for pxssh processes
    that are remotely running aide updates on other hosts.
    It will need to periodically update its own attributes 
    as processes finish so they can be checked on by an outer loop.
    """
    
    def __init__(self, pxssh_connection, hostname, ip):
        if not isinstance(pxssh_connection, pxssh.pxssh):
            raise TypeError('This class must be provided with a pxssh.pxssh object to act as a wrapper for!')
        else:
            self.ip = ip
            self.hostname = hostname
            self.connection = pxssh_connection
            self.state = 'Not Started'
            self.running_task = None
            self.time_out = 3600
            self.time_taken = 0
            self.connection.timeout = 30

            self.aide_update_finished = False
            self.db_overwrite_finished = False
            self.aide_check_finished = False

            self.aide_update_failed = False
            self.db_overwrite_failed = False
            self.aide_check_failed = False

            self.aide_update_succeeded = False
            self.db_overwrite_succeeded = False
            self.aide_check_succeeded = False
            self.last_task_finished = None
            self.timer = Timer()

    @classmethod
    def print_state(self):
        print(f'{self.hostname}\n\tstate:\n\t\t{self.state}\n\t\t\tLast task completed:\n\t\t\t\t{self.last_task_finished}')
    
    @classmethod
    def sudo(self):
        self.timer.start()
        self.timer.start_split()
        self.connection.sendline('sudo -i')
        index = self.connection.expect('password for ansible:', pexpect.EOF, pexpect.TIMEOUT)
        if index == 0:
            self.connection.sendline(ansible_password)
            time.sleep(5)
            try:
                if self.connection.expect('[Ppassword]', '.*', pexpect.EOF, pexpect.TIMEOUT) == 0:
                    raise Exception(f'sudo password is incorrect for {self.hostname}')
            except (pexpect.exceptions.EOF, pexpect.exceptions.TIMEOUT):
                raise Exception(f'sudo failed for {self.hostname}!')
            else:
                self.connection.set_unique_prompt()
                self.connection.send(chr(10))
                if self.connection.prompt():
                    self.timer.stop_split()
                    return True
                else:
                    print('sudo suceeded, but the program was unable to match the shell prompt!')
                    exit(3)
        else:
            raise Exception(f'An unexpected error ocurred while trying to sudo as ansible! password prompt was never recieved!') 

    @classmethod
    def run_aide_update(self):
        self.timer.start_split()
        self.state = 'Running'
        self.running_task = 'Aide Update'
        self.connection.send(chr(10))
        self.connection.prompt()
        self.connection.sendline('aide --update')
        while self.connection.isalive():
            if self.connection.prompt(timeout = 30):
                self.aide_update_finished = True
                self.aide_update_succeeded = True
                self.last_task_finished = 'run_aide_update'
                self.timer.stop_split()
                return True
            else:
                self.time_taken +=30
                if self.time_taken >= self.time_out:
                    self.aide_update_finished = True
                    self.aide_update_failed = True
                    self.state = 'Failed'
                    return False
                self.connection.send(chr(10))
    
    @classmethod
    def overwrite_db(self):
        self.timer.start_split()
        self.running_task = 'Overwriting old aide db'
        self.connection.send(chr(10))
        self.connection.prompt()
        self.connection.sendline('/usr/bin/cp -f /var/lib/aide/aide.db.new.gz /var/lib/aide/aide.db.gz')
        if self.connection.prompt():
            self.db_overwrite_finished = True
            self.db_overwrite_succeeded = True
            self.last_task_finished = 'overwrite_db'
            self.timer.stop_split()
            return True
        else:
            self.db_overwrite_finished = True 
            self.db_overwrite_failed = True
            self.state = 'Failed'

    @classmethod
    def run_aide_check(self):
        self.timer.start_split()
        self.time_taken = 0
        self.running_task = 'Running aide check'
        self.connection.send(chr(10))
        self.connection.prompt()
        self.connection.sendline('aide --check')
        while self.connection.isalive():
            if self.connection.prompt(timeout = 30):
                self.aide_update_finished = True
                self.aide_update_succeeded = True
                self.state = 'Finished'
                self.last_task_finished = 'run_aide_check'
                self.timer.stop_split()
                return True 
            else:
                self.time_taken +=30
                if self.time_taken >= self.time_out:
                    self.aide_check_finished = True
                    self.aide_check_failed = True
                    self.state = 'Failed'
                    return False
                self.connection.send(chr(10))

    @classmethod
    def get_task_times(self):
        return self.timer.split_starts, self.timer.split_stops, self.timer.get_split_total(), self.timer.time_elapsed


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

def get_password_from_vault(this_vault_password, vault_location = '/etc/ansible/group_vars/all/vault.yml', password_to_find = 'vault_ansible'):
    ansible_vault_process = pexpect.spawn(f'ansible-vault view {vault_location}', encoding='utf-8')
    ansible_vault_process.expect('Vault password:')
    ansible_vault_process.sendline(this_vault_password)
    try:
        decryption_result = ansible_vault_process.expect([f'{password_to_find}:', 'ERROR! Decryption failed'])
        if decryption_result == 1:
            print('\033[01;31mVault Password was not correct!\033[00m')
            exit(1)
    except pexpect.exceptions.EOF:
        print('\033[01;31mThe password you are looking for is not in this vault!\033[00m')
        exit(2)
    else:
        ansible_vault_process.expect('\n')
        ansible_password = ansible_vault_process.before 
        ansible_password = ansible_password.split('\x1b')[0].strip()
        ansible_vault_process.close()
        ansible_password = re.sub('"', '', ansible_password)
        return ansible_password

def manage_aide_wrapper(aide_wrapper, child_connection):
    if not isinstance(aide_wrapper, aide_pxssh_wrapper_class):
        raise TypeError('This function must be given an aide_pxsssh_wrapper_class instance')
    else:
        child_connection.send(f'Starting SSH Login for {aide_wrapper.hostname}')
        if aide_wrapper.conneciton.login(server = aide_wrapper.ip, username = 'ansible', password = ansible_password):
            child_connection.send(f'Successfully logged into {aide_wrapper.hostname}')
            if aide_wrapper.sudo():
                while True:
                    child_connection.send(aide_wrapper.print_state())
                    if aide_wrapper.state == 'Not Started':
                        aide_wrapper.run_aide_update()
                    elif aide_wrapper.state == 'Running':
                        if aide_wrapper.aide_update_succeeded == True and aide_wrapper.last_task_finished == 'run_aide_update':
                            aide_wrapper.overwrite_db()
                        elif aide_wrapper.db_overwrite_succeeded == True and aide_wrapper.last_task_finished == 'overwrite_db':
                            aide_wrapper.run_aide_check()
                        elif aide_wrapper.aide_check_succeeded == True and aide_wrapper.last_task_finished == 'run_aide_check':
                            print(f'\033[01;32mAide checks and updates have finished on {aide_wrapper.hostname}!\033[00m')
                    time.sleep(15)
            else:
                print(f'sudo attempt for ansible failed on {aide_wrapper.hostname}')
        else:
            child_connection.connection.send(f'Unsucessful login attempt on {aide_wrapper.hostname}')
            aide_wrapper.connection.close()
            return False

def main():
    if __name__ == '__main__':
        if len(sys.argv) >= 2:
            password_name = sys.argv[1]
        else:
            password_name = 'vault_ansible'
        vault_file = '/etc/ansible/group_vars/all/vault.yml'

        mp.set_start_method('fork')

        global ansible_password 
        ansible_password = get_password_from_vault(this_vault_password = getpass("Please enter the vault password: "), vault_location = vault_file, password_to_find = password_name)

        sat01_ssh = pxssh.pxssh(encoding = 'utf-8')
        wks01_ssh = pxssh.pxssh(encoding = 'utf-8')
        wks02_ssh = pxssh.pxssh(encoding = 'utf-8')
        cms01_ssh = pxssh.pxssh(encoding = 'utf-8')

        sat01_logfile = open('/tmp/sat01_aide.log', 'w')
        wks01_logfile = open('/tmp/wks01_aide.log', 'w')
        wks02_logfile = open('/tmp/wks02_aide.log', 'w')
        cms01_logfile = open('/tmp/cms01_aide.log', 'w')

        sat01_ssh.logfile_read = sat01_logfile
        wks01_ssh.logfile_read = wks01_logfile
        wks02_ssh.logfile_read = wks02_logfile
        cms01_ssh.logfile_read = cms01_logfile

        last_octets = ['66', '65', '85', '87']
        subnet_ip = get_ip()
        subnet_ip = re.sub('[0-9]{1,3}$', '', subnet_ip)

        sat01_wrapper = aide_pxssh_wrapper_class(pxssh_connection = sat01_ssh, hostname = 'sat01', ip = f'{subnet_ip}86')
        wks01_wrapper = aide_pxssh_wrapper_class(pxssh_connection = wks01_ssh, hostname = 'wks01', ip = f'{subnet_ip}65')
        wks02_wrapper = aide_pxssh_wrapper_class(pxssh_connection = wks02_ssh, hostname = 'wks02', ip = f'{subnet_ip}66')
        cms01_wrapper = aide_pxssh_wrapper_class(pxssh_connection = cms01_ssh, hostname = 'cms01', ip = f'{subnet_ip}85')
        ssh_wrappers = (sat01_wrapper, wks01_wrapper, wks02_wrapper, cms01_wrapper)

        ssh_processes = []
        connection_states = []
        parent_pipes = []
        child_pipes = []

        for index in range(0,4):
            connection_states[index] = ssh_wrappers[index].print_state()
            parent_pipes[index], child_pipes[index] = mp.Pipe()
            ssh_processes[index] = mp.Process(target = manage_aide_wrapper, args=(ssh_wrappers[index], child_pipes[index],))
            ssh_processes[index].start()

        time.sleep(10)
        for index in range(0,4):
            print(ssh_wrappers[index].hostname)
            print(parent_pipes[index].recv())
            print(parent_pipes[index].recv())
            print('')
        
        while True:
            for index in range(0,4):
                print(ssh_wrappers[index].hostname)
                if parent_pipes[index].poll(0.1):
                    connection_states[index] = parent_pipes[index].recv()
                print(connection_states[index])
            time.sleep(15.1)
main()






#def send_line(ssh_connections = [], line_to_send):
#    for ssh_connection in ssh_connections:
#        ssh_connection.sendline(line_to_send)
#
#def expect_output(ssh_connections = [], pattern_to_expect, time_out = 30):
#    for ssh_connection in ssh_connections:
#        ssh_connection.expect(pattern_to_expect, timeout = time_out) 
#
#def maintain_connections(ssh_connections = [], time_out = 900):
#    time_waited = 0
#    while time_waited <= time_out:
#        for ssh_connection in ssh_connections:
#            ssh_connection.sendline()
#        time.sleep(30)
#        time_waited+=30
#
#def check_connections(ssh_connections = []):
#    connections_finished = 0
#    for ssh_connection in ssh_connections:
#        if ssh_connection.isalive():
#            try:
#                ssh_connection.prompt()
#            except pexpect.exceptions.TIMEOUT:
#                continue
#            else:
#                connections_finished += 1
#    return connections_finished


#print(get_password_from_vault(vault_location = vault_file, this_vault_password = vault_password))

#my_ssh = pxssh.pxssh(encoding='utf-8', logfile = open('/tmp/ssh.log', 'w'))

#password = getpass()
#my_ssh.login(username = 'jharmon', server = '127.0.0.1', ssh_key = '/home/jharmon/.ssh/id_rsa')

#my_wrapper = aide_pxssh_wrapper_class(my_ssh)
#
#
#print(my_wrapper.aide_check_failed)
#print(my_wrapper.state)
#print(my_wrapper.connection.isalive())
#my_wrapper.connection.close()