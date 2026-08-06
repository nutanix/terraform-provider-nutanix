package clusters_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameHost = "nutanix_host.acctest-managed"

func TestAccNutanixHost_WithCategory(t *testing.T) {
	imgName := fmt.Sprintf("test-acc-host-image-%s", acctest.RandString(3))
	vmName := fmt.Sprintf("test-acc-host-vm-%s", acctest.RandString(3))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixHostStillExists,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixHostConfigWithCategory(imgName, vmName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixHostExists(resourceNameHost),
					resource.TestCheckResourceAttrSet(resourceNameHost, "host_id"),
					resource.TestCheckResourceAttrSet(resourceNameHost, "name"),
					resource.TestCheckResourceAttr(resourceNameHost, "categories.#", "1"),
					resource.TestCheckResourceAttr(resourceNameHost, "categories.0.name", "Environment"),
					resource.TestCheckResourceAttr(resourceNameHost, "categories.0.value", "Production"),
				),
			},
			{
				Config: testAccNutanixHostConfigWithCategoryUpdate(imgName, vmName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixHostExists(resourceNameHost),
					resource.TestCheckResourceAttr(resourceNameHost, "categories.#", "1"),
					resource.TestCheckResourceAttr(resourceNameHost, "categories.0.name", "Environment"),
					resource.TestCheckResourceAttr(resourceNameHost, "categories.0.value", "Staging"),
				),
			},
			{
				ResourceName:      resourceNameHost,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckNutanixHostExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		conn := acc.TestAccProvider.Meta().(*conns.Client)
		if _, err := conn.API.V3.GetHost(rs.Primary.ID); err != nil {
			return fmt.Errorf("error fetching host (%s): %s", rs.Primary.ID, err)
		}

		return nil
	}
}

func testAccCheckNutanixHostStillExists(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_host" {
			continue
		}

		host, err := conn.API.V3.GetHost(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("host (%s) should still exist after destroying nutanix_host, but got: %s", rs.Primary.ID, err)
		}
		if len(host.Metadata.Categories) != 0 {
			return fmt.Errorf("expected categories on host (%s) to be cleared after destroy, got: %v", rs.Primary.ID, host.Metadata.Categories)
		}
	}

	return nil
}

func testAccNutanixHostConfigWithCategory(imgName, vmName string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
			cluster1 = "${data.nutanix_clusters.clusters.entities.0.service_list.0 == "PRISM_CENTRAL"
			? data.nutanix_clusters.clusters.entities.1.metadata.uuid : data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
		}

		resource "nutanix_image" "acctest-host-image" {
			name        = "%[1]s"
			source_uri  = "http://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img"
			description = "tiny image used to discover a host uuid for the nutanix_host acctest"
		}

		resource "nutanix_virtual_machine" "acctest-host-vm" {
			name                  = "%[2]s"
			cluster_uuid          = local.cluster1
			num_vcpus_per_socket  = 1
			num_sockets           = 1
			memory_size_mib       = 186

			disk_list {
				data_source_reference = {
					kind = "image"
					uuid = nutanix_image.acctest-host-image.id
				}

				device_properties {
					disk_address = {
						device_index = 0
						adapter_type = "IDE"
					}
					device_type = "CDROM"
				}
			}
		}

		resource "nutanix_host" "acctest-managed" {
			host_id = nutanix_virtual_machine.acctest-host-vm.host_reference.uuid

			categories {
				name  = "Environment"
				value = "Production"
			}
		}
	`, imgName, vmName)
}

func testAccNutanixHostConfigWithCategoryUpdate(imgName, vmName string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
			cluster1 = "${data.nutanix_clusters.clusters.entities.0.service_list.0 == "PRISM_CENTRAL"
			? data.nutanix_clusters.clusters.entities.1.metadata.uuid : data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
		}

		resource "nutanix_image" "acctest-host-image" {
			name        = "%[1]s"
			source_uri  = "http://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img"
			description = "tiny image used to discover a host uuid for the nutanix_host acctest"
		}

		resource "nutanix_virtual_machine" "acctest-host-vm" {
			name                  = "%[2]s"
			cluster_uuid          = local.cluster1
			num_vcpus_per_socket  = 1
			num_sockets           = 1
			memory_size_mib       = 186

			disk_list {
				data_source_reference = {
					kind = "image"
					uuid = nutanix_image.acctest-host-image.id
				}

				device_properties {
					disk_address = {
						device_index = 0
						adapter_type = "IDE"
					}
					device_type = "CDROM"
				}
			}
		}

		resource "nutanix_host" "acctest-managed" {
			host_id = nutanix_virtual_machine.acctest-host-vm.host_reference.uuid

			categories {
				name  = "Environment"
				value = "Staging"
			}
		}
	`, imgName, vmName)
}
