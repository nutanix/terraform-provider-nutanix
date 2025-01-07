package foundation_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccFoundationImageNodesResource(t *testing.T) {
	name := "batch1"
	resourcePath := "nutanix_foundation_image_nodes." + name

	// get file file path to config having nodes info
	path, _ := os.Getwd()
	filepath := path + "/../../../test_foundation_config.json"

	// using block 0 in the test_foundation_config.json for this testcase
	blockNum := 0

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testImageNodesResource(filepath, blockNum, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourcePath, "session_id"),
					resource.TestCheckResourceAttr(resourcePath, "cluster_urls.0.cluster_name", foundationVars.Blocks[0].Nodes[0].HypervisorHostname),
				),
			},
		},
	})
}

// Checks negative scenario for a given invalid nos file name
func TestAccFoundationImageNodesResource_InvalidNosError(t *testing.T) {
	name := "batch1"

	// get file file path to config having nodes info
	path, _ := os.Getwd()
	filepath := path + "/../../../test_foundation_config.json"

	// using block 0 in the test_foundation_config.json for this testcase
	blockNum := 0

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testImageNodesResourceInvalidNosError(filepath, blockNum, name),
				ExpectError: regexp.MustCompile("Node imaging process failed due to error: Couldn't find nos_package at"),
			},
		},
	})
}

func testImageNodesResource(filepath string, blockNum int, name string) string {
	return fmt.Sprintf(`
	data "nutanix_foundation_nos_packages" "nos"{}

	data "nutanix_foundation_hypervisor_isos" "hypervisor" {}

	locals{
		config = (jsondecode(file("%s"))).blocks[%v]
	}

	resource "nutanix_foundation_image_nodes" "%s" {
	  timeouts {
		create = "80m"
	  }
	  cvm_gateway = local.config.cvm_gateway
	  hypervisor_gateway = local.config.hypervisor_gateway
	  cvm_netmask = local.config.cvm_netmask
	  hypervisor_netmask = local.config.cvm_netmask
	  ipmi_user = local.config.ipmi_user
	  nos_package = data.nutanix_foundation_nos_packages.nos.entities[2]
	  hypervisor_iso{
	  kvm{
	  	filename = data.nutanix_foundation_hypervisor_isos.hypervisor.kvm.3.filename
		checksum = "Needs_sha256_checksum"
	  	}
	  }
	  blocks {
		// manual mode using ipmi (Bare Metal, AOS or DOS nodes)
		nodes{
			cvm_ip = local.config.nodes[0].cvm_ip
			hypervisor_ip = local.config.nodes[0].hypervisor_ip
			hypervisor = local.config.nodes[0].hypervisor
			hypervisor_hostname = local.config.nodes[0].hypervisor_hostname
			ipmi_ip = local.config.nodes[0].ipmi_ip
			ipmi_netmask = local.config.nodes[0].ipmi_netmask
			ipmi_gateway = local.config.nodes[0].ipmi_gateway
			ipmi_user = local.config.nodes[0].ipmi_user
			ipmi_password = local.config.nodes[0].ipmi_password
			node_position = local.config.nodes[0].node_position
		}
		// using cvm (AOS or DOS nodes)
		nodes{
			cvm_ip = local.config.nodes[1].cvm_ip
			hypervisor_ip = local.config.nodes[1].hypervisor_ip
			hypervisor = local.config.nodes[1].hypervisor
			hypervisor_hostname = local.config.nodes[1].hypervisor_hostname
			ipmi_ip = local.config.nodes[1].ipmi_ip
			ipv6_address = local.config.nodes[1].ipv6_address
			current_network_interface = local.config.nodes[1].current_network_interface
			node_position = local.config.nodes[1].node_position
			ipmi_netmask = local.config.nodes[1].ipmi_netmask
			ipmi_gateway = local.config.nodes[1].ipmi_gateway
			ipmi_password = local.config.nodes[1].ipmi_password
		}
		// using cvm (AOS or DOS nodes)
		nodes{
			cvm_ip = local.config.nodes[2].cvm_ip
			hypervisor_ip = local.config.nodes[2].hypervisor_ip
			hypervisor = local.config.nodes[2].hypervisor
			hypervisor_hostname = local.config.nodes[2].hypervisor_hostname
			ipmi_ip = local.config.nodes[2].ipmi_ip
			ipv6_address = local.config.nodes[2].ipv6_address
			current_network_interface = local.config.nodes[2].current_network_interface
			node_position = local.config.nodes[2].node_position
			device_hint = "vm_installer"
			ipmi_netmask = local.config.nodes[2].ipmi_netmask
			ipmi_gateway = local.config.nodes[2].ipmi_gateway
		}
		block_id = local.config.block_id
	  }
	  clusters {
		cluster_members = [
			local.config.nodes[0].cvm_ip,
			local.config.nodes[1].cvm_ip,
			local.config.nodes[2].cvm_ip
		]
		redundancy_factor = 2
		cluster_name =  local.config.nodes[0].hypervisor_hostname
		timezone = "Asia/Kolkata"
	  }
	}`, filepath, blockNum, name)
}

func testImageNodesResourceInvalidNosError(filepath string, blockNum int, name string) string {
	return fmt.Sprintf(`
	locals{
		config = (jsondecode(file("%s"))).blocks[%v]
	}
	
	data "nutanix_foundation_hypervisor_isos" "hypervisor" {}

	resource "nutanix_foundation_image_nodes" "%s" {
	  timeouts {
		create = "80m"
	  }
	  cvm_gateway = local.config.cvm_gateway
	  hypervisor_gateway = local.config.hypervisor_gateway
	  cvm_netmask = local.config.cvm_netmask
	  hypervisor_netmask = local.config.cvm_netmask
	  ipmi_user = local.config.ipmi_user
	  nos_package = "ironman"
	  hypervisor{
	  	kvm {
			filename = data.nutanix_foundation_hypervisor_isos.hypervisor.kvm.0.filename
			checksum = "Needs_sha256_checksum"
		}
	  }
	  blocks {
		// manual mode using ipmi (Bare Metal, AOS or DOS nodes)
		nodes{
			cvm_ip = local.config.nodes[0].cvm_ip
			hypervisor_ip = local.config.nodes[0].hypervisor_ip
			hypervisor = local.config.nodes[0].hypervisor
			hypervisor_hostname = local.config.nodes[0].hypervisor_hostname
			ipmi_ip = local.config.nodes[0].ipmi_ip
			ipmi_netmask = local.config.nodes[0].ipmi_netmask
			ipmi_gateway = local.config.nodes[0].ipmi_gateway
			ipmi_user = local.config.nodes[0].ipmi_user
			ipmi_password = local.config.nodes[0].ipmi_password
			node_position = local.config.nodes[0].node_position
		}
		// using cvm (AOS or DOS nodes)
		nodes{
			cvm_ip = local.config.nodes[1].cvm_ip
			hypervisor_ip = local.config.nodes[1].hypervisor_ip
			hypervisor = local.config.nodes[1].hypervisor
			hypervisor_hostname = local.config.nodes[1].hypervisor_hostname
			ipmi_ip = local.config.nodes[1].ipmi_ip
			// ipv6_address = local.config.nodes[1].ipv6_address
			current_network_interface = local.config.nodes[1].current_network_interface
			node_position = local.config.nodes[1].node_position
			device_hint = "vm_installer"
		}
		// using cvm (AOS or DOS nodes)
		nodes{
			cvm_ip = local.config.nodes[2].cvm_ip
			hypervisor_ip = local.config.nodes[2].hypervisor_ip
			hypervisor = local.config.nodes[2].hypervisor
			hypervisor_hostname = local.config.nodes[2].hypervisor_hostname
			ipmi_ip = local.config.nodes[2].ipmi_ip
			// ipv6_address = local.config.nodes[2].ipv6_address
			current_network_interface = local.config.nodes[2].current_network_interface
			node_position = local.config.nodes[2].node_position
			device_hint = "vm_installer"
		}
		block_id = local.config.block_id
	  }
	  clusters {
		cluster_members = [
			local.config.nodes[0].cvm_ip,
			local.config.nodes[1].cvm_ip,
			local.config.nodes[2].cvm_ip
		]
		redundancy_factor = 2
		cluster_name =  local.config.nodes[0].hypervisor_hostname
	  }
	}`, filepath, blockNum, name)
}
