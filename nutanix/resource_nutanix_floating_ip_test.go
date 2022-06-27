package nutanix

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const resourceNameFloatingIP = "nutanix_floating_ip.acctest-managed"

func TestAccNutanixFloatingIP_basic(t *testing.T) {
	r := randIntBetween(31, 40)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixFloatingIPConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameFloatingIP, "external_subnet_reference_uuid"),
				),
			},
		},
	})
}

func TestAccNutanixFloatingIP_WithVPC(t *testing.T) {
	r := randIntBetween(31, 40)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixFloatingIPConfigWithVpc(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameFloatingIP, "external_subnet_reference_uuid"),
					resource.TestCheckResourceAttrSet(resourceNameFloatingIP, "vpc_reference_uuid"),
					resource.TestCheckResourceAttr(resourceNameFloatingIP, "private_ip", "10.3.3.6"),
				),
			},
		},
	})
}

func testAccCheckNutanixFloatingIPDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*Client)
	ctx := context.TODO()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_floating_ip" {
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

func testAccNutanixFloatingIPConfig(r int) string {
	return fmt.Sprintf(`

	data "nutanix_clusters" "clusters" {}
	
	locals {
		cluster1 = [
		for cluster in data.nutanix_clusters.clusters.entities :
		cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
		][0]
	}
	
	resource "nutanix_subnet" "sub-test" {
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

	resource "nutanix_floating_ip" "acctest-managed" {
		external_subnet_reference_uuid = resource.nutanix_subnet.sub-test.id
	}
	`, r)
}

func testAccNutanixFloatingIPConfigWithVpc(r int) string {
	return fmt.Sprintf(`

	data "nutanix_clusters" "clusters" {}
	
	locals {
		cluster1 = [
		for cluster in data.nutanix_clusters.clusters.entities :
		cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
		][0]
	}
	
	resource "nutanix_subnet" "sub-test" {
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

	resource "nutanix_vpc" "test-vpc" {
		name = "acctest-vpc-%[1]d"
	  
	  
		external_subnet_reference_uuid = [
		  resource.nutanix_subnet.sub-test.id
		]
	  
		common_domain_name_server_ip_list{
				ip = "8.8.8.9"
		}
	  
		externally_routable_prefix_list{
		  ip=  "172.31.0.0"
		  prefix_length= 16
		}
	  }

	resource "nutanix_floating_ip" "acctest-managed" {
		external_subnet_reference_uuid = resource.nutanix_subnet.sub-test.id
		vpc_reference_uuid= resource.nutanix_vpc.test-vpc.id
		private_ip = "10.3.3.6"
	}
	`, r)
}
