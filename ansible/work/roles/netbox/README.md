Role Name
=========

This role is designed to be given a series of nutanix vms and translate them to netbox, it also relies on some of the host_vars for that vm.

Requirements
------------

Any pre-requisites that may not be covered by Ansible itself or the role should be mentioned here. For instance, if the role uses the EC2 module, it may be a good idea to mention in this section that the boto package is required.

Role Variables
--------------

This role relies on you haveing a nutanix vm response body called nutanix_vm tied to each host

Host/Group/Inventory vars used
c2bmc_mgt_if from hostvars/inventory_hostname/ip_vars.yml
Netbox_Tenant from inventory/name.yml


Dependencies
------------

A list of other roles hosted on Galaxy should go here, plus any details in regards to parameters that may need to be set for other roles, or variables that are used from other roles.

Example Playbook
----------------

Including an example of how to use your role (for instance, with variables passed in as parameters) is always nice for users too:

    - hosts: servers
      roles:
         - { role: username.rolename, x: 42 }

License
-------

BSD

Author Information
------------------

John Harmon - Integration cube, if I've left the team and no one else has looked into how this role works, goood luck lmao
