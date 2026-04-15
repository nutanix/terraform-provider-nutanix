variable "nutanix_username" {
  type        = string
  description = "Nutanix Prism username"
}

variable "nutanix_password" {
  type        = string
  sensitive   = true
  description = "Nutanix Prism password"
}

variable "nutanix_endpoint" {
  type        = string
  description = "Nutanix Prism endpoint (IP or FQDN)"
}

variable "nutanix_port" {
  type        = string
  default     = "9440"
  description = "Nutanix Prism port"
}

# Entity datasource
variable "entity_ext_id" {
  type        = string
  description = "Ext ID of the IAM entity to fetch (e.g. from authorization policy entities)"
  default     = ""
}

# List entities datasource
variable "entities_limit" {
  type        = number
  description = "Max number of entities to return (1-100)"
  default     = 50
}

variable "entities_filter" {
  type        = string
  description = "OData filter for listing entities"
  default     = ""
}

variable "entities_order_by" {
  type        = string
  description = "OData orderby for listing entities (e.g. name asc)"
  default     = ""
}

# SAML IDP datasource (get by id)
variable "saml_idp_ext_id" {
  type        = string
  description = "Ext ID of the SAML Identity Provider to fetch"
  default     = ""
}

# SAML IDPs list datasource
variable "saml_idps_limit" {
  type        = number
  description = "Max number of SAML Identity Providers to return"
  default     = 10
}

variable "saml_idps_filter" {
  type        = string
  description = "OData filter for listing SAML Identity Providers"
  default     = ""
}

# SAML Identity Provider resource
variable "create_saml_idp" {
  type        = bool
  description = "Set to true to create a SAML Identity Provider resource"
  default     = false
}

variable "saml_idp_name" {
  type        = string
  description = "Name of the SAML Identity Provider"
  default     = "example_idp"
}

variable "saml_idp_entity_id" {
  type        = string
  description = "Entity ID for IdP metadata"
  default     = "entity_id"
}

variable "saml_idp_login_url" {
  type        = string
  description = "Login URL for IdP metadata"
  default     = "https://idp.example.com/login"
}

variable "saml_idp_logout_url" {
  type        = string
  description = "Logout URL for IdP metadata"
  default     = ""
}

variable "saml_idp_error_url" {
  type        = string
  description = "Error URL for IdP metadata"
  default     = ""
}

variable "saml_idp_certificate" {
  type        = string
  description = "Certificate for IdP metadata"
  default     = ""
}

variable "saml_idp_username_attribute" {
  type        = string
  description = "SAML assertion username attribute"
  default     = "username"
}

variable "saml_idp_email_attribute" {
  type        = string
  description = "SAML assertion email attribute"
  default     = "email"
}

variable "saml_idp_groups_attribute" {
  type        = string
  description = "SAML assertion groups attribute"
  default     = "groups"
}

variable "saml_idp_groups_delim" {
  type        = string
  description = "Delimiter for groups attribute"
  default     = ","
}

variable "saml_idp_entity_issuer" {
  type        = string
  description = "Entity issuer for SAML authnRequest"
  default     = ""
}

variable "saml_idp_is_signed_authn_req_enabled" {
  type        = bool
  description = "Whether to sign SAML authnRequests"
  default     = false
}

variable "saml_idp_custom_attributes" {
  type        = list(string)
  description = "Custom SAML attribute names"
  default     = []
}
