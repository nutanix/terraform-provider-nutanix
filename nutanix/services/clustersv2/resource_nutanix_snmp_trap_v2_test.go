package clustersv2_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const (
	resourceNameSnmpTrap          = "nutanix_snmp_trap_v2.test"
	dataSourceNameSnmpConfigCheck = "data.nutanix_snmp_config_v2.after_create"
)

// TestAccV2NutanixSnmpTrapResource_Basic
//   - Step 1 creates a V2 trap with a randomized address + port pair to
//     avoid collisions with traps already configured on the cluster, and
//     verifies all attributes (community_string, protocol, port, version,
//     address.ipv4.value, ext_id) are returned by the API.
//   - Step 2 bumps the port to exercise UpdateSnmpTrapById.
//   - CheckDestroy verifies DeleteSnmpTrapById actually removed the trap
//     so the next test run starts from a clean SNMP config.
func TestAccV2NutanixSnmpTrapResource_Basic(t *testing.T) {
	octet := acctest.RandIntRange(100, 200)
	trapIP := fmt.Sprintf("10.0.0.%d", octet)
	port := acctest.RandIntRange(30000, 39000)
	portUpdated := port + 1
	community := "public"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixSnmpTrapDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create SNMP V2 trap.
			{
				Config: testSnmpTrapV2Config(trapIP, community, port, "UDP"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameSnmpTrap, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameSnmpTrap, "cluster_ext_id"),
					resource.TestCheckResourceAttr(resourceNameSnmpTrap, "version", "V2"),
					resource.TestCheckResourceAttr(resourceNameSnmpTrap, "address.0.ipv4.0.value", trapIP),
					resource.TestCheckResourceAttr(resourceNameSnmpTrap, "community_string", community),
					resource.TestCheckResourceAttr(resourceNameSnmpTrap, "protocol", "UDP"),
					resource.TestCheckResourceAttr(resourceNameSnmpTrap, "port", fmt.Sprintf("%d", port)),
				),
			},
			// Step 2: Update the port to exercise UpdateSnmpTrapById.
			{
				Config: testSnmpTrapV2Config(trapIP, community, portUpdated, "UDP"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameSnmpTrap, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameSnmpTrap, "version", "V2"),
					resource.TestCheckResourceAttr(resourceNameSnmpTrap, "address.0.ipv4.0.value", trapIP),
					resource.TestCheckResourceAttr(resourceNameSnmpTrap, "community_string", community),
					resource.TestCheckResourceAttr(resourceNameSnmpTrap, "port", fmt.Sprintf("%d", portUpdated)),
				),
			},
		},
	})
}

// TestAccV2NutanixSnmpTrapResource_UsingSnmpUser exercises the V3 trap
// creation path, which (unlike V2) requires an existing SNMP user on the
// same cluster and references it by username.
//
//   - Provisions an inline nutanix_snmp_user_v2 with full auth+priv
//     (SHA / AES) so the trap has a valid V3 binding.
//   - Creates a V3 trap pointing at that user via the `username`
//     attribute, with a randomized IP last-octet and port to avoid
//     colliding with leftover state from earlier runs.
//   - On test exit testAccCheckNutanixSnmpTrapDestroy verifies the
//     trap was actually removed; testAccCheckNutanixSnmpUserDestroy
//     (registered via composition) verifies the inline user was too.
//
// The depends_on link from trap → user is implicit because the trap
// references nutanix_snmp_user_v2.test.username, which gives Terraform
// the ordering it needs to create the user before the trap and tear
// them down in reverse on destroy.
func TestAccV2NutanixSnmpTrapResource_UsingSnmpUser(t *testing.T) {
	r := acctest.RandInt()
	username := fmt.Sprintf("tf-acc-snmp-user-%d", r)
	octet := acctest.RandIntRange(100, 200)
	trapIP := fmt.Sprintf("10.0.0.%d", octet)
	port := acctest.RandIntRange(30000, 39000)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		CheckDestroy: resource.ComposeAggregateTestCheckFunc(
			testAccCheckNutanixSnmpTrapDestroy,
			testAccCheckNutanixSnmpUserDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testSnmpTrapV3Config(username, trapIP, port, "UDP"),
				Check: resource.ComposeTestCheckFunc(
					// Resource-level checks: round-trip via GetSnmpTrapById.
					resource.TestCheckResourceAttrSet(resourceNameSnmpTrap, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameSnmpTrap, "cluster_ext_id"),
					resource.TestCheckResourceAttr(resourceNameSnmpTrap, "version", "V3"),
					resource.TestCheckResourceAttr(resourceNameSnmpTrap, "address.0.ipv4.0.value", trapIP),
					resource.TestCheckResourceAttr(resourceNameSnmpTrap, "username", username),
					resource.TestCheckResourceAttr(resourceNameSnmpTrap, "protocol", "UDP"),
					resource.TestCheckResourceAttr(resourceNameSnmpTrap, "port", fmt.Sprintf("%d", port)),
					// V3 traps must NOT carry a community_string — make
					// sure we didn't accidentally leak the V2 attribute
					// through the resource state.
					resource.TestCheckResourceAttr(resourceNameSnmpTrap, "community_string", ""),
					// Cluster-level check: the trap must appear in the full
					// SNMP config of the cluster, with the right
					// version/address/username/protocol/port tuple.
					resource.TestCheckResourceAttrSet(dataSourceNameSnmpConfigCheck, "traps.#"),
					testCheckSnmpConfigContainsTrap(
						dataSourceNameSnmpConfigCheck,
						"V3", trapIP, username, "UDP", port,
					),
				),
			},
		},
	})
}

func testSnmpTrapV2Config(trapIP, community string, port int, protocol string) string {
	return fmt.Sprintf(`
data "nutanix_clusters_v2" "aos" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}

locals {
  clusterExtID = data.nutanix_clusters_v2.aos.cluster_entities[0].ext_id
}

resource "nutanix_snmp_trap_v2" "test" {
  cluster_ext_id   = local.clusterExtID
  version          = "V2"
  community_string = "%s"
  port             = %d
  protocol         = "%s"

  address {
    ipv4 {
      value = "%s"
    }
  }
}
`, community, port, protocol, trapIP)
}

// testCheckSnmpConfigContainsTrap walks `traps.*` on the given SNMP
// config data source, picks the entry that matches the (version,
// ipv4.value) identity tuple, and asserts username/protocol/port match
// the expected values. Returns an error if no matching entry exists or
// if any field doesn't match — this is what catches both "trap was never
// created on the cluster" and "trap exists but with wrong attributes"
// regressions.
func testCheckSnmpConfigContainsTrap(dataSource, version, ipv4, username, protocol string, port int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[dataSource]
		if !ok {
			return fmt.Errorf("data source %s not found in state", dataSource)
		}
		attrs := rs.Primary.Attributes

		count, err := strconv.Atoi(attrs["traps.#"])
		if err != nil {
			return fmt.Errorf("data source %s missing traps.#: %w", dataSource, err)
		}

		for i := 0; i < count; i++ {
			prefix := fmt.Sprintf("traps.%d.", i)
			if attrs[prefix+"version"] != version {
				continue
			}
			if attrs[prefix+"address.0.ipv4.0.value"] != ipv4 {
				continue
			}
			// Found the matching trap — now assert the remaining fields.
			if got := attrs[prefix+"username"]; got != username {
				return fmt.Errorf("traps[%d]: expected username=%q, got %q", i, username, got)
			}
			if got := attrs[prefix+"protocol"]; got != protocol {
				return fmt.Errorf("traps[%d]: expected protocol=%q, got %q", i, protocol, got)
			}
			gotPort, err := strconv.Atoi(attrs[prefix+"port"])
			if err != nil {
				return fmt.Errorf("traps[%d]: port not parseable as int (%q): %w", i, attrs[prefix+"port"], err)
			}
			if gotPort != port {
				return fmt.Errorf("traps[%d]: expected port=%d, got %d", i, port, gotPort)
			}
			return nil
		}
		return fmt.Errorf("trap with version=%q ipv4=%q not found in %s.traps (count=%d)",
			version, ipv4, dataSource, count)
	}
}

func testSnmpTrapV3Config(username, trapIP string, port int, protocol string) string {
	return fmt.Sprintf(`
data "nutanix_clusters_v2" "aos" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}

locals {
  clusterExtID = data.nutanix_clusters_v2.aos.cluster_entities[0].ext_id
}

resource "nutanix_snmp_user_v2" "test" {
  cluster_ext_id = local.clusterExtID
  username       = "%[1]s"
  auth_type      = "SHA"
  auth_key       = "auth-key-original-1234"
  priv_type      = "AES"
  priv_key       = "priv-key-original-1234"
}

resource "nutanix_snmp_trap_v2" "test" {
  cluster_ext_id = local.clusterExtID
  version        = "V3"
  username       = nutanix_snmp_user_v2.test.username
  port           = %[2]d
  protocol       = "%[3]s"

  address {
    ipv4 {
      value = "%[4]s"
    }
  }
}

# Pulls the cluster-wide SNMP config so the test can assert the V3 trap
# actually shows up in the traps list with the right attributes.
# depends_on is required because data sources are otherwise read at
# plan-time, before the trap has been created.
data "nutanix_snmp_config_v2" "after_create" {
  cluster_ext_id = local.clusterExtID
  depends_on     = [nutanix_snmp_trap_v2.test]
}
`, username, port, protocol, trapIP)
}
