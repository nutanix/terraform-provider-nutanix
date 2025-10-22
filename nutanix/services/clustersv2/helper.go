package clustersv2

import (
	"fmt"

	"github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	import4 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/common/v1/config"
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
func DiffNodes(existing, newNodes []config.NodeListItemReference) (
	added, removed []config.NodeListItemReference, changed []ChangedPair) {
	matchedNew := make(map[int]bool)

	for _, old := range existing {
		oldCopy := old // ✅ local copy to safely take address
		var found bool

		for j, newNode := range newNodes {
			if matchedNew[j] {
				continue
			}

			newCopy := newNode // ✅ local copy to safely take address

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
			removed = append(removed, oldCopy)
		}
	}

	for j, newNode := range newNodes {
		if !matchedNew[j] {
			added = append(added, newNode)
		}
	}

	return added, removed, changed
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
