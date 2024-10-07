package dataprotectionv2_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameRecoveryPointReplicate = "nutanix_recovery_point_replicate_v2.test"

func TestAccNutanixRecoveryPointReplicateV2Resource_basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("terraform-test-recovery-point-%d", r)

	// End time is two week later
	expirationTime := time.Now().Add(14 * 24 * time.Hour)

	expirationTimeFormatted := expirationTime.UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRecoveryPointReplicateResourceConfig(name, expirationTimeFormatted),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRecoveryPointReplicate, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameRecoveryPointReplicate, "pc_ext_id", testVars.DataProtection.PcExtID),
					resource.TestCheckResourceAttr(resourceNameRecoveryPointReplicate, "cluster_ext_id", testVars.DataProtection.ClusterExtID),
					resource.TestCheckResourceAttrSet(resourceNameRecoveryPointReplicate, "replicated_rp_ext_id"),
				),
			},
		},
	})
}

func testRecoveryPointReplicateResourceConfig(name, expirationTime string) string {
	return testRecoveryPointsResourceConfigWithVolumeGroupRecoveryPoints(name, expirationTime) + fmt.Sprintf(`
	locals{
		config = (jsondecode(file("%[1]s")))
		data_protection = local.config.data_protection			
	}
	resource "nutanix_recovery_point_replicate_v2" "test" {
	  ext_id         = nutanix_recovery_points_v2.test.id
	  cluster_ext_id = "000620a9-8183-2553-1fc3-ac1f6b6029c1"
	  pc_ext_id      = "63bebabf-744c-48ff-a6d7-cb028707f972"
	  depends_on     = [nutanix_recovery_points_v2.test]
	}`, filepath)
}
