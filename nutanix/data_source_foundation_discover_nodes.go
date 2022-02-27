package nutanix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFoundationDiscoverNodes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFoundationDiscoverNodesRead,
		Schema: map[string]*schema.Schema{
			"entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"model": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"chasis_n": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"block_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nodes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"foundation_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ipv6_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"node_uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"current_network_interface": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"node_position": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"hypervisor": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"configured": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"nos_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cluster_id": {
										Type:     schema.TypeInt, //check for type
										Computed: true,
									},
									"current_cvm_vlan_tag": {
										Type:     schema.TypeInt, //check for type
										Computed: true,
									},
									"hypervisor_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cvm_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"model": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceFoundationDiscoverNodesRead(ctx context.Context, d *schema.ResourceData, mata interface{}) diag.Diagnostics {
	return nil
}
