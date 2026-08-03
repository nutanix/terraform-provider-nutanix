
output "NUTANIX_ENDPOINT" {
  value = local.pc.endpoint
}

output "NUTANIX_STORAGE_CONTAINER" {
  value = local.default_container_uuid
}

######## IAM OUTPUTS ########
# API key secret is only returned at creation time; a GET on the key does not
# return it. So the output must read from the resource, not a data source.
output "api_key_value" {
  description = "The generated API Key. Store this securely (only available once at create)."
  value       = try(nutanix_user_key_v2.admin_key.key_details[0].api_key_details[0].api_key, null)
  sensitive   = true
}

output "admin_sa_ext_id" {
  description = "ext_id of the Terraform Admin service-account user."
  value       = nutanix_users_v2.admin_sa.ext_id
}

output "admin_key_ext_id" {
  description = "ext_id of the generated API key."
  value       = nutanix_user_key_v2.admin_key.ext_id
}


