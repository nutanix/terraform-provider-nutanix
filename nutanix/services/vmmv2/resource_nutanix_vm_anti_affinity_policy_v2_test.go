package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameVMAntiAffinityPolicy = "nutanix_vm_anti_affinity_policy_v2.test"

func TestAccV2NutanixVMAntiAffinityPolicyResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-anti-affinity-policy-%d", r)
	desc := "test vm anti affinity policy description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVMAntiAffinityPolicyV2Config(name, desc, 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVMAntiAffinityPolicy, "name", name),
					resource.TestCheckResourceAttr(resourceNameVMAntiAffinityPolicy, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVMAntiAffinityPolicy, "categories.#", "1"),
				),
			},
		},
	})
}

func TestAccV2NutanixVMAntiAffinityPolicyResource_WithUpdate(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-anti-affinity-policy-%d", r)
	updatedName := fmt.Sprintf("test-vm-anti-affinity-policy-%d-updated", r)
	desc := "test vm anti affinity policy description"
	updatedDesc := "test vm anti affinity policy description updated"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVMAntiAffinityPolicyV2Config(name, desc, 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVMAntiAffinityPolicy, "name", name),
					resource.TestCheckResourceAttr(resourceNameVMAntiAffinityPolicy, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVMAntiAffinityPolicy, "categories.#", "1"),
				),
			},
			{
				Config: testVMAntiAffinityPolicyV2Config(updatedName, updatedDesc, 2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVMAntiAffinityPolicy, "name", updatedName),
					resource.TestCheckResourceAttr(resourceNameVMAntiAffinityPolicy, "description", updatedDesc),
					resource.TestCheckResourceAttr(resourceNameVMAntiAffinityPolicy, "categories.#", "2"),
				),
			},
		},
	})
}

func testVMAntiAffinityPolicyV2Config(name, desc string, count int) string {
	return fmt.Sprintf(`
		resource "nutanix_category_v2" "vm_category" {
			count = %[3]d
		    key = "vm-anti-affinity-vm-category"
			value = "vm-anti-affinity-vm-category-value-${count.index}"
		}

		resource "nutanix_vm_anti_affinity_policy_v2" "test" {
			name = "%[1]s"
			description = "%[2]s"
			categories = nutanix_category_v2.vm_category[*].id
		}
`, name, desc, count)
}
