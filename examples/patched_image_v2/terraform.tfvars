# define values to the variables to be used in terraform file
nutanix_username = "admin"
nutanix_password = "password"
nutanix_endpoint = "10.xx.xx.xx"

claim_token_name        = "example-patched-image-token"
claim_token_expiry_time = "2030-01-01T00:00:00Z"

patched_image_name      = "example-patched-image"
host_type               = "AHV"
local_host_image_ext_id = "00000000-0000-0000-0000-000000000000"
node_ext_id             = "00000000-0000-0000-0000-000000000001"
hostname                = "patched-host-01"
node_ip                 = "10.xx.xx.xx"
node_gateway            = "10.xx.xx.1"
