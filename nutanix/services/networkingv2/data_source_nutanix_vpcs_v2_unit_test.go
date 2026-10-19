package networkingv2

import (
	"testing"

	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v17/models/networking/v4/config"
)

func TestUnitFlattenVPCsEntities(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = flattenVPCsEntities(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVPCsEntities([]import1.Vpc{})
		_ = got
	})

}
