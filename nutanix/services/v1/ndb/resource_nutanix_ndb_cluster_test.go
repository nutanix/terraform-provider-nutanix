package ndb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNDBCluster = "nutanix_ndb_cluster.acctest-managed"

func TestAccEra_Clusterbasic(t *testing.T) {
	r := acc.RandIntBetween(25, 35)
	name := fmt.Sprintf("testcluster-%d", r)
	updatedName := fmt.Sprintf("testcluster-updated-%d", r)
	desc := "this is cluster desc"
	updatedDesc := "updated description for cluster"
	storageContainer := acc.TestVars.NDB.RegisterClusterInfo.StorageContainer
	clusterIP := acc.TestVars.NDB.RegisterClusterInfo.ClusterIP
	username := acc.TestVars.NDB.RegisterClusterInfo.Username
	password := acc.TestVars.NDB.RegisterClusterInfo.Password
	staticIP := acc.TestVars.NDB.RegisterClusterInfo.StaticIP
	subnetMask := acc.TestVars.NDB.RegisterClusterInfo.SubnetMask
	gateway := acc.TestVars.NDB.RegisterClusterInfo.Gateway
	dns := acc.TestVars.NDB.RegisterClusterInfo.DNS
	ntp := acc.TestVars.NDB.RegisterClusterInfo.NTP
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccEraPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEraClusterConfig(name, desc, clusterIP, username, password, staticIP, subnetMask, gateway, dns, ntp, storageContainer),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNDBCluster, "name", name),
					resource.TestCheckResourceAttr(resourceNDBCluster, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNDBCluster, "unique_name"),
					resource.TestCheckResourceAttr(resourceNDBCluster, "cloud_type", "NTNX"),
					resource.TestCheckResourceAttr(resourceNDBCluster, "status", "UP"),
					resource.TestCheckResourceAttr(resourceNDBCluster, "healthy", "true"),
					resource.TestCheckResourceAttrSet(resourceNDBCluster, "properties.#"),
					resource.TestCheckResourceAttr(resourceNDBCluster, "hypervisor_type", "AHV"),
					resource.TestCheckResourceAttr(resourceNDBCluster, "version", "v2"),
				),
			},
			{
				Config: testAccEraClusterConfig(updatedName, updatedDesc, clusterIP, username, password, staticIP, subnetMask, gateway, dns, ntp, storageContainer),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNDBCluster, "name", updatedName),
					resource.TestCheckResourceAttr(resourceNDBCluster, "description", updatedDesc),
					resource.TestCheckResourceAttrSet(resourceNDBCluster, "unique_name"),
					resource.TestCheckResourceAttr(resourceNDBCluster, "cloud_type", "NTNX"),
					resource.TestCheckResourceAttr(resourceNDBCluster, "status", "UP"),
					resource.TestCheckResourceAttr(resourceNDBCluster, "healthy", "true"),
					resource.TestCheckResourceAttrSet(resourceNDBCluster, "properties.#"),
					resource.TestCheckResourceAttr(resourceNDBCluster, "hypervisor_type", "AHV"),
					resource.TestCheckResourceAttr(resourceNDBCluster, "version", "v2"),
				),
			},
		},
	})
}

func testAccEraClusterConfig(name, desc, cluster, user, pass, static, mask, gateway, dns, ntp, container string) string {
	return fmt.Sprintf(
		`
		resource "nutanix_ndb_cluster" "acctest-managed" {
			name= "%[1]s"
			description = "%[2]s"
			cluster_ip = "%[3]s"
			username= "%[4]s"
			password = "%[5]s"
			storage_container = "%[11]s"
			agent_network_info{
			  dns = "%[9]s"
			  ntp = "%[10]s"
			}
			networks_info{
			  type = "DHCP"
			  network_info{
				  vlan_name = "vlan_static"
				  static_ip = "%[6]s"
				  gateway = "%[8]s"
				  subnet_mask="%[7]s"
			  }
			  access_type = [
				  "PRISM",
				  "DSIP",
				  "DBSERVER"
				]
			}
		  }
		`, name, desc, cluster, user, pass, static, mask, gateway, dns, ntp, container,
	)
}
