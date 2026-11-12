package networkingv2

import (
	"testing"

	"github.com/stretchr/testify/require"

	networkingCommon "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/common/v1/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
)

func TestUnitFlattenVpcScope(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = flattenVpcScope(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVpcScope(func() *import1.VpcScope { var v import1.VpcScope; return &v }())
		_ = got
	})

}

func TestUnitFlattenSupportedMultipleExternalSubnetType(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = flattenSupportedMultipleExternalSubnetType(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenSupportedMultipleExternalSubnetType(func() *import1.SupportedMultipleExternalSubnetType {
			var v import1.SupportedMultipleExternalSubnetType
			return &v
		}())
		_ = got
	})

}

func TestUnitFlattenKubernetesClusters(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenKubernetesClusters(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenKubernetesClusters([]import1.KubernetesCluster{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenKubernetesClusters([]import1.KubernetesCluster{import1.KubernetesCluster{}})
		_ = got
	})

}

func TestUnitExpandKubernetesClusters(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandKubernetesClusters(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandKubernetesClusters([]interface{}{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandKubernetesClusters([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandGatewayNodesSelector(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandGatewayNodesSelector(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandGatewayNodesSelector([]interface{}{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandGatewayNodesSelector([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandKubernetesEntitySelector(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandKubernetesEntitySelector(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandKubernetesEntitySelector([]interface{}{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandKubernetesEntitySelector([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandMatchLabels(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandMatchLabels(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandMatchLabels([]interface{}{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandMatchLabels([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandPodNetwork(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandPodNetwork(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandPodNetwork([]interface{}{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandPodNetwork([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitFlattenLabelSelector(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenLabelSelector(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenLabelSelector(&import1.GatewayNodesSelector{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenLabelSelector(&import1.GatewayNodesSelector{})
		_ = got
	})

}

func TestUnitFlattenKubernetesEntitySelector(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenKubernetesEntitySelector(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenKubernetesEntitySelector(&import1.KubernetesEntitySelector{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenKubernetesEntitySelector(&import1.KubernetesEntitySelector{})
		_ = got
	})

}

func TestUnitFlattenMatchLabels(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenMatchLabels(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenMatchLabels([]networkingCommon.KVStringPair{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenMatchLabels([]networkingCommon.KVStringPair{networkingCommon.KVStringPair{}})
		_ = got
	})

}

func TestUnitFlattenPodNetwork(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenPodNetwork(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenPodNetwork(&import1.PodNetwork{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenPodNetwork(&import1.PodNetwork{})
		_ = got
	})

}
