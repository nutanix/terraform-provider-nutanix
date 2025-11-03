package licensingv2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceClusterLicenseRecommendations = "data.nutanix_cluster_license_recommendations_v2.get_recommendations"

func TestLicensingDataSourceClusterLicenseRecommendationsV2(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceClusterLicenseRecommendationsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceClusterLicenseRecommendations, "acceptances.0.accepted_by.company_name", "Nutanix"),
					resource.TestCheckResourceAttr(datasourceClusterLicenseRecommendations, "acceptances.0.accepted_by.job_title", "MTS"),
					resource.TestCheckResourceAttr(datasourceClusterLicenseRecommendations, "acceptances.0.accepted_by.login_id", "admin"),
					resource.TestCheckResourceAttr(datasourceClusterLicenseRecommendations, "acceptances.0.accepted_by.user_name", "Nutanix"),
				),
			},
		},
	})
}

func testDataSourceClusterLicenseRecommendationsConfig() string {
	return `
	data "nutanix_cluster_license_recommendations_v2" "get_recommendations" {}
  `
}