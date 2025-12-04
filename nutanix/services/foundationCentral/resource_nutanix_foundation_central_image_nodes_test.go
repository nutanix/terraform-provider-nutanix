package fc_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccFCImageNodesResource(t *testing.T) {
	name := "batch2"
	resourcePath := "nutanix_foundation_central_image_cluster." + name
	clusterName := "test_cluster"
	// get file file path to config having nodes info
	path, _ := os.Getwd()
	filepath := path + "/../test_foundation_config.json"

	// using block 1 in the test_foundation_config.json for this testcase
	blockNum := 1
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testFCImageNodesResource(filepath, blockNum, name, clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourcePath, "cluster_name", clusterName),
				),
			},
		},
	})
}

func testFCImageNodesResource(filepath string, blockNum int, name, clusterName string) string {
	return fmt.Sprintf(`
	locals{
		config = (jsondecode(file("%[1]s"))).blocks[%[2]v]
	}
	
	resource "nutanix_foundation_central_image_cluster" "%[3]s" {
	  aos_package_url = local.config.aos_package_url

	  node_list{
			cvm_gateway = local.config.cvm_gateway
			cvm_netmask = local.config.cvm_netmask
			cvm_ip = local.config.nodes[0].cvm_ip
			hypervisor_gateway = local.config.hypervisor_gateway
			hypervisor_netmask = local.config.hypervisor_netmask
			hypervisor_ip = local.config.nodes[0].hypervisor_ip
			hypervisor_hostname = local.config.nodes[0].hypervisor_hostname
			imaged_node_uuid = local.config.nodes[0].imaged_node_uuid
			use_existing_network_settings = local.config.use_existing_network_settings
			ipmi_ip = local.config.nodes[0].ipmi_ip
			ipmi_netmask = local.config.nodes[0].ipmi_netmask
			ipmi_gateway = local.config.nodes[0].ipmi_gateway
			image_now = local.config.image_now
			hypervisor_type = local.config.nodes[0].hypervisor_type
		}
		node_list{
			cvm_gateway = local.config.cvm_gateway
			cvm_netmask = local.config.cvm_netmask
			cvm_ip = local.config.nodes[1].cvm_ip
			hypervisor_gateway = local.config.hypervisor_gateway
			hypervisor_netmask = local.config.hypervisor_netmask
			hypervisor_ip = local.config.nodes[1].hypervisor_ip
			hypervisor_hostname = local.config.nodes[1].hypervisor_hostname
			imaged_node_uuid = local.config.nodes[1].imaged_node_uuid
			use_existing_network_settings = local.config.use_existing_network_settings
			ipmi_ip = local.config.nodes[1].ipmi_ip
			ipmi_netmask = local.config.nodes[1].ipmi_netmask
			ipmi_gateway = local.config.nodes[1].ipmi_gateway
			image_now = local.config.image_now
			hypervisor_type = local.config.nodes[1].hypervisor_type
		}
		node_list{
			cvm_gateway = local.config.cvm_gateway
			cvm_netmask = local.config.cvm_netmask
			cvm_ip = local.config.nodes[2].cvm_ip
			hypervisor_gateway = local.config.hypervisor_gateway
			hypervisor_netmask = local.config.hypervisor_netmask
			hypervisor_ip = local.config.nodes[2].hypervisor_ip
			hypervisor_hostname = local.config.nodes[2].hypervisor_hostname
			imaged_node_uuid = local.config.nodes[2].imaged_node_uuid
			use_existing_network_settings = local.config.use_existing_network_settings
			ipmi_ip = local.config.nodes[2].ipmi_ip
			ipmi_netmask = local.config.nodes[2].ipmi_netmask
			ipmi_gateway = local.config.nodes[2].ipmi_gateway
			image_now = local.config.image_now
			hypervisor_type = local.config.nodes[2].hypervisor_type
		}
		common_network_settings{
			cvm_dns_servers = [local.config.common_network_settings.cvm_dns_servers[0]]
			hypervisor_dns_servers = [local.config.common_network_settings.hypervisor_dns_servers[0]]
			cvm_ntp_servers = [local.config.common_network_settings.cvm_ntp_servers[0]]
			hypervisor_ntp_servers = [local.config.common_network_settings.hypervisor_ntp_servers[0]]
		}
		redundancy_factor = 2
		cluster_name =  "%[4]s"
	}`, filepath, blockNum, name, clusterName)
}
