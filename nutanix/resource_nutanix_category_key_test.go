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
)

func TestAccNutanixCategoryKey_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixCategoryKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixCategoryKeyConfig(acctest.RandIntRange(0, 500)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixCategoryKeyExists("nutanix_category_key.test"),
				),
			},
		},
	})
}

func testAccCheckNutanixCategoryKeyExists(n string) resource.TestCheckFunc {
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

func testAccCheckNutanixCategoryKeyDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_category_key" {
			continue
		}
		for {
			_, err := conn.API.V3.GetCategoryKey(rs.Primary.ID)
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

func testAccNutanixCategoryKeyConfig(r int) string {
	return fmt.Sprintf(`
resource "nutanix_category_key" "test"{
    name = "app-support-%d"
    description = "App Support CategoryKey"
}
`, r)
}

func Test_resourceNutanixCategoryKey(t *testing.T) {
	tests := []struct {
		name string
		want *schema.Resource
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resourceNutanixCategoryKey(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resourceNutanixCategoryKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_resourceNutanixCategoryKeyCreateOrUpdate(t *testing.T) {
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
			if err := resourceNutanixCategoryKeyCreateOrUpdate(tt.args.resourceData, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("resourceNutanixCategoryKeyCreateOrUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_resourceNutanixCategoryKeyRead(t *testing.T) {
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
			if err := resourceNutanixCategoryKeyRead(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("resourceNutanixCategoryKeyRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_resourceNutanixCategoryKeyDelete(t *testing.T) {
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
			if err := resourceNutanixCategoryKeyDelete(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("resourceNutanixCategoryKeyDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getCategoryKeySchema(t *testing.T) {
	tests := []struct {
		name string
		want map[string]*schema.Schema
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCategoryKeySchema(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCategoryKeySchema() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testAccCheckNutanixCategoryKeyExists(t *testing.T) {
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
			if got := testAccCheckNutanixCategoryKeyExists(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("testAccCheckNutanixCategoryKeyExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testAccCheckNutanixCategoryKeyDestroy(t *testing.T) {
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
			if err := testAccCheckNutanixCategoryKeyDestroy(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("testAccCheckNutanixCategoryKeyDestroy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_testAccNutanixCategoryKeyConfig(t *testing.T) {
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
			if got := testAccNutanixCategoryKeyConfig(tt.args.r); got != tt.want {
				t.Errorf("testAccNutanixCategoryKeyConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
