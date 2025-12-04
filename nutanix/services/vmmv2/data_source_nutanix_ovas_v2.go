package vmmv2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/content"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixOvasV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNutanixOvasV2Read,
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
			"ovas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixOvaV2(),
			},
		},
	}
}

func datasourceNutanixOvasV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

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
	resp, err := conn.OvasAPIInstance.ListOvas(page, limit, filter, orderBy, selects)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving OVA list: %w", err))
	}

	if resp.Data == nil {
		if err := d.Set("ovas", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No Data found",
			Detail:   "The API returned an empty list of OVA.",
		}}
	}

	getResp := resp.Data.GetValue().([]import1.Ova)

	if err := d.Set("ovas", flattenOvaEntities(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenOvaEntities(ovas []import1.Ova) []interface{} {
	if len(ovas) > 0 {
		ovaList := make([]interface{}, len(ovas))
		for k, v := range ovas {
			ova := make(map[string]interface{})
			if v.ExtId != nil {
				ova["ext_id"] = *v.ExtId
			}
			if v.Name != nil {
				ova["name"] = *v.Name
			}
			if v.Checksum != nil {
				ova["checksum"] = flattenOneOfOvaChecksum(v.Checksum)
			}
			if v.SizeBytes != nil {
				ova["size_bytes"] = int(*v.SizeBytes)
			}
			if v.Source != nil {
				ova["source"] = flattenOneOfOvaSource(v.Source)
			}
			if v.CreatedBy != nil {
				ova["created_by"] = flattenCreatedBy(v.CreatedBy)
			}
			if v.ClusterLocationExtIds != nil {
				ova["cluster_location_ext_ids"] = v.ClusterLocationExtIds
			}
			if v.ParentVm != nil {
				ova["parent_vm"] = *v.ParentVm
			}
			if v.VmConfig != nil {
				fields, diags := extractVMConfigFields(*v.VmConfig)
				if diags.HasError() {
					return nil
				}
				ova["vm_config"] = []interface{}{fields}
			}
			if v.DiskFormat != nil {
				ova["disk_format"] = flattenOvaDiskFormat(v.DiskFormat)
			}
			if v.CreateTime != nil {
				t := v.CreateTime
				ova["create_time"] = t.String()
			}
			if v.LastUpdateTime != nil {
				t := v.LastUpdateTime
				ova["last_update_time"] = t.String()
			}
			ovaList[k] = ova
		}
		return ovaList
	}
	return nil
}
