package networkingv2

import (
	"testing"

	import1 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/microseg/v4/config"
)

func TestUnitFlattenNetworkSecurityPolicy(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = flattenNetworkSecurityPolicy(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenNetworkSecurityPolicy([]import1.NetworkSecurityPolicy{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenNetworkSecurityPolicy([]import1.NetworkSecurityPolicy{import1.NetworkSecurityPolicy{}})
		_ = got
	})

}
