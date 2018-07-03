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

func TestAccNutanixSubnetsDataSource_basic(t *testing.T) {
	//Skipped because this test didn't pass in GCP environment
	if isGCPEnvironment() {
		t.Skip()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetsDataSourceConfig(acctest.RandIntRange(0, 500)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixSubnetsExists("data.nutanix_subnets.test1"),
				),
			},
		},
	})
}

func testAccCheckNutanixSubnetsExists(n string) resource.TestCheckFunc {
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

func testAccSubnetsDataSourceConfig(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {
  metadata = {
    length = 2
  }
}

resource "nutanix_subnet" "test" {
	name = "dou_vlan0_test_%d"

	cluster_reference = {
	  kind = "cluster"
	  uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  	}

	vlan_id = %d
	subnet_type = "VLAN"

	prefix_length = 24
	default_gateway_ip = "192.168.0.1"
	subnet_ip = "192.168.0.0"
	#ip_config_pool_list_ranges = ["192.168.0.5", "192.168.0.100"]

	dhcp_options {
		boot_file_name = "bootfile"
		tftp_server_name = "192.168.0.252"
		domain_name = "nutanix"
	}

	dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
	dhcp_domain_search_list = ["nutanix.com", "calm.io"]

}

data "nutanix_subnet" "test" {
	subnet_id = "${nutanix_subnet.test.id}"
}

data "nutanix_subnets" "test1" {
	metadata {
		length = 1
	}
}`, r, r)
}

func Test_dataSourceNutanixSubnets(t *testing.T) {
	tests := []struct {
		name string
		want *schema.Resource
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dataSourceNutanixSubnets(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dataSourceNutanixSubnets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataSourceNutanixSubnetsRead(t *testing.T) {
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
			if err := dataSourceNutanixSubnetsRead(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("dataSourceNutanixSubnetsRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getDataSourceSubnetsSchema(t *testing.T) {
	tests := []struct {
		name string
		want map[string]*schema.Schema
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDataSourceSubnetsSchema(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDataSourceSubnetsSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testAccCheckNutanixSubnetsExists(t *testing.T) {
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
			if got := testAccCheckNutanixSubnetsExists(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("testAccCheckNutanixSubnetsExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testAccSubnetsDataSourceConfig(t *testing.T) {
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
			if got := testAccSubnetsDataSourceConfig(tt.args.r); got != tt.want {
				t.Errorf("testAccSubnetsDataSourceConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
