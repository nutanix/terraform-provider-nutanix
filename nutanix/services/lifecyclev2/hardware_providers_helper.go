package lifecyclev2

import (
	commonConfig "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/common/v1/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/config"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// -----------------------------------------------------------------------------
// Expand helpers (schema -> SDK) for the hardware-provider Connection resource.
// -----------------------------------------------------------------------------

// expandConnectionAccessDetails converts the schema representation of a connection's
// access details into the SDK []ConnectionDetail slice.
func expandConnectionAccessDetails(pr []interface{}) []import1.ConnectionDetail {
	if len(pr) == 0 {
		return nil
	}
	details := make([]import1.ConnectionDetail, len(pr))
	for i, item := range pr {
		val := item.(map[string]interface{})
		detail := import1.ConnectionDetail{}
		if auth, ok := val["auth"]; ok {
			detail.Auth = expandConnectionDetailAuth(auth.([]interface{}))
		}
		if endpoint, ok := val["endpoint"]; ok {
			detail.Endpoint = expandConnectionDetailEndpoint(endpoint.([]interface{}))
		}
		details[i] = detail
	}
	return details
}

// expandConnectionDetailAuth builds the OneOfConnectionDetailAuth wrapper from the
// nested schema block. Exactly one of api_key_auth / basic_auth should be set.
func expandConnectionDetailAuth(pr []interface{}) *import1.OneOfConnectionDetailAuth {
	if len(pr) == 0 || pr[0] == nil {
		return nil
	}
	val := pr[0].(map[string]interface{})
	auth := import1.NewOneOfConnectionDetailAuth()

	if apiKey, ok := val["api_key_auth"]; ok && len(apiKey.([]interface{})) > 0 {
		apiKeyObj := import1.NewApiKeyAuth()
		akMap := apiKey.([]interface{})[0].(map[string]interface{})
		if v, ok := akMap["api_key_id"]; ok && v.(string) != "" {
			apiKeyObj.ApiKeyId = utils.StringPtr(v.(string))
		}
		if v, ok := akMap["api_key_secret"]; ok && v.(string) != "" {
			apiKeyObj.ApiKeySecret = utils.StringPtr(v.(string))
		}
		if err := auth.SetValue(*apiKeyObj); err == nil {
			return auth
		}
	}

	if basic, ok := val["basic_auth"]; ok && len(basic.([]interface{})) > 0 {
		basicObj := import1.NewBasicAuth()
		bMap := basic.([]interface{})[0].(map[string]interface{})
		if v, ok := bMap["username"]; ok && v.(string) != "" {
			basicObj.Username = utils.StringPtr(v.(string))
		}
		if v, ok := bMap["password"]; ok && v.(string) != "" {
			basicObj.Password = utils.StringPtr(v.(string))
		}
		if err := auth.SetValue(*basicObj); err == nil {
			return auth
		}
	}

	return nil
}

// expandConnectionDetailEndpoint builds the OneOfConnectionDetailEndpoint wrapper.
// Exactly one of url_endpoint / ip_range_endpoint / ip_address_endpoint should be set.
func expandConnectionDetailEndpoint(pr []interface{}) *import1.OneOfConnectionDetailEndpoint {
	if len(pr) == 0 || pr[0] == nil {
		return nil
	}
	val := pr[0].(map[string]interface{})
	endpoint := import1.NewOneOfConnectionDetailEndpoint()

	if urlEp, ok := val["url_endpoint"]; ok && len(urlEp.([]interface{})) > 0 {
		urlObj := import1.NewUrlEndpoint()
		uMap := urlEp.([]interface{})[0].(map[string]interface{})
		if v, ok := uMap["url"]; ok && v.(string) != "" {
			urlObj.Url = utils.StringPtr(v.(string))
		}
		if err := endpoint.SetValue(*urlObj); err == nil {
			return endpoint
		}
	}

	if rangeEp, ok := val["ip_range_endpoint"]; ok && len(rangeEp.([]interface{})) > 0 {
		rangeObj := import1.NewIpRangeEndpoint()
		rMap := rangeEp.([]interface{})[0].(map[string]interface{})
		if v, ok := rMap["ip_ranges"]; ok {
			rangeObj.IpRanges = expandIPRanges(v.([]interface{}))
		}
		if err := endpoint.SetValue(*rangeObj); err == nil {
			return endpoint
		}
	}

	if addrEp, ok := val["ip_address_endpoint"]; ok && len(addrEp.([]interface{})) > 0 {
		addrObj := import1.NewIpAddressEndpoint()
		aMap := addrEp.([]interface{})[0].(map[string]interface{})
		if v, ok := aMap["ip_addresses"]; ok {
			addrObj.IpAddresses = expandIPAddressesOrFQDN(v.([]interface{}))
		}
		if err := endpoint.SetValue(*addrObj); err == nil {
			return endpoint
		}
	}

	return nil
}

// expandIPRanges converts the schema representation of IP ranges to the SDK slice.
func expandIPRanges(pr []interface{}) []import1.IpRange {
	if len(pr) == 0 {
		return nil
	}
	ranges := make([]import1.IpRange, len(pr))
	for i, item := range pr {
		val := item.(map[string]interface{})
		ipRange := import1.IpRange{}
		if v, ok := val["start_ip"]; ok {
			ipRange.StartIp = expandIPAddress(v.([]interface{}))
		}
		if v, ok := val["end_ip"]; ok {
			ipRange.EndIp = expandIPAddress(v.([]interface{}))
		}
		ranges[i] = ipRange
	}
	return ranges
}

// expandIPAddress converts a single ipv4/ipv6 schema block to the SDK IPAddress.
func expandIPAddress(pr []interface{}) *commonConfig.IPAddress {
	if len(pr) == 0 || pr[0] == nil {
		return nil
	}
	val := pr[0].(map[string]interface{})
	ip := commonConfig.NewIPAddress()
	if v, ok := val["ipv4"]; ok && len(v.([]interface{})) > 0 {
		ip.Ipv4 = expandIPv4Address(v.([]interface{}))
	}
	if v, ok := val["ipv6"]; ok && len(v.([]interface{})) > 0 {
		ip.Ipv6 = expandIPv6Address(v.([]interface{}))
	}
	if ip.Ipv4 == nil && ip.Ipv6 == nil {
		return nil
	}
	return ip
}

// expandIPAddressesOrFQDN converts the schema representation of ip_addresses to the SDK slice.
func expandIPAddressesOrFQDN(pr []interface{}) []commonConfig.IPAddressOrFQDN {
	if len(pr) == 0 {
		return nil
	}
	addrs := make([]commonConfig.IPAddressOrFQDN, 0, len(pr))
	for _, item := range pr {
		val := item.(map[string]interface{})
		addr := commonConfig.IPAddressOrFQDN{}
		if v, ok := val["ipv4"]; ok && len(v.([]interface{})) > 0 {
			addr.Ipv4 = expandIPv4Address(v.([]interface{}))
		}
		if v, ok := val["ipv6"]; ok && len(v.([]interface{})) > 0 {
			addr.Ipv6 = expandIPv6Address(v.([]interface{}))
		}
		if v, ok := val["fqdn"]; ok && len(v.([]interface{})) > 0 {
			fMap := v.([]interface{})[0].(map[string]interface{})
			fqdn := commonConfig.NewFQDN()
			if fv, ok := fMap["value"]; ok && fv.(string) != "" {
				fqdn.Value = utils.StringPtr(fv.(string))
			}
			addr.Fqdn = fqdn
		}
		addrs = append(addrs, addr)
	}
	return addrs
}

// expandIPv4Address converts a value/prefix_length schema block to the SDK IPv4Address.
func expandIPv4Address(pr []interface{}) *commonConfig.IPv4Address {
	if len(pr) == 0 || pr[0] == nil {
		return nil
	}
	val := pr[0].(map[string]interface{})
	ip := commonConfig.NewIPv4Address()
	if v, ok := val["value"]; ok && v.(string) != "" {
		ip.Value = utils.StringPtr(v.(string))
	}
	if v, ok := val["prefix_length"]; ok {
		ip.PrefixLength = utils.IntPtr(v.(int))
	}
	return ip
}

// expandIPv6Address converts a value/prefix_length schema block to the SDK IPv6Address.
func expandIPv6Address(pr []interface{}) *commonConfig.IPv6Address {
	if len(pr) == 0 || pr[0] == nil {
		return nil
	}
	val := pr[0].(map[string]interface{})
	ip := commonConfig.NewIPv6Address()
	if v, ok := val["value"]; ok && v.(string) != "" {
		ip.Value = utils.StringPtr(v.(string))
	}
	if v, ok := val["prefix_length"]; ok {
		ip.PrefixLength = utils.IntPtr(v.(int))
	}
	return ip
}

// -----------------------------------------------------------------------------
// Flatten helpers (SDK -> schema).
// -----------------------------------------------------------------------------

// flattenConnection flattens a single Connection into a map for d.Set / list elements.
func flattenConnection(conn *import1.Connection) map[string]interface{} {
	if conn == nil {
		return nil
	}
	m := make(map[string]interface{})
	m["tenant_id"] = utils.StringValue(conn.TenantId)
	m["ext_id"] = utils.StringValue(conn.ExtId)
	m["links"] = common.FlattenLinks(conn.Links)
	m["name"] = utils.StringValue(conn.Name)
	m["region"] = utils.StringValue(conn.Region)
	m["access_details"] = flattenConnectionAccessDetails(conn.AccessDetails)
	if conn.DeploymentType != nil {
		m["deployment_type"] = conn.DeploymentType.GetName()
	}
	if conn.CreatedTime != nil {
		m["created_time"] = conn.CreatedTime.String()
	}
	return m
}

// flattenConnectionAccessDetails flattens the []ConnectionDetail slice.
func flattenConnectionAccessDetails(details []import1.ConnectionDetail) []map[string]interface{} {
	if len(details) == 0 {
		return nil
	}
	out := make([]map[string]interface{}, len(details))
	for i, d := range details {
		m := make(map[string]interface{})
		m["auth"] = flattenConnectionDetailAuth(d.Auth)
		m["endpoint"] = flattenConnectionDetailEndpoint(d.Endpoint)
		out[i] = m
	}
	return out
}

// flattenConnectionDetailAuth flattens the OneOfConnectionDetailAuth wrapper.
func flattenConnectionDetailAuth(auth *import1.OneOfConnectionDetailAuth) []map[string]interface{} {
	if auth == nil || auth.ObjectType_ == nil {
		return nil
	}
	m := make(map[string]interface{})
	value := auth.GetValue()
	switch v := value.(type) {
	case import1.ApiKeyAuth:
		m["api_key_auth"] = []map[string]interface{}{
			{
				"api_key_id":     utils.StringValue(v.ApiKeyId),
				"api_key_secret": utils.StringValue(v.ApiKeySecret),
			},
		}
	case import1.BasicAuth:
		m["basic_auth"] = []map[string]interface{}{
			{
				"username": utils.StringValue(v.Username),
				"password": utils.StringValue(v.Password),
			},
		}
	default:
		return nil
	}
	return []map[string]interface{}{m}
}

// flattenConnectionDetailEndpoint flattens the OneOfConnectionDetailEndpoint wrapper.
func flattenConnectionDetailEndpoint(endpoint *import1.OneOfConnectionDetailEndpoint) []map[string]interface{} {
	if endpoint == nil || endpoint.ObjectType_ == nil {
		return nil
	}
	m := make(map[string]interface{})
	value := endpoint.GetValue()
	switch v := value.(type) {
	case import1.UrlEndpoint:
		m["url_endpoint"] = []map[string]interface{}{
			{"url": utils.StringValue(v.Url)},
		}
	case import1.IpRangeEndpoint:
		m["ip_range_endpoint"] = []map[string]interface{}{
			{"ip_ranges": flattenIPRanges(v.IpRanges)},
		}
	case import1.IpAddressEndpoint:
		m["ip_address_endpoint"] = []map[string]interface{}{
			{"ip_addresses": flattenIPAddressesOrFQDN(v.IpAddresses)},
		}
	default:
		return nil
	}
	return []map[string]interface{}{m}
}

// flattenIPRanges flattens []IpRange.
func flattenIPRanges(ranges []import1.IpRange) []map[string]interface{} {
	if len(ranges) == 0 {
		return nil
	}
	out := make([]map[string]interface{}, len(ranges))
	for i, r := range ranges {
		m := make(map[string]interface{})
		m["start_ip"] = flattenIPAddress(r.StartIp)
		m["end_ip"] = flattenIPAddress(r.EndIp)
		out[i] = m
	}
	return out
}

// flattenIPAddress flattens a single IPAddress into an ipv4/ipv6 block.
func flattenIPAddress(ip *commonConfig.IPAddress) []map[string]interface{} {
	if ip == nil {
		return nil
	}
	m := make(map[string]interface{})
	if ip.Ipv4 != nil {
		m["ipv4"] = flattenIPv4Address(ip.Ipv4)
	}
	if ip.Ipv6 != nil {
		m["ipv6"] = flattenIPv6Address(ip.Ipv6)
	}
	return []map[string]interface{}{m}
}

// flattenIPAddressesOrFQDN flattens []IPAddressOrFQDN.
func flattenIPAddressesOrFQDN(addrs []commonConfig.IPAddressOrFQDN) []map[string]interface{} {
	if len(addrs) == 0 {
		return nil
	}
	out := make([]map[string]interface{}, len(addrs))
	for i, a := range addrs {
		m := make(map[string]interface{})
		if a.Ipv4 != nil {
			m["ipv4"] = flattenIPv4Address(a.Ipv4)
		}
		if a.Ipv6 != nil {
			m["ipv6"] = flattenIPv6Address(a.Ipv6)
		}
		if a.Fqdn != nil {
			m["fqdn"] = []map[string]interface{}{
				{"value": utils.StringValue(a.Fqdn.Value)},
			}
		}
		out[i] = m
	}
	return out
}

// flattenIPv4Address flattens an IPv4Address into a value/prefix_length block.
func flattenIPv4Address(ip *commonConfig.IPv4Address) []map[string]interface{} {
	if ip == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"value":         utils.StringValue(ip.Value),
			"prefix_length": utils.IntValue(ip.PrefixLength),
		},
	}
}

// flattenIPv6Address flattens an IPv6Address into a value/prefix_length block.
func flattenIPv6Address(ip *commonConfig.IPv6Address) []map[string]interface{} {
	if ip == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"value":         utils.StringValue(ip.Value),
			"prefix_length": utils.IntValue(ip.PrefixLength),
		},
	}
}

// flattenPoolGroups flattens []PoolGroup used by IP/MAC/server-identity pools.
func flattenPoolGroups(groups []import1.PoolGroup) []map[string]interface{} {
	if len(groups) == 0 {
		return nil
	}
	out := make([]map[string]interface{}, len(groups))
	for i, g := range groups {
		out[i] = map[string]interface{}{
			"id": utils.StringValue(g.Id),
		}
	}
	return out
}

// flattenNodeIdentifiers flattens []NodeIdentifier used by discovered nodes.
func flattenNodeIdentifiers(identifiers []import1.NodeIdentifier) []map[string]interface{} {
	if len(identifiers) == 0 {
		return nil
	}
	out := make([]map[string]interface{}, len(identifiers))
	for i, id := range identifiers {
		m := map[string]interface{}{
			"value": utils.StringValue(id.Value),
		}
		if id.Type != nil {
			m["type"] = id.Type.GetName()
		}
		out[i] = m
	}
	return out
}

// flattenCPUInfo flattens the CpuInfo block used by discovered nodes.
func flattenCPUInfo(cpu *import1.CpuInfo) []map[string]interface{} {
	if cpu == nil {
		return nil
	}
	var capacityGHz float64
	if cpu.CapacityGHz != nil {
		capacityGHz = float64(*cpu.CapacityGHz)
	}
	return []map[string]interface{}{
		{
			"capacity_ghz":       capacityGHz,
			"logical_core_count": utils.IntValue(cpu.LogicalCoreCount),
			"manufacturer":       utils.StringValue(cpu.Manufacturer),
			"model":              utils.StringValue(cpu.Model),
			"socket_count":       utils.IntValue(cpu.SocketCount),
		},
	}
}

// flattenGroups flattens []Group used by node provider data.
func flattenGroups(groups []import1.Group) []map[string]interface{} {
	if len(groups) == 0 {
		return nil
	}
	out := make([]map[string]interface{}, len(groups))
	for i, g := range groups {
		out[i] = map[string]interface{}{
			"id":   utils.StringValue(g.Id),
			"name": utils.StringValue(g.Name),
		}
	}
	return out
}

// flattenKVStringPairs flattens []KVStringPair used by node provider data tags.
func flattenKVStringPairs(tags []commonConfig.KVStringPair) []map[string]interface{} {
	if len(tags) == 0 {
		return nil
	}
	out := make([]map[string]interface{}, len(tags))
	for i, t := range tags {
		out[i] = map[string]interface{}{
			"name":  utils.StringValue(t.Name),
			"value": utils.StringValue(t.Value),
		}
	}
	return out
}

// flattenDisplayNameMapping flattens the DisplayNameMapping block for node provider data.
func flattenDisplayNameMapping(m *import1.DisplayNameMapping) []map[string]interface{} {
	if m == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"domain": utils.StringValue(m.Domain),
			"groups": utils.StringValue(m.Groups),
			"mode":   utils.StringValue(m.Mode),
			"tags":   utils.StringValue(m.Tags),
		},
	}
}

// flattenNodeProviderData flattens the NodeProviderData block for discovered nodes.
func flattenNodeProviderData(pd *import1.NodeProviderData) []map[string]interface{} {
	if pd == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"display_name_mapping": flattenDisplayNameMapping(pd.DisplayNameMapping),
			"domain":               utils.StringValue(pd.Domain),
			"groups":               flattenGroups(pd.Groups),
			"is_available":         utils.BoolValue(pd.IsAvailable),
			"is_configured":        utils.BoolValue(pd.IsConfigured),
			"mode":                 utils.StringValue(pd.Mode),
			"name":                 utils.StringValue(pd.Name),
			"tags":                 flattenKVStringPairs(pd.Tags),
		},
	}
}

// flattenManagedBy flattens the []ManagedBy enum slice into string names.
func flattenManagedBy(managedBy []import1.ManagedBy) []string {
	if len(managedBy) == 0 {
		return nil
	}
	out := make([]string, len(managedBy))
	for i := range managedBy {
		out[i] = managedBy[i].GetName()
	}
	return out
}

// flattenDiscoveredNode flattens a DiscoveredNode into a map for d.Set / list elements.
func flattenDiscoveredNode(node *import1.DiscoveredNode) map[string]interface{} {
	if node == nil {
		return nil
	}
	m := make(map[string]interface{})
	m["tenant_id"] = utils.StringValue(node.TenantId)
	m["ext_id"] = utils.StringValue(node.ExtId)
	m["links"] = common.FlattenLinks(node.Links)
	m["bmc_ip"] = flattenIPAddress(node.BmcIp)
	m["cpu_info"] = flattenCPUInfo(node.CpuInfo)
	m["identifiers"] = flattenNodeIdentifiers(node.Identifiers)
	m["managed_by"] = flattenManagedBy(node.ManagedBy)
	m["manufacturer"] = utils.StringValue(node.Manufacturer)
	m["memory_gb"] = utils.IntValue(node.MemoryGB)
	m["model"] = utils.StringValue(node.Model)
	m["provider_data"] = flattenNodeProviderData(node.ProviderData)
	if node.LastRefreshedTime != nil {
		m["last_refreshed_time"] = node.LastRefreshedTime.String()
	}
	return m
}

// flattenIPPool flattens an IpPool into a map for d.Set / list elements.
func flattenIPPool(pool *import1.IpPool) map[string]interface{} {
	if pool == nil {
		return nil
	}
	m := make(map[string]interface{})
	m["tenant_id"] = utils.StringValue(pool.TenantId)
	m["ext_id"] = utils.StringValue(pool.ExtId)
	m["links"] = common.FlattenLinks(pool.Links)
	m["name"] = utils.StringValue(pool.Name)
	m["reference_id"] = utils.StringValue(pool.ReferenceId)
	m["ipv4_available_count"] = utils.IntValue(pool.Ipv4AvailableCount)
	m["ipv6_available_count"] = utils.IntValue(pool.Ipv6AvailableCount)
	m["groups"] = flattenPoolGroups(pool.Groups)
	if pool.DisplayNameMapping != nil {
		m["display_name_mapping"] = []map[string]interface{}{
			{"group": utils.StringValue(pool.DisplayNameMapping.Group)},
		}
	}
	if pool.LastUpdatedTime != nil {
		m["last_updated_time"] = pool.LastUpdatedTime.String()
	}
	return m
}

// flattenMacPool flattens a MacPool into a map for d.Set / list elements.
func flattenMacPool(pool *import1.MacPool) map[string]interface{} {
	if pool == nil {
		return nil
	}
	m := make(map[string]interface{})
	m["tenant_id"] = utils.StringValue(pool.TenantId)
	m["ext_id"] = utils.StringValue(pool.ExtId)
	m["links"] = common.FlattenLinks(pool.Links)
	m["name"] = utils.StringValue(pool.Name)
	m["reference_id"] = utils.StringValue(pool.ReferenceId)
	m["available_count"] = utils.IntValue(pool.AvailableCount)
	m["groups"] = flattenPoolGroups(pool.Groups)
	if pool.DisplayNameMapping != nil {
		m["display_name_mapping"] = []map[string]interface{}{
			{"group": utils.StringValue(pool.DisplayNameMapping.Group)},
		}
	}
	if pool.LastUpdatedTime != nil {
		m["last_updated_time"] = pool.LastUpdatedTime.String()
	}
	return m
}

// flattenServerIdentityPool flattens a ServerIdentityPool into a map for d.Set / list elements.
func flattenServerIdentityPool(pool *import1.ServerIdentityPool) map[string]interface{} {
	if pool == nil {
		return nil
	}
	m := make(map[string]interface{})
	m["tenant_id"] = utils.StringValue(pool.TenantId)
	m["ext_id"] = utils.StringValue(pool.ExtId)
	m["links"] = common.FlattenLinks(pool.Links)
	m["name"] = utils.StringValue(pool.Name)
	m["reference_id"] = utils.StringValue(pool.ReferenceId)
	m["available_count"] = utils.IntValue(pool.AvailableCount)
	m["groups"] = flattenPoolGroups(pool.Groups)
	if pool.DisplayNameMapping != nil {
		m["display_name_mapping"] = []map[string]interface{}{
			{"group": utils.StringValue(pool.DisplayNameMapping.Group)},
		}
	}
	if pool.LastUpdatedTime != nil {
		m["last_updated_time"] = pool.LastUpdatedTime.String()
	}
	return m
}

// flattenHardwareProviderAuthTypes flattens the []AuthType enum slice into string names.
func flattenHardwareProviderAuthTypes(authTypes []import1.AuthType) []string {
	if len(authTypes) == 0 {
		return nil
	}
	out := make([]string, len(authTypes))
	for i := range authTypes {
		out[i] = authTypes[i].GetName()
	}
	return out
}

// flattenHardwareProvider flattens a HardwareProvider into a map for d.Set / list elements.
func flattenHardwareProvider(hp *import1.HardwareProvider) map[string]interface{} {
	if hp == nil {
		return nil
	}
	m := make(map[string]interface{})
	m["tenant_id"] = utils.StringValue(hp.TenantId)
	m["ext_id"] = utils.StringValue(hp.ExtId)
	m["links"] = common.FlattenLinks(hp.Links)
	m["name"] = utils.StringValue(hp.Name)
	m["vendor"] = utils.StringValue(hp.Vendor)
	m["available_connection_count"] = utils.IntValue(hp.AvailableConnectionCount)
	m["auth_types"] = flattenHardwareProviderAuthTypes(hp.AuthTypes)
	if hp.Type != nil {
		m["type"] = hp.Type.GetName()
	}
	return m
}
