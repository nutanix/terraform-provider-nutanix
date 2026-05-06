package volumesv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceVolumeGroupMetadata = "data.nutanix_volume_group_metadata_v2.test"

func TestAccV2NutanixVolumeGroupMetadataDataSource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("terraform-test-volume-group-metadata-%d", r)
	desc := "terraform test volume group metadata description"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeGroupMetadataDataSourceConfig(name, desc),
				Check: resource.ComposeTestCheckFunc(
					testAndCheckComputedValues(dataSourceVolumeGroupMetadata),
				),
			},
		},
	})
}

func testAccVolumeGroupMetadataDataSourceConfig(name, desc string) string {
	return testAccVolumeGroupResourceConfig(name, desc) + `
		data "nutanix_volume_group_metadata_v2" "test" {
			volume_group_ext_id = resource.nutanix_volume_group_v2.test.id
			depends_on = [resource.nutanix_volume_group_v2.test]
		}
	`
}
