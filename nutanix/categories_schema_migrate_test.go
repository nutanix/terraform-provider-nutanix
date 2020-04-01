package nutanix

import (
	"testing"

	"github.com/hashicorp/terraform/terraform"
)

func TestResourceNutanixCategoriesteUpgradeV0(t *testing.T) {
	cases := map[string]struct {
		StateVersion int
		Attributes   map[string]string
		Expected     map[string]string
		Meta         interface{}
	}{
		"v0_without_values": {
			StateVersion: 0,
			Attributes: map[string]string{
				"categories.%": "0",
			},
			Expected: map[string]string{
				"categories.#": "0",
			},
		},
		"v0_with_values": {
			StateVersion: 0,
			Attributes: map[string]string{
				"categories.%":          "2",
				"categories.os_type":    "ubuntu",
				"categories.os_version": "18.04",
			},
			Expected: map[string]string{
				"categories.#":       "2",
				"categories.0.name":  "os_type",
				"categories.0.value": "ubuntu",
				"categories.1.name":  "os_version",
				"categories.1.value": "18.04",
			},
		},
		"v0_categories_not_set": {
			StateVersion: 0,
			Attributes: map[string]string{
				"name": "test-name",
			},
			Expected: map[string]string{
				"name": "test-name",
			},
		},
		"v0_multiple_categories": {
			StateVersion: 0,
			Attributes: map[string]string{
				"categories.%":          "3",
				"categories.os_type":    "ubuntu",
				"categories.os_version": "18.04",
				"categories.tier":       "application",
				"categories.test":       "test-value",
			},
			Expected: map[string]string{
				"categories.#":       "3",
				"categories.0.name":  "os_type",
				"categories.0.value": "ubuntu",
				"categories.1.name":  "os_version",
				"categories.1.value": "18.04",
				"categories.2.name":  "test",
				"categories.2.value": "test-value",
				"categories.3.name":  "tier",
				"categories.3.value": "application",
			},
		},
		"v0_already_migrated": {
			StateVersion: 0,
			Attributes: map[string]string{
				"categories.#":       "3",
				"categories.0.name":  "os_type",
				"categories.0.value": "ubuntu",
				"categories.1.name":  "os_version",
				"categories.1.value": "18.04",
				"categories.2.name":  "tier",
				"categories.2.value": "application",
			},
			Expected: map[string]string{
				"categories.#":       "3",
				"categories.0.name":  "os_type",
				"categories.0.value": "ubuntu",
				"categories.1.name":  "os_version",
				"categories.1.value": "18.04",
				"categories.2.name":  "tier",
				"categories.2.value": "application",
			},
		},
		"v0_empty_value": {
			StateVersion: 0,
			Attributes: map[string]string{
				"categories.%":          "3",
				"categories.os_type":    "",
				"categories.os_version": "",
				"categories.tier":       "",
			},
			Expected: map[string]string{
				"categories.#":       "3",
				"categories.0.name":  "os_type",
				"categories.0.value": "",
				"categories.1.name":  "os_version",
				"categories.1.value": "",
				"categories.2.name":  "tier",
				"categories.2.value": "",
			},
		},
	}

	for tn, tc := range cases {
		is := &terraform.InstanceState{
			ID:         "i-abc123",
			Attributes: tc.Attributes,
		}
		is, err := resourceNutanixCategoriesMigrateState(
			tc.StateVersion, is, tc.Meta)

		if err != nil {
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		for k, v := range tc.Expected {
			if is.Attributes[k] != v {
				t.Fatalf(
					"bad: %s\n\n expected: %#v -> %#v\n got: %#v -> %#v\n in: %#v",
					tn, k, v, k, is.Attributes[k], is.Attributes)
			}
		}
	}
}
