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

func TestAccV2NutanixVMAntiAffinityPolicyResource_WithProjectAssociation(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-vm-anti-affinity-policy-%d", r)
	desc := "test vm anti affinity policy description"
	projectName := fmt.Sprintf("test-vm-afp-pa-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVMAntiAffinityPolicyV2ConfigWithProjectAssociation(name, desc, projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVMAntiAffinityPolicy, "name", name),
					resource.TestCheckResourceAttr(resourceNameVMAntiAffinityPolicy, "description", desc),
					resource.TestCheckResourceAttr(resourceNameVMAntiAffinityPolicy, "categories.#", "2"),
					resource.TestCheckResourceAttrPair(resourceNameVMAntiAffinityPolicy, "project_ext_id", "nutanix_project_v2.test", "ext_id"),
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

func testVMAntiAffinityPolicyV2ConfigWithProjectAssociation(name, desc, projectName string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 = [
				for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
				cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
			][0]
		}

		resource "nutanix_project_v2" "test" {
			name        = "%[4]s"
			project_id  = "%[4]s"
			description = "vm anti affinity policy project association test"
		}

		resource "nutanix_resource_group_v2" "rg" {
			name           = "tf-vmaap-pa-rg-%[4]s"
			project_ext_id = nutanix_project_v2.test.ext_id
			placement_targets {
				cluster_ext_id = local.cluster0
			}
			# Ignore changes to placement_targets to avoid perpetual diffs after apply.
			lifecycle {
				ignore_changes = [placement_targets]
			}
		}

		resource "nutanix_category_v2" "vm_category" {
			count = %[3]d
		    key = "vm-anti-affinity-vm-category-${count.index}"
			value = "vm-anti-affinity-vm-category-value-${count.index}"
			project_ext_id = nutanix_project_v2.test.ext_id
		}

		resource "nutanix_vm_anti_affinity_policy_v2" "test" {
			name = "%[1]s"
			description = "%[2]s"
			categories = nutanix_category_v2.vm_category[*].id
			project_ext_id = nutanix_project_v2.test.ext_id
			depends_on = [nutanix_resource_group_v2.rg]
		}
	`, name, desc, 2, projectName)
}
