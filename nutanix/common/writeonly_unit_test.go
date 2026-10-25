//go:build unit

package common

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-cty/cty/convert"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestGetWriteOnlyString_PlainString(t *testing.T) {
	r := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				WriteOnly: true,
			},
		},
	}
	d := r.Data(&terraform.InstanceState{
		RawConfig: cty.ObjectVal(map[string]cty.Value{
			"password": cty.StringVal("s3cret"),
		}),
	})
	got, ok := GetWriteOnlyString(d, "password")
	if !ok || got != "s3cret" {
		t.Fatalf("GetWriteOnlyString: got (%q, %v), want (s3cret, true)", got, ok)
	}
}

func TestGetWriteOnlyString_ConvertibleNonString(t *testing.T) {
	num := cty.NumberIntVal(42)
	converted, err := convert.Convert(num, cty.String)
	if err != nil {
		t.Fatalf("convert sanity: %v", err)
	}
	if converted.AsString() != "42" {
		t.Fatalf("convert sanity: got %q", converted.AsString())
	}

	r := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:      schema.TypeString,
				Optional:  true,
				WriteOnly: true,
			},
		},
	}
	d := r.Data(&terraform.InstanceState{
		RawConfig: cty.ObjectVal(map[string]cty.Value{
			"token": num,
		}),
	})
	got, ok := GetWriteOnlyString(d, "token")
	if !ok || got != "42" {
		t.Fatalf("GetWriteOnlyString convertible: got (%q, %v), want (42, true)", got, ok)
	}
}

func TestGetWriteOnlyString_DynamicPseudoType(t *testing.T) {
	// Values sourced via jsondecode/locals arrive as DynamicPseudoType but
	// hold a concrete string value. GetWriteOnlyString must convert via
	// cty/convert.Convert before AsString to read these correctly.
	r := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"private_key": {
				Type:      schema.TypeString,
				Optional:  true,
				WriteOnly: true,
			},
			"ca_chain": {
				Type:      schema.TypeString,
				Optional:  true,
				WriteOnly: true,
			},
		},
	}

	// Unknown DynamicPseudoType value - should return false.
	d := r.Data(&terraform.InstanceState{
		RawConfig: cty.ObjectVal(map[string]cty.Value{
			"private_key": cty.DynamicVal,
		}),
	})
	if got, ok := GetWriteOnlyString(d, "private_key"); ok {
		t.Fatalf("GetWriteOnlyString unknown dynamic: got (%q, %v), want (\"\", false)", got, ok)
	}

	// Known DynamicPseudoType holding a string. Use convert to produce
	// DynamicPseudoType values (cty.NewValue is unavailable in older versions).
	pk, err := convert.Convert(cty.StringVal("super-secret-key"), cty.DynamicPseudoType)
	if err != nil {
		t.Fatalf("convert to dynamic: %v", err)
	}
	ca, err := convert.Convert(cty.StringVal("ca-chain-data"), cty.DynamicPseudoType)
	if err != nil {
		t.Fatalf("convert to dynamic: %v", err)
	}
	d2 := r.Data(&terraform.InstanceState{
		RawConfig: cty.ObjectVal(map[string]cty.Value{
			"private_key": pk,
			"ca_chain":    ca,
		}),
	})
	if got, ok := GetWriteOnlyString(d2, "private_key"); !ok || got != "super-secret-key" {
		t.Fatalf("GetWriteOnlyString dynamic string: got (%q, %v), want (super-secret-key, true)", got, ok)
	}
	if got, ok := GetWriteOnlyString(d2, "ca_chain"); !ok || got != "ca-chain-data" {
		t.Fatalf("GetWriteOnlyString dynamic string: got (%q, %v), want (ca-chain-data, true)", got, ok)
	}
}

func TestGetWriteOnlyString_NestedDynamicPseudoType(t *testing.T) {
	// Nested path where an intermediate value is DynamicPseudoType (e.g. from
	// jsondecode). The path helper must convert before navigating.
	r := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"credential": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"password": {
							Type:      schema.TypeString,
							Optional:  true,
							WriteOnly: true,
						},
					},
				},
			},
		},
	}

	innerObj := cty.ObjectVal(map[string]cty.Value{
		"password": func() cty.Value {
			v, err := convert.Convert(cty.StringVal("nested-secret"), cty.DynamicPseudoType)
			if err != nil {
				t.Fatalf("convert: %v", err)
			}
			return v
		}(),
	})
	dynList, err := convert.Convert(cty.ListVal([]cty.Value{innerObj}), cty.DynamicPseudoType)
	if err != nil {
		t.Fatalf("convert list to dynamic: %v", err)
	}

	d := r.Data(&terraform.InstanceState{
		RawConfig: cty.ObjectVal(map[string]cty.Value{
			"credential": dynList,
		}),
	})
	if got, ok := GetWriteOnlyString(d, "credential.0.password"); !ok || got != "nested-secret" {
		t.Fatalf("GetWriteOnlyString nested dynamic: got (%q, %v), want (nested-secret, true)", got, ok)
	}
}

func TestGetWriteOnlyString_Null(t *testing.T) {
	r := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				WriteOnly: true,
			},
		},
	}
	d := r.Data(&terraform.InstanceState{
		RawConfig: cty.ObjectVal(map[string]cty.Value{
			"password": cty.NullVal(cty.String),
		}),
	})
	if got, ok := GetWriteOnlyString(d, "password"); ok || got != "" {
		t.Fatalf("GetWriteOnlyString null: got (%q, %v), want (\"\", false)", got, ok)
	}
}
