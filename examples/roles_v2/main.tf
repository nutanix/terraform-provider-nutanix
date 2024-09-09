provider "nutanix" {
  username     = var.user
  password     = var.password
  endpoint     = var.endpoint
  insecure     = var.insecure
  port         = var.port
  wait_timeout = 60
}

# Create role
resource "nutanix_roles_v2" "test" {
  display_name = "test_role"
  description  = "creat a test role using terraform"
  operations = var.operations
}

# list Roles
data "nutanix_roles_v2" "test"{}

# get a specific role by id
data "nutanix_role_v2" "test" {
  ext_id = resource.nutanix_roles_v2.test.id
}

