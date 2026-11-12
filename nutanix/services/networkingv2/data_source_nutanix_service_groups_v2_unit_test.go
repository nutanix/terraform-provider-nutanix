package networkingv2

import (
	"testing"

	import1 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/microseg/v4/config"
)

func TestUnitFlattenServiceGroupsEntities(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = flattenServiceGroupsEntities(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenServiceGroupsEntities([]import1.ServiceGroup{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenServiceGroupsEntities([]import1.ServiceGroup{import1.ServiceGroup{}})
		_ = got
	})

}
