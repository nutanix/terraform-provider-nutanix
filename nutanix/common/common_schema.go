package common

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

const ipv4PrefixLengthDefaultValue = 32
const ipv6PrefixLengthDefaultValue = 128

// LinksSchema returns a schema definition for a list of links, each containing 'rel' and 'href' fields.
func LinksSchema() *schema.Schema {
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

// S
func SchemaForIPList(includeFQDN bool) *schema.Resource {
	schemaMap := map[string]*schema.Schema{
		"ipv4": SchemaForValuePrefixLengthResource(ipv4PrefixLengthDefaultValue),
		"ipv6": SchemaForValuePrefixLengthResource(ipv6PrefixLengthDefaultValue),
	}

	if includeFQDN {
		schemaMap["fqdn"] = &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"value": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		}
	}

	return &schema.Resource{
		Schema: schemaMap,
	}
}

func SchemaForValuePrefixLengthResource(defaultPrefixLength int) *schema.Schema {
	return &schema.Schema{
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
					Default:  defaultPrefixLength,
				},
			},
		},
	}
}
