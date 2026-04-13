package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameVMAntiAffinityPolicy = "data.nutanix_vm_anti_affinity_policy_v2.test"

func TestAccV2NutanixVMAntiAffinityPolicyDatasource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-anti-affinity-policy-%d", r)
	desc := "test vm anti affinity policy description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVMAntiAffinityPolicyDataSourceConfigV2(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameVMAntiAffinityPolicy, "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameVMAntiAffinityPolicy, "name", name),
					resource.TestCheckResourceAttr(datasourceNameVMAntiAffinityPolicy, "description", desc),
					resource.TestCheckResourceAttrSet(datasourceNameVMAntiAffinityPolicy, "create_time"),
					resource.TestCheckResourceAttrSet(datasourceNameVMAntiAffinityPolicy, "update_time"),
					resource.TestCheckResourceAttr(datasourceNameVMAntiAffinityPolicy, "categories.#", "1"),
				),
			},
		},
	})
}

func testVMAntiAffinityPolicyDataSourceConfigV2(name, desc string) string {
	return fmt.Sprintf(`
		resource "nutanix_category_v2" "vm_category" {
		    key = "vm-anti-affinity-vm-category"
			value = "vm-anti-affinity-vm-category-value"
		}

		resource "nutanix_vm_anti_affinity_policy_v2" "test" {
			name = "%[1]s"
			description = "%[2]s"
			categories = [ nutanix_category_v2.vm_category.id ]
		}

		data "nutanix_vm_anti_affinity_policy_v2" "test"{
			ext_id = resource.nutanix_vm_anti_affinity_policy_v2.test.id
		}
`, name, desc)
}
