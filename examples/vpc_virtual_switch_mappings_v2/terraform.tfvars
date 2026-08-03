#define values to the variables to be used in terraform file
nutanix_username = "admin"
nutanix_password = "password"
nutanix_endpoint = "10.xx.xx.xx"
cluster_ext_id   = "00000000-0000-0000-0000-000000000000"
host_ext_id      = "00000000-0000-0000-0000-000000000001"
host_nic         = "eth0"
project_ext_id   = "00000000-0000-0000-0000-000000000002"

# Used by the `from_existing_bridge` resource in main.tf. Must name an OVS
# bridge that already exists on the host (e.g. created via
# `manage_ovs --bridge_name br1 create_single_bridge`).
existing_bridge_name = "br1"
