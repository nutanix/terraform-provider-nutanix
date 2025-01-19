package prismv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixListRestorablePcsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixListRestorablePcsV2Read,
		Schema: map[string]*schema.Schema{
			"restorable_source_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"page": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"restorable_pcs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixFetchPcV2(),
			},
		},
	}
}

func DatasourceNutanixListRestorablePcsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).PrismAPI
	var filter *string
	var page, limit *int

	if pagef, ok := d.GetOk("page"); ok {
		page = utils.IntPtr(pagef.(int))
	}
	if limitf, ok := d.GetOk("limit"); ok {
		limit = utils.IntPtr(limitf.(int))
	}
	if filterf, ok := d.GetOk("filter"); ok {
		filter = utils.StringPtr(filterf.(string))
	}

	restoreSourceExtID := d.Get("restorable_source_ext_id").(string)

	resp, err := conn.DomainManagerBackupsAPIInstance.ListRestorableDomainManagers(utils.StringPtr(restoreSourceExtID), page, limit, filter)
	if err != nil {
		return diag.Errorf("Error while Listing Restorable Domain Managers configurations Details: %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("restorable_pcs", []map[string]interface{}{}); err != nil {
			return diag.Errorf("Error setting Restorable pcs: %v", err)
		}
	}
	pcs := resp.Data.GetValue().([]config.DomainManager)
	if err := d.Set("restorable_pcs", flattenPcs(pcs)); err != nil {
		return diag.Errorf("Error setting pcs: %v", err)
	}

	d.SetId(utils.GenUUID())

	return nil
}
