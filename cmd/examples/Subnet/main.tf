provider "nutanix" { 
    username = "admin"
    password = "Nutanix/1234"
    endpoint = "10.5.80.30"
    insecure = true
}

resource "nutanix_subnet" "my-image" {
    name = "sarath_vlan0"
    vlan_id = 0 
    description = "Sarath Vlan 0"
    ip_config {
	prefix_length = 24
	default_gateway_ip = "192.168.0.1"
	pool_range = ["192.168.0.5 192.168.0.100"]
	subnet_ip = "192.168.0.0"
    }
    dhcp_options {
	boot_file_name = "bootfile"
	dhcp_server_address_host = "192.168.0.251"
	domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
	domain_search_list = ["nutanix.com", "calm.io"]
	tftp_server_name = "192.168.0.252"
	domain_name = "nutanix"
    }
}

