package iamv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameEntity = "data.nutanix_iam_entity_v2.test"

func TestAccV2NutanixEntityDatasource_GetEntityById(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testEntityDatasourceV2Config(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceNameEntity, "ext_id", datasourceNameEntities, "entities.0.ext_id"),
					resource.TestCheckResourceAttrPair(datasourceNameEntity, "name", datasourceNameEntities, "entities.0.name"),
					resource.TestCheckResourceAttrPair(datasourceNameEntity, "description", datasourceNameEntities, "entities.0.description"),
					resource.TestCheckResourceAttrPair(datasourceNameEntity, "display_name", datasourceNameEntities, "entities.0.display_name"),
					resource.TestCheckResourceAttrPair(datasourceNameEntity, "client_name", datasourceNameEntities, "entities.0.client_name"),
					resource.TestCheckResourceAttrPair(datasourceNameEntity, "search_url", datasourceNameEntities, "entities.0.search_url"),
					resource.TestCheckResourceAttrPair(datasourceNameEntity, "created_time", datasourceNameEntities, "entities.0.created_time"),
					resource.TestCheckResourceAttrPair(datasourceNameEntity, "last_updated_time", datasourceNameEntities, "entities.0.last_updated_time"),
					resource.TestCheckResourceAttrPair(datasourceNameEntity, "created_by", datasourceNameEntities, "entities.0.created_by"),
					resource.TestCheckResourceAttrPair(datasourceNameEntity, "attribute_list.#", datasourceNameEntities, "entities.0.attribute_list.#"),
					resource.TestCheckResourceAttrPair(datasourceNameEntity, "attribute_list.0.display_name", datasourceNameEntities, "entities.0.attribute_list.0.display_name"),
					resource.TestCheckResourceAttrPair(datasourceNameEntity, "attribute_list.0.supported_operator.#", datasourceNameEntities, "entities.0.attribute_list.0.supported_operator.#"),
					resource.TestCheckResourceAttrPair(datasourceNameEntity, "attribute_list.0.attribute_values.#", datasourceNameEntities, "entities.0.attribute_list.0.attribute_values.#"),
					resource.TestCheckResourceAttrPair(datasourceNameEntity, "is_logical_and_supported_for_attributes", datasourceNameEntities, "entities.0.is_logical_and_supported_for_attributes"),
				),
			},
		},
	})
}

func testEntityDatasourceV2Config(configPath string) string {
	return `

data "nutanix_iam_entities_v2" "test" {
  limit   = 1
}

data "nutanix_iam_entity_v2" "test" {
  ext_id = data.nutanix_iam_entities_v2.test.entities[0].ext_id
}
	`
}
