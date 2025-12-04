package objectstoresv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	objectsCommon "github.com/nutanix/ntnx-api-golang-clients/objects-go-client/v4/models/common/v1/config"
	"github.com/nutanix/ntnx-api-golang-clients/objects-go-client/v4/models/objects/v4/config"
	objectPrismConfig "github.com/nutanix/ntnx-api-golang-clients/objects-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixObjectStoreCertificateV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixObjectStoreCertificateV2Create,
		ReadContext:   ResourceNutanixObjectStoreCertificateV2Read,
		UpdateContext: ResourceNutanixObjectStoreCertificateV2Update,
		DeleteContext: ResourceNutanixObjectStoreCertificateV2Delete,
		Schema: map[string]*schema.Schema{
			"object_store_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"path": {
				Type:     schema.TypeString,
				Required: true,
			},
			// computed attributes
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
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": linksSchema(),
			"metadata": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: metadataSchema(),
				},
			},
		},
	}
}

func ResourceNutanixObjectStoreCertificateV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ObjectStoreAPI

	objectStoreExtID := d.Get("object_store_ext_id").(string)

	readResp, err := conn.ObjectStoresAPIInstance.GetObjectstoreById(utils.StringPtr(objectStoreExtID))
	if err != nil {
		return diag.Errorf("Error reading object store: %s", err)
	}

	// Extract E-Tag Header
	args := make(map[string]interface{})
	etagValue := conn.ObjectStoresAPIInstance.ApiClient.GetEtag(readResp)
	args["If-Match"] = utils.StringPtr(etagValue)

	filePath := d.Get("path").(string)

	resp, err := conn.ObjectStoresAPIInstance.CreateCertificate(utils.StringPtr(objectStoreExtID), utils.StringPtr(filePath), args)
	if err != nil {
		return diag.Errorf("Error creating object store certificate: %s", err)
	}

	TaskRef := resp.Data.GetValue().(objectPrismConfig.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the cluster to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for object store certificate to be created : %s", err)
	}

	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while create object store certificate task: %s", err)
	}

	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Object store certificate task details: %s", string(aJSON))

	// Get created object store extID from TASK API
	objectStoreCertificateExtID := taskDetails.EntitiesAffected[1].ExtId
	d.SetId(*objectStoreCertificateExtID)

	return ResourceNutanixObjectStoreCertificateV2Read(ctx, d, meta)
}

func ResourceNutanixObjectStoreCertificateV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ObjectStoreAPI

	objectStoreExtID := d.Get("object_store_ext_id").(string)

	resp, err := conn.ObjectStoresAPIInstance.GetCertificateById(utils.StringPtr(objectStoreExtID), utils.StringPtr(d.Id()))
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

	return nil
}

func ResourceNutanixObjectStoreCertificateV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return ResourceNutanixObjectStoreCertificateV2Create(ctx, d, meta)
}

func ResourceNutanixObjectStoreCertificateV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

// expandFQDNs expands the FQDNs from the schema into the API model
func expandFqdn(fqdnList interface{}) []objectsCommon.FQDN {
	if fqdnList == nil {
		return nil
	}
	fqdnListSchema := fqdnList.([]interface{})
	fqdnListAPI := make([]objectsCommon.FQDN, 0, len(fqdnListSchema))
	for _, fqdn := range fqdnListSchema {
		fqdnMap := fqdn.(map[string]interface{})
		fqdnAPI := objectsCommon.FQDN{
			Value: utils.StringPtr(fqdnMap["value"].(string)),
		}
		fqdnListAPI = append(fqdnListAPI, fqdnAPI)
	}
	return fqdnListAPI
}

// flattenFQDNs flattens the FQDNs from the API model into the schema
func flattenFQDNs(fQDN []objectsCommon.FQDN) []map[string]interface{} {
	if len(fQDN) == 0 {
		return nil
	}
	fqdnList := make([]map[string]interface{}, 0, len(fQDN))
	for _, fqdn := range fQDN {
		fqdnList = append(fqdnList, map[string]interface{}{
			"value": utils.StringValue(fqdn.Value),
		})
	}
	return fqdnList
}
