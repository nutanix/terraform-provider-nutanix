package nutanix

import (
	"reflect"
	"testing"
)

func TestResourceNutanixCategoriesteUpgradeV0(t *testing.T) {
	cases := map[string]struct {
		Attributes map[string]interface{}
		Expected   map[string]interface{}
		Meta       interface{}
	}{
		"v0_without_values": {
			Attributes: map[string]interface{}{
				"categories": make(map[string]interface{}),
			},
			Expected: map[string]interface{}{
				"categories": make([]interface{}, 0),
			},
		},
		"v0_with_values": {
			Attributes: map[string]interface{}{
				"categories": map[string]interface{}{
					"os_type":    "ubuntu",
					"os_version": "18.04",
				},
			},
			Expected: map[string]interface{}{

				"categories": []interface{}{
					map[string]interface{}{"name": "os_type", "value": "ubuntu"},
					map[string]interface{}{"name": "os_version", "value": "18.04"},
				},
			},
		},
		"v0_categories_not_set": {
			Attributes: map[string]interface{}{
				"name": "test-name",
			},
			Expected: map[string]interface{}{
				"name": "test-name",
			},
		},
		"v0_multiple_categories": {
			Attributes: map[string]interface{}{
				"categories": map[string]interface{}{
					"os_type":    "ubuntu",
					"os_version": "18.04",
					"tier":       "application",
					"test":       "test-value",
				},
			},
			Expected: map[string]interface{}{
				"categories": []interface{}{
					map[string]interface{}{"name": "os_type",
						"value": "ubuntu"},
					map[string]interface{}{"name": "os_version",
						"value": "18.04"},
					map[string]interface{}{"name": "test",
						"value": "test-value"},
					map[string]interface{}{"name": "tier",
						"value": "application"},
				},
			},
		},
		"v0_already_migrated": {
			Attributes: map[string]interface{}{
				"categories": []interface{}{
					map[string]interface{}{"name": "os_type",
						"value": "ubuntu"},
					map[string]interface{}{"name": "os_version",
						"value": "18.04"},
					map[string]interface{}{"name": "tier",
						"value": "application"},
				},
			},
			Expected: map[string]interface{}{
				"categories": []interface{}{
					map[string]interface{}{"name": "os_type",
						"value": "ubuntu"},
					map[string]interface{}{"name": "os_version",
						"value": "18.04"},
					map[string]interface{}{"name": "tier",
						"value": "application"},
				},
			},
		},
		"v0_empty_value": {
			Attributes: map[string]interface{}{
				"categories": map[string]interface{}{
					"os_type":    "",
					"os_version": "",
					"tier":       "",
					"test":       "",
				},
			},
			Expected: map[string]interface{}{
				"categories": []interface{}{
					map[string]interface{}{"name": "os_type",
						"value": ""},
					map[string]interface{}{"name": "os_version",
						"value": ""},
					map[string]interface{}{"name": "test",
						"value": ""},
					map[string]interface{}{"name": "tier",
						"value": ""},
				},
			},
		},
	}

	for tn, tc := range cases {
		is, err := resourceNutanixCategoriesMigrateState(tc.Attributes, tc.Meta)

		if err != nil {
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		for k, v := range tc.Expected {
			if !reflect.DeepEqual(is[k], v) {
				t.Fatalf(
					"bad: %s\n\n expected: %#v -> %#v\n got: %#v -> %#v\n in: %#v",
					tn, k, v, k, is[k], is)
			}
		}
	}
}
