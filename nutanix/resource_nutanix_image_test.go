package nutanix

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-nutanix/client/v3"
)

func TestAccNutanixImage_basic(t *testing.T) {
	r := rand.Int31()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixImageConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixImageExists("nutanix_image.test"),
				),
			},
			{
				Config: testAccNutanixImageConfigUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixImageExists("nutanix_image.test"),
				),
			},
		},
	})
}

// Skipping this test for now, since it is difficult to implement on the CI
func TestAccNutanixImage_basic_uploadLocal(t *testing.T) {

	t.Skip()

	r := "path-to-file"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixImageLocalConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixImageExists("nutanix_image.test"),
				),
			},
			{
				Config: testAccNutanixImageLocalConfigUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixImageExists("nutanix_image.test"),
				),
			},
		},
	})
}

func testAccCheckNutanixImageExists(n string) resource.TestCheckFunc {
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

func testAccCheckNutanixImageDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_image" {
			continue
		}
		for {
			_, err := conn.API.V3.GetImage(rs.Primary.ID)
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

func testAccNutanixImageConfig(r int32) string {
	return fmt.Sprintf(`
resource "nutanix_image" "test" {
  name        = "Ubuntu"
  description = "Ubuntu"
  source_uri  = "http://archive.ubuntu.com/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/mini.iso"
}
`)
}

func testAccNutanixImageConfigUpdate(r int32) string {
	return fmt.Sprintf(`
resource "nutanix_image" "test" {
  name        = "Ubuntu Updated"
  description = "Ubuntu Updated"
  source_uri  = "http://archive.ubuntu.com/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/mini.iso"
}
`)
}

func testAccNutanixImageLocalConfig(r string) string {
	return fmt.Sprintf(`
resource "nutanix_image" "test" {
  name        = "Ubuntu"
  description = "Ubuntu"
  source_path  = "%s"
}
`, r)
}

func testAccNutanixImageLocalConfigUpdate(r string) string {
	return fmt.Sprintf(`
resource "nutanix_image" "test" {
  name        = "Ubuntu Updated"
  description = "Ubuntu Updated"
  source_path  = "%s"
}
`, r)
}

func Test_resourceNutanixImage(t *testing.T) {
	tests := []struct {
		name string
		want *schema.Resource
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resourceNutanixImage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resourceNutanixImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_resourceNutanixImageCreate(t *testing.T) {
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
			if err := resourceNutanixImageCreate(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("resourceNutanixImageCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_resourceNutanixImageRead(t *testing.T) {
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
			if err := resourceNutanixImageRead(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("resourceNutanixImageRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_resourceNutanixImageUpdate(t *testing.T) {
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
			if err := resourceNutanixImageUpdate(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("resourceNutanixImageUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_resourceNutanixImageDelete(t *testing.T) {
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
			if err := resourceNutanixImageDelete(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("resourceNutanixImageDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getImageSchema(t *testing.T) {
	tests := []struct {
		name string
		want map[string]*schema.Schema
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getImageSchema(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getImageSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getImageMetadaAttributes(t *testing.T) {
	type args struct {
		d        *schema.ResourceData
		metadata *v3.ImageMetadata
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
			if err := getImageMetadaAttributes(tt.args.d, tt.args.metadata); (err != nil) != tt.wantErr {
				t.Errorf("getImageMetadaAttributes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getImageResource(t *testing.T) {
	type args struct {
		d     *schema.ResourceData
		image *v3.ImageResources
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
			if err := getImageResource(tt.args.d, tt.args.image); (err != nil) != tt.wantErr {
				t.Errorf("getImageResource() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_resourceNutanixImageExists(t *testing.T) {
	type args struct {
		conn *v3.Client
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    *string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resourceNutanixImageExists(tt.args.conn, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("resourceNutanixImageExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("resourceNutanixImageExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_imageStateRefreshFunc(t *testing.T) {
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
			if got := imageStateRefreshFunc(tt.args.client, tt.args.uuid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("imageStateRefreshFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testAccCheckNutanixImageExists(t *testing.T) {
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
			if got := testAccCheckNutanixImageExists(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("testAccCheckNutanixImageExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testAccCheckNutanixImageDestroy(t *testing.T) {
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
			if err := testAccCheckNutanixImageDestroy(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("testAccCheckNutanixImageDestroy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_testAccNutanixImageConfig(t *testing.T) {
	type args struct {
		r int32
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
			if got := testAccNutanixImageConfig(tt.args.r); got != tt.want {
				t.Errorf("testAccNutanixImageConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testAccNutanixImageConfigUpdate(t *testing.T) {
	type args struct {
		r int32
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
			if got := testAccNutanixImageConfigUpdate(tt.args.r); got != tt.want {
				t.Errorf("testAccNutanixImageConfigUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testAccNutanixImageLocalConfig(t *testing.T) {
	type args struct {
		r string
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
			if got := testAccNutanixImageLocalConfig(tt.args.r); got != tt.want {
				t.Errorf("testAccNutanixImageLocalConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testAccNutanixImageLocalConfigUpdate(t *testing.T) {
	type args struct {
		r string
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
			if got := testAccNutanixImageLocalConfigUpdate(tt.args.r); got != tt.want {
				t.Errorf("testAccNutanixImageLocalConfigUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}
