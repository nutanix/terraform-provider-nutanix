package vmmv2

import (
	"testing"

	"github.com/stretchr/testify/require"

	import7 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/images/config"
)

func TestUnitFlattenEffectiveRateLimitPolicies(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenEffectiveRateLimitPolicies(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenEffectiveRateLimitPolicies([]import7.EffectiveRateLimitPolicy{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenEffectiveRateLimitPolicies([]import7.EffectiveRateLimitPolicy{import7.EffectiveRateLimitPolicy{}})
		_ = got
	})

}
