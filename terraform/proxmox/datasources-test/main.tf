terraform {

    required_providers {

        proxmox = {
            source = "bpg/proxmox"
            version = "0.60.1"
        }
        random = {
            source = "hashicorp/random" 
            version = ">= 3.0"
        }
        ansible = {
            source = "ansible/ansible"
            version = ">= 1.3.0"
        }
    }
}
  provider "proxmox" {
    insecure  = true
    endpoint = var.proxmox_endpoint
    username = "${var.proxmox_username}@pam"
    password = var.proxmox_password
    ssh {
      agent = true
      username = split("@", var.proxmox_username)[0]
      private_key = file(var.proxmox_ssh_private_key_file)
    }
  }
  
data "proxmox_virtual_environment_datastores" "proxmox" {
    node_name = "proxmox"
}


locals {
    vm_pool = [for datastore in data.proxmox_virtual_environment_datastores.proxmox.datastore_ids : datastore if datastore == "vm_pool_1"]
}

# resource "proxmox_virtual_environment_file" "iso_file" {
#     content_type = "iso"
#     node_name = "proxmox"
#     datastore_id = "local"
#     source_file {
#         path = "Rocky-9.3-x86_64-minimal.iso"
#     }

# }


resource "proxmox_virtual_environment_vm" "terraform-vm" {
    name = "terraform-vm"
    description = "Terraform VM"
    tags = ["terraform"]
    node_name = "proxmox"
    disk {
        datastore_id = local.vm_pool[0]
        # file_id = "local:iso/Rocky-9.3-x86_64-minimal.iso"
        file_format = "raw"
        interface = "scsi0"
        size = 40
        iothread = true
        aio = "native" 

    }

    cdrom {
        enabled = true
        file_id = "local:iso/Rocky-9.3-x86_64-minimal.iso"
        interface = "ide0"
    }
    cpu {
        cores = 2
        type = "x86-64-v2-AES"
    }
    memory {
        dedicated = 2048
    }
    # initialization {
    #     datastore_id = local.vm_pool[0]
    #     ip_config {
    #         ipv4 {
    #             address = "dhcp"
    #         }
    #     }
    #     user_account {
    #         keys = [trimspace(file("${pathexpand("~")}/.ssh/authorized_keys"))]
    #         username = "cloud_user"
    #         password = var.proxmox_password
    #     }
    # }
}

output "datastores" {
    value = local.vm_pool
}

output "username" {
    value = split("@", var.proxmox_username)
}