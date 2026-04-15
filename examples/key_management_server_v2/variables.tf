#define the type of variables to be used in terraform file
variable "nutanix_username" {
  type = string
}
variable "nutanix_password" {
  type = string 
}
variable "nutanix_endpoint" {
  type = string
}
variable "nutanix_port" {
  type = string
}


variable "access_information" {
  type = object({
    endpoint_url           = string
    key_id                 = string
    tenant_id              = string
    client_id              = string
    client_secret          = string
    credential_expiry_date = string
  })
}
