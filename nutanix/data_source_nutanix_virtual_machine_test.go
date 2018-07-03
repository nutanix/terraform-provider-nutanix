package nutanix

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func TestAccNutanixVirtualMachineDataSource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVMDataSourceConfig(acctest.RandIntRange(0, 500)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_virtual_machine.nutanix_virtual_machine", "num_vcpus_per_socket", "1"),
					resource.TestCheckResourceAttr(
						"data.nutanix_virtual_machine.nutanix_virtual_machine", "num_sockets", "1"),
				),
			},
		},
	})
}

func testAccVMDataSourceConfig(r int) string {
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
  name = "test-dou-%d"

  cluster_reference = {
	  kind = "cluster"
	  uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }

  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 186
  power_state          = "ON"
}

data "nutanix_virtual_machine" "nutanix_virtual_machine" {
	vm_id = "${nutanix_virtual_machine.vm1.id}"
}
`, r)
}

func Test_dataSourceNutanixVirtualMachine(t *testing.T) {
	tests := []struct {
		name string
		want *schema.Resource
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dataSourceNutanixVirtualMachine(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dataSourceNutanixVirtualMachine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataSourceNutanixVirtualMachineRead(t *testing.T) {
	type args struct {
		d    *schema.ResourceData
		meta interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := dataSourceNutanixVirtualMachineRead(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("dataSourceNutanixVirtualMachineRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getDataSourceVMSchema(t *testing.T) {
	tests := []struct {
		name string
		want map[string]*schema.Schema
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDataSourceVMSchema(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDataSourceVMSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testAccVMDataSourceConfig(t *testing.T) {
	type args struct {
		r int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := testAccVMDataSourceConfig(tt.args.r); got != tt.want {
				t.Errorf("testAccVMDataSourceConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
