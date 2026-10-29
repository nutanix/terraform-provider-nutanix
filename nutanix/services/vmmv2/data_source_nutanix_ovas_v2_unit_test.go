package vmmv2

import (
	"testing"

	import1 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/content"
)

func TestUnitFlattenOvaEntities(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = flattenOvaEntities(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenOvaEntities([]import1.Ova{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenOvaEntities([]import1.Ova{import1.Ova{}})
		_ = got
	})

}
