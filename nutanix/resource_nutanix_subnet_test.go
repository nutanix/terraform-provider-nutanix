package nutanix

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-nutanix/client/v3"
)

func TestAccNutanixSubnet_basic(t *testing.T) {
	//Skipped because this test didn't pass in GCP environment
	if isGCPEnvironment() {
		t.Skip()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixSubnetConfig(acctest.RandIntRange(0, 500)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixSubnetExists("nutanix_subnet.next-iac-managed"),
				),
			},
			{
				Config: testAccNutanixSubnetConfigUpdate(acctest.RandIntRange(0, 500)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixSubnetExists("nutanix_subnet.next-iac-managed"),
				),
			},
		},
	})
}

func testAccCheckNutanixSubnetExists(n string) resource.TestCheckFunc {
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

func testAccCheckNutanixSubnetDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_subnet" {
			continue
		}
		if _, err := resourceNutanixSubnetExists(conn.API, rs.Primary.ID); err != nil {
			if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
				return nil
			}
			return err
		}
	}

	return nil
}

func testAccNutanixSubnetConfig(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {
  metadata = {
    length = 2
  }
}

output "cluster" {
  value = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
}

resource "nutanix_subnet" "next-iac-managed" {
  # What cluster will this VLAN live on?
  cluster_reference = {
	kind = "cluster"
	uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }

  # General Information
  name        = "next-iac-managed-%d"
  vlan_id     = %d
  subnet_type = "VLAN"

  # Managed L3 Networks
  # This bit is only needed if you intend to turn on IPAM
  prefix_length = 20

  default_gateway_ip = "10.5.80.1"
  subnet_ip          = "10.5.80.0"

  #dhcp_options {
  #    boot_file_name   = "bootfile"
  #    tftp_server_name = "1.2.3.200"
  #    domain_name      = "nutanix"
  #}

  dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
  dhcp_domain_search_list      = ["nutanix.com", "eng.nutanix.com"]
}
`, r, r)
}

func testAccNutanixSubnetConfigUpdate(r int) string {
	return fmt.Sprintf(`
data "nutanix_clusters" "clusters" {
  metadata = {
    length = 2
  }
}

output "cluster" {
  value = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
}

resource "nutanix_subnet" "next-iac-managed" {
  # What cluster will this VLAN live on?
  cluster_reference = {
	kind = "cluster"
	uuid = "${data.nutanix_clusters.clusters.entities.0.metadata.uuid}"
  }

  # General Information
  name        = "next-iac-managed-%d-Updated"
  vlan_id     = %d
  subnet_type = "VLAN"

  # Managed L3 Networks
  # This bit is only needed if you intend to turn on IPAM
  prefix_length = 20

  default_gateway_ip = "10.5.80.1"
  subnet_ip          = "10.5.80.0"

  #dhcp_options {
  #    boot_file_name   = "bootfile"
  #    tftp_server_name = "1.2.3.200"
  #    domain_name      = "nutanix"
  #}

  dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
  dhcp_domain_search_list      = ["nutanix.com", "eng.nutanix.com"]
}
`, r, r)
}
func Test_resourceNutanixSubnet(t *testing.T) {
	tests := []struct {
		name string
		want *schema.Resource
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resourceNutanixSubnet(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resourceNutanixSubnet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_resourceNutanixSubnetCreate(t *testing.T) {
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
			if err := resourceNutanixSubnetCreate(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("resourceNutanixSubnetCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_resourceNutanixSubnetRead(t *testing.T) {
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
			if err := resourceNutanixSubnetRead(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("resourceNutanixSubnetRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_resourceNutanixSubnetUpdate(t *testing.T) {
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
			if err := resourceNutanixSubnetUpdate(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("resourceNutanixSubnetUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_resourceNutanixSubnetDelete(t *testing.T) {
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
			if err := resourceNutanixSubnetDelete(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("resourceNutanixSubnetDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_resourceNutanixSubnetExists(t *testing.T) {
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
			got, err := resourceNutanixSubnetExists(tt.args.conn, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("resourceNutanixSubnetExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("resourceNutanixSubnetExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getSubnetResources(t *testing.T) {
	type args struct {
		d      *schema.ResourceData
		subnet *v3.SubnetResources
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
			if err := getSubnetResources(tt.args.d, tt.args.subnet); (err != nil) != tt.wantErr {
				t.Errorf("getSubnetResources() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_subnetStateRefreshFunc(t *testing.T) {
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
			if got := subnetStateRefreshFunc(tt.args.client, tt.args.uuid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("subnetStateRefreshFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getSubnetSchema(t *testing.T) {
	tests := []struct {
		name string
		want map[string]*schema.Schema
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSubnetSchema(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSubnetSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testAccCheckNutanixSubnetExists(t *testing.T) {
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
			if got := testAccCheckNutanixSubnetExists(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("testAccCheckNutanixSubnetExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_testAccCheckNutanixSubnetDestroy(t *testing.T) {
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
			if err := testAccCheckNutanixSubnetDestroy(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("testAccCheckNutanixSubnetDestroy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_testAccNutanixSubnetConfig(t *testing.T) {
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
			if got := testAccNutanixSubnetConfig(tt.args.r); got != tt.want {
				t.Errorf("testAccNutanixSubnetConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
