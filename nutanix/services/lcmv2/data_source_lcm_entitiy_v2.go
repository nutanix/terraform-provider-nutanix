package lcmv2

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lcmConfigPkg "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/common/v1/config"
	"github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/common"
	lcmEntityPkg "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/resources"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixLcmEntityV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixLcmEntityV2Create,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": schemaForLinks(),
			"entity_class": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"entity_model": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"entity_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"entity_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hardware_family": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"entity_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"location_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"target_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_updated_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"device_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"entity_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemaForKeyValuePairs(),
			},
			"child_entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Description: "List of String",
					Type:        schema.TypeString,
				},
			},
			"available_versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemaForAvailableVersions(),
			},
			"sub_entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemaForSubEntities(),
			},
			"cluster_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hardware_vendor": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DatasourceNutanixLcmEntityV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LcmAPI
	extID := d.Get("ext_id").(string)

	resp, err := conn.LcmEntitiesAPIInstance.GetEntityById(utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while fetching the Lcm etity : %v", err)
	}

	lcmEntity := resp.Data.GetValue().(lcmEntityPkg.Entity)
	if err := d.Set("tenant_id", lcmEntity.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(lcmEntity.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("entity_class", lcmEntity.EntityClass); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("entity_model", lcmEntity.EntityModel); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("entity_type", flattenEntityTypes(lcmEntity.EntityType)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("entity_version", lcmEntity.EntityVersion); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("hardware_family", lcmEntity.HardwareFamily); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("entity_description", lcmEntity.EntityDescription); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("location_info", flattenLocationInfo(lcmEntity.LocationInfo)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("target_version", lcmEntity.TargetVersion); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("last_updated_time", flattenTime(lcmEntity.LastUpdatedTime)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("device_id", lcmEntity.DeviceId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("group_uuid", lcmEntity.GroupUuid); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("entity_details", flattenKeyValuePairs(lcmEntity.EntityDetails)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("child_entities", lcmEntity.ChildEntities); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("available_versions", flattenAvailableVersions(lcmEntity.AvailableVersions)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("sub_entities", flattenSubEntities(lcmEntity.SubEntities)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_ext_id", lcmEntity.ClusterExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("hardware_vendor", lcmEntity.HardwareVendor); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(lcmEntity.ExtId))
	return nil
}

func schemaForKeyValuePairs() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"value": schemaForValue(),
		},
	}
}

func schemaForValue() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"string": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"integer": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"boolean": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"string_list": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"object": {
					Type:     schema.TypeMap,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"map_of_strings": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"map": {
								Type:     schema.TypeMap,
								Computed: true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
						},
					},
				},
				"integer_list": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeInt,
					},
				},
			},
		},
	}
}

func schemaForAvailableVersions() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"available_version_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"order": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"disablement_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"release_notes": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"release_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"child_entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Description: "List of String",
					Type:        schema.TypeString,
				},
			},
			"group_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dependencies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"links": schemaForLinks(),
						"entity_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"entity_model": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"entity_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"entity_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hardware_family": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dependent_versions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     schemaForKeyValuePairs(),
						},
					},
				},
			},
		},
	}
}

func schemaForSubEntities() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": schemaForLinks(),
			"entity_class": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"entity_model": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"entity_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"entity_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hardware_family": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// flatten Funcs
func flattenLocationInfo(locationInfo *common.LocationInfo) []map[string]interface{} {
	if locationInfo == nil {
		return nil
	}

	locationInfoList := make([]map[string]interface{}, 0)

	locationInfoMap := map[string]interface{}{
		"uuid":          locationInfo.Uuid,
		"location_type": locationInfo.LocationType.GetName(),
	}
	locationInfoList = append(locationInfoList, locationInfoMap)

	return locationInfoList
}

func flattenAvailableVersions(availableVersions []lcmEntityPkg.AvailableVersion) []map[string]interface{} {
	if len(availableVersions) == 0 {
		return nil
	}

	availableVersionsList := make([]map[string]interface{}, 0)

	for _, availableVersion := range availableVersions {
		availableVersionMap := map[string]interface{}{
			"version":                availableVersion.Version,
			"status":                 availableVersion.Status.GetName(),
			"is_enabled":             availableVersion.IsEnabled,
			"available_version_uuid": availableVersion.AvailableVersionUuid,
			"order":                  availableVersion.Order,
			"disablement_reason":     availableVersion.DisablementReason,
			"release_notes":          availableVersion.ReleaseNotes,
			"release_date":           availableVersion.ReleaseDate,
			"custom_message":         availableVersion.CustomMessage,
			"child_entities":         availableVersion.ChildEntities,
			"group_uuid":             availableVersion.GroupUuid,
			"dependencies":           flattenDependencies(availableVersion.Dependencies),
		}
		availableVersionsList = append(availableVersionsList, availableVersionMap)
	}

	return availableVersionsList
}

func flattenDependencies(dependencies []lcmEntityPkg.DependentEntity) []map[string]interface{} {
	if len(dependencies) == 0 {
		return nil
	}

	dependenciesList := make([]map[string]interface{}, 0)

	for _, dependency := range dependencies {
		dependencyMap := map[string]interface{}{
			"tenant_id":          dependency.TenantId,
			"ext_id":             dependency.ExtId,
			"links":              flattenLinks(dependency.Links),
			"entity_class":       dependency.EntityClass,
			"entity_model":       dependency.EntityModel,
			"entity_type":        dependency.EntityType,
			"entity_version":     dependency.EntityVersion,
			"hardware_family":    dependency.HardwareFamily,
			"dependent_versions": flattenKeyValuePairs(dependency.DependentVersions),
		}
		dependenciesList = append(dependenciesList, dependencyMap)
	}

	return dependenciesList
}

func flattenKeyValuePairs(dependentVersions []lcmConfigPkg.KVPair) []map[string]interface{} {
	if len(dependentVersions) == 0 {
		return nil
	}

	dependentVersionsList := make([]map[string]interface{}, 0)

	for _, dependentVersion := range dependentVersions {
		dependentVersionMap := map[string]interface{}{
			"name":  dependentVersion.Name,
			"value": flattenKVValue(dependentVersion.Value),
		}
		dependentVersionsList = append(dependentVersionsList, dependentVersionMap)
	}

	return dependentVersionsList
}

func flattenKVValue(value interface{}) []interface{} {
	valueMap := make(map[string]interface{})
	switch v := value.(type) {
	case string:
		valueMap["string"] = v
	case int:
		valueMap["integer"] = v
	case bool:
		valueMap["boolean"] = v
	case []string:
		valueMap["string_list"] = v
	case []int:
		valueMap["integer_list"] = v
	case map[string]string:
		valueMap["object"] = v

	case []lcmConfigPkg.MapOfStringWrapper:
		mapOfStrings := make([]interface{}, len(v))
		for i, m := range v {
			mapOfStrings[i] = m
		}

		valueMap["map_of_strings"] = mapOfStrings
	default:
		log.Printf("[WARN] Unknown type %T", v)
		return nil
	}
	return []interface{}{valueMap}
}

func flattenSubEntities(subEntities []common.EntityBaseModel) []map[string]interface{} {
	if len(subEntities) == 0 {
		return nil
	}

	subEntitiesList := make([]map[string]interface{}, 0)

	for _, subEntity := range subEntities {
		subEntityMap := map[string]interface{}{
			"tenant_id":       subEntity.TenantId,
			"ext_id":          subEntity.ExtId,
			"links":           flattenLinks(subEntity.Links),
			"entity_class":    subEntity.EntityClass,
			"entity_model":    subEntity.EntityModel,
			"entity_type":     flattenEntityTypes(subEntity.EntityType),
			"entity_version":  subEntity.EntityVersion,
			"hardware_family": subEntity.HardwareFamily,
		}
		subEntitiesList = append(subEntitiesList, subEntityMap)
	}

	return subEntitiesList
}

func flattenEntityTypes(entityType *common.EntityType) string {
	return entityType.GetName()
}

func flattenTime(inTime *time.Time) string {
	if inTime != nil {
		return inTime.UTC().Format(time.RFC3339)
	}
	return ""
}
