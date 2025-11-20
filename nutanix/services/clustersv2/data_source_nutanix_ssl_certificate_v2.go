package clustersv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixSSLCertificateV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixSSLCertificateV2Read,
		Schema: map[string]*schema.Schema{
			"cluster_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"passphrase": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"public_certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ca_chain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_key_algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DatasourceNutanixSSLCertificateV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI

	clusterExtID := d.Get("cluster_ext_id").(string)

	resp, err := conn.SSLCertificateAPI.GetSSLCertificate(utils.StringPtr(clusterExtID))
	if err != nil {
		return diag.Errorf("error while fetching SSL certificate: %v", err)
	}

	if resp.Data == nil {
		d.SetId(resource.UniqueId())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No SSL certificate found.",
			Detail:   "The API returned no SSL certificate data.",
		}}
	}

	value := resp.Data.GetValue()

	// Log the response type for debugging
	log.Printf("[DEBUG] SSL certificate response type: %T, value: %v", value, value)

	// Try to cast to SSLCertificate
	sslCert, ok := value.(import1.SSLCertificate)
	if !ok {
		// If it's a string, handle it as a simple certificate
		if certStr, ok := value.(string); ok {
			if err := d.Set("certificate", certStr); err != nil {
				return diag.FromErr(err)
			}
			d.SetId(resource.UniqueId())
			return nil
		}
		return diag.Errorf("unexpected response type: expected import1.SSLCertificate or string, got %T", value)
	}

	aJSON, _ := json.MarshalIndent(sslCert, "", "  ")
	log.Printf("[DEBUG] SSLCertificate struct received: %s", string(aJSON))

	if err := d.Set("passphrase", sslCert.Passphrase); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("private_key", sslCert.PrivateKey); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("public_certificate", sslCert.PublicCertificate); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ca_chain", sslCert.CaChain); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("private_key_algorithm", sslCert.PrivateKeyAlgorithm.GetName()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}
