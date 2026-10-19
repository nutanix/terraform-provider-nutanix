package networkingv2

import (
	"testing"

	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v17/models/networking/v4/config"
)

func TestUnitFlattenFloatingIPsEntities(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = flattenFloatingIPsEntities(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenFloatingIPsEntities([]import1.FloatingIp{})
		_ = got
	})

}
