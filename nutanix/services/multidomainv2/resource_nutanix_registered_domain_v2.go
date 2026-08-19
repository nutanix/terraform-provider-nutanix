package multidomainv2

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import7 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/models/multidomain/v4/common"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/models/multidomain/v4/config"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/models/multidomain/v4/request/registereddomains"
	multidomainPrism "github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	commonUtils "github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/multidomain"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixRegisteredDomainV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixRegisteredDomainV2Create,
		ReadContext:   ResourceNutanixRegisteredDomainV2Read,
		UpdateContext: ResourceNutanixRegisteredDomainV2Update,
		DeleteContext: ResourceNutanixRegisteredDomainV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dns_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"fqdn": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"domain_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"location_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"credential_issuer": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"api_credential_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"VALID", "NEAR_EXPIRY", "EXPIRED",
				}, false),
			},
			"registration_state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"UNREGISTERED", "REGISTERING", "REGISTERED", "UNREGISTERING",
					"REGISTRATION_ERROR", "UNREGISTRATION_ERROR", "DRAFT",
				}, false),
			},
			"registered_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"platform_data": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     schemaRegisteredDomainPlatformData(),
			},
			"registration_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     schemaRegisteredDomainRegistrationConfig(),
			},
			"tunnel": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     schemaRegisteredDomainTunnel(),
			},
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": schemaForLinks(),
		},
	}
}

func schemaRegisteredDomainPlatformData() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cluster_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"connectivity_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"pc_env_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"pc_provider_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"pc_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func schemaRegisteredDomainRegistrationConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"domain_ext_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"target_url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tenant_ext_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"credentials": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_key": {
							Type:      schema.TypeString,
							Optional:  true,
							Computed:  true,
							Sensitive: true,
						},
						"key_id": {
							Type:      schema.TypeString,
							Optional:  true,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},
		},
	}
}

func schemaRegisteredDomainTunnel() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"provider": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"IKAT",
				}, false),
			},
		},
	}
}

func ResourceNutanixRegisteredDomainV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NCDomainManagementAPI

	body := expandRegisteredDomain(d)

	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] Create RegisteredDomain Body: %s", string(aJSON))

	createReq := import1.CreateRegisteredDomainRequest{
		Body: body,
	}
	resp, err := conn.RegisteredDomains.CreateRegisteredDomain(ctx, &createReq)
	if err != nil {
		return diag.Errorf("error creating RegisteredDomain: %v", err)
	}

	taskRef, ok := resp.Data.GetValue().(multidomainPrism.TaskReference)
	if !ok {
		return diag.Errorf("create registered domain response did not contain task reference")
	}
	taskUUID := taskRef.ExtId
	log.Printf("[DEBUG] Task UUID: %s", utils.StringValue(taskUUID))

	// Wait for the task to complete
	// Nutanix Central Task API is not yet GA'ed/ not avilable in SDK. So depending on the list api to check if the domain registration is completed.
	listReq := import1.ListRegisteredDomainsRequest{
		Filter_: utils.StringPtr(fmt.Sprintf("name eq '%s'", utils.StringValue(body.Name))),
		Limit_:  utils.IntPtr(1),
	}
	domainName := utils.StringValue(body.Name)
	timeout := time.NewTimer(3 * time.Minute)
	defer timeout.Stop()
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		log.Printf("[DEBUG] Polling registered domains for name %q", domainName)
		listResp, err := conn.RegisteredDomains.ListRegisteredDomains(ctx, &listReq)
		if err != nil {
			return diag.Errorf("error while listing registered domains: %v", err)
		}

		if listResp.Data != nil {
			registeredDomains, ok := listResp.Data.GetValue().([]config.RegisteredDomain)
			if !ok {
				return diag.Errorf("unexpected registered domains list response")
			}
			for _, registeredDomain := range registeredDomains {
				if utils.StringValue(registeredDomain.Name) == domainName {
					d.SetId(utils.StringValue(registeredDomain.ExtId))
					return ResourceNutanixRegisteredDomainV2Read(ctx, d, meta)
				}
			}
		}

		select {
		case <-ctx.Done():
			return diag.FromErr(ctx.Err())
		case <-timeout.C:
			return diag.Errorf("registered domain %q did not appear within 3 minutes", domainName)
		case <-ticker.C:
		}
	}
}

func ResourceNutanixRegisteredDomainV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NCDomainManagementAPI

	extID := d.Id()
	getReq := import1.GetRegisteredDomainByIdRequest{
		ExtId: utils.StringPtr(extID),
	}
	resp, err := conn.RegisteredDomains.GetRegisteredDomainById(ctx, &getReq)
	if err != nil {
		return diag.Errorf("error while fetching registered domain: %v", err)
	}

	rd := resp.Data.GetValue().(config.RegisteredDomain)
	return diag.FromErr(setRegisteredDomainToResourceData(d, &rd))
}

func ResourceNutanixRegisteredDomainV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NCDomainManagementAPI
	if d.HasChange("domain_ext_id") {
		return diag.Errorf("domain_ext_id is not allowed to be updated")
	}

	extID := d.Id()
	getReq := import1.GetRegisteredDomainByIdRequest{
		ExtId: utils.StringPtr(extID),
	}
	resp, err := conn.RegisteredDomains.GetRegisteredDomainById(ctx, &getReq)
	if err != nil {
		return diag.Errorf("error while fetching registered domain: %v", err)
	}
	getSpec := resp.Data.GetValue().(config.RegisteredDomain)
	updateSpecNew := config.NewRegisteredDomain()
	if d.HasChange("name") {
		updateSpecNew.Name = utils.StringPtr(d.Get("name").(string))
	} else {
		updateSpecNew.Name = getSpec.Name
	}
	if d.HasChange("location_ext_id") {
		updateSpecNew.LocationExtId = utils.StringPtr(d.Get("location_ext_id").(string))
	} else {
		updateSpecNew.LocationExtId = getSpec.LocationExtId
	}

	// Domain ExtId is same as ExtId.
	updateSpecNew.DomainExtId = getSpec.ExtId

	getRespBeforeUpdate, err := conn.RegisteredDomains.GetRegisteredDomainById(ctx, &getReq)
	if err != nil {
		return diag.Errorf("error while fetching registered domain: %v", err)
	}
	args := make(map[string]interface{})
	etagValue := conn.APIClientInstance.GetEtag(getRespBeforeUpdate)
	args["If-Match"] = utils.StringPtr(etagValue)

	updateReq := import1.UpdateRegisteredDomainByIdRequest{
		ExtId: utils.StringPtr(extID),
		Body:  updateSpecNew,
	}

	aJSON, _ := json.MarshalIndent(updateSpecNew, "", "  ")
	log.Printf("[DEBUG] Update RegisteredDomain Body: %s", string(aJSON))
	updateResp, err := conn.RegisteredDomains.UpdateRegisteredDomainById(ctx, &updateReq, args)
	if err != nil {
		return diag.Errorf("error updating RegisteredDomain: %v", err)
	}

	taskRef, ok := updateResp.Data.GetValue().(multidomainPrism.TaskReference)
	if !ok {
		return diag.Errorf("update registered domain response did not contain task reference")
	}
	taskUUID := taskRef.ExtId
	log.Printf("[DEBUG] Task UUID: %s", utils.StringValue(taskUUID))

	if err := waitForRegisteredDomainTask(ctx, conn, utils.StringValue(taskUUID), "update"); err != nil {
		return diag.FromErr(err)
	}

	return ResourceNutanixRegisteredDomainV2Read(ctx, d, meta)
}

func ResourceNutanixRegisteredDomainV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NCDomainManagementAPI

	extID := d.Id()
	getReq := import1.GetRegisteredDomainByIdRequest{
		ExtId: utils.StringPtr(extID),
	}
	getResp, err := conn.RegisteredDomains.GetRegisteredDomainById(ctx, &getReq)
	if err != nil {
		return diag.Errorf("error fetching RegisteredDomain before delete: %v", err)
	}
	args := make(map[string]interface{})
	etagValue := conn.APIClientInstance.GetEtag(getResp)
	args["If-Match"] = utils.StringPtr(etagValue)

	deleteReq := import1.DeleteRegisteredDomainByIdRequest{
		ExtId: utils.StringPtr(extID),
	}
	resp, err := conn.RegisteredDomains.DeleteRegisteredDomainById(ctx, &deleteReq, args)
	if err != nil {
		return diag.Errorf("error deleting RegisteredDomain: %v", err)
	}

	taskRef, ok := resp.Data.GetValue().(multidomainPrism.TaskReference)
	if !ok {
		return diag.Errorf("delete registered domain response did not contain task reference")
	}
	taskUUID := taskRef.ExtId
	log.Printf("[DEBUG] Task UUID: %s", utils.StringValue(taskUUID))

	if err := waitForRegisteredDomainTask(ctx, conn, utils.StringValue(taskUUID), "delete"); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func waitForRegisteredDomainTask(
	ctx context.Context,
	conn *multidomain.NCDomainManagementClient,
	taskUUID string,
	operation string,
) error {
	log.Printf("[DEBUG] Waiting for registered domain task %s", taskUUID)
	taskStatus, err := conn.Tasks.WaitForTerminalState(ctx, taskUUID, 10*time.Second)
	if err != nil {
		return fmt.Errorf("waiting for registered domain %s task %q: %w", operation, taskUUID, err)
	}
	if !strings.EqualFold(taskStatus.State(), "SUCCEEDED") {
		return fmt.Errorf(
			"registered domain %s task %q ended with status %s: %s",
			operation,
			taskUUID,
			taskStatus.State(),
			taskStatus.Data.SubStateInfo.ErrorDetail,
		)
	}
	return nil
}

func expandRegisteredDomain(d *schema.ResourceData) *config.RegisteredDomain {
	body := config.NewRegisteredDomain()
	if v, ok := d.GetOk("name"); ok {
		body.Name = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("dns_name"); ok {
		body.DnsName = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("fqdn"); ok {
		body.Fqdn = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("domain_ext_id"); ok {
		body.DomainExtId = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("location_ext_id"); ok {
		body.LocationExtId = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("credential_issuer"); ok {
		body.CredentialIssuer = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("created_by"); ok {
		body.CreatedBy = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("api_credential_status"); ok {
		body.ApiCredentialStatus = commonUtils.ExpandEnum[config.ApiCredentialStatus](v.(string))
	}
	if v, ok := d.GetOk("registration_state"); ok {
		body.RegistrationState = commonUtils.ExpandEnum[import7.RegistrationState](v.(string))
	}
	if v, ok := d.GetOk("registered_time"); ok {
		body.RegisteredTime = expandRegisteredDomainTime(v.(string))
	}
	if v, ok := d.GetOk("platform_data"); ok {
		body.PlatformData = expandRegisteredDomainPlatformData(v.([]interface{}))
	}
	if v, ok := d.GetOk("registration_config"); ok {
		body.RegistrationConfig = expandRegisteredDomainRegistrationConfig(v.([]interface{}))
	}
	if v, ok := d.GetOk("tunnel"); ok {
		body.Tunnel = expandRegisteredDomainTunnel(v.([]interface{}))
	}
	return body
}

func expandRegisteredDomainTime(in string) *time.Time {
	if in == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, in)
	if err != nil {
		log.Printf("[WARN] unable to parse registered_time %q: %v", in, err)
		return nil
	}
	return &t
}

func expandRegisteredDomainPlatformData(in []interface{}) *config.PlatformData {
	if len(in) == 0 || in[0] == nil {
		return nil
	}
	m, ok := in[0].(map[string]interface{})
	if !ok {
		return nil
	}
	out := config.NewPlatformData()
	if v, ok := m["cluster_count"].(int); ok && v != 0 {
		out.ClusterCount = utils.IntPtr(v)
	}
	if v, ok := m["connectivity_status"].(string); ok && v != "" {
		out.ConnectivityStatus = utils.StringPtr(v)
	}
	if v, ok := m["pc_env_type"].(string); ok && v != "" {
		out.PcEnvType = utils.StringPtr(v)
	}
	if v, ok := m["pc_provider_type"].(string); ok && v != "" {
		out.PcProviderType = utils.StringPtr(v)
	}
	if v, ok := m["pc_version"].(string); ok && v != "" {
		out.PcVersion = utils.StringPtr(v)
	}
	return out
}

func expandRegisteredDomainRegistrationConfig(in []interface{}) *config.RegistrationConfig {
	if len(in) == 0 || in[0] == nil {
		return nil
	}
	m, ok := in[0].(map[string]interface{})
	if !ok {
		return nil
	}
	out := config.NewRegistrationConfig()
	if v, ok := m["domain_ext_id"].(string); ok && v != "" {
		out.DomainExtId = utils.StringPtr(v)
	}
	if v, ok := m["target_url"].(string); ok && v != "" {
		out.TargetUrl = utils.StringPtr(v)
	}
	if v, ok := m["tenant_ext_id"].(string); ok && v != "" {
		out.TenantExtId = utils.StringPtr(v)
	}
	if v, ok := m["credentials"].([]interface{}); ok && len(v) > 0 {
		out.Credentials = expandRegistrationCredentials(v)
	}
	return out
}

func expandRegistrationCredentials(in []interface{}) *import7.RegistrationCredentials {
	if len(in) == 0 || in[0] == nil {
		return nil
	}
	m, ok := in[0].(map[string]interface{})
	if !ok {
		return nil
	}
	out := import7.NewRegistrationCredentials()
	if v, ok := m["api_key"].(string); ok && v != "" {
		out.ApiKey = utils.StringPtr(v)
	}
	if v, ok := m["key_id"].(string); ok && v != "" {
		out.KeyId = utils.StringPtr(v)
	}
	return out
}

func expandRegisteredDomainTunnel(in []interface{}) *config.DomainTunnel {
	if len(in) == 0 || in[0] == nil {
		return nil
	}
	m, ok := in[0].(map[string]interface{})
	if !ok {
		return nil
	}
	out := config.NewDomainTunnel()
	if v, ok := m["endpoint"].(string); ok && v != "" {
		out.Endpoint = utils.StringPtr(v)
	}
	if v, ok := m["provider"].(string); ok && v != "" {
		out.Provider = commonUtils.ExpandEnum[import7.TunnelProvider](v)
	}
	return out
}

// setRegisteredDomainToResourceData maps a RegisteredDomain onto the resource /
// singular datasource schema. Both surfaces share the identical field set.
func setRegisteredDomainToResourceData(d *schema.ResourceData, rd *config.RegisteredDomain) error {
	if err := d.Set("name", utils.StringValue(rd.Name)); err != nil {
		return err
	}
	if err := d.Set("dns_name", utils.StringValue(rd.DnsName)); err != nil {
		return err
	}
	if err := d.Set("fqdn", utils.StringValue(rd.Fqdn)); err != nil {
		return err
	}
	if err := d.Set("domain_ext_id", utils.StringValue(rd.ExtId)); err != nil {
		return err
	}
	if err := d.Set("location_ext_id", utils.StringValue(rd.LocationExtId)); err != nil {
		return err
	}
	if err := d.Set("credential_issuer", utils.StringValue(rd.CredentialIssuer)); err != nil {
		return err
	}
	if err := d.Set("created_by", utils.StringValue(rd.CreatedBy)); err != nil {
		return err
	}
	if err := d.Set("api_credential_status", flattenApiCredentialStatus(rd.ApiCredentialStatus)); err != nil {
		return err
	}
	if err := d.Set("registration_state", flattenRegistrationState(rd.RegistrationState)); err != nil {
		return err
	}
	if err := d.Set("registered_time", utils.TimeStringValue(rd.RegisteredTime)); err != nil {
		return err
	}
	if err := d.Set("platform_data", flattenRegisteredDomainPlatformData(rd.PlatformData)); err != nil {
		return err
	}
	if err := d.Set("registration_config", flattenRegisteredDomainRegistrationConfig(rd.RegistrationConfig)); err != nil {
		return err
	}
	if err := d.Set("tunnel", flattenRegisteredDomainTunnel(rd.Tunnel)); err != nil {
		return err
	}
	if err := d.Set("ext_id", utils.StringValue(rd.ExtId)); err != nil {
		return err
	}
	if err := d.Set("tenant_id", utils.StringValue(rd.TenantId)); err != nil {
		return err
	}
	if err := d.Set("links", flattenLinks(rd.Links)); err != nil {
		return err
	}
	return nil
}

func flattenApiCredentialStatus(status *config.ApiCredentialStatus) string {
	if status == nil {
		return ""
	}
	return status.GetName()
}

func flattenRegistrationState(state *import7.RegistrationState) string {
	if state == nil {
		return ""
	}
	return state.GetName()
}

func flattenTunnelProvider(provider *import7.TunnelProvider) string {
	if provider == nil {
		return ""
	}
	return provider.GetName()
}

func flattenRegisteredDomainPlatformData(pd *config.PlatformData) []map[string]interface{} {
	if pd == nil {
		return []map[string]interface{}{}
	}
	return []map[string]interface{}{
		{
			"cluster_count":       utils.IntValue(pd.ClusterCount),
			"connectivity_status": utils.StringValue(pd.ConnectivityStatus),
			"pc_env_type":         utils.StringValue(pd.PcEnvType),
			"pc_provider_type":    utils.StringValue(pd.PcProviderType),
			"pc_version":          utils.StringValue(pd.PcVersion),
		},
	}
}

func flattenRegisteredDomainRegistrationConfig(rc *config.RegistrationConfig) []map[string]interface{} {
	if rc == nil {
		return []map[string]interface{}{}
	}
	return []map[string]interface{}{
		{
			"domain_ext_id": utils.StringValue(rc.DomainExtId),
			"target_url":    utils.StringValue(rc.TargetUrl),
			"tenant_ext_id": utils.StringValue(rc.TenantExtId),
			"credentials":   flattenRegistrationCredentials(rc.Credentials),
		},
	}
}

func flattenRegistrationCredentials(rc *import7.RegistrationCredentials) []map[string]interface{} {
	if rc == nil {
		return []map[string]interface{}{}
	}
	return []map[string]interface{}{
		{
			"api_key": utils.StringValue(rc.ApiKey),
			"key_id":  utils.StringValue(rc.KeyId),
		},
	}
}

func flattenRegisteredDomainTunnel(t *config.DomainTunnel) []map[string]interface{} {
	if t == nil {
		return []map[string]interface{}{}
	}
	return []map[string]interface{}{
		{
			"endpoint": utils.StringValue(t.Endpoint),
			"provider": flattenTunnelProvider(t.Provider),
		},
	}
}
