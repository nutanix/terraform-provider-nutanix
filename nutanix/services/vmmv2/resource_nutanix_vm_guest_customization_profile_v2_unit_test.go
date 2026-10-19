package vmmv2

import (
	"testing"
)

func TestUnitExpandVmGcProfileConfig(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandVmGcProfileConfig(map[string]interface{}{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandVmGcProfileConfig(map[string]interface{}{"sysprep_config": []interface{}{map[string]interface{}{"ext_id": "ext-1"}}, "customization": []interface{}{map[string]interface{}{"ext_id": "ext-1"}}, "sysprep_params": []interface{}{map[string]interface{}{"ext_id": "ext-1"}}, "answer_file": []interface{}{map[string]interface{}{"ext_id": "ext-1"}}, "unattend_xml": "UNKNOWN"})
		_ = got
	})

}

func TestUnitExpandSysprepParams(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandSysprepParams(map[string]interface{}{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandSysprepParams(map[string]interface{}{"name": "UNKNOWN"})
		_ = got
	})

}
