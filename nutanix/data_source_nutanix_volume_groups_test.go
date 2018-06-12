package nutanix

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func TestAccNutanixVolumeGroupsDataSource_basic(t *testing.T) {
	// skipping as this API is not yet GA (will GA in upcoming AOS release)
	t.Skip()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeGroupsDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_volume_groups.test", "entities.#", "2"),
				),
			},
		},
	})
}

// Lookup based on InstanceID
const testAccVolumeGroupsDataSourceConfig = `
resource "nutanix_volume_group" "test" {
  name        = "VG Test"
  description = "VG Test Description"
  
}

resource "nutanix_volume_group" "test-1" {
  name        = "VG Test-1"
  description = "VG Test-1 Description"
  
}

data "nutanix_volume_groups" "test" {
	metadata = {
		length = 2
	}
}
`

func Test_dataSourceNutanixVolumeGroups(t *testing.T) {
	tests := []struct {
		name string
		want *schema.Resource
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dataSourceNutanixVolumeGroups(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dataSourceNutanixVolumeGroups() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataSourceNutanixVolumeGroupsRead(t *testing.T) {
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
			if err := dataSourceNutanixVolumeGroupsRead(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("dataSourceNutanixVolumeGroupsRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getDataSourceVolumeGroupsSchema(t *testing.T) {
	tests := []struct {
		name string
		want map[string]*schema.Schema
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDataSourceVolumeGroupsSchema(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDataSourceVolumeGroupsSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}
