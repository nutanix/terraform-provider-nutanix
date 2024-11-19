package vmmv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameVm = "data.nutanix_virtual_machines_v2.test"

func TestAccNutanixVmsDataSourceV2_List(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVMDataSourceConfigV4Vms(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameVm, "vms.#"),
				),
			},
		},
	})
}

func TestAccNutanixVmsDataSourceV2_ListWithFilters(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVMDataSourceConfigV4VmsWithFilters(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameVm, "vms.#"),
					resource.TestCheckResourceAttr(datasourceNameVm, "limit", "2"),
				),
			},
		},
	})
}

func TestAccNutanixVmsDataSourceV2_ListWithFilterName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVMDataSourceConfigV4VmsWithFilterName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameVm, "vms.#"),
					resource.TestCheckResourceAttr(datasourceNameVm, "vms.0.name", "test-vm-filter"),
					resource.TestCheckResourceAttr(datasourceNameVm, "vms.0.num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(datasourceNameVm, "vms.0.num_sockets", "1"),
				),
			},
		},
	})
}

func testAccVMDataSourceConfigV4Vms() string {
	return (`
		data "nutanix_virtual_machines_v2" "test" {
		}
`)
}

func testAccVMDataSourceConfigV4VmsWithFilters() string {
	return (`
		data "nutanix_virtual_machines_v2" "test" {
			page=0
			limit=2
		}
`)
}

func testAccVMDataSourceConfigV4VmsWithFilterName() string {
	return (`

	data "nutanix_clusters" "clusters" {}

	locals {
	cluster0 = data.nutanix_clusters.clusters.entities[0].metadata.uuid
	}

	resource "nutanix_virtual_machine_v2" "test"{
		name= "test-vm-filter"
		num_cores_per_socket = 1
		num_sockets = 1
		cluster {
			ext_id = local.cluster0
		}
	}
		data "nutanix_virtual_machines_v2" "test" {
			filter = "name eq 'test-vm-filter'"
			depends_on = [
				resource.nutanix_virtual_machine_v2.test
			]
		}
`)
}
