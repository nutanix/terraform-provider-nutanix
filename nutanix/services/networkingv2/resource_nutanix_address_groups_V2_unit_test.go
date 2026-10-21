package networkingv2

import (
	"testing"
)

func TestUnitExpandIPv4AddressList(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = expandIPv4AddressList(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandIPv4AddressList([]interface{}{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIPv4AddressList([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandIPv4Range(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = expandIPv4Range(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandIPv4Range([]interface{}{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIPv4Range([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}
