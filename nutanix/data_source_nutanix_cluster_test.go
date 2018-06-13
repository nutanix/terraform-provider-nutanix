package nutanix

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func TestAccNutanixClusterDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.nutanix_cluster.cluster", "id"),
				),
			},
		},
	})
}

const testAccClusterDataSourceConfig = `
data "nutanix_clusters" "clusters" {
	metadata = {
		length = 2
	}
}


data "nutanix_cluster" "cluster" {
	cluster_id = "${data.nutanix_clusters.clusters.entities.1.metadata.uuid}"
}`

func Test_dataSourceNutanixCluster(t *testing.T) {
	tests := []struct {
		name string
		want *schema.Resource
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dataSourceNutanixCluster(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dataSourceNutanixCluster() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataSourceNutanixClusterRead(t *testing.T) {
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
			if err := dataSourceNutanixClusterRead(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("dataSourceNutanixClusterRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getDataSourceClusterSchema(t *testing.T) {
	tests := []struct {
		name string
		want map[string]*schema.Schema
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDataSourceClusterSchema(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDataSourceClusterSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}
