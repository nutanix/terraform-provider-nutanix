package nutanix

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNutanixImagesDataSource_basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesDataSourceConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixImagesExists("data.nutanix_images.test1"),
				),
			},
		},
	})
}

func testAccCheckNutanixImagesExists(n string) resource.TestCheckFunc {
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

func testAccImagesDataSourceConfig(r int) string {
	return fmt.Sprintf(`
resource "nutanix_image" "test" {
  name        = "Ubuntu"
  description = "Ubuntu"
  source_uri  = "http://archive.ubuntu.com/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/mini.iso"
}


data "nutanix_images" "test1" {
	metadata {
		length = 1
	}
}`)
}

func Test_dataSourceNutanixImages(t *testing.T) {
	tests := []struct {
		name string
		want *schema.Resource
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dataSourceNutanixImages(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dataSourceNutanixImages() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataSourceNutanixImagesRead(t *testing.T) {
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
			if err := dataSourceNutanixImagesRead(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("dataSourceNutanixImagesRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getDataSourceImagesSchema(t *testing.T) {
	tests := []struct {
		name string
		want map[string]*schema.Schema
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDataSourceImagesSchema(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDataSourceImagesSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testAccCheckNutanixImagesExists(t *testing.T) {
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
			if got := testAccCheckNutanixImagesExists(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("testAccCheckNutanixImagesExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testAccImagesDataSourceConfig(t *testing.T) {
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
			if got := testAccImagesDataSourceConfig(tt.args.r); got != tt.want {
				t.Errorf("testAccImagesDataSourceConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
