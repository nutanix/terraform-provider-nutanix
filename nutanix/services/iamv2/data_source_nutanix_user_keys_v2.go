package iamv2

import (
	"context"
	"log"
	"time"

	import3 "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/iam/v4/authn"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
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
	var userExtID *string
	if v, ok := d.GetOk("user_ext_id"); ok {
		userExtID = utils.StringPtr(v.(string))
	}

	// initialize query params
	var filter, orderBy, selects *string
	var page, limit *int
	if pagef, ok := d.GetOk("page"); ok {
		page = utils.IntPtr(pagef.(int))
	} else {
		page = nil
	}
	if limitf, ok := d.GetOk("limit"); ok {
		limit = utils.IntPtr(limitf.(int))
	} else {
		limit = nil
	}
	if filterf, ok := d.GetOk("filter"); ok {
		filter = utils.StringPtr(filterf.(string))
	} else {
		filter = nil
	}
	if order, ok := d.GetOk("order_by"); ok {
		orderBy = utils.StringPtr(order.(string))
	} else {
		orderBy = nil
	}
	if selectf, ok := d.GetOk("select"); ok {
		selects = utils.StringPtr(selectf.(string))
	} else {
		selects = nil
	}
	resp, err := conn.UsersAPIInstance.ListUserKeys(userExtID, page, limit, filter, orderBy, selects)
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
