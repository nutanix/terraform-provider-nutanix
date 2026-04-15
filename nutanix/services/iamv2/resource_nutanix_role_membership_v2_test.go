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
	datasourcRoleMembershipswithFilter := "data.nutanix_role_memberships_v2.get_role_memberships_with_filter"
	datasourceRoleMembership := "data.nutanix_role_membership_v2.get_role_membership_by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acc.TestAccPreCheck(t) },
		Providers:         acc.TestAccProviders,
		CheckDestroy:      testAccCheckNutanixRoleMembershipV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixRoleMembershipV2Config(),
				Check: resource.ComposeTestCheckFunc(
					// Dump entire state for debugging
						func(s *terraform.State) error {
							for name, rs := range s.RootModule().Resources {
									t.Logf("Resource: %s (Type: %s, ID: %s)", name, rs.Type, rs.Primary.ID)
									for k, v := range rs.Primary.Attributes {
											t.Logf("  %s = %s", k, v)
									}
							}
							return nil
					},
					resource.TestCheckResourceAttrSet(resourceNameProjectAdmin, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameProjectAdmin, "role_ext_id"),
					resource.TestCheckResourceAttr(resourceNameProjectAdmin, "identity_type", "USER"),
					resource.TestCheckResourceAttrPair(
						resourceNameProjectAdmin, "identity_ext_id",
						"nutanix_users_v2.user", "ext_id",
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
					resource.TestCheckResourceAttrPair(
						resourceNameDeveloper, "identity_ext_id",
						"nutanix_user_groups_v2.usergroup", "ext_id",
					),
					resource.TestCheckResourceAttrPair(
						resourceNameDeveloper, "project_ext_id",
						"nutanix_project_v2.test", "ext_id",
					),
					resource.TestCheckResourceAttrSet(resourceNameDeveloper, "created_by"),
					resource.TestCheckResourceAttrSet(resourceNameDeveloper, "created_time"),
					resource.TestCheckResourceAttrSet(resourceNameDeveloper, "last_updated_time"),
          
					// Validate the role memberships with filter
					resource.TestCheckResourceAttrSet(datasourcRoleMembershipswithFilter, "role_memberships.#"),
					resource.TestCheckResourceAttr(datasourcRoleMembershipswithFilter, "role_memberships.#", "2"),
					resource.TestCheckResourceAttr(datasourcRoleMembershipswithFilter, "role_memberships.0.identity_type",
						"USER",
					),
					resource.TestCheckResourceAttr(
						datasourcRoleMembershipswithFilter, "role_memberships.1.identity_type", "GROUP",
					),
					resource.TestCheckResourceAttrPair(
						datasourcRoleMembershipswithFilter, "role_memberships.0.identity_ext_id",
						"nutanix_users_v2.user", "ext_id",
					),
					resource.TestCheckResourceAttrPair(
						datasourcRoleMembershipswithFilter, "role_memberships.1.identity_ext_id",
						"nutanix_user_groups_v2.usergroup", "ext_id",
					),
					resource.TestCheckResourceAttrPair(
						datasourcRoleMembershipswithFilter, "role_memberships.0.project_ext_id",
						"nutanix_project_v2.test", "ext_id",
					),
					resource.TestCheckResourceAttrPair(
						datasourcRoleMembershipswithFilter, "role_memberships.1.project_ext_id",
						"nutanix_project_v2.test", "ext_id",
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
					resource.TestCheckResourceAttrPair(
						datasourceRoleMembership, "identity_ext_id",
						"nutanix_users_v2.user", "ext_id",
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

func TestAccNutanixRoleMembershipV2Resource_import(t *testing.T) {
	resourceName := "nutanix_role_membership_v2.project_admin_role"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acc.TestAccPreCheck(t) },
		CheckDestroy:      testAccCheckNutanixRoleMembershipV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixRoleMembershipV2Config(),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
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
	return fmt.Sprintf(`
	data "nutanix_roles_v2" "roles" {}

	locals {
		config = jsondecode(file("%s"))
	  project_admin_role_ext_id = [
    for role in data.nutanix_roles_v2.roles.roles :
    role.ext_id if role.display_name == "Project Admin"
  ][0]
	  developer_role_ext_id = [
    for role in data.nutanix_roles_v2.roles.roles :
    role.ext_id if role.display_name == "Developer"
  ][0]
	  idp_ext_id = local.config.iam.users.directory_service_id
	}

	resource "nutanix_project_v2" "test" {
		name = "test"
		project_id = "test"
		description = "test"
	}

	resource "nutanix_users_v2" "user" {
		username = "tf_project_user"
		user_type = "LDAP"
		idp_id = local.idp_ext_id
	}

	resource "nutanix_user_groups_v2" "usergroup" {
		group_type = "LDAP"
		idp_id = local.idp_ext_id
		name = "tf_project_group"
		distinguished_name = "cn=tf_project_group,ou=group,dc=devtest,dc=local"
	}
  
	resource "nutanix_role_membership_v2" "project_admin_role" {
		role_ext_id      = local.project_admin_role_ext_id
		identity_type    = "USER"
		identity_ext_id  = nutanix_users_v2.user.ext_id
		idp_ext_id       = local.idp_ext_id
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
		identity_ext_id  = nutanix_user_groups_v2.usergroup.ext_id
		idp_ext_id       = local.idp_ext_id
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
	}
	
	data "nutanix_role_membership_summary_v2" "get_role_membership_summary" {
		filter = "extId eq '${nutanix_project_v2.test.ext_id}'"
	}
	
`, filepath)
}
