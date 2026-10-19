package networkingv2

import (
	"testing"

	"github.com/stretchr/testify/require"

	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v17/models/networking/v4/config"
)

func TestUnitFlattenPolicies(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenPolicies(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenPolicies([]import1.RoutingPolicyRule{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenPolicies([]import1.RoutingPolicyRule{import1.RoutingPolicyRule{}})
		_ = got
	})

}

func TestUnitFlattenPolicyMatch(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenPolicyMatch(&import1.RoutingPolicyMatchCondition{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenPolicyMatch(&import1.RoutingPolicyMatchCondition{})
		_ = got
	})

}

func TestUnitFlattenAddressTypeObject(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenAddressTypeObject(&import1.AddressTypeObject{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenAddressTypeObject(&import1.AddressTypeObject{})
		_ = got
	})

}

func TestUnitFlattenIPSubnet(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenIPSubnet(&import1.IPSubnet{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenIPSubnet(&import1.IPSubnet{})
		_ = got
	})

}

func TestUnitFlattenOneOfRoutingPolicyMatchConditionProtocolParameters(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenOneOfRoutingPolicyMatchConditionProtocolParameters(nil))
	})
}

func TestUnitFlattenRoutingPolicyAction(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenRoutingPolicyAction(&import1.RoutingPolicyAction{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenRoutingPolicyAction(&import1.RoutingPolicyAction{})
		_ = got
	})

}

func TestUnitFlattenRerouteParam(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenRerouteParam(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenRerouteParam([]import1.RerouteParam{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenRerouteParam([]import1.RerouteParam{import1.RerouteParam{}})
		_ = got
	})

}

func TestUnitFlattenPortRange(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenPortRange(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenPortRange([]import1.PortRange{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenPortRange([]import1.PortRange{import1.PortRange{}})
		_ = got
	})

}

func TestUnitFlattenAddressType(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenAddressType(func() *import1.AddressType { var v import1.AddressType; return &v }())
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenAddressType(func() *import1.AddressType { v := import1.AddressType(1); return &v }())
		_ = got
	})

}

func TestUnitFlattenRoutingPolicyActionType(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenRoutingPolicyActionType(func() *import1.RoutingPolicyActionType { var v import1.RoutingPolicyActionType; return &v }())
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenRoutingPolicyActionType(func() *import1.RoutingPolicyActionType { v := import1.RoutingPolicyActionType(1); return &v }())
		_ = got
	})

}

func TestUnitFlattenProtocolType(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenProtocolType(func() *import1.ProtocolType { var v import1.ProtocolType; return &v }())
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenProtocolType(func() *import1.ProtocolType { v := import1.ProtocolType(1); return &v }())
		_ = got
	})

}

func TestUnitFlattenVpcName(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVpcName(&import1.VpcName{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenVpcName(&import1.VpcName{})
		_ = got
	})

}

func TestUnitFlattenRerouteFallbackAction(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenRerouteFallbackAction(func() *import1.RerouteFallbackAction { var v import1.RerouteFallbackAction; return &v }())
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenRerouteFallbackAction(func() *import1.RerouteFallbackAction { v := import1.RerouteFallbackAction(1); return &v }())
		_ = got
	})

}
