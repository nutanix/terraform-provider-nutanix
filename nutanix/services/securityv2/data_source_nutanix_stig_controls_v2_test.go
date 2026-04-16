package securityv2_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
)

const datasourceNameStigs = "data.nutanix_stigs_v2.test"

func TestAccV2NutanixStigsControlsDatasource_Basic(t *testing.T) {
	limit := 3

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testStigControlsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameStigs, "stigs.#"),
					common.CheckAttributeLength(datasourceNameStigs, "stigs", 1),
				),
			},
			{
				Config: testStigControlsLimitConfig(limit),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameStigs, "stigs.#"),
					common.CheckAttributeLengthEqual(datasourceNameStigs, "stigs", limit),
				),
			},
		},
	})
}

func TestAccV2NutanixStigsControlsDatasource_Filtered(t *testing.T) {
	statusFilter := "status eq Security.Report.StigStatus'APPLICABLE'"
	severityFilter := "severity eq Security.Report.Severity'HIGH'"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testStigControlsFilterConfig(statusFilter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameStigs, "stigs.#"),
					common.CheckAttributeLength(datasourceNameStigs, "stigs", 1),
					resource.TestCheckResourceAttr(datasourceNameStigs, "filter", statusFilter),

					func(s *terraform.State) error {
						kmsAttributes := s.RootModule().Resources[datasourceNameStigs].Primary.Attributes

						stigs := kmsAttributes["stigs.#"]
						//convert stigsCount to int
						stigsCount, err := strconv.Atoi(stigs)
						if err != nil {
							return fmt.Errorf("failed to convert stigs.# to int: %v", err)
						}

						// loop through all returned stigs and make sure all status type is APPLICABLE
						for i := 0; i < stigsCount; i++ {
							if kmsAttributes[fmt.Sprintf("stigs.%d.status", i)] != "APPLICABLE" {
								return fmt.Errorf("expected status of stig %d to be %q, got %q", i, "APPLICABLE", kmsAttributes[fmt.Sprintf("stigs.%d.status", i)])
							}
						}
						return nil
					},
				),
			},
			{
				Config: testStigControlsFilterConfig(severityFilter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameStigs, "stigs.#"),
					common.CheckAttributeLength(datasourceNameStigs, "stigs", 1),
					resource.TestCheckResourceAttr(datasourceNameStigs, "filter", severityFilter),
					func(s *terraform.State) error {
						kmsAttributes := s.RootModule().Resources[datasourceNameStigs].Primary.Attributes

						stigs := kmsAttributes["stigs.#"]
						//convert stigsCount to int
						stigsCount, err := strconv.Atoi(stigs)
						if err != nil {
							return fmt.Errorf("failed to convert stigs.# to int: %v", err)
						}

						// loop through all returned stigs and make sure all severity type is HIGH
						for i := 0; i < stigsCount; i++ {
							if kmsAttributes[fmt.Sprintf("stigs.%d.severity", i)] != "HIGH" {
								return fmt.Errorf("expected severity of stig %d to be %q, got %q", i, "HIGH", kmsAttributes[fmt.Sprintf("stigs.%d.severity", i)])
							}
						}
						return nil
					},
				),
			},
		},
	})
}

func testStigControlsConfig() string {
	return `
data "nutanix_stigs_v2" "test" {}
`
}

func testStigControlsLimitConfig(limit int) string {
	return `
data "nutanix_stigs_v2" "test" {
  limit = ` + fmt.Sprintf("%d", limit) + `
}
`
}

func testStigControlsFilterConfig(filter string) string {
	return `
data "nutanix_stigs_v2" "test" {
  filter = "` + filter + `"
}
`
}
