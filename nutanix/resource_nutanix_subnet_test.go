package nutanix

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const resourceNameSubnet = "nutanix_subnet.acctest-managed"

func TestAccNutanixSubnet_basic(t *testing.T) {
	r := acctest.RandIntRange(3500, 3900)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixSubnetConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixSubnetExists(resourceNameSubnet),
					resource.TestCheckResourceAttr(resourceNameSubnet, "name", "acctest-managed"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "description", "Description of my unit test VLAN"),
				),
			},
			{
				ResourceName:            resourceNameSubnet,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"description"},
			},
		},
	})
}

func TestAccNutanixSubnet_Update(t *testing.T) {
	r := acctest.RandIntRange(3500, 3900)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixSubnetConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixSubnetExists(resourceNameSubnet),
					resource.TestCheckResourceAttr(resourceNameSubnet, "name", "acctest-managed"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "description", "Description of my unit test VLAN"),
				),
			},
			{
				Config: testAccNutanixSubnetConfigUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixSubnetExists(resourceNameSubnet),
					resource.TestCheckResourceAttr(resourceNameSubnet, "name", "acctest-managed-updateName"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "description", "Description of my unit test VLAN updated"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "subnet_type", "VLAN"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "subnet_ip", "10.250.141.0"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "default_gateway_ip", "10.250.141.1"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "dhcp_options.tftp_server_name", "10.250.141.200"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "dhcp_server_address.ip", "10.250.141.254"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "dhcp_domain_name_server_list.1", "4.2.2.3"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "dhcp_domain_search_list.1", "terraform.uptated.test.com"),
				),
			},
			{
				ResourceName:            resourceNameSubnet,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"description"},
			},
		},
	})
}

func TestAccNutanixSubnet_WithCategory(t *testing.T) {
	r := acctest.RandIntRange(3500, 3900)
	resourceName := "nutanix_subnet.acctest-managed-categories"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixSubnetConfigWithCategory(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixSubnetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "categories.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "categories.environment-terraform", "production"),
				),
			},
			{
				Config: testAccNutanixSubnetConfigWithCategoryUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixSubnetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "categories.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "categories.environment-terraform", "staging"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"description"},
			},
		},
	})
}

func TestAccNutanixSubnet_withIpPoolListRanges(t *testing.T) {
	r := acctest.RandIntRange(3500, 3900)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config:             testAccNutanixSubnetConfigWithIPPoolListRanges(r),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccNutanixSubnet_withIpPoolListRangesErrored(t *testing.T) {
	r := acctest.RandIntRange(3500, 3900)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNutanixSubnetConfigWithIPPoolListRangesErrored(r),
				ExpectError: regexp.MustCompile("please see https://developer.nutanix.com/reference/prism_central/v3/#definitions-ip_pool"),
			},
		},
	})
}

func testAccCheckNutanixSubnetExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		return nil
	}
}

func testAccCheckNutanixSubnetDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_subnet" {
			continue
		}
		if _, err := resourceNutanixSubnetExists(conn.API, rs.Primary.ID); err != nil {
			if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
				return nil
			}
			return err
		}
	}

	return nil
}

func testAccNutanixSubnetConfig(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}

output "cluster" {
  value = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
}

resource "nutanix_subnet" "acctest-managed" {
  # What cluster will this VLAN live on?
  cluster_reference = {
	kind = "cluster"
	uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }

  # General Information for subnet
	name        = "acctest-managed"
	description = "Description of my unit test VLAN"
  vlan_id     = %d
	subnet_type = "VLAN"

  # Provision a Managed L3 Network
  # This bit is only needed if you intend to turn on AHV's IPAM
	subnet_ip          = "10.250.140.0"
  default_gateway_ip = "10.250.140.1"
  prefix_length = 24
  dhcp_options {
		boot_file_name   = "bootfile"
		domain_name      = "nutanix"
		tftp_server_name = "10.250.140.200"
	}
	dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
	dhcp_domain_search_list      = ["terraform.nutanix.com", "terraform.unit.test.com"]
}
`, r)
}

func testAccNutanixSubnetConfigUpdate(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}

output "cluster" {
  value = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
}

resource "nutanix_subnet" "acctest-managed" {
  # What cluster will this VLAN live on?
  cluster_reference = {
	kind = "cluster"
	uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }

  # General Information for subnet
	name        = "acctest-managed-updateName"
	description = "Description of my unit test VLAN updated"
  vlan_id     = %d
	subnet_type = "VLAN"

  # Provision a Managed L3 Network
  # This bit is only needed if you intend to turn on AHV's IPAM
	subnet_ip          = "10.250.141.0"
  default_gateway_ip = "10.250.141.1"
  prefix_length = 24
  dhcp_options {
		boot_file_name   = "bootfile"
		domain_name      = "nutanix"
		tftp_server_name = "10.250.141.200"
	}

	dhcp_server_address {
		ip = "10.250.141.254"
	}

	dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.3"]
	dhcp_domain_search_list      = ["terraform.nutanix.com", "terraform.uptated.test.com"]
}
`, r)
}

func testAccNutanixSubnetConfigWithCategory(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}

output "cluster" {
  value = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
}

resource "nutanix_subnet" "acctest-managed-categories" {
  # What cluster will this VLAN live on?
  cluster_reference = {
	kind = "cluster"
	uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }

  # General Information for subnet
	name        = "acctest-managed"
  vlan_id     = %d
	subnet_type = "VLAN"

  # Provision a Managed L3 Network
  # This bit is only needed if you intend to turn on AHV's IPAM
	subnet_ip          = "10.250.140.0"
  default_gateway_ip = "10.250.140.1"
  prefix_length = 24
  dhcp_options {
		boot_file_name   = "bootfile"
		domain_name      = "nutanix"
		tftp_server_name = "10.250.140.200"
	}
	dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
	dhcp_domain_search_list      = ["terraform.nutanix.com", "terraform.unit.test.com"]

	categories {
		environment-terraform = "production"
	}
}
`, r)
}

func testAccNutanixSubnetConfigWithCategoryUpdate(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}

output "cluster" {
  value = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
}

resource "nutanix_subnet" "acctest-managed-categories" {
  # What cluster will this VLAN live on?
  cluster_reference = {
	kind = "cluster"
	uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }

  # General Information for subnet
	name        = "acctest-managed"
  vlan_id     = %d
	subnet_type = "VLAN"

  # Provision a Managed L3 Network
  # This bit is only needed if you intend to turn on AHV's IPAM
	subnet_ip          = "10.250.140.0"
  default_gateway_ip = "10.250.140.1"
  prefix_length = 24
  dhcp_options {
		boot_file_name   = "bootfile"
		domain_name      = "nutanix"
		tftp_server_name = "10.250.140.200"
	}
	dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
	dhcp_domain_search_list      = ["terraform.nutanix.com", "terraform.unit.test.com"]

	categories {
		environment-terraform = "staging"
	}
}
`, r)
}

func testAccNutanixSubnetConfigWithIPPoolListRanges(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}

output "cluster" {
  value = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
}

resource "nutanix_subnet" "acctest-managed-categories" {
  # What cluster will this VLAN live on?
  cluster_reference = {
	kind = "cluster"
	uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }

  # General Information for subnet
	name        = "acctest-managed"
  vlan_id     = %d
	subnet_type = "VLAN"

  # Provision a Managed L3 Network
  # This bit is only needed if you intend to turn on AHV's IPAM
	subnet_ip          = "10.250.140.0"
  default_gateway_ip = "10.250.140.1"
  prefix_length = 24
  dhcp_options {
		boot_file_name   = "bootfile"
		domain_name      = "nutanix"
		tftp_server_name = "10.250.140.200"
	}

	ip_config_pool_list_ranges= [
    "10.250.140.110 10.250.140.250"
  ]

	dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
	dhcp_domain_search_list      = ["terraform.nutanix.com", "terraform.unit.test.com"]
}
`, r)
}

func testAccNutanixSubnetConfigWithIPPoolListRangesErrored(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}

output "cluster" {
  value = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
}

resource "nutanix_subnet" "acctest-managed-categories" {
  # What cluster will this VLAN live on?
  cluster_reference = {
	kind = "cluster"
	uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }

  # General Information for subnet
	name        = "acctest-managed"
  vlan_id     = %d
	subnet_type = "VLAN"

  # Provision a Managed L3 Network
  # This bit is only needed if you intend to turn on AHV's IPAM
	subnet_ip          = "10.250.140.0"
  default_gateway_ip = "10.250.140.1"
  prefix_length = 24
  dhcp_options {
		boot_file_name   = "bootfile"
		domain_name      = "nutanix"
		tftp_server_name = "10.250.140.200"
	}

	ip_config_pool_list_ranges= [
    "10.250.140.110" #bad configuration
  ]

	dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
	dhcp_domain_search_list      = ["terraform.nutanix.com", "terraform.unit.test.com"]
}
`, r)
}
