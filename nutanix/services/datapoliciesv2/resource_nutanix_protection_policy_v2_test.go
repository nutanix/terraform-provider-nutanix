package datapoliciesv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameProtectionPolicy = "nutanix_protection_policy_v2.test"

func TestAccV2NutanixProtectionPolicyResource_Synchronous(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-protection-policy-%d", r)
	description := "terraform test protection policy CRUD"

	updateName := fmt.Sprintf("tf-test-protection-policy-%d-update", r)
	updateDescription := "terraform test protection policy CRUD update"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testProtectionPolicyV2CheckDestroy,
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
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.start_time", "23h:54m"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.source_location_label", "target"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.remote_location_label", "source"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.recovery_point_objective_time_seconds", "0"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.recovery_point_type", "CRASH_CONSISTENT"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.sync_replication_auto_suspend_timeout_seconds", "10"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.start_time", "23h:54m"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.0.label", "source"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.0.is_primary", "true"),
					resource.TestCheckResourceAttrSet(resourceNameProtectionPolicy, "replication_locations.0.domain_manager_ext_id"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.1.domain_manager_ext_id", testVars.AvailabilityZone.PcExtID),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.1.is_primary", "false"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.1.label", "target"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "category_ids.#", "1"),
				),
			},
			// update
			{
				Config: testProtectionPolicyResourceUpdateConfig(updateName, updateDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameProtectionPolicy, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "name", updateName),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "description", updateDescription),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.source_location_label", "source-updated"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.remote_location_label", "target-updated"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.recovery_point_objective_time_seconds", "60"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.recovery_point_type", "APPLICATION_CONSISTENT"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.sync_replication_auto_suspend_timeout_seconds", "20"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.retention.0.auto_rollup_retention.0.local.0.frequency", "2"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.retention.0.auto_rollup_retention.0.local.0.snapshot_interval_type", "WEEKLY"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.retention.0.auto_rollup_retention.0.remote.0.frequency", "1"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.retention.0.auto_rollup_retention.0.remote.0.snapshot_interval_type", "DAILY"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.start_time", "15h:19m"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.source_location_label", "target-updated"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.remote_location_label", "source-updated"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.recovery_point_objective_time_seconds", "60"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.recovery_point_type", "APPLICATION_CONSISTENT"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.sync_replication_auto_suspend_timeout_seconds", "30"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.retention.0.auto_rollup_retention.0.local.0.frequency", "1"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.retention.0.auto_rollup_retention.0.local.0.snapshot_interval_type", "DAILY"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.retention.0.auto_rollup_retention.0.remote.0.frequency", "2"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.retention.0.auto_rollup_retention.0.remote.0.snapshot_interval_type", "WEEKLY"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.start_time", "15h:19m"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.0.label", "source-updated"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.0.is_primary", "true"),
					resource.TestCheckResourceAttrSet(resourceNameProtectionPolicy, "replication_locations.0.domain_manager_ext_id"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.1.domain_manager_ext_id", testVars.AvailabilityZone.PcExtID),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.1.is_primary", "false"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.1.label", "target-updated"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "category_ids.#", "1"),
				),
			},
		},
	})
}

func TestAccV2NutanixProtectionPolicyResource_LinearRetention(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-protection-policy-%d", r)
	description := "terraform test protection policy CRUD"

	nameUpdated := fmt.Sprintf("tf-test-protection-policy-%d-update", r)
	descriptionUpdated := "terraform test protection policy CRUD update"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testProtectionPolicyV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testProtectionPolicyResourceConfigLinearRetentionConfig(name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameProtectionPolicy, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "name", name),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "description", description),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.source_location_label", "0"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.remote_location_label", "1"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.recovery_point_objective_time_seconds", "7200"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.recovery_point_type", "CRASH_CONSISTENT"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.start_time", "23h:54m"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.retention.0.linear_retention.0.local", "1"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.retention.0.linear_retention.0.remote", "1"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.source_location_label", "1"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.remote_location_label", "0"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.recovery_point_objective_time_seconds", "7200"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.recovery_point_type", "CRASH_CONSISTENT"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.retention.0.linear_retention.0.local", "1"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.retention.0.linear_retention.0.remote", "1"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.0.label", "0"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.0.is_primary", "true"),
					resource.TestCheckResourceAttrSet(resourceNameProtectionPolicy, "replication_locations.0.domain_manager_ext_id"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.1.label", "1"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.1.is_primary", "false"),
					resource.TestCheckResourceAttrSet(resourceNameProtectionPolicy, "replication_locations.1.domain_manager_ext_id"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "category_ids.#", "1"),
				),
			},
			{
				Config: testProtectionPolicyResourceConfigLinearRetentionUpdateConfig(nameUpdated, descriptionUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameProtectionPolicy, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "description", descriptionUpdated),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.source_location_label", "0-updated"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.remote_location_label", "1-updated"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.recovery_point_objective_time_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.recovery_point_type", "APPLICATION_CONSISTENT"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.start_time", "15h:19m"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.retention.0.linear_retention.0.local", "2"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.retention.0.linear_retention.0.remote", "2"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.source_location_label", "1-updated"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.remote_location_label", "0-updated"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.recovery_point_objective_time_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.recovery_point_type", "APPLICATION_CONSISTENT"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.start_time", "15h:19m"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.retention.0.linear_retention.0.local", "2"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.retention.0.linear_retention.0.remote", "2"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.0.label", "0-updated"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.0.is_primary", "true"),
					resource.TestCheckResourceAttrSet(resourceNameProtectionPolicy, "replication_locations.0.domain_manager_ext_id"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.1.label", "1-updated"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.1.is_primary", "false"),
					resource.TestCheckResourceAttrSet(resourceNameProtectionPolicy, "replication_locations.1.domain_manager_ext_id"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "category_ids.#", "1"),
				),
			},
		},
	})
}

func TestAccV2NutanixProtectionPolicyResource_AutoRollupRetention(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-protection-policy-%d", r)
	description := "terraform test protection policy CRUD auto rollup retention"

	nameUpdated := fmt.Sprintf("tf-test-protection-policy-%d-update", r)
	descriptionUpdated := "terraform test protection policy CRUD update auto rollup retention"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testProtectionPolicyV2CheckDestroy,
		Steps: []resource.TestStep{

			{
				Config: testProtectionPolicyResourceConfigAutoRollupRetentionConfig(name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameProtectionPolicy, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "name", name),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "description", description),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.source_location_label", "source"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.remote_location_label", "target"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.recovery_point_objective_time_seconds", "60"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.recovery_point_type", "CRASH_CONSISTENT"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.sync_replication_auto_suspend_timeout_seconds", "20"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.retention.0.auto_rollup_retention.0.local.0.frequency", "2"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.retention.0.auto_rollup_retention.0.local.0.snapshot_interval_type", "WEEKLY"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.retention.0.auto_rollup_retention.0.remote.0.frequency", "1"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.retention.0.auto_rollup_retention.0.remote.0.snapshot_interval_type", "DAILY"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.start_time", "18h:10m"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.source_location_label", "target"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.remote_location_label", "source"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.recovery_point_objective_time_seconds", "60"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.recovery_point_type", "CRASH_CONSISTENT"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.sync_replication_auto_suspend_timeout_seconds", "30"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.retention.0.auto_rollup_retention.0.local.0.frequency", "1"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.retention.0.auto_rollup_retention.0.local.0.snapshot_interval_type", "DAILY"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.retention.0.auto_rollup_retention.0.remote.0.frequency", "2"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.retention.0.auto_rollup_retention.0.remote.0.snapshot_interval_type", "WEEKLY"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.start_time", "18h:10m"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.0.label", "source"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.0.is_primary", "true"),
					resource.TestCheckResourceAttrSet(resourceNameProtectionPolicy, "replication_locations.0.domain_manager_ext_id"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.1.domain_manager_ext_id", testVars.AvailabilityZone.PcExtID),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.1.is_primary", "false"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.1.label", "target"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "category_ids.#", "1"),
				),
			},
			{
				Config: testProtectionPolicyResourceConfigAutoRollupRetentionUpdateConfig(nameUpdated, descriptionUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameProtectionPolicy, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "description", descriptionUpdated),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.source_location_label", "source-updated"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.remote_location_label", "target-updated"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.recovery_point_objective_time_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.recovery_point_type", "APPLICATION_CONSISTENT"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.sync_replication_auto_suspend_timeout_seconds", "90"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.retention.0.auto_rollup_retention.0.local.0.snapshot_interval_type", "WEEKLY"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.retention.0.auto_rollup_retention.0.local.0.frequency", "3"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.retention.0.auto_rollup_retention.0.remote.0.snapshot_interval_type", "DAILY"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.retention.0.auto_rollup_retention.0.remote.0.frequency", "2"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.start_time", "13h:08m"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.source_location_label", "target-updated"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.remote_location_label", "source-updated"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.recovery_point_objective_time_seconds", "3600"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.recovery_point_type", "APPLICATION_CONSISTENT"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.sync_replication_auto_suspend_timeout_seconds", "120"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.retention.0.auto_rollup_retention.0.local.0.frequency", "2"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.retention.0.auto_rollup_retention.0.local.0.snapshot_interval_type", "DAILY"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.retention.0.auto_rollup_retention.0.remote.0.frequency", "3"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.retention.0.auto_rollup_retention.0.remote.0.snapshot_interval_type", "WEEKLY"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.1.schedule.0.start_time", "13h:08m"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.0.label", "source-updated"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.0.is_primary", "true"),
					resource.TestCheckResourceAttrSet(resourceNameProtectionPolicy, "replication_locations.0.domain_manager_ext_id"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.1.domain_manager_ext_id", testVars.AvailabilityZone.PcExtID),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.1.is_primary", "false"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_locations.1.label", "target-updated"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "category_ids.#", "1"),
				),
			},
		},
	})
}

func testProtectionPolicyResourceConfig(name, description string) string {
	return fmt.Sprintf(`
# List domain Managers
data "nutanix_pcs_v2" "pcs-list" {}

locals {
	config = jsondecode(file("%[3]s"))
  	availability_zone = local.config.availability_zone
}

# Create Category
resource "nutanix_category_v2" "test" {
  key = "tf-test-category-synchronous-protection-policy"
  value = "category_synchronous_protection_policy"
  description = "category for synchronous protection policy "
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
      start_time									= "23h:54m"
    }
  }
  replication_configurations {
    source_location_label = "target"
    remote_location_label = "source"
    schedule {
      recovery_point_objective_time_seconds         = 0
      recovery_point_type                           = "CRASH_CONSISTENT"
      sync_replication_auto_suspend_timeout_seconds = 10
      start_time									= "23h:54m"
    }
  }

  replication_locations {
    domain_manager_ext_id = data.nutanix_pcs_v2.pcs-list.pcs[0].ext_id
    label                 = "source"
    is_primary            = true
  }
  replication_locations {
    domain_manager_ext_id = local.availability_zone.pc_ext_id
    label                 = "target"
    is_primary            = false
  }

  category_ids = [nutanix_category_v2.test.id]
}
`, name, description, filepath)
}

func testProtectionPolicyResourceUpdateConfig(name, description string) string {
	return fmt.Sprintf(`

# List domain Managers
data "nutanix_pcs_v2" "pcs-list" {}

locals {
	config = jsondecode(file("%[3]s"))
  	availability_zone = local.config.availability_zone
}

# Create Category
resource "nutanix_category_v2" "test" {
  key = "tf-test-category-synchronous-protection-policy"
  value = "category_synchronous_protection_policy"
  description = "category for synchronous protection policy "
}

resource "nutanix_protection_policy_v2" "test" {
  name        = "%[1]s"
  description = "%[2]s"

  replication_configurations {
    source_location_label = "source-updated"
    remote_location_label = "target-updated"
    schedule {
      recovery_point_objective_time_seconds         = 60
      recovery_point_type                           = "APPLICATION_CONSISTENT"
      sync_replication_auto_suspend_timeout_seconds = 20
      start_time									= "15h:19m"
      retention {
        auto_rollup_retention {
          local {
            snapshot_interval_type = "WEEKLY"
            frequency              = 2
          }
          remote {
            snapshot_interval_type = "DAILY"
            frequency              = 1
          }
        }
      }
    }
  }
  replication_configurations {
    source_location_label = "target-updated"
    remote_location_label = "source-updated"
    schedule {
      recovery_point_objective_time_seconds         = 60
      recovery_point_type                           = "APPLICATION_CONSISTENT"
      sync_replication_auto_suspend_timeout_seconds = 30
      start_time									= "15h:19m"
      retention {
        auto_rollup_retention {
          local {
            snapshot_interval_type = "DAILY"
            frequency              = 1
          }
          remote {
            snapshot_interval_type = "WEEKLY"
            frequency              = 2
          }
        }
      }
    }
  }

  replication_locations {
    domain_manager_ext_id = data.nutanix_pcs_v2.pcs-list.pcs[0].ext_id
    label                 = "source-updated"
    is_primary            = true
  }
  replication_locations {
    domain_manager_ext_id = local.availability_zone.pc_ext_id
    label                 = "target-updated"
    is_primary            = false
  }

  category_ids = [ nutanix_category_v2.test.id ]
}
`, name, description, filepath)
}

func testProtectionPolicyResourceConfigLinearRetentionConfig(name, description string) string {
	return fmt.Sprintf(`
# List domain Managers
data "nutanix_pcs_v2" "pcs-list" {}

# list Clusters
data "nutanix_clusters_v2" "clusters" {}

locals {
	clusterExtId = [
		  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
		  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
	][0]
	config = jsondecode(file("%[3]s"))
  	availability_zone = local.config.availability_zone
}

# Create Category
resource "nutanix_category_v2" "test" {
  key = "tf-test-category-linear-retention-protection-policy"
  value = "category_linear_retention_protection_policy"
  description = "category for linea retention protection policy"
}

resource "nutanix_protection_policy_v2" "test" {
  name        = "%[1]s"
  description = "%[2]s"

 replication_configurations {
    source_location_label = "0"
    remote_location_label = "1"
    schedule {
      recovery_point_objective_time_seconds = 7200
      recovery_point_type                   = "CRASH_CONSISTENT"
	  start_time							= "23h:54m"
      retention {
        linear_retention {
          local  = 1
          remote = 1
        }
      }
    }
  }
  replication_configurations {
    source_location_label = "1"
    remote_location_label = "0"
    schedule {
      recovery_point_objective_time_seconds = 7200
      recovery_point_type                   = "CRASH_CONSISTENT"
	  start_time							= "23h:54m"
      retention {
        linear_retention {
          local  = 1
          remote = 1
        }
      }
    }
  }

  replication_locations {
    domain_manager_ext_id = data.nutanix_pcs_v2.pcs-list.pcs[0].ext_id
    label                 = "0"
    is_primary            = true
    replication_sub_location {
      cluster_ext_ids {
        cluster_ext_ids = [local.clusterExtId]
      }
    }
  }
  replication_locations {
    domain_manager_ext_id = local.availability_zone.pc_ext_id
    label                 = "1"
    is_primary            = false
  }

  category_ids = [ nutanix_category_v2.test.id ]
}`, name, description, filepath)
}

func testProtectionPolicyResourceConfigLinearRetentionUpdateConfig(name, description string) string {
	return fmt.Sprintf(`
# List domain Managers
data "nutanix_pcs_v2" "pcs-list" {}

# list Clusters
data "nutanix_clusters_v2" "clusters" {}

locals {
	clusterExtId = [
		  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
		  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
	][0]
	config = jsondecode(file("%[3]s"))
  	availability_zone = local.config.availability_zone
}

# Create Category
resource "nutanix_category_v2" "test" {
  key = "tf-test-category-linear-retention-protection-policy"
  value = "category_linear_retention_protection_policy"
  description = "category for linea retention protection policy"
}

resource "nutanix_protection_policy_v2" "test" {
  name        = "%[1]s"
  description = "%[2]s"

 replication_configurations {
    source_location_label = "0-updated"
    remote_location_label = "1-updated"
    schedule {
      recovery_point_objective_time_seconds = 3600
      recovery_point_type                   = "APPLICATION_CONSISTENT"
	  start_time							= "15h:19m"
      retention {
        linear_retention {
          local  = 2
          remote = 2
        }
      }
    }
  }
  replication_configurations {
    source_location_label = "1-updated"
    remote_location_label = "0-updated"
    schedule {
      recovery_point_objective_time_seconds = 3600
      recovery_point_type                   = "APPLICATION_CONSISTENT"
	  start_time							= "15h:19m"
      retention {
        linear_retention {
          local  = 2
          remote = 2
        }
      }
    }
  }

  replication_locations {
    domain_manager_ext_id = data.nutanix_pcs_v2.pcs-list.pcs[0].ext_id
    label                 = "0-updated"
    is_primary            = true
    replication_sub_location {
      cluster_ext_ids {
        cluster_ext_ids = [local.clusterExtId]
      }
    }
  }
  replication_locations {
    domain_manager_ext_id = local.availability_zone.pc_ext_id
    label                 = "1-updated"
    is_primary            = false
  }

  category_ids = [ nutanix_category_v2.test.id ]
}`, name, description, filepath)
}

func testProtectionPolicyResourceConfigAutoRollupRetentionConfig(name, description string) string {
	return fmt.Sprintf(`
# List domain Managers
data "nutanix_pcs_v2" "pcs-list" {}

locals {
	config = jsondecode(file("%[3]s"))
  	availability_zone = local.config.availability_zone
}

# Create Category
resource "nutanix_category_v2" "test" {
  key = "tf-test-category-auto-rollup-retention-protection-policy"
  value = "category_auto_rollup_retention_protection_policy"
  description = "category for auto rollup retention protection policy "
}

resource "nutanix_protection_policy_v2" "test" {
  name        = "%[1]s"
  description = "%[2]s"

  replication_configurations {
    source_location_label = "source"
    remote_location_label = "target"
    schedule {
      recovery_point_objective_time_seconds         = 60
      recovery_point_type                           = "CRASH_CONSISTENT"
      sync_replication_auto_suspend_timeout_seconds = 20
      start_time = "18h:10m"
      retention {
        auto_rollup_retention {
          local {
            snapshot_interval_type = "WEEKLY"
            frequency              = 2
          }
          remote {
            snapshot_interval_type = "DAILY"
            frequency              = 1
          }
        }
      }
    }
  }
  replication_configurations {
    source_location_label = "target"
    remote_location_label = "source"
    schedule {
      recovery_point_objective_time_seconds         = 60
      recovery_point_type                           = "CRASH_CONSISTENT"
      sync_replication_auto_suspend_timeout_seconds = 30
      start_time = "18h:10m"
      retention {
        auto_rollup_retention {
          local {
            snapshot_interval_type = "DAILY"
            frequency              = 1
          }
          remote {
            snapshot_interval_type = "WEEKLY"
            frequency              = 2
          }
        }
      }
    }
  }

  replication_locations {
    domain_manager_ext_id = data.nutanix_pcs_v2.pcs-list.pcs[0].ext_id
    label                 = "source"
    is_primary            = true
  }
  replication_locations {
    domain_manager_ext_id = local.availability_zone.pc_ext_id
    label                 = "target"
    is_primary            = false
  }

  category_ids = [ nutanix_category_v2.test.id ]
}
`, name, description, filepath)
}

func testProtectionPolicyResourceConfigAutoRollupRetentionUpdateConfig(name, description string) string {
	return fmt.Sprintf(`
# List domain Managers
data "nutanix_pcs_v2" "pcs-list" {}

locals {
	config = jsondecode(file("%[3]s"))
  	availability_zone = local.config.availability_zone
}

# Create Category
resource "nutanix_category_v2" "test" {
  key = "tf-test-category-auto-rollup-retention-protection-policy"
  value = "category_auto_rollup_retention_protection_policy"
  description = "category for auto rollup retention protection policy "
}

resource "nutanix_protection_policy_v2" "test" {
  name        = "%[1]s"
  description = "%[2]s"

  replication_configurations {
    source_location_label = "source-updated"
    remote_location_label = "target-updated"
    schedule {
      recovery_point_objective_time_seconds         = 3600
      recovery_point_type                           = "APPLICATION_CONSISTENT"
      sync_replication_auto_suspend_timeout_seconds = 90
      start_time = "13h:08m"
      retention {
        auto_rollup_retention {
          local {
            snapshot_interval_type = "WEEKLY"
            frequency              = 3
          }
          remote {
            snapshot_interval_type = "DAILY"
            frequency              = 2
          }
        }
      }
    }
  }
  replication_configurations {
    source_location_label = "target-updated"
    remote_location_label = "source-updated"
    schedule {
      recovery_point_objective_time_seconds         = 3600
      recovery_point_type                           = "APPLICATION_CONSISTENT"
      sync_replication_auto_suspend_timeout_seconds = 120
      start_time = "13h:08m"
      retention {
        auto_rollup_retention {
          local {
            snapshot_interval_type = "DAILY"
            frequency              = 2
          }
          remote {
            snapshot_interval_type = "WEEKLY"
            frequency              = 3
          }
        }
      }
    }
  }

  replication_locations {
    domain_manager_ext_id = data.nutanix_pcs_v2.pcs-list.pcs[0].ext_id
    label                 = "source-updated"
    is_primary            = true
  }
  replication_locations {
    domain_manager_ext_id = local.availability_zone.pc_ext_id
    label                 = "target-updated"
    is_primary            = false
  }

  category_ids = [ nutanix_category_v2.test.id ]
}
`, name, description, filepath)
}
