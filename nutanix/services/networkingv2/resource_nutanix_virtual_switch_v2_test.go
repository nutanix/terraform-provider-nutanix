package networkingv2_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	networkingVsReq "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/request/virtualswitches"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// File layout:
//   1. Positive test cases    (resource lifecycle + datasource reads)
//   2. Negative test cases    (plan-time + API-side rejection assertions)
//   3. Helpers                (CheckDestroy)
//   4. Config builders        (one Sprintf-style builder per test)
//
// Test plan:
//   Positive:
//     - TestAccV2NutanixVirtualSwitch_Minimum:              Create with minimum fields.
//     - TestAccV2NutanixVirtualSwitch_WithAllAttributes:    Create with the full set of user-configurable attributes.
//     - TestAccV2NutanixVirtualSwitch_FromExistingBridge:   Migrate an existing host-side OVS bridge into a Virtual Switch (requires lab prep).
//     - TestAccV2NutanixVirtualSwitch_ShareUnshareLifecycle: Create with no shares, add a share, swap shares, then remove all shares.
//     - TestAccV2NutanixVirtualSwitch_Basic:                Create + update + import + read-via-datasource.
//     - TestAccV2NutanixVirtualSwitchDatasource_List:       Read list of virtual switches.
//     - TestAccV2NutanixNodeSchedulableStatuses:            Read node schedulable statuses.
//     - TestAccV2NutanixVpcVirtualSwitchMappings:           Read VPC virtual switch mappings.
//   Negative:
//     - TestAccV2NutanixVirtualSwitch_MissingName:                              Plan-time `Missing required argument`.
//     - TestAccV2NutanixVirtualSwitch_MissingBondMode:                          API-side rejection when bondMode is absent.
//     - TestAccV2NutanixVirtualSwitch_MissingClusters:                          API-side rejection when clusters[] is absent.
//     - TestAccV2NutanixVirtualSwitch_FromExistingBridge_AlreadyInUse:          Two-step: first migrate succeeds; second migrate of the same bridge is rejected.
//     - TestAccV2NutanixVirtualSwitch_FromExistingBridge_WrongClusterReference: Migrate with a bogus cluster UUID; task is expected to fail.
//     - TestAccV2NutanixVirtualSwitch_FromExistingBridge_MissingClusterReference: Migrate with no `clusters[0].ext_id`; API is expected to reject.
//
//   Negative scenarios deliberately NOT translated from Ansible:
//     - "ext_id + existing_bridge_name mutually exclusive": the resource's
//       root-level `ext_id` is `Computed: true` only, so users cannot set it
//       from configuration -- the conflict cannot be expressed in TF.
//     - "bond_mode/igmp_spec mutually exclusive with existing_bridge_name":
//       the provider's migrate path currently only logs a `[WARN]` for each
//       ignored field (see warnIgnoredMigrateFields). It does not reject the
//       apply. Adding such a test would require first hardening the provider
//       to error in CustomizeDiff or Create.

// =============================================================================
// 1. Positive test cases
// =============================================================================

// TestAccV2NutanixVirtualSwitch_Minimum: Create with minimum fields.
func TestAccV2NutanixVirtualSwitch_Minimum(t *testing.T) {
	resourceName := "nutanix_virtual_switch_v2.test"
	rName := acctest.RandomWithPrefix("tf-vs-test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualSwitchDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualSwitchMinimumConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "ext_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

// TestAccV2NutanixVirtualSwitch_WithAllAttributes: Create a virtual switch
// populating every user-configurable attribute the schema exposes (description,
// bond_mode, mtu, cluster+host ext_ids, cluster-level vlan_identifier, and the
// igmp_spec block). `is_default` is intentionally omitted because the schema
// marks it Computed-only and Terraform will reject setting it from config.
func TestAccV2NutanixVirtualSwitch_WithAllAttributes(t *testing.T) {
	resourceName := "nutanix_virtual_switch_v2.test"
	rName := acctest.RandomWithPrefix("tf-vs-test")
	description := "Virtual switch created by integration tests with all attributes"
	updatedDescription := "Updated description for virtual switch with all attributes"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualSwitchDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualSwitchWithAllAttributesConfig(rName, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "ext_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"_all"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "bond_mode", "NONE"),
					resource.TestCheckResourceAttr(resourceName, "mtu", "1500"),
					resource.TestCheckResourceAttr(resourceName, "clusters.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "clusters.0.ext_id"),
					resource.TestCheckResourceAttr(resourceName, "clusters.0.vlan_identifier", "0"),
					resource.TestCheckResourceAttr(resourceName, "clusters.0.hosts.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "clusters.0.hosts.0.ext_id"),
					resource.TestCheckResourceAttr(resourceName, "igmp_spec.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "igmp_spec.0.is_snooping_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "igmp_spec.0.snooping_timeout", "300"),
				),
			},
			{
				Config: testAccVirtualSwitchWithAllAttributesUpdatedConfig(rName, updatedDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "ext_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"_updated"),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
					resource.TestCheckResourceAttr(resourceName, "bond_mode", "NONE"),
					resource.TestCheckResourceAttr(resourceName, "mtu", "1500"),
					resource.TestCheckResourceAttr(resourceName, "is_default", "false"),
					resource.TestCheckResourceAttr(resourceName, "clusters.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "clusters.0.ext_id"),
					resource.TestCheckResourceAttr(resourceName, "clusters.0.vlan_identifier", "0"),
					resource.TestCheckResourceAttr(resourceName, "clusters.0.hosts.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "clusters.0.hosts.0.ext_id"),
					resource.TestCheckResourceAttr(resourceName, "igmp_spec.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "igmp_spec.0.is_snooping_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "igmp_spec.0.snooping_timeout", "300"),
					resource.TestCheckResourceAttr(resourceName, "igmp_spec.0.querier_spec.0.is_querier_enabled", "true"),
				),
			},
		},
	})
}

// TestAccV2NutanixVirtualSwitch_FromExistingBridge: Create a virtual switch
// by migrating an existing host-side OVS bridge.
//
// The resource's CreateContext routes to the BridgesAPI `$actions/migrate`
// endpoint when `clusters[].existing_bridge_name` is set, because the
// standard POST /api/networking/v4.3/config/virtual-switches endpoint
// silently ignores any pre-existing bridge hint and auto-allocates a fresh
// brN. Two attributes are asserted to confirm the migrate path actually
// fired:
//   - clusters.0.existing_bridge_name: the create-time-only input, preserved
//     across reads by the overlay in resourceNutanixVirtualSwitchV2Read.
//   - clusters.0.hosts.0.internal_bridge_name: the API's read-back, which
//     must equal the requested bridge name (the whole point of migrate).
//
// The bridge is created at test setup via createBridgeOnPE (SSH + manage_ovs)
// and removed at test teardown via deleteBridgeOnPE, so no lab-side prep is
// required.
func TestAccV2NutanixVirtualSwitch_FromExistingBridge(t *testing.T) {
	bridgeName := createBridgeOnPE(t)
	t.Cleanup(func() { deleteBridgeOnPE(t, bridgeName) })
	resourceName := "nutanix_virtual_switch_v2.test"
	rName := acctest.RandomWithPrefix("tf-vs-test")
	description := "Virtual switch created from existing bridge"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualSwitchDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualSwitchFromExistingBridgeConfig(rName, description, bridgeName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "ext_id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"_existing"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "clusters.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "clusters.0.ext_id"),
					resource.TestCheckResourceAttr(resourceName, "clusters.0.existing_bridge_name", bridgeName),
					resource.TestCheckResourceAttr(resourceName, "clusters.0.hosts.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "clusters.0.hosts.0.ext_id"),
					resource.TestCheckResourceAttr(resourceName, "clusters.0.hosts.0.internal_bridge_name", bridgeName),
				),
			},
		},
	})
}

// TestAccV2NutanixVirtualSwitch_ShareUnshareLifecycle exercises the full
// native share/unshare lifecycle now handled directly by the
// nutanix_virtual_switch_v2 resource (the standalone
// nutanix_share_virtual_switch_v2 resource was removed in favor of the
// shared_with_projects attribute):
//
//	Step 1: create the Virtual Switch with no shares          -> shared_with_projects empty.
//	Step 2: add a share with project 1                        -> ShareVirtualSwitchById fires once.
//	Step 3: swap the share to project 2                       -> one Unshare (proj1) + one Share (proj2).
//	Step 4: clear shared_with_projects entirely               -> UnshareVirtualSwitchById fires for proj2.
//
// Two nutanix_project resources are created so the share/unshare endpoints
// have real project UUIDs to operate on. CheckDestroy guarantees the VS is
// gone (and thus unshared) at the end of the test.
func TestAccV2NutanixVirtualSwitch_ShareUnshareLifecycle(t *testing.T) {
	resourceName := "nutanix_virtual_switch_v2.test_vs"
	rName := acctest.RandomWithPrefix("tf-vs-share")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualSwitchDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create Virtual Switch (No Shares).
			{
				Config: testAccVirtualSwitchShareStep1Config(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "ext_id"),
					resource.TestCheckResourceAttr(resourceName, "name", "acc-test-vswitch-"+rName),
					resource.TestCheckResourceAttr(resourceName, "shared_with_projects.#", "0"),
				),
			},
			// Step 2: Update to Add Share (share with project 1).
			{
				Config: testAccVirtualSwitchShareStep2Config(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "shared_with_projects.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "shared_with_projects.0", "nutanix_project_v2.test_project_1", "ext_id"),
				),
			},
			// Step 3: Update to Swap Shares (share project 2, unshare project 1).
			{
				Config: testAccVirtualSwitchShareStep3Config(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "shared_with_projects.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "shared_with_projects.0", "nutanix_project_v2.test_project_2", "ext_id"),
				),
			},
			// Step 4: Update to Remove All Shares.
			{
				Config: testAccVirtualSwitchShareStep4Config(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "shared_with_projects.#", "0"),
				),
			},
		},
	})
}

func TestAccV2NutanixVirtualSwitchDatasource_List(t *testing.T) {
	datasourceName := "data.nutanix_virtual_switches_v2.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualSwitchesListConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "virtual_switches.#"),
				),
			},
		},
	})
}

// TestAccV2NutanixVirtualSwitchDatasource_GetByExtId reads a freshly-created
// Virtual Switch through the single-object datasource (by ext_id) and asserts
// every attribute round-trips, mirroring the Ansible "Fetch virtual switch
// using external ID" info-module test.
func TestAccV2NutanixVirtualSwitchDatasource_GetByExtId(t *testing.T) {
	datasourceName := "data.nutanix_virtual_switch_v2.test"
	rName := acctest.RandomWithPrefix("tf-vs-test")
	description := "Updated description for virtual switch with all attributes"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualSwitchDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualSwitchInfoByExtIdConfig(rName, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "ext_id"),
					resource.TestCheckResourceAttr(datasourceName, "name", rName+"_updated"),
					resource.TestCheckResourceAttr(datasourceName, "description", description),
					resource.TestCheckResourceAttr(datasourceName, "bond_mode", "NONE"),
					resource.TestCheckResourceAttr(datasourceName, "mtu", "1500"),
					resource.TestCheckResourceAttr(datasourceName, "is_default", "false"),
					resource.TestCheckResourceAttr(datasourceName, "clusters.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceName, "clusters.0.ext_id"),
					resource.TestCheckResourceAttr(datasourceName, "clusters.0.vlan_identifier", "0"),
					resource.TestCheckResourceAttr(datasourceName, "clusters.0.hosts.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "igmp_spec.0.is_snooping_enabled", "true"),
					resource.TestCheckResourceAttr(datasourceName, "igmp_spec.0.snooping_timeout", "300"),
					resource.TestCheckResourceAttr(datasourceName, "igmp_spec.0.querier_spec.0.is_querier_enabled", "true"),
				),
			},
		},
	})
}

// TestAccV2NutanixVirtualSwitchDatasource_GetByWrongExtId expects the single
// datasource read to fail when given a syntactically-valid but non-existent
// ext_id, mirroring the Ansible "Fetch virtual switch using wrong external ID"
// negative info-module test.
func TestAccV2NutanixVirtualSwitchDatasource_GetByWrongExtId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccVirtualSwitchInfoByWrongExtIdConfig(),
				ExpectError: regexp.MustCompile(`(?is)(error|not.?found|invalid|VALIDATION_FAILED)`),
			},
		},
	})
}

// TestAccV2NutanixVirtualSwitchDatasource_ListWithFilters exercises the list
// datasource with a $filter and with a $limit, mirroring the Ansible
// "List virtual switches with filter" and "with limit" info-module tests.
func TestAccV2NutanixVirtualSwitchDatasource_ListWithFilters(t *testing.T) {
	datasourceName := "data.nutanix_virtual_switches_v2.test"
	rName := acctest.RandomWithPrefix("tf-vs-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualSwitchDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualSwitchListWithFilterConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "virtual_switches.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "virtual_switches.0.name", rName),
				),
			},
			{
				Config: testAccVirtualSwitchListWithLimitConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "virtual_switches.#", "1"),
				),
			},
		},
	})
}

func TestAccV2NutanixNodeSchedulableStatuses(t *testing.T) {
	datasourceName := "data.nutanix_node_schedulable_statuses_v2.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNodeSchedulableStatusesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "node_schedulable_statuses.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixVpcVirtualSwitchMappings(t *testing.T) {
	datasourceName := "data.nutanix_vpc_virtual_switch_mappings_v2.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcVirtualSwitchMappingsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "vpc_virtual_switch_mappings.#"),
				),
			},
		},
	})
}

// =============================================================================
// 2. Negative test cases
// =============================================================================

// TestAccV2NutanixVirtualSwitch_MissingName: `name` is Required in the schema.
// Omitting it must fail at plan time with Terraform core's "Missing required
// argument" diagnostic before any API call is made.
func TestAccV2NutanixVirtualSwitch_MissingName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccVirtualSwitchMissingNameConfig(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

// TestAccV2NutanixVirtualSwitch_MissingBondMode: `bond_mode` is Optional+Computed
// in the Terraform schema, so the plan succeeds and the provider posts to Prism
// without `bondMode`. The API is expected to reject the create request because
// bondMode is a required field on the server side.
//
// The regex is intentionally permissive: it matches the field name
// (`bondMode` / `bond_mode` / `bond mode`) or a generic validation/required
// diagnostic. The config supplies a valid `name` and a `clusters` block, so
// the only field missing from the request is `bondMode` — any "missing"
// error must therefore be about bondMode, not something else.
//
// If Prism's exact wording is known, tighten this regex to match it.
func TestAccV2NutanixVirtualSwitch_MissingBondMode(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-vs-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualSwitchDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccVirtualSwitchMissingBondModeConfig(rName),
				ExpectError: regexp.MustCompile(`(?is)(bond.?mode|VALIDATION_FAILED|InvalidArgument|required)`),
			},
		},
	})
}

// TestAccV2NutanixVirtualSwitch_MissingClusters: `clusters` is Optional+Computed
// in the schema, so the plan succeeds and the provider posts to Prism without a
// `clusters` array. The API is expected to reject the create because a virtual
// switch needs at least one cluster.
//
// The regex matches the field name (`cluster` / `clusters`) or a generic
// validation/required diagnostic. The config provides `name` and `bond_mode`,
// so the only field absent from the request is `clusters` — any "missing"/
// "required" wording in the server response is therefore about `clusters` and
// not some other field.
//
// Tighten this regex once Prism's exact wording is known.
func TestAccV2NutanixVirtualSwitch_MissingClusters(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-vs-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualSwitchDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccVirtualSwitchMissingClustersConfig(rName),
				ExpectError: regexp.MustCompile(`(?is)(clusters?|VALIDATION_FAILED|InvalidArgument|required)`),
			},
		},
	})
}

// TestAccV2NutanixVirtualSwitch_FromExistingBridge_AlreadyInUse: a bridge can
// only back one Virtual Switch at a time, so the second migrate attempt for
// the same bridge must be rejected. Modelled as two steps:
//
//	Step 1: declare a single VS that migrates the test-created bridge.
//	        Apply must succeed -- the bridge is now owned by the first VS.
//	Step 2: keep the first VS and add a second VS pointing at the same
//	        bridge name. Apply must fail; the framework will roll back to
//	        the prior state (first VS still in state) and CheckDestroy
//	        cleans up the first VS at the end of the test.
//
// The bridge is created at test setup via createBridgeOnPE (SSH + manage_ovs)
// and removed at test teardown via deleteBridgeOnPE.
//
// Regex is intentionally permissive: Prism's exact wording for "bridge
// already attached" / "in use" / "task failed" is not pinned. The provider
// wraps API errors with `error while migrating bridge ...` (synchronous
// rejection) or `error waiting for Virtual Switch (...) to create`
// (asynchronous task failure); the pattern catches both wrappers plus
// likely server-side wording.
func TestAccV2NutanixVirtualSwitch_FromExistingBridge_AlreadyInUse(t *testing.T) {
	bridgeName := createBridgeOnPE(t)
	t.Cleanup(func() { deleteBridgeOnPE(t, bridgeName) })
	rName := acctest.RandomWithPrefix("tf-vs-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualSwitchDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualSwitchFromExistingBridgeAlreadyInUseStep1Config(rName, bridgeName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("nutanix_virtual_switch_v2.first", "ext_id"),
					resource.TestCheckResourceAttr("nutanix_virtual_switch_v2.first", "clusters.0.existing_bridge_name", bridgeName),
				),
			},
			{
				Config: testAccVirtualSwitchFromExistingBridgeAlreadyInUseStep2Config(rName, bridgeName),
				ExpectError: regexp.MustCompile(
					`(?is)(error while migrating bridge|error waiting for Virtual Switch|already|in use|attach|owned|migrated|VALIDATION_FAILED|task.*fail)`),
			},
		},
	})
}

// TestAccV2NutanixVirtualSwitch_FromExistingBridge_WrongClusterReference:
// route to the migrate endpoint with a syntactically-valid but
// non-existent cluster UUID. The synchronous migrate POST may accept the
// request, but the resulting task fails because Prism cannot resolve the
// cluster reference.
//
// We deliberately do NOT use a real cluster UUID + real bridge here: a
// "wrong cluster" scenario should fail before the bridge state is
// touched, so there is nothing to clean up after the test. CheckDestroy
// is still set as a safety net.
//
// The UUID used (`04595a7b-010c-4a89-58dd-060c94e9a1f0`) is the same
// sentinel the Ansible negative test uses, kept verbatim so failure
// signatures stay comparable across the two test suites.
func TestAccV2NutanixVirtualSwitch_FromExistingBridge_WrongClusterReference(t *testing.T) {
	bridgeName := "br99"
	rName := acctest.RandomWithPrefix("tf-vs-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualSwitchDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualSwitchFromExistingBridgeWrongClusterRefConfig(rName, bridgeName),
				ExpectError: regexp.MustCompile(
					`(?is)(error while migrating bridge|error waiting for Virtual Switch|cluster|not.?found|invalid|VALIDATION_FAILED|task.*fail)`),
			},
		},
	})
}

// TestAccV2NutanixVirtualSwitch_FromExistingBridge_MissingClusterReference:
// declare a `clusters` block that sets `existing_bridge_name` but omits
// `ext_id`. The provider's detectMigrationFromExistingBridge returns
// `matched=true` with `clusterRef=""`, so the migrate request body is
// constructed WITHOUT a `clusterReference` field (see
// resourceNutanixVirtualSwitchV2CreateFromBridge: `if clusterRef != ""`).
// Prism is expected to reject either synchronously (validation) or
// asynchronously (task failure).
func TestAccV2NutanixVirtualSwitch_FromExistingBridge_MissingClusterReference(t *testing.T) {
	bridgeName := "br99"
	rName := acctest.RandomWithPrefix("tf-vs-test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualSwitchDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualSwitchFromExistingBridgeMissingClusterRefConfig(rName, bridgeName),
				ExpectError: regexp.MustCompile(
					`(?is)(error while migrating bridge|error waiting for Virtual Switch|cluster.?reference|required|missing|VALIDATION_FAILED|task.*fail)`),
			},
		},
	})
}

// =============================================================================
// 3. Helpers
// =============================================================================

// createBridgeOnPE provisions the next free OVS bridge (e.g. br5) on the PE
// cluster and returns its name.
//
// Implementation note: the discover → create → verify sequence runs as a LOCAL
// bash script that calls sshpass three times (one SSH connection per step),
// mirroring the approach used in temp/main.tf. Running a single multi-line
// script remotely with `set -euo pipefail` is unreliable because
// `source /etc/profile` on the CVM can return non-zero in a non-login shell,
// which kills the whole session under `set -e` before manage_ovs is reached.
// Separate SSH calls avoid this: `source /etc/profile; manage_ovs ...` runs
// inside its own remote shell that is NOT subject to the local script's set -e.
//
// Credentials: pe_ip / ssh_pe_username / ssh_pe_password in test_config_v2.json.
// The test is skipped if any credential is absent.
// sshpass(1) must be installed on the test machine:
//
//	macOS:  brew install hudochenkov/sshpass/sshpass
//	Debian: apt install sshpass
//
// Call immediately after:
//
//	bridgeName := createBridgeOnPE(t)
//	t.Cleanup(func() { deleteBridgeOnPE(t, bridgeName) })
func createBridgeOnPE(t *testing.T) string {
	t.Helper()

	peIP := testVars.PeIP
	peUser := testVars.SshPeUsername
	pePass := testVars.SshPePassword
	if peIP == "" || peUser == "" || pePass == "" {
		t.Skip("PE SSH credentials not configured in test_config_v2.json " +
			"(pe_ip / ssh_pe_username / ssh_pe_password)")
	}

	// -q suppresses SSH-level warnings ("Warning: Permanently added...").
	// The Nutanix CVM banner ("Nutanix Controller VM") is sent to stdout by
	// the server and cannot be suppressed client-side, so Go parses only the
	// last ^br[0-9]+$ line from stdout as a safety net.
	sshBase := fmt.Sprintf(
		"sshpass -e ssh -q -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null %s@%s",
		peUser, peIP,
	)
	hostssh := "/usr/local/nutanix/cluster/bin/hostssh"
	manageOvs := "/usr/local/nutanix/cluster/bin/manage_ovs"
	getBridges := hostssh + " 'ovs-vsctl list-br'"

	// Each step is a separate SSH call so that `source /etc/profile` in step 2
	// runs in its own remote shell and cannot abort the local script under
	// set -e.  All diagnostic output is routed to stderr (captured in
	// `debugBuf`).  Only `echo "$next"` in step 3 goes to stdout.
	//
	// IMPORTANT: step 2 uses `>/dev/stderr 2>&1`, NOT `2>&1 >&2`.
	// With `2>&1 >&2`: fd2→stdout first, then fd1→(now stdout) = both on stdout.
	// With `>/dev/stderr 2>&1`: fd1→stderr first, then fd2→(now stderr) = both on stderr.
	localScript := fmt.Sprintf(`set -euo pipefail

echo "[step 1] listing existing bridges on %[4]s" >&2
# hostssh relays the bridge list on stderr, so merge it (2>&1) before grep.
# grep '^br[0-9]+$' discards the SSH banner, hostssh headers, and brN.local /
# br.* names, leaving only real numbered bridges.
bridges=$(%[1]s "%[2]s" 2>&1 | tr -d '\r' | grep -E '^br[0-9]+$' || true)
echo "[step 1] raw bridges: ${bridges:-<none>}" >&2
if [ -z "$bridges" ]; then
  max=-1
else
  max=$(printf '%%s\n' "$bridges" | sed 's/^br//' | sort -n | tail -1)
fi
next="br$((max + 1))"
echo "[step 1] next bridge name: $next" >&2

echo "[step 2] creating bridge $next on %[4]s" >&2
%[1]s "source /etc/profile; %[3]s --bridge_name $next create_single_bridge" >/dev/stderr 2>&1
echo "[step 2] manage_ovs returned $?" >&2

echo "[step 3] verifying $next appears in bridge list (up to 60 s)" >&2
for i in $(seq 1 12); do
  echo "[step 3] poll $i/12 ..." >&2
  if %[1]s "%[2]s" 2>&1 | tr -d '\r' | grep -qE "^${next}$"; then
    echo "[step 3] bridge $next verified" >&2
    echo "$next"
    exit 0
  fi
  sleep 5
done
echo "[step 3] bridge $next not visible on %[4]s after 60 s" >&2
exit 1
`, sshBase, getBridges, manageOvs, peIP)

	var stderr bytes.Buffer
	cmd := exec.Command("bash", "-c", localScript)
	cmd.Env = append(os.Environ(), "SSHPASS="+pePass)
	cmd.Stderr = &stderr

	out, err := cmd.Output()
	t.Logf("createBridgeOnPE: stderr output:\n%s", stderr.String())
	if err != nil {
		t.Fatalf("createBridgeOnPE: failed on %s: %v\nstdout: %q\nstderr:\n%s",
			peIP, err, strings.TrimSpace(string(out)), stderr.String())
	}

	// Parse the last ^br[0-9]+$ line from stdout. `echo "$next"` is the only
	// intentional stdout write, but the CVM banner may precede it.
	bridgeRe := regexp.MustCompile(`^br[0-9]+$`)
	var name string
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if bridgeRe.MatchString(line) {
			name = line
		}
	}
	if name == "" {
		t.Fatalf("createBridgeOnPE: no bridge name found in stdout: %q\nstderr:\n%s",
			string(out), stderr.String())
	}
	t.Logf("createBridgeOnPE: created bridge %s ", name)
	return name
}

// deleteBridgeOnPE removes the named OVS bridge from the PE cluster via
// hostssh + ovs-vsctl del-br. Errors are logged but do not fail the test:
// if the VS destroy step already removed the bridge through the VS delete API,
// this is a no-op (the remote `|| true` absorbs the error).
func deleteBridgeOnPE(t *testing.T, bridgeName string) {
	t.Helper()

	peIP := testVars.PeIP
	peUser := testVars.SshPeUsername
	pePass := testVars.SshPePassword

	t.Logf("deleteBridgeOnPE: removing %s", bridgeName)

	sshBase := fmt.Sprintf(
		"sshpass -e ssh -q -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null %s@%s",
		peUser, peIP,
	)
	hostssh := "/usr/local/nutanix/cluster/bin/hostssh"

	localScript := fmt.Sprintf(`%s "%s 'ovs-vsctl del-br %s'" >/dev/stderr 2>&1 || true`,
		sshBase, hostssh, bridgeName)

	var stderr bytes.Buffer
	cmd := exec.Command("bash", "-c", localScript)
	cmd.Env = append(os.Environ(), "SSHPASS="+pePass)
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Logf("deleteBridgeOnPE: warning: could not delete %s: %v\nstderr:\n%s",
			bridgeName, err, stderr.String())
	} else {
		t.Logf("deleteBridgeOnPE: deleted %s", bridgeName)
	}
}

func testAccCheckNutanixVirtualSwitchDestroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_virtual_switch_v2" {
			continue
		}

		getReq := networkingVsReq.GetVirtualSwitchByIdRequest{
			ExtId: utils.StringPtr(rs.Primary.ID),
		}
		_, err := conn.NetworkingAPI.VirtualSwitchAPI.GetVirtualSwitchById(ctx, &getReq)
		if err == nil {
			return fmt.Errorf("virtual switch %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

// =============================================================================
// 4. Config builders (ordered to match their tests above)
// =============================================================================

// -- Positive --

func testAccVirtualSwitchMinimumConfig(name string) string {
	return fmt.Sprintf(`
# Fetch cluster and host information for virtual switch creation
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

resource "nutanix_virtual_switch_v2" "test" {
  name = "%[1]s"
  bond_mode = "NONE"
  clusters {
    ext_id = local.cluster_ext_id
    hosts {
      ext_id = local.host_ext_id
    }
  }
}
`, name)
}

func testAccVirtualSwitchWithAllAttributesConfig(name, description string) string {
	return fmt.Sprintf(`
# Fetch cluster and host information for virtual switch creation
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

resource "nutanix_virtual_switch_v2" "test" {
  name        = "%[1]s_all"
  description = "%[2]s"
  bond_mode   = "NONE"
  mtu         = 1500
  clusters {
    ext_id          = local.cluster_ext_id
    vlan_identifier = 0
    hosts {
      ext_id = local.host_ext_id
    }
  }
  igmp_spec {
    is_snooping_enabled = false
    snooping_timeout    = 300
  }
}
`, name, description)
}

// testAccVirtualSwitchWithAllAttributesUpdatedConfig renames the VS, updates
// its description, and enables the IGMP querier, mirroring the Ansible "Update
// virtual switch with all attributes" task. `is_default` is omitted
// (Computed-only in the schema) but asserted to remain false.
//
// querier_spec.vlan_id_list is intentionally NOT set: the API only accepts
// VLANs that are configured on the switch (via subnets), and these tests do
// not provision a subnet. Setting it (even to [0]) is rejected with
// "VLAN(s) [...] are not configured to this Virtual Switch".
func testAccVirtualSwitchWithAllAttributesUpdatedConfig(name, description string) string {
	return fmt.Sprintf(`
# Fetch cluster and host information for virtual switch update
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

resource "nutanix_virtual_switch_v2" "test" {
  name        = "%[1]s_updated"
  description = "%[2]s"
  bond_mode   = "NONE"
  mtu         = 1500
  clusters {
    ext_id          = local.cluster_ext_id
    vlan_identifier = 0
    hosts {
      ext_id = local.host_ext_id
    }
  }
  igmp_spec {
    is_snooping_enabled = true
    snooping_timeout    = 300
    querier_spec {
      is_querier_enabled = true
    }
  }
}
`, name, description)
}

// testAccVirtualSwitchFromExistingBridgeConfig wires cluster UUID + host
// UUID from existing v2 data sources and plugs the lab-prepared bridge
// name into clusters[0].existing_bridge_name. That triggers the migrate
// code path in the provider; without it the create call would route to
// the standard endpoint and the API would silently auto-allocate a new
// bridge.
//
// bond_mode, mtu, igmp_spec, etc. are deliberately omitted because the
// migrate API does not honor them. Setting them would still work (the
// provider only logs a warning), but leaving them out keeps the wire
// payload minimal and the test focused.
func testAccVirtualSwitchFromExistingBridgeConfig(name, description, bridgeName string) string {
	return fmt.Sprintf(`
# Discover cluster + host UUIDs from existing v2 data sources. The bridge
# is prepared out-of-band (lab prep) and its name is injected from the test
# config via testVars.Networking.BridgeName.
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

resource "nutanix_virtual_switch_v2" "test" {
  name        = "%[1]s_existing"
  description = "%[2]s"
  clusters {
    ext_id               = local.cluster_ext_id
    existing_bridge_name = "%[3]s"
    hosts {
      ext_id = local.host_ext_id
    }
  }
}
`, name, description, bridgeName)
}

// -- Share / unshare lifecycle --

// testAccVirtualSwitchShareBaseConfig builds the common scaffolding for the
// share/unshare lifecycle test: cluster + host data sources, an optional pair
// of project resources (projectsBlock), and the Virtual Switch resource with
// an optional shared_with_projects argument (sharedBlock).
func testAccVirtualSwitchShareBaseConfig(name, projectsBlock, sharedBlock string) string {
	return fmt.Sprintf(`
# Fetch cluster and host information for virtual switch creation.
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
%[2]s
resource "nutanix_virtual_switch_v2" "test_vs" {
  name        = "acc-test-vswitch-%[1]s"
  description = "Virtual Switch created by Terraform Acc Test"
  # bond_mode NONE: ACTIVE_BACKUP requires >= 2 uplink ports per host, which
  # this config does not specify. This test exercises share/unshare, not bonding.
  bond_mode   = "NONE"
  clusters {
    ext_id = local.cluster_ext_id
    hosts {
      ext_id = local.host_ext_id
    }
  }
%[3]s
}
`, name, projectsBlock, sharedBlock)
}

// testAccVirtualSwitchShareProjectsBlock declares the two projects the
// Virtual Switch is shared with across steps 2-4.
func testAccVirtualSwitchShareProjectsBlock(name string) string {
	return fmt.Sprintf(`
resource "nutanix_project_v2" "test_project_1" {
  name        = "acc-test-proj1-%[1]s"
  project_id  = "acc-test-proj1-%[1]s"
  description = "Test Project 1"
}

resource "nutanix_project_v2" "test_project_2" {
  name        = "acc-test-proj2-%[1]s"
  project_id  = "acc-test-proj2-%[1]s"
  description = "Test Project 2"
}
`, name)
}

// Step 1: create the Virtual Switch with no projects and no shares.
func testAccVirtualSwitchShareStep1Config(name string) string {
	return testAccVirtualSwitchShareBaseConfig(name, "", "")
}

// Step 2: declare both projects and share the Virtual Switch with project 1.
func testAccVirtualSwitchShareStep2Config(name string) string {
	shared := `  shared_with_projects = [
    nutanix_project_v2.test_project_1.ext_id
  ]`
	return testAccVirtualSwitchShareBaseConfig(name, testAccVirtualSwitchShareProjectsBlock(name), shared)
}

// Step 3: swap the share -- unshare project 1, share project 2.
func testAccVirtualSwitchShareStep3Config(name string) string {
	shared := `  shared_with_projects = [
    nutanix_project_v2.test_project_2.ext_id
  ]`
	return testAccVirtualSwitchShareBaseConfig(name, testAccVirtualSwitchShareProjectsBlock(name), shared)
}

// Step 4: clear shared_with_projects entirely (unshare project 2).
func testAccVirtualSwitchShareStep4Config(name string) string {
	shared := `  shared_with_projects = []`
	return testAccVirtualSwitchShareBaseConfig(name, testAccVirtualSwitchShareProjectsBlock(name), shared)
}

func testAccVirtualSwitchConfig(name, description string) string {
	return fmt.Sprintf(`
# Fetch cluster and host information for virtual switch creation
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

resource "nutanix_virtual_switch_v2" "test" {
  name          = "%[1]s"
  description   = "%[2]s"
  bond_mode     = "NONE"
  mtu           = 9000
  is_quick_mode = false
  clusters {
    ext_id          = local.cluster_ext_id
    vlan_identifier = 100
    gateway_ip_address {
      value         = "192.168.1.1"
      prefix_length = 24
    }
    hosts {
      ext_id        = local.host_ext_id
      host_nics     = ["eth2", "eth3"]
      active_uplink = "eth2"
      ip_address {
        ip {
          value         = "192.168.1.99"
          prefix_length = 24
        }
        prefix_length = 24
      }
    }
  }
  igmp_spec {
    is_snooping_enabled = true
    snooping_timeout    = 300
    querier_spec {
      is_querier_enabled = false
    }
  }
}
`, name, description)
}

func testAccVirtualSwitchConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "nutanix_virtual_switch_v2" "test" {
  name        = "%s"
  description = "updated description"
}
`, name)
}

func testAccVirtualSwitchConfigWithDatasource(name string) string {
	return fmt.Sprintf(`
resource "nutanix_virtual_switch_v2" "test" {
  name        = "%s"
  description = "updated description"
}

data "nutanix_virtual_switch_v2" "test" {
  ext_id = nutanix_virtual_switch_v2.test.id
}
`, name)
}

func testAccVirtualSwitchesListConfig() string {
	return `
data "nutanix_virtual_switches_v2" "test" {}
`
}

// testAccVirtualSwitchInfoByExtIdConfig creates a fully-populated Virtual
// Switch (incl. IGMP querier) and reads it back through the single-object
// datasource keyed on the resource's ext_id.
func testAccVirtualSwitchInfoByExtIdConfig(name, description string) string {
	return fmt.Sprintf(`
# Fetch cluster and host information for virtual switch creation
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

resource "nutanix_virtual_switch_v2" "test" {
  name        = "%[1]s_updated"
  description = "%[2]s"
  bond_mode   = "NONE"
  mtu         = 1500
  clusters {
    ext_id          = local.cluster_ext_id
    vlan_identifier = 0
    hosts {
      ext_id = local.host_ext_id
    }
  }
  igmp_spec {
    is_snooping_enabled = true
    snooping_timeout    = 300
    querier_spec {
      is_querier_enabled = true
    }
  }
}

data "nutanix_virtual_switch_v2" "test" {
  ext_id = nutanix_virtual_switch_v2.test.id
}
`, name, description)
}

// testAccVirtualSwitchInfoByWrongExtIdConfig reads the single datasource with a
// bogus ext_id (same sentinel UUID used by the other negative tests).
func testAccVirtualSwitchInfoByWrongExtIdConfig() string {
	return `
data "nutanix_virtual_switch_v2" "test" {
  ext_id = "04595a7b-010c-4a89-58dd-060c94e9a1f0"
}
`
}

// testAccVirtualSwitchListWithFilterConfig creates a VS with a unique name and
// lists virtual switches filtered to that exact name (expecting one result).
func testAccVirtualSwitchListWithFilterConfig(name string) string {
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

resource "nutanix_virtual_switch_v2" "test" {
  name      = "%[1]s"
  bond_mode = "NONE"
  clusters {
    ext_id = local.cluster_ext_id
    hosts {
      ext_id = local.host_ext_id
    }
  }
}

data "nutanix_virtual_switches_v2" "test" {
  filter     = "name eq '%[1]s'"
  depends_on = [nutanix_virtual_switch_v2.test]
}
`, name)
}

// testAccVirtualSwitchListWithLimitConfig keeps the same VS in state and lists
// virtual switches capped to a single result via $limit.
func testAccVirtualSwitchListWithLimitConfig(name string) string {
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

resource "nutanix_virtual_switch_v2" "test" {
  name      = "%[1]s"
  bond_mode = "NONE"
  clusters {
    ext_id = local.cluster_ext_id
    hosts {
      ext_id = local.host_ext_id
    }
  }
}

data "nutanix_virtual_switches_v2" "test" {
  limit      = 1
  depends_on = [nutanix_virtual_switch_v2.test]
}
`, name)
}

func testAccNodeSchedulableStatusesConfig() string {
	return `
data "nutanix_clusters_v2" "test" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}

data "nutanix_node_schedulable_statuses_v2" "test" {
  cluster_id = data.nutanix_clusters_v2.test.cluster_entities[0].ext_id
}
`
}

func testAccVpcVirtualSwitchMappingsConfig() string {
	return `
data "nutanix_vpc_virtual_switch_mappings_v2" "test" {}
`
}

// -- Negative --

// testAccVirtualSwitchMissingNameConfig deliberately omits the Required `name`
// argument so Terraform's schema validation rejects the plan. No data sources
// or other attributes are included because the failure must come from the
// missing-name validation, not from any API/network interaction.
func testAccVirtualSwitchMissingNameConfig() string {
	return `
resource "nutanix_virtual_switch_v2" "test" {
  description = "missing name negative test"
}
`
}

// testAccVirtualSwitchMissingBondModeConfig supplies every other field the
// API needs (name and a valid cluster+host reference) and intentionally
// omits `bond_mode`. The plan succeeds; the create request reaches Prism
// without a `bondMode` field; the API is expected to reject it.
func testAccVirtualSwitchMissingBondModeConfig(name string) string {
	return fmt.Sprintf(`
# Fetch cluster and host information for virtual switch creation
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

resource "nutanix_virtual_switch_v2" "test" {
  name = "%[1]s"
  clusters {
    ext_id = local.cluster_ext_id
    hosts {
      ext_id = local.host_ext_id
    }
  }
}
`, name)
}

// testAccVirtualSwitchMissingClustersConfig supplies `name` and `bond_mode`
// and intentionally omits the `clusters` block. No data sources are
// referenced because without `clusters` there is nothing to look up; the
// failure must come from the API rejecting a create with no cluster targets.
func testAccVirtualSwitchMissingClustersConfig(name string) string {
	return fmt.Sprintf(`
resource "nutanix_virtual_switch_v2" "test" {
  name      = "%[1]s"
  bond_mode = "NONE"
}
`, name)
}

// testAccVirtualSwitchFromExistingBridgeAlreadyInUseStep1Config: a single
// VS that migrates the lab-prepared bridge. Used as the pre-condition for
// the "already in use" negative test -- after this step the bridge is
// owned by `first`, and a second migrate against the same bridge must
// fail. Resource is named `first` (not `test`) so the step-2 config can
// add a `second` resource without renaming this one.
func testAccVirtualSwitchFromExistingBridgeAlreadyInUseStep1Config(name, bridgeName string) string {
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

resource "nutanix_virtual_switch_v2" "first" {
  name        = "%[1]s_first"
  description = "first VS taking ownership of the bridge"
  clusters {
    ext_id               = local.cluster_ext_id
    existing_bridge_name = "%[2]s"
    hosts {
      ext_id = local.host_ext_id
    }
  }
}
`, name, bridgeName)
}

// testAccVirtualSwitchFromExistingBridgeAlreadyInUseStep2Config: keeps the
// `first` VS exactly as step 1 declared it and adds a `second` VS that
// tries to migrate the SAME bridge. depends_on is set explicitly so that
// even though the first resource is already in state, the framework
// records the ordering dependency cleanly. The apply of this step must
// fail at the second resource's create call.
func testAccVirtualSwitchFromExistingBridgeAlreadyInUseStep2Config(name, bridgeName string) string {
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

resource "nutanix_virtual_switch_v2" "first" {
  name        = "%[1]s_first"
  description = "first VS taking ownership of the bridge"
  clusters {
    ext_id               = local.cluster_ext_id
    existing_bridge_name = "%[2]s"
    hosts {
      ext_id = local.host_ext_id
    }
  }
}

resource "nutanix_virtual_switch_v2" "second" {
  name        = "%[1]s_second"
  description = "second VS attempting to migrate the same bridge (must fail)"
  clusters {
    ext_id               = local.cluster_ext_id
    existing_bridge_name = "%[2]s"
    hosts {
      ext_id = local.host_ext_id
    }
  }
  depends_on = [nutanix_virtual_switch_v2.first]
}
`, name, bridgeName)
}

// testAccVirtualSwitchFromExistingBridgeWrongClusterRefConfig hard-codes a
// syntactically-valid but non-existent cluster UUID (same sentinel as the
// Ansible equivalent) and supplies a real bridge name. Host UUID is
// discovered from the hosts data source -- it has to be a real host or
// the schema's cluster {} block won't even build a valid migrate payload.
func testAccVirtualSwitchFromExistingBridgeWrongClusterRefConfig(name, bridgeName string) string {
	return fmt.Sprintf(`
data "nutanix_hosts_v2" "test" {}

locals {
  host_ext_id = data.nutanix_hosts_v2.test.host_entities[0].ext_id
}

resource "nutanix_virtual_switch_v2" "test" {
  name        = "%[1]s_wrong_cluster"
  description = "negative test: bogus cluster UUID"
  clusters {
    ext_id               = "04595a7b-010c-4a89-58dd-060c94e9a1f0"
    existing_bridge_name = "%[2]s"
    hosts {
      ext_id = local.host_ext_id
    }
  }
}
`, name, bridgeName)
}

// testAccVirtualSwitchFromExistingBridgeMissingClusterRefConfig declares a
// `clusters` block (so the provider enters the migrate code path) but
// omits `ext_id`, causing detectMigrationFromExistingBridge to return an
// empty clusterRef. The migrate body is sent without `clusterReference`
// and Prism is expected to reject it.
func testAccVirtualSwitchFromExistingBridgeMissingClusterRefConfig(name, bridgeName string) string {
	return fmt.Sprintf(`
data "nutanix_hosts_v2" "test" {}

locals {
  host_ext_id = data.nutanix_hosts_v2.test.host_entities[0].ext_id
}

resource "nutanix_virtual_switch_v2" "test" {
  name        = "%[1]s_no_cluster"
  description = "negative test: clusters[0].ext_id omitted"
  clusters {
    existing_bridge_name = "%[2]s"
    hosts {
      ext_id = local.host_ext_id
    }
  }
}
`, name, bridgeName)
}
