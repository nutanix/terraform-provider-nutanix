package networkingv2

import (
	"testing"

	"github.com/stretchr/testify/require"

	networkingConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v17/models/networking/v4/config"
)

func TestUnitFlattenVirtualSwitchesEntities(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenVirtualSwitchesEntities(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVirtualSwitchesEntities([]networkingConfig.VirtualSwitch{})
		require.Nil(t, got)
	})

}
