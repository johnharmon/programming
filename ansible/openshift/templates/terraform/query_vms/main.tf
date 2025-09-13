terraform {
  required_providers {

    proxmox = {
      source  = "bpg/proxmox"
      version = "0.81.0"
    }
    random = {
      source  = "hashicorp/random"
      version = ">= 3.0"
    }
    ansible = {
      source  = "ansible/ansible"
      version = ">= 1.3.0"
    }
  }
}
provider "proxmox" {
  insecure = true
  endpoint = var.proxmox_endpoint
  username = var.proxmox_username
  password = var.proxmox_password
  ssh {
    agent       = true
    username    = var.proxmox_username
    private_key = file(var.proxmox_ssh_private_key_file)
  }
}

data "proxmox_virtual_environment_vms" "proxmox" {
}

output "existing_vms" {
  value = data.proxmox_virtual_environment_vms.proxmox
}
