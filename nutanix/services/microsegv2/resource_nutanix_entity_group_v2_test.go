package microsegv2_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccNutanixEntityGroupV2Resource_Basic(t *testing.T) {
	r := acctest.RandIntRange(1, 100)
	name := fmt.Sprintf("tf-entity-group-%d", r)
	description := fmt.Sprintf("tf-entity-group-%d_desc", r)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testEntityGroupV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntityGroupV2ResourceConfig(r, name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameEntityGroupV2, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameEntityGroupV2, "name", name),
					resource.TestCheckResourceAttr(resourceNameEntityGroupV2, "description", description),
					resource.TestCheckResourceAttr(resourceNameEntityGroupV2, "allowed_config.0.entities.0.selected_by", "CATEGORY_EXT_ID"),
					resource.TestCheckResourceAttr(resourceNameEntityGroupV2, "allowed_config.0.entities.0.type", "VM"),
					resource.TestCheckResourceAttr(resourceNameEntityGroupV2, "allowed_config.0.entities.0.reference_ext_ids.#", "2"),
					resource.TestCheckResourceAttrPair(resourceNameEntityGroupV2, "allowed_config.0.entities.0.reference_ext_ids.0", "nutanix_category_v2.categories.0", "id"),
					resource.TestCheckResourceAttrPair(resourceNameEntityGroupV2, "allowed_config.0.entities.0.reference_ext_ids.1", "nutanix_category_v2.categories.1", "id"),
					// Computed flags added in the Flow Management feature set.
					resource.TestCheckResourceAttrSet(resourceNameEntityGroupV2, "is_shared_with_all_projects"),
					resource.TestCheckResourceAttrSet(resourceNameEntityGroupV2, "is_system_defined"),
				),
			},
			{
				Config: testAccEntityGroupV2ResourceUpdateConfig(r, name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameEntityGroupV2, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameEntityGroupV2, "name", name+"-updated"),
					resource.TestCheckResourceAttr(resourceNameEntityGroupV2, "description", description+" updated"),
				),
			},
		},
	})
}

func TestAccNutanixEntityGroupV2Resource_WithoutName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccEntityGroupV2ResourceConfigWithoutName(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccNutanixEntityGroupV2Resource_WithWrongReferenceExtIds(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-entity-group-wrong-ref-%d", r)
	description := "entity_group_wrong_ref_ext_ids_desc"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccEntityGroupV2ResourceConfigWithWrongReferenceExtIds(name, description),
				ExpectError: regexp.MustCompile("categories do not exist in project"),
			},
		},
	})
}

func testAccEntityGroupV2ResourceConfig(r int, name, description string) string {
	return fmt.Sprintf(`
resource "nutanix_category_v2" "categories" {
  count       = 2
  key         = "tf_entity_group_%[1]d_${count.index}_key"
  value       = "tf_entity_group_%[1]d_${count.index}_value"
  description = "tf_entity_group_%[1]d_${count.index}_description"
}

resource "nutanix_entity_group_v2" "test" {
  name        = "%[2]s"
  description = "%[3]s"

  allowed_config {
    entities {
      type             = "VM"
      selected_by      = "CATEGORY_EXT_ID"
      reference_ext_ids = [
        nutanix_category_v2.categories[0].id,
        nutanix_category_v2.categories[1].id
      ]
    }
  }
}
`, r, name, description)
}

func testAccEntityGroupV2ResourceUpdateConfig(r int, name, description string) string {
	return fmt.Sprintf(`
resource "nutanix_category_v2" "categories" {
  count       = 2
  key         = "tf_entity_group_%[1]d_${count.index}_key"
  value       = "tf_entity_group_%[1]d_${count.index}_value"
  description = "tf_entity_group_%[1]d_${count.index}_description"
}

resource "nutanix_entity_group_v2" "test" {
  name        = "%s-updated"
  description = "%s updated"

  allowed_config {
    entities {
      type             = "VM"
      selected_by      = "CATEGORY_EXT_ID"
      reference_ext_ids = [
        nutanix_category_v2.categories[0].id,
        nutanix_category_v2.categories[1].id
      ]
    }
  }
}
`, r, name, description)
}

func testAccEntityGroupV2ResourceConfigWithoutName() string {
	return `
resource "nutanix_entity_group_v2" "test" {
  description = "entity_group_without_name_desc"
}
`
}

func TestAccNutanixEntityGroupV2Resource_ProjectAssociation(t *testing.T) {
	r := acctest.RandIntRange(101, 200)
	name := fmt.Sprintf("tf-eg-projassoc-%d", r)
	description := "entity group project association test"
	projectName := fmt.Sprintf("tf-eg-pa-proj-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testEntityGroupProjectAssociationConfig(r, name, description, projectName, "", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceNameEntityGroupV2, "project_ext_id", "nutanix_project_v2.test", "ext_id"),
					resource.TestCheckResourceAttrPair("data.nutanix_entity_group_v2.test", "project_ext_id", "nutanix_project_v2.test", "ext_id"),
					resource.TestCheckResourceAttr("data.nutanix_entity_groups_v2.test", "entity_groups.#", "1"),
					resource.TestCheckResourceAttrPair("data.nutanix_entity_groups_v2.test", "entity_groups.0.ext_id", resourceNameEntityGroupV2, "ext_id"),
					resource.TestCheckResourceAttrPair("data.nutanix_entity_groups_v2.test", "entity_groups.0.project_ext_id", "nutanix_project_v2.test", "ext_id"),
				),
			},
			{
				Config:      testEntityGroupProjectAssociationConfig(r, name, description, projectName, "00000000-0000-0000-0000-000000000000", true),
				ExpectError: regexp.MustCompile("Update of project_ext_id is not supported"),
			},
			{
				Config: testEntityGroupProjectAssociationCleanupConfig(r, projectName),
			},
		},
	})
}

// TestAccNutanixEntityGroupV2Resource_FqdnSelection covers the FQDN based
// selection added in the Flow Management feature set: an ADDRESS_GROUP entity
// selected BY FQDN_VALUES, populating the new "fqdns" list. The supported
// combination per the SDK is "ADDRESS_GROUP BY FQDN_VALUES". The update step
// changes the FQDN list to exercise the update path.
func TestAccNutanixEntityGroupV2Resource_FqdnSelection(t *testing.T) {
	r := acctest.RandIntRange(201, 300)
	name := fmt.Sprintf("tf-eg-fqdn-%d", r)
	description := "entity group fqdn selection test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testEntityGroupV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntityGroupV2FqdnConfig(name, description, `["example.com", "api.example.com"]`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameEntityGroupV2, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameEntityGroupV2, "name", name),
					resource.TestCheckResourceAttr(resourceNameEntityGroupV2, "allowed_config.0.entities.0.type", "ADDRESS_GROUP"),
					resource.TestCheckResourceAttr(resourceNameEntityGroupV2, "allowed_config.0.entities.0.selected_by", "FQDN_VALUES"),
					resource.TestCheckResourceAttr(resourceNameEntityGroupV2, "allowed_config.0.entities.0.fqdns.#", "2"),
					resource.TestCheckResourceAttr(resourceNameEntityGroupV2, "allowed_config.0.entities.0.fqdns.0", "example.com"),
					resource.TestCheckResourceAttr(resourceNameEntityGroupV2, "allowed_config.0.entities.0.fqdns.1", "api.example.com"),
				),
			},
			{
				Config: testAccEntityGroupV2FqdnConfig(name, description, `["example.org"]`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameEntityGroupV2, "allowed_config.0.entities.0.fqdns.#", "1"),
					resource.TestCheckResourceAttr(resourceNameEntityGroupV2, "allowed_config.0.entities.0.fqdns.0", "example.org"),
				),
			},
		},
	})
}

// TestAccNutanixEntityGroupV2Resource_RegexSelection covers the pattern based
// selection added in the Flow Management feature set: a VM entity selected BY
// REGEX, populating the new "reference_string" and "match_criteria". The
// supported combination per the SDK is "VM BY REGEX". Note reference_string is
// a plain literal, not a raw regex: the backend rejects special characters
// (" ' * , [ ] ? % $ - MIC-30525) and match_criteria (STARTS_WITH, CONTAINS,
// ...) supplies the matching semantics. The update step changes both values.
func TestAccNutanixEntityGroupV2Resource_RegexSelection(t *testing.T) {
	// REGEX (kVmNameByRegex) selection is only accepted by the backend when Flow
	// flex mode is enabled (MIC/NETWORKING: "selection types kVmByUuid,
	// kVmNameByRegex, and kSubnetByUuid are only allowed when flex mode is
	// enabled"). That requires a standalone (SMSP) Flow Controller deployment,
	// which the shared acceptance PC does not have, so gate this test behind an
	// explicit opt-in env var pointed at a flex-enabled PC.
	if os.Getenv("NUTANIX_FLEX_MODE_ENABLED") == "" {
		t.Skip("skipping: REGEX entity selection requires Flow flex mode (SMSP); set NUTANIX_FLEX_MODE_ENABLED to run against a flex-enabled PC")
	}

	r := acctest.RandIntRange(301, 400)
	name := fmt.Sprintf("tf-eg-regex-%d", r)
	description := "entity group regex selection test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testEntityGroupV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntityGroupV2RegexConfig(name, description, "web-prod", "STARTS_WITH"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameEntityGroupV2, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameEntityGroupV2, "name", name),
					resource.TestCheckResourceAttr(resourceNameEntityGroupV2, "allowed_config.0.entities.0.type", "VM"),
					resource.TestCheckResourceAttr(resourceNameEntityGroupV2, "allowed_config.0.entities.0.selected_by", "REGEX"),
					resource.TestCheckResourceAttr(resourceNameEntityGroupV2, "allowed_config.0.entities.0.reference_string", "web-prod"),
					resource.TestCheckResourceAttr(resourceNameEntityGroupV2, "allowed_config.0.entities.0.match_criteria", "STARTS_WITH"),
				),
			},
			{
				Config: testAccEntityGroupV2RegexConfig(name, description, "db-node", "CONTAINS"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameEntityGroupV2, "allowed_config.0.entities.0.reference_string", "db-node"),
					resource.TestCheckResourceAttr(resourceNameEntityGroupV2, "allowed_config.0.entities.0.match_criteria", "CONTAINS"),
				),
			},
		},
	})
}

func egProjectExtIDLine(override string) string {
	if override == "" {
		return `project_ext_id = nutanix_project_v2.test.ext_id`
	}
	return fmt.Sprintf(`project_ext_id = "%s"`, override)
}

func testEntityGroupProjectAssociationConfig(r int, name, description, projectName, projectExtIDOverride string, shareCategories bool) string {
	shareBlock := `shared_with_projects = []`
	if shareCategories {
		shareBlock = `shared_with_projects = [nutanix_project_v2.test.ext_id]`
	}
	return fmt.Sprintf(`
resource "nutanix_project_v2" "test" {
  name        = "%[4]s"
  project_id  = "%[4]s"
  description = "project association test"
}

resource "nutanix_category_v2" "categories" {
  count       = 2
  key         = "tf_eg_pa_%[1]d_${count.index}_key"
  value       = "tf_eg_pa_%[1]d_${count.index}_value"
  description = "tf_eg_pa_%[1]d_${count.index}_description"
  %[6]s
}

resource "nutanix_entity_group_v2" "test" {
  name        = "%[2]s"
  description = "%[3]s"

  allowed_config {
    entities {
      type              = "VM"
      selected_by       = "CATEGORY_EXT_ID"
      reference_ext_ids = [
        nutanix_category_v2.categories[0].id,
        nutanix_category_v2.categories[1].id
      ]
    }
  }
  %[5]s
  depends_on = [nutanix_project_v2.test]
}

data "nutanix_entity_group_v2" "test" {
  ext_id     = nutanix_entity_group_v2.test.ext_id
  depends_on = [nutanix_entity_group_v2.test]
}

data "nutanix_entity_groups_v2" "test" {
  filter     = "name eq '${nutanix_entity_group_v2.test.name}'"
  depends_on = [nutanix_entity_group_v2.test]
}
`, r, name, description, projectName, egProjectExtIDLine(projectExtIDOverride), shareBlock)
}

func testEntityGroupProjectAssociationCleanupConfig(r int, projectName string) string {
	return fmt.Sprintf(`
resource "nutanix_project_v2" "test" {
  name        = "%[2]s"
  project_id  = "%[2]s"
  description = "project association test"
}

resource "nutanix_category_v2" "categories" {
  count       = 2
  key         = "tf_eg_pa_%[1]d_${count.index}_key"
  value       = "tf_eg_pa_%[1]d_${count.index}_value"
  description = "tf_eg_pa_%[1]d_${count.index}_description"
  shared_with_projects = []
}
`, r, projectName)
}

func testAccEntityGroupV2FqdnConfig(name, description, fqdns string) string {
	return fmt.Sprintf(`
resource "nutanix_entity_group_v2" "test" {
  name        = "%[1]s"
  description = "%[2]s"

  allowed_config {
    entities {
      type        = "ADDRESS_GROUP"
      selected_by = "FQDN_VALUES"
      fqdns       = %[3]s
    }
  }
}
`, name, description, fqdns)
}

func testAccEntityGroupV2RegexConfig(name, description, referenceString, matchCriteria string) string {
	return fmt.Sprintf(`
resource "nutanix_entity_group_v2" "test" {
  name        = "%[1]s"
  description = "%[2]s"

  allowed_config {
    entities {
      type             = "VM"
      selected_by      = "REGEX"
      reference_string = "%[3]s"
      match_criteria   = "%[4]s"
    }
  }
}
`, name, description, referenceString, matchCriteria)
}

func testAccEntityGroupV2ResourceConfigWithWrongReferenceExtIds(name, description string) string {
	return fmt.Sprintf(`
resource "nutanix_entity_group_v2" "test" {
  name        = "%s"
  description = "%s"

  allowed_config {
    entities {
      type            = "VM"
      selected_by     = "CATEGORY_EXT_ID"
      reference_ext_ids = [
        "83cbf00d-782f-4efc-87c8-4129f5942aaa",
        "83cbf00d-782f-4efc-87c8-4129f5942bbb"
      ]
    }
  }
}
`, name, description)
}
