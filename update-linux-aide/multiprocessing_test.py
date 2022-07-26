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



class my_timer_class():
    def __init__(self, seconds, name):
        self.seconds = seconds
        self.name = name
        self.sleep_started = False
    def sleep(self):
        self.sleep_started = True
        time.sleep(self.seconds)
        return

my_timers = []



my_times = (1000, 1000, 1000, 1000)
my_names = ('one', 'two', 'three', 'four')
for index in range(0,4):
    my_timers.append(my_timer_class(seconds = my_times[index], name = my_names[index] ))

my_processes = []

for my_timer in my_timers:
    my_processes.append(mp.Process(target = my_timer.sleep, args=()))

for process in my_processes:
    process.start()

slept_for = 0
parent_connections = []
child_connections = []
while True:
    print(f'\033[01;32mSlept for: {slept_for} out of 1000 seconds...')
    for index in range(0,4):

        print(f'\033[01;36m{my_timers[index].name}:\033[00m')
        #print(f'\033[01;37mSleep started: {my_timers[index].sleep_started}\033[00m')
        print(f'\033[01;34m{my_processes[index].is_alive()}\033[00m')
    time.sleep(5)
    slept_for += 5


    