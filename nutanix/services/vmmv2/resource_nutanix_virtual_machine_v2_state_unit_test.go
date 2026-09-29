// Unit tests for backfilling the wait_for_ip_* defaults on refresh, so upgrading
// the provider or importing a VM does not plan an in-place update for them.
package vmmv2

import (
	"context"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// refreshData returns ResourceData as Read sees it on a refresh: the given state
// attributes (all others null), a raw state built from them, and no raw config.
func refreshData(t *testing.T, r *schema.Resource, attrs map[string]cty.Value) *schema.ResourceData {
	t.Helper()
	ty := r.CoreConfigSchema().ImpliedType()
	vals := map[string]cty.Value{}
	for name, attrTy := range ty.AttributeTypes() {
		vals[name] = cty.NullVal(attrTy)
	}
	for name, v := range attrs {
		vals[name] = v
	}
	rawState := cty.ObjectVal(vals)
	state, err := r.ShimInstanceStateFromValue(rawState)
	if err != nil {
		t.Fatal(err)
	}
	state.RawState = rawState
	return r.Data(state)
}

// State written by a provider version without the wait_for_ip_* arguments, or
// an imported VM, gets the defaults on refresh, so the next plan has no diff for them.
func TestSetWaitForIPDefaultsIfUnsetMissingFromState(t *testing.T) {
	r := ResourceNutanixVirtualMachineV2()
	d := refreshData(t, r, map[string]cty.Value{"id": cty.StringVal("vm-1"), "name": cty.StringVal("vm")})
	if err := setWaitForIPDefaultsIfUnset(d); err != nil {
		t.Fatal(err)
	}
	state := d.State()
	if got := state.Attributes["wait_for_ip_timeout"]; got != "5" {
		t.Errorf("wait_for_ip_timeout = %q, want \"5\"", got)
	}
	if got := state.Attributes["wait_for_ip_routable"]; got != "true" {
		t.Errorf("wait_for_ip_routable = %q, want \"true\"", got)
	}

	// With the defaults in state, a config that omits both arguments plans no change to them.
	diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(map[string]interface{}{"name": "vm"}), nil)
	if err != nil {
		t.Fatal(err)
	}
	if diff != nil {
		for _, name := range []string{"wait_for_ip_timeout", "wait_for_ip_routable"} {
			// The legacy diff also lists attributes whose value is unchanged.
			if a, ok := diff.Attributes[name]; ok && (a.Old != a.New || a.NewRemoved) {
				t.Errorf("plan diff for %s: %q -> %q", name, a.Old, a.New)
			}
		}
	}
}

// During create or update (raw config present) nothing is backfilled: the
// planned values are already in d, and overwriting them would make the applied
// state disagree with the plan.
func TestSetWaitForIPDefaultsIfUnsetSkipsApply(t *testing.T) {
	r := ResourceNutanixVirtualMachineV2()
	d := refreshData(t, r, map[string]cty.Value{"id": cty.StringVal("vm-1")})
	ty := r.CoreConfigSchema().ImpliedType()
	vals := map[string]cty.Value{}
	for name, attrTy := range ty.AttributeTypes() {
		vals[name] = cty.NullVal(attrTy)
	}
	vals["name"] = cty.StringVal("vm")
	state := d.State()
	state.RawState = d.GetRawState()
	state.RawConfig = cty.ObjectVal(vals)
	d = r.Data(state)

	if err := setWaitForIPDefaultsIfUnset(d); err != nil {
		t.Fatal(err)
	}
	if got, ok := d.State().Attributes["wait_for_ip_timeout"]; ok && got != "" {
		t.Errorf("wait_for_ip_timeout = %q, want it left unset", got)
	}
}

// Values already in state, including ones that differ from the defaults, are left alone.
func TestSetWaitForIPDefaultsIfUnsetKeepsStateValues(t *testing.T) {
	r := ResourceNutanixVirtualMachineV2()
	d := refreshData(t, r, map[string]cty.Value{
		"id":                   cty.StringVal("vm-1"),
		"wait_for_ip_timeout":  cty.NumberIntVal(0),
		"wait_for_ip_routable": cty.False,
	})
	if err := setWaitForIPDefaultsIfUnset(d); err != nil {
		t.Fatal(err)
	}
	if got := d.Get("wait_for_ip_timeout").(int); got != 0 {
		t.Errorf("wait_for_ip_timeout = %d, want 0 (left as set)", got)
	}
	if got := d.Get("wait_for_ip_routable").(bool); got {
		t.Error("wait_for_ip_routable = true, want false (left as set)")
	}
}
