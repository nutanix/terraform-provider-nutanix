package vmmv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameVM = "data.nutanix_virtual_machines_v2.test"

func TestAccV2NutanixVmsDatasource_List(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVMDataSourceConfigV4Vms(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameVM, "vms.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixVmsDatasource_ListWithFilters(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVMDataSourceConfigV4VmsWithFilters(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameVM, "vms.#"),
					resource.TestCheckResourceAttr(datasourceNameVM, "limit", "2"),
				),
			},
		},
	})
}

func TestAccV2NutanixVmsDatasource_ListWithFilterName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVMDataSourceConfigV4VmsWithFilterName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameVM, "vms.#"),
					resource.TestCheckResourceAttr(datasourceNameVM, "vms.0.name", "tf-test-vm-filter"),
					resource.TestCheckResourceAttr(datasourceNameVM, "vms.0.num_cores_per_socket", "1"),
					resource.TestCheckResourceAttr(datasourceNameVM, "vms.0.num_sockets", "1"),
				),
			},
		},
	})
}

func TestAccV2NutanixVmsDatasource_ListWithInvalidFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVMDataSourceConfigV4VmsWithInvalidFilter(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameVM, "vms.#", "0"),
				),
			},
		},
	})
}

func testAccVMDataSourceConfigV4VmsWithInvalidFilter() string {
	return `
		data "nutanix_virtual_machines_v2" "test" {
			filter = "name eq 'invalid'"
		}
	`
}

func testAccVMDataSourceConfigV4Vms() string {
	return `
		data "nutanix_virtual_machines_v2" "test" {
		}
`
}

func testAccVMDataSourceConfigV4VmsWithFilters() string {
	return `
		data "nutanix_virtual_machines_v2" "test" {
			page=0
			limit=2
		}
`
}

func testAccVMDataSourceConfigV4VmsWithFilterName() string {
	return `

	data "nutanix_clusters_v2" "clusters" {}

	locals {
	cluster0 =  [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
			][0]
	}

	resource "nutanix_virtual_machine_v2" "test"{
		name= "tf-test-vm-filter"
		num_cores_per_socket = 1
		num_sockets = 1
		cluster {
			ext_id = local.cluster0
		}
	}
		data "nutanix_virtual_machines_v2" "test" {
			filter = "name eq 'tf-test-vm-filter'"
			depends_on = [
				resource.nutanix_virtual_machine_v2.test
			]
		}
`
}
