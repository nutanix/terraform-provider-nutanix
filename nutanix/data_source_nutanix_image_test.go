package nutanix

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func TestAccNutanixImageDataSource_basic(t *testing.T) {

	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccImageDataSourceConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_image.test", "description", "Ubuntu mini ISO"),
					resource.TestCheckResourceAttr(
						"data.nutanix_image.test",
						"source_uri",
						"http://archive.ubuntu.com/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/mini.iso"),
				),
			},
		},
	})
}

func testAccImageDataSourceConfig(rNumber int) string {
	return fmt.Sprintf(`
resource "nutanix_image" "test" {
  name        = "Ubuntu-%d"
  description = "Ubuntu mini ISO"
  source_uri  = "http://archive.ubuntu.com/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/mini.iso"
}


data "nutanix_image" "test" {
	image_id = "${nutanix_image.test.id}"
}
`, rNumber)
}

func Test_dataSourceNutanixImage(t *testing.T) {
	tests := []struct {
		name string
		want *schema.Resource
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dataSourceNutanixImage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dataSourceNutanixImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataSourceNutanixImageRead(t *testing.T) {
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
			if err := dataSourceNutanixImageRead(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("dataSourceNutanixImageRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getDataSourceImageSchema(t *testing.T) {
	tests := []struct {
		name string
		want map[string]*schema.Schema
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDataSourceImageSchema(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDataSourceImageSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testAccImageDataSourceConfig(t *testing.T) {
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
			if got := testAccImageDataSourceConfig(tt.args.r); got != tt.want {
				t.Errorf("testAccImageDataSourceConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
