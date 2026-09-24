package microsegv2_test

/*
Test Plan — nutanix_network_security_policy_import_v2

This is an action resource that triggers bulk import of network security
policies from a data file. It is backed by an async task and has no
GetById Read endpoint.

Scenarios:
1. End-to-end lifecycle without project scope — provision a network security
   policy (and its category), export it to a file, destroy the original policy,
   then import the exported file (ntnx_project_ext_id unset) and verify the policy
   is recreated (imported_policy_ext_ids is populated). A dry-run import is also
   exercised to confirm it creates nothing. The exported file and imported
   entities are cleaned up.
2. End-to-end lifecycle with project scope — same round-trip but the policy is
   scoped to a project and the import is performed with ntnx_project_ext_id set.
3. Missing required field (path) — expect "Missing required argument".
4. Invalid path — supply a path that is not a regular file and expect a
   "not a valid file" validation error.
5. Non-existent file — supply a missing file path and expect the same
   "not a valid file" validation error.

// Import not supported for this resource.
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

const importExportFilePath = "/tmp/test_policy_import_export.bin"
const importExportProjectFilePath = "/tmp/test_policy_import_export_project.bin"

// TestAccV2NutanixNetworkSecurityPolicyImport_Basic exercises the full export/import
// round-trip without project scope: create a policy, export it to a file, destroy the
// original policy, and import the file back so the policy is recreated on the cluster.
func TestAccV2NutanixNetworkSecurityPolicyImport_Basic(t *testing.T) {
	t.Cleanup(func() { _ = os.Remove(importExportFilePath) })
	n := acctest.RandIntBetween(1, 100000)

	// Captured during the import step so the teardown step can verify these
	// entities were actually removed from the cluster. The import recreates the
	// category from the exported file outside of Terraform's state, so it is located
	// (and deleted) by name (key) rather than by a captured ext_id.
	var importedPolicyIDs []string
	categoryKey := fmt.Sprintf("tf-test-cat-%d", n)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testNetworkSecurityPolicyImportV2CheckDestroy,
		Steps: []resource.TestStep{
			// 1. Create the network security policy (and its category) and export it to a file.
			{
				Config: testAccNetworkSecurityPolicyImportV2ConfigExport(n),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("nutanix_network_security_policy_export_v2.export", "task_ext_id"),
					resource.TestCheckResourceAttr("nutanix_network_security_policy_export_v2.export", "file_path", importExportFilePath),
					testAccCheckExportFileValid("nutanix_network_security_policy_export_v2.export"),
				),
			},
			// 2. Delete everything from step 1 (policy, export action, and category) so the
			//    cluster is clean before the import. The exported file persists on disk.
			{
				Config: testAccNetworkSecurityPolicyImportV2ConfigNoResources(),
			},
			// 3. Dry-run import: the cluster only validates the file and must not create any
			//    entities. No policies are tracked and the category must not exist afterwards.
			{
				Config: testAccNetworkSecurityPolicyImportV2ConfigImport(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNetworkSecurityPolicyImportV2, "task_ext_id"),
					resource.TestCheckResourceAttr(resourceNameNetworkSecurityPolicyImportV2, "dryrun", "true"),
					resource.TestCheckResourceAttr(resourceNameNetworkSecurityPolicyImportV2, "imported_policy_ext_ids.#", "0"),
					testAccCheckImportEntitiesNotCreated(categoryKey),
				),
			},
			// 4. Real import: re-creates the policy (and its category) on the cluster, then
			//    captures the imported policy ext_ids for the teardown verification.
			{
				Config: testAccNetworkSecurityPolicyImportV2ConfigImport(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNetworkSecurityPolicyImportV2, "task_ext_id"),
					resource.TestCheckResourceAttr(resourceNameNetworkSecurityPolicyImportV2, "path", importExportFilePath),
					// No project scope: ntnx_project_ext_id is not set.
					resource.TestCheckNoResourceAttr(resourceNameNetworkSecurityPolicyImportV2, "ntnx_project_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameNetworkSecurityPolicyImportV2, "imported_policy_ext_ids.0"),
					testAccCaptureImportedPolicyIDs(resourceNameNetworkSecurityPolicyImportV2, &importedPolicyIDs),
				),
			},
			// 5. Tear everything down (no resources left in config): the import resource deletes
			//    the imported policy first, then the category the import created (located by
			//    name) is removed. Verify both are gone.
			{
				Config: testAccNetworkSecurityPolicyImportV2ConfigNoResources(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImportedEntitiesDestroyed(&importedPolicyIDs, categoryKey),
				),
			},
		},
	})
}

// TestAccV2NutanixNetworkSecurityPolicyImport_WithProjectScope exercises the export/import
// round-trip for a policy scoped to a project. A project is created and kept alive across the
// steps so the import (performed with ntnx_project_ext_id set) can reference it.
func TestAccV2NutanixNetworkSecurityPolicyImport_WithProjectScope(t *testing.T) {
	t.Cleanup(func() { _ = os.Remove(importExportProjectFilePath) })
	n := acctest.RandIntBetween(1, 100000)

	var importedPolicyIDs []string
	categoryKey := fmt.Sprintf("tf-test-cat-%d", n)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testNetworkSecurityPolicyImportV2CheckDestroy,
		Steps: []resource.TestStep{
			// 1. Create a project, a category shared with it, and a policy scoped to the project,
			//    then export the policy to a file.
			{
				Config: testAccNSPImportV2ProjectConfigExport(n),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("nutanix_network_security_policy_export_v2.export", "task_ext_id"),
					resource.TestCheckResourceAttr("nutanix_network_security_policy_export_v2.export", "file_path", importExportProjectFilePath),
					testAccCheckExportFileValid("nutanix_network_security_policy_export_v2.export"),
				),
			},
			// 2. Destroy the policy and category but keep the project so the import can scope to it.
			{
				Config: testAccNSPImportV2ProjectConfigProjectOnly(n),
			},
			// 3. Import the exported file with ntnx_project_ext_id set; verify the policy is
			//    recreated and the project scope is recorded.
			{
				Config: testAccNSPImportV2ProjectConfigImport(n),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameNetworkSecurityPolicyImportV2, "task_ext_id"),
					resource.TestCheckResourceAttr(resourceNameNetworkSecurityPolicyImportV2, "path", importExportProjectFilePath),
					resource.TestCheckResourceAttrPair(resourceNameNetworkSecurityPolicyImportV2, "ntnx_project_ext_id", "nutanix_project_v2.test", "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameNetworkSecurityPolicyImportV2, "imported_policy_ext_ids.0"),
					testAccCaptureImportedPolicyIDs(resourceNameNetworkSecurityPolicyImportV2, &importedPolicyIDs),
				),
			},
			// 4. Destroy the import resource (which deletes the imported policy) but keep the
			//    project. Verify the policy is gone and remove the import-created category by
			//    name. The category is unshared from the project before deletion, otherwise the
			//    category (and later the project) cannot be deleted while the share exists.
			{
				Config: testAccNSPImportV2ProjectConfigProjectOnly(n),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImportedEntitiesDestroyed(&importedPolicyIDs, categoryKey),
				),
			},
			// 5. Tear down the project now that no categories are shared with it.
			{
				Config: testAccNetworkSecurityPolicyImportV2ConfigNoResources(),
			},
		},
	})
}

func TestAccV2NutanixNetworkSecurityPolicyImport_MissingPath(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNetworkSecurityPolicyImportV2ConfigMissingPath(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccV2NutanixNetworkSecurityPolicyImport_InvalidPath(t *testing.T) {
	missingFile := "/nonexistent/invalid/path.json"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNetworkSecurityPolicyImportV2Config(missingFile, false, ""),
				ExpectError: regexp.MustCompile(fmt.Sprintf("The provided path '%s' is not a valid file", regexp.QuoteMeta(missingFile))),
			},
		},
	})
}

// Importing from a non-existent file must fail with a clear "not a valid file" error.
func TestAccV2NutanixNetworkSecurityPolicyImport_NonExistentFile(t *testing.T) {
	missingFile := fmt.Sprintf("/tmp/missing_export_file_%d.bin", acctest.RandIntBetween(1, 100000))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNetworkSecurityPolicyImportV2Config(missingFile, false, ""),
				ExpectError: regexp.MustCompile(fmt.Sprintf("The provided path '%s' is not a valid file", regexp.QuoteMeta(missingFile))),
			},
		},
	})
}

func testAccNetworkSecurityPolicyImportV2Config(path string, purgePolicies bool, projectExtID string) string {
	projectLine := ""
	if projectExtID != "" {
		projectLine = fmt.Sprintf(`ntnx_project_ext_id = "%s"`, projectExtID)
	}
	return fmt.Sprintf(`
resource "nutanix_network_security_policy_import_v2" "test" {
  path               = "%s"
  ntnx_purge_policies = %t
  %s
}
`, path, purgePolicies, projectLine)
}

func testAccNetworkSecurityPolicyImportV2ConfigMissingPath() string {
	return `
resource "nutanix_network_security_policy_import_v2" "test" {
  ntnx_purge_policies = true
}
`
}

// nspImportTestCategory returns the category block shared across the end-to-end import
// steps. It must be byte-identical between steps so the category is not churned and the
// imported policy can re-reference it by the same ext_id.
func nspImportTestCategory(n int) string {
	return fmt.Sprintf(`
resource "nutanix_category_v2" "cat1" {
  key   = "tf-test-cat-%[1]d"
  value = "tf-test-cat-%[1]d-value"
}
`, n)
}

// Step 1: a network security policy plus an export action that writes the policy to a file.
func testAccNetworkSecurityPolicyImportV2ConfigExport(n int) string {
	return nspImportTestCategory(n) + fmt.Sprintf(`
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

resource "nutanix_network_security_policy_export_v2" "export" {
  policy_ext_ids = [nutanix_network_security_policy_v2.test-1.id]
  file_path      = "%[1]s"
  depends_on     = [nutanix_network_security_policy_v2.test-1]
}
`, importExportFilePath, n)
}

// Imports the file produced in step 1. With dryrun=true the cluster only validates the file
// and creates nothing; with dryrun=false it actually recreates the policy on the cluster.
func testAccNetworkSecurityPolicyImportV2ConfigImport(dryrun bool) string {
	return fmt.Sprintf(`
resource "nutanix_network_security_policy_import_v2" "test" {
  path                = "%[1]s"
  ntnx_purge_policies = false
  dryrun              = %[2]t
}
`, importExportFilePath, dryrun)
}

// testAccNetworkSecurityPolicyImportV2ConfigNoResources returns a configuration that declares
// no managed resources. Applying it destroys everything from the previous step. It is used both
// to clean the cluster before the import (deleting the policy and category) and to tear down the
// imported policy and category afterwards. The string is intentionally non-empty (a comment),
// because the test framework rejects a truly empty config with "Unsupported test mode".
func testAccNetworkSecurityPolicyImportV2ConfigNoResources() string {
	return "# no resources: applying this destroys everything created by the previous step\n"
}

// nspImportV2ProjectBlock returns the project block shared across the project-scoped import
// steps. It is kept byte-identical between steps so the project persists while the policy and
// category are destroyed and re-imported around it.
func nspImportV2ProjectBlock(n int) string {
	return fmt.Sprintf(`
resource "nutanix_project_v2" "test" {
  name        = "tf-proj-%[1]d"
  project_id  = "tf-proj-%[1]d"
  description = "import project scope test"
}
`, n)
}

// Step 1 (project scope): a project (kept alive for the later import), a category, a policy,
// and an export action that writes the policy to a file. The project scope is exercised on the
// import side (step 3), so here the category is not shared and the policy is not project-scoped
// to keep the pre-import teardown (step 2) clean.
func testAccNSPImportV2ProjectConfigExport(n int) string {
	return nspImportV2ProjectBlock(n) + fmt.Sprintf(`
resource "nutanix_category_v2" "cat1" {
  key   = "tf-test-cat-%[2]d"
  value = "tf-test-cat-%[2]d-value"
}

resource "nutanix_network_security_policy_v2" "test-1" {
  name        = "tf-nsp-test-%[2]d"
  description = "Network Security Policy Test for project scoped import"
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

resource "nutanix_network_security_policy_export_v2" "export" {
  policy_ext_ids = [nutanix_network_security_policy_v2.test-1.id]
  file_path      = "%[1]s"
  depends_on     = [nutanix_network_security_policy_v2.test-1]
}
`, importExportProjectFilePath, n)
}

// Step 2 (project scope): keep only the project so the import can reference it; the policy and
// category from step 1 are destroyed while the exported file persists on disk.
func testAccNSPImportV2ProjectConfigProjectOnly(n int) string {
	return nspImportV2ProjectBlock(n)
}

// Step 3 (project scope): import the exported file scoped to the project.
func testAccNSPImportV2ProjectConfigImport(n int) string {
	return nspImportV2ProjectBlock(n) + fmt.Sprintf(`
resource "nutanix_network_security_policy_import_v2" "test" {
  path                = "%[1]s"
  ntnx_purge_policies = false
  ntnx_project_ext_id = nutanix_project_v2.test.ext_id
  depends_on          = [nutanix_project_v2.test]
}
`, importExportProjectFilePath)
}
