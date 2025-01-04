#!/bin/python3  
from tkinter import W
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
import shutil 


class fileActionInfo():
    def __init__(self, file_source, needs_scp = False, scp_destination = None, ownership = None, permissions = None ):
        self.source = file_source
        self.filename = re.split('/', file_source)[-1]
        self.needs_scp = needs_scp
        self.scp_destination = scp_destination
        self.ownership = ownership 
        self.permissions = permissions 
        self.absolute_path = os.path.realpath(self.source)
        if needs_scp:
            self.remote_info = re.split('@|:', scp_destination )
            self.hostname = self.remote_info[1]
            self.username = self.remote_info[0]
            self.destination = self.remote_info[0]
            self.dirpath = re.sub(self.filename, self.destination, '')
    

def accept_fingerprint(spawn, host): # Conditionally accept new server fingerprint based on user response
    if isinstance(spawn, pexpect.spawn):
        response = input(f'\033[01;34mThe host @{host}does not have a recognized fingerprint. Would you like to add the fingerprint to the known hosts [yes/no]? ')
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

def create_file_list():
    file_list = f'/run/media/{os.getlogin()}/CSS/css/repository/scripts/file_list.txt'
    with open(file_list, 'rb') as file:
        clean_version = re.sub(b'\r\n', b'\n', file.read())
        clean_version = re.sub(b'^\\[0-9]{1,3}\\[0-9]{1,3}\\[0-9]{1,3}', b'', clean_version)
    open(file_list, 'wb').write(clean_version)
    with open(file_list, 'r') as file:
        files = file.readlines()
        file_objects = []
        for line in files:
            line_info = line.split()
            source = line_info[0]
            if len(line_info) == 1:
                file_objects.append(fileActionInfo(file_source = source))
            else:
                needs_scp = True
                scp_destination = line_info[2]
                ownership = line_info[3]
                permissions = line_info[4]
                file_objects.append(fileActionInfo(file_source = source, needs_scp = needs_scp, scp_destination = scp_destination, ownership = ownership, permissions = permissions))
    return file_objects

def scp_file(fileInfo): # copy file to remote server based on local source, remote destination, server, username and hostname
    if not os.path.exists(fileInfo.source):
        print(f'\033[01;31mFile not found, the file: {fileInfo.source} must Exist!\033[00m')
    global times_password_prompted
    times_password_prompted=0
    scp_process = pexpect.spawn(f'scp {fileInfo.source} {fileInfo.username}@{fileInfo.hostname}:{fileInfo.destination}', encoding='utf-8') # use pexpect to create a child scp process that can be controlled and passing in all the necessary values
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
                scp_process.sendline(ansible_password)
                scp_process.logfile_read = open(log_file, 'a')
                scp_index = scp_process.expect([f'scp: {fileInfo.destination}: No such file or directory', f'{fileInfo.filename}', pexpect.EOF, pexpect.TIMEOUT])
                if scp_index == 0:
                    print(f'\033[01;31mUnable to copy  {fileInfo.source} to {fileInfo.hostname}, please ensure that the directory structure: {fileInfo.dirpath} exists and is writable!\033[00m')
                    exit (1)
                elif scp_index == 1:
                    print(f'\033[01;32mScp of {fileInfo.filename} to {fileInfo.hostname} was successful!\033[00m')
        elif scp_index == 1: # add esxi fingerprint to ~/.ssh/known_hosts if it isn't there already for some reason
            accept_fingerprint(scp_process)
        elif scp_index == 2:
            scp_process.close(force=True)
            break
        elif scp_index == 3:
            print('Connection timed out with the scp process. Logs are located at ' +  log_file)
            exit(1)
    scp_process.close()

def chown_and_chmod(fileInfo):
    sshc = pxssh.psxxh(encoding = 'utf-8')
    sshc.login(server = fileInfo.hostname, username = 'ansible', password = ansible_password)
    sshc.sendline('sudo -i')
    sshc.sendline(ansible_password)
    sshc.sendline(f'cd {fileInfo.dirpath}')
    sshc.sendline(f'chown {fileInfo.ownership} {fileInfo.filename}')
    sshc.sendline(f'chmod {fileInfo.permissions} {fileInfo.filename}')


def copy_local_files(file_list):
    if not os.path.exists('/tmp/css'):
        os.makedir('/tmp/css')
    for file in file_list:
        if os.path.exists(file.absolute_path):
            shutil.copy(file.absolute_path, '/tmp/css/')
        else:
            print(f'file: {file.absolute_path} does not exist!')

def set_local_permissions(dir_path = '/tmp/css', permissions = 0o755, recursive = False):
    os.chmod(dir_path, permissions) 
    for dirpath, dirnames, filenames in os.walk(dir_path):
        os.chmod(dirpath, permissions)
        for file in filenames:
            os.chmod(os.path.join(dirpath, file), permissions)

username = 'ansible'
ansible_password = getpass.gepass('Please enter the ansible password: ')