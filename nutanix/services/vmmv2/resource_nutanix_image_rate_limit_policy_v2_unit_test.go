package vmmv2

import (
	"testing"

	"github.com/stretchr/testify/require"

	import7 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/images/config"
)

func TestUnitExpandRateLimitClusterEntityFilter(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandRateLimitClusterEntityFilter(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandRateLimitClusterEntityFilter([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitFlattenRateLimitClusterEntityFilter(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenRateLimitClusterEntityFilter(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenRateLimitClusterEntityFilter(&import7.Filter{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenRateLimitClusterEntityFilter(&import7.Filter{})
		_ = got
	})

}
