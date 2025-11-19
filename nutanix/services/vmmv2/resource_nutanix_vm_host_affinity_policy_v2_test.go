package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameVMHostAffinityPolicy = "nutanix_vm_host_affinity_policy_v2.test"

func TestAccV2NutanixVMHostAffinityPolicyResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-host-affinity-policy-%d", r)
	desc := "test vm host affinity policy description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVMHostAffinityPolicyV2Config(name, desc, 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVMHostAffinityPolicy, "name", name),
					resource.TestCheckResourceAttr(resourceNameVMHostAffinityPolicy, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVMHostAffinityPolicy, "vm_categories.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVMHostAffinityPolicy, "host_categories.#", "1"),
				),
			},
		},
	})
}

func TestAccV2NutanixVMHostAffinityPolicyResource_WithUpdate(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-host-affinity-policy-%d", r)
	updatedName := fmt.Sprintf("test-vm-host-affinity-policy-%d-updated", r)
	desc := "test vm host affinity policy description"
	updatedDesc := "test vm host affinity policy description updated"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVMHostAffinityPolicyV2Config(name, desc, 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVMHostAffinityPolicy, "name", name),
					resource.TestCheckResourceAttr(resourceNameVMHostAffinityPolicy, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVMHostAffinityPolicy, "vm_categories.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVMHostAffinityPolicy, "host_categories.#", "1"),
				),
			},
			{
				Config: testVMHostAffinityPolicyV2Config(updatedName, updatedDesc, 2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVMHostAffinityPolicy, "name", updatedName),
					resource.TestCheckResourceAttr(resourceNameVMHostAffinityPolicy, "description", updatedDesc),
					resource.TestCheckResourceAttr(resourceNameVMHostAffinityPolicy, "vm_categories.#", "2"),
					resource.TestCheckResourceAttr(resourceNameVMHostAffinityPolicy, "host_categories.#", "2"),
				),
			},
		},
	})
}

func testVMHostAffinityPolicyV2Config(name, desc string, count int) string {
	return fmt.Sprintf(`
		resource "nutanix_category_v2" "vm_category" {
			count = %[3]d
		    key = "vm-host-affinity-vm-category"
			value = "vm-host-affinity-vm-category-value-${count.index}"
		}
		resource "nutanix_category_v2" "host_category" {
			count = %[3]d
		    key = "vm-host-affinity-host-category"
			value = "vm-host-affinity-host-category-value-${count.index}"
		}
		resource "nutanix_vm_host_affinity_policy_v2" "test" {
			name = "%[1]s"
			description = "%[2]s"
			vm_categories = nutanix_category_v2.vm_category[*].id
			host_categories = nutanix_category_v2.host_category[*].id
		}
`, name, desc, count)
}
