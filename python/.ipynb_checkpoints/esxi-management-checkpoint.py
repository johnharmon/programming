#!/bin/python3

import re
import time
import ssl
import pyVmomi
from pyVim.connect import SmartConnect 
import getpass
import os
import socket 
import sys
import datetime
import calendar 

def create_socket(set_timeout = None):
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    if set_timeout:
        sock.settimeout(int(set_timeout))
    return sock

def get_site_id():
    return re.sub(socket.gethostname(), 'wks01|wks02|esxi01|ise01|cms01|sat01', '')

def calculate_build():
    date_builds_started = datetime.date(2022, 1, 1)
    this_year = date_builds_started.today().year
    this_year_started = datetime.date(this_year, 1, 1)
    today = int(time.time())
    this_year_started_epoch = calendar.timegm(this_year_started.timetuple()) 
    seconds_diff = today - this_year_started_epoch
    num_days = int(seconds_diff/60/60/24)
    num_builds = int(num_days/60)
    if re.match('dptc|tpt1|tpt2', socket.gethostname()):
        build_number = num_builds+1
    else:
        if num_builds == 0:
            num_builds += 1
        build_number = num_builds
    build = f'{this_year}.{build_number}0'
    return build

def parse_args():
    global stage
    global build
    stage = 'Pre'
    build = calculate_build()
    for index in range(0,len(sys.argv)):
        if re.match('--stage|-s', sys.argv[index]):
            try:
                stage = sys.argv[index+1]
            except (IndexError):
                pass
            else:
                index += 1
        elif re.match('--build|-b', sys.argv[index]):
            try:
                build = sys.argv[index+1]
            except (IndexError):
                pass
            else:
                index += 1
    return

def get_subnet():
    s = create_socket(set_timeout = 1)
    try:
        s.connect(('10.254.254.254', 1))
        IP = s.getsockname()[0]
    except Exception:
        IP = '127.0.0.1'
    finally:
        s.close()
    return re.sub(IP, '\.{0-9}[1,3]$', '')

def connect():
    s = ssl.SSLContext()
    s.verify_mode = ssl.CERT_NONE
    try:
        connection = SmartConnect(host = f'{site_subnet}.80', pwd = getpass.getpass(), sslContext = s, user = 'root')
    except:
        print('Connection to the esxi server was unsuccessful. This is likely due to the host itself being down or a wrong password being given.')
    else:
        return connection

def get_vms(sc_connection):
    try:
        vm_list = sc_connection.content.rootFolder.childEntity[0].vmFolder.childEntity
    except:
        print('The conneciton was unable to access the vm folder on the ESXi host.')
        return False
    else:
        return vm_list

def get_vm_names(vm_list):
    vm_names = [vm.name for vm in vm_list]
    return vm_names

def power_off_vms(vm_list):
    failed_vms = []
    vms_have_failed = False
    for vm in vm_list:
        if vm.runtime.powerState == 'poweredOff':
            vm_list.remove(vm)
            try:
                vm.ShutdownGuest()
            except (AttributeError):
                print(f'VM: {vm.name} does not have VMware tools installed and as such is not able to perform a graceful shutdown')
                response = input('Would you like to send the poweroff command to the vm? This is not generally recommended as it can corrupt the virtual hard disk')
                if re.match('y|Y|yes|Yes|YES', response):
                    vm.PowerOff()
                else:
                    failed_vms += vm
                    vms_have_failed = True
    vms_on = True
    timeout = 600 
    time_elapsed = 0
    while vms_on:
        if time_elapsed > timeout:
            vms_still_on = [vm.name for vm in vm_list]
            raise TimeoutError('Script timed out waiting for all vms to power on. The list of vms still on is: ' + str(*vms_still_on))
        vms_on = False
        for vm in vm_list:
            if vm.runtime.powerState == 'poweredOn':
                vms_on = True
            elif vm.runtime.powerState == 'poweredOff':
                vm_list.remove(vm)
        time.sleep(5)
        time_elapsed+=5
    return vms_have_failed, failed_vms

def wait_for_dom_and_cms(cms, dom):
    timeout = 900
    time_elapsed = 0
    dom01_sock = create_socket(set_timeout = 1)
    cms01_sock = create_socket(set_timeout = 1)
    unconnected_vms = [cms, dom]
    while len(unconnected_vms) >= 1:
        for vm in unconnected_vms:
            if vm.name == f'{site_id}dom01':
                try: 
                    dom01_sock.connect((f'{site_subnet}.80', int(389)))
                except (TimeoutError):
                    dom01_sock.shutdown(socket.SHUT_RDWR)
                    dom01_sock = create_socket(set_timeout = 1)
                else:
                    unconnected_vms.remove(vm)
            elif vm.name == f'{site_id}cms01':
                try:
                    cms01_sock.connect((f'{site_subnet}.80', int(53)))
                except (TimeoutError):
                    cms01_sock.shutdown(socket.SHUT_RDWR)
                    cms01_sock = create_socket(set_timeout = 1)
                else:
                    unconnected_vms.remove(vm)
        if len(unconnected_vms) < 1:
            break
        else:
            time.sleep(5)
            time_elapsed += 5
            if time_elapsed > timeout:
                return [False, unconnected_vms]
    return True

def power_on_vms(vm_list, site_id):
    for vm in vm_list:
        if vm.runtime.powerState == 'poweredOn':
            vm_list.remove(vm)
        elif vm.name == f'{site_id}dom01':
            dom01 = vm
            vm.PowerOn()
        elif vm.name == f'{site_id}cms01':
            cms01 = vm
            vm.PowerOn()
    dom_cms_reboot = wait_for_dom_and_cms( cms = cms01, dom = dom01)
    if dom_cms_reboot[0]:
        timeout = 300
        time_elapsed = 0
        vm_list.remove(dom01)
        vm_list.remove(cms01)
        for vm in vm_list:
            vm.PowerOn()
        vms_still_booting = True
        while vms_still_booting:
            vms_still_booting = False
            for vm in vm_list:
                if re.match('cpp1|sat01', vm.name):
                    sock = create_socket(set_timeout = 1)
                    try:
                        sock.connect((vm.guest.GetIpAddress(), int(22)))
                    except (TimeoutError):
                        sock.shutdown(socket.SHUT_RDWR)
                        vms_still_booting = True
                    else:
                        vm_list.remove(vm)
                elif re.match('fwm01|wsus01|lapc', vm.name):
                    sock = create_socket(set_timeout = 1)
                    try:
                        sock.connect((vm.guest.GetIpAddress(), int(3389)))
                    except (TimeoutError):
                        sock.shutdown(socket.SHUT_RDWR)
                        vms_still_booting = True
                elif re.match('ise01|vcsa01', vm.name):
                    sock = create_socket(set_timeout = 1)
                    try:
                         sock.connect((f'{vm.guest.GetIpAddress()}', 443))
                    except (TimeoutError):
                        sock.shutdown(socket.SHUT_RDWR)
                else:
                    vm_list.remove(vm)
            time.sleep(5)
            time_elapsed += 5
            try:
                if time_elapsed > timeout:
                    raise (TimeoutError)
            except (TimeoutError):
                print('Timeout of 300 seconds reached while trying to power on the rest of the vms!')
                return False
            else:
                pass
    else:
        failed_vms = [vm.name for vm in dom_cms_reboot[1]] 
        print(str(*failed_vms) + ' failed to reboot!')
        return False, failed_vms

def take_snapshots(vm_list):
    snapshots_failed = False
    try:
        powerOff_result = power_off_vms(vm_list)
    except:
        pass
    if not powerOff_result[0]:
        failed_names = [vm.name for vm in powerOff_result[1]]
        print('Warning, These VMs failed to power off, and so snapshots will not be taken:\n' + str(*failed_names))
        [vm_list.remove(vm) for vm in powerOff_result[1]]  
    for vm in vm_list:
        failed_snapshots = []
        try:
            vm.TakeSnapshot(name = f'{stage}-{build}', memory = False, quiesce = False)
            time.sleep(0.2)
        except:
            print(f'VM: {vm.name} was unable to have a snapshot taken! please investigate further in the web interface @{site_subnet}.80')
            failed_snapshots.append(vm)
            snapshots_failed = True
    return snapshots_failed, failed_snapshots

def enter_maintenance_mode(sc_connection):
    vm_list = get_vms(sc_connection)
    try:
        power_off_vms(vm_list)
    except:
        pass
    dc = sc_connection.content.rootFoder.childEntity[0]
    hs = dc.hostFolder.childEntity[0].host[0]
    try:
        hs.EnterMaintenanceMode_Task(timeout = 10)
    except:
        print('Something went worng and the host could not be put into maintenance mode!')
        return False
    else:
        print('Host has successfully been put into maintenance mode!')
    return True

def exit_maintenance_mode(sc_connection):
    dc = sc_connection.content.rootFoder.childEntity[0]
    hs = dc.hostFolder.childEntity[0].host[0]
    try:
        hs.ExitMaintenanceMode_Task(timeout = 10)
    except:
        print('Something went worng and the host could not exit maintenance mode!')
        exit(1)
    else:
        print('Host has successfully exited maintenance mode!')
    return True

def check_container_status():
    cms_ip = f'{site_subnet}.85'
    containers = {}
    containers['splunk_deploy_server'] = 8000
    containers['splunk_enterprise_search'] = 8444
    containers['nessus_scanner'] = 8445
    containers['tenable_sc'] = 8834
    containers_good = True
    for container in containers.keys():
        container_sock = create_socket()
        try:
            container_sock.connect((cms_ip, containers[container]))
        except(TimeoutError):
            print(f'Container: {container} is not currently listening on port: {containers[container]}')
            containers_good = False
        else:
            print(f'Continer {container} is up and listening on port: {containers[container]}')
    return containers_good

global site_id        
global site_subnet
site_id = get_site_id()
site_subnet = get_subnet()