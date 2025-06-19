package objectstoresv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/objects-go-client/v4/models/objects/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixObjectStoreCertificatesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixObjectStoreCertificatesV2Read,
		Schema: map[string]*schema.Schema{
			"object_store_ext_id": {
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
			"select": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"certificates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixObjectStoreCertificateV2(),
			},
		},
	}
}

func DatasourceNutanixObjectStoreCertificatesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ObjectStoreAPI

	// initialize query params
	var filter, selects *string
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
	if selectf, ok := d.GetOk("select"); ok {
		selects = utils.StringPtr(selectf.(string))
	} else {
		selects = nil
	}

	// get object store ext id
	objectStoreExtID := d.Get("object_store_ext_id").(string)

	// list certificates
	listResp, err := conn.ObjectStoresAPIInstance.ListCertificatesByObjectstoreId(utils.StringPtr(objectStoreExtID), page, limit, filter, selects)
	if err != nil {
		return diag.Errorf("error while fetching object stores : %v", err)
	}

	if listResp.Data == nil {
		if err := d.Set("certificates", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}
	} else {
		certificates := listResp.Data.GetValue().([]config.Certificate)

		if err := d.Set("certificates", flattenCertificates(certificates)); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(utils.GenUUID())

	return nil
}

func flattenCertificates(certificates []config.Certificate) interface{} {
	if len(certificates) == 0 {
		return nil
	}
	certificatesList := make([]map[string]interface{}, 0, len(certificates))
	for _, certificate := range certificates {
		certificateMap := make(map[string]interface{})
		certificateMap["tenant_id"] = certificate.TenantId
		certificateMap["ext_id"] = certificate.ExtId
		certificateMap["links"] = flattenLinks(certificate.Links)
		certificateMap["metadata"] = flattenMetadata(certificate.Metadata)
		certificateMap["alternate_fqdns"] = flattenFQDNs(certificate.AlternateFqdns)
		certificateMap["alternate_ips"] = flattenIPAddress(certificate.AlternateIps)
		certificatesList = append(certificatesList, certificateMap)
	}
	return certificatesList
}
