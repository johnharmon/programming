variable "proxmox_node_name"  {
  description = "The name of the proxmox node"
  default = "proxmox"
  type = string 
}

variable "proxmox_username" {
  description = "The username for the proxmox node"
  type = string
  default = "harmonj" 
}

variable "proxmox_password" {
  description = "The password for the proxmox node"
  type = string
  sensitive = true
}

variable "proxmox_endpoint" {
  description = "The endpoint for the proxmox node"
  type = string
  default = "https://proxmox.harmonlab.com:8006/"
}

variable "proxmox_ssh_private_key_file" {
    description = "The path to the private key file for the proxmox node"
    type = string
    default = "~/.ssh/id_rsa"
}

# variable "ansible_vault_password" {
#     description = "The password for the ansible vault"
#     type = string
#     sensitive = true
# }