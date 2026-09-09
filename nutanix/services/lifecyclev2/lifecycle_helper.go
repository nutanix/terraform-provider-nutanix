package lifecyclev2

import (
	commonConfig "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/common/v1/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/common/v1/response"
	import1 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/config"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// flattenLifecycleLinks converts the SDK HATEOAS links into the Terraform
// representation shared by every lifecyclev2 resource and datasource.
func flattenLifecycleLinks(links []import2.ApiLink) []map[string]interface{} {
	if len(links) == 0 {
		return nil
	}
	out := make([]map[string]interface{}, 0, len(links))
	for _, link := range links {
		out = append(out, map[string]interface{}{
			"href": utils.StringValue(link.Href),
			"rel":  utils.StringValue(link.Rel),
		})
	}
	return out
}

// flattenClaimTokenNode flattens a single ClaimTokenNode into the Terraform map
// used by the node datasources.
func flattenClaimTokenNode(node *import1.ClaimTokenNode) map[string]interface{} {
	if node == nil {
		return nil
	}
	m := map[string]interface{}{
		"ext_id":          utils.StringValue(node.ExtId),
		"tenant_id":       utils.StringValue(node.TenantId),
		"host_version":    utils.StringValue(node.HostVersion),
		"manufacturer":    utils.StringValue(node.Manufacturer),
		"model":           utils.StringValue(node.Model),
		"links":           flattenLifecycleLinks(node.Links),
		"cpu_info":        flattenCPUInfo(node.CpuInfo),
		"identifiers":     flattenNodeIdentifiers(node.Identifiers),
		"managed_by":      flattenClaimTokenNodeManagers(node.ManagedBy),
		"network_details": flattenNodeNetworkDetails(node.NetworkDetails),
	}
	if node.CreatedTime != nil {
		m["created_time"] = node.CreatedTime.String()
	}
	if node.MemoryGB != nil {
		m["memory_gb"] = utils.IntValue(node.MemoryGB)
	}
	if node.HostType != nil {
		m["host_type"] = node.HostType.GetName()
	}
	return m
}

func flattenClaimTokenNodes(nodes []import1.ClaimTokenNode) []map[string]interface{} {
	if len(nodes) == 0 {
		return nil
	}
	out := make([]map[string]interface{}, 0, len(nodes))
	for i := range nodes {
		out = append(out, flattenClaimTokenNode(&nodes[i]))
	}
	return out
}

func flattenCPUInfo(cpu *import1.CpuInfo) []map[string]interface{} {
	if cpu == nil {
		return nil
	}
	m := map[string]interface{}{
		"manufacturer": utils.StringValue(cpu.Manufacturer),
		"model":        utils.StringValue(cpu.Model),
	}
	if cpu.CapacityGHz != nil {
		m["capacity_ghz"] = float64(*cpu.CapacityGHz)
	}
	if cpu.LogicalCoreCount != nil {
		m["logical_core_count"] = utils.IntValue(cpu.LogicalCoreCount)
	}
	if cpu.SocketCount != nil {
		m["socket_count"] = utils.IntValue(cpu.SocketCount)
	}
	return []map[string]interface{}{m}
}

func flattenNodeIdentifiers(identifiers []import1.NodeIdentifier) []map[string]interface{} {
	if len(identifiers) == 0 {
		return nil
	}
	out := make([]map[string]interface{}, 0, len(identifiers))
	for _, id := range identifiers {
		m := map[string]interface{}{
			"value": utils.StringValue(id.Value),
		}
		if id.Type != nil {
			m["type"] = id.Type.GetName()
		}
		out = append(out, m)
	}
	return out
}

func flattenClaimTokenNodeManagers(managers []import1.ClaimTokenNodeManager) []string {
	if len(managers) == 0 {
		return nil
	}
	out := make([]string, 0, len(managers))
	for _, mgr := range managers {
		out = append(out, mgr.GetName())
	}
	return out
}

func flattenNodeNetworkDetails(details *import1.NodeNetworkDetails) []map[string]interface{} {
	if details == nil {
		return nil
	}
	m := map[string]interface{}{
		"bmc":  flattenNetworkDetails(details.Bmc),
		"cvm":  flattenComponentNetworkDetails(details.Cvm),
		"host": flattenComponentNetworkDetails(details.Host),
	}
	return []map[string]interface{}{m}
}

func flattenComponentNetworkDetails(details *import1.ComponentNetworkDetails) []map[string]interface{} {
	if details == nil {
		return nil
	}
	m := map[string]interface{}{
		"management_network": flattenNetworkDetails(details.ManagementNetwork),
	}
	return []map[string]interface{}{m}
}

func flattenNetworkDetails(details *import1.NetworkDetails) []map[string]interface{} {
	if details == nil {
		return nil
	}
	m := map[string]interface{}{
		"gateway": flattenIPAddress(details.Gateway),
		"ip":      flattenIPAddress(details.Ip),
	}
	if details.VlanId != nil {
		m["vlan_id"] = utils.IntValue(details.VlanId)
	}
	return []map[string]interface{}{m}
}

func flattenIPAddress(ip *commonConfig.IPAddress) []map[string]interface{} {
	if ip == nil {
		return nil
	}
	m := map[string]interface{}{}
	if ip.Ipv4 != nil {
		v := map[string]interface{}{
			"value": utils.StringValue(ip.Ipv4.Value),
		}
		if ip.Ipv4.PrefixLength != nil {
			v["prefix_length"] = utils.IntValue(ip.Ipv4.PrefixLength)
		}
		m["ipv4"] = []map[string]interface{}{v}
	}
	if ip.Ipv6 != nil {
		v := map[string]interface{}{
			"value": utils.StringValue(ip.Ipv6.Value),
		}
		if ip.Ipv6.PrefixLength != nil {
			v["prefix_length"] = utils.IntValue(ip.Ipv6.PrefixLength)
		}
		m["ipv6"] = []map[string]interface{}{v}
	}
	return []map[string]interface{}{m}
}
