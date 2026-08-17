package tenancyv2

import (
	import6 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/tenancy-go-client/v17/models/common/v1/config"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/tenancy-go-client/v17/models/tenancy/v4/config"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// expandLocalUserInfo builds a LocalUserInfo spec from the admin_info nested block.
func expandLocalUserInfo(input []interface{}) *import1.LocalUserInfo {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	m := input[0].(map[string]interface{})
	adminInfo := &import1.LocalUserInfo{}
	if v, ok := m["email"]; ok && v.(string) != "" {
		adminInfo.Email = utils.StringPtr(v.(string))
	}
	if v, ok := m["first_name"]; ok && v.(string) != "" {
		adminInfo.FirstName = utils.StringPtr(v.(string))
	}
	if v, ok := m["last_name"]; ok && v.(string) != "" {
		adminInfo.LastName = utils.StringPtr(v.(string))
	}
	if v, ok := m["password"]; ok && v.(string) != "" {
		adminInfo.Password = utils.StringPtr(v.(string))
	}
	if v, ok := m["username"]; ok && v.(string) != "" {
		adminInfo.Username = utils.StringPtr(v.(string))
	}
	return adminInfo
}

// flattenLocalUserInfo maps a LocalUserInfo response to a nested block. The password
// is never returned by the API, so it is intentionally not flattened.
func flattenLocalUserInfo(adminInfo *import1.LocalUserInfo) []map[string]interface{} {
	if adminInfo == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"email":      utils.StringValue(adminInfo.Email),
			"first_name": utils.StringValue(adminInfo.FirstName),
			"last_name":  utils.StringValue(adminInfo.LastName),
			"username":   utils.StringValue(adminInfo.Username),
		},
	}
}

// expandQuota builds a Quota spec from a nested block.
func expandQuota(input []interface{}) *import1.Quota {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	m := input[0].(map[string]interface{})
	quota := &import1.Quota{}
	if v, ok := m["ext_id"]; ok && v.(string) != "" {
		quota.ExtId = utils.StringPtr(v.(string))
	}
	if v, ok := m["memory"]; ok {
		quota.Memory = utils.Int64Ptr(int64(v.(int)))
	}
	if v, ok := m["storage"]; ok {
		quota.Storage = utils.Int64Ptr(int64(v.(int)))
	}
	if v, ok := m["vcpu"]; ok {
		quota.Vcpu = utils.IntPtr(v.(int))
	}
	return quota
}

// flattenQuota maps a Quota response to a nested block.
func flattenQuota(quota *import1.Quota) []map[string]interface{} {
	if quota == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"ext_id":  utils.StringValue(quota.ExtId),
			"memory":  utils.Int64Value(quota.Memory),
			"storage": utils.Int64Value(quota.Storage),
			"vcpu":    utils.IntValue(quota.Vcpu),
		},
	}
}

// flattenTenantStateInfo maps a TenantStateInfo response to a nested block.
func flattenTenantStateInfo(stateInfo *import1.TenantStateInfo) []map[string]interface{} {
	if stateInfo == nil {
		return nil
	}
	m := make(map[string]interface{})
	if stateInfo.State != nil {
		m["state"] = stateInfo.State.GetName()
	}
	return []map[string]interface{}{m}
}

// flattenDomainAllocationStateInfo maps a DomainAllocationStateInfo response to a nested block.
func flattenDomainAllocationStateInfo(stateInfo *import1.DomainAllocationStateInfo) []map[string]interface{} {
	if stateInfo == nil {
		return nil
	}
	m := make(map[string]interface{})
	if stateInfo.State != nil {
		m["state"] = stateInfo.State.GetName()
	}
	return []map[string]interface{}{m}
}

// flattenExternalNetworkConfigStateInfo maps an ExternalNetworkConfigStateInfo response to a nested block.
func flattenExternalNetworkConfigStateInfo(stateInfo *import1.ExternalNetworkConfigStateInfo) []map[string]interface{} {
	if stateInfo == nil {
		return nil
	}
	m := make(map[string]interface{})
	if stateInfo.State != nil {
		m["state"] = stateInfo.State.GetName()
	}
	return []map[string]interface{}{m}
}

// expandProjects builds a list of Project specs from a nested block list.
func expandProjects(input []interface{}) []import1.Project {
	if len(input) == 0 {
		return nil
	}
	projects := make([]import1.Project, 0, len(input))
	for _, item := range input {
		if item == nil {
			continue
		}
		m := item.(map[string]interface{})
		project := import1.Project{}
		if v, ok := m["ext_id"]; ok && v.(string) != "" {
			project.ExtId = utils.StringPtr(v.(string))
		}
		if v, ok := m["name"]; ok && v.(string) != "" {
			project.Name = utils.StringPtr(v.(string))
		}
		if v, ok := m["quota"]; ok {
			project.Quota = expandQuota(v.([]interface{}))
		}
		if v, ok := m["resource_group"]; ok {
			project.ResourceGroup = expandResourceGroup(v.([]interface{}))
		}
		projects = append(projects, project)
	}
	return projects
}

// flattenProjects maps a list of Project responses to a nested block list.
func flattenProjects(projects []import1.Project) []map[string]interface{} {
	if len(projects) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(projects))
	for _, project := range projects {
		m := map[string]interface{}{
			"ext_id":         utils.StringValue(project.ExtId),
			"name":           utils.StringValue(project.Name),
			"quota":          flattenQuota(project.Quota),
			"resource_group": flattenResourceGroup(project.ResourceGroup),
		}
		if project.State != nil {
			m["state"] = project.State.GetName()
		}
		result = append(result, m)
	}
	return result
}

// expandResourceGroup builds a ResourceGroup spec from a nested block.
func expandResourceGroup(input []interface{}) *import1.ResourceGroup {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	m := input[0].(map[string]interface{})
	resourceGroup := &import1.ResourceGroup{}
	if v, ok := m["uuid"]; ok && v.(string) != "" {
		resourceGroup.Uuid = utils.StringPtr(v.(string))
	}
	if v, ok := m["placement_targets"]; ok {
		resourceGroup.PlacementTargets = expandPlacementTargets(v.([]interface{}))
	}
	return resourceGroup
}

// flattenResourceGroup maps a ResourceGroup response to a nested block.
func flattenResourceGroup(resourceGroup *import1.ResourceGroup) []map[string]interface{} {
	if resourceGroup == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"uuid":              utils.StringValue(resourceGroup.Uuid),
			"placement_targets": flattenPlacementTargets(resourceGroup.PlacementTargets),
		},
	}
}

// expandPlacementTargets builds a list of PlacementTarget specs from a nested block list.
func expandPlacementTargets(input []interface{}) []import1.PlacementTarget {
	if len(input) == 0 {
		return nil
	}
	targets := make([]import1.PlacementTarget, 0, len(input))
	for _, item := range input {
		if item == nil {
			continue
		}
		m := item.(map[string]interface{})
		target := import1.PlacementTarget{}
		if v, ok := m["cluster_uuid"]; ok && v.(string) != "" {
			target.ClusterUuid = utils.StringPtr(v.(string))
		}
		if v, ok := m["storage_container_uuids"]; ok {
			target.StorageContainerUuids = common.ExpandListOfString(v.([]interface{}))
		}
		targets = append(targets, target)
	}
	return targets
}

// flattenPlacementTargets maps a list of PlacementTarget responses to a nested block list.
func flattenPlacementTargets(targets []import1.PlacementTarget) []map[string]interface{} {
	if len(targets) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(targets))
	for _, target := range targets {
		result = append(result, map[string]interface{}{
			"cluster_uuid":            utils.StringValue(target.ClusterUuid),
			"storage_container_uuids": target.StorageContainerUuids,
		})
	}
	return result
}

// expandExternalNetworkConnection builds an ExternalNetworkConnection spec from a nested block.
func expandExternalNetworkConnection(input []interface{}) *import1.ExternalNetworkConnection {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	m := input[0].(map[string]interface{})
	externalNetworkConfig := &import1.ExternalNetworkConnection{}
	if v, ok := m["config"]; ok {
		externalNetworkConfig.Config = expandExternalNetworkConnConfig(v.([]interface{}))
	}
	return externalNetworkConfig
}

// flattenExternalNetworkConnection maps an ExternalNetworkConnection response to a nested block.
func flattenExternalNetworkConnection(externalNetworkConfig *import1.ExternalNetworkConnection) []map[string]interface{} {
	if externalNetworkConfig == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"config":     flattenExternalNetworkConnConfig(externalNetworkConfig.Config),
			"ext_id":     utils.StringValue(externalNetworkConfig.ExtId),
			"tenant_id":  utils.StringValue(externalNetworkConfig.TenantId),
			"state_info": flattenExternalNetworkConfigStateInfo(externalNetworkConfig.StateInfo),
		},
	}
}

// expandExternalNetworkConnConfig builds an ExternalNetworkConnConfig spec from a nested block.
func expandExternalNetworkConnConfig(input []interface{}) *import1.ExternalNetworkConnConfig {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	m := input[0].(map[string]interface{})
	config := &import1.ExternalNetworkConnConfig{}
	if v, ok := m["nat"]; ok {
		config.Nat = expandNatConfig(v.([]interface{}))
	}
	if v, ok := m["no_nat"]; ok {
		config.NoNat = expandNoNatConfig(v.([]interface{}))
	}
	return config
}

// flattenExternalNetworkConnConfig maps an ExternalNetworkConnConfig response to a nested block.
func flattenExternalNetworkConnConfig(config *import1.ExternalNetworkConnConfig) []map[string]interface{} {
	if config == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"nat":    flattenNatConfig(config.Nat),
			"no_nat": flattenNoNatConfig(config.NoNat),
		},
	}
}

// expandNatConfig builds a NatConfig spec from a nested block.
func expandNatConfig(input []interface{}) *import1.NatConfig {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	m := input[0].(map[string]interface{})
	nat := &import1.NatConfig{}
	if v, ok := m["transit_vpc_ext_id"]; ok && v.(string) != "" {
		nat.TransitVpcExtId = utils.StringPtr(v.(string))
	}
	if v, ok := m["connections"]; ok {
		nat.Connections = expandNatConnectionConfigs(v.([]interface{}))
	}
	return nat
}

// flattenNatConfig maps a NatConfig response to a nested block.
func flattenNatConfig(nat *import1.NatConfig) []map[string]interface{} {
	if nat == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"transit_vpc_ext_id": utils.StringValue(nat.TransitVpcExtId),
			"connections":        flattenNatConnectionConfigs(nat.Connections),
		},
	}
}

// expandNatConnectionConfigs builds a list of NatConnectionConfig specs from a nested block list.
func expandNatConnectionConfigs(input []interface{}) []import1.NatConnectionConfig {
	if len(input) == 0 {
		return nil
	}
	connections := make([]import1.NatConnectionConfig, 0, len(input))
	for _, item := range input {
		if item == nil {
			continue
		}
		m := item.(map[string]interface{})
		connection := import1.NatConnectionConfig{}
		if v, ok := m["floating_ip_count"]; ok {
			connection.FloatingIpCount = utils.IntPtr(v.(int))
		}
		if v, ok := m["overlay_subnet_ext_id"]; ok && v.(string) != "" {
			connection.OverlaySubnetExtId = utils.StringPtr(v.(string))
		}
		connections = append(connections, connection)
	}
	return connections
}

// flattenNatConnectionConfigs maps a list of NatConnectionConfig responses to a nested block list.
func flattenNatConnectionConfigs(connections []import1.NatConnectionConfig) []map[string]interface{} {
	if len(connections) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(connections))
	for _, connection := range connections {
		result = append(result, map[string]interface{}{
			"floating_ip_count":     utils.IntValue(connection.FloatingIpCount),
			"overlay_subnet_ext_id": utils.StringValue(connection.OverlaySubnetExtId),
		})
	}
	return result
}

// expandNoNatConfig builds a NoNatConfig spec from a nested block.
func expandNoNatConfig(input []interface{}) *import1.NoNatConfig {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	m := input[0].(map[string]interface{})
	noNat := &import1.NoNatConfig{}
	if v, ok := m["overlay_subnet_ext_id"]; ok && v.(string) != "" {
		noNat.OverlaySubnetExtId = utils.StringPtr(v.(string))
	}
	if v, ok := m["transit_vpc_ext_id"]; ok && v.(string) != "" {
		noNat.TransitVpcExtId = utils.StringPtr(v.(string))
	}
	if v, ok := m["externally_routable_prefixes"]; ok {
		noNat.ExternallyRoutablePrefixes = expandIPSubnets(v.([]interface{}))
	}
	return noNat
}

// flattenNoNatConfig maps a NoNatConfig response to a nested block.
func flattenNoNatConfig(noNat *import1.NoNatConfig) []map[string]interface{} {
	if noNat == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"overlay_subnet_ext_id":        utils.StringValue(noNat.OverlaySubnetExtId),
			"transit_vpc_ext_id":           utils.StringValue(noNat.TransitVpcExtId),
			"externally_routable_prefixes": flattenIPSubnets(noNat.ExternallyRoutablePrefixes),
		},
	}
}

// expandIPSubnets builds a list of IPSubnet specs from a nested block list.
func expandIPSubnets(input []interface{}) []import1.IPSubnet {
	if len(input) == 0 {
		return nil
	}
	subnets := make([]import1.IPSubnet, 0, len(input))
	for _, item := range input {
		if item == nil {
			continue
		}
		m := item.(map[string]interface{})
		subnet := import1.IPSubnet{}
		if v, ok := m["ipv4"]; ok {
			subnet.Ipv4 = expandIPv4Subnet(v.([]interface{}))
		}
		if v, ok := m["ipv6"]; ok {
			subnet.Ipv6 = expandIPv6Subnet(v.([]interface{}))
		}
		subnets = append(subnets, subnet)
	}
	return subnets
}

// flattenIPSubnets maps a list of IPSubnet responses to a nested block list.
func flattenIPSubnets(subnets []import1.IPSubnet) []map[string]interface{} {
	if len(subnets) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(subnets))
	for _, subnet := range subnets {
		result = append(result, map[string]interface{}{
			"ipv4": flattenIPv4Subnet(subnet.Ipv4),
			"ipv6": flattenIPv6Subnet(subnet.Ipv6),
		})
	}
	return result
}

// expandIPv4Subnet builds an IPv4Subnet spec from a nested block.
func expandIPv4Subnet(input []interface{}) *import1.IPv4Subnet {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	m := input[0].(map[string]interface{})
	subnet := &import1.IPv4Subnet{}
	if v, ok := m["prefix_length"]; ok {
		subnet.PrefixLength = utils.IntPtr(v.(int))
	}
	if v, ok := m["ip"]; ok {
		subnet.Ip = expandIPv4Address(v.([]interface{}))
	}
	return subnet
}

// flattenIPv4Subnet maps an IPv4Subnet response to a nested block.
func flattenIPv4Subnet(subnet *import1.IPv4Subnet) []map[string]interface{} {
	if subnet == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"prefix_length": utils.IntValue(subnet.PrefixLength),
			"ip":            flattenIPv4Address(subnet.Ip),
		},
	}
}

// expandIPv6Subnet builds an IPv6Subnet spec from a nested block.
func expandIPv6Subnet(input []interface{}) *import1.IPv6Subnet {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	m := input[0].(map[string]interface{})
	subnet := &import1.IPv6Subnet{}
	if v, ok := m["prefix_length"]; ok {
		subnet.PrefixLength = utils.IntPtr(v.(int))
	}
	if v, ok := m["ip"]; ok {
		subnet.Ip = expandIPv6Address(v.([]interface{}))
	}
	return subnet
}

// flattenIPv6Subnet maps an IPv6Subnet response to a nested block.
func flattenIPv6Subnet(subnet *import1.IPv6Subnet) []map[string]interface{} {
	if subnet == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"prefix_length": utils.IntValue(subnet.PrefixLength),
			"ip":            flattenIPv6Address(subnet.Ip),
		},
	}
}

// expandIPv4Address builds an IPv4Address spec from a nested block.
func expandIPv4Address(input []interface{}) *import6.IPv4Address {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	m := input[0].(map[string]interface{})
	address := &import6.IPv4Address{}
	if v, ok := m["value"]; ok && v.(string) != "" {
		address.Value = utils.StringPtr(v.(string))
	}
	if v, ok := m["prefix_length"]; ok {
		address.PrefixLength = utils.IntPtr(v.(int))
	}
	return address
}

// flattenIPv4Address maps an IPv4Address response to a nested block.
func flattenIPv4Address(address *import6.IPv4Address) []map[string]interface{} {
	if address == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"value":         utils.StringValue(address.Value),
			"prefix_length": utils.IntValue(address.PrefixLength),
		},
	}
}

// expandIPv6Address builds an IPv6Address spec from a nested block.
func expandIPv6Address(input []interface{}) *import6.IPv6Address {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	m := input[0].(map[string]interface{})
	address := &import6.IPv6Address{}
	if v, ok := m["value"]; ok && v.(string) != "" {
		address.Value = utils.StringPtr(v.(string))
	}
	if v, ok := m["prefix_length"]; ok {
		address.PrefixLength = utils.IntPtr(v.(int))
	}
	return address
}

// flattenIPv6Address maps an IPv6Address response to a nested block.
func flattenIPv6Address(address *import6.IPv6Address) []map[string]interface{} {
	if address == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"value":         utils.StringValue(address.Value),
			"prefix_length": utils.IntValue(address.PrefixLength),
		},
	}
}
