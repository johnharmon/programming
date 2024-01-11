#!/bin/python3
# -*- coding: utf-8 -*-
#This is a script (eventually to be a module) designed to provide a datastructure representing the available IP space of a given suite.
# It will break it out by node type (msn, user-int, user-ext, elements) creating a dictionary representing each different subnet,
# and a list of lists representing available IPs in consecutive ranges
from ansible.module_utils.basic import AnsibleModule
import requests
import json
import re 
import ipaddress
import yaml 
import sys 
import os
import argparse

def define_params():
    parser = argparse.ArgumentParser(description='This is a script (eventually to be a module) designed to provide a datastructure representing the IP space of a given suite.')
    # parser.add_argument('--inventory', type=str, mandatory=True)
    # parser.add_argument('--suite_vars', type=str, mandatory=True)
    parser.add_argument('--host_identifier', type=str, required=False, default='ansible_host:')
    # inventory = args.inventory
    # suite_vars = json.loads(args.suite_vars)
    # return {'inventory_file': inventory, 'suite_vars': suite_vars, 'host_identifier': args.host_identifier}
    inventory = '/home/jharmon/programming/python/ansible-modules/test-files/inventory.yml'
    parser.add_argument('--initial_inventory', type=bool, required=False, default=False)
    parser.add_argument('--subnets', type=str, required=False)
    parser.add_argument('--exclusions', type=str, required=False, default=None) # string representing a json list of ips to exclude, assumes /24 and is given only 4th octets, will default to bottom 10 and top 20 of each subnet
    args = parser.parse_args()
    return inventory, args, 'ansible_host:'

def parse_suite_vars(suite_vars):
    subnets: {
        'mission': {
            'msn': {
            'octet1': suite_vars['mission']['msn']['octet1'] if suite_vars['mission']['msn']['octet1'] else '10',
            'octet2': suite_vars['mission']['msn']['octet2'] if suite_vars['mission']['msn']['octet2'] else '177',
            'octet3': suite_vars['mission']['msn']['octet3'] if suite_vars['mission']['msn']['octet3'] else suite_vars['subnet'],
            'octet4': '0',
            'mask': suite_vars['mission']['msn']['mask'] if suite_vars['mission']['msn']['mask'] else '16'
            },
            'mgt': {
            'octet1': suite_vars['mission']['mgt']['octet1'] if suite_vars['mission']['mgt']['octet1'] else '192',
            'octet2': suite_vars['mission']['mgt']['octet2'] if suite_vars['mission']['mgt']['octet2'] else '168',
            'octet3': suite_vars['mission']['mgt']['octet3'] if suite_vars['mission']['mgt']['octet3'] else suite_vars['subnet'],
            'octet4': '0',
            'mask': suite_vars['mission']['mgt']['mask'] if suite_vars['mission']['mgt']['mask'] else '24'
            },
            'int': {
            'octet1': suite_vars['mission']['int']['octet1'] if suite_vars['mission']['int']['octet1'] else '192',
            'octet2': suite_vars['mission']['int']['octet2'] if suite_vars['mission']['int']['octet2'] else '168',
            'octet3': suite_vars['mission']['int']['octet3'] if suite_vars['mission']['int']['octet3'] else '1',
            'octet4': '0',
            'mask': suite_vars['mission']['int']['mask'] if suite_vars['mission']['int']['mask'] else '24'
            },
            'ext': {
            'octet1': suite_vars['mission']['ext']['octet1'] if suite_vars['mission']['ext']['octet1'] else '192',
            'octet2': suite_vars['mission']['ext']['octet2'] if suite_vars['mission']['ext']['octet2'] else '168',
            'octet3': suite_vars['mission']['ext']['octet3'] if suite_vars['mission']['ext']['octet3'] else '2',
            'octet4': '0',
            'mask': suite_vars['mission']['ext']['mask'] if suite_vars['mission']['ext']['mask'] else '24'
            }
        }
    }
    return subnets

def typecast_list(type, list):
    return [type(item) for item in list]

def build_exclusion_list(exclusions=None):
    exclusion_list = list()
    if exclusions:
        for exclusion in json.loads(exclusions):
            exclusion_list.append(exclusion)
    else:
        for i in range(11):
            exclusion_list.append(i)
        for i in range(236, 256):
            exclusion_list.append(i)
    return exclusion_list

def generate_lines(filename):
    with open(filename, 'r') as f:
        for line in f:
            yield line

def initial_inventory(subnet_list, exclusion_list):
    inventory = dict()
    for subnet in json.loads(subnet_list):
        inventory[subnet] = list()
        inventory[subnet].append(list([int(str(ip).split('.')[-1]) for ip in ipaddress.IPv4Network(subnet) if int(str(ip).split('.')[-1]) not in exclusion_list]))
        #inventory[subnet] = [item for sublist in inventory[subnet] for item in sublist if item not in exclusion_list]
        #inventory[subnet][0].sort()
    return inventory

def parse_inventory(inventory_file, host_identifier, exclusion_list):
    host_regex = re.compile(host_identifier.strip().strip("'").strip('"'))
    ip_addresses = dict()
    for line in generate_lines(inventory_file):
       rematch = re.match(host_regex, line.strip().strip('"'))
       if rematch:
            print('re was matched')
            ip = line.split(rematch.group(0))[1].strip().strip('"')
            subnet = '.'.join(ip.split('.')[:-1]) + '.0/24'
            subnet = subnet.strip('"')
            if not subnet in ip_addresses.keys():
                ip_addresses[subnet] =  set()
            ip_addresses[subnet].add(ip)
    for subnet in ip_addresses.keys():
        exclusions = typecast_list(str, exclusion_list)
        exclusions = ['.'.join(subnet.split('.')[0:-1]) + f'.{exclusion_ip}' for exclusion_ip in exclusions]
        ip_addresses[subnet].update(exclusions)
        print(ip_addresses[subnet])
        available_ips = set([str(ip) for ip in ipaddress.IPv4Network(subnet)]) - ip_addresses[subnet]
        ip_addresses[subnet] = list(available_ips)
        ip_addresses[subnet].sort(key=lambda ip: tuple(map(int, ip.split('.'))))
    return ip_addresses

# given a subnet, break it up into lists of consecutive IPs
def breakup_subnet_spaces(subnet):
    subnet_spaces = list()
    subnet_4th_octets = [int(ip.split('.')[-1]) for ip in subnet]
    new_ip = False
    subnet_index = 0
    for ip in range(1, 255):
        if ip in subnet_4th_octets:
            if not new_ip:
                print('appending item')
                subnet_spaces.append(list())
                subnet_spaces[subnet_index].append(ip)
                new_ip = True
                subnet_index += 1
            else:
                subnet_spaces[subnet_index -1].append(ip)
        else:
            new_ip = False
    return subnet_spaces

def main():
    inventory, args, host_identifier = define_params()
    ip_addresses = dict()
    exclusion_list = build_exclusion_list(args.exclusions)
    if args.initial_inventory:
        ip_addresses = initial_inventory(args.subnets, exclusion_list)
        # for ip_address in ip_addresses:
        #     ip_addresses[ip_address] = breakup_subnet_spaces(ip_addresses[ip_address])
    else:
        ip_addresses = parse_inventory(inventory, host_identifier, exclusion_list)
        for ip_address in ip_addresses.keys():
            ip_addresses[ip_address] = breakup_subnet_spaces(ip_addresses[ip_address])
    print(yaml.dump(ip_addresses))

if __name__ == '__main__':
    main()