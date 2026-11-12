package vmmv2

import (
	"testing"
)

func TestUnitExpandEntityFilter(t *testing.T) {
	t.Parallel()

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandEntityFilter([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}
