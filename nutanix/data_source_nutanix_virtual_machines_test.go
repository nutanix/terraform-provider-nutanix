package nutanix

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/v3"
)

func TestAccNutanixVMSDataSource_basic(t *testing.T) {
	//Skipped because this test didn't pass in GCP environment
	if isGCPEnvironment() {
		t.Skip()
	}

	rInt := acctest.RandIntRange(0, 500)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVMSSDataSourceConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_virtual_machines.unittest", "entities.#", "2"),
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
  name = "unittest-dou-vm1"

  cluster_reference = {
	  kind = "cluster"
	  uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }

  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 186
  power_state          = "ON"

}

resource "nutanix_virtual_machine" "vm2" {
  name = "unittest-dou-vm2"

  cluster_reference = {
	  kind = "cluster"
	  uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }

  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 186
  power_state          = "ON"

}

data "nutanix_virtual_machines" "unittest" {
	metadata = {
		length = 2
	}
}`)
}

func Test_dataSourceNutanixVirtualMachines(t *testing.T) {
	tests := []struct {
		name string
		want *schema.Resource
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dataSourceNutanixVirtualMachines(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dataSourceNutanixVirtualMachines() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataSourceNutanixVirtualMachinesRead(t *testing.T) {
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
			if err := dataSourceNutanixVirtualMachinesRead(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("dataSourceNutanixVirtualMachinesRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_setGPUList(t *testing.T) {
	type args struct {
		gpu []*v3.VMGpuOutputStatus
	}
	tests := []struct {
		name string
		args args
		want []map[string]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := setGPUList(tt.args.gpu); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setGPUList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setNutanixGuestTools(t *testing.T) {
	type args struct {
		guest *v3.GuestToolsStatus
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := setNutanixGuestTools(tt.args.guest); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setNutanixGuestTools() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setNicList(t *testing.T) {
	type args struct {
		nics []*v3.VMNicOutputStatus
	}
	tests := []struct {
		name string
		args args
		want []map[string]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := setNicList(tt.args.nics); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setNicList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getDataSourceVMSSchema(t *testing.T) {
	tests := []struct {
		name string
		want map[string]*schema.Schema
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDataSourceVMSSchema(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDataSourceVMSSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testAccVMSSDataSourceConfig(t *testing.T) {
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
			if got := testAccVMSSDataSourceConfig(tt.args.r); got != tt.want {
				t.Errorf("testAccVMSSDataSourceConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
