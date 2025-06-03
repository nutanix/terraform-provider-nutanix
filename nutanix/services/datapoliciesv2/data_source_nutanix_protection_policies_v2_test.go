package datapoliciesv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceNameProtectionPolicies = "data.nutanix_protection_policies_v2.test"

func TestAccV2NutanixProtectionPoliciesDatasource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-protection-policy-%d", r)
	description := "terraform test protection policy CRUD"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testProtectionPolicyV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testProtectionPolicyResourceConfig(name, description) + testProtectionPoliciesDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameProtectionPolicies, "protection_policies.#"),
					checkAttributeLength(dataSourceNameProtectionPolicies, "protection_policies", 1),
				),
			},
		},
	})
}

func TestAccV2NutanixProtectionPoliciesDatasource_WithFilter(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-protection-policy-%d", r)
	description := "terraform test protection policy CRUD"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testProtectionPolicyV2CheckDestroy,
		Steps: []resource.TestStep{

			{
				Config: testProtectionPolicyResourceConfig(name, description) + testProtectionPoliciesDatasourceConfigWithFilter(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameProtectionPolicies, "protection_policies.#"),
					checkAttributeLengthEqual(dataSourceNameProtectionPolicies, "protection_policies", 1),

					resource.TestCheckResourceAttrSet(dataSourceNameProtectionPolicies, "protection_policies.0.ext_id"),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicies, "protection_policies.0.name", name),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicies, "protection_policies.0.description", description),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicies, "protection_policies.0.replication_configurations.0.source_location_label", "source"),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicies, "protection_policies.0.replication_configurations.0.remote_location_label", "target"),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicies, "protection_policies.0.replication_configurations.0.schedule.0.recovery_point_objective_time_seconds", "0"),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicies, "protection_policies.0.replication_configurations.0.schedule.0.recovery_point_type", "CRASH_CONSISTENT"),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicies, "protection_policies.0.replication_configurations.0.schedule.0.sync_replication_auto_suspend_timeout_seconds", "10"),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicies, "protection_policies.0.replication_configurations.1.source_location_label", "target"),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicies, "protection_policies.0.replication_configurations.1.remote_location_label", "source"),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicies, "protection_policies.0.replication_configurations.1.schedule.0.recovery_point_objective_time_seconds", "0"),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicies, "protection_policies.0.replication_configurations.1.schedule.0.recovery_point_type", "CRASH_CONSISTENT"),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicies, "protection_policies.0.replication_configurations.1.schedule.0.sync_replication_auto_suspend_timeout_seconds", "10"),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicies, "protection_policies.0.replication_locations.0.label", "source"),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicies, "protection_policies.0.replication_locations.0.is_primary", "true"),
					resource.TestCheckResourceAttrSet(dataSourceNameProtectionPolicies, "protection_policies.0.replication_locations.0.domain_manager_ext_id"),
					resource.TestCheckResourceAttrSet(dataSourceNameProtectionPolicies, "protection_policies.0.replication_locations.1.domain_manager_ext_id"),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicies, "protection_policies.0.replication_locations.1.is_primary", "false"),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicies, "protection_policies.0.replication_locations.1.label", "target"),
				),
			},
		},
	})
}

func TestAccV2NutanixProtectionPoliciesDatasource_WithLimit(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-protection-policy-%d", r)
	description := "terraform test protection policy CRUD"

	limit := 1
	page := 0

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testProtectionPolicyV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testProtectionPolicyResourceConfig(name, description) + testProtectionPoliciesDatasourceConfigWithLimit(limit, page),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameProtectionPolicies, "protection_policies.#"),
					checkAttributeLengthEqual(dataSourceNameProtectionPolicies, "protection_policies", limit),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicies, "limit", fmt.Sprintf("%d", limit)),
					resource.TestCheckResourceAttr(dataSourceNameProtectionPolicies, "page", fmt.Sprintf("%d", page)),
				),
			},
		},
	})
}

func TestAccV2NutanixProtectionPoliciesDatasource_WithInvalidFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testProtectionPolicyV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testProtectionPoliciesDatasourceConfigWithInvalidFilter(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameProtectionPolicies, "protection_policies.#"),
					checkAttributeLengthEqual(dataSourceNameProtectionPolicies, "protection_policies", 0),
				),
			},
		},
	})
}

func testProtectionPoliciesDatasourceConfig() string {
	return `

data "nutanix_protection_policies_v2" "test" {
	depends_on = [nutanix_protection_policy_v2.test]
}

`
}

func testProtectionPoliciesDatasourceConfigWithFilter(name string) string {
	return fmt.Sprintf(`

data "nutanix_protection_policies_v2" "test" {
	filter = "name eq '%s'"
	depends_on = [nutanix_protection_policy_v2.test]
}

`, name)
}

func testProtectionPoliciesDatasourceConfigWithLimit(limit, page int) string {
	return fmt.Sprintf(`
data "nutanix_protection_policies_v2" "test" {
	limit = %d
	page = %d
	depends_on = [nutanix_protection_policy_v2.test]
}

`, limit, page)
}

func testProtectionPoliciesDatasourceConfigWithInvalidFilter() string {
	return `

data "nutanix_protection_policies_v2" "test" {
	filter = "name eq 'invalid_filter'"
}

`
}
