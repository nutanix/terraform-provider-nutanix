package datapoliciesv2_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"testing"
)

const resourceNameProtectionPolicy = "nutanix_protection_policy_v2.test"

func TestAccV2NutanixProtectionPolicyResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-protection-policy-%d", r)
	description := "terraform test protection policy CRUD"

	updateName := fmt.Sprintf("tf-test-protection-policy-%d-update", r)
	updateDescription := "terraform test protection policy CRUD update"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testProtectionPolicyResourceConfig(name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameProtectionPolicy, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "name", name),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "description", description),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.source_location_label", "source"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.remote_location_label", "target"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.recovery_point_objective_time_seconds", "0"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.recovery_point_type", "CRASH_CONSISTENT"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.sync_replication_auto_suspend_timeout_seconds", "10"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.source_location_label", "target"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.remote_location_label", "source"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.recovery_point_objective_time_seconds", "0"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.recovery_point_type", "CRASH_CONSISTENT"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.sync_replication_auto_suspend_timeout_seconds", "10"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.0.label", "source"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.0.is_primary", "true"),
					resource.TestCheckResourceAttrSet(resourceNameProtectionPolicy, "replication_locations.0.domain_manager_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameProtectionPolicy, "replication_locations.1.domain_manager_ext_id"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.1.is_primary", "false"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.1.label", "target"),
				),
			},
			// update
			{
				Config: testProtectionPolicyResourceUpdateConfig(updateName, updateDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameProtectionPolicy, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "name", updateName),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "description", updateDescription),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.source_location_label", "source"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.remote_location_label", "target"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.recovery_point_objective_time_seconds", "0"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.recovery_point_type", "APPLICATION_CONSISTENT"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.sync_replication_auto_suspend_timeout_seconds", "20"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.source_location_label", "target"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.remote_location_label", "source"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.recovery_point_objective_time_seconds", "0"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.recovery_point_type", "APPLICATION_CONSISTENT"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.sync_replication_auto_suspend_timeout_seconds", "20"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.0.label", "source"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.0.is_primary", "true"),
					resource.TestCheckResourceAttrSet(resourceNameProtectionPolicy, "replication_locations.0.domain_manager_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameProtectionPolicy, "replication_locations.1.domain_manager_ext_id"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.1.is_primary", "false"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.1.label", "target"),
				),
			},
		},
	})
}

func testProtectionPolicyResourceConfig(name, description string) string {
	return fmt.Sprintf(`
# List domain Managers
data "nutanix_domain_managers_v2" "pcs" {}

# List categories
data "nutanix_categories_v2" "categories" {}

# list Clusters 
data "nutanix_clusters_v2" "clusters" {}

locals {
	clusterExtId = [
		  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
		  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
	][0]
}

resource "nutanix_protection_policy_v2" "test" {
  name        = "%[1]s"
  description = "%[2]s"

  replication_configurations {
    source_location_label = "source"
    remote_location_label = "target"
    schedule {
      recovery_point_objective_time_seconds         = 0
      recovery_point_type                           = "CRASH_CONSISTENT"
      sync_replication_auto_suspend_timeout_seconds = 10
    }
  }
  replication_configurations {
    source_location_label = "target"
    remote_location_label = "source"
    schedule {
      recovery_point_objective_time_seconds         = 0
      recovery_point_type                           = "CRASH_CONSISTENT"
      sync_replication_auto_suspend_timeout_seconds = 10
    }
  }

  replication_locations {
    domain_manager_ext_id = "425cd2d4-32e0-4c2d-a026-31d81fa4c805"
    label                 = "source"
    is_primary            = true
  }
  replication_locations {
    domain_manager_ext_id = data.nutanix_domain_managers_v2.pcs.domain_managers[0].ext_id
    label                 = "target"
    is_primary            = false
  }

  category_ids = [data.nutanix_categories_v2.categories.categories.0.ext_id]
}
`, name, description)
}

func testProtectionPolicyResourceUpdateConfig(name, description string) string {
	return fmt.Sprintf(`
# List domain Managers
data "nutanix_domain_managers_v2" "pcs" {}

# List categories
data "nutanix_categories_v2" "categories" {}

# list Clusters 
data "nutanix_clusters_v2" "clusters" {}

locals {
	clusterExtId = [
		  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
		  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
	][0]
}

resource "nutanix_protection_policy_v2" "test" {
  name        = "%[1]s"
  description = "%[2]s"

  replication_configurations {
    source_location_label = "source"
    remote_location_label = "target"
    schedule {
      recovery_point_objective_time_seconds         = 0
      recovery_point_type                           = "APPLICATION_CONSISTENT"
      sync_replication_auto_suspend_timeout_seconds = 20
    }
  }
  replication_configurations {
    source_location_label = "target"
    remote_location_label = "source"
    schedule {
      recovery_point_objective_time_seconds         = 0
      recovery_point_type                           = "APPLICATION_CONSISTENT"
      sync_replication_auto_suspend_timeout_seconds = 20
    }
  }

  replication_locations {
    domain_manager_ext_id = "425cd2d4-32e0-4c2d-a026-31d81fa4c805"
    label                 = "source"
    is_primary            = true
  }
  replication_locations {
    domain_manager_ext_id = data.nutanix_domain_managers_v2.pcs.domain_managers[0].ext_id
    label                 = "target"
    is_primary            = false
  }

  category_ids = [data.nutanix_categories_v2.categories.categories.0.ext_id,data.nutanix_categories_v2.categories.categories.1.ext_id]
}
`, name, description)
}
