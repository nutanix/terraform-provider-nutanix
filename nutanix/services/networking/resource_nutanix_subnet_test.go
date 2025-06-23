package networking_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const resourceNameSubnet = "nutanix_subnet.acctest-managed"

func TestAccNutanixSubnet_basic(t *testing.T) {
	r := randIntBetween(31, 40)
	subnetName := fmt.Sprintf("acctest-managed-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixSubnetConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixSubnetExists(resourceNameSubnet),
					resource.TestCheckResourceAttr(resourceNameSubnet, "name", subnetName),
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
	r := randIntBetween(41, 50)
	subnetName := fmt.Sprintf("acctest-managed-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixSubnetConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixSubnetExists(resourceNameSubnet),
					resource.TestCheckResourceAttr(resourceNameSubnet, "name", subnetName),
					resource.TestCheckResourceAttr(resourceNameSubnet, "description", "Description of my unit test VLAN"),
				),
			},
			{
				Config: testAccNutanixSubnetConfigUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixSubnetExists(resourceNameSubnet),
					resource.TestCheckResourceAttr(resourceNameSubnet, "name", fmt.Sprintf("acctest-managed-updateName-%d", r)),
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
	r := randIntBetween(51, 60)
	resourceName := "nutanix_subnet.acctest-managed-categories"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixSubnetConfigWithCategory(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixSubnetExists(resourceName),
					testAccCheckNutanixCategories(resourceName),
					resource.TestCheckResourceAttr(resourceName, "categories.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "categories.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "categories.0.value"),
					resource.TestCheckResourceAttr(resourceName, "categories.0.name", "Environment"),
					resource.TestCheckResourceAttr(resourceName, "categories.0.value", "Production"),
				),
			},
			{
				Config: testAccNutanixSubnetConfigWithCategoryUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixSubnetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "categories.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "categories.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "categories.0.value"),
					resource.TestCheckResourceAttr(resourceName, "categories.0.name", "Environment"),
					resource.TestCheckResourceAttr(resourceName, "categories.0.value", "Staging"),
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

func TestAccNutanixSubnet_ExternalSubnet(t *testing.T) {
	r := randIntBetween(331, 340)
	subnetName := fmt.Sprintf("acctest-managed-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixSubnetConfigExternalSubnet(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixSubnetExists(resourceNameSubnet),
					resource.TestCheckResourceAttr(resourceNameSubnet, "name", subnetName),
					resource.TestCheckResourceAttr(resourceNameSubnet, "description", "Description of my unit test VLAN"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "subnet_type", "VLAN"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "subnet_ip", "10.250.140.0"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "default_gateway_ip", "10.250.140.1"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "is_external", "true"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "enable_nat", "true"),
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

func TestAccNutanixSubnet_ExternalSubnetWithNoNAT(t *testing.T) {
	r := randIntBetween(41, 50)
	subnetName := fmt.Sprintf("acctest-managed-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixSubnetConfigExternalSubnetNoNAT(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixSubnetExists(resourceNameSubnet),
					resource.TestCheckResourceAttr(resourceNameSubnet, "name", subnetName),
					resource.TestCheckResourceAttr(resourceNameSubnet, "description", "Description of my unit test VLAN"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "subnet_type", "VLAN"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "subnet_ip", "10.250.140.0"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "default_gateway_ip", "10.250.140.1"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "is_external", "true"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "enable_nat", "false"),
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

func TestAccNutanixSubnet_OverlaySubnet(t *testing.T) {
	r := randIntBetween(51, 60)
	subnetName := fmt.Sprintf("acctest-managed-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixSubnetConfigOverlaySubnet(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixSubnetExists(resourceNameSubnet),
					resource.TestCheckResourceAttr(resourceNameSubnet, "name", subnetName),
					resource.TestCheckResourceAttr(resourceNameSubnet, "description", "Description of my unit test OVERLAY"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "subnet_type", "OVERLAY"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "subnet_ip", "10.250.140.0"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "default_gateway_ip", "10.250.140.1"),
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

func TestAccNutanixSubnet_withIpPoolListRanges(t *testing.T) {
	r := randIntBetween(61, 70)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
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
	r := randIntBetween(71, 80)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNutanixSubnetConfigWithIPPoolListRangesErrored(r),
				ExpectError: regexp.MustCompile("please see https://developer.nutanix.com/reference/prism_central/v3/#definitions-ip_pool"),
			},
		},
	})
}

func TestAccNutanixSubnet_nameDuplicated(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetResourceConfigNameDuplicated(randIntBetween(21, 30)),
				// ExpectError: regexp.MustCompile("subnet already with name"),
			},
		},
	})
}

func TestAccNutanixSubnet_WithVlan0(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixSubnetConfig(0),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixSubnetExists(resourceNameSubnet),
					resource.TestCheckResourceAttr(resourceNameSubnet, "name", "acctest-managed-0"),
					resource.TestCheckResourceAttr(resourceNameSubnet, "description", "Description of my unit test VLAN"),
				),
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

func testAccCheckNutanixCategories(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		if val, ok := rs.Primary.Attributes["categories.0.name"]; !ok || val == "" {
			return fmt.Errorf("%s: manual Attribute '%s' expected to be set", n, "categories.0.name")
		}

		return nil
	}
}

func testAccCheckNutanixSubnetDestroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)

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

func resourceNutanixSubnetExists(conn *v3.Client, name string) (*string, error) {
	var subnetUUID *string

	filter := fmt.Sprintf("name==%s", name)
	subnetList, err := conn.V3.ListAllSubnet(filter, nil)
	if err != nil {
		return nil, err
	}

	for _, subnet := range subnetList.Entities {
		if utils.StringValue(subnet.Status.Name) == name {
			subnetUUID = subnet.Metadata.UUID
		}
	}
	return subnetUUID, nil
}

func testAccNutanixSubnetConfig(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}

locals {
	cluster1 = [
	for cluster in data.nutanix_clusters.clusters.entities :
	cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
	][0]
}

resource "nutanix_subnet" "acctest-managed" {
  	cluster_uuid = local.cluster1
	name        = "acctest-managed-%[1]d"
	description = "Description of my unit test VLAN"
	vlan_id     = %[1]d
	subnet_type = "VLAN"
	subnet_ip          = "10.250.140.0"
  default_gateway_ip = "10.250.140.1"
  prefix_length = 24
  dhcp_options = {
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

locals {
	cluster1 = [
	for cluster in data.nutanix_clusters.clusters.entities :
	cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
	][0]
}
resource "nutanix_subnet" "acctest-managed" {
  # What cluster will this VLAN live on?
  cluster_uuid = local.cluster1

  # General Information for subnet
	name        = "acctest-managed-updateName-%[1]d"
	description = "Description of my unit test VLAN updated"
  vlan_id     = %[1]d
	subnet_type = "VLAN"

  # Provision a Managed L3 Network
  # This bit is only needed if you intend to turn on AHV's IPAM
	subnet_ip          = "10.250.141.0"
  default_gateway_ip = "10.250.141.1"
  prefix_length = 24
  dhcp_options = {
		boot_file_name   = "bootfile"
		domain_name      = "nutanix"
		tftp_server_name = "10.250.141.200"
	}

	dhcp_server_address = {
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

locals {
	cluster1 = [
	for cluster in data.nutanix_clusters.clusters.entities :
	cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
	][0]
}

resource "nutanix_subnet" "acctest-managed-categories" {
  # What cluster will this VLAN live on?
  cluster_uuid = local.cluster1

  # General Information for subnet
	name        = "acctest-managed-%[1]d"
    vlan_id     = %[1]d
	subnet_type = "VLAN"

  # Provision a Managed L3 Network
  # This bit is only needed if you intend to turn on AHV's IPAM
	subnet_ip          = "10.250.140.0"
  default_gateway_ip = "10.250.140.1"
  prefix_length = 24
  dhcp_options = {
		boot_file_name   = "bootfile"
		domain_name      = "nutanix"
		tftp_server_name = "10.250.140.200"
	}
	dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
	dhcp_domain_search_list      = ["terraform.nutanix.com", "terraform.unit.test.com"]

	categories {
		name = "Environment"
		value = "Production"
	}
}
`, r)
}

func testAccNutanixSubnetConfigWithCategoryUpdate(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}

locals {
	cluster1 = [
	for cluster in data.nutanix_clusters.clusters.entities :
	cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
	][0]
}

resource "nutanix_subnet" "acctest-managed-categories" {
  # What cluster will this VLAN live on?
  cluster_uuid = local.cluster1

  # General Information for subnet
	name        = "acctest-managed-%[1]d"
    vlan_id     = %[1]d
	subnet_type = "VLAN"

  # Provision a Managed L3 Network
  # This bit is only needed if you intend to turn on AHV's IPAM
	subnet_ip          = "10.250.140.0"
  default_gateway_ip = "10.250.140.1"
  prefix_length = 24
  dhcp_options = {
		boot_file_name   = "bootfile"
		domain_name      = "nutanix"
		tftp_server_name = "10.250.140.200"
	}
	dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
	dhcp_domain_search_list      = ["terraform.nutanix.com", "terraform.unit.test.com"]

	categories {
		name = "Environment"
		value = "Staging"
	}
}
`, r)
}

func testAccNutanixSubnetConfigWithIPPoolListRanges(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}

locals {
	cluster1 = [
	for cluster in data.nutanix_clusters.clusters.entities :
	cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
	][0]
}

resource "nutanix_subnet" "acctest-managed-categories" {
  # What cluster will this VLAN live on?
  cluster_uuid = local.cluster1

  # General Information for subnet
	name        = "acctest-managed"
    vlan_id     = %d
	subnet_type = "VLAN"

  # Provision a Managed L3 Network
  # This bit is only needed if you intend to turn on AHV's IPAM
	subnet_ip          = "10.250.140.0"
  default_gateway_ip = "10.250.140.1"
  prefix_length = 24
  dhcp_options = {
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

locals {
	cluster1 = [
	for cluster in data.nutanix_clusters.clusters.entities :
	cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
	][0]
}

resource "nutanix_subnet" "acctest-managed-categories" {
  # What cluster will this VLAN live on?
  cluster_uuid = local.cluster1

  # General Information for subnet
	name        = "acctest-managed"
  vlan_id     = %d
	subnet_type = "VLAN"

  # Provision a Managed L3 Network
  # This bit is only needed if you intend to turn on AHV's IPAM
	subnet_ip          = "10.250.140.0"
  default_gateway_ip = "10.250.140.1"
  prefix_length = 24
  dhcp_options = {
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

func testAccSubnetResourceConfigNameDuplicated(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}

locals {
	cluster1 = [
	for cluster in data.nutanix_clusters.clusters.entities :
	cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
	][0]
}

resource "nutanix_subnet" "test" {
	name = "dou_vlan0_test_%d"
	cluster_uuid = local.cluster1
	vlan_id = %d
	subnet_type = "VLAN"

	prefix_length = 24
	default_gateway_ip = "192.168.0.1"
	subnet_ip = "192.168.0.0"
	ip_config_pool_list_ranges = ["192.168.0.10 192.168.0.100"]

	dhcp_options = {
		boot_file_name   = "bootfile"
		domain_name      = "nutanix"
		tftp_server_name = "10.250.140.200"
	}

	dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
	dhcp_domain_search_list      = ["terraform.nutanix.com", "terraform.unit.test.com"]
}

resource "nutanix_subnet" "test1" {
	name = nutanix_subnet.test.name
	cluster_uuid= local.cluster1
	vlan_id = %d
	subnet_type = "VLAN"
	prefix_length = 24
	default_gateway_ip = "192.168.0.1"
	subnet_ip = "192.168.0.0"
	#ip_config_pool_list_ranges = ["192.168.0.5", "192.168.0.100"]

	dhcp_options = {
		boot_file_name   = "bootfile"
		domain_name      = "nutanix"
		tftp_server_name = "10.250.140.200"
	}

	dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
	dhcp_domain_search_list      = ["terraform.nutanix.com", "terraform.unit.test.com"]
}
`, r, r, r+2)
}

func testAccNutanixSubnetConfigExternalSubnet(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}

locals {
	cluster1 = [
	for cluster in data.nutanix_clusters.clusters.entities :
	cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
	][0]
}

resource "nutanix_subnet" "acctest-managed" {
  	cluster_uuid = local.cluster1
	name        = "acctest-managed-%[1]d"
	description = "Description of my unit test VLAN"
	vlan_id     = %[1]d
	subnet_type = "VLAN"
	subnet_ip          = "10.250.140.0"
  	default_gateway_ip = "10.250.140.1"
	ip_config_pool_list_ranges = ["10.250.140.15 10.250.140.30"]

  	prefix_length = 24
	is_external = true
}
`, r)
}

func testAccNutanixSubnetConfigExternalSubnetNoNAT(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}

locals {
	cluster1 = [
	for cluster in data.nutanix_clusters.clusters.entities :
	cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
	][0]
}

resource "nutanix_subnet" "acctest-managed" {
  	cluster_uuid = local.cluster1
	name        = "acctest-managed-%[1]d"
	description = "Description of my unit test VLAN"
	vlan_id     = %[1]d
	subnet_type = "VLAN"
	subnet_ip          = "10.250.140.0"
  	default_gateway_ip = "10.250.140.1"
	ip_config_pool_list_ranges = ["10.250.140.15 10.250.140.30"]

  	prefix_length = 24
	is_external = true
	enable_nat = false
}
`, r)
}

func testAccNutanixSubnetConfigOverlaySubnet(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}

locals {
	cluster1 = [
	for cluster in data.nutanix_clusters.clusters.entities :
	cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
	][0]
}

resource "nutanix_subnet" "sub-ext" {
  	cluster_uuid = local.cluster1
	name        = "acctest-managed-ext-%[1]d"
	description = "Description of my unit test VLAN"
	vlan_id     = %[1]d
	subnet_type = "VLAN"
	subnet_ip          = "10.250.142.0"
  	default_gateway_ip = "10.250.142.1"
	ip_config_pool_list_ranges = ["10.250.142.15 10.250.142.30"]

  	prefix_length = 24
	is_external = true
	enable_nat = false
}

resource "nutanix_vpc" "acctest-managed-vpc" {
	name = "acctest-managed-%[1]d"
	external_subnet_reference_uuid = [
	  resource.nutanix_subnet.sub-ext.id
	]
	common_domain_name_server_ip_list{
		ip = "8.8.8.9"
	}
	externally_routable_prefix_list{
	  ip=  "172.50.0.0"
	  prefix_length= 16
	}
  }

  resource "nutanix_subnet" "acctest-managed" {
	name        = "acctest-managed-%[1]d"
	description = "Description of my unit test OVERLAY"
	vpc_reference_uuid = resource.nutanix_vpc.acctest-managed-vpc.id
	subnet_type = "OVERLAY"
	subnet_ip          = "10.250.140.0"
	default_gateway_ip = "10.250.140.1"
	ip_config_pool_list_ranges = ["10.250.140.15 10.250.140.30"]
	prefix_length = 24
}

`, r)
}
