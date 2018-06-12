package nutanix

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNutanixCategoryValue_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixCategoryValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixCategoryValueConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixCategoryValueExists("nutanix_category_value.test"),
				),
			},
		},
	})
}

func testAccCheckNutanixCategoryValueExists(n string) resource.TestCheckFunc {
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

func testAccCheckNutanixCategoryValueDestroy(s *terraform.State) error {

	conn := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_category_value" {
			continue
		}
		for {
			_, err := conn.API.V3.GetCategoryValue(rs.Primary.Attributes["name"], rs.Primary.ID)
			if err != nil {
				if strings.Contains(fmt.Sprint(err), "CATEGORY_NAME_VALUE_MISMATCH") {
					return nil
				}
				return err
			}
			time.Sleep(3000 * time.Millisecond)
		}

	}

	return nil
}

func testAccNutanixCategoryValueConfig() string {
	return `
resource "nutanix_category_key" "test-category-key"{
    name = "app-suppport-1"
	description = "App Support Category Key"
}


resource "nutanix_category_value" "test"{
    name = "${nutanix_category_key.test-category-key.id}"
	description = "Test Category Value"
	value = "test-value"
}
`
}

func Test_resourceNutanixCategoryValue(t *testing.T) {
	tests := []struct {
		name string
		want *schema.Resource
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resourceNutanixCategoryValue(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resourceNutanixCategoryValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_resourceNutanixCategoryValueCreateOrUpdate(t *testing.T) {
	type args struct {
		resourceData *schema.ResourceData
		meta         interface{}
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
			if err := resourceNutanixCategoryValueCreateOrUpdate(tt.args.resourceData, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("resourceNutanixCategoryValueCreateOrUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_resourceNutanixCategoryValueRead(t *testing.T) {
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
			if err := resourceNutanixCategoryValueRead(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("resourceNutanixCategoryValueRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_resourceNutanixCategoryValueDelete(t *testing.T) {
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
			if err := resourceNutanixCategoryValueDelete(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("resourceNutanixCategoryValueDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getCategoryValueSchema(t *testing.T) {
	tests := []struct {
		name string
		want map[string]*schema.Schema
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCategoryValueSchema(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCategoryValueSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testAccCheckNutanixCategoryValueExists(t *testing.T) {
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
			if got := testAccCheckNutanixCategoryValueExists(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("testAccCheckNutanixCategoryValueExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testAccCheckNutanixCategoryValueDestroy(t *testing.T) {
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
			if err := testAccCheckNutanixCategoryValueDestroy(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("testAccCheckNutanixCategoryValueDestroy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_testAccNutanixCategoryValueConfig(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := testAccNutanixCategoryValueConfig(); got != tt.want {
				t.Errorf("testAccNutanixCategoryValueConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
