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
const exportSelectedProjectFilePath = "/tmp/test_policy_export_selected_project.bin"
const exportAllProjectFilePath = "/tmp/test_policy_export_all_project.bin"

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

// Export specific policies that belong to a project: both policy_ext_ids and project_ext_id
// are set, so only the policies matching both the IDs and the project context are exported.
func TestAccV2NutanixNetworkSecurityPolicyExport_SelectedWithProject(t *testing.T) {
	t.Cleanup(func() { _ = os.Remove(exportSelectedProjectFilePath) })
	r := acctest.RandIntBetween(1, 1000)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testNetworkSecurityPolicyExportV2CheckDestroy,
		Steps: []resource.TestStep{
			// 1. Export the project-scoped policy.
			{
				Config: testAccNetworkSecurityPolicyExportV2ConfigSelectedWithProject(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNetworkSecurityPolicyExportV2, "task_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameNetworkSecurityPolicyExportV2, "exported_payload"),
					resource.TestCheckResourceAttr(resourceNameNetworkSecurityPolicyExportV2, "policy_ext_ids.#", "1"),
					resource.TestCheckResourceAttrPair(resourceNameNetworkSecurityPolicyExportV2, "project_ext_id", "nutanix_project_v2.test", "ext_id"),
					resource.TestCheckResourceAttr(resourceNameNetworkSecurityPolicyExportV2, "file_path", exportSelectedProjectFilePath),
					testAccCheckExportFileValid(resourceNameNetworkSecurityPolicyExportV2),
				),
			},
			// 2. Remove the policy (and export) so the category no longer has active project
			//    associations. The category stays shared with the project at this point.
			{
				Config: projectCategoryConfig(r, true),
			},
			// 3. Unshare the category from the project before teardown. A category that is still
			//    shared with a project cannot be deleted, which would also block the project's
			//    deletion; and the unshare itself only succeeds once the policy association is
			//    gone (step 2). Setting shared_with_projects = [] triggers the unshare.
			{
				Config: projectCategoryConfig(r, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nutanix_category_v2.cat1", "shared_with_projects.#", "0"),
				),
			},
		},
	})
}

// Export all policies that belong to a project: policy_ext_ids is omitted but project_ext_id
// is set, so every policy scoped to that project is exported.
func TestAccV2NutanixNetworkSecurityPolicyExport_AllWithProject(t *testing.T) {
	t.Cleanup(func() { _ = os.Remove(exportAllProjectFilePath) })
	r := acctest.RandIntBetween(1, 1000)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testNetworkSecurityPolicyExportV2CheckDestroy,
		Steps: []resource.TestStep{
			// 1. Export every policy scoped to the project.
			{
				Config: testAccNetworkSecurityPolicyExportV2ConfigAllWithProject(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNetworkSecurityPolicyExportV2, "task_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameNetworkSecurityPolicyExportV2, "exported_payload"),
					resource.TestCheckResourceAttr(resourceNameNetworkSecurityPolicyExportV2, "policy_ext_ids.#", "0"),
					resource.TestCheckResourceAttrPair(resourceNameNetworkSecurityPolicyExportV2, "project_ext_id", "nutanix_project_v2.test", "ext_id"),
					resource.TestCheckResourceAttr(resourceNameNetworkSecurityPolicyExportV2, "file_path", exportAllProjectFilePath),
					testAccCheckExportFileValid(resourceNameNetworkSecurityPolicyExportV2),
				),
			},
			// 2. Remove the policy (and export) so the category no longer has active project
			//    associations. The category stays shared with the project at this point.
			{
				Config: projectCategoryConfig(r, true),
			},
			// 3. Unshare the category from the project before teardown (see note above); the
			//    unshare only succeeds once the policy association is gone (step 2).
			{
				Config: projectCategoryConfig(r, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nutanix_category_v2.cat1", "shared_with_projects.#", "0"),
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

// projectCategoryConfig returns a project and a category that is optionally shared with that
// project. When shared is false the category is unshared (shared_with_projects = []), which is
// required before the category and the project can be deleted. It is used on its own for the
// teardown steps (after the policy has been removed) and as the base for the scoped policy.
func projectCategoryConfig(r int, shared bool) string {
	shareLine := "shared_with_projects = []"
	if shared {
		shareLine = "shared_with_projects = [nutanix_project_v2.test.ext_id]"
	}
	return fmt.Sprintf(`
resource "nutanix_project_v2" "test" {
  name        = "tf-proj-%[1]d"
  project_id  = "tf-proj-%[1]d"
  description = "export project scope test"
}

resource "nutanix_category_v2" "cat1" {
  key   = "tf-test-cat-%[1]d"
  value = "tf-test-cat-%[1]d-value"
  %[2]s
}
`, r, shareLine)
}

// projectScopedPolicyConfig returns a project, a category shared with it, and a network
// security policy scoped to the project so the policy actually belongs to a project context.
func projectScopedPolicyConfig(r int) string {
	return projectCategoryConfig(r, true) + fmt.Sprintf(`
resource "nutanix_network_security_policy_v2" "test-1" {
  name           = "tf-nsp-test-%[1]d"
  description    = "Network Security Policy Test for project scoped export"
  state          = "MONITOR"
  type           = "APPLICATION"
  scope          = "ALL_VLAN"
  project_ext_id = nutanix_project_v2.test.ext_id
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
  depends_on = [nutanix_project_v2.test]
}
`, r)
}

func testAccNetworkSecurityPolicyExportV2ConfigSelectedWithProject(r int) string {
	return projectScopedPolicyConfig(r) + fmt.Sprintf(`
# export selected policies scoped to a project
resource "nutanix_network_security_policy_export_v2" "test" {
  policy_ext_ids = [nutanix_network_security_policy_v2.test-1.id]
  project_ext_id = nutanix_project_v2.test.ext_id
  file_path      = "%[1]s"
  depends_on     = [nutanix_network_security_policy_v2.test-1]
}
`, exportSelectedProjectFilePath)
}

func testAccNetworkSecurityPolicyExportV2ConfigAllWithProject(r int) string {
	return projectScopedPolicyConfig(r) + fmt.Sprintf(`
# export all policies scoped to a project (policy_ext_ids omitted)
resource "nutanix_network_security_policy_export_v2" "test" {
  project_ext_id = nutanix_project_v2.test.ext_id
  file_path      = "%[1]s"
  depends_on     = [nutanix_network_security_policy_v2.test-1]
}
`, exportAllProjectFilePath)
}
