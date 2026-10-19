package clustersv2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	config "github.com/nutanix-core/ntnx-api-golang-sdk-internal/clustermgmt-go-client/v17/models/clustermgmt/v4/config"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/clustermgmt-go-client/v17/models/clustermgmt/v4/config"
	import4 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/clustermgmt-go-client/v17/models/common/v1/config"
	commonResp "github.com/nutanix-core/ntnx-api-golang-sdk-internal/clustermgmt-go-client/v17/models/common/v1/response"
	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitClustersV2Schemas(t *testing.T) {
	t.Parallel()

	require.NoError(t, DatasourceNutanixClusterEntitiesV2().InternalValidate(DatasourceNutanixClusterEntitiesV2().Schema, false))
	require.NoError(t, DatasourceNutanixClusterEntityV2().InternalValidate(DatasourceNutanixClusterEntityV2().Schema, false))
	require.NoError(t, DatasourceNutanixClusterProfileV2().InternalValidate(DatasourceNutanixClusterProfileV2().Schema, false))
	require.NoError(t, DatasourceNutanixClusterProfilesV2().InternalValidate(DatasourceNutanixClusterProfilesV2().Schema, false))
	require.NoError(t, DatasourceNutanixHostEntitiesV2().InternalValidate(DatasourceNutanixHostEntitiesV2().Schema, false))
	require.NoError(t, DatasourceNutanixHostEntityV2().InternalValidate(DatasourceNutanixHostEntityV2().Schema, false))
	require.NoError(t, DatasourceNutanixSnmpConfigV2().InternalValidate(DatasourceNutanixSnmpConfigV2().Schema, false))
	require.NoError(t, DatasourceNutanixSnmpTrapV2().InternalValidate(DatasourceNutanixSnmpTrapV2().Schema, false))
	require.NoError(t, DatasourceNutanixSnmpUserV2().InternalValidate(DatasourceNutanixSnmpUserV2().Schema, false))
	require.NoError(t, DatasourceNutanixSSLCertificateV2().InternalValidate(DatasourceNutanixSSLCertificateV2().Schema, false))
	require.NoError(t, ResourceNutanixClusterCategoriesV2().InternalValidate(ResourceNutanixClusterCategoriesV2().Schema, true))
	require.NoError(t, ResourceNutanixClusterDiscoverUnconfiguredNodesV2().InternalValidate(ResourceNutanixClusterDiscoverUnconfiguredNodesV2().Schema, true))
	require.NoError(t, ResourceNutanixClusterAddNodeV2().InternalValidate(ResourceNutanixClusterAddNodeV2().Schema, true))
	require.NoError(t, ResourceNutanixClusterV2().InternalValidate(ResourceNutanixClusterV2().Schema, true))
	require.NoError(t, ResourceNutanixClusterUnconfiguredNodeNetworkV2().InternalValidate(ResourceNutanixClusterUnconfiguredNodeNetworkV2().Schema, true))
	require.NoError(t, ResourceNutanixClusterProfileV2().InternalValidate(ResourceNutanixClusterProfileV2().Schema, true))
	require.NoError(t, ResourceNutanixSnmpConfigV2().InternalValidate(ResourceNutanixSnmpConfigV2().Schema, true))
	require.NoError(t, ResourceNutanixSnmpTrapV2().InternalValidate(ResourceNutanixSnmpTrapV2().Schema, true))
	require.NoError(t, ResourceNutanixSnmpUserV2().InternalValidate(ResourceNutanixSnmpUserV2().Schema, true))
	require.NoError(t, ResourceNutanixSSLCertificateV2().InternalValidate(ResourceNutanixSSLCertificateV2().Schema, true))
	require.NoError(t, schemaForLinks().Elem.(*schema.Resource).InternalValidate(schemaForLinks().Elem.(*schema.Resource).Schema, true))
	require.NoError(t, schemaForValuePrefixLength().Elem.(*schema.Resource).InternalValidate(schemaForValuePrefixLength().Elem.(*schema.Resource).Schema, true))
	require.NoError(t, schemaForIPAddress(false).Elem.(*schema.Resource).InternalValidate(schemaForIPAddress(false).Elem.(*schema.Resource).Schema, true))
}

func TestUnitIPHelpers(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, "", ipToKey(nil))
		require.Equal(t, "UNKNOWN", getIPType(nil))
		require.True(t, ipv4Equal(nil, nil))
		require.True(t, ipv6Equal(nil, nil))
		require.True(t, ipAddressEqual(nil, nil))
		require.True(t, nodeEqual(nil, nil))
		require.Empty(t, nodeKeyCandidates(nil))
	})
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, "", ipToKey(&import4.IPAddress{}))
		require.Equal(t, "UNKNOWN", getIPType(&import4.IPAddress{}))
		require.False(t, ipv4Equal(&import4.IPv4Address{}, nil))
		require.Empty(t, nodeKeyCandidates(&config.NodeListItemReference{}))
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		ip := &import4.IPAddress{Ipv4: &import4.IPv4Address{Value: utils.StringPtr("10.0.0.1"), PrefixLength: utils.IntPtr(32)}}
		require.Equal(t, "ipv4:10.0.0.1/32", ipToKey(ip))
		require.Equal(t, "IPV4", getIPType(ip))
		require.True(t, ipv4Equal(ip.Ipv4, ip.Ipv4))
		n := &config.NodeListItemReference{NodeUuid: utils.StringPtr("n1"), ControllerVmIp: ip}
		require.Contains(t, nodeKeyCandidates(n), "uuid:n1")
		require.True(t, nodeEqual(n, n))
	})
}

func TestUnitFindCategories(t *testing.T) {
	t.Parallel()
	t.Run("nil/empty old", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []string{}, FindRemovedCategories(nil, []string{"a"}))
		require.Equal(t, []string{}, FindAddedCategories([]string{"a"}, nil))
	})
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []string{}, FindRemovedCategories([]string{}, []string{}))
		require.Equal(t, []string{}, FindAddedCategories([]string{}, []string{}))
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []string{"a"}, FindRemovedCategories([]string{"a", "b"}, []string{"b"}))
		require.Equal(t, []string{"c"}, FindAddedCategories([]string{"a"}, []string{"a", "c"}))
	})
}

func TestUnitFlattenLinksHelper(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenLinks(nil))
	})
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenLinks([]commonResp.ApiLink{}))
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenLinks([]commonResp.ApiLink{{
			Href: utils.StringPtr("https://example.test"),
			Rel:  utils.StringPtr("self"),
		}})
		require.Equal(t, []map[string]interface{}{{"href": "https://example.test", "rel": "self"}}, got)
	})
}

func TestUnitAuthorizedPublicKeyHash(t *testing.T) {
	t.Parallel()
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, schema.HashString("-"), authorizedPublicKeyHash(map[string]interface{}{}))
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, schema.HashString("n-k"), authorizedPublicKeyHash(map[string]interface{}{"name": "n", "key": "k"}))
	})
}

func TestUnitClustersV2FlattenExpandHelpers(t *testing.T) {
	t.Parallel()

	t.Run("flattenClusterEntities", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenClusterEntities([]import1.Cluster{{}})
		})
	})

	t.Run("flattenNodeReference", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenNodeReference(&import1.NodeReference{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenNodeReference(&import1.NodeReference{})
		})
	})

	t.Run("flattenNodeListItemReference", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenNodeListItemReference([]import1.NodeListItemReference{{}})
		})
	})

	t.Run("flattenClusterNetworkReference", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenClusterNetworkReference(&import1.ClusterNetworkReference{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenClusterNetworkReference(&import1.ClusterNetworkReference{})
		})
	})

	t.Run("flattenHTTPProxyWhiteList", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenHTTPProxyWhiteList([]import1.HttpProxyWhiteListConfig{{}})
		})
	})

	t.Run("flattenTargetType", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			var v import1.HttpProxyWhiteListTargetType
			_ = flattenTargetType(&v)
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			var v import1.HttpProxyWhiteListTargetType
			_ = flattenTargetType(&v)
		})
	})

	t.Run("flattenHTTPProxyList", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenHTTPProxyList([]import1.HttpProxyConfig{{}})
		})
	})

	t.Run("flattenProxyTypes", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			var v import1.HttpProxyType
			_ = flattenProxyTypes([]import1.HttpProxyType{v})
		})
	})

	t.Run("flattenClusterConfigReference", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenClusterConfigReference(&import1.ClusterConfigReference{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenClusterConfigReference(&import1.ClusterConfigReference{})
		})
	})

	t.Run("flattenPulseStatus", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenPulseStatus(&import1.PulseStatus{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenPulseStatus(&import1.PulseStatus{})
		})
	})

	t.Run("flattenPulseScrubbingLevel", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			var v import1.PIIScrubbingLevel
			_ = flattenPulseScrubbingLevel(&v)
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			var v import1.PIIScrubbingLevel
			_ = flattenPulseScrubbingLevel(&v)
		})
	})

	t.Run("flattenIPAddress", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenIPAddress(&import4.IPAddress{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenIPAddress(&import4.IPAddress{})
		})
	})

	t.Run("flattenIPAddressOrFQDN", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenIPAddressOrFQDN([]import4.IPAddressOrFQDN{{}})
		})
	})

	t.Run("flattenSMTPServerRef", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenSMTPServerRef(&import1.SmtpServerRef{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenSMTPServerRef(&import1.SmtpServerRef{})
		})
	})

	t.Run("flattenBackplaneNetworkParams", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenBackplaneNetworkParams(&import1.BackplaneNetworkParams{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenBackplaneNetworkParams(&import1.BackplaneNetworkParams{})
		})
	})

	t.Run("flattenManagementServerRef", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenManagementServerRef(&import1.ManagementServerRef{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenManagementServerRef(&import1.ManagementServerRef{})
		})
	})

	t.Run("flattenIPv4Address", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenIPv4Address(&import4.IPv4Address{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenIPv4Address(&import4.IPv4Address{})
		})
	})

	t.Run("flattenIPv6Address", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenIPv6Address(&import4.IPv6Address{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenIPv6Address(&import4.IPv6Address{})
		})
	})

	t.Run("flattenFQDN", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenFQDN(&import4.FQDN{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenFQDN(&import4.FQDN{})
		})
	})

	t.Run("flattenBuildReference", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenBuildReference(&import1.BuildReference{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenBuildReference(&import1.BuildReference{})
		})
	})

	t.Run("flattenPublicKey", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenPublicKey([]import1.PublicKey{{}})
		})
	})

	t.Run("flattenSoftwareMapReference", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenSoftwareMapReference([]import1.SoftwareMapReference{{}})
		})
	})

	t.Run("flattenFaultToleranceState", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenFaultToleranceState(&import1.FaultToleranceState{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenFaultToleranceState(&import1.FaultToleranceState{})
		})
	})

	t.Run("flattenRedundancyStatus", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenRedundancyStatus(&import1.RedundancyStatusDetails{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenRedundancyStatus(&import1.RedundancyStatusDetails{})
		})
	})

	t.Run("flattenClusterFaultTolerance", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			var v import1.ClusterFaultToleranceRef
			_ = flattenClusterFaultTolerance(&v)
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			var v import1.ClusterFaultToleranceRef
			_ = flattenClusterFaultTolerance(&v)
		})
	})

	t.Run("flattenHypervisorType", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			var v import1.HypervisorType
			_ = flattenHypervisorType([]import1.HypervisorType{v})
		})
	})

	t.Run("flattenClusterFunctionRef", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			var v import1.ClusterFunctionRef
			_ = flattenClusterFunctionRef([]import1.ClusterFunctionRef{v})
		})
	})

	t.Run("flattenClusterArchReference", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			var v import1.ClusterArchReference
			_ = flattenClusterArchReference(&v)
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			var v import1.ClusterArchReference
			_ = flattenClusterArchReference(&v)
		})
	})

	t.Run("flattenOperationMode", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			var v import1.OperationMode
			_ = flattenOperationMode(&v)
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			var v import1.OperationMode
			_ = flattenOperationMode(&v)
		})
	})

	t.Run("flattenEncryptionStatus", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			var v import1.EncryptionStatus
			_ = flattenEncryptionStatus(&v)
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			var v import1.EncryptionStatus
			_ = flattenEncryptionStatus(&v)
		})
	})

	t.Run("flattenEncryptionOptionInfo", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			var v import1.EncryptionOptionInfo
			_ = flattenEncryptionOptionInfo([]import1.EncryptionOptionInfo{v})
		})
	})

	t.Run("flattenEncryptionScopeInfo", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			var v import1.EncryptionScopeInfo
			_ = flattenEncryptionScopeInfo([]import1.EncryptionScopeInfo{v})
		})
	})

	t.Run("flattenKeyManagementServerType", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			var v import1.KeyManagementServerType
			_ = flattenKeyManagementServerType(&v)
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			var v import1.KeyManagementServerType
			_ = flattenKeyManagementServerType(&v)
		})
	})

	t.Run("flattenDomainAwarenessLevel", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			var v import1.DomainAwarenessLevel
			_ = flattenDomainAwarenessLevel(&v)
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			var v import1.DomainAwarenessLevel
			_ = flattenDomainAwarenessLevel(&v)
		})
	})

	t.Run("flattenUpgradeStatus", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			var v import1.UpgradeStatus
			_ = flattenUpgradeStatus(&v)
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			var v import1.UpgradeStatus
			_ = flattenUpgradeStatus(&v)
		})
	})

	t.Run("flattenSoftwareTypeRef", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			var v import1.SoftwareTypeRef
			_ = flattenSoftwareTypeRef(&v)
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			var v import1.SoftwareTypeRef
			_ = flattenSoftwareTypeRef(&v)
		})
	})

	t.Run("flattenClusterProfileClusterList", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenClusterProfileClusterList([]config.ManagedCluster{{}})
		})
	})

	t.Run("flattenClusterProfiles", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenClusterProfiles([]import1.ClusterProfile{{}})
		})
	})

	t.Run("flattenHostEntities", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenHostEntities([]import1.Host{{}})
		})
	})

	t.Run("flattenNodeStatus", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			var v import1.NodeStatus
			_ = flattenNodeStatus(&v)
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			var v import1.NodeStatus
			_ = flattenNodeStatus(&v)
		})
	})

	t.Run("flattenHostTypeEnum", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			var v import1.HostTypeEnum
			_ = flattenHostTypeEnum(&v)
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			var v import1.HostTypeEnum
			_ = flattenHostTypeEnum(&v)
		})
	})

	t.Run("flattenHypervisorReference", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenHypervisorReference(&import1.HypervisorReference{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenHypervisorReference(&import1.HypervisorReference{})
		})
	})

	t.Run("flattenHypervisorState", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			var v import1.HypervisorState
			_ = flattenHypervisorState(&v)
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			var v import1.HypervisorState
			_ = flattenHypervisorState(&v)
		})
	})

	t.Run("flattenHostHypervisorType", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			var v import1.HypervisorType
			_ = flattenHostHypervisorType(&v)
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			var v import1.HypervisorType
			_ = flattenHostHypervisorType(&v)
		})
	})

	t.Run("flattenAcropolisConnectionState", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			var v import1.AcropolisConnectionState
			_ = flattenAcropolisConnectionState(&v)
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			var v import1.AcropolisConnectionState
			_ = flattenAcropolisConnectionState(&v)
		})
	})

	t.Run("flattenClusterReference", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenClusterReference(&import1.ClusterReference{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenClusterReference(&import1.ClusterReference{})
		})
	})

	t.Run("flattenControllerVMReference", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenControllerVMReference(&import1.ControllerVmReference{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenControllerVMReference(&import1.ControllerVmReference{})
		})
	})

	t.Run("flattenIPAddressList", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenIPAddressList([]import4.IPAddress{{}})
		})
	})

	t.Run("flattenIpmiReference", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenIpmiReference(&import1.IpmiReference{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenIpmiReference(&import1.IpmiReference{})
		})
	})

	t.Run("flattenDiskReference", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenDiskReference([]import1.DiskReference{{}})
		})
	})

	t.Run("flattenKeyManagementDeviceToCertStatusInfo", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenKeyManagementDeviceToCertStatusInfo([]import1.KeyManagementDeviceToCertStatusInfo{{}})
		})
	})

	t.Run("flattenStorageTierReference", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			var v import1.StorageTierReference
			_ = flattenStorageTierReference(&v)
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			var v import1.StorageTierReference
			_ = flattenStorageTierReference(&v)
		})
	})

	t.Run("flattenSnmpTransports", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenSnmpTransports([]config.SnmpTransport{{}})
		})
	})

	t.Run("flattenSnmpTraps", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenSnmpTraps([]config.SnmpTrap{{}})
		})
	})

	t.Run("flattenSnmpUsers", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenSnmpUsers([]config.SnmpUser{{}})
		})
	})

	t.Run("flattenUnconfiguredNodes", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenUnconfiguredNodes([]config.UnconfiguredNodeListItem{{}})
		})
	})

	t.Run("flattenAttributes", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenAttributes(&config.UnconfiguredNodeAttributeMap{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenAttributes(&config.UnconfiguredNodeAttributeMap{})
		})
	})

	t.Run("expandIPFilterList", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandIPFilterList([]interface{}{map[string]interface{}{"name": "n", "value": "v", "email_address": "a@b.c", "type": "PLAIN", "is_enabled": true}})
		})
	})

	t.Run("expandClusterNodeParams", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			d := schema.TestResourceDataRaw(t, ResourceNutanixClusterAddNodeV2().Schema, map[string]interface{}{})
			require.Nil(t, expandClusterNodeParams(d, nil))
		})
	})

	t.Run("expandBlockList", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandBlockList(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandBlockList([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandBlockList([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandNodeList", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandNodeList(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandNodeList([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandNodeList([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandComputeNodeList", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandComputeNodeList(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandComputeNodeList([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandComputeNodeList([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandKeyValueMap", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandKeyValueMap(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandKeyValueMap([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandKeyValueMap([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandNetworks", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandNetworks(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandNetworks([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandNetworks([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandUplink", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandUplink(nil))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandUplink([]interface{}{map[string]interface{}{"name": "n", "value": "v", "email_address": "a@b.c", "type": "PLAIN", "is_enabled": true}})
		})
	})

	t.Run("expandUplinkParams", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandUplinkParams(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandUplinkParams([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandUplinkParams([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandHypervisorIsos", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandHypervisorIsos([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandBundleInfo", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandBundleInfo(nil))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandBundleInfo([]interface{}{map[string]interface{}{"name": "n", "value": "v", "email_address": "a@b.c", "type": "PLAIN", "is_enabled": true}})
		})
	})

	t.Run("expandClusterConfigParams", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandClusterConfigParams(nil))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandClusterConfigParams([]interface{}{map[string]interface{}{"name": "n", "value": "v", "email_address": "a@b.c", "type": "PLAIN", "is_enabled": true}})
		})
	})

	t.Run("expandHyperv", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandHyperv(nil))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandHyperv([]interface{}{map[string]interface{}{"name": "n", "value": "v", "email_address": "a@b.c", "type": "PLAIN", "is_enabled": true}})
		})
	})

	t.Run("expandDetails", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandDetails(nil))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandDetails([]interface{}{map[string]interface{}{"name": "n", "value": "v", "email_address": "a@b.c", "type": "PLAIN", "is_enabled": true}})
		})
	})

	t.Run("expandExtraParams", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandExtraParams(nil))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandExtraParams([]interface{}{map[string]interface{}{"name": "n", "value": "v", "email_address": "a@b.c", "type": "PLAIN", "is_enabled": true}})
		})
	})

	t.Run("expandNodeReference", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandNodeReference(nil))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandNodeReference([]interface{}{map[string]interface{}{"name": "n", "value": "v", "email_address": "a@b.c", "type": "PLAIN", "is_enabled": true}})
		})
	})

	t.Run("expandNodeListItemReference", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandNodeListItemReference(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandNodeListItemReference([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandNodeListItemReference([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandClusterNetworkReference", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandClusterNetworkReference(nil))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandClusterNetworkReference([]interface{}{map[string]interface{}{"name": "n", "value": "v", "email_address": "a@b.c", "type": "PLAIN", "is_enabled": true}})
		})
	})

	t.Run("expandHTTPProxyWhiteList", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandHTTPProxyWhiteList(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandHTTPProxyWhiteList([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandHTTPProxyWhiteList([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandHTTPProxyList", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandHTTPProxyList(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandHTTPProxyList([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandHTTPProxyList([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandClusterConfigReference", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			d := schema.TestResourceDataRaw(t, ResourceNutanixClusterV2().Schema, map[string]interface{}{})
			require.Nil(t, expandClusterConfigReference(nil, d))
		})
	})

	t.Run("expandPulseStatus", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandPulseStatus(nil))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandPulseStatus([]interface{}{map[string]interface{}{"name": "n", "value": "v", "email_address": "a@b.c", "type": "PLAIN", "is_enabled": true}})
		})
	})

	t.Run("expandIPAddress", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandIPAddress(nil))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandIPAddress([]interface{}{map[string]interface{}{"name": "n", "value": "v", "email_address": "a@b.c", "type": "PLAIN", "is_enabled": true}})
		})
	})

	t.Run("expandIPAddressOrFQDN", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandIPAddressOrFQDN(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandIPAddressOrFQDN([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandIPAddressOrFQDN([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandIPv4Address", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandIPv4Address([]interface{}{map[string]interface{}{"name": "n", "value": "v", "email_address": "a@b.c", "type": "PLAIN", "is_enabled": true}})
		})
	})

	t.Run("expandIPv6Address", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandIPv6Address([]interface{}{map[string]interface{}{"name": "n", "value": "v", "email_address": "a@b.c", "type": "PLAIN", "is_enabled": true}})
		})
	})

	t.Run("expandSMTPServerRef", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandSMTPServerRef([]interface{}{map[string]interface{}{"name": "n", "value": "v", "email_address": "a@b.c", "type": "PLAIN", "is_enabled": true}})
		})
	})

	t.Run("expandBackplaneNetworkParams", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandBackplaneNetworkParams([]interface{}{map[string]interface{}{"name": "n", "value": "v", "email_address": "a@b.c", "type": "PLAIN", "is_enabled": true}})
		})
	})

	t.Run("expandManagementServerRef", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandManagementServerRef(nil))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandManagementServerRef([]interface{}{map[string]interface{}{"name": "n", "value": "v", "email_address": "a@b.c", "type": "PLAIN", "is_enabled": true}})
		})
	})

	t.Run("expandSMTPNetwork", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandSMTPNetwork(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandSMTPNetwork([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandSMTPNetwork([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandFQDN", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandFQDN([]interface{}{map[string]interface{}{"name": "n", "value": "v", "email_address": "a@b.c", "type": "PLAIN", "is_enabled": true}})
		})
	})

	t.Run("expandBuildReference", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandBuildReference(nil))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandBuildReference([]interface{}{map[string]interface{}{"name": "n", "value": "v", "email_address": "a@b.c", "type": "PLAIN", "is_enabled": true}})
		})
	})

	t.Run("expandPublicKey", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandPublicKey(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandPublicKey([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandPublicKey([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandFaultToleranceState", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandFaultToleranceState(nil))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandFaultToleranceState([]interface{}{map[string]interface{}{"name": "n", "value": "v", "email_address": "a@b.c", "type": "PLAIN", "is_enabled": true}})
		})
	})

	t.Run("flattenNodesNetworkDetails", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenNodesNetworkDetails(config.NodeNetworkingDetails{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenNodesNetworkDetails(config.NodeNetworkingDetails{})
		})
	})

	t.Run("flattenUplinkList", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenUplinkList([]config.NameMacRef{{}})
		})
	})

	t.Run("expandNodeListNetworkingDetails", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandNodeListNetworkingDetails(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandNodeListNetworkingDetails([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandNodeListNetworkingDetails([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandSNMPConfig", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandSNMPConfig(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandSNMPConfig([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandSNMPConfig([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandRsyslogServerList", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandRsyslogServerList(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandRsyslogServerList([]interface{}{}))
		})
	})

	t.Run("flattenSnmpConfig", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenSnmpConfig(&config.SnmpConfig{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenSnmpConfig(&config.SnmpConfig{})
		})
	})

	t.Run("flattenRsyslogServerList", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenRsyslogServerList([]config.RsyslogServer{{}})
		})
	})

	t.Run("expandIPAddressOrFQDNSingle", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandIPAddressOrFQDNSingle([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandNtpServerConfigList", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandNtpServerConfigList(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandNtpServerConfigList([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandNtpServerConfigList([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandHttpProxyConfig", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandHttpProxyConfig(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandHttpProxyConfig([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandHttpProxyConfig([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandFaultToleranceConfig", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandFaultToleranceConfig(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandFaultToleranceConfig([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandFaultToleranceConfig([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandRebuildReservationConfig", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandRebuildReservationConfig(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandRebuildReservationConfig([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandRebuildReservationConfig([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("expandResilientCapacityWarningThresholdConfig", func(t *testing.T) {
		t.Parallel()
		t.Run("nil", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandResilientCapacityWarningThresholdConfig(nil))
		})
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			require.Nil(t, expandResilientCapacityWarningThresholdConfig([]interface{}{}))
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = expandResilientCapacityWarningThresholdConfig([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})
		})
	})

	t.Run("flattenNtpServerConfigList", func(t *testing.T) {
		t.Parallel()
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenNtpServerConfigList([]config.NtpServerConfig{{}})
		})
	})

	t.Run("flattenHttpProxyConfig", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenHttpProxyConfig(&config.HttpProxyConfiguration{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenHttpProxyConfig(&config.HttpProxyConfiguration{})
		})
	})

	t.Run("flattenFaultToleranceConfig", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenFaultToleranceConfig(&config.FaultToleranceConfiguration{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenFaultToleranceConfig(&config.FaultToleranceConfiguration{})
		})
	})

	t.Run("flattenRebuildReservationConfig", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenRebuildReservationConfig(&config.RebuildReservationConfig{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenRebuildReservationConfig(&config.RebuildReservationConfig{})
		})
	})

	t.Run("flattenResilientCapacityWarningThresholdConfig", func(t *testing.T) {
		t.Parallel()
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			_ = flattenResilientCapacityWarningThresholdConfig(&config.ResilientCapacityWarningThresholdConfig{})
		})
		t.Run("populated", func(t *testing.T) {
			t.Parallel()
			_ = flattenResilientCapacityWarningThresholdConfig(&config.ResilientCapacityWarningThresholdConfig{})
		})
	})

}
