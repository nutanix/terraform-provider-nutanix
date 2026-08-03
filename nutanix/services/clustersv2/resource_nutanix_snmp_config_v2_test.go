package clustersv2_test

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	cmgmtConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/clustermgmt-go-client/v17/models/clustermgmt/v4/config"
	cmgmtRequest "github.com/nutanix-core/ntnx-api-golang-sdk-internal/clustermgmt-go-client/v17/models/clustermgmt/v4/request/clusters"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const (
	resourceNameSnmpConfig          = "nutanix_snmp_config_v2.test"
	dataSourceNameInitialSnmpState  = "data.nutanix_snmp_config_v2.initial"
	dataSourceNameAfterAddSnmpState = "data.nutanix_snmp_config_v2.test"
)

// snmpConfigInitial* package-level variables capture the cluster's SNMP
// is_enabled value before TestAccV2NutanixSnmpConfigResource_Basic runs so
// that testAccCheckNutanixSnmpConfigStatusDestroy can restore the cluster
// to its pre-test state regardless of how the test exited. Guarded by a
// mutex so concurrent tests in the same package don't race on it.
var (
	snmpConfigInitialMu      sync.Mutex
	snmpConfigInitialCaptred bool
	snmpConfigInitialEnabled bool
	snmpConfigInitialCluster string
)

// snmpTransport* package-level variables capture the cluster's SNMP
// transports list before TestAccV2NutanixSnmpConfigResource_Transport runs.
// They're populated by Step 0 (captureSnmpTransportInitialState) and
// consumed by both the after-add verification check and the CheckDestroy
// hook.
var (
	snmpTransportInitialMu       sync.Mutex
	snmpTransportInitialCaptured bool
	snmpTransportInitialCluster  string
	snmpTransportInitialCount    int
	snmpTransportTestPort        int
	snmpTransportTestProtocol    string
)

// TestAccV2NutanixSnmpConfigResource_Basic exercises the *status mode* of
// the unified nutanix_snmp_config_v2 resource: only `is_enabled` is set,
// so the resource manages the cluster-wide SNMP enabled flag via
// UpdateSnmpStatus.
//
//   - Step 0 captures the cluster's initial SNMP is_enabled value via the
//     data source (no resource yet, so no mutation).
//   - Step 1 creates the resource with is_enabled=true to exercise Create
//     and the post-task convergence wait.
//   - Step 2 flips to false to exercise Update in the disable direction.
//   - Step 3 flips back to true to exercise Update in the enable direction.
//   - On test completion CheckDestroy runs Delete (a state-only drop in
//     status mode) and then testAccCheckNutanixSnmpConfigStatusDestroy
//     restores is_enabled to whatever Step 0 captured.
//
// We intentionally do NOT toggle relative to local.current_snmp_status:
// that local reads from a data source which is re-evaluated on every
// refresh, so the toggle expression flips again after each apply and the
// config can never be idempotent.
func TestAccV2NutanixSnmpConfigResource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixSnmpConfigStatusDestroy,
		Steps: []resource.TestStep{
			// Step 0: capture initial state, no resource yet.
			{
				Config: testSnmpConfigCaptureOnlyConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameInitialSnmpState, "cluster_ext_id"),
					resource.TestCheckResourceAttrSet(dataSourceNameInitialSnmpState, "is_enabled"),
					captureInitialSnmpState(),
				),
			},
			// Step 1: create the resource, explicitly enabling SNMP.
			{
				Config: testSnmpConfigStatusResourceConfig("true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameSnmpConfig, "id"),
					resource.TestCheckResourceAttrSet(resourceNameSnmpConfig, "cluster_ext_id"),
					resource.TestCheckResourceAttr(resourceNameSnmpConfig, "is_enabled", "true"),
				),
			},
			// Step 2: update to disable SNMP.
			{
				Config: testSnmpConfigStatusResourceConfig("false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSnmpConfig, "is_enabled", "false"),
				),
			},
			// Step 3: update to re-enable SNMP (exercises both directions
			// of the post-task convergence wait).
			{
				Config: testSnmpConfigStatusResourceConfig("true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSnmpConfig, "is_enabled", "true"),
				),
			},
		},
	})
}

// TestAccV2NutanixSnmpConfigResource_Transport exercises the *transport
// mode* of the unified nutanix_snmp_config_v2 resource: `port` and
// `protocol` are set, so the resource manages a single SNMP transport on
// the cluster via AddSnmpTransport / RemoveSnmpTransport. is_enabled is
// left out of the HCL config in this test so the resource doesn't touch
// the cluster's global SNMP toggle.
//
//   - Step 0 captures the cluster ext_id, the current transports list size,
//     and computes a dynamic (port, protocol) pair guaranteed not to clash
//     with an existing transport using
//     port = max(existing transports.port + [161]) + 1, protocol = "UDP".
//   - Step 1 adds the transport and asserts the post-add SNMP config
//     exposes exactly one extra transport whose (port, protocol) match
//     the test values.
//   - Step 2 drops the resource from the config (forcing
//     RemoveSnmpTransport on apply) and polls the SDK directly until the
//     cluster's SnmpConfig.transports list reflects the removal — covers
//     the parent-task / child-reconciliation race where the Prism task
//     can flip to SUCCEEDED before the transports list is updated.
//   - Step 3 re-introduces data.nutanix_snmp_config_v2.test as a fresh
//     instance so Terraform reads it from scratch, giving us a faithful
//     post-remove snapshot in state to assert against.
//   - On test completion testAccCheckNutanixSnmpConfigTransportDestroy
//     re-fetches the SNMP config via the SDK and verifies the env is
//     back to its captured initial shape.
func TestAccV2NutanixSnmpConfigResource_Transport(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixSnmpConfigTransportDestroy,
		Steps: []resource.TestStep{
			// Step 0: capture initial transports state, no resource yet.
			{
				Config: testSnmpConfigCaptureOnlyConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameInitialSnmpState, "cluster_ext_id"),
					captureSnmpTransportInitialState(),
				),
			},
			// Step 1: add the transport and verify it shows up in the
			// post-add SNMP config.
			{
				Config: testSnmpConfigTransportResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameSnmpConfig, "id"),
					resource.TestCheckResourceAttrSet(resourceNameSnmpConfig, "cluster_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameSnmpConfig, "port"),
					resource.TestCheckResourceAttr(resourceNameSnmpConfig, "protocol", "UDP"),
					verifySnmpTransportAdded(),
				),
			},
			// Step 2: explicit reset — drop the resource (and its frozen
			// port helper) from the config so the framework destroys
			// them on apply, then poll the SDK directly until the
			// cluster's SnmpConfig.transports list reflects the
			// removal. We can't rely on the data source here because
			// nothing it depends_on has changed and its inputs are
			// stable, so Terraform would not re-read it and state would
			// still show the just-removed transport.
			{
				Config: testSnmpConfigTransportRemoveConfig(),
				Check: resource.ComposeTestCheckFunc(
					verifySnmpTransportRemoved(),
				),
			},
			// Step 3: re-introduce data.nutanix_snmp_config_v2.test as a
			// brand-new data source instance. Because Step 2 dropped
			// the previous one, Terraform will read this one from
			// scratch on apply, giving us a faithful post-remove
			// snapshot in state. By the time we get here Step 2's
			// SDK-poll has already confirmed cluster-side convergence,
			// so the read is guaranteed accurate.
			{
				Config: testSnmpConfigTransportPostRemoveConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameAfterAddSnmpState, "cluster_ext_id"),
					verifySnmpTransportPostRemoveStateMatchesInitial(),
				),
			},
		},
	})
}

// captureInitialSnmpState reads the cluster_ext_id and is_enabled
// attributes from the data.nutanix_snmp_config_v2.initial state entry
// (populated by Step 0) and stores them in the package-level variables
// consumed by testAccCheckNutanixSnmpConfigStatusDestroy.
func captureInitialSnmpState() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[dataSourceNameInitialSnmpState]
		if !ok {
			return fmt.Errorf("%q not found in state", dataSourceNameInitialSnmpState)
		}
		clusterExtID := rs.Primary.Attributes["cluster_ext_id"]
		if clusterExtID == "" {
			return fmt.Errorf("%q has empty cluster_ext_id", dataSourceNameInitialSnmpState)
		}
		isEnabledStr := rs.Primary.Attributes["is_enabled"]
		isEnabled, err := strconv.ParseBool(isEnabledStr)
		if err != nil {
			return fmt.Errorf("parsing initial is_enabled (%q): %w", isEnabledStr, err)
		}

		snmpConfigInitialMu.Lock()
		snmpConfigInitialCluster = clusterExtID
		snmpConfigInitialEnabled = isEnabled
		snmpConfigInitialCaptred = true
		snmpConfigInitialMu.Unlock()

		log.Printf("[INFO] captured initial SNMP state: cluster=%s is_enabled=%t", clusterExtID, isEnabled)
		return nil
	}
}

// captureSnmpTransportInitialState reads the current transports list from
// the data.nutanix_snmp_config_v2.initial state entry, stores the cluster
// ext_id and the size of the list in package vars, and computes the
// (test_port, test_protocol) the resource step will use:
// protocol = "UDP", port = max(existing transport ports + [161]) + 1.
func captureSnmpTransportInitialState() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[dataSourceNameInitialSnmpState]
		if !ok {
			return fmt.Errorf("%q not found in state", dataSourceNameInitialSnmpState)
		}
		clusterExtID := rs.Primary.Attributes["cluster_ext_id"]
		if clusterExtID == "" {
			return fmt.Errorf("%q has empty cluster_ext_id", dataSourceNameInitialSnmpState)
		}
		countStr := rs.Primary.Attributes["transports.#"]
		count, err := strconv.Atoi(countStr)
		if err != nil {
			return fmt.Errorf("parsing initial transports count (%q): %w", countStr, err)
		}

		// max(existing ports, 161) + 1
		maxPort := 161
		for i := 0; i < count; i++ {
			portStr := rs.Primary.Attributes[fmt.Sprintf("transports.%d.port", i)]
			if p, err := strconv.Atoi(portStr); err == nil && p > maxPort {
				maxPort = p
			}
		}
		testPort := maxPort + 1
		testProto := "UDP"

		snmpTransportInitialMu.Lock()
		snmpTransportInitialCluster = clusterExtID
		snmpTransportInitialCount = count
		snmpTransportTestPort = testPort
		snmpTransportTestProtocol = testProto
		snmpTransportInitialCaptured = true
		snmpTransportInitialMu.Unlock()

		log.Printf("[INFO] captured initial SNMP transports: cluster=%s count=%d test_port=%d test_protocol=%s",
			clusterExtID, count, testPort, testProto)
		return nil
	}
}

// verifySnmpTransportAdded reads data.nutanix_snmp_config_v2.test (the
// post-add snapshot of the cluster's SNMP config) from state and asserts:
//
//  1. transports list size == initial + 1
//  2. there is exactly one entry whose (port, protocol) match the
//     captured test values
func verifySnmpTransportAdded() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		snmpTransportInitialMu.Lock()
		captured := snmpTransportInitialCaptured
		initCount := snmpTransportInitialCount
		port := snmpTransportTestPort
		proto := snmpTransportTestProtocol
		snmpTransportInitialMu.Unlock()
		if !captured {
			return fmt.Errorf("initial SNMP transport state was not captured by Step 0")
		}

		rs, ok := s.RootModule().Resources[dataSourceNameAfterAddSnmpState]
		if !ok {
			return fmt.Errorf("%q not found in state", dataSourceNameAfterAddSnmpState)
		}
		cnt, err := strconv.Atoi(rs.Primary.Attributes["transports.#"])
		if err != nil {
			return fmt.Errorf("parsing after-add transports count: %w", err)
		}
		if cnt != initCount+1 {
			return fmt.Errorf("expected %d transports after add (initial %d + 1), got %d", initCount+1, initCount, cnt)
		}

		matches := 0
		for i := 0; i < cnt; i++ {
			p := rs.Primary.Attributes[fmt.Sprintf("transports.%d.port", i)]
			pr := rs.Primary.Attributes[fmt.Sprintf("transports.%d.protocol", i)]
			if p == strconv.Itoa(port) && pr == proto {
				matches++
			}
		}
		if matches != 1 {
			return fmt.Errorf("expected exactly 1 transport matching (port=%d, protocol=%s) in after-add config, found %d", port, proto, matches)
		}
		log.Printf("[INFO] verified after-add SNMP config: %d transports total, includes (port=%d, protocol=%s)", cnt, port, proto)
		return nil
	}
}

// snmpTransportRemovePollInterval / snmpTransportRemovePollTimeout: how
// often / how long verifySnmpTransportRemoved re-fetches the SNMP config
// from the cluster while waiting for RemoveSnmpTransport to converge
// server-side.
const (
	snmpTransportRemovePollInterval = 2 * time.Second
	snmpTransportRemovePollTimeout  = 60 * time.Second
)

// verifySnmpTransportRemoved is the Step 2 counterpart of
// verifySnmpTransportAdded. It does NOT read the
// data.nutanix_snmp_config_v2.test data source from state — Terraform only
// re-reads a data source when an input changes (or via depends_on), and
// in Step 2 the resource we used to depend on no longer exists, so the
// transports list in state would be stale and would still report the
// just-removed transport.
//
// Instead, we go straight to the SDK and poll the cluster's SNMP config
// until either:
//
//  1. transports list size is back to the captured initial count AND no
//     entry remains whose (port, protocol) match the captured test
//     values, or
//  2. snmpTransportRemovePollTimeout elapses (in which case the test
//     fails with the last observed mismatch).
//
// This polling loop also handles the parent-task / child-reconciliation
// race we cope with in updateSnmpStatus: the RemoveSnmpTransport task
// can flip to SUCCEEDED before the SnmpConfig.transports list is updated
// server-side, so a single post-apply fetch can still see the removed
// transport.
func verifySnmpTransportRemoved() resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		snmpTransportInitialMu.Lock()
		captured := snmpTransportInitialCaptured
		clusterExtID := snmpTransportInitialCluster
		initCount := snmpTransportInitialCount
		port := snmpTransportTestPort
		proto := snmpTransportTestProtocol
		snmpTransportInitialMu.Unlock()
		if !captured {
			return fmt.Errorf("initial SNMP transport state was not captured by Step 0")
		}

		conn := acc.TestAccProvider.Meta().(*conns.Client).ClusterAPI
		ctx := context.Background()

		var lastErr error
		deadline := time.Now().Add(snmpTransportRemovePollTimeout)
		for attempt := 1; ; attempt++ {
			getReq := cmgmtRequest.GetSnmpConfigByClusterIdRequest{
				ClusterExtId: utils.StringPtr(clusterExtID),
			}
			resp, err := conn.ClusterEntityAPI.GetSnmpConfigByClusterId(ctx, &getReq)
			if err != nil {
				lastErr = fmt.Errorf("read SNMP config for cluster %s: %w", clusterExtID, err)
			} else {
				cfg, ok := resp.Data.GetValue().(cmgmtConfig.SnmpConfig)
				if !ok {
					lastErr = fmt.Errorf("unexpected SNMP config response type for cluster %s", clusterExtID)
				} else {
					actualCount := len(cfg.Transports)
					found := false
					for _, tr := range cfg.Transports {
						if tr.Port == nil || tr.Protocol == nil {
							continue
						}
						if utils.IntValue(tr.Port) == port && common.FlattenPtrEnum(tr.Protocol) == proto {
							found = true
							break
						}
					}
					if actualCount == initCount && !found {
						log.Printf("[INFO] verified after-remove SNMP config (attempt %d): %d transports total, (port=%d, protocol=%s) gone",
							attempt, actualCount, port, proto)
						return nil
					}
					lastErr = fmt.Errorf("transports count=%d (initial=%d), (port=%d, protocol=%s) still present=%t",
						actualCount, initCount, port, proto, found)
				}
			}

			if time.Now().After(deadline) {
				return fmt.Errorf("SNMP transport remove did not converge for cluster %s within %s: %w",
					clusterExtID, snmpTransportRemovePollTimeout, lastErr)
			}
			log.Printf("[DEBUG] post-remove SNMP config not yet converged (attempt %d): %v; retrying in %s",
				attempt, lastErr, snmpTransportRemovePollInterval)
			time.Sleep(snmpTransportRemovePollInterval)
		}
	}
}

// verifySnmpTransportPostRemoveStateMatchesInitial reads the freshly-read
// data.nutanix_snmp_config_v2.test snapshot from Terraform state (Step 3
// re-declares the data source so the read is genuinely post-remove) and
// asserts the cluster is back to its pre-test shape:
//
//  1. transports list size == captured initial count.
//  2. no entry whose (port, protocol) match the captured test values.
//
// This is the state-level companion to verifySnmpTransportRemoved's SDK
// polling check: the SDK poll guarantees the cluster converged, and this
// check guarantees Terraform's view of that convergence is the post-remove
// shape we expect.
func verifySnmpTransportPostRemoveStateMatchesInitial() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		snmpTransportInitialMu.Lock()
		captured := snmpTransportInitialCaptured
		initCount := snmpTransportInitialCount
		port := snmpTransportTestPort
		proto := snmpTransportTestProtocol
		snmpTransportInitialMu.Unlock()
		if !captured {
			return fmt.Errorf("initial SNMP transport state was not captured by Step 0")
		}

		rs, ok := s.RootModule().Resources[dataSourceNameAfterAddSnmpState]
		if !ok {
			return fmt.Errorf("%q not found in state", dataSourceNameAfterAddSnmpState)
		}
		cnt, err := strconv.Atoi(rs.Primary.Attributes["transports.#"])
		if err != nil {
			return fmt.Errorf("parsing post-remove transports count: %w", err)
		}
		if cnt != initCount {
			return fmt.Errorf("expected %d transports after remove (initial), got %d", initCount, cnt)
		}
		for i := 0; i < cnt; i++ {
			p := rs.Primary.Attributes[fmt.Sprintf("transports.%d.port", i)]
			pr := rs.Primary.Attributes[fmt.Sprintf("transports.%d.protocol", i)]
			if p == strconv.Itoa(port) && pr == proto {
				return fmt.Errorf("transport (port=%d, protocol=%s) still present in SNMP config after remove", port, proto)
			}
		}
		log.Printf("[INFO] verified post-remove SNMP config from state: %d transports total, (port=%d, protocol=%s) gone",
			cnt, port, proto)
		return nil
	}
}

// testAccCheckNutanixSnmpConfigStatusDestroy is the CheckDestroy hook for
// TestAccV2NutanixSnmpConfigResource_Basic. It compares the cluster's
// current SNMP is_enabled value against the value captured in Step 0 and,
// if they differ, calls UpdateSnmpStatus to restore it. This ensures the
// test is non-destructive: the cluster is left in the exact SNMP
// enabled/disabled state it was in before the test ran.
func testAccCheckNutanixSnmpConfigStatusDestroy(_ *terraform.State) error {
	snmpConfigInitialMu.Lock()
	captured := snmpConfigInitialCaptred
	clusterExtID := snmpConfigInitialCluster
	desired := snmpConfigInitialEnabled
	snmpConfigInitialMu.Unlock()

	if !captured {
		log.Printf("[WARN] CheckDestroy: initial SNMP state was not captured; skipping restore")
		return nil
	}

	conn := acc.TestAccProvider.Meta().(*conns.Client).ClusterAPI
	ctx := context.Background()

	getReq := cmgmtRequest.GetSnmpConfigByClusterIdRequest{
		ClusterExtId: utils.StringPtr(clusterExtID),
	}
	getResp, err := conn.ClusterEntityAPI.GetSnmpConfigByClusterId(ctx, &getReq)
	if err != nil {
		return fmt.Errorf("CheckDestroy: read SNMP config for cluster %s: %w", clusterExtID, err)
	}
	cfg, ok := getResp.Data.GetValue().(cmgmtConfig.SnmpConfig)
	if !ok {
		return fmt.Errorf("CheckDestroy: unexpected SNMP config response type for cluster %s", clusterExtID)
	}
	actual := utils.BoolValue(cfg.IsEnabled)
	if actual == desired {
		log.Printf("[INFO] CheckDestroy: cluster %s SNMP is_enabled already at initial value (%t); no restore needed", clusterExtID, desired)
		return nil
	}

	log.Printf("[INFO] CheckDestroy: restoring cluster %s SNMP is_enabled %t -> %t", clusterExtID, actual, desired)
	args := make(map[string]interface{})
	if etag := conn.ClusterEntityAPI.ApiClient.GetEtag(getResp); etag != "" {
		args["If-Match"] = utils.StringPtr(etag)
	}
	body := cmgmtConfig.NewSnmpStatusParam()
	body.IsEnabled = utils.BoolPtr(desired)
	updReq := cmgmtRequest.UpdateSnmpStatusRequest{
		ClusterExtId: utils.StringPtr(clusterExtID),
		Body:         body,
	}
	if _, err := conn.ClusterEntityAPI.UpdateSnmpStatus(ctx, &updReq, args); err != nil {
		return fmt.Errorf("CheckDestroy: restore SNMP is_enabled to %t for cluster %s: %w", desired, clusterExtID, err)
	}
	return nil
}

// testAccCheckNutanixSnmpConfigTransportDestroy is the CheckDestroy hook
// for TestAccV2NutanixSnmpConfigResource_Transport. It uses the SDK
// directly to query the cluster's SNMP config and asserts:
//
//  1. the transports list is back to its captured initial size
//  2. the captured (test_port, test_protocol) tuple is gone
//
// Acts as a final SDK-level safety net regardless of how the test
// exited.
func testAccCheckNutanixSnmpConfigTransportDestroy(_ *terraform.State) error {
	snmpTransportInitialMu.Lock()
	captured := snmpTransportInitialCaptured
	clusterExtID := snmpTransportInitialCluster
	initCount := snmpTransportInitialCount
	port := snmpTransportTestPort
	proto := snmpTransportTestProtocol
	snmpTransportInitialMu.Unlock()
	if !captured {
		log.Printf("[WARN] CheckDestroy: SNMP transport initial state not captured; skipping destroy verification")
		return nil
	}

	conn := acc.TestAccProvider.Meta().(*conns.Client).ClusterAPI
	ctx := context.Background()

	getReq := cmgmtRequest.GetSnmpConfigByClusterIdRequest{
		ClusterExtId: utils.StringPtr(clusterExtID),
	}
	resp, err := conn.ClusterEntityAPI.GetSnmpConfigByClusterId(ctx, &getReq)
	if err != nil {
		return fmt.Errorf("CheckDestroy: read SNMP config for cluster %s: %w", clusterExtID, err)
	}
	cfg, ok := resp.Data.GetValue().(cmgmtConfig.SnmpConfig)
	if !ok {
		return fmt.Errorf("CheckDestroy: unexpected SNMP config response type for cluster %s", clusterExtID)
	}

	actualCount := len(cfg.Transports)
	if actualCount != initCount {
		return fmt.Errorf("CheckDestroy: expected %d transports after destroy, got %d", initCount, actualCount)
	}
	for _, tr := range cfg.Transports {
		if tr.Port == nil || tr.Protocol == nil {
			continue
		}
		if utils.IntValue(tr.Port) == port && common.FlattenPtrEnum(tr.Protocol) == proto {
			return fmt.Errorf("CheckDestroy: transport (port=%d, protocol=%s) still present after destroy", port, proto)
		}
	}
	log.Printf("[INFO] CheckDestroy: SNMP transports back to initial count=%d, (port=%d, protocol=%s) removed", actualCount, port, proto)
	return nil
}

// testSnmpConfigCaptureOnlyConfig renders an HCL config that only declares
// the data sources needed to discover the cluster ext_id and read the
// initial SNMP config. No resource is created, so applying it doesn't
// mutate the cluster — leaving the data source's `is_enabled` and
// transports list faithful to the cluster's pre-test state. Shared by
// both the Basic and Transport tests' Step 0.
func testSnmpConfigCaptureOnlyConfig() string {
	return `
data "nutanix_clusters_v2" "aos" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}

locals {
  clusterExtID = data.nutanix_clusters_v2.aos.cluster_entities[0].ext_id
}

data "nutanix_snmp_config_v2" "initial" {
  cluster_ext_id = local.clusterExtID
}
`
}

// testSnmpConfigStatusResourceConfig renders the HCL for the status-mode
// (Basic) test. The `isEnabledExpr` argument is interpolated verbatim
// into the resource's is_enabled attribute, which lets tests pass either
// a literal (`"true"` / `"false"`) or a derived expression like
// `"!local.current_snmp_status"`.
func testSnmpConfigStatusResourceConfig(isEnabledExpr string) string {
	return `
data "nutanix_clusters_v2" "aos" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}

locals {
  clusterExtID = data.nutanix_clusters_v2.aos.cluster_entities[0].ext_id
}

data "nutanix_snmp_config_v2" "initial" {
  cluster_ext_id = local.clusterExtID
}

locals {
  current_snmp_status     = tobool(data.nutanix_snmp_config_v2.initial.is_enabled)
  current_snmp_transports = try(data.nutanix_snmp_config_v2.initial.transports, [])
}

# Status-mode: only is_enabled is set. The unified resource recognises
# this and only issues UpdateSnmpStatus, leaving any existing transports
# untouched.
resource "nutanix_snmp_config_v2" "test" {
  cluster_ext_id = local.clusterExtID
  is_enabled     = ` + isEnabledExpr + `
}
`
}

// testSnmpConfigTransportResourceConfig renders the HCL for the transport-
// mode test's Step 1: it adds an SNMP transport using a dynamically-
// computed port, then re-reads the SNMP config so the verification check
// can assert the transport is present.
//
// Why terraform_data.frozen_test_port:
//
//	local.test_transport_port = max(existing transports.port + [161]) + 1.
//	After our resource adds a transport, data.nutanix_snmp_config_v2.initial
//	re-reads on the next refresh and now contains the new transport, so
//	the local would jump by 1 again on every plan and the resource would
//	drift indefinitely. Pinning the value through terraform_data with
//	`ignore_changes = [input]` freezes it at create-time, breaking the
//	feedback loop and making the config idempotent.
//
// is_enabled is intentionally omitted so the resource is in pure
// transport mode — the global SNMP toggle is left at whatever the
// cluster had before the test.
func testSnmpConfigTransportResourceConfig() string {
	return `
data "nutanix_clusters_v2" "aos" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}

locals {
  clusterExtID = data.nutanix_clusters_v2.aos.cluster_entities[0].ext_id
}

data "nutanix_snmp_config_v2" "initial" {
  cluster_ext_id = local.clusterExtID
}

locals {
  current_snmp_transports = try(data.nutanix_snmp_config_v2.initial.transports, [])
  test_transport_protocol = "UDP"
  test_transport_port     = max(concat([for t in local.current_snmp_transports : t.port], [161])...) + 1
}

resource "terraform_data" "frozen_test_port" {
  input = local.test_transport_port

  lifecycle {
    ignore_changes = [input]
  }
}

# Transport-mode: port + protocol are set, is_enabled is omitted. The
# unified resource recognises this and only issues AddSnmpTransport on
# create / RemoveSnmpTransport on delete, leaving the global SNMP
# enabled flag alone.
resource "nutanix_snmp_config_v2" "test" {
  cluster_ext_id = local.clusterExtID
  port           = terraform_data.frozen_test_port.output
  protocol       = local.test_transport_protocol
}

data "nutanix_snmp_config_v2" "test" {
  cluster_ext_id = local.clusterExtID
  depends_on     = [nutanix_snmp_config_v2.test]
}
`
}

// testSnmpConfigTransportRemoveConfig is the transport test's Step 2 HCL:
// it intentionally drops the nutanix_snmp_config_v2.test resource (and
// its terraform_data port pinner) from the config so the framework's
// apply phase calls RemoveSnmpTransport.
//
// We deliberately do NOT include data.nutanix_snmp_config_v2.test here.
// Terraform only re-reads a data source when its inputs change or it
// depends_on something that changed; in this step the resource it used
// to depend on is gone and the cluster_ext_id input is stable, so the
// data source's state would be carried over from Step 1 and the
// transports list would be stale (still showing the just-removed
// transport). verifySnmpTransportRemoved goes straight to the SDK with
// a polling wait instead.
func testSnmpConfigTransportRemoveConfig() string {
	return `
data "nutanix_clusters_v2" "aos" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}

locals {
  clusterExtID = data.nutanix_clusters_v2.aos.cluster_entities[0].ext_id
}
`
}

// testSnmpConfigTransportPostRemoveConfig is the transport test's Step 3
// HCL: it re-introduces data.nutanix_snmp_config_v2.test (which Step 2
// deliberately dropped) so that Terraform reads it fresh on apply,
// giving us a faithful post-remove snapshot of the cluster's SNMP
// config in state for verifySnmpTransportPostRemoveStateMatchesInitial
// to assert against.
//
// Splitting the read into its own step (rather than carrying the data
// source through Step 2) is the only way to force a re-read: keeping
// the same data source instance with stable inputs lets Terraform reuse
// the pre-remove value from prior state.
func testSnmpConfigTransportPostRemoveConfig() string {
	return testSnmpConfigTransportRemoveConfig() + `
data "nutanix_snmp_config_v2" "test" {
  cluster_ext_id = local.clusterExtID
}
`
}
