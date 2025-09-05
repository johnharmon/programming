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
    username    = split("@", var.proxmox_username)[0]
    private_key = file(pathexpand(var.proxmox_ssh_private_key_file))
  }
}

data "proxmox_virtual_environment_datastores" "proxmox" {
  node_name = "proxmox"
}

data "proxmox_virtual_environment_vms" "proxmox" {
  node_name = "proxmox"

}

locals {
  vm_pool        = [for datastore in data.proxmox_virtual_environment_datastores.proxmox.datastores : datastore if datastore.id == "vm_pool_1"]
  isos_datastore = [for datastore in data.proxmox_virtual_environment_datastores.proxmox.datastores : datastore if datastore.id == "local"][0]
}

locals {
  total_cpu_cores = 72
  total_mem_gb    = 512
}



resource "proxmox_virtual_environment_download_file" "openshift_installer_iso" {
  content_type        = "iso"
  overwrite_unmanaged = true
  datastore_id        = local.isos_datastore.id
  file_name           = var.openshift_install_iso_name
  node_name           = var.proxmox_node_name
  url                 = var.iso_url
}


resource "proxmox_virtual_environment_vm" "openshift_nodes" {
  count       = var.openshift_nodes
  name        = "ocpn-${count.index + 1}"
  description = "Openshift Node"
  tags        = ["terraform", "openshift", "worker", "master"]
  node_name   = var.proxmox_node_name

  network_device {
    enabled     = true
    mac_address = var.mac_addresses[count.index]
  }

  cdrom {
    file_id = proxmox_virtual_environment_download_file.openshift_installer_iso.id
  }

  # disk {
  #datastore_id = local.vm_pool[0].id
  #import_from  = proxmox_virtual_environment_download_file.openshift_installer_iso.id
  #interface    = "scsi0"
  #}

  disk {
    datastore_id = local.vm_pool[0].id
    interface    = "virtio0"
    size         = 101
    aio          = "native"
    serial       = "OCP-BOOT01"
  }

  disk { 
    datastore_id = local.vm_pool[0].id
    interface = "scsi0"
    size = var.openshift_storage_disk_size
    aio = "native" 
    serial = "OCP-DATA-1"
    ssd = true
  }
  boot_order = [
    "virtio0",
    "ide3"
  ]

  cpu {
    cores = floor((local.total_cpu_cores * var.openshift_cpu_multiplier) / var.openshift_nodes)
    type  = "host"
  }
  memory {
    dedicated = 128000
  }
  #  initialization {
  #    datastore_id = local.vm_pool[0]
  #    ip_config {
  #      ipv4 {
  #        address = format("192.168.86.%d", 101 + count.index)
  #      }
  #    }
}

output "datastores" {
  value = local.vm_pool
}

output "username" {
  value = split("@", var.proxmox_username)
}

output "openshift_nodes" {
  value = proxmox_virtual_environment_vm.openshift_nodes
}

output "vms" {
  value = data.proxmox_virtual_environment_vms.proxmox.vms
}
