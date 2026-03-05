package objectstoresv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/objects-go-client/v17/models/objects/v4/config"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/objects-go-client/v17/models/objects/v4/request/objectstores"
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

	objectStoreExtID := d.Get("object_store_ext_id").(string)

	listCertificatesRequest := import1.ListCertificatesByObjectstoreIdRequest{
		ObjectStoreExtId: utils.StringPtr(objectStoreExtID),
	}

	if v, ok := d.GetOk("page"); ok {
		listCertificatesRequest.Page_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("limit"); ok {
		listCertificatesRequest.Limit_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("filter"); ok {
		listCertificatesRequest.Filter_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("select"); ok {
		listCertificatesRequest.Select_ = utils.StringPtr(v.(string))
	}

	listResp, err := conn.ObjectStoresAPIInstance.ListCertificatesByObjectstoreId(ctx, &listCertificatesRequest)
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
