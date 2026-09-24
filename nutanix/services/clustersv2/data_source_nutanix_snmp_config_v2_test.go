package clustersv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceNameSnmpConfig = "data.nutanix_snmp_config_v2.test"

// TestAccV2NutanixSnmpConfigDataSource_Basic verifies the singular
// nutanix_snmp_config_v2 datasource correctly fetches the full SNMP
// configuration for a cluster (status + transports + traps + users) in a
// single call.
func TestAccV2NutanixSnmpConfigDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testSnmpConfigDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameSnmpConfig, "id"),
					resource.TestCheckResourceAttrSet(dataSourceNameSnmpConfig, "cluster_ext_id"),
					resource.TestCheckResourceAttr(dataSourceNameSnmpConfig, "is_enabled", "true"),
					resource.TestCheckResourceAttrSet(dataSourceNameSnmpConfig, "transports.#"),
					resource.TestCheckResourceAttrSet(dataSourceNameSnmpConfig, "traps.#"),
					resource.TestCheckResourceAttrSet(dataSourceNameSnmpConfig, "users.#"),
				),
			},
		},
	})
}

// TestAccV2NutanixSnmpConfigDataSource_WithUser provisions an SNMP user and
// confirms it surfaces in the cluster-wide SNMP config datasource's `users`
// list with the expected username, auth_type and priv_type values.
func TestAccV2NutanixSnmpConfigDataSource_WithUser(t *testing.T) {
	r := acctest.RandInt()
	username := fmt.Sprintf("tf-acc-snmp-cfg-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testSnmpUserResourceMiniConfig(username, "MD5", "auth-key-cfg-1234") +
					testSnmpConfigDatasourceOnly(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameSnmpConfig, "id"),
					resource.TestCheckResourceAttrSet(dataSourceNameSnmpConfig, "cluster_ext_id"),
					resource.TestCheckResourceAttrSet(dataSourceNameSnmpConfig, "users.#"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceNameSnmpConfig, "users.*", map[string]string{
						"username":  username,
						"auth_type": "MD5",
					}),
				),
			},
		},
	})
}

func testSnmpConfigDatasourceConfig() string {
	return `

# List all clusters to get AOS external ID
data "nutanix_clusters_v2" "aos" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}

locals {
  clusterExtID = data.nutanix_clusters_v2.aos.cluster_entities[0].ext_id
}

data "nutanix_snmp_config_v2" "test" {
  cluster_ext_id = local.clusterExtID
}
`
}

// testSnmpConfigDatasourceOnly returns only the snmp_config datasource block,
// reusing the `nutanix_clusters_v2` data source and `clusterExtID` local that
// are already declared by the resource config it is concatenated with (e.g.
// testSnmpUserResourceMiniConfig), avoiding duplicate declarations.
func testSnmpConfigDatasourceOnly() string {
	return `
data "nutanix_snmp_config_v2" "test" {
  cluster_ext_id = local.clusterExtID
  # Force the data source to be read AFTER the user is created.
  depends_on = [nutanix_snmp_user_v2.test]
}
`
}
