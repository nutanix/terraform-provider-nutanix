package prism_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccNutanixProject_basic(t *testing.T) {
	resourceName := "nutanix_project.project_test"

	subnetName := acctest.RandomWithPrefix("test-subnateName")
	name := acctest.RandomWithPrefix("test-project-name-dou")
	description := acctest.RandomWithPrefix("test-project-desc-dou")
	categoryName := "Environment"
	categoryVal := "Staging"

	updateName := acctest.RandomWithPrefix("test-project-name-dou")
	updateDescription := acctest.RandomWithPrefix("test-project-desc-dou")
	updateCategoryName := "Environment"
	updateCategoryVal := "Production"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixProjectConfig(subnetName, name, description, categoryName, categoryVal),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixProjectExists(&resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "categories.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "api_version", "3.1"),
				),
			},
			{
				Config: testAccNutanixProjectConfig(subnetName, updateName, updateDescription, updateCategoryName, updateCategoryVal),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixProjectExists(&resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", updateDescription),
					resource.TestCheckResourceAttr(resourceName, "categories.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "api_version", "3.1"),
				),
			},
		},
	})
}

func TestAccNutanixProject_importBasic(t *testing.T) {
	resourceName := "nutanix_project.project_test"

	subnetName := acctest.RandomWithPrefix("test-subnateName")
	name := acctest.RandomWithPrefix("test-project-name-dou")
	description := acctest.RandomWithPrefix("test-project-desc-dou")
	categoryName := "Environment"
	categoryVal := "Staging"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixProjectConfig(subnetName, name, description, categoryName, categoryVal),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckNutanixProjectImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNutanixProject_withInternal(t *testing.T) {
	resourceName := "nutanix_project.project_test"

	subnetName := acctest.RandomWithPrefix("test-subnateName")
	name := acctest.RandomWithPrefix("test-project-name-dou")
	description := acctest.RandomWithPrefix("test-project-desc-dou")
	categoryName := "Environment"
	categoryVal := "Staging"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixProjectInternalConfig(subnetName, name, description, categoryName, categoryVal),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "categories.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "api_version", "3.1"),
					resource.TestCheckResourceAttr(resourceName, "subnet_reference_list.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "cluster_reference_list.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpc_reference_list.#", "1"),
				),
			},
		},
	})
}

func TestAccNutanixProject_withInternalUpdate(t *testing.T) {
	resourceName := "nutanix_project.project_test"
	subnetName := acctest.RandomWithPrefix("test-subnateName")
	name := acctest.RandomWithPrefix("test-project-name-dou")
	description := acctest.RandomWithPrefix("test-project-desc-dou")

	updatedName := acctest.RandomWithPrefix("test-project-updated")
	updateDes := acctest.RandomWithPrefix("test-desc-got-updated")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixProjectInternalConfigUpdate(subnetName, name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "api_version", "3.1"),
					resource.TestCheckResourceAttr(resourceName, "subnet_reference_list.#", "1"),
				),
			},
			{
				Config: testAccNutanixProjectInternalConfigUpdate(subnetName, updatedName, updateDes),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", updateDes),
					resource.TestCheckResourceAttr(resourceName, "api_version", "3.1"),
					resource.TestCheckResourceAttr(resourceName, "subnet_reference_list.#", "1"),
				),
			},
		},
	})
}

func TestAccNutanixProject_withInternalWithACP(t *testing.T) {
	resourceName := "nutanix_project.project_test"

	subnetName := acctest.RandomWithPrefix("test-subnateName")
	name := acctest.RandomWithPrefix("test-project-name-dou")
	description := acctest.RandomWithPrefix("test-project-desc-dou")
	categoryName := "Environment"
	categoryVal := "Staging"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixProjectInternalConfigWithACP(subnetName, name, description, categoryName, categoryVal),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "categories.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "api_version", "3.1"),
					resource.TestCheckResourceAttr(resourceName, "subnet_reference_list.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "acp.#", "1"),
				),
			},
		},
	})
}

func TestAccNutanixProject_withInternalWithACPUserGroup(t *testing.T) {
	resourceName := "nutanix_project.project_test"

	subnetName := acctest.RandomWithPrefix("test-subnateName")
	name := acctest.RandomWithPrefix("test-project-name-dou")
	description := acctest.RandomWithPrefix("test-project-desc-dou")
	categoryName := "Environment"
	categoryVal := "Staging"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixProjectInternalConfigWithACPUserGroup(subnetName, name, description, categoryName, categoryVal),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "categories.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "api_version", "3.1"),
					resource.TestCheckResourceAttr(resourceName, "subnet_reference_list.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "acp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "external_user_group_reference_list.#", "1"),
				),
			},
		},
	})
}

func TestAccNutanixProject_ACPOrderAndNestedRefs_NoPlanDiff(t *testing.T) {
	resourceName := "nutanix_project.projects"

	name := acctest.RandomWithPrefix("tf-acc-project-acp-order")
	description := "project description"
	// Move these from constants -> local test variables (per request)
	projectAdminUserName := "ssptest1@qa.nucalm.io"
	developerUserName := "ssptest2@qa.nucalm.io"
	extraUserName := "ssptest3@qa.nucalm.io"
	backupAdminUserName := "ssptest4@qa.nucalm.io"
	roleDeveloperName := "Developer"
	roleProjectAdminName := "Project Admin"
	roleBackupAdminName := "Backup Admin"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixProjectDestroy,
		Steps: []resource.TestStep{
			// 1) Create the project (ACP order: Developer, Project Admin)
			{
				PreConfig: func() {
					fmt.Println("Step 1: creating project")
				},
				Config: testAccNutanixProjectACPOrderConfig(name, description, false, false, false, projectAdminUserName, developerUserName, extraUserName, backupAdminUserName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixProjectExists(&resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "use_project_internal", "true"),
					resource.TestCheckResourceAttr(resourceName, "api_version", "3.1"),
					resource.TestCheckResourceAttr(resourceName, "acp.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "user_reference_list.#", "2"),
					testAccCheckProjectACPHasRole(resourceName, roleDeveloperName),
					testAccCheckProjectACPHasRole(resourceName, roleProjectAdminName),
					testAccCheckProjectACPCountsByRole(resourceName, roleDeveloperName, 1, 0),
					testAccCheckProjectACPCountsByRole(resourceName, roleProjectAdminName, 1, 0),
				),
			},
			// 2) Plan, no changes (same config)
			{
				PreConfig: func() {
					fmt.Println("Step 2: planning project, same config")
				},
				Config:             testAccNutanixProjectACPOrderConfig(name, description, false, false, false, projectAdminUserName, developerUserName, extraUserName, backupAdminUserName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			// 3) Change ACP order in TF config, plan should still be clean (ACP order: Project Admin, Developer)
			{
				PreConfig: func() {
					fmt.Println("Step 3: planning project, change ACP order")
				},
				Config:             testAccNutanixProjectACPOrderConfig(name, description, true, false, false, projectAdminUserName, developerUserName, extraUserName, backupAdminUserName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			// 4) Plan, no changes (repeat)
			{
				PreConfig: func() {
					fmt.Println("Step 4: planning project, same config")
				},
				Config:             testAccNutanixProjectACPOrderConfig(name, description, true, false, false, projectAdminUserName, developerUserName, extraUserName, backupAdminUserName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			// 5) Add a new ACP in the middle between the two ACP blocks (Developer, Backup Admin, Project Admin)
			{
				PreConfig: func() {
					fmt.Println("Step 5: applying config with Backup Admin ACP inserted in the middle")
				},
				Config: testAccNutanixProjectACPOrderConfig(name, description, false, false, true, projectAdminUserName, developerUserName, extraUserName, backupAdminUserName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixProjectExists(&resourceName),
					resource.TestCheckResourceAttr(resourceName, "acp.#", "3"),
					testAccCheckProjectACPHasRole(resourceName, roleDeveloperName),
					testAccCheckProjectACPHasRole(resourceName, roleBackupAdminName),
					testAccCheckProjectACPHasRole(resourceName, roleProjectAdminName),
					testAccCheckProjectACPCountsByRole(resourceName, roleBackupAdminName, 1, 0),
				),
			},
			// 6) Plan, no changes after inserting the middle ACP
			{
				PreConfig: func() {
					fmt.Println("Step 6: planning project, same config (with Backup Admin ACP)")
				},
				Config:             testAccNutanixProjectACPOrderConfig(name, description, false, false, true, projectAdminUserName, developerUserName, extraUserName, backupAdminUserName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			// 7) Add user_reference_list, external_user_group_reference_list,
			//    acp[Project Admin].user_reference_list, acp[Developer].user_group_reference_list
			{
				PreConfig: func() {
					fmt.Println("Step 7: planning project, add user_reference_list, external_user_group_reference_list, acp[Project Admin].user_reference_list, acp[Developer].user_group_reference_list")
				},
				Config: testAccNutanixProjectACPOrderConfig(name, description, true, true, false, projectAdminUserName, developerUserName, extraUserName, backupAdminUserName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixProjectExists(&resourceName),
					resource.TestCheckResourceAttr(resourceName, "acp.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "user_reference_list.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "external_user_group_reference_list.#", "1"),
					testAccCheckProjectACPCountsByRole(resourceName, roleDeveloperName, 1, 1),
					testAccCheckProjectACPCountsByRole(resourceName, roleProjectAdminName, 2, 0),
				),
			},
			// 8) Plan, no changes
			{
				PreConfig: func() {
					fmt.Println("Step 8: planning project, same config")
				},
				Config:             testAccNutanixProjectACPOrderConfig(name, description, true, true, false, projectAdminUserName, developerUserName, extraUserName, backupAdminUserName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			// 9) Change ACP order again, plan should still be clean (ACP order: Developer, Project Admin)
			{
				PreConfig: func() {
					fmt.Println("Step 9: planning project, change ACP order")
				},
				Config:             testAccNutanixProjectACPOrderConfig(name, description, false, true, false, projectAdminUserName, developerUserName, extraUserName, backupAdminUserName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			// 10) Plan, no changes
			{
				PreConfig: func() {
					fmt.Println("Step 10: planning project, same config")
				},
				Config:             testAccNutanixProjectACPOrderConfig(name, description, false, true, false, projectAdminUserName, developerUserName, extraUserName, backupAdminUserName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func testAccCheckNutanixProjectImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return rs.Primary.ID, nil
	}
}

func testAccCheckNutanixProjectExists(resourceName *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[*resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", *resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		return nil
	}
}

func testAccCheckNutanixProjectDestroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_project" {
			continue
		}
		for {
			_, err := conn.API.V3.GetProject(rs.Primary.ID)
			if err != nil {
				if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
					return nil
				}
				return err
			}
			time.Sleep(3000 * time.Millisecond)
		}
	}
	return nil
}

func testAccCheckProjectACPHasRole(resourceName, roleName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		acpCountStr := rs.Primary.Attributes["acp.#"]
		acpCount, err := strconv.Atoi(acpCountStr)
		if err != nil {
			return fmt.Errorf("invalid acp.# %q: %w", acpCountStr, err)
		}

		for i := 0; i < acpCount; i++ {
			if rs.Primary.Attributes[fmt.Sprintf("acp.%d.role_reference.0.name", i)] == roleName {
				return nil
			}
		}

		return fmt.Errorf("acp role not found in state: role name %s", roleName)
	}
}

func testAccCheckProjectACPCountsByRole(resourceName, roleName string, expectedUserRefs, expectedUserGroupRefs int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		acpCountStr := rs.Primary.Attributes["acp.#"]
		acpCount, err := strconv.Atoi(acpCountStr)
		if err != nil {
			return fmt.Errorf("invalid acp.# %q: %w", acpCountStr, err)
		}

		for i := 0; i < acpCount; i++ {
			if rs.Primary.Attributes[fmt.Sprintf("acp.%d.role_reference.0.name", i)] != roleName {
				continue
			}

			userRefsStr := rs.Primary.Attributes[fmt.Sprintf("acp.%d.user_reference_list.#", i)]
			userGroupsStr := rs.Primary.Attributes[fmt.Sprintf("acp.%d.user_group_reference_list.#", i)]

			userRefs, err := strconv.Atoi(userRefsStr)
			if err != nil {
				return fmt.Errorf("invalid acp.%d.user_reference_list.# %q: %w", i, userRefsStr, err)
			}
			userGroups, err := strconv.Atoi(userGroupsStr)
			if err != nil {
				return fmt.Errorf("invalid acp.%d.user_group_reference_list.# %q: %w", i, userGroupsStr, err)
			}

			if userRefs != expectedUserRefs {
				return fmt.Errorf("expected role %s to have %d user_reference_list entries, got %d", roleName, expectedUserRefs, userRefs)
			}
			if userGroups != expectedUserGroupRefs {
				return fmt.Errorf("expected role %s to have %d user_group_reference_list entries, got %d", roleName, expectedUserGroupRefs, userGroups)
			}
			return nil
		}

		return fmt.Errorf("acp role not found in state: role name %s", roleName)
	}
}

// testAccNutanixProjectACPOrderConfig reproduces issue #1042:
// - Create project with ACP blocks, then ensure that reordering ACP blocks does NOT cause drift.
// - Then add extra user/group references and ensure drift-free plans after applying.
func testAccNutanixProjectACPOrderConfig(name, description string, acpProjectAdminFirst, includeExtraRefs, includeBackupAdminACP bool, projectAdminUserName, developerUserName, extraUserName, backupAdminUserName string) string {
	// Base users (always present) - values are derived from data sources in locals below
	usersBlock := `
  # Project Admin User
  user_reference_list {
    name = local.user1_name
    kind = "user"
    uuid = local.user1_uuid
  }

  # Developer User
  user_reference_list {
    name = local.user2_name
    kind = "user"
    uuid = local.user2_uuid
  }
`

	extraRefsBlock := ""
	if includeExtraRefs {
		extraRefsBlock = `
  user_reference_list {
    name = local.user3_name
    kind = "user"
    uuid = local.user3_uuid
  }

  external_user_group_reference_list {
    kind = "user_group"
    name = local.ug1_dn
    uuid = local.ug1_uuid
  }
`
	}

	// ACPs (same semantics, configurable order) - role UUIDs come from role data sources, refs from locals
	backupAdminACP := `
  # Backup Admin ACP
  acp {
    role_reference {
      kind = "role"
      uuid = data.nutanix_role.backup_admin.id
      name = "Backup Admin"
    }
    user_reference_list {
      name = local.user4_name
      kind = "user"
      uuid = local.user4_uuid
    }
  }
`

	projectAdminACP := fmt.Sprintf(`
  # Project Admin ACP
  acp {
    role_reference {
      kind = "role"
      uuid = data.nutanix_role.project_admin.id
      name = "Project Admin"
    }
    user_reference_list {
      name = local.user1_name
      kind = "user"
      uuid = local.user1_uuid
    }
%s
  }
`, func() string {
		if !includeExtraRefs {
			return ""
		}
		return `    user_reference_list {
      name = local.user3_name
      kind = "user"
      uuid = local.user3_uuid
    }
`
	}())

	developerACP := fmt.Sprintf(`
  # Developer ACP
  acp {
    role_reference {
      kind = "role"
      uuid = data.nutanix_role.developer.id
      name = "Developer"
    }
    user_reference_list {
      name = local.user2_name
      kind = "user"
      uuid = local.user2_uuid
    }
%s
  }
`, func() string {
		if !includeExtraRefs {
			return ""
		}
		return `    user_group_reference_list {
      kind = "user_group"
      name = local.ug1_dn
      uuid = local.ug1_uuid
    }
`
	}())

	acpBlock := developerACP + projectAdminACP
	if acpProjectAdminFirst {
		acpBlock = projectAdminACP + developerACP
	}
	if includeBackupAdminACP {
		// Insert Backup Admin ACP in the middle between the two blocks (preserve outer order)
		if acpProjectAdminFirst {
			acpBlock = projectAdminACP + backupAdminACP + developerACP
		} else {
			acpBlock = developerACP + backupAdminACP + projectAdminACP
		}
	}

	// This config mirrors temp/issues/1042/main.tf:
	// - v3 data sources: nutanix_clusters / nutanix_subnets / nutanix_users / nutanix_user_groups / nutanix_role
	// - locals compute cluster UUID, subnet UUID, user UUIDs, and user-group UUID/DN
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {}

locals {
  cluster_ext_id = [
    for cluster in data.nutanix_clusters.clusters.entities :
    cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
  ][0]
}

data "nutanix_role" "developer" {
  role_name = "Developer"
}

data "nutanix_role" "project_admin" {
  role_name = "Project Admin"
}
%s

data "nutanix_subnets" "test" {
  metadata {
    filter = "name==%s"
  }
}

locals {
  subnet_ext_id = data.nutanix_subnets.test.entities[0].metadata.uuid
}

data "nutanix_user_groups" "user_groups" {}

locals {
  ugs = data.nutanix_user_groups.user_groups.entities

  // same pattern as tf script: filter by DN from test_config.json
  ugs_by_dn = [
    for ug in local.ugs :
    ug if try(ug.directory_service_user_group[0].distinguished_name, "") == %q
  ]

  // Prefer the DN match, but fall back to the first user group if the filter returns empty.
  ug1      = try(local.ugs_by_dn[0], local.ugs[0])
  ug1_uuid = try(local.ug1.metadata.uuid, "")
  ug1_dn   = try(local.ug1.directory_service_user_group[0].distinguished_name, "")
}

data "nutanix_users" "users" {}

locals {
  users = data.nutanix_users.users.entities

  // match users by name (as in tf script)
  user1 = try([for u in local.users : u if u.name == "%s"][0], local.users[0])
  user2 = try([for u in local.users : u if u.name == "%s"][0], local.users[1])
%s

  user1_name = local.user1.name
  user2_name = local.user2.name
  user1_uuid = local.user1.metadata.uuid
  user2_uuid = local.user2.metadata.uuid
%s
}

resource "nutanix_project" "projects" {
  name                 = "%s"
  description          = "%s"
  use_project_internal = true
  api_version          = "3.1"
  cluster_uuid         = local.cluster_ext_id

  cluster_reference_list {
    kind = "cluster"
    uuid = local.cluster_ext_id
  }

  account_reference_list {
    kind = "account"
    uuid = "%s"
  }

  subnet_reference_list {
    kind = "subnet"
    uuid = local.subnet_ext_id
  }

  default_subnet_reference {
    kind = "subnet"
    uuid = local.subnet_ext_id
  }
%s
%s
%s
}
`, func() string {
		if !includeBackupAdminACP {
			return ""
		}
		return `
data "nutanix_role" "backup_admin" {
  role_name = "Backup Admin"
}
`
	}(), testVars.SubnetName, testVars.UserGroupWithDistinguishedName[1].DistinguishedName, projectAdminUserName, developerUserName, func() string {
		if !includeBackupAdminACP {
			return ""
		}
		return fmt.Sprintf(`
  user4 = try([for u in local.users : u if u.name == "%s"][0], local.users[0])
  user4_name = local.user4.name
  user4_uuid = local.user4.metadata.uuid
`, backupAdminUserName)
	}(), func() string {
		if !includeExtraRefs {
			return ""
		}
		// third user and its derived locals
		return fmt.Sprintf(`
  user3 = try([for u in local.users : u if u.name == "%s"][0], local.users[2])
  user3_name = local.user3.name
  user3_uuid = local.user3.metadata.uuid
`, extraUserName)
	}(), name, description, testVars.AccountUUID, usersBlock, extraRefsBlock, acpBlock)
}

func testAccNutanixProjectConfig(subnetName, name, description, categoryName, categoryVal string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
			cluster1 = [
				for cluster in data.nutanix_clusters.clusters.entities :
				cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
			][0]
		}

		resource "nutanix_subnet" "subnet" {
			cluster_uuid       = local.cluster1
			name               = "%s"
			description        = "Description of my unit test VLAN"
			vlan_id            = 31
			subnet_type        = "VLAN"
			subnet_ip          = "10.250.140.0"
			default_gateway_ip = "10.250.140.1"
			prefix_length      = 24

			dhcp_options = {
				boot_file_name   = "bootfile"
				domain_name      = "nutanix"
				tftp_server_name = "10.250.140.200"
			}

			dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
			dhcp_domain_search_list      = ["terraform.nutanix.com", "terraform.unit.test.com"]
		}

		resource "nutanix_project" "project_test" {
			name        = "%s"
			description = "%s"

			categories {
				name  = "%s"
				value = "%s"
			}

			default_subnet_reference {
				uuid = nutanix_subnet.subnet.metadata.uuid
			}
			subnet_reference_list{
				uuid = nutanix_subnet.subnet.metadata.uuid
			}

			api_version = "3.1"
		}
	`, subnetName, name, description, categoryName, categoryVal)
}

func testAccNutanixProjectInternalConfig(subnetName, name, description, categoryName, categoryVal string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
			cluster1 = [
				for cluster in data.nutanix_clusters.clusters.entities :
				cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
			][0]
		}

		resource "nutanix_subnet" "subnet" {
			cluster_uuid       = local.cluster1
			name               = "%s"
			description        = "Description of my unit test VLAN"
			vlan_id            = 31
			subnet_type        = "VLAN"
			subnet_ip          = "10.250.140.0"
			default_gateway_ip = "10.250.140.1"
			prefix_length      = 24

			dhcp_options = {
				boot_file_name   = "bootfile"
				domain_name      = "nutanix"
				tftp_server_name = "10.250.140.200"
			}

			dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
			dhcp_domain_search_list      = ["terraform.nutanix.com", "terraform.unit.test.com"]
		}

		resource "nutanix_subnet" "overlay-subnet" {
			cluster_uuid = local.cluster1
			name        = "acctest-subnet-updated"
			description = "Description of my unit test VLAN"
			vlan_id     = 876
			subnet_type = "VLAN"
			subnet_ip          = "10.250.144.0"
		  default_gateway_ip = "10.250.144.1"
		  prefix_length = 24
		  is_external = true
		  ip_config_pool_list_ranges = ["10.250.144.10 10.250.144.20"]
		}

		resource "nutanix_vpc" "acctest-managed" {
			depends_on = [
				resource.nutanix_subnet.overlay-subnet
			]
			name = "acctest-managed-vpc"

			external_subnet_reference_name = [
			  "acctest-subnet-updated"
			]

			common_domain_name_server_ip_list{
					ip = "8.8.8.9"
			}

			externally_routable_prefix_list{
			  ip=  "172.30.0.0"
			  prefix_length= 16
			}
			externally_routable_prefix_list{
				ip=  "172.34.0.0"
				prefix_length= 16
			  }
		  }

		resource "nutanix_project" "project_test" {
			name        = "%s"
			description = "%s"

			categories {
				name  = "%s"
				value = "%s"
			}

			default_subnet_reference {
				uuid = nutanix_subnet.subnet.metadata.uuid
			}

			use_project_internal = true

			api_version = "3.1"

			subnet_reference_list{
				kind="subnet"
				name=nutanix_subnet.subnet.name
				uuid=nutanix_subnet.subnet.metadata.uuid
			}
			subnet_reference_list{
				kind="subnet"
				name=nutanix_subnet.overlay-subnet.name
				uuid=nutanix_subnet.overlay-subnet.id
			}
			cluster_reference_list{
				kind="cluster"
				uuid=local.cluster1
			}
			vpc_reference_list{
				kind="vpc"
				uuid= nutanix_vpc.acctest-managed.id
			}
		}
	`, subnetName, name, description, categoryName, categoryVal)
}

func testAccNutanixProjectInternalConfigUpdate(subnetName, name, description string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
			cluster1 = [
				for cluster in data.nutanix_clusters.clusters.entities :
				cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
			][0]
		}

		resource "nutanix_subnet" "subnet" {
			cluster_uuid       = local.cluster1
			name               = "%s"
			description        = "Description of my unit test VLAN"
			vlan_id            = 31
			subnet_type        = "VLAN"
			subnet_ip          = "10.250.140.0"
			default_gateway_ip = "10.250.140.1"
			prefix_length      = 24

			dhcp_options = {
				boot_file_name   = "bootfile"
				domain_name      = "nutanix"
				tftp_server_name = "10.250.140.200"
			}

			dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
			dhcp_domain_search_list      = ["terraform.nutanix.com", "terraform.unit.test.com"]
		}

		resource "nutanix_subnet" "overlay-subnet" {
			cluster_uuid = local.cluster1
			name        = "acctest-subnet-updated"
			description = "Description of my unit test VLAN"
			vlan_id     = 876
			subnet_type = "VLAN"
			subnet_ip          = "10.250.144.0"
		  default_gateway_ip = "10.250.144.1"
		  prefix_length = 24
		  is_external = true
		  ip_config_pool_list_ranges = ["10.250.144.10 10.250.144.20"]
		}

		resource "nutanix_project" "project_test" {
			name        = "%s"
			description = "%s"

			default_subnet_reference {
				uuid = nutanix_subnet.subnet.metadata.uuid
			}

			use_project_internal = true

			api_version = "3.1"

			subnet_reference_list{
				kind="subnet"
				uuid=nutanix_subnet.subnet.metadata.uuid
			}
		}
	`, subnetName, name, description)
}

func testAccNutanixProjectInternalConfigWithACP(subnetName, name, description, categoryName, categoryVal string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
			cluster1 = [
				for cluster in data.nutanix_clusters.clusters.entities :
				cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
			][0]
		}

		resource "nutanix_subnet" "subnet" {
			cluster_uuid       = local.cluster1
			name               = "%[1]s"
			description        = "Description of my unit test VLAN"
			vlan_id            = 31
			subnet_type        = "VLAN"
			subnet_ip          = "10.250.140.0"
			default_gateway_ip = "10.250.140.1"
			prefix_length      = 24

			dhcp_options = {
				boot_file_name   = "bootfile"
				domain_name      = "nutanix"
				tftp_server_name = "10.250.140.200"
			}

			dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
			dhcp_domain_search_list      = ["terraform.nutanix.com", "terraform.unit.test.com"]
		}

		resource "nutanix_role" "test" {
			name        = "project-role-acctest"
			description = "description role"
			permission_reference_list {
				kind = "permission"
				uuid = "%[6]s"
			}
		}

		resource "nutanix_project" "project_test" {
			name        = "%[2]s"
			description = "%[3]s"
			cluster_uuid = local.cluster1
			categories {
				name  = "%[4]s"
				value = "%[5]s"
			}

			default_subnet_reference {
				uuid = nutanix_subnet.subnet.metadata.uuid
			}

			use_project_internal = true

			api_version = "3.1"

			subnet_reference_list{
				kind="subnet"
				uuid=nutanix_subnet.subnet.metadata.uuid
			}

			user_reference_list{
			uuid = "00000000-0000-0000-0000-000000000000"
			name = "admin"
			}

			acp{
				name="nuCalmAcp-97c623"

				role_reference {
					kind = "role"
					uuid = nutanix_role.test.id
				}

				user_reference_list{
					uuid = "00000000-0000-0000-0000-000000000000"
					name = "admin"
					kind = "user"
				}

				description= "untitledAcp-54acc50f-ab94-640a-5f06-5c855cc09539"
			}
		}
	`, subnetName, name, description, categoryName, categoryVal, testVars.Permissions[0].UUID)
}

func testAccNutanixProjectInternalConfigWithACPUserGroup(subnetName, name, description, categoryName, categoryVal string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
			cluster1 = [
				for cluster in data.nutanix_clusters.clusters.entities :
				cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
			][0]
		}

		resource "nutanix_subnet" "subnet" {
			cluster_uuid       = local.cluster1
			name               = "%[1]s"
			description        = "Description of my unit test VLAN"
			vlan_id            = 31
			subnet_type        = "VLAN"
			subnet_ip          = "10.250.140.0"
			default_gateway_ip = "10.250.140.1"
			prefix_length      = 24

			dhcp_options = {
				boot_file_name   = "bootfile"
				domain_name      = "nutanix"
				tftp_server_name = "10.250.140.200"
			}

			dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
			dhcp_domain_search_list      = ["terraform.nutanix.com", "terraform.unit.test.com"]
		}

		resource "nutanix_role" "test" {
			name        = "project-role-acctest"
			description = "description role"
			permission_reference_list {
				kind = "permission"
				uuid = "%[6]s"
			}
		}

		resource "nutanix_user_groups" "acctest-managed" {
			directory_service_user_group {
				distinguished_name = "%[7]s"
			}
		}

		resource "nutanix_project" "project_test" {
			name        = "%[2]s"
			description = "%[3]s"
			cluster_uuid = local.cluster1
			categories {
				name  = "%[4]s"
				value = "%[5]s"
			}

			default_subnet_reference {
				uuid = nutanix_subnet.subnet.metadata.uuid
			}

			use_project_internal = true

			api_version = "3.1"

			subnet_reference_list{
				kind="subnet"
				uuid=nutanix_subnet.subnet.metadata.uuid
			}

			external_user_group_reference_list {
				name= "%[7]s"
			   	kind= "user_group"
			   	uuid= nutanix_user_groups.acctest-managed.id
			}

			acp{
				name="nuCalmAcp-97c623"

				role_reference {
					kind = "role"
					uuid = nutanix_role.test.id
				}

				user_group_reference_list {
					name= "%[7]s"
					kind= "user_group"
					uuid= nutanix_user_groups.acctest-managed.id
				}

				description= "untitledAcp-54acc50f-ab94-640a-5f06-5c855cc09539"
			}
		}
	`, subnetName, name, description, categoryName, categoryVal, testVars.Permissions[0].UUID, testVars.UserGroupWithDistinguishedName[3].DistinguishedName)
}
