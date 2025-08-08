variable "proxmox_node_name" {
  description = "The name of the proxmox node"
  default     = "proxmox"
  type        = string
}

variable "openshift_nodes" {
  description = "number of combined kubernetes control/worker nodes to create"
  default     = 3
  type        = number
}

variable "openshift_install_iso_name" {
  description = "name of the iso to install openshift from"
  default     = "openshift-install.iso"
  type        = string
}

variable "iso_url" {
  default = "http://dev01/isos/"
  type    = string
}

variable "proxmox_username" {
  description = "The username for the proxmox node"
  type        = string
  default     = "harmonj@pam"
}

variable "proxmox_password" {
  description = "The password for the proxmox node"
  type        = string
  sensitive   = true
}

variable "proxmox_endpoint" {
  description = "The endpoint for the proxmox node"
  type        = string
  default     = "https://proxmox.harmonlab.com:8006/"
}

variable "proxmox_ssh_private_key_file" {
  description = "The path to the private key file for the proxmox node"
  type        = string
  default     = "~/.ssh/id_rsa"
}


