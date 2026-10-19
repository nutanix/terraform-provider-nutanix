package vmmv2

import (
	"testing"

	"github.com/stretchr/testify/require"

	vmmConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/config"
)

func TestUnitFlattenTemplatePlacementPolicies(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenTemplatePlacementPolicies(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenTemplatePlacementPolicies([]vmmConfig.TemplatePlacementPolicy{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenTemplatePlacementPolicies([]vmmConfig.TemplatePlacementPolicy{vmmConfig.TemplatePlacementPolicy{}})
		_ = got
	})

}
