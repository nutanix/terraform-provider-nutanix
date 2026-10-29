package vmmv2

import (
	"testing"

	import7 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/images/config"
)

func TestUnitFlattenPlacementType(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenPlacementType(func() *import7.PlacementType { var v import7.PlacementType; return &v }())
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenPlacementType(func() *import7.PlacementType { v := import7.PlacementType(1); return &v }())
		_ = got
	})

}

func TestUnitFlattenEnforcementState(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenEnforcementState(func() *import7.EnforcementState { var v import7.EnforcementState; return &v }())
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenEnforcementState(func() *import7.EnforcementState { v := import7.EnforcementState(1); return &v }())
		_ = got
	})

}

func TestUnitFlattenEntityFilter(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenEntityFilter(&import7.Filter{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenEntityFilter(&import7.Filter{})
		_ = got
	})

}

func TestUnitFlattenType(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenType(func() *import7.FilterMatchType { var v import7.FilterMatchType; return &v }())
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenType(func() *import7.FilterMatchType { v := import7.FilterMatchType(1); return &v }())
		_ = got
	})

}
