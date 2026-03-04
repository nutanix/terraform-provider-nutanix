package lcmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lcmEntityPkg "github.com/nutanix-core/ntnx-api-golang-sdk-internal/lifecycle-go-client/v17/models/lifecycle/v4/resources"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/lifecycle-go-client/v17/models/lifecycle/v4/request/entities"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixLcmEntitiesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixLcmEntitiesV2Create,
		Schema: map[string]*schema.Schema{
			"page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"select": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixLcmEntityV2(),
			},
		},
	}
}

func DatasourceNutanixLcmEntitiesV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LcmAPI

	listEntitiesRequest := import1.ListEntitiesRequest{}

	if v, ok := d.GetOk("page"); ok {
		listEntitiesRequest.Page_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("limit"); ok {
		listEntitiesRequest.Limit_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("filter"); ok {
		listEntitiesRequest.Filter_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("order_by"); ok {
		listEntitiesRequest.Orderby_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("select"); ok {
		listEntitiesRequest.Select_ = utils.StringPtr(v.(string))
	}

	resp, err := conn.LcmEntitiesAPIInstance.ListEntities(ctx, &listEntitiesRequest)
	if err != nil {
		return diag.Errorf("error while listing the Lcm entities : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("entities", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of LCM entities.",
		}}
	}

	entities := resp.Data.GetValue().([]lcmEntityPkg.Entity)
	if err := d.Set("entities", flattenLcmEntities(entities)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.GenUUID())
	return nil
}

func flattenLcmEntities(entities []lcmEntityPkg.Entity) []map[string]interface{} {
	if len(entities) == 0 {
		return nil
	}

	flattenedEntities := make([]map[string]interface{}, 0)

	for _, entity := range entities {
		flattenedEntity := map[string]interface{}{
			"tenant_id":          entity.TenantId,
			"ext_id":             entity.ExtId,
			"links":              flattenLinks(entity.Links),
			"entity_class":       entity.EntityClass,
			"entity_model":       entity.EntityModel,
			"entity_type":        flattenEntityTypes(entity.EntityType),
			"entity_version":     entity.EntityVersion,
			"hardware_family":    entity.HardwareFamily,
			"entity_description": entity.EntityDescription,
			"location_info":      flattenLocationInfo(entity.LocationInfo),
			"target_version":     entity.TargetVersion,
			"last_updated_time":  flattenTime(entity.LastUpdatedTime),
			"device_id":          entity.DeviceId,
			"group_uuid":         entity.GroupUuid,
			"entity_details":     flattenKeyValuePairs(entity.EntityDetails),
			"child_entities":     entity.ChildEntities,
			"available_versions": flattenAvailableVersions(entity.AvailableVersions),
			"sub_entities":       flattenSubEntities(entity.SubEntities),
			"cluster_ext_id":     entity.ClusterExtId,
			"hardware_vendor":    entity.HardwareVendor,
		}
		flattenedEntities = append(flattenedEntities, flattenedEntity)
	}

	return flattenedEntities
}
