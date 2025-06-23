package networking_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameVpc = "nutanix_vpc.acctest-managed"

func TestAccNutanixVpc_basic(t *testing.T) {
	r := randIntBetween(400, 410)
	vpcName := fmt.Sprintf("acctest-managed-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixVpcConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVPCExists(resourceNameVpc),
					resource.TestCheckResourceAttr(resourceNameVpc, "name", vpcName),
					resource.TestCheckResourceAttr(resourceNameVpc, "common_domain_name_server_ip_list.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVpc, "externally_routable_prefix_list.#", "1"),
				),
			},
		},
	})
}

func TestAccNutanixVpc_Update(t *testing.T) {
	r := randIntBetween(411, 420)
	vpcName := fmt.Sprintf("acctest-managed-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixVpcConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVPCExists(resourceNameVpc),
					resource.TestCheckResourceAttr(resourceNameVpc, "name", vpcName),
					resource.TestCheckResourceAttr(resourceNameVpc, "common_domain_name_server_ip_list.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVpc, "externally_routable_prefix_list.#", "1"),
				),
			},
			{
				Config: testAccNutanixVpcConfigUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVPCExists(resourceNameVpc),
					resource.TestCheckResourceAttr(resourceNameVpc, "name", fmt.Sprintf("acctest-managed-updateName-%d", r)),
					resource.TestCheckResourceAttr(resourceNameVpc, "common_domain_name_server_ip_list.#", "2"),
					resource.TestCheckResourceAttr(resourceNameVpc, "externally_routable_prefix_list.#", "2"),
					resource.TestCheckResourceAttr(resourceNameVpc, "common_domain_name_server_ip_list.0.ip", "8.8.8.8"),
					resource.TestCheckResourceAttr(resourceNameVpc, "common_domain_name_server_ip_list.1.ip", "8.8.8.9"),
					resource.TestCheckResourceAttr(resourceNameVpc, "externally_routable_prefix_list.0.prefix_length", "16"),
					resource.TestCheckResourceAttr(resourceNameVpc, "externally_routable_prefix_list.1.prefix_length", "24"),
				),
			},
		},
	})
}

func TestAccNutanixVpc_WithSubnetName(t *testing.T) {
	r := randIntBetween(321, 430)
	vpcName := fmt.Sprintf("acctest-managed-%d", r)
	subnetName := fmt.Sprintf("acctest-subnet-%d", r)
	updateSubnetName := fmt.Sprintf("acctest-subnet-updated-%d", r)
	updatedVpcName := fmt.Sprintf("acctest-managed-updated-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixVpcConfigWithSubnetName(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVPCExists(resourceNameVpc),
					resource.TestCheckResourceAttr(resourceNameVpc, "name", vpcName),
					resource.TestCheckResourceAttr(resourceNameVpc, "external_subnet_reference_name.0", subnetName),
					resource.TestCheckResourceAttr(resourceNameVpc, "common_domain_name_server_ip_list.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVpc, "externally_routable_prefix_list.#", "1"),
				),
			},
			{
				Config: testAccNutanixVpcConfigWithSubnetNameUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVPCExists(resourceNameVpc),
					resource.TestCheckResourceAttr(resourceNameVpc, "name", updatedVpcName),
					resource.TestCheckResourceAttr(resourceNameVpc, "external_subnet_reference_name.0", updateSubnetName),
					resource.TestCheckResourceAttr(resourceNameVpc, "common_domain_name_server_ip_list.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVpc, "externally_routable_prefix_list.#", "2"),
				),
			},
		},
	})
}

func TestAccNutanixVpc_WithSubnetUUID(t *testing.T) {
	r := randIntBetween(431, 440)
	vpcName := fmt.Sprintf("acctest-managed-%d", r)
	updatedVpcName := fmt.Sprintf("acctest-managed-updateName-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixVpcConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVPCExists(resourceNameVpc),
					resource.TestCheckResourceAttr(resourceNameVpc, "name", vpcName),
					resource.TestCheckResourceAttr(resourceNameVpc, "common_domain_name_server_ip_list.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVpc, "externally_routable_prefix_list.#", "1"),
				),
			},
			{
				Config: testAccNutanixVpcConfigUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVPCExists(resourceNameVpc),
					resource.TestCheckResourceAttr(resourceNameVpc, "name", updatedVpcName),
					resource.TestCheckResourceAttr(resourceNameVpc, "common_domain_name_server_ip_list.#", "2"),
					resource.TestCheckResourceAttr(resourceNameVpc, "externally_routable_prefix_list.#", "2"),
				),
			},
		},
	})
}

func testAccCheckNutanixVpcDestroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	ctx := context.TODO()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_vpc" {
			continue
		}
		if _, err := conn.API.V3.GetVPC(ctx, rs.Primary.ID); err != nil {
			if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
				return nil
			}
			return err
		}
	}

	return nil
}

func testAccCheckNutanixVPCExists(n string) resource.TestCheckFunc {
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

func testAccNutanixVpcConfig(r int) string {
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
	  is_external = true
	  ip_config_pool_list_ranges = ["10.250.140.10 10.250.140.20"]
	}

	resource "nutanix_vpc" "acctest-managed" {
		name = "acctest-managed-%[1]d"


		external_subnet_reference_uuid = [
		  resource.nutanix_subnet.acctest-managed.id
		]

		common_domain_name_server_ip_list{
				ip = "8.8.8.9"
		}

		externally_routable_prefix_list{
		  ip=  "172.30.0.0"
		  prefix_length= 16
		}
	  }
	`, r)
}

func testAccNutanixVpcConfigUpdate(r int) string {
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
	is_external = true
	ip_config_pool_list_ranges = ["10.250.140.10 10.250.140.20"]
}

resource "nutanix_vpc" "acctest-managed" {
	name = "acctest-managed-updateName-%[1]d"


	external_subnet_reference_uuid = [
	  resource.nutanix_subnet.acctest-managed.id
	]

	common_domain_name_server_ip_list{
			ip = "8.8.8.8"
	}
	common_domain_name_server_ip_list{
			ip = "8.8.8.9"
	}

	externally_routable_prefix_list{
	  ip=  "172.51.0.0"
	  prefix_length= 16
	}
	externally_routable_prefix_list{
		ip=  "192.31.0.0"
		prefix_length= 24
	  }
  }


`, r)
}

func testAccNutanixVpcConfigWithSubnetName(r int) string {
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
		name        = "acctest-subnet-%[1]d"
		description = "Description of my unit test VLAN"
		vlan_id     = %[1]d
		subnet_type = "VLAN"
		subnet_ip          = "10.250.144.0"
	  default_gateway_ip = "10.250.144.1"
	  prefix_length = 24
	  is_external = true
	  ip_config_pool_list_ranges = ["10.250.144.10 10.250.144.20"]
	}

	resource "nutanix_vpc" "acctest-managed" {
		depends_on = [
			resource.nutanix_subnet.acctest-managed
		]
		name = "acctest-managed-%[1]d"

		external_subnet_reference_name = [
		  "acctest-subnet-%[1]d"
		]

		common_domain_name_server_ip_list{
				ip = "8.8.8.9"
		}

		externally_routable_prefix_list{
		  ip=  "172.30.0.0"
		  prefix_length= 16
		}
	  }
	`, r)
}

func testAccNutanixVpcConfigWithSubnetNameUpdate(r int) string {
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
		name        = "acctest-subnet-updated-%[1]d"
		description = "Description of my unit test VLAN"
		vlan_id     = %[1]d
		subnet_type = "VLAN"
		subnet_ip          = "10.250.144.0"
	  default_gateway_ip = "10.250.144.1"
	  prefix_length = 24
	  is_external = true
	  ip_config_pool_list_ranges = ["10.250.144.10 10.250.144.20"]
	}

	resource "nutanix_vpc" "acctest-managed" {
		depends_on = [
			resource.nutanix_subnet.acctest-managed
		]
		name = "acctest-managed-updated-%[1]d"

		external_subnet_reference_name = [
		  "acctest-subnet-updated-%[1]d"
		]

		common_domain_name_server_ip_list{
				ip = "8.8.8.9"
		}

		externally_routable_prefix_list{
		  ip=  "172.30.0.0"
		  prefix_length= 16
		}
		externally_routable_prefix_list{
			ip=  "172.34.0.0"
			prefix_length= 16
		  }
	  }
	`, r)
}
