package networkingv2

import (
	"testing"

	"github.com/stretchr/testify/require"

	networkingConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v17/models/networking/v4/config"
)

func TestUnitFlattenVpcVirtualSwitchMappings(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenVpcVirtualSwitchMappings(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVpcVirtualSwitchMappings([]networkingConfig.VpcVirtualSwitchMapping{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenVpcVirtualSwitchMappings([]networkingConfig.VpcVirtualSwitchMapping{networkingConfig.VpcVirtualSwitchMapping{}})
		_ = got
	})

}
