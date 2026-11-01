package vmmv2

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
)

func TestUnitFlattenVmGuestCustomizationProfileEntities(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenVmGuestCustomizationProfileEntities(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVmGuestCustomizationProfileEntities([]config.VmGuestCustomizationProfile{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenVmGuestCustomizationProfileEntities([]config.VmGuestCustomizationProfile{config.VmGuestCustomizationProfile{}})
		_ = got
	})

}
