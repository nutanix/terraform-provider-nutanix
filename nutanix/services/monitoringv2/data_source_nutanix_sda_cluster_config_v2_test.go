package monitoringv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameClusterConfig = "data.nutanix_sda_cluster_config_v2.test"

func TestAccV2NutanixSdaClusterConfigDatasource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSdaClusterConfigDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameClusterConfig, "ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameClusterConfig, "system_defined_policy_ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameClusterConfig, "is_enabled"),
				),
			},
		},
	})
}

func testAccSdaClusterConfigDatasourceConfig() string {
	return `
		data "nutanix_sda_cluster_configs_v2" "policies" {
			system_defined_policy_ext_id = data.nutanix_sda_cluster_configs_v2.list.cluster_configs.0.system_defined_policy_ext_id
		}

		data "nutanix_sda_cluster_configs_v2" "list" {
			system_defined_policy_ext_id = local.sda_policy_ext_id
		}

		locals {
			sda_policy_ext_id = data.nutanix_sda_cluster_configs_v2.all.cluster_configs.0.system_defined_policy_ext_id
		}

		data "nutanix_sda_cluster_configs_v2" "all" {
			system_defined_policy_ext_id = "` + testVars.SystemDefinedPolicyExtID + `"
		}

		data "nutanix_sda_cluster_config_v2" "test" {
			system_defined_policy_ext_id = "` + testVars.SystemDefinedPolicyExtID + `"
			ext_id                       = "` + testVars.ClusterExtID + `"
		}
	`
}
