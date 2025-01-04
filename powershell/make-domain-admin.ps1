#!/bin/pwsh

Add-ADGroupMember -Identity "Enterprise Administrators" -Members "@@{ent_dod_admin.username}@@"


$netbox_plugin = 'https://docs.ansible.com/ansible/latest/collections/netbox/netbox/nb_inventory_inventory.html'