#!/bin/python3
import re
import json
import yaml 
import os
import sys 
import configparser 
import typing 
import typing_extensions
import jinja2 

def get_var_dict_key(group: str = None):
    # Stands for mission_groups
    mg = ['msn_high', 'msn_low', 'ele', 'usr']

    var_group_keys = {
        mg[0]: 'msn-v1',
        mg[1]: 'msn-low-v1',
        mg[2]: 'element',
        mg[3]: 'user'
    }

    for idx in range(0,len(mg)):
        if re.search(var_group_keys[mg[idx]], group):
            return var_group_keys[mg[idx]]
    raise NotImplementedError(f'The group: {group} not supported by this script')

def pull_hosts_from_vars(var_group: str = None, vars: yaml = None):
    var_group_dict = vars[var_group]
    return var_group_dict

def get_inventory_differential(inventory: configparser = None, group: str = None, vars_inventory: dict = None):

    existing_inventory = inventory.get(f'{group}')

