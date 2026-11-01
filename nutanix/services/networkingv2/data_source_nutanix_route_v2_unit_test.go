package networkingv2

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
)

func TestUnitFlattenDestination(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenDestination(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenDestination(&config.IPSubnet{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenDestination(&config.IPSubnet{})
		_ = got
	})

}

func TestUnitFlattenNextHops(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenNextHops(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenNextHops([]config.Nexthop{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenNextHops([]config.Nexthop{config.Nexthop{}})
		_ = got
	})

}

func TestUnitFlattenNextHopType(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenNextHopType(func() *config.NexthopType { var v config.NexthopType; return &v }())
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenNextHopType(func() *config.NexthopType { v := config.NexthopType(1); return &v }())
		_ = got
	})

}

func TestUnitFlattenRouteType(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenRouteType(func() *config.RouteType { var v config.RouteType; return &v }())
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenRouteType(func() *config.RouteType { v := config.RouteType(1); return &v }())
		_ = got
	})

}
