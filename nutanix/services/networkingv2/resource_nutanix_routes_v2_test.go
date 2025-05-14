package networkingv2_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccV2NutanixRoutesResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("terraform-test-route-%d", r)
	desc := "test terraform route description"

	resourceSubnet := "nutanix_subnet_v2.test"
	resourceVpc1 := "nutanix_vpc_v2.test-1"
	resourceVpc2 := "nutanix_vpc_v2.test-2"
	resourceRouteTable1 := "data.nutanix_route_tables_v2.rt_vpc1"
	resourceRouteTable2 := "data.nutanix_route_tables_v2.rt_vpc2"
	resourceRoute1 := "nutanix_routes_v2.test-1"
	resourceRoute2 := "nutanix_routes_v2.test-2"

	//goland:noinspection GoDeprecation
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// Create subnet
			{
				Config: testRouteSubnetConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceSubnet, "name", fmt.Sprintf("terraform_test_route_subnet_%d", r)),
					resource.TestCheckResourceAttr(resourceSubnet, "description", "terraform test subnet to test create route"),
					resource.TestCheckResourceAttr(resourceSubnet, "subnet_type", "VLAN"),
					resource.TestCheckResourceAttr(resourceSubnet, "network_id", strconv.Itoa(testVars.Networking.Subnets.VlanID)),
					resource.TestCheckResourceAttr(resourceSubnet, "is_external", "true"),
					resource.TestCheckResourceAttr(resourceSubnet, "ip_config.0.ipv4.0.ip_subnet.0.ip.0.value", testVars.Networking.Subnets.NetworkIP),
					resource.TestCheckResourceAttr(resourceSubnet, "ip_config.0.ipv4.0.ip_subnet.0.prefix_length", strconv.Itoa(testVars.Networking.Subnets.NetworkPrefix)),
					resource.TestCheckResourceAttr(resourceSubnet, "ip_config.0.ipv4.0.default_gateway_ip.0.value", testVars.Networking.Subnets.GatewayIP),
					resource.TestCheckResourceAttr(resourceSubnet, "ip_config.0.ipv4.0.pool_list.0.start_ip.0.value", testVars.Networking.Subnets.DHCP.StartIP),
					resource.TestCheckResourceAttr(resourceSubnet, "ip_config.0.ipv4.0.pool_list.0.end_ip.0.value", testVars.Networking.Subnets.DHCP.EndIP),
				),
			},
			// Create VPC 1
			{
				Config: testRouteVpc1Config(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVpc1, "name", fmt.Sprintf("terraform_test_vpc_%d", r)),
					resource.TestCheckResourceAttr(resourceVpc1, "description", "terraform test vpc 1 to test create route"),
					resource.TestCheckResourceAttrSet(resourceVpc1, "external_subnets.0.subnet_reference"),
				),
			},
			// Create VPC 2
			{
				Config: testRouteVpc2Config(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceVpc2, "name", fmt.Sprintf("terraform_test_vpc_%d", r)),
					resource.TestCheckResourceAttr(resourceVpc2, "description", "terraform test vpc 2 to test create route"),
					resource.TestCheckResourceAttrSet(resourceVpc2, "external_subnets.0.subnet_reference"),
				),
			},
			// Get route table info for VPC 1
			{
				Config: testRouteTableInfoVpc1Config(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRouteTable1, "route_tables.#", "1"),
				),
			},
			// Get route table info for VPC 2
			{
				Config: testRouteTableInfoVpc2Config(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRouteTable2, "route_tables.#", "1"),
				),
			},
			// Create route 1
			{
				Config: testRoute1Config(name, desc, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRoute1, "name", name),
					resource.TestCheckResourceAttr(resourceRoute1, "description", desc),
					resource.TestCheckResourceAttrSet(resourceRoute1, "vpc_reference"),
					resource.TestCheckResourceAttrSet(resourceRoute1, "route_table_ext_id"),
					resource.TestCheckResourceAttr(resourceRoute1, "destination.0.ipv4.0.ip.0.value", "10.0.0.2"),
					resource.TestCheckResourceAttr(resourceRoute1, "destination.0.ipv4.0.prefix_length", "32"),
					resource.TestCheckResourceAttr(resourceRoute1, "next_hop.0.next_hop_type", "EXTERNAL_SUBNET"),
					resource.TestCheckResourceAttrSet(resourceRoute1, "next_hop.0.next_hop_reference"),
					resource.TestCheckResourceAttrSet(resourceRoute1, "metadata.0.owner_reference_id"),
					resource.TestCheckResourceAttr(resourceRoute1, "metadata.0.project_reference_id", testVars.Networking.Subnets.ProjectID),
					resource.TestCheckResourceAttr(resourceRoute1, "route_type", "STATIC"),
				),
			},
			// Create route 2
			{
				Config: testRoute2Config(name, desc, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRoute2, "name", name),
					resource.TestCheckResourceAttr(resourceRoute2, "description", desc),
					resource.TestCheckResourceAttrSet(resourceRoute2, "vpc_reference"),
					resource.TestCheckResourceAttrSet(resourceRoute2, "route_table_ext_id"),
					resource.TestCheckResourceAttr(resourceRoute2, "destination.0.ipv4.0.ip.0.value", "10.0.0.3"),
					resource.TestCheckResourceAttr(resourceRoute2, "destination.0.ipv4.0.prefix_length", "32"),
					resource.TestCheckResourceAttr(resourceRoute2, "next_hop.0.next_hop_type", "EXTERNAL_SUBNET"),
					resource.TestCheckResourceAttrSet(resourceRoute2, "next_hop.0.next_hop_reference"),
					resource.TestCheckResourceAttrSet(resourceRoute2, "metadata.0.owner_reference_id"),
					resource.TestCheckResourceAttr(resourceRoute2, "metadata.0.project_reference_id", testVars.Networking.Subnets.ProjectID),
					resource.TestCheckResourceAttr(resourceRoute2, "route_type", "STATIC"),
				),
			},
			// Update route 1
			{
				Config: testRoute1UpdateConfig(name, desc, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRoute1, "name", name+"_updated"),
					resource.TestCheckResourceAttr(resourceRoute1, "description", desc+"_updated"),
					resource.TestCheckResourceAttrSet(resourceRoute1, "vpc_reference"),
					resource.TestCheckResourceAttrSet(resourceRoute1, "route_table_ext_id"),
					resource.TestCheckResourceAttr(resourceRoute1, "destination.0.ipv4.0.ip.0.value", "10.0.0.4"),
					resource.TestCheckResourceAttr(resourceRoute1, "destination.0.ipv4.0.prefix_length", "32"),
					resource.TestCheckResourceAttr(resourceRoute1, "next_hop.0.next_hop_type", "EXTERNAL_SUBNET"),
					resource.TestCheckResourceAttrSet(resourceRoute1, "next_hop.0.next_hop_reference"),
					resource.TestCheckResourceAttrSet(resourceRoute1, "metadata.0.owner_reference_id"),
					resource.TestCheckResourceAttr(resourceRoute1, "metadata.0.project_reference_id", testVars.Networking.Subnets.ProjectID),
					resource.TestCheckResourceAttr(resourceRoute1, "route_type", "STATIC"),
				),
			},
			// Update route 2
			{
				Config: testRoute2UpdateConfig(name, desc, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRoute2, "name", name+"_updated"),
					resource.TestCheckResourceAttr(resourceRoute2, "description", desc+"_updated"),
					resource.TestCheckResourceAttrSet(resourceRoute2, "vpc_reference"),
					resource.TestCheckResourceAttrSet(resourceRoute2, "route_table_ext_id"),
					resource.TestCheckResourceAttr(resourceRoute2, "destination.0.ipv4.0.ip.0.value", "10.0.0.5"),
					resource.TestCheckResourceAttr(resourceRoute2, "destination.0.ipv4.0.prefix_length", "32"),
					resource.TestCheckResourceAttr(resourceRoute2, "next_hop.0.next_hop_type", "EXTERNAL_SUBNET"),
					resource.TestCheckResourceAttrSet(resourceRoute2, "next_hop.0.next_hop_reference"),
					resource.TestCheckResourceAttrSet(resourceRoute2, "metadata.0.owner_reference_id"),
					resource.TestCheckResourceAttr(resourceRoute2, "metadata.0.project_reference_id", testVars.Networking.Subnets.ProjectID),
					resource.TestCheckResourceAttr(resourceRoute2, "route_type", "STATIC"),
				),
			},
		},
	})
}

func testRouteSubnetConfig(r int) string {
	return fmt.Sprintf(`
	data "nutanix_clusters_v2" "clusters" {}

	locals {
	  config  = (jsondecode(file("%[1]s")))
	  subnets = local.config.networking.subnets
	  cluster1 = [
		for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    		cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
	  ][0]
	}

	resource "nutanix_subnet_v2" "test" {
	  name              = "terraform_test_route_subnet_%[2]d"
	  description       = "terraform test subnet to test create route"
	  cluster_reference = local.cluster1
	  subnet_type       = "VLAN"
	  network_id        = local.subnets.vlan_id
	  is_external       = true
	  ip_config {
		ipv4 {
		  ip_subnet {
			ip {
			  value = local.subnets.network_ip
			}
			prefix_length = local.subnets.network_prefix
		  }
		  default_gateway_ip {
			value = local.subnets.gateway_ip
		  }
		  pool_list {
			start_ip {
			  value = local.subnets.dhcp.start_ip
			}
			end_ip {
			  value = local.subnets.dhcp.end_ip
			}
		  }
		}
	  }
	}

`, filepath, r)
}

func testRouteVpc1Config(r int) string {
	return testRouteSubnetConfig(r) + fmt.Sprintf(`

	resource "nutanix_vpc_v2" "test-1" {
	  name        = "terraform_test_vpc_%[1]d"
	  description = "terraform test vpc 1 to test create route"
	  external_subnets {
		subnet_reference = nutanix_subnet_v2.test.id
	  }
	  depends_on = [nutanix_subnet_v2.test]
	}

`, r)
}

func testRouteVpc2Config(r int) string {
	return testRouteSubnetConfig(r) + fmt.Sprintf(`
		resource "nutanix_vpc_v2" "test-2" {
		  name        = "terraform_test_vpc_%[1]d"
		  description = "terraform test vpc 2 to test create route"
		  external_subnets {
			subnet_reference = nutanix_subnet_v2.test.id
		  }
		  depends_on = [nutanix_subnet_v2.test]
		}
	`, r)
}

func testRouteTableInfoVpc1Config(r int) string {
	return testRouteVpc1Config(r) + `
		data "nutanix_route_tables_v2" "rt_vpc1" {
		  filter     = "vpcReference eq '${nutanix_vpc_v2.test-1.id}'"
  		  depends_on = [nutanix_vpc_v2.test-1]
		}
	`
}

func testRouteTableInfoVpc2Config(r int) string {
	return testRouteVpc2Config(r) + `
		data "nutanix_route_tables_v2" "rt_vpc2" {
		  filter = "vpcReference eq '${nutanix_vpc_v2.test-2.id}'"
		  depends_on = [nutanix_vpc_v2.test-2]
		}
	`
}

func testRoute1Config(name, desc string, r int) string {
	return testRouteTableInfoVpc1Config(r) + fmt.Sprintf(`

	resource "nutanix_routes_v2" "test-1" {
	  name               = "%[1]s"
	  description        = "%[2]s"
	  vpc_reference      = nutanix_vpc_v2.test-1.id
	  route_table_ext_id = data.nutanix_route_tables_v2.rt_vpc1.route_tables[0].ext_id
	  destination {
		ipv4 {
		  ip {
			value = "10.0.0.2"
		  }
		  prefix_length = 32
		}
	  }
	  next_hop {
		next_hop_type      = "EXTERNAL_SUBNET"
		next_hop_reference = nutanix_subnet_v2.test.id
	  }
	  metadata {
		owner_reference_id   = nutanix_vpc_v2.test-1.id
		project_reference_id = local.subnets.project_id
	  }
	  route_type = "STATIC"
	}
	`, name, desc)
}

func testRoute2Config(name, desc string, r int) string {
	return testRouteTableInfoVpc2Config(r) + fmt.Sprintf(`
	resource "nutanix_routes_v2" "test-2" {
	  name               = "%[1]s"
	  description        = "%[2]s"
	  vpc_reference      = nutanix_vpc_v2.test-2.id
	  route_table_ext_id = data.nutanix_route_tables_v2.rt_vpc2.route_tables[0].ext_id
	  destination {
		ipv4 {
		  ip {
			value = "10.0.0.3"
		  }
		  prefix_length = 32
		}
	  }
	  next_hop {
		next_hop_type      = "EXTERNAL_SUBNET"
		next_hop_reference = nutanix_subnet_v2.test.id
	  }
	  metadata {
		owner_reference_id   = nutanix_vpc_v2.test-2.id
		project_reference_id = local.subnets.project_id
	  }
	  route_type = "STATIC"
	}
	`, name, desc)
}

func testRoute1UpdateConfig(name, desc string, r int) string {
	return testRouteTableInfoVpc1Config(r) + fmt.Sprintf(`
	resource "nutanix_routes_v2" "test-1" {
	  name               = "%[1]s_updated"
	  description        = "%[2]s_updated"
	  vpc_reference      = nutanix_vpc_v2.test-1.id
	  route_table_ext_id = data.nutanix_route_tables_v2.rt_vpc1.route_tables[0].ext_id
	  destination {
		ipv4 {
		  ip {
			value = "10.0.0.4"
		  }
		  prefix_length = 32
		}
	  }
	  next_hop {
		next_hop_type      = "EXTERNAL_SUBNET"
		next_hop_reference = nutanix_subnet_v2.test.id
	  }
	  metadata {
		owner_reference_id   = nutanix_vpc_v2.test-1.id
		project_reference_id = local.subnets.project_id
	  }
	  route_type = "STATIC"
	}
	`, name, desc)
}

func testRoute2UpdateConfig(name, desc string, r int) string {
	return testRouteTableInfoVpc2Config(r) + fmt.Sprintf(`
	resource "nutanix_routes_v2" "test-2" {
	  name               = "%[1]s_updated"
	  description        = "%[2]s_updated"
	  vpc_reference      = nutanix_vpc_v2.test-2.id
	  route_table_ext_id = data.nutanix_route_tables_v2.rt_vpc2.route_tables[0].ext_id
	  destination {
		ipv4 {
		  ip {
			value = "10.0.0.5"
		  }
		  prefix_length = 32
		}
	  }
	  next_hop {
		next_hop_type      = "EXTERNAL_SUBNET"
		next_hop_reference = nutanix_subnet_v2.test.id
	  }
	  metadata {
		owner_reference_id   = nutanix_vpc_v2.test-2.id
		project_reference_id = local.subnets.project_id
	  }
	  route_type = "STATIC"
	}
	`, name, desc)
}
