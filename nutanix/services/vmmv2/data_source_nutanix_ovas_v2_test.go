package vmmv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameOvas = "data.nutanix_ovas_v2.test"

func TestAccV2NutanixOvaDatasource_ListAllOvas(t *testing.T) {
	r := acctest.RandIntRange(1, 999)
	vmName := fmt.Sprintf("tf-test-vm-ova-%d", r)
	vmDescription := "VM for OVA terraform testing"
	ovaName := fmt.Sprintf("tf-test-ova-%d", r)

	config := testOvaResourceConfigCreateOvaFromVM(vmName, vmDescription, ovaName)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// List all Ovas
			{
				Config: config + testOvasDatasourceConfigListAllOvas(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameOvas, "ovas.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixOvaDatasource_ListAllOvasWithFilter(t *testing.T) {
	r := acctest.RandIntRange(1, 999)
	vmName := fmt.Sprintf("tf-test-vm-ova-%d", r)
	vmDescription := "VM for OVA terraform testing"
	ovaName := fmt.Sprintf("tf-test-ova-%d", r)

	config := testOvaResourceConfigCreateOvaFromVM(vmName, vmDescription, ovaName)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// Filter Ovas by name
			{
				Config: config + testOvasDatasourceConfigFilterOvasByName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameOvas, "ovas.#"),
					// ova checks
					resource.TestCheckResourceAttrPair(resourceNameOva, "ext_id", datasourceNameOvas, "ovas.0.ext_id"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "size_bytes", datasourceNameOvas, "ovas.0.size_bytes"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "create_time", datasourceNameOvas, "ovas.0.create_time"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "name", datasourceNameOvas, "ovas.0.name"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "parent_vm", datasourceNameOvas, "ovas.0.parent_vm"),
					resource.TestCheckResourceAttrPair(resourceNameOva, "disk_format", datasourceNameOvas, "ovas.0.disk_format"),
				),
			},
		},
	})
}

func TestAccV2NutanixOvaDatasource_ListAllOvasWithLimit(t *testing.T) {
	r := acctest.RandIntRange(1, 999)
	vmName := fmt.Sprintf("tf-test-vm-ova-%d", r)
	vmDescription := "VM for OVA terraform testing"
	ovaName := fmt.Sprintf("tf-test-ova-%d", r)

	config := testOvaResourceConfigCreateOvaFromVM(vmName, vmDescription, ovaName)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// Filter Ovas by name
			{
				Config: config + testOvasDatasourceConfigLimit(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameOvas, "ovas.#"),
					resource.TestCheckResourceAttr(datasourceNameOvas, "ovas.#", "1"),
				),
			},
		},
	})
}
func testOvasDatasourceConfigListAllOvas() string {
	return `
data "nutanix_ovas_v2" "test" {}
`
}

func testOvasDatasourceConfigFilterOvasByName() string {
	return `
data "nutanix_ovas_v2" "test" {
	filter = "name eq '${nutanix_ova_v2.test.name}'"
}
`
}

func testOvasDatasourceConfigLimit() string {
	return `
data "nutanix_ovas_v2" "test" {
	limit = 1
	depends_on = [nutanix_ova_v2.test]
}
`
}
