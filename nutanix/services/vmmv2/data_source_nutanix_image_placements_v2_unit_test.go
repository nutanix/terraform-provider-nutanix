package vmmv2

import (
	"testing"

	import7 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/images/config"
)

func TestUnitFlattenPlacementPolicyEntities(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = flattenPlacementPolicyEntities(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenPlacementPolicyEntities([]import7.PlacementPolicy{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenPlacementPolicyEntities([]import7.PlacementPolicy{import7.PlacementPolicy{}})
		_ = got
	})

}
