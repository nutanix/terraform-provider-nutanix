package cluster_managementv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameAddSnmpTransport = "nutanix_add_snmp_transport_v2.test"

func TestAccV2NutanixAddSnmpTransportResource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAddSnmpTransportConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameAddSnmpTransport, "cluster_ext_id"),
					resource.TestCheckResourceAttr(resourceNameAddSnmpTransport, "port", "162"),
					resource.TestCheckResourceAttr(resourceNameAddSnmpTransport, "protocol", "UDP"),
				),
			},
		},
	})
}

func testAccAddSnmpTransportConfig() string {
	return fmt.Sprintf(`
	data "nutanix_clusters_v2" "clusters" {}

	locals {
	  cluster_ext_id = [
		for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
		cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
	  ][0]
	}

	resource "nutanix_add_snmp_transport_v2" "test" {
	  cluster_ext_id = local.cluster_ext_id
	  port           = 162
	  protocol       = "UDP"
	  depends_on     = [data.nutanix_clusters_v2.clusters]
	}
`)
}
