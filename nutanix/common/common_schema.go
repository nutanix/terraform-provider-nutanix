package common

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

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
