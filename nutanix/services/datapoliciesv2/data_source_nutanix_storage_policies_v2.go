package datapoliciesv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/datapolicies-go-client/v4/models/datapolicies/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DataSourceNutanixStoragePoliciesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixStoragePoliciesV2Read,
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
			"storage_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DataSourceNutanixStoragePolicyV2(),
			},
			"total_available_results": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceNutanixStoragePoliciesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).DataPoliciesAPI

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
	resp, err := conn.StoragePolicies.ListStoragePolicies(page, limit, filter, orderBy, selects)
	if err != nil {
		return diag.Errorf("error while fetching storage policies: %v", err)
	}
	if resp.Data == nil {
		if err := d.Set("storage_policies", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(utils.GenUUID())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of storage policies.",
		}}
	}
	getResp := resp.Data.GetValue().([]import1.StoragePolicy)
	if err := d.Set("storage_policies", flattenStoragePolicies(getResp)); err != nil {
		return diag.FromErr(err)
	}
	if resp.Metadata != nil && resp.Metadata.TotalAvailableResults != nil {
		if err := d.Set("total_available_results", *resp.Metadata.TotalAvailableResults); err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(utils.GenUUID())
	return nil
}

func flattenStoragePolicies(storagePolicies []import1.StoragePolicy) []map[string]interface{} {
	if len(storagePolicies) == 0 {
		return []map[string]interface{}{}
	}
	storagePoliciesList := make([]map[string]interface{}, 0)
	for _, storagePolicy := range storagePolicies {
		storagePoliciesList = append(storagePoliciesList, flattenStoragePolicy(storagePolicy))
	}

	return storagePoliciesList
}

func flattenStoragePolicy(storagePolicy import1.StoragePolicy) map[string]interface{} {
	return map[string]interface{}{
		"tenant_id":            utils.StringValue(storagePolicy.TenantId),
		"ext_id":               utils.StringValue(storagePolicy.ExtId),
		"links":                flattenLinks(storagePolicy.Links),
		"name":                 utils.StringValue(storagePolicy.Name),
		"category_ext_ids":     storagePolicy.CategoryExtIds,
		"compression_spec":     flattenCompressionSpec(storagePolicy.CompressionSpec),
		"encryption_spec":      flattenEncryptionSpec(storagePolicy.EncryptionSpec),
		"qos_spec":             flattenQosSpec(storagePolicy.QosSpec),
		"fault_tolerance_spec": flattenFaultToleranceSpec(storagePolicy.FaultToleranceSpec),
		"policy_type":          storagePolicy.PolicyType.GetName(),
	}
}
