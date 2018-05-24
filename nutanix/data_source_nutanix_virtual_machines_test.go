package nutanix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccNutanixVMSDataSource_basic(t *testing.T) {

	rInt := acctest.RandIntRange(0, 500)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVMSSDataSourceConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_virtual_machines.basic_web", "entities.#", "2"),
				),
			},
		},
	})
}

// Lookup based on InstanceID
func testAccVMSSDataSourceConfig(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {
  metadata = {
    length = 2
  }
}

output "cluster" {
  value = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
}

resource "nutanix_virtual_machine" "vm1" {
  name = "test-dou-vm1"

  cluster_reference = {
	  kind = "cluster"
	  uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }

  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 2048
  power_state          = "ON"

  nic_list = [{
    subnet_reference = {
      kind = "subnet"
      uuid = "${nutanix_subnet.test.id}"
    }

    ip_endpoint_list = {
			ip = "192.168.0.10"
			type = "ASSIGNED"
    }
  }]
}

resource "nutanix_virtual_machine" "vm2" {
  name = "test-dou-vm2"

  cluster_reference = {
	  kind = "cluster"
	  uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }

  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 2048
  power_state          = "ON"

  nic_list = [{
    subnet_reference = {
      kind = "subnet"
      uuid = "${nutanix_subnet.test.id}"
    }

    ip_endpoint_list = {
			ip = "192.168.0.11"
			type = "ASSIGNED"
    }
  }]
}

resource "nutanix_subnet" "test" {
  name        = "dou_vlan0_test"

  cluster_reference = {
	  kind = "cluster"
	  uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }

  vlan_id     = %d
  subnet_type = "VLAN"

  prefix_length      = 24
  default_gateway_ip = "192.168.0.1"
  subnet_ip          = "192.168.0.0"

  dhcp_options {
    boot_file_name   = "bootfile"
    tftp_server_name = "192.168.0.252"
    domain_name      = "nutanix"
  }

  dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
  dhcp_domain_search_list      = ["nutanix.com", "calm.io"]
}

data "nutanix_virtual_machines" "basic_web" {
	metadata = {
		length = 2
	}
}`, r)
}
