# define the type of variables to be used in terraform file
variable "nutanix_username" {
  description = "Username for the Nutanix Prism Central."
  type        = string
}

variable "nutanix_password" {
  description = "Password for the Nutanix Prism Central."
  type        = string
  sensitive   = true
}

variable "nutanix_endpoint" {
  description = "Prism Central endpoint (IP address or FQDN)."
  type        = string
}

variable "claim_token_name" {
  description = "Name of the claim token used to create the patched image."
  type        = string
}

variable "claim_token_expiry_time" {
  description = "Expiry time of the claim token in RFC3339 format."
  type        = string
}

variable "patched_image_name" {
  description = "Name of the patched image."
  type        = string
}

variable "host_type" {
  description = "Type of the host installed on the node (AHV or ESX)."
  type        = string
}

variable "local_host_image_ext_id" {
  description = "External ID of the uploaded host image used for patching."
  type        = string
}

variable "node_ext_id" {
  description = "External ID of the node to patch."
  type        = string
}

variable "hostname" {
  description = "Hostname of the hypervisor or host OS."
  type        = string
}

variable "node_ip" {
  description = "Management IPv4 address of the node."
  type        = string
}

variable "node_gateway" {
  description = "Management network IPv4 gateway of the node."
  type        = string
}
