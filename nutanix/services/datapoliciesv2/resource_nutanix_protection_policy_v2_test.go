package datapoliciesv2_test

import (
	"fmt"
	"regexp"
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
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		ProtoV5ProviderFactories: acc.TestAccProtoV5ProviderFactories,
		CheckDestroy:             testProtectionPolicyV2CheckDestroy,
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

func TestAccV2NutanixProtectionPolicyResource_ProjectAssociation(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-pp-projassoc-%d", r)
	description := "protection policy project association test"
	projectName := fmt.Sprintf("tf-pp-pa-proj-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		ProtoV5ProviderFactories: acc.TestAccProtoV5ProviderFactories,
		CheckDestroy:             testProtectionPolicyV2CheckDestroy,
		Steps: []resource.TestStep{
			// Negative: associating a protection policy with a non-default (user) project
			// is not supported by the platform (FEAT-17448, planned for Kronos).
			{
				Config:      testProtectionPolicyProjectAssociationConfig(name, description, projectName, "", true),
				ExpectError: regexp.MustCompile("non default project is not supported"),
			},
			// Positive: the only supported association is the default project, using a
			// category the default project can access (shared_with_projects = []).
			{
				Config: testProtectionPolicyProjectAssociationConfig(name, description, projectName, "00000000-0000-0000-0000-000000000000", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "project_ext_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("data.nutanix_protection_policy_v2.test", "project_ext_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("data.nutanix_protection_policies_v2.test", "protection_policies.#", "1"),
					resource.TestCheckResourceAttrPair("data.nutanix_protection_policies_v2.test", "protection_policies.0.ext_id", resourceNameProtectionPolicy, "ext_id"),
					resource.TestCheckResourceAttr("data.nutanix_protection_policies_v2.test", "protection_policies.0.project_ext_id", "00000000-0000-0000-0000-000000000000"),
				),
			},
			// Negative: changing project_ext_id after creation is rejected by the provider.
			{
				Config:      testProtectionPolicyProjectAssociationConfig(name, description, projectName, "", false),
				ExpectError: regexp.MustCompile("Update of project_ext_id is not supported"),
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
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		ProtoV5ProviderFactories: acc.TestAccProtoV5ProviderFactories,
		CheckDestroy:             testProtectionPolicyV2CheckDestroy,
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
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		ProtoV5ProviderFactories: acc.TestAccProtoV5ProviderFactories,
		CheckDestroy:             testProtectionPolicyV2CheckDestroy,
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

// Case 1: Creating a synchronous protection policy with is_replication_paused = true
// must fail; replication paused can only be set through an update request.
func TestAccV2NutanixProtectionPolicyResource_SyncReplicationPausedCreateError(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-pp-sync-paused-create-%d", r)
	description := "sync pp with is_replication_paused=true on create must fail"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		ProtoV5ProviderFactories: acc.TestAccProtoV5ProviderFactories,
		CheckDestroy:             testProtectionPolicyV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testProtectionPolicySyncScheduleConfig(name, description, "is_replication_paused = true"),
				ExpectError: regexp.MustCompile("replication paused is not supported"),
			},
		},
	})
}

// Case 2: Creating a synchronous protection policy with latest_recovery_point_retention_seconds
// (even 0) must fail; the field is not allowed for synchronous replication.
func TestAccV2NutanixProtectionPolicyResource_SyncLatestRecoveryPointRetentionCreateError(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-pp-sync-lrprs-create-%d", r)
	description := "sync pp with latest_recovery_point_retention_seconds must fail"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		ProtoV5ProviderFactories: acc.TestAccProtoV5ProviderFactories,
		CheckDestroy:             testProtectionPolicyV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testProtectionPolicySyncScheduleConfig(name, description, "latest_recovery_point_retention_seconds = 0"),
				ExpectError: regexp.MustCompile("Latest recovery point retention seconds cannot be specified for a synchronous replication"),
			},
		},
	})
}

// Case 3: Create a synchronous protection policy with is_replication_paused = false, then
// update it to true. Pausing is only supported for synchronous replication via update.
func TestAccV2NutanixProtectionPolicyResource_SyncReplicationPausedUpdate(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-pp-sync-paused-update-%d", r)
	description := "sync pp is_replication_paused create false then update true"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		ProtoV5ProviderFactories: acc.TestAccProtoV5ProviderFactories,
		CheckDestroy:             testProtectionPolicyV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testProtectionPolicySyncScheduleConfig(name, description, "is_replication_paused = false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameProtectionPolicy, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.is_replication_paused", "false"),
				),
			},
			{
				Config: testProtectionPolicySyncScheduleConfig(name, description, "is_replication_paused = true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.is_replication_paused", "true"),
				),
			},
		},
	})
}

// Case 4: is_replication_paused is not supported for asynchronous replication. Create a
// synchronous policy, pause it via update, then attempt to convert it to asynchronous while
// still paused (which must fail), and finally recover by unpausing.
func TestAccV2NutanixProtectionPolicyResource_ReplicationPausedAsyncError(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-pp-paused-async-%d", r)
	description := "is_replication_paused not supported for async replication"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		ProtoV5ProviderFactories: acc.TestAccProtoV5ProviderFactories,
		CheckDestroy:             testProtectionPolicyV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testProtectionPolicySyncScheduleConfig(name, description, "is_replication_paused = false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.is_replication_paused", "false"),
				),
			},
			{
				Config: testProtectionPolicySyncScheduleConfig(name, description, "is_replication_paused = true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.is_replication_paused", "true"),
				),
			},
			{
				Config:      testProtectionPolicyAsyncScheduleConfig(name, description, "is_replication_paused = true", false),
				ExpectError: regexp.MustCompile("Pause replication is not supported for asynchronous or near synchronous replications"),
			},
			{
				Config: testProtectionPolicySyncScheduleConfig(name, description, "is_replication_paused = false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.is_replication_paused", "false"),
				),
			},
		},
	})
}

// Cases 5, 6 and 7: latest_recovery_point_retention_seconds lifecycle on an asynchronous
// policy - create with 0, update to 3600 and back to 0 - and verify the data source reflects
// the value.
func TestAccV2NutanixProtectionPolicyResource_AsyncLatestRecoveryPointRetentionLifecycle(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-pp-async-lrprs-%d", r)
	description := "async pp latest_recovery_point_retention_seconds lifecycle"

	dataSourceName := "data.nutanix_protection_policy_v2.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		ProtoV5ProviderFactories: acc.TestAccProtoV5ProviderFactories,
		CheckDestroy:             testProtectionPolicyV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testProtectionPolicyAsyncScheduleConfig(name, description, "latest_recovery_point_retention_seconds = 0", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameProtectionPolicy, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.latest_recovery_point_retention_seconds", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "replication_configurations.0.schedule.0.latest_recovery_point_retention_seconds", "0"),
				),
			},
			{
				Config: testProtectionPolicyAsyncScheduleConfig(name, description, "latest_recovery_point_retention_seconds = 3600", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.latest_recovery_point_retention_seconds", "3600"),
					resource.TestCheckResourceAttr(dataSourceName, "replication_configurations.0.schedule.0.latest_recovery_point_retention_seconds", "3600"),
				),
			},
			{
				Config: testProtectionPolicyAsyncScheduleConfig(name, description, "latest_recovery_point_retention_seconds = 0", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameProtectionPolicy, "replication_configurations.0.schedule.0.latest_recovery_point_retention_seconds", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "replication_configurations.0.schedule.0.latest_recovery_point_retention_seconds", "0"),
				),
			},
		},
	})
}

func ppProjectExtIDLine(override string) string {
	if override == "" {
		return `project_ext_id = nutanix_project_v2.test.ext_id`
	}
	return fmt.Sprintf(`project_ext_id = "%s"`, override)
}

func testProtectionPolicyProjectAssociationConfig(name, description, projectName, projectExtIDOverride string, shareCategories bool) string {
	shareBlock := `shared_with_projects = []`
	if shareCategories {
		shareBlock = `shared_with_projects = [nutanix_project_v2.test.ext_id]`
	}
	return fmt.Sprintf(`
data "nutanix_pcs_v2" "pcs-list" {}

locals {
  config            = jsondecode(file("%[3]s"))
  availability_zone = local.config.availability_zone
}

resource "nutanix_project_v2" "test" {
  name        = "%[4]s"
  project_id  = "%[4]s"
  description = "project association test"
}

resource "nutanix_category_v2" "test" {
  key         = "tf-pp-pa-cat-%[4]s"
  value       = "pp_pa_category_value"
  description = "category for protection policy project association"
  %[6]s
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
      start_time                                    = "23h:54m"
    }
  }
  replication_configurations {
    source_location_label = "target"
    remote_location_label = "source"
    schedule {
      recovery_point_objective_time_seconds         = 0
      recovery_point_type                           = "CRASH_CONSISTENT"
      sync_replication_auto_suspend_timeout_seconds = 10
      start_time                                    = "23h:54m"
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
  %[5]s
  depends_on = [nutanix_project_v2.test]
}

data "nutanix_protection_policy_v2" "test" {
  ext_id     = nutanix_protection_policy_v2.test.ext_id
  depends_on = [nutanix_protection_policy_v2.test]
}

data "nutanix_protection_policies_v2" "test" {
  filter     = "name eq '${nutanix_protection_policy_v2.test.name}'"
  depends_on = [nutanix_protection_policy_v2.test]
}
`, name, description, filepath, projectName, ppProjectExtIDLine(projectExtIDOverride), shareBlock)
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

// testProtectionPolicySyncScheduleConfig builds a two-way synchronous protection policy where
// scheduleExtra is injected into each schedule block (e.g. is_replication_paused or
// latest_recovery_point_retention_seconds), used by the paused/retention scenarios.
func testProtectionPolicySyncScheduleConfig(name, description, scheduleExtra string) string {
	return fmt.Sprintf(`
data "nutanix_pcs_v2" "pcs-list" {}

locals {
  config            = jsondecode(file("%[3]s"))
  availability_zone = local.config.availability_zone
}

resource "nutanix_category_v2" "test" {
  key         = "tf-test-cat-pp-sched"
  value       = "category_pp_sched"
  description = "category for protection policy schedule tests"
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
      start_time                                    = "23h:54m"
      %[4]s
    }
  }
  replication_configurations {
    source_location_label = "target"
    remote_location_label = "source"
    schedule {
      recovery_point_objective_time_seconds         = 0
      recovery_point_type                           = "CRASH_CONSISTENT"
      sync_replication_auto_suspend_timeout_seconds = 10
      start_time                                    = "23h:54m"
      %[4]s
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
`, name, description, filepath, scheduleExtra)
}

// testProtectionPolicyAsyncScheduleConfig builds a two-way asynchronous protection policy
// (RPO = 3600 with auto rollup retention) where scheduleExtra is injected into each schedule
// block. When withDataSource is true, a single protection policy data source is included.
func testProtectionPolicyAsyncScheduleConfig(name, description, scheduleExtra string, withDataSource bool) string {
	dataSource := ""
	if withDataSource {
		dataSource = `
data "nutanix_protection_policy_v2" "test" {
  ext_id     = nutanix_protection_policy_v2.test.ext_id
  depends_on = [nutanix_protection_policy_v2.test]
}
`
	}
	return fmt.Sprintf(`
data "nutanix_pcs_v2" "pcs-list" {}

locals {
  config            = jsondecode(file("%[3]s"))
  availability_zone = local.config.availability_zone
}

resource "nutanix_category_v2" "test" {
  key         = "tf-test-cat-pp-sched"
  value       = "category_pp_sched"
  description = "category for protection policy schedule tests"
}

resource "nutanix_protection_policy_v2" "test" {
  name        = "%[1]s"
  description = "%[2]s"

  replication_configurations {
    source_location_label = "source"
    remote_location_label = "target"
    schedule {
      recovery_point_objective_time_seconds = 3600
      recovery_point_type                   = "CRASH_CONSISTENT"
      start_time                            = "23h:54m"
      retention {
        auto_rollup_retention {
          local {
            snapshot_interval_type = "DAILY"
            frequency              = 1
          }
          remote {
            snapshot_interval_type = "DAILY"
            frequency              = 1
          }
        }
      }
      %[4]s
    }
  }
  replication_configurations {
    source_location_label = "target"
    remote_location_label = "source"
    schedule {
      recovery_point_objective_time_seconds = 3600
      recovery_point_type                   = "CRASH_CONSISTENT"
      start_time                            = "23h:54m"
      retention {
        auto_rollup_retention {
          local {
            snapshot_interval_type = "DAILY"
            frequency              = 1
          }
          remote {
            snapshot_interval_type = "DAILY"
            frequency              = 1
          }
        }
      }
      %[4]s
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
%[5]s
`, name, description, filepath, scheduleExtra, dataSource)
}
