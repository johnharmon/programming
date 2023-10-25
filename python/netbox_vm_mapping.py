#!/bin/python3

import requests, json, re, typing
import vms

def create_resource(api, resource_name):
    body = {
        'name': resource_name
    }
    response = requests.post(url = api.url, headers = {'Content-Type': 'application/json', 'Authorization': api.auth }, data = json.dumps(body))
    if response.status_code != 201:
        return False
    else:
        return response.json()

def generate_query_string(**kwargs):
    query_string = '?'
    for key, value in kwargs.items():
        kvp = f'{key}={value.replace(" ", "+")}&'
        query_string += kvp 
    query_string = query_string.strip('&')
    return query_string

def check_resource_exists(api, resource_name, mode = 'create'):
    response = requests.get(url = api.url)
    query = generate_query_string(name = resource_name)
    response = requests.get(url = f'{api.url}/{query}', headers = {'Accept': 'application/json', 'Authorization': api.auth})
    if response.status_code != 201:
        if mode == 'create':
            response = create_resource(api, resource_name)
            return response
        else:
            return False
    else:
        return response.json()[0]

def map_values_to_ids(vm_data, api):
    for key in vm_data.keys():

        if key == 'role':
            api.kind = 'devices/device-role'
            result = check_resource_exists(api, resource_name = vm_data[key])
            if result:
                vm_data[key] = result['id']

        elif key == 'name':
            api.kind = 'virtualization/virtual-machine'
            result = check_resource_exists(api, resource_name = vm_data[key])
            if result:
                vm_data[key] = result['id']

        elif key == 'cluster':
            api.kind = 'virtualization/clusters'
            result = check_resource_exists(api, resource_name = vm_data[key], mode = 'check')
            if result:
                vm_data[key] = result['id']
            else:
                raise ValueError(f'Resource of type {key} at endpoint {api.kind} must be created manually!')

        elif key == 'tenant':
            api.kind = 'tenancy/tenants'
            result = check_resource_exists(api, resource_name = vm_data[key], mode = 'check')
            if result:
                vm_data[key] = result['id']
            else:
                raise ValueError(f'Resource of type {key} at endpoint {api.kind} must be created manually!')
        else: 
            pass
    return vm_data
        