package networkingv2

import (
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import1 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
	commonConfig "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networkingCommon/v1/config"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

var nicProfileCapabilityTypeAllowed = []string{
	"SRIOV",
	"DP_OFFLOAD",
	"PCIE_PASSTHROUGH",
}

var nicProfileOperatingModeAllowed = []string{
	"ETHERNET",
	"INFINIBAND",
}

var nicProfileOwnerTypeAllowed = []string{
	"USER",
	"SYSTEM",
}

func nicProfileCapabilityConfigSchema(required bool, computed bool) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: required,
		Optional: !required,
		Computed: computed,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"capability_type": {
					Type:         schema.TypeString,
					Required:     required,
					Optional:     !required,
					Computed:     computed,
					ValidateFunc: validation.StringInSlice(nicProfileCapabilityTypeAllowed, false),
				},
			},
		},
	}
}

func nicProfileHostNicReferencesSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"associated_vm_nic_references": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"compliance_status": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"ext_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"num_vfs": {
					Type:     schema.TypeInt,
					Computed: true,
				},
			},
		},
	}
}

func nicProfileHostNicExtIDsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Computed: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

// nicProfileMetadataSchema returns the schema for the `metadata` block.
// On the data source paths (extIDRequired=true) the block is purely
// computed; on the resource path it is user-supplied and the API
// surface for it is fully writeable, so we expose it as Optional only —
// flagging it Computed would suppress legitimate user-driven diffs on
// fields like categoryIds and *ReferenceId values that are part of the
// PUT body.
func nicProfileMetadataSchema(extIDRequired bool) *schema.Schema {
	if extIDRequired {
		return &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: DatasourceMetadataSchemaV2(),
			},
		}
	}
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: DatasourceMetadataSchemaV2(),
		},
	}
}

func nicProfileSchema(extIDRequired bool) map[string]*schema.Schema {
	operatingModeSchema := &schema.Schema{
		Type:     schema.TypeString,
		Optional: !extIDRequired,
		Computed: true,
	}
	if !extIDRequired {
		operatingModeSchema.ValidateFunc = validation.StringInSlice(nicProfileOperatingModeAllowed, false)
	}

	ownerTypeSchema := &schema.Schema{
		Type:     schema.TypeString,
		Optional: !extIDRequired,
		Computed: true,
	}
	if !extIDRequired {
		ownerTypeSchema.ValidateFunc = validation.StringInSlice(nicProfileOwnerTypeAllowed, false)
	}

	return map[string]*schema.Schema{
		"ext_id": {
			Type:     schema.TypeString,
			Required: extIDRequired,
			Optional: !extIDRequired,
			Computed: !extIDRequired,
		},
		"name": {
			Type:     schema.TypeString,
			Required: !extIDRequired,
			Computed: extIDRequired,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: !extIDRequired,
			Computed: true,
		},
		"capability_config":   nicProfileCapabilityConfigSchema(!extIDRequired, extIDRequired),
		"host_nic_ext_ids":    nicProfileHostNicExtIDsSchema(),
		"host_nic_references": nicProfileHostNicReferencesSchema(),
		"links": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"href": {Type: schema.TypeString, Computed: true},
					"rel":  {Type: schema.TypeString, Computed: true},
				},
			},
		},
		"metadata": nicProfileMetadataSchema(extIDRequired),
		"nic_family": {
			Type:     schema.TypeString,
			Required: !extIDRequired,
			Optional: false,
			Computed: extIDRequired,
		},
		"operating_mode": operatingModeSchema,
		"owner_type":     ownerTypeSchema,
		"project_ext_id": {
			Type:     schema.TypeString,
			Optional: !extIDRequired,
			Computed: true,
		},
		"tenant_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func flattenNicProfileCapabilityConfig(cfg *import1.CapabilityConfig) []map[string]interface{} {
	if cfg == nil {
		return nil
	}

	return []map[string]interface{}{{
		"capability_type": common.FlattenPtrEnum(cfg.CapabilityType),
	}}
}

func flattenNicProfileHostNicReferences(items []import1.HostNicReference) []map[string]interface{} {
	if len(items) == 0 {
		return nil
	}

	out := make([]map[string]interface{}, len(items))
	for i, item := range items {
		out[i] = map[string]interface{}{
			"associated_vm_nic_references": item.AssociatedVmNicReferences,
			"compliance_status":            common.FlattenPtrEnum(item.ComplianceStatus),
			"ext_id":                       utils.StringValue(item.ExtId),
			"num_vfs":                      utils.IntValue(item.NumVFs),
		}
	}
	return out
}

func flattenNicProfileHostNicExtIDs(items []import1.HostNicReference) []string {
	if len(items) == 0 {
		return nil
	}

	out := make([]string, 0, len(items))
	for _, item := range items {
		if extID := utils.StringValue(item.ExtId); extID != "" {
			out = append(out, extID)
		}
	}
	if len(out) == 0 {
		return nil
	}

	sort.Strings(out)
	return out
}

func flattenNicProfileEntity(m map[string]interface{}, item import1.NicProfile) map[string]interface{} {
	m["ext_id"] = utils.StringValue(item.ExtId)
	m["name"] = utils.StringValue(item.Name)
	m["description"] = utils.StringValue(item.Description)
	m["capability_config"] = flattenNicProfileCapabilityConfig(item.CapabilityConfig)
	m["host_nic_ext_ids"] = flattenNicProfileHostNicExtIDs(item.HostNicReferences)
	m["host_nic_references"] = flattenNicProfileHostNicReferences(item.HostNicReferences)
	m["links"] = flattenLinks(item.Links)
	m["metadata"] = flattenMetadata(item.Metadata)
	m["nic_family"] = utils.StringValue(item.NicFamily)
	m["operating_mode"] = common.FlattenPtrEnum(item.OperatingMode)
	m["owner_type"] = common.FlattenPtrEnum(item.OwnerType)
	m["project_ext_id"] = utils.StringValue(item.ProjectExtId)
	m["tenant_id"] = utils.StringValue(item.TenantId)
	return m
}

func flattenNicProfileProjectionEntity(m map[string]interface{}, item import1.NicProfileProjection) map[string]interface{} {
	m["ext_id"] = utils.StringValue(item.ExtId)
	m["name"] = utils.StringValue(item.Name)
	m["description"] = utils.StringValue(item.Description)
	m["capability_config"] = flattenNicProfileCapabilityConfig(item.CapabilityConfig)
	m["host_nic_ext_ids"] = flattenNicProfileHostNicExtIDs(item.HostNicReferences)
	m["host_nic_references"] = flattenNicProfileHostNicReferences(item.HostNicReferences)
	m["links"] = flattenLinks(item.Links)
	m["metadata"] = flattenMetadata(item.Metadata)
	m["nic_family"] = utils.StringValue(item.NicFamily)
	m["operating_mode"] = common.FlattenPtrEnum(item.OperatingMode)
	m["owner_type"] = common.FlattenPtrEnum(item.OwnerType)
	m["project_ext_id"] = utils.StringValue(item.ProjectExtId)
	m["tenant_id"] = utils.StringValue(item.TenantId)
	return m
}

func expandNicProfileCapabilityConfig(v interface{}) *import1.CapabilityConfig {
	if v == nil {
		return nil
	}

	items := v.([]interface{})
	if len(items) == 0 || items[0] == nil {
		return nil
	}

	cfgMap := items[0].(map[string]interface{})
	cfg := &import1.CapabilityConfig{}
	if capabilityType, ok := cfgMap["capability_type"].(string); ok && capabilityType != "" {
		cfg.CapabilityType = common.ExpandEnum[import1.CapabilityType](capabilityType)
	}
	return cfg
}

func expandNicProfileOperatingMode(value string) *commonConfig.OperatingMode {
	if value == "" {
		return nil
	}
	return common.ExpandEnum[commonConfig.OperatingMode](value)
}

func expandNicProfileOwnerType(value string) *import1.NicProfileOwnerType {
	if value == "" {
		return nil
	}
	return common.ExpandEnum[import1.NicProfileOwnerType](value)
}

func expandNicProfile(item map[string]interface{}) import1.NicProfile {
	profile := import1.NicProfile{}

	if v, ok := item["name"].(string); ok && v != "" {
		profile.Name = utils.StringPtr(v)
	}
	if v, ok := item["description"].(string); ok && v != "" {
		profile.Description = utils.StringPtr(v)
	}
	if v, ok := item["capability_config"]; ok {
		profile.CapabilityConfig = expandNicProfileCapabilityConfig(v)
	}
	if v, ok := item["metadata"]; ok {
		profile.Metadata = expandMetadata(v.([]interface{}))
	}
	if v, ok := item["nic_family"].(string); ok && v != "" {
		profile.NicFamily = utils.StringPtr(v)
	}
	if v, ok := item["operating_mode"].(string); ok && v != "" {
		profile.OperatingMode = expandNicProfileOperatingMode(v)
	}
	if v, ok := item["owner_type"].(string); ok && v != "" {
		profile.OwnerType = expandNicProfileOwnerType(v)
	}
	if v, ok := item["project_ext_id"].(string); ok && v != "" {
		profile.ProjectExtId = utils.StringPtr(v)
	}

	return profile
}

func expandStringSet(v interface{}) []string {
	if v == nil {
		return nil
	}

	var items []interface{}
	switch val := v.(type) {
	case *schema.Set:
		items = val.List()
	case []interface{}:
		items = val
	default:
		return nil
	}

	out := make([]string, 0, len(items))
	for _, item := range items {
		if s, ok := item.(string); ok && s != "" {
			out = append(out, s)
		}
	}
	if len(out) == 0 {
		return nil
	}

	sort.Strings(out)
	return out
}
