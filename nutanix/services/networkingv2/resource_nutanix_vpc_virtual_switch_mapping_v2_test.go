package networkingv2_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	networkingConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v17/models/networking/v4/config"
	networkingMappingReq "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v17/models/networking/v4/request/vpcvirtualswitchmappings"
	import4 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v17/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const resourceNameVpcVsMapping = "nutanix_vpc_virtual_switch_mapping_v2.test"

// sentinelWrongVsUUID is a syntactically valid but non-existent virtual switch
// UUID, kept verbatim from the Ansible negative task.
const sentinelWrongVsUUID = "04595a7b-010c-4a89-58dd-060c94e9a1f0"

func TestAccV2NutanixVirtualSwitchMapping_BasicCreateAndUpdate(t *testing.T) {
	rName := acctest.RandString(12)
	vsName := fmt.Sprintf("vs_acc_%s_vpc_map", rName)
	catName := fmt.Sprintf("acc%s", rName)
	// Randomize the third octet of the host bridge subnet so reruns (and any
	// leftover virtual switches from previously-failed runs) do not collide on
	// the host IP. The same value is reused across all steps so the virtual
	// switch is not recreated between them.
	subnetID := acctest.RandIntRange(10, 250)

	// initialMappings holds the global East-West mapping config captured before
	// this test sets any mapping. The "set" API replaces the entire config and
	// has no per-entry delete, so the mappings are restored at the end to unpin
	// the virtual switch this test creates (otherwise it cannot be destroyed --
	// the API reports it as "used for Advanced Networking").
	var initialMappings []networkingConfig.VpcVirtualSwitchMapping

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualSwitchDestroy,
		Steps: []resource.TestStep{
			// Step 1: create only the prerequisites (virtual switch + categories)
			// and snapshot the mappings that exist before we set ours.
			{
				Config: testAccVpcVirtualSwitchMappingPrereqs(vsName, catName, subnetID),
				Check:  captureInitialVpcMappings(&initialMappings),
			},
			// Step 2: create the mapping with all traffic permitted + both categories.
			{
				Config: testAccVpcVirtualSwitchMappingConfig(
					vsName, catName, true,
					"[nutanix_category_v2.cat1.id, nutanix_category_v2.cat2.id]",
					subnetID,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVpcVsMapping, "mappings.#", "1"),
					resource.TestCheckResourceAttrSet(resourceNameVpcVsMapping, "mappings.0.virtual_switch_uuid"),
					resource.TestCheckResourceAttrPair(
						resourceNameVpcVsMapping, "mappings.0.virtual_switch_uuid",
						"nutanix_virtual_switch_v2.test", "id",
					),
					resource.TestCheckResourceAttr(resourceNameVpcVsMapping, "mappings.0.cluster_uuids.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVpcVsMapping, "mappings.0.is_all_traffic_permitted", "true"),
					resource.TestCheckResourceAttr(resourceNameVpcVsMapping, "mappings.0.metadata.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVpcVsMapping, "mappings.0.metadata.0.category_ids.#", "2"),
				),
			},
			// Step 3: update -> recreate (ForceNew) with one category and traffic off.
			{
				Config: testAccVpcVirtualSwitchMappingConfig(
					vsName, catName, false,
					"[nutanix_category_v2.cat1.id]",
					subnetID,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameVpcVsMapping, "mappings.#", "1"),
					resource.TestCheckResourceAttrSet(resourceNameVpcVsMapping, "mappings.0.virtual_switch_uuid"),
					resource.TestCheckResourceAttr(resourceNameVpcVsMapping, "mappings.0.cluster_uuids.#", "1"),
					resource.TestCheckResourceAttr(resourceNameVpcVsMapping, "mappings.0.is_all_traffic_permitted", "false"),
					resource.TestCheckResourceAttr(resourceNameVpcVsMapping, "mappings.0.metadata.0.category_ids.#", "1"),
				),
			},
			// Step 4: drop the mapping resource and restore the initial mappings so
			// the virtual switch is unpinned before the framework destroys it.
			{
				Config: testAccVpcVirtualSwitchMappingPrereqs(vsName, catName, subnetID),
				Check:  restoreVpcMappings(&initialMappings),
			},
		},
	})
}

func TestAccV2NutanixVirtualSwitchMapping_NegativeWrongVsUuid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcVirtualSwitchMappingWrongVsUUIDConfig(),
				ExpectError: regexp.MustCompile(
					`(?is)(error while creating VPC Virtual Switch Mapping|error waiting for VPC Virtual Switch Mapping|Task Failed|FAILED|not.?found|invalid|VALIDATION_FAILED)`),
			},
		},
	})
}

func TestAccV2NutanixVirtualSwitchMapping_DataSourceList(t *testing.T) {
	rName := acctest.RandString(12)
	vsName := fmt.Sprintf("vs_acc_%s_vpc_map_ds", rName)
	catName := fmt.Sprintf("accds%s", rName)
	datasourceName := "data.nutanix_vpc_virtual_switch_mappings_v2.test"
	subnetID := acctest.RandIntRange(10, 250)

	var initialMappings []networkingConfig.VpcVirtualSwitchMapping

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualSwitchDestroy,
		Steps: []resource.TestStep{
			// Step 1: prerequisites only; snapshot the initial mappings.
			{
				Config: testAccVpcVirtualSwitchMappingPrereqs(vsName, catName, subnetID),
				Check:  captureInitialVpcMappings(&initialMappings),
			},
			// Step 2: create the mapping and read it back via the data source.
			{
				Config: testAccVpcVirtualSwitchMappingsDataSourceConfig(
					vsName, catName,
					"[nutanix_category_v2.cat1.id, nutanix_category_v2.cat2.id]",
					subnetID,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "vpc_virtual_switch_mappings.#"),
					testCheckVpcMappingDataSourceContainsVS(
						datasourceName, "nutanix_virtual_switch_v2.test",
					),
				),
			},
			// Step 3: restore the initial mappings to unpin the virtual switch.
			{
				Config: testAccVpcVirtualSwitchMappingPrereqs(vsName, catName, subnetID),
				Check:  restoreVpcMappings(&initialMappings),
			},
		},
	})
}

// captureInitialVpcMappings lists the current VPC virtual switch mappings via
// the SDK client and stores their writable fields (virtual_switch_uuid /
// cluster_uuids / is_all_traffic_permitted) into dst, to be restored on
// teardown. Run this before the test sets any mapping, because the "set" API
// replaces the entire East-West config.
func captureInitialVpcMappings(dst *[]networkingConfig.VpcVirtualSwitchMapping) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		client := acc.TestAccProvider.Meta().(*conns.Client)
		listReq := networkingMappingReq.ListVpcVirtualSwitchMappingsRequest{}
		resp, err := client.NetworkingAPI.VpcVirtualSwitchMappingsAPI.ListVpcVirtualSwitchMappings(context.Background(), &listReq)
		if err != nil {
			return fmt.Errorf("listing initial VPC virtual switch mappings: %v", err)
		}
		out := make([]networkingConfig.VpcVirtualSwitchMapping, 0)
		if resp != nil && resp.Data != nil {
			for _, m := range resp.Data.GetValue().([]networkingConfig.VpcVirtualSwitchMapping) {
				out = append(out, networkingConfig.VpcVirtualSwitchMapping{
					VirtualSwitchUuid:     m.VirtualSwitchUuid,
					ClusterUuids:          m.ClusterUuids,
					IsAllTrafficPermitted: m.IsAllTrafficPermitted,
				})
			}
		}
		*dst = out
		return nil
	}
}

// restoreVpcMappings re-sets the global East-West config to the captured
// initial state via the SDK client and waits for the task to succeed. This
// removes the mapping this test added, unpinning the virtual switch so the
// framework can destroy it.
func restoreVpcMappings(src *[]networkingConfig.VpcVirtualSwitchMapping) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		client := acc.TestAccProvider.Meta().(*conns.Client)
		ctx := context.Background()

		body := make([]networkingConfig.VpcVirtualSwitchMapping, 0)
		if src != nil {
			body = *src
		}
		req := networkingMappingReq.CreateVpcVirtualSwitchMappingRequest{Body: &body}
		resp, err := client.NetworkingAPI.VpcVirtualSwitchMappingsAPI.CreateVpcVirtualSwitchMapping(ctx, &req)
		if err != nil {
			return fmt.Errorf("restoring VPC virtual switch mappings: %v", err)
		}

		taskRef := resp.Data.GetValue().(import4.TaskReference)
		stateConf := &resource.StateChangeConf{
			Pending: []string{"PENDING", "RUNNING", "QUEUED"},
			Target:  []string{"SUCCEEDED"},
			Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, client.PrismAPI, utils.StringValue(taskRef.ExtId)),
			Timeout: 10 * time.Minute,
		}
		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for VPC virtual switch mapping restore task: %v", err)
		}
		return nil
	}
}

// testCheckVpcMappingDataSourceContainsVS verifies the listed mappings contain
// an entry whose virtual_switch_uuid matches the given virtual switch's ID. The
// data source returns the full set of mappings, so the index is discovered by
// iterating.
func testCheckVpcMappingDataSourceContainsVS(dsName, vsResourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		vsRes, ok := s.RootModule().Resources[vsResourceName]
		if !ok {
			return fmt.Errorf("virtual switch resource %q not found in state", vsResourceName)
		}
		vsID := vsRes.Primary.ID
		if vsID == "" {
			return fmt.Errorf("virtual switch resource %q has an empty ID", vsResourceName)
		}

		dsRes, ok := s.RootModule().Resources[dsName]
		if !ok {
			return fmt.Errorf("data source %q not found in state", dsName)
		}
		attrs := dsRes.Primary.Attributes

		count, err := strconv.Atoi(attrs["vpc_virtual_switch_mappings.#"])
		if err != nil {
			return fmt.Errorf("could not parse vpc_virtual_switch_mappings.# from data source %q: %v", dsName, err)
		}

		for i := 0; i < count; i++ {
			key := fmt.Sprintf("vpc_virtual_switch_mappings.%d.virtual_switch_uuid", i)
			if attrs[key] == vsID {
				return nil
			}
		}
		return fmt.Errorf(
			"virtual switch UUID %s not found among the %d mapping(s) returned by data source %q",
			vsID, count, dsName)
	}
}

// testAccVpcVirtualSwitchMappingPrereqs emits the shared prerequisites. subnetID
// is the randomized third octet of the host bridge subnet (192.168.<subnetID>.0/24)
// used to keep the host IP unique across runs.
func testAccVpcVirtualSwitchMappingPrereqs(vsName, catName string, subnetID int) string {
	return fmt.Sprintf(`
data "nutanix_clusters_v2" "test" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}
# Only fetch hosts that belong to the selected cluster (server-side OData filter
# on cluster/uuid), so host_entities[0] is guaranteed to be in cluster_ext_id.
# The backend rejects a VS whose host UUID is not part of the target cluster
# ("host config list contains unknown host UUIDs ... in cluster ...").
data "nutanix_hosts_v2" "test" {
  filter = "cluster/uuid eq '${data.nutanix_clusters_v2.test.cluster_entities[0].ext_id}'"
}

locals {
  cluster_ext_id = data.nutanix_clusters_v2.test.cluster_entities[0].ext_id
  host_ext_id    = data.nutanix_hosts_v2.test.host_entities[0].ext_id
}

resource "nutanix_category_v2" "cat1" {
  key         = "%[2]s_key1"
  value       = "%[2]s_value1"
  description = "Category 1 for VPC virtual switch mapping tests"
}

resource "nutanix_category_v2" "cat2" {
  key         = "%[2]s_key2"
  value       = "%[2]s_value2"
  description = "Category 2 for VPC virtual switch mapping tests"
}

resource "nutanix_virtual_switch_v2" "test" {
  name      = "%[1]s"
  bond_mode = "NONE"
  clusters {
    ext_id          = local.cluster_ext_id
    vlan_identifier = 100
    gateway_ip_address {
      value         = "192.168.%[3]d.1"
      prefix_length = 32
    }
    hosts {
      ext_id = local.host_ext_id
      ip_address {
        ip {
          value         = "192.168.%[3]d.10"
          prefix_length = 32
        }
        prefix_length = 24
      }
    }
  }
}
`, vsName, catName, subnetID)
}

// testAccVpcVirtualSwitchMappingConfig builds the prerequisites plus the
// mapping resource. categoryIDsHCL is the raw HCL list expression for
// metadata.category_ids, letting one builder serve both the create and update
// step.
func testAccVpcVirtualSwitchMappingConfig(vsName, catName string, isAllTrafficPermitted bool, categoryIDsHCL string, subnetID int) string {
	return testAccVpcVirtualSwitchMappingPrereqs(vsName, catName, subnetID) + fmt.Sprintf(`
resource "nutanix_vpc_virtual_switch_mapping_v2" "test" {
  mappings {
    virtual_switch_uuid      = nutanix_virtual_switch_v2.test.id
    cluster_uuids            = [local.cluster_ext_id]
    is_all_traffic_permitted = %[1]t

    metadata {
      category_ids = %[2]s
    }
  }
}
`, isAllTrafficPermitted, categoryIDsHCL)
}

func testAccVpcVirtualSwitchMappingWrongVsUUIDConfig() string {
	return fmt.Sprintf(`
data "nutanix_clusters_v2" "test" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}

locals {
  cluster_ext_id = data.nutanix_clusters_v2.test.cluster_entities[0].ext_id
}

resource "nutanix_vpc_virtual_switch_mapping_v2" "test" {
  mappings {
    virtual_switch_uuid      = "%[1]s"
    cluster_uuids            = [local.cluster_ext_id]
    is_all_traffic_permitted = true
  }
}
`, sentinelWrongVsUUID)
}

func testAccVpcVirtualSwitchMappingsDataSourceConfig(vsName, catName, categoryIDsHCL string, subnetID int) string {
	return testAccVpcVirtualSwitchMappingConfig(vsName, catName, true, categoryIDsHCL, subnetID) + `
data "nutanix_vpc_virtual_switch_mappings_v2" "test" {
  depends_on = [nutanix_vpc_virtual_switch_mapping_v2.test]
}
`
}
