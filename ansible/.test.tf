terraform {
  required_providers {
    proxmox = {
      source = "bpg/proxmox"
      version = ">= 0.61.0"
    }
    random = {
    source "hashicorp/random" {
      version = ">- 3.0"
    }
  }
  ansible = {
    source = "ansible/ansible"
    version = ">= 1.3.0"
  }
}
provider "proxmox" {
  insecure  = true
  endpoint = var.proxmox_endpoint
  username = var.proxmox_username
  password = var.proxmox_password
  ssh {
    agent = true
    username = split("@", var.proxmox_username)[0]
    private_key = file(var.proxmox_ssh_private_key)
  }
}

/* DATA SOURCES */

data "proxmox_virtual_environment_vms" "vm_list" {
    node_name = "proxmox"
    filter {
        name = "name" 
        regex = true 
        values = [".*(?i)freeipa.*"]
    }
}

data "proxmox_virtual_environment_datastores" "datastores" {
    node_name = "proxmox"
}



/* LOCALS */

locals {
    vm_pool = [for datastore in data.proxmox_virtual_environment_datastores.proxmox.datastore_ids : datastore if datastore == "vm_pool_1"]
}


/* LOCALS /*


/* VM RESOURCES */

resource "proxmox_virtual_environment_vm" "freeipa" {
    name = "freeipa"
    description = "Freeipa server"
    tags = ["freeipa", "terraform"]
    node_name = "proxmox"
    cdrom {
        datastore_id = "vm_pool_1"
        file_id = proxmox_virtual_environment_file.local_kube_image.id 
        interface = "ide0"
        enabled = true
    }
    disk {
        datastore_id = "vm_pool_1"
        interface = "scsi0"
        size = 20
        aio = "native"
        file_format = "raw"
    }
    cpu {
        cores = 2
        type = "x86-64-v2-AES"
    }
    memory {
        dedicated = 2048
    }
    network_device {
        bridge = "vmbr0"
        model = "virtio"
    }
    initialization {
        ip_config {
            ipv4 {
                address = "dhcp"
            }
        }
        dns {
            domain = "harmonlab.com" 
            servers = ["192.168.86.4"]
        }
        user_account {
        keys = [trimspace(file("${pathexpand("~")}/.ssh/authorized_keys"))]
        username = "cloud_user"
        password = var.proxmox_password
        }
    }

}
