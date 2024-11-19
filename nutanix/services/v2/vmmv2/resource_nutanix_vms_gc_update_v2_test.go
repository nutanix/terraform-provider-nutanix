package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameVmGCUpdate = "nutanix_vm_gc_update_v2.test"

func TestAccNutanixVmsGCUpdateV2Resource_Basic(t *testing.T) {
	// t.Skip("Skipping test as it requires GCUpdate")
	r := acctest.RandInt()
	vmName := fmt.Sprintf("tf-test-vm-%d", r)
	// stateOn := "power_on"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmPreEnvConfig(r) + testVmConfig(vmName) + testVmsGCUpdateV2Config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameVmGCUpdate, "ext_id"),
				),
			},
		},
	})
}

func testVmsGCUpdateV2Config() string {
	return `
		resource "nutanix_vm_gc_update_v2" "test" {
			ext_id = resource.nutanix_virtual_machine_v2.test-vm.id
			config{
				cloud_init{
					cloud_init_script{
						user_data{
							value="${local.gs}"		
						}
					}
				}
			}	  
		}
		  
`
}

func testVmConfig(vmName string) string {
	return fmt.Sprintf(`
		resource "nutanix_virtual_machine_v2" "test-vm" {
		  name                 = "%[1]s"
		  description          = "vm for testing "
		  num_cores_per_socket = 1
		  num_sockets          = 1
		  memory_size_bytes    = 4 * 1024 * 1024 * 1024
		  cluster {
			ext_id = local.clusterUUID
		  }
					 
				  
		  power_state = "OFF"

		  lifecycle {
			ignore_changes = [guest_tools, nics, cd_roms]
		  }
		
		  depends_on = [data.nutanix_clusters.clusters, data.nutanix_storage_containers_v2.sc]
		}
`, vmName)
}
