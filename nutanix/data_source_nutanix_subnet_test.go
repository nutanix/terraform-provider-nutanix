package nutanix

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func TestAccNutanixSubnetDataSource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetDataSourceConfig(acctest.RandIntRange(0, 500)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_subnet.test", "prefix_length", "24"),
					resource.TestCheckResourceAttr(
						"data.nutanix_subnet.test", "subnet_type", "VLAN"),
				),
			},
		},
	})
}

func testAccSubnetDataSourceConfig(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {
  metadata = {
    length = 2
  }
}

output "cluster" {
  value = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
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
`, r, r)
}

func Test_dataSourceNutanixSubnet(t *testing.T) {
	tests := []struct {
		name string
		want *schema.Resource
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dataSourceNutanixSubnet(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dataSourceNutanixSubnet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataSourceNutanixSubnetRead(t *testing.T) {
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
			if err := dataSourceNutanixSubnetRead(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("dataSourceNutanixSubnetRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getDataSourceSubnetSchema(t *testing.T) {
	tests := []struct {
		name string
		want map[string]*schema.Schema
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDataSourceSubnetSchema(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDataSourceSubnetSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testAccSubnetDataSourceConfig(t *testing.T) {
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
			if got := testAccSubnetDataSourceConfig(tt.args.r); got != tt.want {
				t.Errorf("testAccSubnetDataSourceConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
