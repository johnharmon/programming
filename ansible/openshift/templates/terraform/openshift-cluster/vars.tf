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


variable "mac_addresses" {
  description = "A list of mac addresses to be used as the primary interface address and node identifier"
  type        = list(string)
  validation {
    condition     = length(var.mac_addresses) > 0
    error_message = "You must provide at least one mac address for the rendezvous host"
  }
}

variable "openshift_install_iso_name" {
  description = "name of the iso to install openshift from"
  default     = "agent.x86_64.iso"
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

# variable "ansible_vault_password" {
#     description = "The password for the ansible vault"
#     type = string
#     sensitive = true
# }

