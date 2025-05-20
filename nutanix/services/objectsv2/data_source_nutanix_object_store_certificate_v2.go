package objectstoresv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/objects-go-client/v4/models/objects/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixObjectStoreCertificateV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixObjectStoreCertificateV2Read,
		Schema: map[string]*schema.Schema{
			"object_store_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rel": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"metadata": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: metadataSchema(),
				},
			},
			"alternate_fqdns": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"alternate_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4": SchemaForValuePrefixLength(ipv4PrefixLengthDefaultValue),
						"ipv6": SchemaForValuePrefixLength(ipv6PrefixLengthDefaultValue),
					},
				},
			},
		},
	}
}

func DatasourceNutanixObjectStoreCertificateV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ObjectStoreAPI

	objectStoreExtID := d.Get("object_store_ext_id").(string)
	certificateExtID := d.Get("ext_id").(string)

	resp, err := conn.ObjectStoresAPIInstance.GetCertificateById(utils.StringPtr(objectStoreExtID), utils.StringPtr(certificateExtID))
	if err != nil {
		return diag.Errorf("Error reading object store certificate : %s", err)
	}

	certificate := resp.Data.GetValue().(config.Certificate)

	if err := d.Set("tenant_id", certificate.TenantId); err != nil {
		return diag.Errorf("Error setting tenant_id: %s", err)
	}
	if err := d.Set("ext_id", certificate.ExtId); err != nil {
		return diag.Errorf("Error setting ext_id: %s", err)
	}
	if err := d.Set("links", flattenLinks(certificate.Links)); err != nil {
		return diag.Errorf("Error setting links: %s", err)
	}
	if err := d.Set("metadata", flattenMetadata(certificate.Metadata)); err != nil {
		return diag.Errorf("Error setting metadata: %s", err)
	}
	if err := d.Set("alternate_fqdns", flattenFQDNs(certificate.AlternateFqdns)); err != nil {
		return diag.Errorf("Error setting alternate_fqdns: %s", err)
	}
	if err := d.Set("alternate_ips", flattenIPAddress(certificate.AlternateIps)); err != nil {
		return diag.Errorf("Error setting alternate_ips: %s", err)
	}

	d.SetId(utils.StringValue(certificate.ExtId))
	return nil
}
