package networking_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccNutanixVPCListDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_vpcs.test", "entities.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_vpcs.test", "api_version"),
				),
			},
		},
	})
}

func TestAccNutanixVPCsDataSource_UUID(t *testing.T) {
	r := randIntBetween(25, 45)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCsDataSourceConfigWithUUID(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nutanix_vpcs.test", "entities.#"),
					resource.TestCheckResourceAttrSet("data.nutanix_vpcs.test", "api_version"),
					func(s *terraform.State) error {
						res := s.RootModule().Resources["data.nutanix_vpcs.test"]
						if res == nil {
							return fmt.Errorf("data.nutanix_vpcs.test not found")
						}

						vpcName := fmt.Sprintf("acctest-managed-%d", r)
						attrs := res.Primary.Attributes

						index, err := findVPCIndex(attrs, vpcName)
						if err != nil {
							return err
						}

						prefix := fmt.Sprintf("entities.%d.spec.0.resources.0", index)

						if got := attrs[fmt.Sprintf("%s.externally_routable_prefix_list.0.prefix_length", prefix)]; got != "16" {
							return fmt.Errorf("expected prefix_length 16, got %q", got)
						}
						if got := attrs[fmt.Sprintf("%s.externally_routable_prefix_list.0.ip", prefix)]; got != "172.36.0.0" {
							return fmt.Errorf("expected ip 172.36.0.0, got %q", got)
						}
						if got := attrs[fmt.Sprintf("%s.common_domain_name_server_ip_list.0.ip", prefix)]; got != "8.8.8.9" {
							return fmt.Errorf("expected DNS ip 8.8.8.9, got %q", got)
						}

						return nil
					},
				),
			},
		},
	})
}

func testAccVPCsDataSourceConfig() string {
	return (`
	data "nutanix_vpcs" "test" {
	}
`)
}

func testAccVPCsDataSourceConfigWithUUID(r int) string {
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

resource "nutanix_vpc" "test" {
	name = "acctest-managed-%[1]d"


	external_subnet_reference_uuid = [
	  resource.nutanix_subnet.acctest-managed.id
	]

	common_domain_name_server_ip_list{
			ip = "8.8.8.9"
	}

	externally_routable_prefix_list{
	  ip=  "172.36.0.0"
	  prefix_length= 16
	}
  }
	data "nutanix_vpcs" "test" {
		depends_on = [
			resource.nutanix_vpc.test
		]
	}
`, r)
}

func findVPCIndex(attributes map[string]string, targetName string) (int, error) {
	prefix := "entities."
	for i := 0; ; i++ {
		nameKey := fmt.Sprintf("%s%d.spec.0.name", prefix, i)
		if name, ok := attributes[nameKey]; ok {
			if name == targetName {
				return i, nil
			}
		} else {
			break // No more VPCs
		}
	}
	return -1, fmt.Errorf("no VPC found with name %q", targetName)
}
