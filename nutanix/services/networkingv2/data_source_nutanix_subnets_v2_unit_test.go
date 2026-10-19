package networkingv2

import (
	"testing"

	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v17/models/networking/v4/config"
)

func TestUnitFlattenSubnetEntities(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = flattenSubnetEntities(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenSubnetEntities([]import1.Subnet{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenSubnetEntities([]import1.Subnet{import1.Subnet{}})
		_ = got
	})

}
