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

func TestAccNutanixVolumeGroup_basic(t *testing.T) {
	// Skipping as this test needs functional work
	t.Skip()
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVolumeGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixVolumeGroupConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVolumeGroupExists("nutanix_volume_group.test_volume"),
				),
			},
			{
				Config: testAccNutanixVolumeGroupConfigUpdate(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVolumeGroupExists("nutanix_volume_group.test_volume"),
				),
			},
		},
	})
}

func testAccCheckNutanixVolumeGroupExists(n string) resource.TestCheckFunc {
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

func testAccCheckNutanixVolumeGroupDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_volume_group" {
			continue
		}
		for {
			_, err := conn.API.V3.GetVolumeGroup(rs.Primary.ID)
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

func testAccNutanixVolumeGroupConfig(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {
  metadata = {
    length = 2
  }
}
output "cluster" {
  value = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
}
resource "nutanix_volume_group" "test_volume" {
	name        = "Test Volume Group %d"
	description = "Tes Volume Group Description"

  cluster_reference = {
	  kind = "cluster"
	  uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }

}
`, r)
}

func testAccNutanixVolumeGroupConfigUpdate(r int) string {
	return fmt.Sprintf(`
resource "nutanix_volume_group" "test_volume" {
	name        = "Test Volume Group %d"
  description = "Tes Volume Group Description Update"
}
`, r)
}

func Test_resourceNutanixVolumeGroup(t *testing.T) {
	tests := []struct {
		name string
		want *schema.Resource
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resourceNutanixVolumeGroup(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resourceNutanixVolumeGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_resourceNutanixVolumeGroupCreate(t *testing.T) {
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
			if err := resourceNutanixVolumeGroupCreate(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("resourceNutanixVolumeGroupCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_resourceNutanixVolumeGroupRead(t *testing.T) {
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
			if err := resourceNutanixVolumeGroupRead(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("resourceNutanixVolumeGroupRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_resourceNutanixVolumeGroupUpdate(t *testing.T) {
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
			if err := resourceNutanixVolumeGroupUpdate(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("resourceNutanixVolumeGroupUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_resourceNutanixVolumeGroupDelete(t *testing.T) {
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
			if err := resourceNutanixVolumeGroupDelete(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("resourceNutanixVolumeGroupDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getVolumeGroupResources(t *testing.T) {
	type args struct {
		d  *schema.ResourceData
		vg *v3.VolumeGroupResources
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
			if err := getVolumeGroupResources(tt.args.d, tt.args.vg); (err != nil) != tt.wantErr {
				t.Errorf("getVolumeGroupResources() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getVGSchema(t *testing.T) {
	tests := []struct {
		name string
		want map[string]*schema.Schema
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getVGSchema(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getVGSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_volumeGroupStateRefreshFunc(t *testing.T) {
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
			if got := volumeGroupStateRefreshFunc(tt.args.client, tt.args.uuid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("volumeGroupStateRefreshFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testAccCheckNutanixVolumeGroupExists(t *testing.T) {
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
			if got := testAccCheckNutanixVolumeGroupExists(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("testAccCheckNutanixVolumeGroupExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testAccCheckNutanixVolumeGroupDestroy(t *testing.T) {
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
			if err := testAccCheckNutanixVolumeGroupDestroy(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("testAccCheckNutanixVolumeGroupDestroy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_testAccNutanixVolumeGroupConfig(t *testing.T) {
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
			if got := testAccNutanixVolumeGroupConfig(tt.args.r); got != tt.want {
				t.Errorf("testAccNutanixVolumeGroupConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testAccNutanixVolumeGroupConfigUpdate(t *testing.T) {
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
			if got := testAccNutanixVolumeGroupConfigUpdate(tt.args.r); got != tt.want {
				t.Errorf("testAccNutanixVolumeGroupConfigUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}
