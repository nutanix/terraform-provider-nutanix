package clustersv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameSnmpUser = "nutanix_snmp_user_v2.test"

// TestAccV2NutanixSnmpUserResource_Basic exercises the full CRUD lifecycle of a
// SNMP user resource and verifies every attribute exposed by the schema is
// either set explicitly or computed in state.
func TestAccV2NutanixSnmpUserResource_MinimumAttributes(t *testing.T) {
	r := acctest.RandInt()
	username := fmt.Sprintf("tf-acc-snmp-user-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixSnmpUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testSnmpUserResourceMiniConfig(username, "MD5", "auth-key-original-1234"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameSnmpUser, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameSnmpUser, "cluster_ext_id"),
					resource.TestCheckResourceAttr(resourceNameSnmpUser, "username", username),
					resource.TestCheckResourceAttr(resourceNameSnmpUser, "auth_type", "MD5"),
				),
			},
			// update
			{
				Config: testSnmpUserResourceMiniConfig(username, "SHA", "auth-key-updated-5678"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameSnmpUser, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameSnmpUser, "username", username),
					resource.TestCheckResourceAttr(resourceNameSnmpUser, "auth_type", "SHA"),
				),
			},
		},
	})
}

func TestAccV2NutanixSnmpUserResource_AllAttributes(t *testing.T) {
	r := acctest.RandInt()
	username := fmt.Sprintf("tf-acc-snmp-user-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixSnmpUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testSnmpUserResourceAllConfig(username, "SHA", "auth-key-original-1234", "AES", "priv-key-original-1234"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameSnmpUser, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameSnmpUser, "cluster_ext_id"),
					resource.TestCheckResourceAttr(resourceNameSnmpUser, "username", username),
					resource.TestCheckResourceAttr(resourceNameSnmpUser, "auth_type", "SHA"),
					resource.TestCheckResourceAttr(resourceNameSnmpUser, "priv_type", "AES"),
				),
			},
			// update
			{
				Config: testSnmpUserResourceAllConfig(username, "MD5", "auth-key-updated-5678", "DES", "priv-key-updated-5678"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameSnmpUser, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameSnmpUser, "username", username),
					resource.TestCheckResourceAttr(resourceNameSnmpUser, "auth_type", "MD5"),
					resource.TestCheckResourceAttr(resourceNameSnmpUser, "priv_type", "DES"),
				),
			},
		},
	})
}

func testSnmpUserResourceMiniConfig(username, authType, authKey string) string {
	return fmt.Sprintf(`
# List all clusters to get AOS external ID
data "nutanix_clusters_v2" "aos" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}

locals {
  clusterExtID = data.nutanix_clusters_v2.aos.cluster_entities[0].ext_id
}

resource "nutanix_snmp_user_v2" "test" {
  cluster_ext_id = local.clusterExtID
  username       = "%s"
  auth_type      = "%s"
  auth_key       = "%s"
}

`, username, authType, authKey)
}

func testSnmpUserResourceAllConfig(username, authType, authKey, privType, privKey string) string {
	return fmt.Sprintf(`

data "nutanix_clusters_v2" "aos" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}

locals {
  clusterExtID = data.nutanix_clusters_v2.aos.cluster_entities[0].ext_id
}

resource "nutanix_snmp_user_v2" "test" {
  cluster_ext_id = local.clusterExtID
  username       = "%s"
  auth_type      = "%s"
  auth_key       = "%s"
  priv_type      = "%s"
  priv_key       = "%s"
}

`, username, authType, authKey, privType, privKey)
}
