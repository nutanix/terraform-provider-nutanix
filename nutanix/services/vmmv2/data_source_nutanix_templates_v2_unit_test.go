package vmmv2

import (
	"testing"

	import5 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/content"
)

func TestUnitFlattenTemplatesEntities(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = flattenTemplatesEntities(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenTemplatesEntities([]import5.Template{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenTemplatesEntities([]import5.Template{import5.Template{}})
		_ = got
	})

}
