package networkingv2

import (
	"testing"

	"github.com/stretchr/testify/require"

	config "github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/models/common/v1/config"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/models/microseg/v4/config"
)

func TestUnitFlattenNetworkSecurityPolicyRule(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = flattenNetworkSecurityPolicyRule(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenNetworkSecurityPolicyRule([]import1.NetworkSecurityPolicyRule{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenNetworkSecurityPolicyRule([]import1.NetworkSecurityPolicyRule{import1.NetworkSecurityPolicyRule{}})
		_ = got
	})

}

func TestUnitFlattenOneOfNetworkSecurityPolicyRuleSpec(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenOneOfNetworkSecurityPolicyRuleSpec(nil))
	})
}

func TestUnitFlattenIPv4AddressMicroSegList(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenIPv4AddressMicroSegList(&config.IPv4Address{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenIPv4AddressMicroSegList(&config.IPv4Address{})
		_ = got
	})

}
