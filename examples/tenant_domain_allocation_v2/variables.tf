#define the type of variables to be used in terraform file
variable "nutanix_username" {
  type        = string
  description = "The username for the Nutanix Prism Central."
}

variable "nutanix_password" {
  type        = string
  description = "The password for the Nutanix Prism Central."
  sensitive   = true
}

variable "nutanix_endpoint" {
  type        = string
  description = "The endpoint (IP/FQDN) for the Nutanix Prism Central."
}

variable "tenant_name" {
  type        = string
  description = "The name of the prerequisite tenant."
}

variable "tenant_admin_password" {
  type        = string
  description = "The password for the tenant local admin user."
  sensitive   = true
}

variable "domain_ext_id" {
  type        = string
  description = "External identifier of the domain to allocate to the tenant."
}

variable "cluster_uuid" {
  type        = string
  description = "The UUID of the cluster used as a placement target."
}

variable "storage_container_uuid" {
  type        = string
  description = "The UUID of the storage container used as a placement target."
}
