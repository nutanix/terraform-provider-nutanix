package cluster_managementv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func schemaForIPAddress(required bool) *schema.Schema {
	s := &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ipv4": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"value": {
								Type:     schema.TypeString,
								Required: true,
							},
							"prefix_length": {
								Type:     schema.TypeInt,
								Optional: true,
								Computed: true,
							},
						},
					},
				},
				"ipv6": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"value": {
								Type:     schema.TypeString,
								Required: true,
							},
							"prefix_length": {
								Type:     schema.TypeInt,
								Optional: true,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
	if required {
		s.Required = true
	} else {
		s.Computed = true
	}
	return s
}

func schemaForIPAddressComputed() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ipv4": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
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
					},
				},
				"ipv6": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
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
					},
				},
			},
		},
	}
}

func schemaForLinks() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"rel": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"href": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}
