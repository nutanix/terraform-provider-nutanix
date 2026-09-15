package microsegv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import2 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/microseg/v4/config"
)

var (
	// matchEntityEnums enumerates the supported MatchEntity values.
	matchEntityEnums = []string{
		import2.MATCHENTITY_VM.GetName(),
	}
	// matchFieldEnums enumerates the supported MatchField values.
	matchFieldEnums = []string{
		import2.MATCHFIELD_NAME.GetName(),
	}
	// matchTypeEnums enumerates the supported MatchType values.
	matchTypeEnums = []string{
		import2.MATCHTYPE_CONTAINS.GetName(),
		import2.MATCHTYPE_ALL.GetName(),
	}
	// adStatusEnums enumerates the supported AdStatus values.
	adStatusEnums = []string{
		import2.ADSTATUS_USABLE.GetName(),
		import2.ADSTATUS_DELETED.GetName(),
		import2.ADSTATUS_DIRECTORY_NOT_CONFIGURED.GetName(),
	}
)

// resourceSchemaForIPAddressOrFQDN returns the editable schema for a domain
// controller entry expressed as an IP address (IPv4/IPv6) or FQDN.
func resourceSchemaForIPAddressOrFQDN() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ipv4": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     resourceSchemaForIPAddressValue(),
			},
			"ipv6": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     resourceSchemaForIPAddressValue(),
			},
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
		},
	}
}

func resourceSchemaForIPAddressValue() *schema.Resource {
	return &schema.Resource{
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
	}
}

// schemaForIPAddressOrFQDN returns the computed-only schema for a domain
// controller entry used by datasources.
func schemaForIPAddressOrFQDN() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ipv4": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemaForIPAddressValue(),
			},
			"ipv6": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemaForIPAddressValue(),
			},
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
		},
	}
}

func schemaForIPAddressValue() *schema.Resource {
	return &schema.Resource{
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
}

// resourceSchemaForMatchingCriteria returns the editable schema for the
// matching criteria used by identity categorization.
func resourceSchemaForMatchingCriteria() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"criteria": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					Description: "Only allowed when match_type is \"CONTAINS\". Must not be set when match_type is \"ALL\".",
				},
				"match_entity": {
					Type:         schema.TypeString,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.StringInSlice(matchEntityEnums, false),
				},
				"match_field": {
					Type:         schema.TypeString,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.StringInSlice(matchFieldEnums, false),
				},
				"match_type": {
					Type:         schema.TypeString,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.StringInSlice(matchTypeEnums, false),
				},
			},
		},
	}
}

// schemaForMatchingCriteria returns the computed-only matching criteria schema.
func schemaForMatchingCriteria() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
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
		},
	}
}

// resourceDirectoryServerConfigSchema is the schema for the Directory Server
// Config resource.
func resourceDirectoryServerConfigSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"directory_service_reference": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"domain_controllers": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem:     resourceSchemaForIPAddressOrFQDN(),
		},
		"is_default_category_enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Can only be set to true when match_type is \"CONTAINS\". Must be false when match_type is \"ALL\".",
		},
		"matching_criterias": resourceSchemaForMatchingCriteria(),
		"should_keep_default_category_on_login": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Can only be true when is_default_category_enabled is true and match_type is \"CONTAINS\". Must be false when match_type is \"ALL\".",
		},
		"project_ext_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"ext_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tenant_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"links": schemaForLinks(),
	}
}

// datasourceDirectoryServerConfigSchema is the computed schema shared by the
// singular and plural Directory Server Config datasources.
func datasourceDirectoryServerConfigSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"directory_service_reference": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"domain_controllers": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     schemaForIPAddressOrFQDN(),
		},
		"is_default_category_enabled": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"matching_criterias": schemaForMatchingCriteria(),
		"should_keep_default_category_on_login": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"project_ext_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ext_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tenant_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"links": schemaForLinks(),
	}
}

// resourceSchemaForAdInfo returns the editable schema for the Active Directory
// object information of a Category Mapping.
func resourceSchemaForAdInfo() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"directory_service_reference": {
					Type:     schema.TypeString,
					Required: true,
				},
				"object_identifier": {
					Type:     schema.TypeString,
					Required: true,
				},
				"object_path": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"status": {
					Type:         schema.TypeString,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.StringInSlice(adStatusEnums, false),
				},
			},
		},
	}
}

// schemaForAdInfo returns the computed-only Active Directory info schema
// used by the singular and plural Category Mapping datasources.
func schemaForAdInfo() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
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
		},
	}
}

// resourceCategoryMappingSchema is the schema for the Category Mapping resource.
func resourceCategoryMappingSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"category_name": {
			Type:     schema.TypeString,
			Default:  "ADGroup",
			Optional: true,
		},
		"category_value": {
			Type:     schema.TypeString,
			Required: true,
		},
		"ad_info": resourceSchemaForAdInfo(),
		"project_ext_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"ext_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tenant_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"links": schemaForLinks(),
	}
}

// datasourceCategoryMappingSchema is the computed schema shared by the singular
// and plural Category Mapping datasources.
func datasourceCategoryMappingSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"category_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"category_value": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ad_info": schemaForAdInfo(),
		"project_ext_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ext_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tenant_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"links": schemaForLinks(),
	}
}
