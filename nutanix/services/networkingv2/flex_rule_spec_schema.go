package networkingv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func flexRuleSpecSchema(computed bool) map[string]*schema.Schema {
	if computed {
		return flexRuleSpecSchemaComputed()
	}
	return flexRuleSpecSchemaResource()
}

func flexRuleSpecSchemaResource() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"action": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"ALLOW", "DENY", "REJECT"}, false),
		},
		"direction": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"IN", "OUT", "IN_OUT"}, false),
		},
		"applied_to_entity_group_references": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"dest_entity_group_references": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"dest_subnet": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"value": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"prefix_length": {
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"src_entity_group_references": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"src_subnet": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"value": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"prefix_length": {
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"is_all_protocol_allowed": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"service_group_references": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"tcp_services": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"start_port": {
						Type:     schema.TypeInt,
						Required: true,
					},
					"end_port": {
						Type:     schema.TypeInt,
						Required: true,
					},
				},
			},
		},
		"udp_services": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"start_port": {
						Type:     schema.TypeInt,
						Required: true,
					},
					"end_port": {
						Type:     schema.TypeInt,
						Required: true,
					},
				},
			},
		},
		"icmp_services": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"is_all_allowed": {
						Type:     schema.TypeBool,
						Optional: true,
						Computed: true,
					},
					"type": {
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
					"code": {
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"icmpv6_services": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"is_all_allowed": {
						Type:     schema.TypeBool,
						Optional: true,
						Computed: true,
					},
					"type": {
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
					"code": {
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"should_allow_any_dst": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"should_allow_any_src": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"network_function_reference": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"priority": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"ip_version": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringInSlice([]string{"IPV4", "IPV6", "IPV4_IPV6"}, false),
		},
		"is_system_rule": {
			Type:     schema.TypeBool,
			Computed: true,
		},
	}
}

func flexRuleSpecSchemaComputed() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"action": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"direction": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"applied_to_entity_group_references": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"dest_entity_group_references": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"dest_subnet": {
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
		"src_entity_group_references": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"src_subnet": {
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
		"is_all_protocol_allowed": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"service_group_references": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"tcp_services": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"start_port": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"end_port": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
		"udp_services": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"start_port": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"end_port": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
		"icmp_services": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"is_all_allowed": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"type": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"code": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
		"icmpv6_services": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"is_all_allowed": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"type": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"code": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
		"should_allow_any_dst": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"should_allow_any_src": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"network_function_reference": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"priority": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"ip_version": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"is_system_rule": {
			Type:     schema.TypeBool,
			Computed: true,
		},
	}
}
