package networkingv2_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	nicprofilesreq "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v17/models/networking/v4/request/nicprofiles"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const (
	resourceNameNicProfileV2     = "nutanix_nic_profile_v2.test"
	datasourceNameNicProfileV2   = "data.nutanix_nic_profile_v2.test"
	datasourceNameNicProfilesV2  = "data.nutanix_nic_profiles_v2.test"
	defaultNicProfileFamily      = "15b3:101d"
	defaultNicProfileDescription = "NIC profile managed by Terraform"
)

// skipNicProfileV2 skips NIC profile v2 acceptance tests. The resource and data
// sources are intentionally not registered in the provider for the current
// release, so their acceptance tests cannot run until it is re-enabled.
func skipNicProfileV2(t *testing.T) {
	t.Helper()
	t.Skip("skipping: nutanix_nic_profile_v2 is deferred and not registered in the provider.")
}

// nicProfileFamily returns the NIC family configured for acceptance testing,
// falling back to a sane default when test_config_v2.json does not set one.
func nicProfileFamily() string {
	if testVars.NicProfile.NicFamily != "" {
		return testVars.NicProfile.NicFamily
	}
	return defaultNicProfileFamily
}

// nicProfileCapability returns the capability_type to use for the host-NIC
// association test. It must match a capability actually supported by the
// configured host NIC (see preEnv/scripts/automate_nic_profile.sh, which
// discovers a capable NIC and records the matching nic_family + capability_type
// + host_nic_ext_ids). Falls back to DP_OFFLOAD to mirror the runbook.
func nicProfileCapability() string {
	switch testVars.NicProfile.CapabilityType {
	case "SRIOV", "DP_OFFLOAD", "PCIE_PASSTHROUGH":
		return testVars.NicProfile.CapabilityType
	default:
		return "DP_OFFLOAD"
	}
}

// TestAccV2NutanixNicProfileResource_Basic covers the full happy path: create a
// SR-IOV NIC profile, read it back through both data sources, mutate the
// mutable fields (name + description) and finally import it.
func TestAccV2NutanixNicProfileResource_Basic(t *testing.T) {
	skipNicProfileV2(t)
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-nic-profile-%d", r)
	updatedName := name + "-updated"
	updatedDesc := defaultNicProfileDescription + " updated"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNicProfileV2Destroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccNicProfileV2Config(name, defaultNicProfileDescription, "SRIOV", "ETHERNET"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNicProfileV2, "id"),
					resource.TestCheckResourceAttrSet(resourceNameNicProfileV2, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameNicProfileV2, "name", name),
					resource.TestCheckResourceAttr(resourceNameNicProfileV2, "description", defaultNicProfileDescription),
					resource.TestCheckResourceAttr(resourceNameNicProfileV2, "nic_family", nicProfileFamily()),
					resource.TestCheckResourceAttr(resourceNameNicProfileV2, "capability_config.#", "1"),
					resource.TestCheckResourceAttr(resourceNameNicProfileV2, "capability_config.0.capability_type", "SRIOV"),
					resource.TestCheckResourceAttr(resourceNameNicProfileV2, "operating_mode", "ETHERNET"),
				),
			},
			// Read via data sources (single by ext_id + list)
			{
				Config: testAccNicProfileV2Config(name, defaultNicProfileDescription, "SRIOV", "ETHERNET") + testAccNicProfileV2DataSources(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceNameNicProfileV2, "ext_id", resourceNameNicProfileV2, "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameNicProfileV2, "name", name),
					resource.TestCheckResourceAttr(datasourceNameNicProfileV2, "nic_family", nicProfileFamily()),
					resource.TestCheckResourceAttr(datasourceNameNicProfileV2, "capability_config.0.capability_type", "SRIOV"),
					resource.TestCheckResourceAttr(datasourceNameNicProfileV2, "operating_mode", "ETHERNET"),
					checkAttributeLength(datasourceNameNicProfilesV2, "nic_profiles", 1),
				),
			},
			// Update mutable fields (name + description)
			{
				Config: testAccNicProfileV2Config(updatedName, updatedDesc, "SRIOV", "ETHERNET"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameNicProfileV2, "name", updatedName),
					resource.TestCheckResourceAttr(resourceNameNicProfileV2, "description", updatedDesc),
					// immutable fields must be unchanged
					resource.TestCheckResourceAttr(resourceNameNicProfileV2, "nic_family", nicProfileFamily()),
					resource.TestCheckResourceAttr(resourceNameNicProfileV2, "capability_config.0.capability_type", "SRIOV"),
					resource.TestCheckResourceAttr(resourceNameNicProfileV2, "operating_mode", "ETHERNET"),
				),
			},
			// Import
			{
				ResourceName:            resourceNameNicProfileV2,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"host_nic_ext_ids"},
			},
		},
	})
}

// TestAccV2NutanixNicProfileResource_DPOffload validates the DP_OFFLOAD
// capability type in Ethernet mode.
func TestAccV2NutanixNicProfileResource_DPOffload(t *testing.T) {
	skipNicProfileV2(t)
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-nic-profile-dpo-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNicProfileV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNicProfileV2Config(name, defaultNicProfileDescription, "DP_OFFLOAD", "ETHERNET"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNicProfileV2, "id"),
					resource.TestCheckResourceAttr(resourceNameNicProfileV2, "capability_config.0.capability_type", "DP_OFFLOAD"),
					resource.TestCheckResourceAttr(resourceNameNicProfileV2, "operating_mode", "ETHERNET"),
				),
			},
		},
	})
}

// TestAccV2NutanixNicProfileResource_InfiniBandNotAllowed verifies that the API
// rejects user-created InfiniBand profiles. InfiniBand operating mode is
// reserved for system-defined profiles only (e.g. ConnectX7_INFINIBAND_*), so
// any attempt to create one through Terraform must surface the Atlas guardrail
// error rather than succeed.
func TestAccV2NutanixNicProfileResource_InfiniBandNotAllowed(t *testing.T) {
	skipNicProfileV2(t)
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-nic-profile-ib-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNicProfileV2Destroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNicProfileV2Config(name, defaultNicProfileDescription, "SRIOV", "INFINIBAND"),
				ExpectError: regexp.MustCompile("Users cannot create NIC profiles with InfiniBand mode"),
			},
		},
	})
}

// TestAccV2NutanixNicProfileDataSource_ByExtID reads a pre-existing NIC profile
// (configured via test_config_v2.json, e.g. DP_OFFLOAD_PROFILE) by its ext_id
// and validates the exported attributes. This does not create any resource.
func TestAccV2NutanixNicProfileDataSource_ByExtID(t *testing.T) {
	skipNicProfileV2(t)
	extID := testVars.NicProfile.ExtID
	if extID == "" {
		t.Skip("skipping: no nic_profile.ext_id configured in test_config_v2.json")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
data "nutanix_nic_profile_v2" "existing" {
  ext_id = "%s"
}
`, extID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.nutanix_nic_profile_v2.existing", "ext_id", extID),
					resource.TestCheckResourceAttr("data.nutanix_nic_profile_v2.existing", "id", extID),
					resource.TestCheckResourceAttr("data.nutanix_nic_profile_v2.existing", "name", testVars.NicProfile.Name),
					resource.TestCheckResourceAttr("data.nutanix_nic_profile_v2.existing", "capability_config.0.capability_type", testVars.NicProfile.CapabilityType),
					resource.TestCheckResourceAttrSet("data.nutanix_nic_profile_v2.existing", "nic_family"),
				),
			},
		},
	})
}

// TestAccV2NutanixNicProfileDataSource_List lists all NIC profiles and verifies
// the configured pre-existing profile is present with the expected attributes.
func TestAccV2NutanixNicProfileDataSource_List(t *testing.T) {
	skipNicProfileV2(t)
	extID := testVars.NicProfile.ExtID
	if extID == "" {
		t.Skip("skipping: no nic_profile.ext_id configured in test_config_v2.json")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
data "nutanix_nic_profiles_v2" "test" {}
`,
				Check: resource.ComposeTestCheckFunc(
					checkAttributeLength(datasourceNameNicProfilesV2, "nic_profiles", 1),
					testAccCheckNicProfileInList(datasourceNameNicProfilesV2, extID, testVars.NicProfile.Name, testVars.NicProfile.CapabilityType),
				),
			},
		},
	})
}

// TestAccV2NutanixNicProfileDataSource_NotFound ensures fetching a NIC profile
// by a non-existent ext_id surfaces an error.
func TestAccV2NutanixNicProfileDataSource_NotFound(t *testing.T) {
	skipNicProfileV2(t)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
data "nutanix_nic_profile_v2" "missing" {
  ext_id = "00000000-0000-0000-0000-000000000000"
}
`,
				ExpectError: regexp.MustCompile("error while fetching NIC profile"),
			},
		},
	})
}

// TestAccV2NutanixNicProfileResource_HostNicAssociation covers associating and
// then disassociating Host NICs from a NIC profile. It is skipped unless Host
// NIC UUIDs are provided in test_config_v2.json.
//
// Atlas enforces (unconditionally, in update_association_task.py) that a Host
// NIC's pci_model_id must equal the profile's nic_family AND that the NIC
// actually supports the requested capability before it can be associated.
// preEnv/scripts/automate_nic_profile.sh discovers a capability-capable Host NIC
// on the cluster and records the matching nic_family + capability_type +
// host_nic_ext_ids into test_config_v2.json, so this test creates the profile
// from those real, matched values and the association succeeds. It is only
// skipped when no capable NIC is configured.
func TestAccV2NutanixNicProfileResource_HostNicAssociation(t *testing.T) {
	skipNicProfileV2(t)
	hostNics := nonEmptyHostNicExtIDs()
	if len(hostNics) == 0 {
		t.Skip("skipping: no nic_profile.host_nic_ext_ids configured in test_config_v2.json")
	}

	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-nic-profile-hostnic-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNicProfileV2Destroy,
		Steps: []resource.TestStep{
			// Create with the Host NICs associated
			{
				Config: testAccNicProfileV2ConfigWithHostNics(name, defaultNicProfileDescription, "DP_OFFLOAD", "ETHERNET", hostNics),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNicProfileV2, "id"),
					resource.TestCheckResourceAttr(resourceNameNicProfileV2, "host_nic_ext_ids.#", fmt.Sprintf("%d", len(hostNics))),
					checkAttributeLengthEqual(resourceNameNicProfileV2, "host_nic_references", len(hostNics)),
				),
			},
			// Disassociate all Host NICs
			{
				Config: testAccNicProfileV2Config(name, defaultNicProfileDescription, "DP_OFFLOAD", "ETHERNET"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameNicProfileV2, "host_nic_ext_ids.#", "0"),
				),
			},
		},
	})
}

// ----------------------------------------------------------------------------
// Negative scenarios
// ----------------------------------------------------------------------------

// TestAccV2NutanixNicProfileResource_InvalidCapabilityType ensures the
// capability_type validation rejects unknown values before any API call.
func TestAccV2NutanixNicProfileResource_InvalidCapabilityType(t *testing.T) {
	skipNicProfileV2(t)
	name := fmt.Sprintf("tf-test-nic-profile-%d", acctest.RandInt())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNicProfileV2Destroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNicProfileV2Config(name, defaultNicProfileDescription, "INVALID_CAP", "ETHERNET"),
				ExpectError: regexp.MustCompile("capability_type to be one of"),
			},
		},
	})
}

// TestAccV2NutanixNicProfileResource_InvalidOperatingMode ensures the
// operating_mode validation rejects unknown values.
func TestAccV2NutanixNicProfileResource_InvalidOperatingMode(t *testing.T) {
	skipNicProfileV2(t)
	name := fmt.Sprintf("tf-test-nic-profile-%d", acctest.RandInt())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNicProfileV2Destroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccNicProfileV2Config(name, defaultNicProfileDescription, "SRIOV", "INVALID_MODE"),
				ExpectError: regexp.MustCompile("operating_mode to be one of"),
			},
		},
	})
}

// TestAccV2NutanixNicProfileResource_InvalidOwnerType ensures the owner_type
// validation rejects unknown values.
func TestAccV2NutanixNicProfileResource_InvalidOwnerType(t *testing.T) {
	skipNicProfileV2(t)
	name := fmt.Sprintf("tf-test-nic-profile-%d", acctest.RandInt())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNicProfileV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "nutanix_nic_profile_v2" "test" {
  name        = "%[1]s"
  description = "%[2]s"
  nic_family  = "%[3]s"
  owner_type  = "INVALID_OWNER"
  capability_config {
    capability_type = "SRIOV"
  }
  operating_mode = "ETHERNET"
}
`, name, defaultNicProfileDescription, nicProfileFamily()),
				ExpectError: regexp.MustCompile("owner_type to be one of"),
			},
		},
	})
}

// TestAccV2NutanixNicProfileResource_MissingNicFamily ensures the required
// nic_family argument is enforced.
func TestAccV2NutanixNicProfileResource_MissingNicFamily(t *testing.T) {
	skipNicProfileV2(t)
	name := fmt.Sprintf("tf-test-nic-profile-%d", acctest.RandInt())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNicProfileV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "nutanix_nic_profile_v2" "test" {
  name = "%[1]s"
  capability_config {
    capability_type = "SRIOV"
  }
  operating_mode = "ETHERNET"
}
`, name),
				ExpectError: regexp.MustCompile(`The argument "nic_family" is required`),
			},
		},
	})
}

// TestAccV2NutanixNicProfileResource_MissingCapabilityConfig ensures the
// required capability_config block is enforced.
func TestAccV2NutanixNicProfileResource_MissingCapabilityConfig(t *testing.T) {
	skipNicProfileV2(t)
	name := fmt.Sprintf("tf-test-nic-profile-%d", acctest.RandInt())
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNicProfileV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "nutanix_nic_profile_v2" "test" {
  name           = "%[1]s"
  nic_family     = "%[2]s"
  operating_mode = "ETHERNET"
}
`, name, nicProfileFamily()),
				ExpectError: regexp.MustCompile("capability_config"),
			},
		},
	})
}

// TestAccV2NutanixNicProfileResource_ImmutableUpdates ensures that updating any
// of the immutable fields (capability_config, nic_family, operating_mode) is
// rejected by the resource with a clear error.
func TestAccV2NutanixNicProfileResource_ImmutableUpdates(t *testing.T) {
	skipNicProfileV2(t)
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-nic-profile-immutable-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNicProfileV2Destroy,
		Steps: []resource.TestStep{
			// Create baseline
			{
				Config: testAccNicProfileV2Config(name, defaultNicProfileDescription, "SRIOV", "ETHERNET"),
				Check:  resource.TestCheckResourceAttrSet(resourceNameNicProfileV2, "id"),
			},
			// capability_config change is not supported
			{
				Config:      testAccNicProfileV2Config(name, defaultNicProfileDescription, "DP_OFFLOAD", "ETHERNET"),
				ExpectError: regexp.MustCompile("Update of capability_config is not supported"),
			},
			// operating_mode change is not supported
			{
				Config:      testAccNicProfileV2Config(name, defaultNicProfileDescription, "SRIOV", "INFINIBAND"),
				ExpectError: regexp.MustCompile("Update of operating_mode is not supported"),
			},
			// nic_family change is not supported
			{
				Config: fmt.Sprintf(`
resource "nutanix_nic_profile_v2" "test" {
  name        = "%[1]s"
  description = "%[2]s"
  nic_family  = "15b3:101f"
  capability_config {
    capability_type = "SRIOV"
  }
  operating_mode = "ETHERNET"
}
`, name, defaultNicProfileDescription),
				ExpectError: regexp.MustCompile("Update of nic_family is not supported"),
			},
		},
	})
}

// ----------------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------------

// testAccCheckNicProfileInList verifies the nic_profiles list data source
// contains an entry matching the given ext_id, and that its name and
// capability_type match the expected values.
func testAccCheckNicProfileInList(dataSourceName, extID, name, capabilityType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("data source not found: %s", dataSourceName)
		}

		countStr, ok := rs.Primary.Attributes["nic_profiles.#"]
		if !ok {
			return fmt.Errorf("nic_profiles.# attribute not found")
		}
		count, err := strconv.Atoi(countStr)
		if err != nil {
			return fmt.Errorf("invalid nic_profiles count: %w", err)
		}

		for i := 0; i < count; i++ {
			prefix := fmt.Sprintf("nic_profiles.%d.", i)
			if rs.Primary.Attributes[prefix+"ext_id"] != extID {
				continue
			}
			if got := rs.Primary.Attributes[prefix+"name"]; got != name {
				return fmt.Errorf("expected name %q for profile %s, got %q", name, extID, got)
			}
			if got := rs.Primary.Attributes[prefix+"capability_config.0.capability_type"]; got != capabilityType {
				return fmt.Errorf("expected capability_type %q for profile %s, got %q", capabilityType, extID, got)
			}
			return nil
		}

		return fmt.Errorf("NIC profile %s not found in data source %s", extID, dataSourceName)
	}
}

func nonEmptyHostNicExtIDs() []string {
	out := make([]string, 0, len(testVars.NicProfile.HostNicExtIDs))
	for _, id := range testVars.NicProfile.HostNicExtIDs {
		if id != "" {
			out = append(out, id)
		}
	}
	return out
}

func testAccNicProfileV2Config(name, description, capabilityType, operatingMode string) string {
	return fmt.Sprintf(`
resource "nutanix_nic_profile_v2" "test" {
  name        = "%[1]s"
  description = "%[2]s"
  nic_family  = "%[3]s"
  capability_config {
    capability_type = "%[4]s"
  }
  operating_mode = "%[5]s"
}
`, name, description, nicProfileFamily(), capabilityType, operatingMode)
}

func testAccNicProfileV2ConfigWithHostNics(name, description, capabilityType, operatingMode string, hostNics []string) string {
	hostNicsBlock := ""
	for _, id := range hostNics {
		hostNicsBlock += fmt.Sprintf("\n    \"%s\",", id)
	}
	return fmt.Sprintf(`
resource "nutanix_nic_profile_v2" "test" {
  name        = "%[1]s"
  description = "%[2]s"
  nic_family  = "%[3]s"
  host_nic_ext_ids = [%[6]s
  ]
  capability_config {
    capability_type = "%[4]s"
  }
  operating_mode = "%[5]s"
}
`, name, description, nicProfileFamily(), capabilityType, operatingMode, hostNicsBlock)
}

func testAccNicProfileV2DataSources() string {
	return `
data "nutanix_nic_profile_v2" "test" {
  ext_id = nutanix_nic_profile_v2.test.id
}

data "nutanix_nic_profiles_v2" "test" {}
`
}

func testAccCheckNicProfileV2Destroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client).NetworkingAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_nic_profile_v2" {
			continue
		}

		getRequest := nicprofilesreq.GetNicProfileByIdRequest{
			ExtId: utils.StringPtr(rs.Primary.ID),
		}
		if _, err := conn.NicProfilesAPI.GetNicProfileById(context.Background(), &getRequest); err != nil {
			// Not found -> successfully destroyed.
			continue
		}
		return fmt.Errorf("NIC profile %s still exists", rs.Primary.ID)
	}

	return nil
}
