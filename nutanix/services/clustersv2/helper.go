package clustersv2

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	import4 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/common/v1/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/prism/v4/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
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
	ShouldSkipAddNode           bool
	ShouldSkipHostNetworking    *bool
	ShouldSkipPreExpandChecks   bool
	ShouldSkipDiscovery         bool
	ShouldSkipImaging           bool
	ShouldValidateRackAwareness bool
	IsNosCompatible             bool
	IsComputeOnly               bool
	IsNeverScheduleable         bool
	IsLightCompute              bool
	HypervisorHostname          string

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

// isFieldSetInRawConfig checks if a specific field was explicitly set in the raw TF config
// for a node_list item matching the given controller VM IP.
func isFieldSetInRawConfig(d *schema.ResourceData, nodeIP string, fieldName string) bool {
	rawConfig := d.GetRawConfig()
	if rawConfig.IsNull() {
		return false
	}

	nodesAttr := rawConfig.GetAttr("nodes")
	if nodesAttr.IsNull() || !nodesAttr.CanIterateElements() {
		return false
	}

	for nodesIt := nodesAttr.ElementIterator(); nodesIt.Next(); {
		_, nodeVal := nodesIt.Element()
		if nodeVal.IsNull() {
			continue
		}

		nodeListAttr := nodeVal.GetAttr("node_list")
		if nodeListAttr.IsNull() || !nodeListAttr.CanIterateElements() {
			continue
		}

		for nodeListIt := nodeListAttr.ElementIterator(); nodeListIt.Next(); {
			_, nodeListItem := nodeListIt.Element()
			if nodeListItem.IsNull() {
				continue
			}

			// Check if this node_list item matches by controller_vm_ip
			controllerVmIpAttr := nodeListItem.GetAttr("controller_vm_ip")
			if controllerVmIpAttr.IsNull() || !controllerVmIpAttr.CanIterateElements() {
				continue
			}

			for cvmIt := controllerVmIpAttr.ElementIterator(); cvmIt.Next(); {
				_, cvmVal := cvmIt.Element()
				if cvmVal.IsNull() {
					continue
				}

				// Check ipv4
				ipv4Attr := cvmVal.GetAttr("ipv4")
				if !ipv4Attr.IsNull() && ipv4Attr.CanIterateElements() {
					for ipv4It := ipv4Attr.ElementIterator(); ipv4It.Next(); {
						_, ipv4Val := ipv4It.Element()
						if !ipv4Val.IsNull() {
							valueAttr := ipv4Val.GetAttr("value")
							if !valueAttr.IsNull() && valueAttr.AsString() == nodeIP {
								// Found matching node, check if field is set
								fieldAttr := nodeListItem.GetAttr(fieldName)
								return !fieldAttr.IsNull()
							}
						}
					}
				}

				// Check ipv6
				ipv6Attr := cvmVal.GetAttr("ipv6")
				if !ipv6Attr.IsNull() && ipv6Attr.CanIterateElements() {
					for ipv6It := ipv6Attr.ElementIterator(); ipv6It.Next(); {
						_, ipv6Val := ipv6It.Element()
						if !ipv6Val.IsNull() {
							valueAttr := ipv6Val.GetAttr("value")
							if !valueAttr.IsNull() && valueAttr.AsString() == nodeIP {
								// Found matching node, check if field is set
								fieldAttr := nodeListItem.GetAttr(fieldName)
								return !fieldAttr.IsNull()
							}
						}
					}
				}
			}
		}
	}

	return false
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
						// Only set ShouldSkipHostNetworking if explicitly set in TF config
						if isFieldSetInRawConfig(d, nodeIP, "should_skip_host_networking") {
							if v, ok := itemMap["should_skip_host_networking"].(bool); ok {
								flags.ShouldSkipHostNetworking = utils.BoolPtr(v)
							}
						}
						if isFieldSetInRawConfig(d, nodeIP, "should_skip_pre_expand_checks") {
							if v, ok := itemMap["should_skip_pre_expand_checks"].(bool); ok {
								flags.ShouldSkipPreExpandChecks = v
							}
						}
						if isFieldSetInRawConfig(d, nodeIP, "should_skip_discovery") {
							if v, ok := itemMap["should_skip_discovery"].(bool); ok {
								flags.ShouldSkipDiscovery = v
							}
						}
						if isFieldSetInRawConfig(d, nodeIP, "should_skip_imaging") {
							if v, ok := itemMap["should_skip_imaging"].(bool); ok {
								flags.ShouldSkipImaging = v
							}
						}
						if isFieldSetInRawConfig(d, nodeIP, "should_validate_rack_awareness") {
							if v, ok := itemMap["should_validate_rack_awareness"].(bool); ok {
								flags.ShouldValidateRackAwareness = v
							}
						}
						if isFieldSetInRawConfig(d, nodeIP, "is_nos_compatible") {
							if v, ok := itemMap["is_nos_compatible"].(bool); ok {
								flags.IsNosCompatible = v
							}
						}
						if isFieldSetInRawConfig(d, nodeIP, "is_compute_only") {
							if v, ok := itemMap["is_compute_only"].(bool); ok {
								flags.IsComputeOnly = v
							}
						}
						if isFieldSetInRawConfig(d, nodeIP, "is_never_scheduleable") {
							if v, ok := itemMap["is_never_scheduleable"].(bool); ok {
								flags.IsNeverScheduleable = v
							}
						}
						if isFieldSetInRawConfig(d, nodeIP, "is_light_compute") {
							if v, ok := itemMap["is_light_compute"].(bool); ok {
								flags.IsLightCompute = v
							}
						}
						if isFieldSetInRawConfig(d, nodeIP, "hypervisor_hostname") {
							if v, ok := itemMap["hypervisor_hostname"].(string); ok {
								flags.HypervisorHostname = v
							}
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

// ############################
// ### Category helpers ###
// ############################

// FindRemovedCategories returns categories that are in old but not in new
func FindRemovedCategories(oldCategories, newCategories []string) []string {
	if len(oldCategories) == 0 {
		return []string{}
	}

	// Create a map of new categories for quick lookup
	newCategoriesMap := make(map[string]bool)
	for _, cat := range newCategories {
		if cat != "" {
			newCategoriesMap[cat] = true
		}
	}

	// Find categories in old that are not in new
	removed := make([]string, 0)
	for _, cat := range oldCategories {
		if cat != "" && !newCategoriesMap[cat] {
			removed = append(removed, cat)
		}
	}

	return removed
}

// FindAddedCategories returns categories that are in new but not in old
func FindAddedCategories(oldCategories, newCategories []string) []string {
	if len(newCategories) == 0 {
		return []string{}
	}

	// Create a map of old categories for quick lookup
	oldCategoriesMap := make(map[string]bool)
	for _, cat := range oldCategories {
		if cat != "" {
			oldCategoriesMap[cat] = true
		}
	}

	// Find categories in new that are not in old
	added := make([]string, 0)
	for _, cat := range newCategories {
		if cat != "" && !oldCategoriesMap[cat] {
			added = append(added, cat)
		}
	}

	return added
}

// UpdateClusterCategories handles category association and disassociation for a cluster
// This is a shared function used by both cluster entity and cluster categories resources
func UpdateClusterCategories(ctx context.Context, d *schema.ResourceData, meta interface{}, clusterExtID string, oldCategoriesRaw, newCategoriesRaw interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI
	taskconn := meta.(*conns.Client).PrismAPI

	// Convert to slices - handles both TypeList and TypeSet
	oldCategoriesList := common.InterfaceToSlice(oldCategoriesRaw)
	newCategoriesList := common.InterfaceToSlice(newCategoriesRaw)

	// Convert to string slices for easier comparison
	oldCategories := common.ExpandListOfString(oldCategoriesList)
	newCategories := common.ExpandListOfString(newCategoriesList)

	// Find categories to disassociate (present in old but not in new)
	categoriesToDisassociate := FindRemovedCategories(oldCategories, newCategories)

	// Find categories to associate (present in new but not in old)
	categoriesToAssociate := FindAddedCategories(oldCategories, newCategories)

	log.Printf("[DEBUG] Category changes - To disassociate: %v, To associate: %v", categoriesToDisassociate, categoriesToAssociate)

	// Disassociate removed categories first
	if len(categoriesToDisassociate) > 0 {
		body := &config.CategoryEntityReferences{
			Categories: categoriesToDisassociate,
		}

		aJSON, _ := json.MarshalIndent(body, "", " ")
		log.Printf("[DEBUG] Disassociate Categories from Cluster Request Body: %s", string(aJSON))

		resp, err := conn.ClusterEntityAPI.DisassociateCategoriesFromCluster(utils.StringPtr(clusterExtID), body)
		if err != nil {
			return diag.Errorf("error while disassociating categories from cluster: %v", err)
		}

		TaskRef := resp.Data.GetValue().(import1.TaskReference)
		taskUUID := TaskRef.ExtId

		// Wait for the categories to be disassociated
		stateConf := &resource.StateChangeConf{
			Pending: []string{"PENDING", "RUNNING", "QUEUED"},
			Target:  []string{"SUCCEEDED"},
			Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
			Timeout: d.Timeout(schema.TimeoutUpdate),
		}

		if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
			resourceUUID, _ := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
			if resourceUUID != nil {
				rUUID := resourceUUID.Data.GetValue().(import2.Task)
				aJSON, _ = json.MarshalIndent(rUUID, "", "  ")
				log.Printf("[DEBUG] Error Disassociate Categories from Cluster Task Details: %s", string(aJSON))
			}
			return diag.Errorf("error waiting for categories to be disassociated from cluster (%s): %s", utils.StringValue(taskUUID), errWaitTask)
		}

		// Get task details
		taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
		if err != nil {
			return diag.Errorf("error while fetching disassociate categories from cluster task: %v", err)
		}

		aJSON, _ = json.MarshalIndent(taskResp, "", "  ")
		log.Printf("[DEBUG] Disassociate categories from cluster task details: %s", string(aJSON))
	}

	// Associate added categories
	if len(categoriesToAssociate) > 0 {
		body := config.CategoryEntityReferences{
			Categories: categoriesToAssociate,
		}

		aJSON, _ := json.MarshalIndent(body, "", " ")
		log.Printf("[DEBUG] Associate Categories to Cluster Request Body: %s", string(aJSON))

		resp, err := conn.ClusterEntityAPI.AssociateCategoriesToCluster(utils.StringPtr(clusterExtID), &body)
		if err != nil {
			return diag.Errorf("error while associating categories to cluster: %v", err)
		}

		TaskRef := resp.Data.GetValue().(import1.TaskReference)
		taskUUID := TaskRef.ExtId

		// Wait for the categories to be associated
		stateConf := &resource.StateChangeConf{
			Pending: []string{"PENDING", "RUNNING", "QUEUED"},
			Target:  []string{"SUCCEEDED"},
			Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
			Timeout: d.Timeout(schema.TimeoutUpdate),
		}

		if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
			return diag.Errorf("error waiting for categories to be associated to the cluster (%s): %s", utils.StringValue(taskUUID), errWaitTask)
		}

		// Get task details
		taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
		if err != nil {
			return diag.Errorf("error while fetching associate categories to cluster task: %v", err)
		}

		aJSON, _ = json.Marshal(taskResp)
		log.Printf("[DEBUG] Associate categories to cluster task details: %s", string(aJSON))
	}

	return nil
}

// ##################################
// ### Authorized Public Key Hash ###
// ##################################
// authorizedPublicKeyHash --- hash function for authorized_public_key set ---
func authorizedPublicKeyHash(v interface{}) int {
	m := v.(map[string]interface{})
	name := ""
	key := ""

	if val, ok := m["name"].(string); ok {
		name = val
	}
	if val, ok := m["key"].(string); ok {
		key = val
	}

	// Combine name and key for hashing
	hashInput := fmt.Sprintf("%s-%s", name, key)
	return schema.HashString(hashInput)
}
