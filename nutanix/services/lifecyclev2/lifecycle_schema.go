package lifecyclev2

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// datasourceLifecycleLinksSchema returns the fully-computed schema for the
// HATEOAS links block shared by lifecyclev2 datasources.
func datasourceLifecycleLinksSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"href": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"rel": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

// datasourceIPAddressSchema returns the fully-computed schema for an IP address
// (ipv4/ipv6 with value + prefix_length) used by the node network details.
func datasourceIPAddressSchema() *schema.Schema {
	valuePrefix := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"prefix_length": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ipv4": {
					Type:     schema.TypeList,
					Computed: true,
					Elem:     valuePrefix,
				},
				"ipv6": {
					Type:     schema.TypeList,
					Computed: true,
					Elem:     valuePrefix,
				},
			},
		},
	}
}

// datasourceNetworkDetailsSchema returns the schema for a NetworkDetails block.
func datasourceNetworkDetailsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"gateway": datasourceIPAddressSchema(),
				"ip":      datasourceIPAddressSchema(),
				"vlan_id": {
					Type:     schema.TypeInt,
					Computed: true,
				},
			},
		},
	}
}

func datasourceComponentNetworkDetailsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"management_network": datasourceNetworkDetailsSchema(),
			},
		},
	}
}

// datasourceClaimTokenNodeSchema returns the fully-computed schema map for a
// claim token node, shared by the singular node datasource and the plural list.
func datasourceClaimTokenNodeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ext_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tenant_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"links": datasourceLifecycleLinksSchema(),
		"created_time": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"host_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"host_version": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"manufacturer": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"model": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"memory_gb": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"managed_by": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"cpu_info": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"capacity_ghz": {
						Type:     schema.TypeFloat,
						Computed: true,
					},
					"logical_core_count": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"manufacturer": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"model": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"socket_count": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
		"identifiers": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"type": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"value": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"network_details": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"bmc":  datasourceNetworkDetailsSchema(),
					"cvm":  datasourceComponentNetworkDetailsSchema(),
					"host": datasourceComponentNetworkDetailsSchema(),
				},
			},
		},
	}
}
