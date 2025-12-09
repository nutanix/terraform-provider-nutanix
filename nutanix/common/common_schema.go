package common

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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

// SchemaForIPList returns a schema definition for a list of IP addresses, including IPv4 and IPv6, and optionally FQDN.
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

func HashIPItem(v interface{}) int {
	m := v.(map[string]interface{})
	var hashKey string

	// Handle IPv4
	if ipv4List, ok := m["ipv4"].([]interface{}); ok && len(ipv4List) > 0 {
		for _, ip4 := range ipv4List {
			ipMap := ip4.(map[string]interface{})
			value := ipMap["value"].(string)
			prefix := 0
			if p, ok := ipMap["prefix_length"]; ok {
				prefix = p.(int)
			}
			hashKey += fmt.Sprintf("ipv4-%s-%d;", value, prefix)
		}
	}

	// Handle IPv6
	if ipv6List, ok := m["ipv6"].([]interface{}); ok && len(ipv6List) > 0 {
		for _, ip6 := range ipv6List {
			ipMap := ip6.(map[string]interface{})
			value := ipMap["value"].(string)
			prefix := 0
			if p, ok := ipMap["prefix_length"]; ok {
				prefix = p.(int)
			}
			hashKey += fmt.Sprintf("ipv6-%s-%d;", value, prefix)
		}
	}

	// Handle FQDN (if present)
	if fqdnList, ok := m["fqdn"].([]interface{}); ok && len(fqdnList) > 0 {
		for _, fq := range fqdnList {
			fqMap := fq.(map[string]interface{})
			fqVal := fqMap["value"].(string)
			hashKey += fmt.Sprintf("fqdn-%s;", fqVal)
		}
	}

	return schema.HashString(hashKey)
}
