package prismv2_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameUnregisterCluster = "nutanix_pc_unregistration_v2.test"

func TestAccV2NutanixUnregisterClusterResource_Unregister_PC_PC(t *testing.T) {
	if testVars.Prism.Unregister.PcIP == "" {
		t.Skip("Skipping test as it requires PC to unregister not provided")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		ProtoV5ProviderFactories: acc.TestAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUnregisterClusterResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameUnregisterCluster, "pc_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameUnregisterCluster, "ext_id"),
				),
			},
		},
	})
}

func testAccUnregisterClusterResourceConfig() string {
	username, password := getBasicAuthForAPINonSupportedTests()
	port, _ := strconv.Atoi(os.Getenv("NUTANIX_PORT"))
	insecure, _ := strconv.ParseBool(os.Getenv("NUTANIX_INSECURE"))
	remotePCProviderConfig := fmt.Sprintf(`
provider "nutanix-2" {
  username = "%[1]s"
  password = "%[2]s"
  endpoint = "%[3]s"
  insecure = %[4]t
  port     = %[5]d
  api_key  = ""  # Force basic auth against the remote PC
}
`, username, password, testVars.Prism.Unregister.PcIP, insecure, port)

	return fmt.Sprintf(`
%[2]s

data "nutanix_clusters_v2" "cls" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'PRISM_CENTRAL')"
}

# Fetch the remote PC ext_id from the remote endpoint (nutanix-2).
data "nutanix_pcs_v2" "remote" {
  provider = nutanix-2
}

data "nutanix_pc_v2" "remote" {
  provider = nutanix-2
  ext_id   = data.nutanix_pcs_v2.remote.pcs.0.ext_id
}

locals {
  pcExtID      = data.nutanix_clusters_v2.cls.cluster_entities.0.ext_id
  config       = jsondecode(file("%[1]s"))
  unregister   = local.config.prism.unregister
  remotePcExtID = data.nutanix_pc_v2.remote.ext_id
}

resource "nutanix_pc_registration_v2" "remote_pc" {
  pc_ext_id = local.pcExtID

  remote_cluster {
    domain_manager_remote_cluster_spec {
      cloud_type = "ONPREM_CLOUD"

      remote_cluster {
        address {
          ipv4 {
            value = local.unregister.pc_ip
          }
        }

        credentials {
          authentication {
            username = local.config.username_for_test
            password = local.config.password_for_test
          }
        }
      }
    }
  }
}

resource "nutanix_pc_unregistration_v2" "test" {
  pc_ext_id  = local.pcExtID       # local pc ext id
  ext_id     = local.remotePcExtID # remote pc ext id (fetched via nutanix-2)
  depends_on = [nutanix_pc_registration_v2.remote_pc]
}
`, filepath, remotePCProviderConfig)
}
