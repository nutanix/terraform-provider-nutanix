package clustersv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceNameSnmpTrap = "data.nutanix_snmp_trap_v2.test"

// TestAccV2NutanixSnmpTrapDataSource_Basic fetches an existing SNMP trap by
// UUID. The trap UUID is provided via test_config_v2.json since trap creation
// is not exposed as a Terraform resource at this time.
func TestAccV2NutanixSnmpTrapDataSource_Basic(t *testing.T) {
	octet := acctest.RandIntRange(100, 200)
	trapIP := fmt.Sprintf("10.0.0.%d", octet)
	port := acctest.RandIntRange(30000, 39000)
	community := "public"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testSnmpTrapV2Config(trapIP, community, port, "UDP") + testSnmpTrapDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameSnmpTrap, "ext_id"),
					resource.TestCheckResourceAttrSet(dataSourceNameSnmpTrap, "cluster_ext_id"),
					resource.TestCheckResourceAttr(dataSourceNameSnmpTrap, "version", "V2"),
					resource.TestCheckResourceAttr(dataSourceNameSnmpTrap, "address.0.ipv4.0.value", trapIP),
					resource.TestCheckResourceAttr(dataSourceNameSnmpTrap, "community_string", community),
					resource.TestCheckResourceAttr(dataSourceNameSnmpTrap, "protocol", "UDP"),
					resource.TestCheckResourceAttr(dataSourceNameSnmpTrap, "port", fmt.Sprintf("%d", port)),
				),
			},
		},
	})
}

func testSnmpTrapDatasourceConfig() string {
	return `
data "nutanix_snmp_trap_v2" "test" {
  cluster_ext_id = local.clusterExtID
  ext_id         = nutanix_snmp_trap_v2.test.id
}
`
}
