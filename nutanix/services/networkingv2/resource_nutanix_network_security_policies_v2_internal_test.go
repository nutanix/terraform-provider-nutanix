package networkingv2

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
)

func appSpec(category, entityGroup cty.Value) cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"application_rule_spec": cty.ListVal([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"secured_group_category_references":    category,
				"secured_group_entity_group_reference": entityGroup,
			}),
		}),
	})
}

func TestCheckSecuredGroupExclusivity_ApplicationRuleSpec(t *testing.T) {
	catList := cty.ListVal([]cty.Value{cty.StringVal("cat-ext-id")})
	emptyCat := cty.ListValEmpty(cty.String)
	nullCat := cty.NullVal(cty.List(cty.String))
	entityGroup := cty.StringVal("eg-ext-id")
	nullEG := cty.NullVal(cty.String)
	emptyEG := cty.StringVal("")

	tests := []struct {
		name      string
		category  cty.Value
		eg        cty.Value
		wantError bool
	}{
		{"scenario C - category only", catList, nullEG, false},
		{"scenario B - entity group only", nullCat, entityGroup, false},
		{"scenario A - both set", catList, entityGroup, true},
		{"neither set (null)", nullCat, nullEG, true},
		{"neither set (empty)", emptyCat, emptyEG, true},
		{"category set, empty entity group", catList, emptyEG, false},
		{"empty category, entity group set", emptyCat, entityGroup, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := appSpec(tt.category, tt.eg)
			err := checkSecuredGroupExclusivity(spec, "application_rule_spec", 0, true)
			if tt.wantError && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantError && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

func TestCheckSecuredGroupExclusivity_NotRequired(t *testing.T) {
	catList := cty.ListVal([]cty.Value{cty.StringVal("cat-ext-id")})
	nullCat := cty.NullVal(cty.List(cty.String))
	entityGroup := cty.StringVal("eg-ext-id")
	nullEG := cty.NullVal(cty.String)

	// When not required, neither being set is allowed.
	spec := cty.ObjectVal(map[string]cty.Value{
		"intra_entity_group_rule_spec": cty.ListVal([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"secured_group_category_references":    nullCat,
				"secured_group_entity_group_reference": nullEG,
			}),
		}),
	})
	if err := checkSecuredGroupExclusivity(spec, "intra_entity_group_rule_spec", 0, false); err != nil {
		t.Fatalf("expected no error when neither set and not required, got %v", err)
	}

	// Both set is still rejected.
	specBoth := cty.ObjectVal(map[string]cty.Value{
		"intra_entity_group_rule_spec": cty.ListVal([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"secured_group_category_references":    catList,
				"secured_group_entity_group_reference": entityGroup,
			}),
		}),
	})
	if err := checkSecuredGroupExclusivity(specBoth, "intra_entity_group_rule_spec", 0, false); err == nil {
		t.Fatalf("expected error when both set, got nil")
	}
}

func TestCheckSecuredGroupExclusivity_UnknownValuesCountAsSet(t *testing.T) {
	// References to other resources are unknown at plan time; both unknown must
	// still be detected as the "both set" conflict.
	unknownCat := cty.UnknownVal(cty.List(cty.String))
	unknownEG := cty.UnknownVal(cty.String)

	spec := appSpec(unknownCat, unknownEG)
	if err := checkSecuredGroupExclusivity(spec, "application_rule_spec", 0, true); err == nil {
		t.Fatalf("expected error when both references are unknown but set, got nil")
	}
}

func TestCheckSecuredGroupExclusivity_BlockAbsent(t *testing.T) {
	// A rule that uses a different spec block must not trigger validation.
	spec := cty.ObjectVal(map[string]cty.Value{
		"application_rule_spec": cty.ListValEmpty(cty.Object(map[string]cty.Type{
			"secured_group_category_references":    cty.List(cty.String),
			"secured_group_entity_group_reference": cty.String,
		})),
	})
	if err := checkSecuredGroupExclusivity(spec, "application_rule_spec", 0, true); err != nil {
		t.Fatalf("expected no error when block absent, got %v", err)
	}
}
