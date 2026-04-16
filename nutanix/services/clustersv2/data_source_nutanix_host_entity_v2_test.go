package clustersv2_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameHostEntity = "data.nutanix_host_v2.test"

func TestAccV2NutanixHostEntityDatasource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testHostEntityDatasourceV4Config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameHostEntity, "cluster.#"),
					resource.TestCheckResourceAttrSet(datasourceNameHostEntity, "cluster.0.uuid"),
					resource.TestCheckResourceAttrSet(datasourceNameHostEntity, "ext_id"),
				),
			},
		},
	})
}

func TestAccV2NutanixHostEntityDatasource_WithNoClsExtId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testHostEntityDatasourceV4WithoutClsExtIDConfig(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccV2NutanixHostEntityDatasource_WithNoHostExtId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testHostEntityDatasourceV4WithoutHostExtIDConfig(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func testHostEntityDatasourceV4Config() string {
	return `
	data "nutanix_hosts_v2" "test" {  
	}

	data "nutanix_host_v2" "test" {
		cluster_ext_id = data.nutanix_hosts_v2.test.host_entities[0].cluster[0].uuid
		ext_id = data.nutanix_hosts_v2.test.host_entities[0].ext_id  
	}
	`
}

func testHostEntityDatasourceV4WithoutClsExtIDConfig() string {
	return `
		data "nutanix_host_v2" "test" {
			ext_id = "00000000-0000-0000-0000-000000000000"
		}
	`
}

func testHostEntityDatasourceV4WithoutHostExtIDConfig() string {
	return `
		data "nutanix_host_v2" "test" {
			cluster_ext_id = "00000000-0000-0000-0000-000000000000"
		}
	`
}
