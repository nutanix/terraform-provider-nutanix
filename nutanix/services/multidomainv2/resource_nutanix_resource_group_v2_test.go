package multidomainv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestAccV2NutanixResourceGroupResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-rg-%d", r)
	proj_name := fmt.Sprintf("tf-test-proj-%d", r)
	updateName := fmt.Sprintf("tf-test-rg-%d-update", r)
	containerName := "SelfServiceContainer"
	newContainerName := fmt.Sprintf("tf-test-sc-%d", r)
	dataSourceGet := "data.nutanix_resource_group_v2.resource_group_get"
	dataSourceList := "data.nutanix_resource_groups_v2.resource_groups"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testResourceGroupV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGroupV2ResourceConfig(name, containerName, proj_name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameResourceGroupV2, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameResourceGroupV2, "name", name),
					resource.TestCheckResourceAttr(resourceNameResourceGroupV2, "placement_targets.#", "1"),
					resource.TestCheckResourceAttr(resourceNameResourceGroupV2, "placement_targets.0.storage_containers.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceGet, "ext_id"),
					resource.TestCheckResourceAttr(dataSourceGet, "name", name),
					resource.TestCheckResourceAttr(dataSourceGet, "placement_targets.#", "1"),
					resource.TestCheckResourceAttr(dataSourceGet, "placement_targets.0.storage_containers.#", "1"),
					resource.TestCheckResourceAttr(dataSourceList, "resource_groups.#", "1"),
					resource.TestCheckResourceAttr(dataSourceList, "resource_groups.0.name", name),
					resource.TestCheckResourceAttr(dataSourceList, "resource_groups.0.placement_targets.#", "1"),
					resource.TestCheckResourceAttr(dataSourceList, "resource_groups.0.placement_targets.0.storage_containers.#", "1"),
					resource.TestCheckResourceAttrPair(dataSourceGet, "placement_targets.0.storage_containers.0", dataSourceList, "resource_groups.0.placement_targets.0.storage_containers.0"),
				),
			},
			{
				Config: testAccResourceGroupV2ResourceUpdateConfig(updateName, containerName, newContainerName, proj_name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameResourceGroupV2, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameResourceGroupV2, "name", updateName),
					resource.TestCheckResourceAttr(resourceNameResourceGroupV2, "placement_targets.#", "1"),
					resource.TestCheckResourceAttr(resourceNameResourceGroupV2, "placement_targets.0.storage_containers.#", "2"),
					resource.TestCheckResourceAttrSet(resourceNameResourceGroupV2, "placement_targets.0.storage_containers.0.ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameResourceGroupV2, "placement_targets.0.storage_containers.1.ext_id"),
				),
			},
		},
	})
}

func TestAccV2NutanixResourceGroupResource_ListWithInvalidFilter(t *testing.T) {
	dataSourceList := "data.nutanix_resource_groups_v2.resource_groups"
	randomUUID := utils.GenUUID()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGroupV2ListWithInvalidFilterConfig(randomUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceList, "resource_groups.#", "0"),
				),
			},
		},
	})
}

func testAccResourceGroupV2ListWithInvalidFilterConfig(uuid string) string {
	return fmt.Sprintf(`
	data "nutanix_resource_groups_v2" "resource_groups" {
		filter = "projectExtId eq '%s'"
	}
`, uuid)
}

func testAccResourceGroupV2ResourceConfig(name string, containerName string, proj_name string) string {
	return fmt.Sprintf(`
	data "nutanix_clusters_v2" "clusters" {}
  data "nutanix_storage_containers_v2" "storage_container" {}
	locals {
		cluster_ext_id = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		][0]
		storage_container_ext_id = [
			for storage_container in data.nutanix_storage_containers_v2.storage_container.storage_containers :
			storage_container.ext_id if storage_container.name == "%s"
		][0]
	}
	resource "nutanix_project_v2" "example" {
		name        = "terra-%s"
		project_id = "terra-%s"
		description = "Example project for multidomain namespace"
  }
  resource "nutanix_resource_group_v2" "test" {
		name           = "%s"
		project_ext_id = nutanix_project_v2.example.ext_id
		placement_targets {
			cluster_ext_id = local.cluster_ext_id
			storage_containers {
				ext_id = local.storage_container_ext_id
			}
		}
	}
	data "nutanix_resource_group_v2" "resource_group_get" {
		ext_id = nutanix_resource_group_v2.test.id
	}
	data "nutanix_resource_groups_v2" "resource_groups" {
		filter = "projectExtId eq '${nutanix_project_v2.example.ext_id}'"
		depends_on = [nutanix_resource_group_v2.test]
	}
`, containerName, proj_name, proj_name, name)
}

// testAccResourceGroupV2ResourceUpdateConfig provisions a brand-new storage
// container and adds it to the placement target alongside the existing
// user-managed container. The multidomain API is additive for a placement
// target's storage containers (an update never removes the existing one), so the
// resource group ends up referencing both.
func testAccResourceGroupV2ResourceUpdateConfig(name string, containerName string, newContainerName string, proj_name string) string {
	return fmt.Sprintf(`
	data "nutanix_clusters_v2" "clusters" {}
	data "nutanix_storage_containers_v2" "storage_container" {}
	locals {
		cluster_ext_id = [
			for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		][0]
		storage_container_ext_id = [
			for storage_container in data.nutanix_storage_containers_v2.storage_container.storage_containers :
			storage_container.ext_id if storage_container.name == "%s"
		][0]
	}
	resource "nutanix_storage_containers_v2" "new_container" {
		name                                     = "%s"
		cluster_ext_id                           = local.cluster_ext_id
		logical_advertised_capacity_bytes        = 1073741824000
		logical_explicit_reserved_capacity_bytes = 20
		replication_factor                       = 1
		erasure_code                             = "OFF"
		is_inline_ec_enabled                     = false
		has_higher_ec_fault_domain_preference    = false
		cache_deduplication                      = "OFF"
		on_disk_dedup                            = "OFF"
		is_compression_enabled                   = true
		is_internal                              = false
		is_software_encryption_enabled           = false
	}
	resource "nutanix_project_v2" "example" {
		name        = "terra-%s"
		project_id = "terra-%s"
		description = "Example project for multidomain namespace"
	}
	resource "nutanix_resource_group_v2" "test" {
		name           = "%s"
		project_ext_id = nutanix_project_v2.example.ext_id
		placement_targets {
			cluster_ext_id = local.cluster_ext_id
			storage_containers {
				ext_id = local.storage_container_ext_id
			}
			storage_containers {
				ext_id = nutanix_storage_containers_v2.new_container.container_ext_id
			}
		}
	}
	data "nutanix_resource_group_v2" "resource_group_get" {
		ext_id = nutanix_resource_group_v2.test.id
	}
	data "nutanix_resource_groups_v2" "resource_groups" {
		filter = "projectExtId eq '${nutanix_project_v2.example.ext_id}'"
		depends_on = [nutanix_resource_group_v2.test]
	}
`, containerName, newContainerName, proj_name, proj_name, name)
}
