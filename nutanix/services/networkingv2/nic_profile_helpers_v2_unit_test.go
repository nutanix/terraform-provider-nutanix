package networkingv2

import (
	"testing"

	"github.com/stretchr/testify/require"

	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v17/models/networking/v4/config"
)

func TestUnitFlattenNicProfileCapabilityConfig(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenNicProfileCapabilityConfig(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenNicProfileCapabilityConfig(&import1.CapabilityConfig{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenNicProfileCapabilityConfig(&import1.CapabilityConfig{})
		_ = got
	})

}

func TestUnitFlattenNicProfileHostNicReferences(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenNicProfileHostNicReferences(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenNicProfileHostNicReferences([]import1.HostNicReference{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenNicProfileHostNicReferences([]import1.HostNicReference{import1.HostNicReference{}})
		_ = got
	})

}

func TestUnitFlattenNicProfileHostNicExtIDs(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenNicProfileHostNicExtIDs(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenNicProfileHostNicExtIDs([]import1.HostNicReference{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenNicProfileHostNicExtIDs([]import1.HostNicReference{import1.HostNicReference{}})
		_ = got
	})

}

func TestUnitFlattenNicProfileEntity(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenNicProfileEntity(map[string]interface{}{}, import1.NicProfile{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenNicProfileEntity(map[string]interface{}{}, import1.NicProfile{})
		_ = got
	})

}

func TestUnitFlattenNicProfileProjectionEntity(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenNicProfileProjectionEntity(map[string]interface{}{}, import1.NicProfileProjection{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenNicProfileProjectionEntity(map[string]interface{}{}, import1.NicProfileProjection{})
		_ = got
	})

}

func TestUnitExpandNicProfileCapabilityConfig(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandNicProfileCapabilityConfig(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandNicProfileCapabilityConfig([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandNicProfileOperatingMode(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandNicProfileOperatingMode("")
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandNicProfileOperatingMode("ACTIVE")
		_ = got
	})

}

func TestUnitExpandNicProfileOwnerType(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandNicProfileOwnerType("")
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandNicProfileOwnerType("ACTIVE")
		_ = got
	})

}

func TestUnitExpandNicProfile(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandNicProfile(map[string]interface{}{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandNicProfile(map[string]interface{}{"name": "UNKNOWN"})
		_ = got
	})

}

func TestUnitExpandStringSet(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandStringSet(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandStringSet("UNKNOWN")
		_ = got
	})

}
