package nutanix

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func TestAccNutanixClustersDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClustersDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nutanix_clusters.basic_web", "entities.#", "2"),
				),
			},
		},
	})
}

// Lookup based on InstanceID
const testAccClustersDataSourceConfig = `
data "nutanix_clusters" "basic_web" {
	metadata = {
		length = 2
	}
}`

func Test_dataSourceNutanixClusters(t *testing.T) {
	tests := []struct {
		name string
		want *schema.Resource
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dataSourceNutanixClusters(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dataSourceNutanixClusters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataSourceNutanixClustersRead(t *testing.T) {
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
			if err := dataSourceNutanixClustersRead(tt.args.d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("dataSourceNutanixClustersRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getDataSourceClustersSchema(t *testing.T) {
	tests := []struct {
		name string
		want map[string]*schema.Schema
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDataSourceClustersSchema(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDataSourceClustersSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}
