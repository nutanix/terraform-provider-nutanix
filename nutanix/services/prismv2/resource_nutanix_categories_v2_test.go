package prismv2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameCategory = "nutanix_category_v2.test"

func TestAccV2NutanixCategoryResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	value := fmt.Sprintf("test category value-%d", r)
	desc := "test category description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCategoryV2Config(r, value, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameCategory, "key", fmt.Sprintf("test-cat-%d", r)),
					resource.TestCheckResourceAttr(resourceNameCategory, "value", value),
					resource.TestCheckResourceAttr(resourceNameCategory, "description", desc),
					resource.TestCheckResourceAttr(resourceNameCategory, "type", "USER"),
					resource.TestCheckResourceAttr(resourceNameCategory, "shared_with_all_projects", "false"),
				),
			},
		},
	})
}

func TestAccV2NutanixCategoryResource_Update(t *testing.T) {
	r := acctest.RandInt()
	value := fmt.Sprintf("test category value-%d", r)
	desc := "test category description"
	updatedValue := fmt.Sprintf("test category value updated-%d", r)
	updateDesc := "test category description updated"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCategoryV2Config(r, value, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameCategory, "key", fmt.Sprintf("test-cat-%d", r)),
					resource.TestCheckResourceAttr(resourceNameCategory, "value", value),
					resource.TestCheckResourceAttr(resourceNameCategory, "description", desc),
					resource.TestCheckResourceAttr(resourceNameCategory, "type", "USER"),
					resource.TestCheckResourceAttr(resourceNameCategory, "shared_with_all_projects", "false"),
				),
			},
			{
				Config: testAccCategoryV2Config(r, updatedValue, updateDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameCategory, "key", fmt.Sprintf("test-cat-%d", r)),
					resource.TestCheckResourceAttr(resourceNameCategory, "value", updatedValue),
					resource.TestCheckResourceAttr(resourceNameCategory, "description", updateDesc),
					resource.TestCheckResourceAttr(resourceNameCategory, "type", "USER"),
					resource.TestCheckResourceAttr(resourceNameCategory, "shared_with_all_projects", "false"),
				),
			},
		},
	})
}

func TestAccV2NutanixCategoryResource_WithNoKey(t *testing.T) {
	r := acctest.RandInt()
	value := fmt.Sprintf("test category value-%d", r)
	desc := "test category description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCategoryV2ConfigWithNoKey(r, value, desc),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

const datasourceNameCategoryProject = "data.nutanix_category_v2.project_test"
const datasourceNameCategoryProjectList = "data.nutanix_categories_v2.list_test"
const defaultProjectUUID = "00000000-0000-0000-0000-000000000000"

func TestAccV2NutanixCategoryResource_DefaultProjectAndSharing(t *testing.T) {
	r := acctest.RandInt()
	value := fmt.Sprintf("test-cat-proj-val-%d", r)
	desc := "category for project sharing test"
	updatedDesc := "category updated with sharing"
	unsharedDesc := "category after unsharing"
	projectName := fmt.Sprintf("tf-cat-share-proj-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCategoryV2ProjectConfig(r, value, desc, projectName, "none"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameCategory, "key", fmt.Sprintf("test-cat-%d", r)),
					resource.TestCheckResourceAttr(resourceNameCategory, "value", value),
					resource.TestCheckResourceAttr(resourceNameCategory, "description", desc),
					resource.TestCheckResourceAttr(resourceNameCategory, "project_ext_id", defaultProjectUUID),
					resource.TestCheckResourceAttr(resourceNameCategory, "shared_with_all_projects", "false"),
					resource.TestCheckResourceAttr(datasourceNameCategoryProject, "project_ext_id", defaultProjectUUID),
					resource.TestCheckResourceAttr(datasourceNameCategoryProject, "shared_with_all_projects", "false"),
				),
			},
			{
				Config: testAccCategoryV2ProjectConfig(r, value, updatedDesc, projectName, "share"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameCategory, "description", updatedDesc),
					resource.TestCheckResourceAttr(resourceNameCategory, "shared_with_projects.#", "1"),
					resource.TestCheckResourceAttrPair(resourceNameCategory, "shared_with_projects.0", "nutanix_project_v2.share_project", "ext_id"),
					resource.TestCheckResourceAttr(resourceNameCategory, "shared_with_all_projects", "false"),
					resource.TestCheckResourceAttr(datasourceNameCategoryProject, "shared_with_projects.#", "1"),
					resource.TestCheckResourceAttrPair(datasourceNameCategoryProject, "shared_with_projects.0", "nutanix_project_v2.share_project", "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameCategoryProject, "shared_with_all_projects", "false"),
					resource.TestCheckResourceAttr(datasourceNameCategoryProjectList, "categories.#", "1"),
					resource.TestCheckResourceAttrPair(datasourceNameCategoryProjectList, "categories.0.ext_id", resourceNameCategory, "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameCategoryProjectList, "categories.0.shared_with_projects.#", "1"),
					resource.TestCheckResourceAttrPair(datasourceNameCategoryProjectList, "categories.0.shared_with_projects.0", "nutanix_project_v2.share_project", "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameCategoryProjectList, "categories.0.shared_with_all_projects", "false"),
				),
			},
			{
				Config: testAccCategoryV2ProjectConfig(r, value, unsharedDesc, projectName, "empty"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameCategory, "description", unsharedDesc),
					resource.TestCheckResourceAttr(resourceNameCategory, "shared_with_projects.#", "0"),
					resource.TestCheckResourceAttr(resourceNameCategory, "shared_with_all_projects", "false"),
					resource.TestCheckResourceAttr(datasourceNameCategoryProject, "shared_with_projects.#", "0"),
					resource.TestCheckResourceAttr(datasourceNameCategoryProject, "shared_with_all_projects", "false"),
					resource.TestCheckResourceAttr(datasourceNameCategoryProjectList, "categories.#", "0"),
				),
			},
		},
	})
}

func TestAccV2NutanixCategoryResource_NonDefaultProjectSharingFails(t *testing.T) {
	r := acctest.RandInt()
	value := fmt.Sprintf("test-cat-nondfl-val-%d", r)
	desc := "category in non-default project"
	project1Name := fmt.Sprintf("tf-cat-proj1-%d", r)
	project2Name := fmt.Sprintf("tf-cat-proj2-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCategoryV2NonDefaultProjectConfig(r, value, desc, project1Name, project2Name, "none"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameCategory, "ext_id"),
					resource.TestCheckResourceAttrPair(resourceNameCategory, "project_ext_id", "nutanix_project_v2.project1", "ext_id"),
				),
			},
			{
				Config:      testAccCategoryV2NonDefaultProjectConfig(r, value, desc, project1Name, project2Name, "share"),
				ExpectError: regexp.MustCompile("Failed to share or unshare category due to category not being owned by the default project"),
			},
		},
	})
}

func testAccCategoryV2ProjectConfig(r int, val, desc, projectName, shareState string) string {
	shareBlock := ""
	switch shareState {
	case "share":
		shareBlock = `shared_with_projects = [nutanix_project_v2.share_project.ext_id]`
	case "empty":
		shareBlock = `shared_with_projects = []`
	}
	listDataSource := ""
	if shareState != "none" {
		listDataSource = `
	data "nutanix_categories_v2" "list_test" {
		filter     = "sharedWithProjects/any(p:p eq '${nutanix_project_v2.share_project.id}')"
		depends_on = [nutanix_category_v2.test]
	}`
	}
	return fmt.Sprintf(`
	resource "nutanix_project_v2" "share_project" {
		name       = "%[4]s"
		project_id = "%[4]s"
		description = "project for category sharing test"
	}

	resource "nutanix_category_v2" "test" {
		key         = "test-cat-%[1]d"
		value       = "%[2]s"
		description = "%[3]s"
		%[5]s
		depends_on = [nutanix_project_v2.share_project]
	}

	data "nutanix_category_v2" "project_test" {
		ext_id = nutanix_category_v2.test.id
	}
	%[6]s
`, r, val, desc, projectName, shareBlock, listDataSource)
}

func testAccCategoryV2NonDefaultProjectConfig(r int, val, desc, proj1, proj2, shareState string) string {
	shareBlock := ""
	switch shareState {
	case "share":
		shareBlock = `shared_with_projects = [nutanix_project_v2.project2.ext_id]`
	case "empty":
		shareBlock = `shared_with_projects = []`
	}
	return fmt.Sprintf(`
	resource "nutanix_project_v2" "project1" {
		name       = "%[4]s"
		project_id = "%[4]s"
		description = "first project for category test"
	}

	resource "nutanix_project_v2" "project2" {
		name       = "%[5]s"
		project_id = "%[5]s"
		description = "second project for category test"
	}

	resource "nutanix_category_v2" "test" {
		key            = "test-cat-%[1]d"
		value          = "%[2]s"
		description    = "%[3]s"
		project_ext_id = nutanix_project_v2.project1.ext_id
		%[6]s
		depends_on = [nutanix_project_v2.project1, nutanix_project_v2.project2]
	}
`, r, val, desc, proj1, proj2, shareBlock)
}

func testAccCategoryV2Config(r int, val, desc string) string {
	return fmt.Sprintf(`
	resource "nutanix_category_v2" "test" {
		key = "test-cat-%d"
		value = "%[2]s"
		description = "%[3]s"
	}
`, r, val, desc)
}

func testAccCategoryV2ConfigWithNoKey(r int, val, desc string) string {
	return fmt.Sprintf(`
	resource "nutanix_category_v2" "test" {
		value = "%[2]s"
		description = "%[3]s"
	  }
`, r, val, desc)
}
