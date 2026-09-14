variable "nutanix_username" {
  description = "Username for the Nutanix Prism Central endpoint."
  type        = string
}

variable "nutanix_password" {
  description = "Password for the Nutanix Prism Central endpoint."
  type        = string
  sensitive   = true
}

variable "nutanix_endpoint" {
  description = "Prism Central endpoint (IP or FQDN)."
  type        = string
}

variable "node_manufacturer" {
  description = "Manufacturer of the node being registered."
  type        = string
}

variable "node_model" {
  description = "Model of the node being registered."
  type        = string
}

variable "node_serial" {
  description = "Serial number that uniquely identifies the node."
  type        = string
}
