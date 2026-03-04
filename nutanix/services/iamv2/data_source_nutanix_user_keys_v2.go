package iamv2

import (
	"context"
	"log"
	"time"

	import3 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/iam-go-client/v17/models/iam/v4/authn"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/iam-go-client/v17/models/iam/v4/request/users"
)

func DatasourceNutanixUserKeysV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceNutanixUserKeysV2Read,
		Schema: map[string]*schema.Schema{
			"user_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
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
			"keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixUserKeyV2(),
			},
		},
	}
}

func DataSourceNutanixUserKeysV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection
	conn := meta.(*conns.Client).IamAPI

	listUserKeysRequest := import1.ListUserKeysRequest{}
	if v, ok := d.GetOk("page"); ok {
		listUserKeysRequest.Page_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("limit"); ok {
		listUserKeysRequest.Limit_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("filter"); ok {
		listUserKeysRequest.Filter_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("order_by"); ok {
		listUserKeysRequest.Orderby_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("select"); ok {
		listUserKeysRequest.Select_ = utils.StringPtr(v.(string))
	}
	resp, err := conn.UsersAPIInstance.ListUserKeys(ctx, &listUserKeysRequest)
	if err != nil {
		return diag.Errorf("error while fetching the user keys: %v", err)
	}
	if resp.Data != nil {
		getResp := resp.Data.GetValue().([]import3.Key)
		log.Printf("getResp: %v", getResp)
		if err := d.Set("keys", flattenKeysEntities(getResp)); err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(utils.GenUUID())
	return nil
}

func flattenKeysEntities(data []import3.Key) []map[string]interface{} {
	flattened := make([]map[string]interface{}, 0)
	for _, item := range data {
		entry := map[string]interface{}{
			"tenant_id":         item.TenantId,
			"ext_id":            item.ExtId,
			"name":              item.Name,
			"description":       item.Description,
			"key_type":          item.KeyType.GetName(),
			"created_time":      flattenTime(item.CreatedTime),
			"last_updated_by":   item.LastUpdatedBy,
			"creation_type":     item.CreationType.GetName(),
			"expiry_time":       flattenTime(item.ExpiryTime),
			"status":            item.Status.GetName(),
			"created_by":        item.CreatedBy,
			"last_updated_time": flattenTime(item.LastUpdatedTime),
			"assigned_to":       item.AssignedTo,
			"last_used_time":    flattenTime(item.LastUsedTime),
			"key_details":       flattenKeyDetails(item.KeyDetails),
		}
		flattened = append(flattened, entry)
	}
	return flattened
}

func flattenTime(inTime *time.Time) string {
	if inTime != nil {
		return inTime.UTC().Format(time.RFC3339)
	}
	return ""
}

func SchemaForLinks() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"href": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"rel": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}
