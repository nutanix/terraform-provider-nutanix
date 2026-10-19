package networkingv2

import (
	"testing"

	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/models/microseg/v4/config"
)

func TestUnitFlattenAddressGroupsEntities(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = flattenAddressGroupsEntities(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenAddressGroupsEntities([]import1.AddressGroup{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenAddressGroupsEntities([]import1.AddressGroup{import1.AddressGroup{}})
		_ = got
	})

}
