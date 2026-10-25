package networkingv2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitExpandNetworkSecurityPolicyRule(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = expandNetworkSecurityPolicyRule(nil, nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandNetworkSecurityPolicyRule(nil, []interface{}{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandNetworkSecurityPolicyRule(nil, []interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandOneOfNetworkSecurityPolicyRuleSpec(t *testing.T) {
	t.Parallel()

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandOneOfNetworkSecurityPolicyRuleSpec(nil, "rules.0.spec", []interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandOneOfMultiEnvIsolationRuleSpecSpec(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = expandOneOfMultiEnvIsolationRuleSpecSpec(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandOneOfMultiEnvIsolationRuleSpecSpec([]interface{}{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandOneOfMultiEnvIsolationRuleSpecSpec([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandIsolationGroup(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = expandIsolationGroup(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandIsolationGroup([]interface{}{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIsolationGroup([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandIPv4AddressMicroseg(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandIPv4AddressMicroseg(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIPv4AddressMicroseg([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}
