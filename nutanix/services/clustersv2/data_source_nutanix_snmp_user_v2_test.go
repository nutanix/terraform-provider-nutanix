package clustersv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceNameSnmpUser = "data.nutanix_snmp_user_v2.test"

// TestAccV2NutanixSnmpUserDataSource_Basic verifies the singular SNMP user
// datasource correctly fetches a user by its UUID and surfaces every attribute
// in state.
func TestAccV2NutanixSnmpUserDataSource_Basic(t *testing.T) {
	r := acctest.RandInt()
	username := fmt.Sprintf("tf-acc-snmp-user-ds-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testSnmpUserResourceMiniConfig(username, "MD5", "auth-key-ds-1234") +
					testSnmpUserDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameSnmpUser, "ext_id"),
					resource.TestCheckResourceAttrSet(dataSourceNameSnmpUser, "cluster_ext_id"),
					resource.TestCheckResourceAttr(dataSourceNameSnmpUser, "username", username),
					resource.TestCheckResourceAttr(dataSourceNameSnmpUser, "auth_type", "MD5"),
				),
			},
		},
	})
}

func testSnmpUserDatasourceConfig() string {
	return `
data "nutanix_snmp_user_v2" "test" {
  cluster_ext_id = nutanix_snmp_user_v2.test.cluster_ext_id
  ext_id         = nutanix_snmp_user_v2.test.ext_id
  depends_on     = [nutanix_snmp_user_v2.test]
}
`
}
