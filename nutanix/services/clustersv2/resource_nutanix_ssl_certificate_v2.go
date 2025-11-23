package clustersv2

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	clustermgmtPrism "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixSSLCertificateV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixSSLCertificateV2Create,
		ReadContext:   ResourceNutanixSSLCertificateV2Read,
		UpdateContext: ResourceNutanixSSLCertificateV2Update,
		DeleteContext: ResourceNutanixSSLCertificateV2Delete,
		Schema: map[string]*schema.Schema{
			"cluster_ext_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"passphrase": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"private_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"public_certificate": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ca_chain": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"private_key_algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(PrivateKeyAlgorithmStrings, false),
			},
		},
	}
}

func ResourceNutanixSSLCertificateV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI

	clusterExtID := d.Get("cluster_ext_id").(string)

	// Extract the etag header
	resp, err := conn.SSLCertificateAPI.GetSSLCertificate(utils.StringPtr(clusterExtID))
	if err != nil {
		return diag.Errorf("error while fetching SSL certificate: %v", err)
	}
	etag := conn.SSLCertificateAPI.ApiClient.GetEtag(resp)
	args := make(map[string]interface{})
	args["If-Match"] = utils.StringPtr(etag)

	// Build the update payload
	updateSpec := config.NewSSLCertificate()

	if value, ok := utils.IsStringSetAndNotEmpty(d.Get("passphrase")); ok {
		updateSpec.Passphrase = utils.StringPtr(value)
	}
	if value, ok := utils.IsStringSetAndNotEmpty(d.Get("private_key")); ok {
		updateSpec.PrivateKey = utils.StringPtr(value)
	}
	if value, ok := utils.IsStringSetAndNotEmpty(d.Get("public_certificate")); ok {
		updateSpec.PublicCertificate = utils.StringPtr(value)
	}
	if value, ok := utils.IsStringSetAndNotEmpty(d.Get("ca_chain")); ok {
		updateSpec.CaChain = utils.StringPtr(value)
	}
	if value, ok := utils.IsStringSetAndNotEmpty(d.Get("private_key_algorithm")); ok {
		updateSpec.PrivateKeyAlgorithm = common.ExpandEnum(value, PrivateKeyAlgorithmMap, "private_key_algorithm")
	}

	// Log the update payload
	aJSON, _ := json.MarshalIndent(updateSpec, "", "  ")
	log.Printf("[DEBUG] SSL certificate update payload: %s", string(aJSON))

	// Call the update API
	updateResp, err := conn.SSLCertificateAPI.UpdateSSLCertificate(utils.StringPtr(clusterExtID), updateSpec, args)
	if err != nil {
		return diag.Errorf("error while updating SSL certificate: %v", err)
	}

	TaskRef := updateResp.Data.GetValue().(clustermgmtPrism.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task
	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the SSL Certificate to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for SSL certificate (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	returnResourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching SSL certificate task: %v", err)
	}
	taskDetails := returnResourceUUID.Data.GetValue().(prismConfig.Task)
	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] SSL certificate update task details: %s", string(aJSON))

	d.SetId(resource.UniqueId())

	return ResourceNutanixSSLCertificateV2Read(ctx, d, meta)
}

func ResourceNutanixSSLCertificateV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI

	clusterExtID := d.Get("cluster_ext_id").(string)

	// Retry logic to handle temporary API unavailability after certificate regeneration
	var resp *import1.GetSSLCertificateApiResponse
	var err error
	maxRetries := 10
	retryDelay := 2 * time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, err = conn.SSLCertificateAPI.GetSSLCertificate(utils.StringPtr(clusterExtID))
		if err == nil {
			break
		}

		if attempt < maxRetries-1 {
			log.Printf("[DEBUG] Attempt %d/%d failed to fetch SSL certificate, retrying in %v: %v", attempt+1, maxRetries, retryDelay, err)
			select {
			case <-ctx.Done():
				return diag.Errorf("context cancelled while fetching SSL certificate: %v", ctx.Err())
			case <-time.After(retryDelay):
				retryDelay = time.Duration(float64(retryDelay) * 1.5) // Exponential backoff
			}
		}
	}

	if err != nil {
		return diag.Errorf("error while fetching SSL certificate after %d attempts: %v", maxRetries, err)
	}

	if resp.Data == nil {
		d.SetId("")
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No SSL certificate found.",
			Detail:   "The API returned no SSL certificate data.",
		}}
	}

	value := resp.Data.GetValue()
	sslCert, ok := value.(config.SSLCertificate)
	if !ok {
		return diag.Errorf("unexpected response type: expected config.SSLCertificate, got %T", value)
	}

	// The API returns public certificate and private key algorithm only
	// Only set public_certificate from API if it wasn't explicitly set in config
	// This prevents unnecessary diffs when user provides base64-encoded value
	// but API returns formatted certificate
	if _, ok := d.GetOk("public_certificate"); !ok {
		if err := d.Set("public_certificate", sslCert.PublicCertificate); err != nil {
			return diag.FromErr(err)
		}
	}
	if sslCert.PrivateKeyAlgorithm != nil {
		if err := d.Set("private_key_algorithm", sslCert.PrivateKeyAlgorithm.GetName()); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func ResourceNutanixSSLCertificateV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return ResourceNutanixSSLCertificateV2Create(ctx, d, meta)
}

func ResourceNutanixSSLCertificateV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
