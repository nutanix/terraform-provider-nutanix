# Service account with Super Admin access + an API key, for tests that need a
# provider api_key credential (NUTANIX_API_KEY). The api key secret is only
# returned at creation time, so it is exposed via the api_key_value output.

# 1. Create the Service Account User
resource "nutanix_users_v2" "admin_sa" {
  username     = "tf-admin-sa"
  display_name = "Terraform Admin SA"
  user_type    = "SERVICE_ACCOUNT"
}

# 2. Look up the "Super Admin" Role
data "nutanix_roles_v2" "admin_role" {
  filter = "displayName eq 'Super Admin'"
}

# 3. Create Authorization Policy (Bind User to Role)
resource "nutanix_authorization_policy_v2" "admin_policy" {
  display_name = "tf-admin-policy"
  description  = "Grants TF Super Admin access to the Service Account"

  role = data.nutanix_roles_v2.admin_role.roles[0].ext_id

  # Identity: assign the Service Account (reserved must be a JSON string)
  identities {
    reserved = jsonencode({
      user = {
        uuid = {
          anyof = [nutanix_users_v2.admin_sa.id]
        }
      }
    })
  }

  # Apply to "All" entities (Global Admin Scope; reserved must be a JSON string)
  entities {
    reserved = jsonencode({
      "*" = {
        "*" = {
          eq = "*"
        }
      }
    })
  }
}

# 4. Generate the API Key
resource "nutanix_user_key_v2" "admin_key" {
  user_ext_id = nutanix_users_v2.admin_sa.ext_id
  key_type    = "API_KEY"
  name        = "tf-admin-access-key"

  depends_on = [nutanix_authorization_policy_v2.admin_policy]
}
