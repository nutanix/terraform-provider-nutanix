package datapoliciesv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceNameProtectionPolicy = "data.nutanix_protection_policy_v2.test"

func TestAccV2NutanixProtectionPolicyDatasource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-protection-policy-%d", r)
	description := "terraform test protection policy CRUD"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testProtectionPolicyV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testProtectionPolicyResourceConfig(name, description) + testProtectionPolicyDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameProtectionPolicy, "ext_id"),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicy, "name", name),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicy, "description", description),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicy, "replication_configurations.0.source_location_label", "source"),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicy, "replication_configurations.0.remote_location_label", "target"),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicy, "replication_configurations.0.schedule.0.recovery_point_objective_time_seconds", "0"),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicy, "replication_configurations.0.schedule.0.recovery_point_type", "CRASH_CONSISTENT"),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicy, "replication_configurations.0.schedule.0.sync_replication_auto_suspend_timeout_seconds", "10"),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicy, "replication_configurations.1.source_location_label", "target"),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicy, "replication_configurations.1.remote_location_label", "source"),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicy, "replication_configurations.1.schedule.0.recovery_point_objective_time_seconds", "0"),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicy, "replication_configurations.1.schedule.0.recovery_point_type", "CRASH_CONSISTENT"),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicy, "replication_configurations.1.schedule.0.sync_replication_auto_suspend_timeout_seconds", "10"),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicy, "replication_locations.0.label", "source"),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicy, "replication_locations.0.is_primary", "true"),
					resource.TestCheckResourceAttrSet(dataSourceNameProtectionPolicy, "replication_locations.0.domain_manager_ext_id"),
					resource.TestCheckResourceAttrSet(dataSourceNameProtectionPolicy, "replication_locations.1.domain_manager_ext_id"),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicy, "replication_locations.1.is_primary", "false"),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicy, "replication_locations.1.label", "target"),
				),
			},
		},
	})
}

func testProtectionPolicyDatasourceConfig() string {
	return `

data "nutanix_protection_policy_v2" "test" {
	ext_id = nutanix_protection_policy_v2.test.id
	depends_on = [nutanix_protection_policy_v2.test]
}

`
}
