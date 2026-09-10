package networkingv2_test

import (
	"fmt"
	"regexp"
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
					resource.TestCheckResourceAttr(resourceSubnet, "project_ext_id", "00000000-0000-0000-0000-000000000000"), // default project
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
					resource.TestCheckResourceAttr(resourceRoute1, "metadata.0.project_reference_id", "00000000-0000-0000-0000-000000000000"),
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
					resource.TestCheckResourceAttr(resourceRoute2, "metadata.0.project_reference_id", "00000000-0000-0000-0000-000000000000"),
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
					resource.TestCheckResourceAttr(resourceRoute1, "metadata.0.project_reference_id", "00000000-0000-0000-0000-000000000000"),
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
					resource.TestCheckResourceAttr(resourceRoute2, "metadata.0.project_reference_id", "00000000-0000-0000-0000-000000000000"),
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
	  name                 = "terraform_test_route_subnet_%[2]d"
	  description          = "terraform test subnet to test create route"
	  cluster_reference    = local.cluster1
	  subnet_type          = "VLAN"
	  network_id           = local.subnets.vlan_id
	  is_external          = true
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
	  lifecycle {
		ignore_changes = [links, snat_ips]
	  }
	}

`, r)
}

func TestAccV2NutanixRoutesResource_ProjectAssociation(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-route-projassoc-%d", r)
	desc := "route project association test"
	projectName := fmt.Sprintf("tf-rt-pa-proj-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRouteProjectAssociationConfig(name, desc, r, projectName, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("nutanix_routes_v2.test-1", "project_ext_id", "nutanix_project_v2.test", "ext_id"),
					resource.TestCheckResourceAttrPair("data.nutanix_route_v2.test", "project_ext_id", "nutanix_project_v2.test", "ext_id"),
					resource.TestCheckResourceAttr("data.nutanix_routes_v2.test", "routes.#", "1"),
					resource.TestCheckResourceAttrPair("data.nutanix_routes_v2.test", "routes.0.ext_id", "nutanix_routes_v2.test-1", "ext_id"),
					resource.TestCheckResourceAttrPair("data.nutanix_routes_v2.test", "routes.0.project_ext_id", "nutanix_project_v2.test", "ext_id"),
				),
			},
			{
				Config:      testRouteProjectAssociationConfig(name, desc, r, projectName, "00000000-0000-0000-0000-000000000000"),
				ExpectError: regexp.MustCompile("Update of project_ext_id is not supported"),
			},
		},
	})
}

func routeProjectExtIDLine(override string) string {
	if override == "" {
		return `project_ext_id = nutanix_project_v2.test.ext_id`
	}
	return fmt.Sprintf(`project_ext_id = "%s"`, override)
}

func testRouteProjectAssociationConfig(name, desc string, r int, projectName, projectExtIDOverride string) string {
	return fmt.Sprintf(`
	data "nutanix_clusters_v2" "clusters" {}

	locals {
		config  = (jsondecode(file("%[5]s")))
		subnets = local.config.networking.subnets
		cluster1 = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		][0]
	}

	resource "nutanix_project_v2" "test" {
		name        = "%[4]s"
		project_id  = "%[4]s"
		description = "project association test"
	}

	resource "nutanix_subnet_v2" "test" {
		name              = "tf-route-pa-subnet-%[3]d"
		description       = "subnet for route project association"
		cluster_reference = local.cluster1
		subnet_type       = "VLAN"
		network_id        = local.subnets.vlan_id
		is_external       = true
		shared_with_projects = [nutanix_project_v2.test.ext_id]
		ip_config {
			ipv4 {
				ip_subnet {
					ip { value = local.subnets.network_ip }
					prefix_length = local.subnets.network_prefix
				}
				default_gateway_ip { value = local.subnets.gateway_ip }
				pool_list {
					start_ip { value = local.subnets.dhcp.start_ip }
					end_ip   { value = local.subnets.dhcp.end_ip }
				}
			}
		}
	}

	resource "nutanix_vpc_v2" "test-1" {
		name        = "tf-route-pa-vpc-%[3]d"
		description = "vpc for route project association"
		external_subnets {
			subnet_reference = nutanix_subnet_v2.test.id
		}
		project_ext_id       = nutanix_project_v2.test.ext_id
		shared_with_projects = [nutanix_project_v2.test.ext_id]
		depends_on           = [nutanix_subnet_v2.test]
		lifecycle {
			ignore_changes = [links, snat_ips, shared_with_projects]
		}
	}

	data "nutanix_route_tables_v2" "rt_vpc1" {
		filter     = "vpcReference eq '${nutanix_vpc_v2.test-1.id}'"
		depends_on = [nutanix_vpc_v2.test-1]
	}

	resource "nutanix_routes_v2" "test-1" {
		name               = "%[1]s"
		description        = "%[2]s"
		vpc_reference      = nutanix_vpc_v2.test-1.id
		route_table_ext_id = data.nutanix_route_tables_v2.rt_vpc1.route_tables[0].ext_id
		destination {
			ipv4 {
				ip { value = "10.0.0.2" }
				prefix_length = 32
			}
		}
		next_hop {
			next_hop_type      = "EXTERNAL_SUBNET"
			next_hop_reference = nutanix_subnet_v2.test.id
		}
		metadata {
			owner_reference_id   = nutanix_vpc_v2.test-1.id
			project_reference_id = nutanix_project_v2.test.ext_id
		}
		route_type = "STATIC"
		%[6]s
		depends_on = [nutanix_project_v2.test]
		lifecycle {
			ignore_changes = [route_table_ext_id]
		}
	}

	data "nutanix_route_v2" "test" {
		route_table_ext_id = data.nutanix_route_tables_v2.rt_vpc1.route_tables[0].ext_id
		ext_id             = nutanix_routes_v2.test-1.id
		depends_on         = [nutanix_routes_v2.test-1]
	}

	data "nutanix_routes_v2" "test" {
		route_table_ext_id = data.nutanix_route_tables_v2.rt_vpc1.route_tables[0].ext_id
		filter             = "name eq '${nutanix_routes_v2.test-1.name}'"
		depends_on         = [nutanix_routes_v2.test-1]
	}
	`, name, desc, r, projectName, filepath, routeProjectExtIDLine(projectExtIDOverride))
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
		  lifecycle {
			ignore_changes = [links, snat_ips]
		  }
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
		project_reference_id = "00000000-0000-0000-0000-000000000000"
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
		project_reference_id = "00000000-0000-0000-0000-000000000000"
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
		project_reference_id = "00000000-0000-0000-0000-000000000000"
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
		project_reference_id = "00000000-0000-0000-0000-000000000000"
	  }
	  route_type = "STATIC"
	}
	`, name, desc)
}

// TestAccV2NutanixRoutesResource_MultipleNextHops covers the "support of
// multiple nexthops" change (next_hop is now a repeatable block backed by
// Nexthops []config.Nexthop). It creates a STATIC route with two
// EXTERNAL_SUBNET nexthops (ECMP across two external subnets), then shrinks it
// to a single nexthop to exercise the list expand/flatten + resize-on-update
// paths that the single-nexthop tests never touch.
func TestAccV2NutanixRoutesResource_MultipleNextHops(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-route-multinh-%d", r)
	desc := "route with multiple nexthops"
	resourceRoute := "nutanix_routes_v2.test-multinh"

	//goland:noinspection GoDeprecation
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// Create a route with two EXTERNAL_SUBNET nexthops.
			{
				Config: testRouteMultiNextHopConfig(name, desc, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRoute, "name", name),
					resource.TestCheckResourceAttr(resourceRoute, "description", desc),
					resource.TestCheckResourceAttrSet(resourceRoute, "vpc_reference"),
					resource.TestCheckResourceAttrSet(resourceRoute, "route_table_ext_id"),
					resource.TestCheckResourceAttr(resourceRoute, "destination.0.ipv4.0.ip.0.value", "10.0.0.2"),
					resource.TestCheckResourceAttr(resourceRoute, "next_hop.#", "2"),
					// Both nexthops are EXTERNAL_SUBNET, so these checks are order-independent.
					resource.TestCheckResourceAttr(resourceRoute, "next_hop.0.next_hop_type", "EXTERNAL_SUBNET"),
					resource.TestCheckResourceAttrSet(resourceRoute, "next_hop.0.next_hop_reference"),
					resource.TestCheckResourceAttr(resourceRoute, "next_hop.1.next_hop_type", "EXTERNAL_SUBNET"),
					resource.TestCheckResourceAttrSet(resourceRoute, "next_hop.1.next_hop_reference"),
					resource.TestCheckResourceAttr(resourceRoute, "route_type", "STATIC"),
				),
			},
			// Update: shrink from two nexthops down to one.
			{
				Config: testRouteMultiNextHopUpdateConfig(name, desc, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceRoute, "name", name+"_updated"),
					resource.TestCheckResourceAttr(resourceRoute, "description", desc+"_updated"),
					resource.TestCheckResourceAttr(resourceRoute, "next_hop.#", "1"),
					resource.TestCheckResourceAttr(resourceRoute, "next_hop.0.next_hop_type", "EXTERNAL_SUBNET"),
					resource.TestCheckResourceAttrSet(resourceRoute, "next_hop.0.next_hop_reference"),
					resource.TestCheckResourceAttr(resourceRoute, "route_type", "STATIC"),
				),
			},
		},
	})
}

// testRouteMultiNHBase provisions two no-NAT external subnets (subnet A on the
// primary VLAN fixture, subnet B on a distinct made-up VLAN/range) and a VPC that
// attaches both, so a route can reference each one as a separate EXTERNAL_SUBNET
// nexthop. Both are created no-NAT because a VPC may attach only 1 NAT external
// subnet but up to 4 no-NAT ones. Subnet B is created (not looked up): the only
// other external subnet in the environment is NAT-enabled, and the gc_subnet
// fixture ("VLAN 225") is a regular, non-external subnet.
func testRouteMultiNHBase(r int) string {
	return fmt.Sprintf(`
	data "nutanix_clusters_v2" "clusters" {}

	locals {
	  config   = jsondecode(file("%[1]s"))
	  subnets  = local.config.networking.subnets
	  cluster1 = [
		for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
		  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
	  ][0]
	}

	resource "nutanix_subnet_v2" "test" {
	  name              = "tf-route-multinh-subnet1-%[2]d"
	  description       = "external subnet A for multi-nexthop route test"
	  cluster_reference = local.cluster1
	  subnet_type       = "VLAN"
	  network_id        = local.subnets.vlan_id
	  is_external       = true
	  # No-NAT: a VPC may attach at most 1 NAT external subnet but up to 4 VLAN
	  # no-NAT external subnets, so both nexthop subnets are created no-NAT.
	  is_nat_enabled    = false
	  ip_config {
		ipv4 {
		  ip_subnet {
			ip { value = local.subnets.network_ip }
			prefix_length = local.subnets.network_prefix
		  }
		  default_gateway_ip { value = local.subnets.gateway_ip }
		  pool_list {
			start_ip { value = local.subnets.dhcp.start_ip }
			end_ip   { value = local.subnets.dhcp.end_ip }
		  }
		}
	  }
	}

	# External subnet B. Created on a distinct, unused VLAN id with a made-up
	# routable range (like the vpcs_v2 tests do) rather than reusing gc_subnet:
	# gc_subnet ("VLAN 225") is a regular non-external subnet, and the only
	# pre-existing external subnet in the environment is NAT-enabled (which cannot
	# be attached to an ONLY_NONAT VPC). A different VLAN id is required because
	# creating a second subnet with the same VLAN id on the same DVS is rejected.
	resource "nutanix_subnet_v2" "test2" {
	  name              = "tf-route-multinh-subnet2-%[2]d"
	  description       = "external subnet B for multi-nexthop route test"
	  cluster_reference = local.cluster1
	  subnet_type       = "VLAN"
	  network_id        = %[3]d
	  is_external       = true
	  is_nat_enabled    = false
	  ip_config {
		ipv4 {
		  ip_subnet {
			ip { value = "192.168.230.0" }
			prefix_length = 24
		  }
		  default_gateway_ip { value = "192.168.230.1" }
		  pool_list {
			start_ip { value = "192.168.230.20" }
			end_ip   { value = "192.168.230.30" }
		  }
		}
	  }
	}

	resource "nutanix_vpc_v2" "test-1" {
	  name        = "tf-route-multinh-vpc-%[2]d"
	  description = "vpc for multi-nexthop route test"
	  # Allow attaching more than one no-NAT external subnet to this VPC.
	  supported_multiple_external_subnet_type = "ONLY_NONAT"
	  external_subnets { subnet_reference = nutanix_subnet_v2.test.id }
	  external_subnets { subnet_reference = nutanix_subnet_v2.test2.id }
	  depends_on = [nutanix_subnet_v2.test, nutanix_subnet_v2.test2]
	  lifecycle {
		ignore_changes = [links, snat_ips]
	  }
	}

	data "nutanix_route_tables_v2" "rt_vpc1" {
	  filter     = "vpcReference eq '${nutanix_vpc_v2.test-1.id}'"
	  depends_on = [nutanix_vpc_v2.test-1]
	}
	`, filepath, r, 900+r%90)
}

func testRouteMultiNextHopConfig(name, desc string, r int) string {
	return testRouteMultiNHBase(r) + fmt.Sprintf(`
	resource "nutanix_routes_v2" "test-multinh" {
	  name               = "%[1]s"
	  description        = "%[2]s"
	  vpc_reference      = nutanix_vpc_v2.test-1.id
	  route_table_ext_id = data.nutanix_route_tables_v2.rt_vpc1.route_tables[0].ext_id
	  destination {
		ipv4 {
		  ip { value = "10.0.0.2" }
		  prefix_length = 32
		}
	  }
	  next_hop {
		next_hop_type      = "EXTERNAL_SUBNET"
		next_hop_reference = nutanix_subnet_v2.test.id
	  }
	  next_hop {
		next_hop_type      = "EXTERNAL_SUBNET"
		next_hop_reference = nutanix_subnet_v2.test2.id
	  }
	  metadata {
		owner_reference_id   = nutanix_vpc_v2.test-1.id
		project_reference_id = "00000000-0000-0000-0000-000000000000"
	  }
	  route_type = "STATIC"
	  lifecycle {
		ignore_changes = [route_table_ext_id]
	  }
	}
	`, name, desc)
}

func testRouteMultiNextHopUpdateConfig(name, desc string, r int) string {
	return testRouteMultiNHBase(r) + fmt.Sprintf(`
	resource "nutanix_routes_v2" "test-multinh" {
	  name               = "%[1]s_updated"
	  description        = "%[2]s_updated"
	  vpc_reference      = nutanix_vpc_v2.test-1.id
	  route_table_ext_id = data.nutanix_route_tables_v2.rt_vpc1.route_tables[0].ext_id
	  destination {
		ipv4 {
		  ip { value = "10.0.0.2" }
		  prefix_length = 32
		}
	  }
	  next_hop {
		next_hop_type      = "EXTERNAL_SUBNET"
		next_hop_reference = nutanix_subnet_v2.test.id
	  }
	  metadata {
		owner_reference_id   = nutanix_vpc_v2.test-1.id
		project_reference_id = "00000000-0000-0000-0000-000000000000"
	  }
	  route_type = "STATIC"
	  lifecycle {
		ignore_changes = [route_table_ext_id]
	  }
	}
	`, name, desc)
}
