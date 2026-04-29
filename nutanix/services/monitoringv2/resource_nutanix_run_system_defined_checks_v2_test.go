package monitoringv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameRunSystemDefinedChecks = "nutanix_run_system_defined_checks_v2.test"

func TestAccV2NutanixRunSystemDefinedChecksResource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRunSystemDefinedChecksConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRunSystemDefinedChecks, "cluster_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameRunSystemDefinedChecks, "task_ext_id"),
					resource.TestCheckResourceAttr(resourceNameRunSystemDefinedChecks, "should_run_all_checks", "true"),
					resource.TestCheckResourceAttr(resourceNameRunSystemDefinedChecks, "additional_recipients.#", "1"),
					resource.TestCheckResourceAttr(resourceNameRunSystemDefinedChecks, "additional_recipients.0", "test@nutanix.com"),
				),
			},
		},
	})
}

func TestAccV2NutanixRunSystemDefinedChecksResource_WithAdditionalRecipients(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRunSystemDefinedChecksConfigWithRecipients(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRunSystemDefinedChecks, "cluster_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameRunSystemDefinedChecks, "task_ext_id"),
					resource.TestCheckResourceAttr(resourceNameRunSystemDefinedChecks, "should_run_all_checks", "true"),
					resource.TestCheckResourceAttr(resourceNameRunSystemDefinedChecks, "additional_recipients.#", "1"),
					resource.TestCheckResourceAttr(resourceNameRunSystemDefinedChecks, "additional_recipients.0", "test@example.com"),
				),
			},
		},
	})
}

func testRunSystemDefinedChecksConfigBasic() string {
	return `
data "nutanix_clusters_v2" "clusters" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}

locals {
  clusterExtID = data.nutanix_clusters_v2.clusters.cluster_entities[0].ext_id
}

resource "nutanix_run_system_defined_checks_v2" "test" {
  cluster_ext_id                                = local.clusterExtID
  should_run_all_checks                         = true
  additional_recipients                         = ["test@nutanix.com"]
}
`
}

func testRunSystemDefinedChecksConfigWithRecipients() string {
	return `
data "nutanix_clusters_v2" "clusters" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}

locals {
  clusterExtID = data.nutanix_clusters_v2.clusters.cluster_entities[0].ext_id
}

resource "nutanix_run_system_defined_checks_v2" "test" {
  cluster_ext_id                                = local.clusterExtID
  should_run_all_checks                         = true
  additional_recipients                         = ["test@example.com"]
}
`
}
