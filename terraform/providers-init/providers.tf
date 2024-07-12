terraform {
  required_providers {
    random = {
      source = "hashicorp/random"
      version = "3.6.2"
   }
      proxmox = {
        source = "bpg/proxmox"
        version = "0.60.1"
      }
    }
  }
