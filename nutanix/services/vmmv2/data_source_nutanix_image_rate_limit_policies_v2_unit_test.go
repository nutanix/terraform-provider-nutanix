package vmmv2

import (
	"testing"

	"github.com/stretchr/testify/require"

	import7 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/images/config"
)

func TestUnitFlattenRateLimitPolicies(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenRateLimitPolicies(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenRateLimitPolicies([]import7.RateLimitPolicy{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenRateLimitPolicies([]import7.RateLimitPolicy{import7.RateLimitPolicy{}})
		_ = got
	})

}
