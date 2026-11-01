package networkingv2

import (
	"testing"

	"github.com/stretchr/testify/require"

	import1 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/microseg/v4/config"
)

func TestUnitFlattenUDPPortRangeSpec(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenUDPPortRangeSpec(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenUDPPortRangeSpec([]import1.UdpPortRangeSpec{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenUDPPortRangeSpec([]import1.UdpPortRangeSpec{import1.UdpPortRangeSpec{}})
		_ = got
	})

}

func TestUnitFlattenTCPPortRangeSpec(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenTCPPortRangeSpec(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenTCPPortRangeSpec([]import1.TcpPortRangeSpec{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenTCPPortRangeSpec([]import1.TcpPortRangeSpec{import1.TcpPortRangeSpec{}})
		_ = got
	})

}

func TestUnitFlattenIcmpTypeCodeSpec(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenIcmpTypeCodeSpec(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenIcmpTypeCodeSpec([]import1.IcmpTypeCodeSpec{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenIcmpTypeCodeSpec([]import1.IcmpTypeCodeSpec{import1.IcmpTypeCodeSpec{}})
		_ = got
	})

}
