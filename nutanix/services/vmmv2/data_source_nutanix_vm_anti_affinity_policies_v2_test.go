package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameVMAntiAffinityPolicies = "data.nutanix_vm_anti_affinity_policies_v2.test"

func TestAccV2NutanixVMAntiAffinityPoliciesDatasource_Basic(t *testing.T) {
	name := "test-vm-anti-affinity-policy-"
	desc := "test vm anti affinity policy description"
	count := "3"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testPreConfig(name, desc, count) + testVMAntiAffinityPoliciesV2(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameVMAntiAffinityPolicies, "policies.#", count),
				),
			},
		},
	})
}

func TestAccV2NutanixVMAntiAffinityPoliciesDatasource_WithFilter(t *testing.T) {
	name := "test-vm-anti-affinity-policy-"
	desc := "test vm anti affinity policy description"
	count := "3"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testPreConfig(name, desc, count) + testVMAntiAffinityPoliciesV2WithFilter(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameVMAntiAffinityPolicies, "policies.0.name", fmt.Sprintf(`%[1]s0`, name)),
					resource.TestCheckResourceAttr(datasourceNameVMAntiAffinityPolicies, "policies.#", "1"),
				),
			},
		},
	})
}

func TestAccV2NutanixVMAntiAffinityPoliciesDatasource_WithInvalidFilters(t *testing.T) {
	name := "test-vm-anti-affinity-policy-"
	desc := "test vm anti affinity policy description"
	count := "3"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testPreConfig(name, desc, count) + testVMAntiAffinityPoliciesV2WithInvalidFilter(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameVMAntiAffinityPolicies, "policies.#", "0"),
				),
			},
		},
	})
}

func testPreConfig(name, desc, count string) string {
	return fmt.Sprintf(`
		resource "nutanix_category_v2" "vm_category" {
		    count = %[3]s
		    key = "vm-anti-affinity-vm-category"
			value = "vm-anti-affinity-vm-category-value-${count.index}"
		}

		resource "nutanix_vm_anti_affinity_policy_v2" "test" {
		    count = %[3]s
			name = "%[1]s${count.index}"
			description = "%[2]s"
			categories = [ nutanix_category_v2.vm_category[count.index].id ]
		}
	`, name, desc, count)
}

func testVMAntiAffinityPoliciesV2() string {
	return `
		data "nutanix_vm_anti_affinity_policies_v2" "test" {
		    depends_on = [
				resource.nutanix_vm_anti_affinity_policy_v2.test
			]
		}
	`
}

func testVMAntiAffinityPoliciesV2WithFilter() string {
	return `
		data "nutanix_vm_anti_affinity_policies_v2" "test" {
			filter="name eq '${nutanix_vm_anti_affinity_policy_v2.test[0].name}'"
		    depends_on = [
				resource.nutanix_vm_anti_affinity_policy_v2.test
			]
		}
	`
}

func testVMAntiAffinityPoliciesV2WithInvalidFilter() string {
	return `
		data "nutanix_vm_anti_affinity_policies_v2" "test" {
			filter = "name eq 'invalid'"
		}
	`
}
