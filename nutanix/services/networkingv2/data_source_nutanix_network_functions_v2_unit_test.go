package networkingv2

import (
	"testing"

	"github.com/stretchr/testify/require"

	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v17/models/networking/v4/config"
)

func TestUnitFlattenNetworkFunctionEntities(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenNetworkFunctionEntities(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenNetworkFunctionEntities([]import1.NetworkFunction{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenNetworkFunctionEntities([]import1.NetworkFunction{import1.NetworkFunction{}})
		_ = got
	})

}
