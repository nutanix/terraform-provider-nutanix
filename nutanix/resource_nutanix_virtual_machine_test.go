package nutanix

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-nutanix/client/v3"
)

func TestAccNutanixVirtualMachine_basic(t *testing.T) {
	r := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixVMConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists("nutanix_virtual_machine.vm1"),
					resource.TestCheckResourceAttr("nutanix_virtual_machine.vm1", "hardware_clock_timezone", "UTC"),
					resource.TestCheckResourceAttr("nutanix_virtual_machine.vm1", "power_state", "ON"),
					resource.TestCheckResourceAttr("nutanix_virtual_machine.vm1", "memory_size_mib", "186"),
					resource.TestCheckResourceAttr("nutanix_virtual_machine.vm1", "num_sockets", "1"),
					resource.TestCheckResourceAttr("nutanix_virtual_machine.vm1", "num_vcpus_per_socket", "1"),
				),
			},
			{
				Config: testAccNutanixVMConfigUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists("nutanix_virtual_machine.vm1"),
					resource.TestCheckResourceAttr("nutanix_virtual_machine.vm1", "hardware_clock_timezone", "UTC"),
					resource.TestCheckResourceAttr("nutanix_virtual_machine.vm1", "power_state", "ON"),
					resource.TestCheckResourceAttr("nutanix_virtual_machine.vm1", "memory_size_mib", "186"),
					resource.TestCheckResourceAttr("nutanix_virtual_machine.vm1", "num_sockets", "1"),
					resource.TestCheckResourceAttr("nutanix_virtual_machine.vm1", "num_vcpus_per_socket", "1"),
				),
			},
		},
	})
}

func testAccCheckNutanixVirtualMachineExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		return nil
	}
}

func testAccCheckNutanixVirtualMachineDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_virtual_machine" {
			continue
		}
		for {
			_, err := conn.API.V3.GetVM(rs.Primary.ID)
			if err != nil {
				if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
					return nil
				}
				return err
			}
			time.Sleep(3000 * time.Millisecond)
		}

	}

	return nil

}

func testAccNutanixVMConfig(r int) string {
	return fmt.Sprint(`
data "nutanix_clusters" "clusters" {
  metadata = {
    length = 2
  }
}
output "cluster" {
  value = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
}
resource "nutanix_virtual_machine" "vm1" {
  name = "test-dou"
  cluster_reference = {
	  kind = "cluster"
	  uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }
  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 186
  power_state          = "ON"
}
`)
}

func testAccNutanixVMConfigUpdate(r int) string {
	return fmt.Sprint(`
data "nutanix_clusters" "clusters" {
  metadata = {
    length = 2
  }
}

output "cluster" {
  value = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
}


resource "nutanix_virtual_machine" "vm1" {
  name = "test-dou-updated"

  cluster_reference = {
	  kind = "cluster"
	  uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }

  num_vcpus_per_socket = 1
  num_sockets          = 1
  memory_size_mib      = 186
  power_state          = "ON"
}
`)
}

func Test_resourceNutanixVirtualMachine(t *testing.T) {
	tests := []struct {
		name string
		want *schema.Resource
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resourceNutanixVirtualMachine(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resourceNutanixVirtualMachine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_resourceNutanixVirtualMachineCreate(t *testing.T) {
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
			if err := resourceNutanixVirtualMachineCreate(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("resourceNutanixVirtualMachineCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_resourceNutanixVirtualMachineRead(t *testing.T) {
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
			if err := resourceNutanixVirtualMachineRead(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("resourceNutanixVirtualMachineRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_resourceNutanixVirtualMachineUpdate(t *testing.T) {
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
			if err := resourceNutanixVirtualMachineUpdate(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("resourceNutanixVirtualMachineUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_resourceNutanixVirtualMachineDelete(t *testing.T) {
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
			if err := resourceNutanixVirtualMachineDelete(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("resourceNutanixVirtualMachineDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_resourceNutanixVirtualMachineExists(t *testing.T) {
	type args struct {
		d    *schema.ResourceData
		meta interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resourceNutanixVirtualMachineExists(tt.args.d, tt.args.meta)
			if (err != nil) != tt.wantErr {
				t.Errorf("resourceNutanixVirtualMachineExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("resourceNutanixVirtualMachineExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getVMResources(t *testing.T) {
	type args struct {
		d  *schema.ResourceData
		vm *v3.VMResources
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
			if err := getVMResources(tt.args.d, tt.args.vm); (err != nil) != tt.wantErr {
				t.Errorf("getVMResources() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_vmStateRefreshFunc(t *testing.T) {
	type args struct {
		client *v3.Client
		uuid   string
	}
	tests := []struct {
		name string
		args args
		want resource.StateRefreshFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := vmStateRefreshFunc(tt.args.client, tt.args.uuid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("vmStateRefreshFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_preFillResUpdateRequest(t *testing.T) {
	type args struct {
		res      *v3.VMResources
		response *v3.VMIntentResponse
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			preFillResUpdateRequest(tt.args.res, tt.args.response)
		})
	}
}

func Test_preFillGTUpdateRequest(t *testing.T) {
	type args struct {
		guestTool *v3.GuestToolsSpec
		response  *v3.VMIntentResponse
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			preFillGTUpdateRequest(tt.args.guestTool, tt.args.response)
		})
	}
}

func Test_preFillGUpdateRequest(t *testing.T) {
	type args struct {
		guest    *v3.GuestCustomization
		response *v3.VMIntentResponse
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			preFillGUpdateRequest(tt.args.guest, tt.args.response)
		})
	}
}

func Test_preFillPWUpdateRequest(t *testing.T) {
	type args struct {
		pw       *v3.VMPowerStateMechanism
		response *v3.VMIntentResponse
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			preFillPWUpdateRequest(tt.args.pw, tt.args.response)
		})
	}
}

func Test_getVMSchema(t *testing.T) {
	tests := []struct {
		name string
		want map[string]*schema.Schema
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getVMSchema(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getVMSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testAccCheckNutanixVirtualMachineExists(t *testing.T) {
	type args struct {
		n string
	}
	tests := []struct {
		name string
		args args
		want resource.TestCheckFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := testAccCheckNutanixVirtualMachineExists(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("testAccCheckNutanixVirtualMachineExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testAccCheckNutanixVirtualMachineDestroy(t *testing.T) {
	type args struct {
		s *terraform.State
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
			if err := testAccCheckNutanixVirtualMachineDestroy(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("testAccCheckNutanixVirtualMachineDestroy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_testAccNutanixVMConfig(t *testing.T) {
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
			if got := testAccNutanixVMConfig(tt.args.r); got != tt.want {
				t.Errorf("testAccNutanixVMConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
