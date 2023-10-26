#!/bin/python3

import requests, json, re, typing
import vms

class netbox_api():
    def __init__(self, url, api_token):
        self.__kind = ''
        self.__root_url = url 
        self._url = self.__root_url
        self.__token = api_token 
        self.auth = f'Token {self.__token}'
        self.headers = {'Content-Type': 'application/json', 'Authorization': self.auth}
    @property
    def kind(self):
        return self.__kind 
    @kind.setter
    def kind(self, kind)
        self.__kind = kind 
        self._url = f'{self.__root_url}/{self.__kind}/'
    @property 
    def url(self):
        return self._url 
    @url.setter
    def url(self, url):
        self.__root_url = url
    def post(self, data):
        response = requests.post(url = self.url, headers = self.headers, data = json.dumps(data))
        return response
    def put(self, data):
        response = requests.put(url = self.url, headers = self.headers, data = json.dumps(data))
        return response
    def get(self):
        response = requests.get(url = self.url, headers = self.headers)
        return response

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
                raise ValueError(f'Resource of type {key} at endpoint {api.url} must be created manually!')

        elif key == 'tenant':
            api.kind = 'tenancy/tenants'
            result = check_resource_exists(api, resource_name = vm_data[key], mode = 'check')
            if result:
                vm_data[key] = result['id']
            else:
                raise ValueError(f'Resource of type {key} at endpoint {api.url} must be created manually!')
        else: 
            pass
    return vm_data
        