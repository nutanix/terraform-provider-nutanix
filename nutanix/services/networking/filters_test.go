package networking

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
)

func TestFilters_replaceFilterPrefixes(t *testing.T) {
	mappings := map[string]string{
		"source":  "target",
		"replace": "replaced",
	}
	originalFilters := []*client.AdditionalFilter{
		{
			Name:   "source",
			Values: []string{"value1"},
		},
		{
			Name:   "dnd",
			Values: []string{"value"},
		},
		{
			Name:   "replace.me.too",
			Values: []string{"value2", "value3"},
		},
	}
	expected := []*client.AdditionalFilter{
		{
			Name:   "target",
			Values: []string{"value1"},
		},
		{
			Name:   "dnd",
			Values: []string{"value"},
		},
		{
			Name:   "replaced.me.too",
			Values: []string{"value2", "value3"},
		},
	}

	actual := ReplaceFilterPrefixes(originalFilters, mappings)

	for i := 0; i < len(expected); i++ {
		if actual[i].Name != expected[i].Name {
			t.Fatalf("Failed to replace filter mapping from %s to %s, actual: %s", originalFilters[i].Name, expected[i].Name, actual[i].Name)
		}
	}
}
