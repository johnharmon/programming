terraform {
  required_providers {
    kubernetes = {
      source = "hashicorp/kubernetes"
      version = "2.25.2"
    }
    proxmox = {
      source = "Telmate/proxmox"
      version = "3.0.1-rc1"
    }
    helm = {
      source = "hashicorp/helm"
      version = "2.12.1"
    }
  }
}

provider "proxmox" {

}

provider "kubernetes" {

}

provider "helm" {

}
