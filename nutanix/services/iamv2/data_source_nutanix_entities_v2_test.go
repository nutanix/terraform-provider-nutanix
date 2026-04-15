package iamv2_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameEntities = "data.nutanix_iam_entities_v2.test"

func TestAccV2NutanixEntitiesDatasource_List(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testEntitiesDatasourceV2Config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameEntities, "id"),
					resource.TestCheckResourceAttrSet(datasourceNameEntities, "entities.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixEntitiesDatasource_ListWithLimit(t *testing.T) {
	limit := 1
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testEntitiesDatasourceV2ConfigWithLimit(limit),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameEntities, "id"),
					resource.TestCheckResourceAttr(datasourceNameEntities, "entities.#", strconv.Itoa(limit)),
				),
			},
		},
	})
}

func TestAccV2NutanixEntitiesDatasource_ListWithFilter(t *testing.T) {
	filter := "name eq 'image'"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testEntitiesDatasourceV2ConfigWithFilter(filter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameEntities, "id"),
					resource.TestCheckResourceAttr(datasourceNameEntities, "entities.#", strconv.Itoa(1)),
				),
			},
		},
	})
}

func testEntitiesDatasourceV2Config() string {
	return `
		data "nutanix_iam_entities_v2" "test" {}
	`
}

func testEntitiesDatasourceV2ConfigWithLimit(limit int) string {
	return fmt.Sprintf(`
		data "nutanix_iam_entities_v2" "test" {
			limit = %d
		}
	`, limit)
}

func testEntitiesDatasourceV2ConfigWithFilter(filter string) string {
	return fmt.Sprintf(`
		data "nutanix_iam_entities_v2" "test" {
			filter = "%s"
		}
	`, filter)
}
