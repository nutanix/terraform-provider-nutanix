package microsegv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func schemaForIPAddressOrFQDN() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"fqdn": {
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
					},
				},
			},
			"ipv4": {
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
			"ipv6": {
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
		},
	}
}

func schemaForIPAddressOrFQDNComputed() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"fqdn": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
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
	}
}

func schemaForMatchingCriteria() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"criteria": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"match_entity": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"match_field": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"match_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func schemaForMatchingCriteriaComputed() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"criteria": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"match_entity": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"match_field": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"match_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func schemaForAdInfo() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"directory_service_reference": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"object_identifier": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"object_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func schemaForAdInfoComputed() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"directory_service_reference": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"object_identifier": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"object_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
