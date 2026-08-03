# External VLAN-0 subnet on the bare deploy PE that the prismv2 deploy-PC test
# uses as its deploy network. Created here (once, on the deploy PE) rather than
# in the test: the deployed PC permanently holds an IP on this subnet, so a
# test-managed subnet would fail to delete at teardown. The deploy test looks it
# up by name via a data source. Mirrors the bare VLAN-0 subnet the test used to
# create inline (no ip_config / cluster_reference).
resource "nutanix_subnet_v2" "deploy-pc-external-subnet" {
  provider    = nutanix.deploy_pe
  name        = local.deploy_pc.subnet_name
  description = "subnet VLAN 0 for deploy PC test managed by Terraform"
  subnet_type = "VLAN"
  network_id  = 0
}

resource "nutanix_subnet_v2" "vlan-255" {
  name              = local.gc_subnet.name
  description       = "subnet VLAN 255 managed by Terraform"
  cluster_reference = local.cluster_ext_id
  subnet_type       = "VLAN"
  network_id        = local.gc_subnet.vlan_id
  ip_config {
    ipv4 {
      ip_subnet {
        ip {
          value = local.gc_subnet.network_ip
        }
        prefix_length = local.gc_subnet.prefix_length
      }
      default_gateway_ip {
        value = local.gc_subnet.gateway_ip
      }
      pool_list {
        start_ip {
          value = local.gc_subnet.start_ip
        }
        end_ip {
          value = local.gc_subnet.end_ip
        }
      }
    }
  }
}




resource "nutanix_subnet_v2" "external-nat-subnet" {
  name              = local.external_nat_subnet.name
  description       = "terraform test integration_test_Ext-Nat1"
  cluster_reference = local.cluster_ext_id
  subnet_type       = "VLAN"
  network_id        = local.external_nat_subnet.network_id
  is_external       = true
  is_nat_enabled    = true
  ip_config {
    ipv4 {
      ip_subnet {
        ip {
          value = local.external_nat_subnet.network_ip
        }
        prefix_length = local.external_nat_subnet.prefix_length
      }
      default_gateway_ip {
        value = local.external_nat_subnet.default_gateway_ip
      }
      pool_list {
        start_ip {
          value = local.external_nat_subnet.start_ip
        }
        end_ip {
          value = local.external_nat_subnet.end_ip
        }
      }
    }
  }
  depends_on = [data.nutanix_clusters_v2.clusters]
}

resource "nutanix_vpc_v2" "integration-test-vpc" {
  name        = "tf-integration_test_vpc"
  description = "terraform test integration_test_vpc"
  external_subnets {
    subnet_reference = nutanix_subnet_v2.external-nat-subnet.id
  }
  depends_on = [nutanix_subnet_v2.external-nat-subnet]
}



resource "nutanix_subnet_v2" "overlay-subnet" {
  name        = local.overlay_subnet.name
  subnet_type = "OVERLAY"
  ip_config {
    ipv4 {
      ip_subnet {
        ip {
          value         = local.overlay_subnet.network_ip
          prefix_length = 32
        }
        prefix_length = local.overlay_subnet.prefix_length
      }
      default_gateway_ip {
        value         = local.overlay_subnet.default_gateway_ip
        prefix_length = 32
      }
    }
  }
  vpc_reference = nutanix_vpc_v2.integration-test-vpc.id
  depends_on    = [nutanix_vpc_v2.integration-test-vpc]
}