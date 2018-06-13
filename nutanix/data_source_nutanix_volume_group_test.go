package nutanix

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func TestAccNutanixVolumeGroupDataSource_basic(t *testing.T) {
	// skipping as this API is not yet GA (will GA in upcoming AOS release)
	t.Skip()

	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeGroupDataSourceConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_volume_group.test", "name", "Ubuntu"),
					resource.TestCheckResourceAttr(
						"data.nutanix_volume_group.test", "description", "VG Test Description"),
				),
			},
		},
	})
}

func testAccVolumeGroupDataSourceConfig(r int) string {
	return fmt.Sprintf(`
resource "nutanix_volume_group" "test" {
  name        = "VG Test"
  description = "VG Test Description"
  
}

data "nutanix_volume_group" "test" {
	volume_group_id = "${nutanix_volume_group.test.id}"
}
`)
}

func Test_dataSourceNutanixVolumeGroup(t *testing.T) {
	tests := []struct {
		name string
		want *schema.Resource
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dataSourceNutanixVolumeGroup(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dataSourceNutanixVolumeGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataSourceNutanixVolumeGroupRead(t *testing.T) {
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
			if err := dataSourceNutanixVolumeGroupRead(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("dataSourceNutanixVolumeGroupRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getDataSourceVolumeGroupSchema(t *testing.T) {
	tests := []struct {
		name string
		want map[string]*schema.Schema
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDataSourceVolumeGroupSchema(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDataSourceVolumeGroupSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testAccVolumeGroupDataSourceConfig(t *testing.T) {
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
			if got := testAccVolumeGroupDataSourceConfig(tt.args.r); got != tt.want {
				t.Errorf("testAccVolumeGroupDataSourceConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
