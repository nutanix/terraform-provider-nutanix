package microsegv2_test

/*
Test Plan — nutanix_network_security_policy_export_v2

This resource exports a point-in-time snapshot of the network security
policies identified by policy_ext_ids. It is backed by an async prepare-export
task followed by a synchronous download of the serialized payload.

Scenarios:
1. Basic lifecycle — create the export resource with a list of policy ext_ids.
   Verify task_ext_id and exported_payload are populated after create.
2. Export all — omit policy_ext_ids so the cluster exports every policy.
3. Selected with project scope — export specific policies that belong to a project
   by setting both policy_ext_ids and project_ext_id.
4. Export all with project scope — omit policy_ext_ids but set project_ext_id so the
   cluster exports every policy scoped to that project.
5. Negative — omit the required file_path and expect a validation error.

// Import not supported for this resource — it is a point-in-time snapshot.
*/

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const exportFilePath = "/tmp/test_policy_export.bin"
const exportAllFilePath = "/tmp/test_policy_export_all.bin"

func TestAccV2NutanixNetworkSecurityPolicyExport_Basic(t *testing.T) {
	t.Cleanup(func() { _ = os.Remove(exportFilePath) })
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testNetworkSecurityPolicyExportV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkSecurityPolicyExportV2Config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNetworkSecurityPolicyExportV2, "task_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameNetworkSecurityPolicyExportV2, "exported_payload"),
					resource.TestCheckResourceAttr(resourceNameNetworkSecurityPolicyExportV2, "policy_ext_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceNameNetworkSecurityPolicyExportV2, "file_path", exportFilePath),
					testAccCheckExportFileValid(resourceNameNetworkSecurityPolicyExportV2),
				),
			},
		},
	})
}

// When policy_ext_ids is omitted, the cluster exports all network security policies.
func TestAccV2NutanixNetworkSecurityPolicyExport_AllPolicies(t *testing.T) {
	t.Cleanup(func() { _ = os.Remove(exportAllFilePath) })
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testNetworkSecurityPolicyExportV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkSecurityPolicyExportV2ConfigAllPolicies(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNetworkSecurityPolicyExportV2, "task_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameNetworkSecurityPolicyExportV2, "exported_payload"),
					resource.TestCheckResourceAttr(resourceNameNetworkSecurityPolicyExportV2, "policy_ext_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceNameNetworkSecurityPolicyExportV2, "file_path", exportAllFilePath),
					testAccCheckExportFileValid(resourceNameNetworkSecurityPolicyExportV2),
				),
			},
		},
	})
}

// Export without the required file_path must fail with a validation error.
func TestAccV2NutanixNetworkSecurityPolicyExport_MissingFilePath(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNetworkSecurityPolicyExportV2ConfigMissingFilePath(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func testAccNetworkSecurityPolicyExportV2ConfigMissingFilePath() string {
	return `
resource "nutanix_network_security_policy_export_v2" "test" {
  policy_ext_ids = ["00000000-0000-0000-0000-000000000001"]
}
`
}

func testAccNetworkSecurityPolicyExportV2Config() string {
	r := acctest.RandIntBetween(1, 1000)
	return fmt.Sprintf(`
resource "nutanix_category_v2" "cat1" {
  key   = "tf-test-cat-%[2]d"
  value = "tf-test-cat-%[2]d-value"
}

# create first nsp with intra group rule
resource "nutanix_network_security_policy_v2" "test-1" {
  name        = "tf-nsp-test-%[2]d"
  description = "Network Security Policy Test for import and export test"
  state       = "MONITOR"
  type        = "APPLICATION"
  scope       = "ALL_VLAN"
  rules {
    description = "intra group rule"
    type        = "INTRA_GROUP"
    spec {
      intra_entity_group_rule_spec {
        secured_group_category_references = [
          nutanix_category_v2.cat1.id,
        ]
        secured_group_action = "ALLOW"
      }
    }
  }
  is_hitlog_enabled = false
  lifecycle {
    ignore_changes = [rules]
  }
}

# export security policy
resource "nutanix_network_security_policy_export_v2" "test" {
  policy_ext_ids = [nutanix_network_security_policy_v2.test-1.id]
  file_path = "%[1]s"
  depends_on = [nutanix_network_security_policy_v2.test-1]
}


`, exportFilePath, r)
}

func testAccNetworkSecurityPolicyExportV2ConfigAllPolicies() string {
	r1 := acctest.RandIntBetween(1, 1000)
	r2 := acctest.RandIntBetween(1001, 2000)
	r3 := acctest.RandIntBetween(2001, 3000)
	return fmt.Sprintf(`
resource "nutanix_category_v2" "cat1" {
  key   = "tf-test-cat-%[2]d"
  value = "tf-test-cat-%[2]d-value"
}

resource "nutanix_category_v2" "cat2" {
  key   = "tf-test-cat-%[3]d"
  value = "tf-test-cat-%[3]d-value"
}

resource "nutanix_category_v2" "cat3" {
  key   = "tf-test-cat-%[4]d"
  value = "tf-test-cat-%[4]d-value"
}

# create first nsp with intra group rule
resource "nutanix_network_security_policy_v2" "test-1" {
  name        = "tf-nsp-test-%[2]d"
  description = "Network Security Policy Test for import and export test"
  state       = "MONITOR"
  type        = "APPLICATION"
  scope       = "ALL_VLAN"
  rules {
    description = "intra group rule"
    type        = "INTRA_GROUP"
    spec {
      intra_entity_group_rule_spec {
        secured_group_category_references = [
          nutanix_category_v2.cat1.id,
        ]
        secured_group_action = "ALLOW"
      }
    }
  }
  is_hitlog_enabled = false
  lifecycle {
    ignore_changes = [rules]
  }
}

# create second nsp with application rule
resource "nutanix_network_security_policy_v2" "test-2" {
  name        = "tf-nsp-test-%[3]d"
  description = "Network Security Policy Test for import and export test"
  state       = "MONITOR"
  type        = "APPLICATION"
  scope       = "ALL_VLAN"
  rules {
    description = "application rule"
    type        = "APPLICATION"
    spec {
      application_rule_spec {
        secured_group_category_references = [
          nutanix_category_v2.cat2.id,
        ]
        src_category_references = [
          nutanix_category_v2.cat3.id,
        ]
        is_all_protocol_allowed = true
      }
    }
  }
  lifecycle {
    ignore_changes = [rules]
  }
}

# export all policies
resource "nutanix_network_security_policy_export_v2" "test" {
  file_path  = "%[1]s"
  depends_on = [nutanix_network_security_policy_v2.test-1, nutanix_network_security_policy_v2.test-2]
}


`, exportAllFilePath, r1, r2, r3)
}

