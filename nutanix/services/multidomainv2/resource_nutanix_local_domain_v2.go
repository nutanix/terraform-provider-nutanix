package multidomainv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	multidomainCommon "github.com/nutanix/ntnx-api-golang-clients/multidomain-go-client/v4/models/multidomain/v4/common"
	"github.com/nutanix/ntnx-api-golang-clients/multidomain-go-client/v4/models/multidomain/v4/management"
	localdomainReq "github.com/nutanix/ntnx-api-golang-clients/multidomain-go-client/v4/models/multidomain/v4/request/localdomain"
	multidomainPrism "github.com/nutanix/ntnx-api-golang-clients/multidomain-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	commonUtils "github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixLocalDomainV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixLocalDomainV2Create,
		ReadContext:   ResourceNutanixLocalDomainV2Read,
		DeleteContext: ResourceNutanixLocalDomainV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"credentials": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_key": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
							// WriteOnly cannot combine with ForceNew.
							WriteOnly: true,
						},
						"key_id": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
							WriteOnly: true,
						},
					},
				},
			},
			"domain_ext_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"identity_provider_ext_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tenant_ext_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connection_status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemaLocalDomainConnectionStatus(),
			},
			"registration_status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemaLocalDomainRegistrationStatus(),
			},
			"location": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemaLocalDomainLocation(),
			},
		},
	}
}

func ResourceNutanixLocalDomainV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MultidomainAPI

	body := management.NewLocalDomainRegistrationSpec()
	if v, ok := d.GetOk("credentials"); ok {
		creds := v.([]interface{})
		if len(creds) > 0 && creds[0] != nil {
			if credMap, ok := creds[0].(map[string]interface{}); ok {
				if apiKey, ok := commonUtils.GetWriteOnlyString(d, "credentials.0.api_key"); ok {
					credMap["api_key"] = apiKey
				}
				if keyID, ok := commonUtils.GetWriteOnlyString(d, "credentials.0.key_id"); ok {
					credMap["key_id"] = keyID
				}
			}
		}
		body.Credentials = expandLocalDomainRegistrationCredentials(creds)
	}
	if v, ok := d.GetOk("domain_ext_id"); ok {
		body.DomainExtId = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("identity_provider_ext_id"); ok {
		body.IdentityProviderExtId = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("target_url"); ok {
		body.TargetUrl = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("tenant_ext_id"); ok {
		body.TenantExtId = utils.StringPtr(v.(string))
	}

	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] Register LocalDomain Body: %s", string(aJSON))

	registerReq := localdomainReq.RegisterLocalDomainRequest{
		Body: body,
	}
	resp, err := conn.LocalDomain.RegisterLocalDomain(ctx, &registerReq)
	if err != nil {
		return diag.Errorf("error registering LocalDomain: %v", err)
	}

	taskRef, ok := resp.Data.GetValue().(multidomainPrism.TaskReference)
	if !ok {
		return diag.Errorf("register local domain response did not contain task reference")
	}
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: commonUtils.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}
	if _, errWait := stateConf.WaitForStateContext(ctx); errWait != nil {
		return diag.Errorf("error waiting for local domain register task (%s): %s", utils.StringValue(taskUUID), errWait)
	}

	// LocalDomain is a singleton entity managed via a single Prism Central; the
	// domain ext_id (Prism Central identifier) is the natural key for the state.
	d.SetId(utils.StringValue(body.DomainExtId))

	return ResourceNutanixLocalDomainV2Read(ctx, d, meta)
}

func ResourceNutanixLocalDomainV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MultidomainAPI

	getReq := localdomainReq.GetLocalDomainRequest{}
	resp, err := conn.LocalDomain.GetLocalDomain(ctx, &getReq)
	if err != nil {
		return diag.Errorf("error while fetching LocalDomain: %v", err)
	}

	domain := resp.Data.GetValue().(management.LocalDomain)
	if err := d.Set("name", utils.StringValue(domain.Name)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("target_url", utils.StringValue(domain.TargetUrl)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("connection_status", flattenLocalDomainConnectionStatus(domain.ConnectionStatus)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("registration_status", flattenLocalDomainRegistrationStatus(domain.RegistrationStatus)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("location", flattenLocalDomainLocation(domain.Location)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ResourceNutanixLocalDomainV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MultidomainAPI

	unregisterReq := localdomainReq.UnregisterLocalDomainRequest{}
	resp, err := conn.LocalDomain.UnregisterLocalDomain(ctx, &unregisterReq)
	if err != nil {
		return diag.Errorf("error unregistering LocalDomain: %v", err)
	}

	taskRef, ok := resp.Data.GetValue().(multidomainPrism.TaskReference)
	if !ok {
		return diag.Errorf("unregister local domain response did not contain task reference")
	}
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: commonUtils.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}
	if _, errWait := stateConf.WaitForStateContext(ctx); errWait != nil {
		return diag.Errorf("error waiting for local domain unregister task (%s): %s", utils.StringValue(taskUUID), errWait)
	}

	d.SetId("")
	return nil
}

func expandLocalDomainRegistrationCredentials(in []interface{}) *multidomainCommon.RegistrationCredentials {
	if len(in) == 0 {
		return nil
	}
	m, ok := in[0].(map[string]interface{})
	if !ok {
		return nil
	}
	creds := multidomainCommon.NewRegistrationCredentials()
	if v, ok := m["api_key"].(string); ok && v != "" {
		creds.ApiKey = utils.StringPtr(v)
	}
	if v, ok := m["key_id"].(string); ok && v != "" {
		creds.KeyId = utils.StringPtr(v)
	}
	return creds
}
