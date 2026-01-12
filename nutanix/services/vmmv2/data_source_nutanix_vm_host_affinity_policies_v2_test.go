package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameVMHostAffinityPolicies = "data.nutanix_vm_host_affinity_policies_v2.test"

func TestAccV2NutanixVMHostAffinityPoliciesDatasource_Basic(t *testing.T) {
	name := "test-vm-host-affinity-policy-"
	desc := "test vm host affinity policy description"
	count := "3"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVMHostAffinityPreConfig(name, desc, count) + testVMHostAffinityPoliciesV2(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameVMHostAffinityPolicies, "policies.#", count),
				),
			},
		},
	})
}

func TestAccV2NutanixVMHostAffinityPoliciesDatasource_WithFilter(t *testing.T) {
	name := "test-vm-host-affinity-policy-"
	desc := "test vm host affinity policy description"
	count := "3"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVMHostAffinityPreConfig(name, desc, count) + testVMHostAffinityPoliciesV2WithFilter(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameVMHostAffinityPolicies, "policies.0.name", fmt.Sprintf(`%[1]s0`, name)),
					resource.TestCheckResourceAttr(datasourceNameVMHostAffinityPolicies, "policies.#", "1"),
				),
			},
		},
	})
}

func TestAccV2NutanixVMHostAffinityPoliciesDatasource_WithInvalidFilters(t *testing.T) {
	name := "test-vm-host-affinity-policy-"
	desc := "test vm host affinity policy description"
	count := "3"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVMHostAffinityPreConfig(name, desc, count) + testVMHostAffinityPoliciesV2WithInvalidFilter(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameVMHostAffinityPolicies, "policies.#", "0"),
				),
			},
		},
	})
}

func testVMHostAffinityPreConfig(name, desc, count string) string {
	return fmt.Sprintf(`
		resource "nutanix_category_v2" "vm_category" {
		    count = %[3]s
		    key = "vm-host-affinity-vm-category"
			value = "vm-host-affinity-vm-category-value-${count.index}"
		}
		resource "nutanix_category_v2" "host_category" {
			count = %[3]s
		    key = "vm-host-affinity-host-category"
			value = "vm-host-affinity-host-category-value-${count.index}"
		}
		resource "nutanix_vm_host_affinity_policy_v2" "test" {
			count = %[3]s
			name = "%[1]s${count.index}"
			description = "%[2]s"
			vm_categories = [ nutanix_category_v2.vm_category[count.index].id ]
			host_categories = [ nutanix_category_v2.host_category[count.index].id ]
		}
	`, name, desc, count)
}

func testVMHostAffinityPoliciesV2() string {
	return `
		data "nutanix_vm_host_affinity_policies_v2" "test" {
		    depends_on = [
				resource.nutanix_vm_host_affinity_policy_v2.test
			]
		}
	`
}

func testVMHostAffinityPoliciesV2WithFilter() string {
	return `
		data "nutanix_vm_host_affinity_policies_v2" "test" {
			filter="name eq '${nutanix_vm_host_affinity_policy_v2.test[0].name}'"
		    depends_on = [
				resource.nutanix_vm_host_affinity_policy_v2.test
			]
		}
	`
}

func testVMHostAffinityPoliciesV2WithInvalidFilter() string {
	return `
		data "nutanix_vm_host_affinity_policies_v2" "test" {
			filter = "name eq 'invalid'"
		}
	`
}
