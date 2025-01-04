#!/bin/python3
import pexpect
import os
from getpass import getpass
import time
import sys
#from pynput import keyboard
import re 

file = open('/home/jharmon/.bashrc', 'r' )
text_to_search = file.read()

pattern=re.compile(r'cdm()')
matches = pattern.finditer(text_to_search)

for match in matches:
    print(match)

#child = pexpect.spawn('/usr/bin/zsh',  encoding='ASCII')
#child.logfile = open('/tmp/python_log', 'w')
##child.sendline('sed l')
##child.sendcontrol('c')
#child.logfile = None
#print('\033[01;03;04;05;37mIMPORTANT: Press Ctrl + ] to return control back to the python script!\033[00m')
##child.interact()
#child.close()
##quit()





