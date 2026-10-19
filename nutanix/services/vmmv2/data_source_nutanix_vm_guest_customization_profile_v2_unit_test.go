package vmmv2

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/ahv/config"
)

func TestUnitFlattenVmGuestCustomizationProfileEntity(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVmGuestCustomizationProfileEntity(config.VmGuestCustomizationProfile{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenVmGuestCustomizationProfileEntity(config.VmGuestCustomizationProfile{})
		_ = got
	})

}

func TestUnitFlattenVmGcProfileConfig(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenVmGcProfileConfig(nil))
	})

}

func TestUnitFlattenSysprepConfig(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenSysprepConfig(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenSysprepConfig(&config.VmGcProfileSysprepConfig{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenSysprepConfig(&config.VmGcProfileSysprepConfig{})
		_ = got
	})

}

func TestUnitFlattenSysprepCustomization(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenSysprepCustomization(nil))
	})

}

func TestUnitFlattenSysprepParams(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenSysprepParams(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenSysprepParams(&config.VmGcProfileSysprepParams{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenSysprepParams(&config.VmGcProfileSysprepParams{})
		_ = got
	})

}

func TestUnitFlattenGeneralSettings(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenGeneralSettings(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenGeneralSettings(&config.VmGcProfileGeneralSettings{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenGeneralSettings(&config.VmGcProfileGeneralSettings{})
		_ = got
	})

}

func TestUnitFlattenComputerName(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenComputerName(nil))
	})

}

func TestUnitFlattenLocaleSettings(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenLocaleSettings(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenLocaleSettings(&config.VmGcProfileLocaleSettings{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenLocaleSettings(&config.VmGcProfileLocaleSettings{})
		_ = got
	})

}

func TestUnitFlattenNetworkSettings(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenNetworkSettings(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenNetworkSettings(&config.VmGcProfileNetworkSettings{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenNetworkSettings(&config.VmGcProfileNetworkSettings{})
		_ = got
	})

}

func TestUnitFlattenWorkgroupOrDomainInfo(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenWorkgroupOrDomainInfo(nil))
	})

}

func TestUnitFlattenAnswerFile(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenAnswerFile(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenAnswerFile(&config.VmGcProfileAnswerFile{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenAnswerFile(&config.VmGcProfileAnswerFile{})
		_ = got
	})

}
