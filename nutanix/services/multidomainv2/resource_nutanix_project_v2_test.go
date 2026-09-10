package multidomainv2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestAccV2NutanixProjectResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-project-%d", r)
	description := "terraform test project CRUD"
	updateDescription := "terraform test project CRUD update"
	resourceNameResourceGroupsV2 := "data.nutanix_resource_groups_v2.get_resource_groups_valid_filter"
	resourceNameRoleMembershipsV2 := "data.nutanix_role_memberships_v2.get_role_memberships_valid_filter"
	resourceNameRoleMembershipSummaryV2 := "data.nutanix_role_membership_summary_v2.get_role_membership_summary_valid_filter"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testProjectV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectV2ResourceConfig(name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameProjectV2, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameProjectV2, "name", name),
					resource.TestCheckResourceAttr(resourceNameProjectV2, "description", description),
					resource.TestCheckResourceAttr(resourceNameProjectV2, "project_id", name),
					resource.TestCheckResourceAttr(resourceNameProjectV2, "is_system_defined", "false"),
					resource.TestCheckResourceAttr(resourceNameProjectV2, "is_default", "false"),
					resource.TestCheckResourceAttr(resourceNameProjectV2, "state", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceNameResourceGroupsV2, "resource_groups.#", "0"),
					resource.TestCheckResourceAttr(resourceNameRoleMembershipsV2, "role_memberships.#", "0"),
					resource.TestCheckResourceAttr(resourceNameRoleMembershipSummaryV2, "summaries.#", "0"),
				),
			},
			{
				Config: testAccProjectV2ResourceConfig(name, updateDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameProjectV2, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameProjectV2, "description", updateDescription),
				),
			},
			{
				Config:      testAccProjectV2ResourceConfig(name+"-updated", updateDescription),
				ExpectError: regexp.MustCompile("Update of name is not supported"),
			},
		},
	})
}

func TestAccV2NutanixProjectResource_DataSource(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-project-%d", r)
	description := "terraform test project datasource"
	datasourceGet := "data.nutanix_project_v2.project_get"
	datasourceList := "data.nutanix_projects_v2.projects_list"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectV2ResourceDataSourceConfig(name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceGet, "ext_id"),
					resource.TestCheckResourceAttr(datasourceGet, "name", name),
					resource.TestCheckResourceAttr(datasourceGet, "description", description),
					resource.TestCheckResourceAttr(datasourceGet, "project_id", name),
					resource.TestCheckResourceAttr(datasourceGet, "is_system_defined", "false"),
					resource.TestCheckResourceAttr(datasourceGet, "is_default", "false"),
					resource.TestCheckResourceAttr(datasourceGet, "state", "ACTIVE"),
					resource.TestCheckResourceAttr(datasourceList, "projects.#", "1"),
					resource.TestCheckResourceAttr(datasourceList, "projects.0.name", name),
					resource.TestCheckResourceAttr(datasourceList, "projects.0.description", description),
					resource.TestCheckResourceAttr(datasourceList, "projects.0.project_id", name),
					resource.TestCheckResourceAttr(datasourceList, "projects.0.is_system_defined", "false"),
					resource.TestCheckResourceAttr(datasourceList, "projects.0.is_default", "false"),
					resource.TestCheckResourceAttr(datasourceList, "projects.0.state", "ACTIVE"),
				),
			},
		},
	})
}
func TestAccV2NutanixProjectResource_ListWithInvalidFilter(t *testing.T) {
	datasourceList := "data.nutanix_projects_v2.projects_list"
	randomUUID := utils.GenUUID()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectV2ListWithInvalidFilterConfig(randomUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceList, "projects.#", "0"),
				),
			},
		},
	})
}

func testAccProjectV2ListWithInvalidFilterConfig(uuid string) string {
	return fmt.Sprintf(`
	data "nutanix_projects_v2" "projects_list" {
		filter = "extId eq '%s'"
	}
`, uuid)
}

func testAccProjectV2ResourceConfig(name, description string) string {
	return fmt.Sprintf(`
	resource "nutanix_project_v2" "test" {
		name        = "%s"
		project_id = "%s" # Id same as name for testing
		description = "%s"
	}

	data "nutanix_resource_groups_v2" "get_resource_groups_valid_filter" {
		filter = "projectExtId eq '${nutanix_project_v2.test.ext_id}'"
	}

	data "nutanix_role_memberships_v2" "get_role_memberships_valid_filter" {
		filter = "projectExtId eq '${nutanix_project_v2.test.ext_id}'"
	}

	data "nutanix_role_membership_summary_v2" "get_role_membership_summary_valid_filter" {
		filter = "extId eq '${nutanix_project_v2.test.ext_id}'"
	}
`, name, name, description)
}

func testAccProjectV2ResourceDataSourceConfig(name, description string) string {
	return fmt.Sprintf(`
	resource "nutanix_project_v2" "project_create" {
		name        = "%s"
		project_id = "%s" # Id same as name for testing
		description = "%s"
	}

	data "nutanix_project_v2" "project_get" {
		ext_id = nutanix_project_v2.project_create.id
		depends_on = [nutanix_project_v2.project_create]
	}

	data "nutanix_projects_v2" "projects_list" {
		filter = "extId eq '${nutanix_project_v2.project_create.id}'"
		depends_on = [nutanix_project_v2.project_create]
	}
`, name, name, description)
}
