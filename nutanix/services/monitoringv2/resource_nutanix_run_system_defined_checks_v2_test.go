package monitoringv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameRunChecks = "nutanix_run_system_defined_checks_v2.test"

func TestAccV2NutanixRunSystemDefinedChecksResource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRunSystemDefinedChecksResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRunChecks, "id"),
					resource.TestCheckResourceAttrSet(resourceNameRunChecks, "task_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameRunChecks, "cluster_ext_id"),
					resource.TestCheckResourceAttr(resourceNameRunChecks, "should_run_all_checks", "true"),
					resource.TestCheckResourceAttr(resourceNameRunChecks, "should_send_report_to_configured_recipients", "false"),
				),
			},
		},
	})
}

func testRunSystemDefinedChecksResourceConfig() string {
	return `
data "nutanix_clusters_v2" "clusters" {}

locals {
  cluster_ext_id = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
}

resource "nutanix_run_system_defined_checks_v2" "test" {
  cluster_ext_id                              = local.cluster_ext_id
  should_run_all_checks                       = true
  should_send_report_to_configured_recipients = false
  additional_recipients                       = ["noreply@nutanix.com"]
}
`
}
