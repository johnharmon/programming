from ansible.module_utils.basic import AnsibleModule
import requests
import json
import re 

def get_vm(module):
    url = f'{module.params["netbox_url"]}/virtualization/virtual-machines/?name={module.params["vm_body"]["name"]}'
    headers = {
        'Authorization': f'Token {module.params["netbox_token"]}',
        "Accept": "application/json" 
        }
    response = requests.get(url, headers=headers)
    return response.json()

def get_cluster(module):
    url = f'{module.params["netbox_url"]}/virtualization/clusters/?name={module.params["vm_body"]["cluster"]}'
    headers = {
        'Authorization': f'Token {module.params["netbox_token"]}',
        "Accept": "application/json" 
        }
    response = requests.get(url, headers=headers)
    return response.json()

def get_tenant(module):
    url = f'{module.params["netbox_url"]}/tenancy/tenants/?name={module.params["vm_body"]["tenant"]}'
    headers = {
        'Authorization': f'Token {module.params["netbox_token"]}',
        "Accept": "application/json" 
        }
    response = requests.get(url, headers=headers)
    return response.json()

def get_platform(module):
    url = f'{module.params["netbox_url"]}/virtualization/platforms/?name={module.params["vm_body"]["platform"]}'
    headers = {
        'Authorization': f'Token {module.params["netbox_token"]}',
        "Accept": "application/json" 
        }
    response = requests.get(url, headers=headers)
    return response.json()

def get_vrf(module):
    query_string = '&'.join([f'name={interface["vm_network"]}' for interface in module.params['interfaces']])
    query_string = '?' + query_string 
    query_string = query_string.strip('&')
    url = f'{module.params["netbox_url"]}/ipam/vrfs/{query_string}'
    headers = {
        'Authorization': f'Token {module.params["netbox_token"]}',
        "Accept": "application/json" 
        }
    response = requests.get(url, headers=headers)
    return response.json()

def get_vlan(module):
    query_string = '&'.join([f'name={interface["vm_network"]}' for interface in module.params['interfaces']])
    query_string = '?' + query_string 
    query_string = query_string.strip('&')
    url = f'{module.params["netbox_url"]}/ipam/vlans/{query_string}'
    headers = {
        'Authorization': f'Token {module.params["netbox_token"]}',
        "Accept": "application/json" 
        }
    response = requests.get(url, headers=headers)
    return response.json()

def get_device_role(module):
    query_string = f'?name={module.params["vm_body"]["device_role"]}'
    url = f'{module.params["netbox_url"]}/dcim/device-roles/{query_string}'
    headers = {
        'Authorization': f'Token {module.params["netbox_token"]}',
        "Accept": "application/json" 
        }
    response = requests.get(url, headers=headers)
    return response.json()

def get_site(module):
    url = f'{module.params["netbox_url"]}/dcim/sites/?name={module.params["vm_body"]["site"]}'
    headers = {
        'Authorization': f'Token {module.params["netbox_token"]}',
        "Accept": "application/json" 
        }
    response = requests.get(url, headers=headers)
    return response.json()

def create_device_role(module):
    url = f'{module.params["netbox_url"]}/dcim/device-roles/'
    headers = {
        'Authorization': f'Token {module.params["netbox_token"]}',
        "Accept": "application/json",
        "Content-Type": "application/json" 
        }

    device_role_body = {
        "name": module.params["vm_body"]["device_role"],
        "slug": module.params["vm_body"]["device_role"],
        "color": "ffffff",
        "vm_role": True,
        "description": None
    }
    response = requests.post(url, headers=headers, data=json.dumps(device_role_body))
    return response.json()

def create_update_vm(module):
    vm_status = get_vm(module)
    if vm_status['count'] == 0:
        method = 'POST'
    else:
        method = 'PUT'
    cluster_result = get_cluster(module)
    platform_result = get_platform(module)
    tenant_result = get_tenant(module)
    device_role = get_device_role(module)
    created_device_role = False
    if cluster_result['count'] == 0:
        module.fail_json(msg=f'Cluster {module.params["vm_body"]["cluster"]} not found')
    elif platform_result['count'] == 0:
        module.fail_json(msg=f'Platform {module.params["vm_body"]["platform"]} not found')
    elif tenant_result['count'] == 0:
        module.fail_json(msg=f'Tenant {module.params["vm_body"]["tenant"]} not found')
    elif device_role['count'] == 0:
        device_role = create_device_role(module)
        created_device_role = True
    cluster_id = cluster_result['results'][0]['id']
    plaform_id = platform_result['results'][0]['id']
    tenant_id = tenant_result['results'][0]['id']
    device_role_id =  device_role['id'] if created_device_role else device_role['results'][0]['id']
    module.params['vm_body']['cluster'] = cluster_id
    module.params['vm_body']['platform'] = plaform_id
    module.params['vm_body']['tenant'] = tenant_id
    module.params['vm_body']['device_role'] = device_role_id
    module.params['vm_body']['status'] = "active"
    url = f'{module.params["netbox_url"]}/virtualization/virtual-machines/'
    headers = {
        'Authorization': f'Token {module.params["netbox_token"]}',
        "Accept": "application/json",
        "Content-Type": "application/json" 
        }
    response = requests.request(method, url, headers=headers, data=json.dumps(module.params['vm_body']))
    module.params['vm_body']['id'] = response.json()['id']
    return module

def create_vm_interfaces(module):
    vlans = get_vlan(module)
    interface_list = []
    for interface, vlan in zip(module.params['interfaces'], vlans):
        interface['virtual_machine'] = module.params['vm_body']['id']
        interface['mode'] = 'access'
        interface['untagged_vlan'] = vlan['id']
        interface['vrf'] = get_vrf(module)['results'][0]['id']
        interface_list.append(interface)
    url = f'{module.params["netbox_url"]}/virtualization/interfaces/'
    headers = {
        'Authorization': f'Token {module.params["netbox_token"]}',
        "Accept": "application/json",
        "Content-Type": "application/json" 
        }
    response = requests.post(url, headers=headers, data=json.dumps(interface_list))
    return response.json, module

def main():
    module_args = dict(
        netbox_url=dict(type='str', required=True),
        netbox_token=dict(type='str', required=True),
        interfaces=dict(type='list', required=False),
        vm_body=dict(type='dict', required=True)
        )
    module = AnsibleModule(argument_spec=module_args, supports_check_mode=True)
    if module.params['vm_body']['name'] == None:
        module.fail_json(msg='VM name is required')
    elif module.params['vm_body']['cluster'] == None:
        module.fail_json(msg='Cluster is required')