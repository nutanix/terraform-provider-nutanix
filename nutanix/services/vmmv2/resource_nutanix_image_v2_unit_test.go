package vmmv2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitExpandOneOfImageChecksum(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandOneOfImageChecksum(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandOneOfImageChecksum([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandOneOfImageSource(t *testing.T) {
	t.Parallel()

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandOneOfImageSource([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandURLBasicAuth(t *testing.T) {
	t.Parallel()

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandURLBasicAuth([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitFlattenStringValue(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = flattenStringValue(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenStringValue([]interface{}{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenStringValue([]interface{}{})
		_ = got
	})

}
