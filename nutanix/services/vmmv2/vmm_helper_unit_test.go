package vmmv2

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/ahv/config"
)

func TestUnitExpandVirtualEthernetNic(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandVirtualEthernetNic(nil, nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandVirtualEthernetNic([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}}, nil)
		_ = got
	})

}

func TestUnitExpandSriovNic(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandSriovNic(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandSriovNic([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandDpOffloadNic(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandDpOffloadNic(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandDpOffloadNic([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitFlattenVirtualEthernetNicAsBackingInfo(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenVirtualEthernetNicAsBackingInfo(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVirtualEthernetNicAsBackingInfo(&config.VirtualEthernetNic{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenVirtualEthernetNicAsBackingInfo(&config.VirtualEthernetNic{})
		_ = got
	})

}

func TestUnitFlattenVirtualEthernetNicModel(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVirtualEthernetNicModel(func() *config.VirtualEthernetNicModel { var v config.VirtualEthernetNicModel; return &v }())
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenVirtualEthernetNicModel(func() *config.VirtualEthernetNicModel { v := config.VirtualEthernetNicModel(1); return &v }())
		_ = got
	})

}

func TestUnitFlattenVirtualEthernetNicNetworkInfo(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenVirtualEthernetNicNetworkInfo(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVirtualEthernetNicNetworkInfo(&config.VirtualEthernetNicNetworkInfo{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenVirtualEthernetNicNetworkInfo(&config.VirtualEthernetNicNetworkInfo{})
		_ = got
	})

}

func TestUnitFlattenSriovNicAsBackingInfo(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenSriovNicAsBackingInfo(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenSriovNicAsBackingInfo(&config.SriovNic{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenSriovNicAsBackingInfo(&config.SriovNic{})
		_ = got
	})

}

func TestUnitFlattenDpOffloadNicAsBackingInfo(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenDpOffloadNicAsBackingInfo(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenDpOffloadNicAsBackingInfo(&config.DpOffloadNic{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenDpOffloadNicAsBackingInfo(&config.DpOffloadNic{})
		_ = got
	})

}

func TestUnitFlattenSriovNicNetworkInfo(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenSriovNicNetworkInfo(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenSriovNicNetworkInfo(&config.SriovNicNetworkInfo{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenSriovNicNetworkInfo(&config.SriovNicNetworkInfo{})
		_ = got
	})

}

func TestUnitFlattenDpOffloadNicNetworkInfo(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenDpOffloadNicNetworkInfo(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenDpOffloadNicNetworkInfo(&config.DpOffloadNicNetworkInfo{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenDpOffloadNicNetworkInfo(&config.DpOffloadNicNetworkInfo{})
		_ = got
	})

}

func TestUnitFlattenHostPcieDeviceReference(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenHostPcieDeviceReference(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenHostPcieDeviceReference(&config.HostPcieDeviceReference{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenHostPcieDeviceReference(&config.HostPcieDeviceReference{})
		_ = got
	})

}

func TestUnitFlattenNicProfileReference(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenNicProfileReference(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenNicProfileReference(&config.NicProfileReference{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenNicProfileReference(&config.NicProfileReference{})
		_ = got
	})

}
