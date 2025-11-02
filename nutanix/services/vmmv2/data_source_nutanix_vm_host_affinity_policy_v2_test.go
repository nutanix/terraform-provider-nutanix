package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameVMHostAffinityPolicy = "data.nutanix_vm_host_affinity_policy_v2.test"

func TestAccV2NutanixVMHostAffinityPolicyDatasource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-host-affinity-policy-%d", r)
	desc := "test vm host affinity policy description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVMHostAffinityPolicyDataSourceConfigV2(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameVMHostAffinityPolicy, "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameVMHostAffinityPolicy, "name", name),
					resource.TestCheckResourceAttr(datasourceNameVMHostAffinityPolicy, "description", desc),
					resource.TestCheckResourceAttrSet(datasourceNameVMHostAffinityPolicy, "create_time"),
					resource.TestCheckResourceAttrSet(datasourceNameVMHostAffinityPolicy, "update_time"),
					resource.TestCheckResourceAttr(datasourceNameVMHostAffinityPolicy, "host_categories.#", "1"),
					resource.TestCheckResourceAttr(datasourceNameVMHostAffinityPolicy, "vm_categories.#", "1"),
				),
			},
		},
	})
}

func testVMHostAffinityPolicyDataSourceConfigV2(name, desc string) string {
	return fmt.Sprintf(`
		resource "nutanix_category_v2" "vm_category" {
		    key = "vm-host-affinity-vm-category"
			value = "vm-host-affinity-vm-category-value"
		}

		resource "nutanix_category_v2" "host_category" {
		    key = "vm-host-affinity-host-category"
			value = "vm-host-affinity-host-category-value"
		}

		resource "nutanix_vm_host_affinity_policy_v2" "test" {
			name = "%[1]s"
			description = "%[2]s"
			host_categories = [ nutanix_category_v2.host_category.id ]
			vm_categories = [ nutanix_category_v2.vm_category.id ]
		}

		data "nutanix_vm_host_affinity_policy_v2" "test"{
			ext_id = resource.nutanix_vm_host_affinity_policy_v2.test.id
		}
`, name, desc)
}
