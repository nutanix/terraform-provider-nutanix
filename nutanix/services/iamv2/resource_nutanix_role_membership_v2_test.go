package iamv2_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/iam-go-client/v17/models/iam/v4/request/rolemembership"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestAccNutanixRoleMembershipV2Resource_basic(t *testing.T) {
	resourceNameProjectAdmin := "nutanix_role_membership_v2.project_admin_role"
	resourceNameDeveloper := "nutanix_role_membership_v2.developer_role"
	datasourceRoleMembershipSummary := "data.nutanix_role_membership_summary_v2.get_role_membership_summary"
	datasourcRoleMembershipswithFilter := "data.nutanix_role_memberships_v2.get_role_memberships"
	datasourceRoleMembership := "data.nutanix_role_membership_v2.get_role_membership_by_id"

	secondaryAD := testVars.Iam.DirectoryServicesMain.SecondaryAD
	userExtID := secondaryAD.DomainUsersUsergroups.Users["ssptest1@qa.nutanix.com"]
	userGroupExtID := secondaryAD.DomainUsersUsergroups.UserGroups["dnd_approval_group_1"]

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixRoleMembershipV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixRoleMembershipV2Config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameProjectAdmin, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameProjectAdmin, "role_ext_id"),
					resource.TestCheckResourceAttr(resourceNameProjectAdmin, "identity_type", "USER"),
					resource.TestCheckResourceAttr(
						resourceNameProjectAdmin, "identity_ext_id",
						userExtID,
					),
					resource.TestCheckResourceAttrPair(
						resourceNameProjectAdmin, "project_ext_id",
						"nutanix_project_v2.test", "ext_id",
					),
					resource.TestCheckResourceAttrSet(resourceNameProjectAdmin, "created_by"),
					resource.TestCheckResourceAttrSet(resourceNameProjectAdmin, "created_time"),
					resource.TestCheckResourceAttrSet(resourceNameProjectAdmin, "last_updated_time"),
					resource.TestCheckResourceAttrSet(resourceNameDeveloper, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameDeveloper, "role_ext_id"),
					resource.TestCheckResourceAttr(resourceNameDeveloper, "identity_type", "GROUP"),
					resource.TestCheckResourceAttr(
						resourceNameDeveloper, "identity_ext_id",
						userGroupExtID,
					),
					resource.TestCheckResourceAttrPair(
						resourceNameDeveloper, "project_ext_id",
						"nutanix_project_v2.test", "ext_id",
					),
					resource.TestCheckResourceAttrSet(resourceNameDeveloper, "created_by"),
					resource.TestCheckResourceAttrSet(resourceNameDeveloper, "created_time"),
					resource.TestCheckResourceAttrSet(resourceNameDeveloper, "last_updated_time"),

					// Validate the role memberships with filter. The list endpoint does
					// not guarantee ordering, so assert the expected USER and GROUP
					// entries exist without depending on their index.
					resource.TestCheckResourceAttrSet(datasourcRoleMembershipswithFilter, "role_memberships.#"),
					resource.TestCheckResourceAttr(datasourcRoleMembershipswithFilter, "role_memberships.#", "2"),
					checkRoleMembershipsContain(
						datasourcRoleMembershipswithFilter, "nutanix_project_v2.test",
						userExtID, userGroupExtID,
					),

					// Validate the rolemembership Summary
					resource.TestCheckResourceAttrSet(datasourceRoleMembershipSummary, "summaries.#"),
					resource.TestCheckResourceAttr(datasourceRoleMembershipSummary, "summaries.#", "1"),
					resource.TestCheckResourceAttr(datasourceRoleMembershipSummary, "summaries.0.users_count", "1"),
					resource.TestCheckResourceAttr(datasourceRoleMembershipSummary, "summaries.0.groups_count", "1"),
					resource.TestCheckResourceAttr(datasourceRoleMembershipSummary, "summaries.0.roles_count", "2"),
					resource.TestCheckResourceAttr(datasourceRoleMembershipSummary, "summaries.0.total_identities_count", "2"),

					// Validate the rolemembership by id
					resource.TestCheckResourceAttr(datasourceRoleMembership, "identity_type", "USER"),
					resource.TestCheckResourceAttr(
						datasourceRoleMembership, "identity_ext_id",
						userExtID,
					),
					resource.TestCheckResourceAttrPair(
						datasourceRoleMembership, "project_ext_id",
						"nutanix_project_v2.test", "ext_id",
					),
				),
			},
		},
	})
}

func TestAccNutanixRoleMembershipV2Resource_ListWithInvalidFilter(t *testing.T) {
	datasourceList := "data.nutanix_role_memberships_v2.get_role_memberships"
	randomUUID := utils.GenUUID()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixRoleMembershipV2ListWithInvalidFilterConfig(randomUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceList, "role_memberships.#", "0"),
				),
			},
		},
	})
}

func testAccNutanixRoleMembershipV2ListWithInvalidFilterConfig(uuid string) string {
	return fmt.Sprintf(`
	data "nutanix_role_memberships_v2" "get_role_memberships" {
		filter = "projectExtId eq '%s'"
	}
`, uuid)
}

func TestAccNutanixRoleMembershipV2Resource_SummaryWithInvalidFilter(t *testing.T) {
	datasourceSummary := "data.nutanix_role_membership_summary_v2.get_role_membership_summary"
	randomUUID := utils.GenUUID()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixRoleMembershipV2SummaryWithInvalidFilterConfig(randomUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceSummary, "summaries.#", "0"),
				),
			},
		},
	})
}

func testAccNutanixRoleMembershipV2SummaryWithInvalidFilterConfig(uuid string) string {
	return fmt.Sprintf(`
	data "nutanix_role_membership_summary_v2" "get_role_membership_summary" {
		filter = "extId eq '%s'"
	}
`, uuid)
}

// checkRoleMembershipsContain verifies the role_memberships list data source
// contains exactly one USER membership (with userExtID) and one GROUP membership
// (with groupExtID), all scoped to the given project, without relying on the
// list ordering returned by the API.
func checkRoleMembershipsContain(dsName, projectResName, userExtID, groupExtID string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[dsName]
		if !ok {
			return fmt.Errorf("data source not found: %s", dsName)
		}
		project, ok := s.RootModule().Resources[projectResName]
		if !ok {
			return fmt.Errorf("project resource not found: %s", projectResName)
		}
		attrs := ds.Primary.Attributes
		projectExtID := project.Primary.Attributes["ext_id"]

		foundUser, foundGroup := false, false
		for i := 0; i < 2; i++ {
			prefix := fmt.Sprintf("role_memberships.%d.", i)
			identityType := attrs[prefix+"identity_type"]
			identityExtID := attrs[prefix+"identity_ext_id"]
			if got := attrs[prefix+"project_ext_id"]; got != projectExtID {
				return fmt.Errorf("role_memberships.%d.project_ext_id = %q, want %q", i, got, projectExtID)
			}
			switch identityType {
			case "USER":
				if identityExtID != userExtID {
					return fmt.Errorf("USER membership identity_ext_id = %q, want %q", identityExtID, userExtID)
				}
				foundUser = true
			case "GROUP":
				if identityExtID != groupExtID {
					return fmt.Errorf("GROUP membership identity_ext_id = %q, want %q", identityExtID, groupExtID)
				}
				foundGroup = true
			default:
				return fmt.Errorf("unexpected identity_type %q at role_memberships.%d", identityType, i)
			}
		}
		if !foundUser || !foundGroup {
			return fmt.Errorf("expected one USER and one GROUP membership, got foundUser=%v foundGroup=%v", foundUser, foundGroup)
		}
		return nil
	}
}

func testAccCheckNutanixRoleMembershipV2Destroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_role_membership_v2" {
			continue
		}
		getRequest := import1.GetRoleMembershipByIdRequest{
			ExtId: utils.StringPtr(rs.Primary.ID),
		}
		_, err := conn.IamAPI.RoleMembershipAPIInstance.GetRoleMembershipById(ctx, &getRequest)
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "not found") || strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
				return nil
			}
		}
		return fmt.Errorf("Role membership still exists: %s", rs.Primary.ID)
	}
	return nil
}

func testAccNutanixRoleMembershipV2Config() string {
	secondaryAD := testVars.Iam.DirectoryServicesMain.SecondaryAD
	userExtID := secondaryAD.DomainUsersUsergroups.Users["ssptest1@qa.nutanix.com"]
	userGroupExtID := secondaryAD.DomainUsersUsergroups.UserGroups["dnd_approval_group_1"]
	idpExtID := secondaryAD.ExtID

	return fmt.Sprintf(`
	data "nutanix_roles_v2" "roles" {}

	locals {
	  project_admin_role_ext_id = [
    for role in data.nutanix_roles_v2.roles.roles :
    role.ext_id if role.display_name == "Project Admin"
  ][0]
	  developer_role_ext_id = [
    for role in data.nutanix_roles_v2.roles.roles :
    role.ext_id if role.display_name == "Developer"
  ][0]
	}

	resource "nutanix_project_v2" "test" {
		name = "test"
		description = "test"
	}

	resource "nutanix_role_membership_v2" "project_admin_role" {
		role_ext_id      = local.project_admin_role_ext_id
		identity_type    = "USER"
		identity_ext_id  = "%s"
		idp_ext_id       = "%s"
		project_ext_id   = nutanix_project_v2.test.ext_id
		scope_template_name = "ProjectsScopeTemplate"
		scope_template_name_values {
			name = "projectExtId"
			value = nutanix_project_v2.test.ext_id
		}
	}
	
	resource "nutanix_role_membership_v2" "developer_role" {
		role_ext_id      = local.developer_role_ext_id
		identity_type    = "GROUP"
		identity_ext_id  = "%s"
		idp_ext_id       = "%s"
		project_ext_id   = nutanix_project_v2.test.ext_id
		scope_template_name = "ProjectsScopeTemplate"
		scope_template_name_values {
			name = "projectExtId"
			value = nutanix_project_v2.test.ext_id
		}
	}

	data "nutanix_role_membership_v2" "get_role_membership_by_id" {
		ext_id = nutanix_role_membership_v2.project_admin_role.ext_id
	}
	
	data "nutanix_role_memberships_v2" "get_role_memberships" {
		filter = "projectExtId eq '${nutanix_project_v2.test.ext_id}'"
		depends_on = [nutanix_role_membership_v2.project_admin_role, nutanix_role_membership_v2.developer_role]
	}

	data "nutanix_role_membership_summary_v2" "get_role_membership_summary" {
		filter = "extId eq '${nutanix_project_v2.test.ext_id}'"
		depends_on = [nutanix_role_membership_v2.project_admin_role, nutanix_role_membership_v2.developer_role]
	}
`, userExtID, idpExtID, userGroupExtID, idpExtID)
}
