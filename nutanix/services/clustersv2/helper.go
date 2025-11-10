package clustersv2

import (
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	import4 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/common/v1/config"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/clusters"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// ############################
// ### ETag header helper ###
// ############################
func getEtagHeader(resp interface{}, conn *clusters.Client) map[string]interface{} {
	// Extract E-Tag Header
	etagValue := conn.ClusterEntityAPI.ApiClient.GetEtag(resp)

	// Extract E-Tag Header
	args := make(map[string]interface{})
	args["If-Match"] = utils.StringPtr(etagValue)

	return args
}

// ###########################
// ### Node list helpers ###
// ###########################
// hashNodeItem --- hash function for node set ---
func hashNodeItem(v interface{}) int {
	m, ok := v.(map[string]interface{})
	if !ok {
		return 0
	}

	// Extract controller_vm_ip list
	controllerVMs, ok := m["controller_vm_ip"].([]interface{})
	if !ok || len(controllerVMs) == 0 {
		return 0
	}

	// Extract first controller_vm_ip map
	ipMap, ok := controllerVMs[0].(map[string]interface{})
	if !ok {
		return 0
	}

	var ipValue string

	// Prefer IPv4 if available
	if ipv4List, ok := ipMap["ipv4"].([]interface{}); ok && len(ipv4List) > 0 {
		if ipv4Map, ok := ipv4List[0].(map[string]interface{}); ok {
			if val, ok := ipv4Map["value"].(string); ok {
				ipValue = val
			}
		}
	}

	// Fall back to IPv6 if IPv4 is missing
	if ipValue == "" {
		if ipv6List, ok := ipMap["ipv6"].([]interface{}); ok && len(ipv6List) > 0 {
			if ipv6Map, ok := ipv6List[0].(map[string]interface{}); ok {
				if val, ok := ipv6Map["value"].(string); ok {
					ipValue = val
				}
			}
		}
	}

	// Compute hash — use FNV for deterministic stable hash
	if ipValue == "" {
		return 0
	}
	h := fnv.New32a()
	_, _ = h.Write([]byte(strings.ToLower(strings.TrimSpace(ipValue))))
	return int(h.Sum32())
}

// --- stringify IP helper used as fallback key ---
func ipToKey(ip *import4.IPAddress) string {
	if ip == nil {
		return ""
	}
	if ip.Ipv4 != nil && utils.StringValue(ip.Ipv4.Value) != "" {
		return fmt.Sprintf("ipv4:%s/%d", utils.StringValue(ip.Ipv4.Value), utils.IntValue(ip.Ipv4.PrefixLength))
	}
	if ip.Ipv6 != nil && utils.StringValue(ip.Ipv6.Value) != "" {
		return fmt.Sprintf("ipv6:%s/%d", utils.StringValue(ip.Ipv6.Value), utils.IntValue(ip.Ipv6.PrefixLength))
	}
	return ""
}

// --- get IP type as string for diagnostics ---
func getIPType(ip *import4.IPAddress) string {
	if ip == nil {
		return "UNKNOWN"
	}
	if ip.Ipv4 != nil && utils.StringValue(ip.Ipv4.Value) != "" {
		return "IPV4"
	}
	if ip.Ipv6 != nil && utils.StringValue(ip.Ipv6.Value) != "" {
		return "IPV6"
	}
	return "UNKNOWN"
}

// --- equality checks (nil-safe, checks meaningful fields) ---
func ipv4Equal(a, b *import4.IPv4Address) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return utils.StringValue(a.Value) == utils.StringValue(b.Value) && utils.IntValue(a.PrefixLength) == utils.IntValue(b.PrefixLength)
}

func ipv6Equal(a, b *import4.IPv6Address) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return utils.StringValue(a.Value) == utils.StringValue(b.Value) && utils.IntValue(a.PrefixLength) == utils.IntValue(b.PrefixLength)
}

func ipAddressEqual(a, b *import4.IPAddress) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return ipv4Equal(a.Ipv4, b.Ipv4) && ipv6Equal(a.Ipv6, b.Ipv6)
}

func nodeEqual(a, b *config.NodeListItemReference) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	// Compare UUID (string) and both IPs
	if utils.StringValue(a.NodeUuid) != "" && utils.StringValue(b.NodeUuid) != "" {
		if utils.StringValue(a.NodeUuid) != utils.StringValue(b.NodeUuid) {
			return false
		}
	}

	if !ipAddressEqual(a.ControllerVmIp, b.ControllerVmIp) {
		return false
	}

	if !ipAddressEqual(a.HostIp, b.HostIp) {
		return false
	}

	// if you want to compare Reserved_/UnknownFields_ you'd add checks here
	return true
}

// --- produce stable key for a node (prefer UUID) ---
func nodeKeyCandidates(n *config.NodeListItemReference) []string {
	keys := []string{}
	if n == nil {
		return keys
	}
	if k := utils.StringValue(n.NodeUuid); k != "" {
		keys = append(keys, "uuid:"+k)
	}
	if s := ipToKey(n.ControllerVmIp); s != "" {
		keys = append(keys, "ctrl:"+s)
	}
	if s := ipToKey(n.HostIp); s != "" {
		keys = append(keys, "host:"+s)
	}
	return keys
}

// ChangedPair Diff result types ---
type ChangedPair struct {
	Old config.NodeListItemReference `json:"old"`
	New config.NodeListItemReference `json:"new"`
}

// DiffNodes returns lists of nodes that were added, removed, or changed.
// existing = current cluster nodes from API
// newNodes = desired nodes from the resource/state
func DiffNodes(d *schema.ResourceData, existing, newNodes []config.NodeListItemReference) (
	added, removed []NodeWithFlags, changed []ChangedPair) {
	matchedNew := make(map[int]bool)

	for _, old := range existing {
		oldCopy := old
		var found bool

		for j, newNode := range newNodes {
			if matchedNew[j] {
				continue
			}

			newCopy := newNode

			match := false
			for _, kOld := range nodeKeyCandidates(&oldCopy) {
				for _, kNew := range nodeKeyCandidates(&newCopy) {
					if kOld == kNew {
						match = true
						break
					}
				}
			}

			if match {
				found = true
				matchedNew[j] = true

				if !nodeEqual(&oldCopy, &newCopy) {
					changed = append(changed, ChangedPair{Old: oldCopy, New: newCopy})
				}
				break
			}
		}

		if !found {
			// Removed node → extract flags for removal (if present)
			flags := extractNodeFlags(d, oldCopy)
			removed = append(removed, NodeWithFlags{
				Node:  oldCopy,
				Flags: flags,
			})
		}
	}

	for j, newNode := range newNodes {
		if !matchedNew[j] {
			flags := extractNodeFlags(d, newNode) // compute flags once here
			added = append(added, NodeWithFlags{
				Node:  newNode,
				Flags: flags,
			})
		}
	}

	return added, removed, changed
}

type NodeFlags struct {
	// Add node flags
	ShouldSkipAddNode         bool
	ShouldSkipHostNetworking  bool
	ShouldSkipPreExpandChecks bool

	// Remove node flags
	ShouldSkipRemove       bool
	ShouldSkipPrechecks    bool
	ShouldSkipUpgradeCheck bool
	SkipSpaceCheck         bool
	ShouldSkipAddCheck     bool
}
type NodeWithFlags struct {
	Node  config.NodeListItemReference
	Flags NodeFlags
}

// extractNodeFlags finds a node from the Terraform diff (nodes) that matches
// the node (by comparing CVM IP) and extracts both add/remove operation flags.
func extractNodeFlags(d *schema.ResourceData, node config.NodeListItemReference) NodeFlags {
	flags := NodeFlags{}

	nodes := common.InterfaceToSlice(d.Get("nodes"))
	if len(nodes) == 0 {
		return flags
	}

	for _, n := range nodes {
		nMap, ok := n.(map[string]interface{})
		if !ok {
			continue
		}

		// ====== Handle add-node flags from node_list ======
		nodeLists := common.InterfaceToSlice(nMap["node_list"])
		for _, nodeItem := range nodeLists {
			itemMap, ok := nodeItem.(map[string]interface{})
			if !ok {
				continue
			}

			controllerIPs := common.InterfaceToSlice(itemMap["controller_vm_ip"])
			for _, ipBlock := range controllerIPs {
				ipMap, _ := ipBlock.(map[string]interface{})

				var ipList []interface{}
				var nodeIP string

				// Determine which IP family to use
				if node.ControllerVmIp != nil {
					if node.ControllerVmIp.Ipv4 != nil {
						ipList = common.InterfaceToSlice(ipMap["ipv4"])
						nodeIP = utils.StringValue(node.ControllerVmIp.Ipv4.Value)
					} else if node.ControllerVmIp.Ipv6 != nil {
						ipList = common.InterfaceToSlice(ipMap["ipv6"])
						nodeIP = utils.StringValue(node.ControllerVmIp.Ipv6.Value)
					}
				}

				// Skip if no IPs available
				if len(ipList) == 0 || nodeIP == "" {
					continue
				}

				// Check if any IP matches the node's CVM IP
				for _, ip := range ipList {
					ipMapItem, ok := ip.(map[string]interface{})
					if !ok {
						continue
					}
					if v, ok := ipMapItem["value"].(string); ok && v == nodeIP {
						if v, ok := itemMap["should_skip_add_node"].(bool); ok {
							flags.ShouldSkipAddNode = v
						}
						if v, ok := itemMap["should_skip_host_networking"].(bool); ok {
							flags.ShouldSkipHostNetworking = v
						}
						if v, ok := itemMap["should_skip_pre_expand_checks"].(bool); ok {
							flags.ShouldSkipPreExpandChecks = v
						}
						break
					}
				}
			}
		}

		// ====== Handle remove-node flags from remove_node_params ======
		removeParams := common.InterfaceToSlice(nMap["remove_node_params"])
		if len(removeParams) > 0 {
			rpMap, ok := removeParams[0].(map[string]interface{})
			if ok {
				if v, ok := rpMap["should_skip_remove"].(bool); ok {
					flags.ShouldSkipRemove = v
				}
				if v, ok := rpMap["should_skip_prechecks"].(bool); ok {
					flags.ShouldSkipPrechecks = v
				}

				extras := common.InterfaceToSlice(rpMap["extra_params"])
				if len(extras) > 0 {
					extraMap, ok := extras[0].(map[string]interface{})
					if ok {
						if v, ok := extraMap["should_skip_upgrade_check"].(bool); ok {
							flags.ShouldSkipUpgradeCheck = v
						}
						if v, ok := extraMap["skip_space_check"].(bool); ok {
							flags.SkipSpaceCheck = v
						}
						if v, ok := extraMap["should_skip_add_check"].(bool); ok {
							flags.ShouldSkipAddCheck = v
						}
					}
				}
			}
		}
	}

	return flags
}

// ############################
// ### Cluster not found error ###
// ############################

// ClusterNotFoundError is returned when a cluster with the specified name is not found.
type ClusterNotFoundError struct {
	Name string
	Err  error
}

func (e *ClusterNotFoundError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("cluster not found: %s, %v", e.Name, e.Err)
	}
	return fmt.Sprintf("cluster not found: %s", e.Name)
}
