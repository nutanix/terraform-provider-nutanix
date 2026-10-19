package vmmv2

import (
	"testing"

	"github.com/stretchr/testify/require"

	vmmAhvConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/ahv/config"
)

func TestUnitFlattenCapabilities(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenCapabilities(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenCapabilities([]vmmAhvConfig.NgtCapability{})
		require.Nil(t, got)
	})

}
