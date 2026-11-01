package vmmv2

import (
	"testing"
)

func TestUnitExpandOneOfCloneVMBootConfig(t *testing.T) {
	t.Parallel()

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandOneOfCloneVMBootConfig([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}
