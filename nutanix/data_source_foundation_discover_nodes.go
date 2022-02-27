package nutanix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/foundation"
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
						"chassis_n": {
							Type:     schema.TypeInt,
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
									"svm_ip": {
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

func dataSourceFoundationDiscoverNodesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	conn := meta.(*Client).FoundationClientAPI

	resp, err := conn.Networking.DiscoverNodes(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	entities := make([]map[string]interface{}, len(*resp))
	for k, v := range *resp {
		entity := make(map[string]interface{})
		entity["model"] = v.Model
		entity["chassis_n"] = v.ChassisN
		entity["block_id"] = v.BlockID

		entity["nodes"] = flattenDiscoveredNodes(v.Nodes)

		entities[k] = entity
	}
	// log.Println(entities)
	setErr := d.Set("entities", entities)
	if setErr != nil {
		return diag.FromErr(err)
	}
	d.SetId(resource.UniqueId())
	return nil
}

func flattenDiscoveredNodes(nodesList []foundation.DiscoveredNode) []map[string]interface{} {
	nodes := make([]map[string]interface{}, len(nodesList))
	for k, v := range nodesList {
		node := make(map[string]interface{})

		node["foundation_version"] = v.FoundationVersion
		node["ipv6_address"] = v.Ipv6Address
		node["node_uuid"] = v.NodeUUID
		node["current_network_interface"] = v.CurrentNetworkInterface
		node["node_position"] = v.NodePosition
		node["hypervisor"] = v.Hypervisor
		node["configured"] = v.Configured
		node["nos_version"] = v.NosVersion

		if v.ClusterID != nil && v.ClusterID != "" {
			node["cluster_id"] = int64((v.ClusterID).(float64))
		} else {
			node["cluster_id"] = int64(0)
		}

		node["current_cvm_vlan_tag"] = v.CurrentCvmVlanTag

		node["hypervisor_version"] = v.HypervisorVersion
		node["svm_ip"] = v.SvmIP
		node["model"] = v.Model

		nodes[k] = node
	}
	return nodes
}
