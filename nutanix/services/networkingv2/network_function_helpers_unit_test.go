package networkingv2

import (
	"testing"

	"github.com/stretchr/testify/require"

	import1 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
)

func TestUnitFlattenDataPlaneHealthCheckConfig(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenDataPlaneHealthCheckConfig(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenDataPlaneHealthCheckConfig(&import1.DataPlaneHealthCheckConfig{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenDataPlaneHealthCheckConfig(&import1.DataPlaneHealthCheckConfig{})
		_ = got
	})

}

func TestUnitExpandDataPlaneHealthCheckConfig(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandDataPlaneHealthCheckConfig(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandDataPlaneHealthCheckConfig([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitFlattenNicPairs(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenNicPairs(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenNicPairs([]import1.NicPair{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenNicPairs([]import1.NicPair{import1.NicPair{}})
		_ = got
	})

}

func TestUnitExpandNicPairs(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandNicPairs(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandNicPairs([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}
