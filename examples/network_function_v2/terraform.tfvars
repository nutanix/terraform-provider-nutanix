# Replace these values for your environment before running the example.
nutanix_username = "CHANGE_ME"
nutanix_password = "CHANGE_ME"
nutanix_endpoint = "10.XX.XX.XX"    
nutanix_port     = 9440

# Leave empty to use the first non-Prism-Central cluster returned by the API.
cluster_name = "CHANGE_ME"

# The image must already exist in Prism Central.
image_name = "CHANGE_ME"

network_function_name        = "tf-network-function-inline"
network_function_description = "Inline network function managed by Terraform"

management_subnet_name          = "tf-network-function-mgmt"
management_subnet_vlan_id       = 887
management_subnet_network       = "10.42.80.0"
management_subnet_prefix_length = 24
management_subnet_gateway       = "10.xx.xx.xx"
management_subnet_pool_start    = "10.xx.xx.xx"
management_subnet_pool_end      = "10.xx.xx.xx"

nf_vm_admin_password = "CHANGE_ME"
