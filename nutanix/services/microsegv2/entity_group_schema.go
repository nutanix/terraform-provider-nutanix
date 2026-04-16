package microsegv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import2 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/microseg/v4/config"
)

var (
	AllowedSelectedByEnums = []string{
		import2.ALLOWEDSELECTBY_IP_VALUES.GetName(),
		import2.ALLOWEDSELECTBY_EXT_ID.GetName(),
		import2.ALLOWEDSELECTBY_CATEGORY_EXT_ID.GetName(),
		import2.ALLOWEDSELECTBY_LABELS.GetName(),
		import2.ALLOWEDSELECTBY_NAME.GetName(),
	}
	AllowedTypeEnums = []string{
		import2.ALLOWEDTYPE_KUBE_NAMESPACE.GetName(),
		import2.ALLOWEDTYPE_SUBNET.GetName(),
		import2.ALLOWEDTYPE_VM.GetName(),
		import2.ALLOWEDTYPE_VPC.GetName(),
		import2.ALLOWEDTYPE_KUBE_SERVICE.GetName(),
		import2.ALLOWEDTYPE_KUBE_CLUSTER.GetName(),
		import2.ALLOWEDTYPE_KUBE_PODS.GetName(),
		import2.ALLOWEDTYPE_ADDRESS_GROUP.GetName(),
	}
	ExceptSelectedByEnums = []string{
		import2.EXCEPTSELECTBY_EXT_ID.GetName(),
		import2.EXCEPTSELECTBY_IP_VALUES.GetName(),
	}
	ExceptTypeEnums = []string{
		import2.EXCEPTTYPE_ADDRESS_GROUP.GetName(),
	}
)

func resourceSchemaForAllowedEntity() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"selected_by": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(AllowedSelectedByEnums, false),
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(AllowedTypeEnums, false),
			},
			"addresses": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4_addresses": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
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
			},
			"ip_ranges": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4_ranges": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_ip": {
										Type:     schema.TypeString,
										Required: true,
									},
									"end_ip": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"kube_entities": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"reference_ext_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceSchemaForExceptConfig() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"entities": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem:     resourceSchemaForExceptEntity(),
				},
			},
		},
	}
}

func resourceSchemaForExceptEntity() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"selected_by": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(ExceptSelectedByEnums, false),
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(ExceptTypeEnums, false),
			},
			"addresses": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4_addresses": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
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
			},
			"ip_ranges": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4_ranges": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_ip": {
										Type:     schema.TypeString,
										Required: true,
									},
									"end_ip": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"reference_ext_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
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

func schemaForAllowedConfig() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"entities": {
					Type:     schema.TypeList,
					Computed: true,
					Elem:     schemaForAllowedEntity(),
				},
			},
		},
	}
}

func schemaForAllowedEntity() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"selected_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4_addresses": {
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
			},
			"ip_ranges": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4_ranges": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"end_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"kube_entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"reference_ext_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func schemaForExceptConfig() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"entities": {
					Type:     schema.TypeList,
					Computed: true,
					Elem:     schemaForExceptEntity(),
				},
			},
		},
	}
}

func schemaForExceptEntity() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"selected_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4_addresses": {
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
			},
			"ip_ranges": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4_ranges": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"end_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"reference_ext_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceSchemaForAllowedConfig() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"entities": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					MaxItems: 3,
					MinItems: 1,
					Elem:     resourceSchemaForAllowedEntity(),
				},
			},
		},
	}
}

func resourceEntityGroupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"allowed_config": resourceSchemaForAllowedConfig(),
		"except_config":  resourceSchemaForExceptConfig(),
		"policy_ext_ids": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"ext_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"creation_time": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"last_update_time": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"links": schemaForLinks(),
		"owner_ext_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tenant_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
