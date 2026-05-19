package ndb

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixNDBOnboarding() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixNDBOnboardingCreate,
		ReadContext:   resourceNutanixNDBOnboardingRead,
		UpdateContext: resourceNutanixNDBOnboardingUpdate,
		DeleteContext: resourceNutanixNDBOnboardingDelete,
		Schema: map[string]*schema.Schema{
			"prism_central_info": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name":        {Type: schema.TypeString, Optional: true},
						"description": {Type: schema.TypeString, Optional: true},
						"ip_address":  {Type: schema.TypeString, Required: true},
						"port":        {Type: schema.TypeInt, Optional: true, Default: 9440},
						"username":    {Type: schema.TypeString, Required: true},
						"password":    {Type: schema.TypeString, Required: true, Sensitive: true},
					},
				},
			},
			"prism_element_info": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name":        {Type: schema.TypeString, Required: true},
						"description": {Type: schema.TypeString, Optional: true},
						"cluster_ip":  {Type: schema.TypeString, Required: true},
						"username":    {Type: schema.TypeString, Required: true},
						"password":    {Type: schema.TypeString, Required: true, Sensitive: true},
						"version":     {Type: schema.TypeString, Optional: true, Default: "v2"},
						"cloud_type":  {Type: schema.TypeString, Optional: true, Default: "NTNX"},
					},
				},
			},
			"ndb_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dns_servers":           {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"ntp_servers":           {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"timezone":              {Type: schema.TypeString, Optional: true, Default: "UTC"},
						"smtp_server_ip_port":   {Type: schema.TypeString, Optional: true},
						"smtp_username":         {Type: schema.TypeString, Optional: true},
						"smtp_password":         {Type: schema.TypeString, Optional: true, Sensitive: true},
						"email_from_address":    {Type: schema.TypeString, Optional: true},
						"smtp_tls_enabled":      {Type: schema.TypeBool, Optional: true},
						"smtp_unsecured":        {Type: schema.TypeBool, Optional: true},
						"apply_smtp_even_empty": {Type: schema.TypeBool, Optional: true, Default: true},
					},
				},
			},
			"storage": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"container_name": {Type: schema.TypeString, Optional: true},
					},
				},
			},
			"network_details": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"skip":                  {Type: schema.TypeBool, Optional: true, Default: true},
						"existing_network_name": {Type: schema.TypeString, Optional: true},
						"vlan_name":             {Type: schema.TypeString, Optional: true},
						"static_ip":             {Type: schema.TypeString, Optional: true},
						"gateway":               {Type: schema.TypeString, Optional: true},
						"subnet_mask":           {Type: schema.TypeString, Optional: true},
					},
				},
			},
			"setup": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"trigger":         {Type: schema.TypeBool, Optional: true, Default: true},
						"timeout_minutes": {Type: schema.TypeInt, Optional: true, Default: 90},
					},
				},
			},
			"enable_full_onboarding": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"selection_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "auto",
				ValidateFunc: validationStringInSlice([]string{"auto", "strict"}),
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operation_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"current_step": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"completed_steps": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"available_storage_containers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"available_dns_servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"available_ntp_servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"available_network_names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"effective_storage_container": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"effective_dns_servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"effective_ntp_servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"effective_network_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"setup_operation_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"setup_progress_percent": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"setup_current_step": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNutanixNDBOnboardingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era
	completed := make([]string, 0)

	// Step 1: Prism Central optional
	if pcInfo, ok := d.GetOk("prism_central_info"); ok && len(pcInfo.([]interface{})) > 0 {
		d.Set("current_step", "step1_prism_central")
		if err := applyPrismCentralStep(ctx, conn.Service, pcInfo.([]interface{})); err != nil {
			return diag.FromErr(fmt.Errorf("step1 prism central failed: %w", err))
		}
		completed = append(completed, "step1_prism_central")
	} else {
		completed = append(completed, "step1_prism_central_skipped")
	}

	// Step 2: Prism Element required
	d.Set("current_step", "step2_prism_element")
	clusterID, opID, err := applyPrismElementStep(ctx, conn.Service, d.Get("prism_element_info").([]interface{}))
	if err != nil {
		return diag.FromErr(fmt.Errorf("step2 prism element failed: %w", err))
	}
	_ = d.Set("cluster_id", clusterID)
	_ = d.Set("operation_id", opID)
	completed = append(completed, "step2_prism_element")
	selectionMode := d.Get("selection_mode").(string)

	discovery, err := discoverOnboardingOptions(ctx, conn.Service, clusterID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("step2 discovery failed: %w", err))
	}
	_ = d.Set("available_storage_containers", discovery.StorageContainers)
	_ = d.Set("available_dns_servers", discovery.DNSServers)
	_ = d.Set("available_ntp_servers", discovery.NTPServers)
	_ = d.Set("available_network_names", discovery.NetworkNames)

	// Default safe mode: stop after Step 2 to avoid partial onboarding states.
	fullFlow := d.Get("enable_full_onboarding").(bool)
	if !fullFlow {
		d.SetId(clusterID)
		completed = append(completed, "step3_ndb_config_skipped")
		completed = append(completed, "step4_storage_skipped")
		completed = append(completed, "step5_network_skipped")
		completed = append(completed, "step6_setup_skipped")
		_ = d.Set("completed_steps", completed)
		_ = d.Set("current_step", "step2_prism_element")
		_ = d.Set("status", "STEP2_COMPLETED")
		log.Printf("NDB onboarding safe mode completed for cluster id %s", clusterID)
		return resourceNutanixNDBOnboardingRead(ctx, d, meta)
	}

	// Step 3: NDB config (user-provided or discovered)
	ndbCfg := resolveNDBConfig(d.Get("ndb_config"), discovery)
	_ = d.Set("effective_dns_servers", ndbCfg.DNSServers)
	_ = d.Set("effective_ntp_servers", ndbCfg.NTPServers)
	d.Set("current_step", "step3_ndb_config")
	if err := applyNDBConfigStep(ctx, conn.Service, ndbCfg); err != nil {
		return diag.FromErr(fmt.Errorf("step3 ndb config failed: %w", err))
	}
	completed = append(completed, "step3_ndb_config")

	// Step 4: storage (user-provided or discovered)
	storageContainer, err := resolveStorageContainer(selectionMode, d.Get("storage"), discovery.StorageContainers)
	if err != nil {
		return diag.FromErr(fmt.Errorf("step4 storage selection failed: %w", err))
	}
	_ = d.Set("effective_storage_container", storageContainer)
	d.Set("current_step", "step4_storage")
	if err := applyStorageStep(ctx, conn.Service, clusterID, storageContainer, d.Get("prism_element_info").([]interface{})); err != nil {
		return diag.FromErr(fmt.Errorf("step4 storage failed: %w", err))
	}
	completed = append(completed, "step4_storage")

	// Step 5: network optional
	resolvedNetwork := resolveNetworkDetails(d.Get("network_details"), discovery)
	_ = d.Set("effective_network_name", resolvedNetwork.NetworkName)
	d.Set("current_step", "step5_network")
	if err := applyNetworkStep(ctx, conn.Service, clusterID, resolvedNetwork); err != nil {
		return diag.FromErr(fmt.Errorf("step5 network failed: %w", err))
	}
	completed = append(completed, "step5_network")

	// Step 6: setup start
	d.Set("current_step", "step6_setup")
	setupOperationID, err := applySetupStep(ctx, conn.Service, clusterID, d.Get("setup"))
	if err != nil {
		return diag.FromErr(fmt.Errorf("step6 setup failed: %w", err))
	}
	_ = d.Set("setup_operation_id", setupOperationID)
	_ = d.Set("setup_progress_percent", 100)
	_ = d.Set("setup_current_step", "")
	_ = d.Set("operation_id", setupOperationID)
	completed = append(completed, "step6_setup")

	d.SetId(clusterID)
	_ = d.Set("completed_steps", completed)
	_ = d.Set("status", "COMPLETED")
	log.Printf("NDB onboarding flow completed for cluster id %s", clusterID)
	return resourceNutanixNDBOnboardingRead(ctx, d, meta)
}

func resourceNutanixNDBOnboardingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era
	if d.Id() == "" {
		return nil
	}
	cluster, err := conn.Service.GetCluster(ctx, d.Id(), "")
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("cluster_id", cluster.ID)
	if cluster.Status != nil {
		_ = d.Set("status", *cluster.Status)
	}
	return nil
}

func resourceNutanixNDBOnboardingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Re-run wizard flow to keep behavior deterministic.
	return resourceNutanixNDBOnboardingCreate(ctx, d, meta)
}

func resourceNutanixNDBOnboardingDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// Onboarding is workflow-oriented; keep delete as state-only removal.
	d.SetId("")
	return nil
}

func applyPrismCentralStep(ctx context.Context, svc era.Service, pcInfo []interface{}) error {
	if len(pcInfo) == 0 {
		return nil
	}
	pc := pcInfo[0].(map[string]interface{})
	body := map[string]interface{}{
		"name":      pc["name"].(string),
		"ipAddress": pc["ip_address"].(string),
		"port":      pc["port"].(int),
		"username":  pc["username"].(string),
		"password":  pc["password"].(string),
	}
	if desc, ok := pc["description"]; ok && desc.(string) != "" {
		body["description"] = desc.(string)
	}
	_, err := svc.CreateCluster(ctx, map[string]interface{}{
		"managementServerInfo": body,
	})
	return err
}

func applyPrismElementStep(ctx context.Context, svc era.Service, peInfo []interface{}) (string, string, error) {
	pe := peInfo[0].(map[string]interface{})
	clusterIP := pe["cluster_ip"].(string)
	if existing, err := svc.GetCluster(ctx, "", pe["name"].(string)); err == nil && existing.ID != nil {
		if uploadErr := uploadClusterCredentialsCompat(ctx, svc, *existing.ID, clusterIP, pe["username"].(string), pe["password"].(string)); uploadErr != nil {
			log.Printf("[WARN] step2 credentials upload (existing cluster) failed for %s: %v", *existing.ID, uploadErr)
		}
		return *existing.ID, "", nil
	}
	basePayload := map[string]interface{}{
		"name":        pe["name"].(string),
		"ipAddresses": []string{clusterIP},
		"cloudType":   pe["cloud_type"].(string),
		"version":     pe["version"].(string),
		"username":    pe["username"].(string),
		"password":    pe["password"].(string),
		"status":      "UP",
		"description": "",
	}
	if desc, ok := pe["description"]; ok {
		basePayload["description"] = desc.(string)
	}

	// Retry ladder for /clusters payload contract differences across NDB versions.
	// First try modern UI variants, then strict 2.11 flat clusterIP variants.
	ladder := []map[string]interface{}{
		copyMap(basePayload),
		removeKeys(copyMap(basePayload), "status"),
		removeKeys(copyMap(basePayload), "description"),
		removeKeys(copyMap(basePayload), "status", "description"),
		{
			"name":      pe["name"].(string),
			"clusterIP": clusterIP,
			"username":  pe["username"].(string),
			"password":  pe["password"].(string),
		},
		{
			"name":      pe["name"].(string),
			"clusterIP": clusterIP,
			"username":  pe["username"].(string),
			"password":  pe["password"].(string),
			"version":   pe["version"].(string),
		},
		{
			"name":      pe["name"].(string),
			"clusterIP": clusterIP,
			"username":  pe["username"].(string),
			"password":  pe["password"].(string),
			"version":   pe["version"].(string),
			"cloudType": pe["cloud_type"].(string),
		},
		{
			"name":       pe["name"].(string),
			"ip_address": clusterIP,
			"username":   pe["username"].(string),
			"password":   pe["password"].(string),
		},
	}

	var (
		resp *era.ProvisionDatabaseResponse
		err  error
	)

	for _, payload := range ladder {
		resp, err = svc.CreateCluster(ctx, payload)
		if err == nil {
			break
		}
		// If contract error indicates old schema, break and try old format below.
		if containsAny(err.Error(),
			"Unrecognized field 'cloudType'",
			"Unrecognized field 'ipAddresses'",
			"Unrecognized field 'name'",
			"Unrecognized field 'status'",
			"Unrecognized field 'description'",
		) {
			break
		}
	}
	if err != nil {
		// Handle older NDB payload contract that expects clusterIP/clusterType/credentialsInfo.
		if containsAny(err.Error(),
			"Unrecognized field 'cloudType'",
			"Unrecognized field 'ipAddresses'",
			"Unrecognized field 'name'",
			"Unrecognized field 'status'",
			"Unrecognized field 'description'",
		) {
			clusterDescription := ""
			if desc, ok := pe["description"]; ok {
				clusterDescription = desc.(string)
			}
			altPayload := map[string]interface{}{
				"clusterName":        pe["name"].(string),
				"clusterDescription": clusterDescription,
				"clusterIP":          clusterIP,
				"clusterType":        pe["cloud_type"].(string),
				"version":            pe["version"].(string),
				"credentialsInfo": []map[string]interface{}{
					{
						"name":  "username",
						"value": pe["username"].(string),
					},
					{
						"name":  "password",
						"value": pe["password"].(string),
					},
				},
			}
			resp, err = svc.CreateCluster(ctx, altPayload)
		}
	}
	if err != nil {
		// Some NDB versions return a generic internal error when cluster already exists.
		// Treat this as idempotent by resolving the cluster by name.
		if strings.Contains(strings.ToLower(err.Error()), "internal error") {
			cluster, er := svc.GetCluster(ctx, "", pe["name"].(string))
			if er == nil && cluster.ID != nil {
				return *cluster.ID, "", nil
			}
		}
		return "", "", err
	}
	if resp.Entityid != "" {
		if uploadErr := uploadClusterCredentialsCompat(ctx, svc, resp.Entityid, clusterIP, pe["username"].(string), pe["password"].(string)); uploadErr != nil {
			log.Printf("[WARN] step2 credentials upload (new cluster) failed for %s: %v", resp.Entityid, uploadErr)
		}
	}
	if resp.Entityid != "" {
		return resp.Entityid, resp.Operationid, nil
	}
	// Fallback lookup by name if API returns sync cluster object semantics.
	cluster, er := svc.GetCluster(ctx, "", pe["name"].(string))
	if er != nil {
		return "", "", er
	}
	if cluster.ID == nil {
		return "", resp.Operationid, nil
	}
	if uploadErr := uploadClusterCredentialsCompat(ctx, svc, *cluster.ID, clusterIP, pe["username"].(string), pe["password"].(string)); uploadErr != nil {
		log.Printf("[WARN] step2 credentials upload (lookup fallback) failed for %s: %v", *cluster.ID, uploadErr)
	}
	return *cluster.ID, resp.Operationid, nil
}

func uploadClusterCredentialsCompat(ctx context.Context, svc era.Service, clusterID, clusterIP, username, password string) error {
	payload := map[string]interface{}{
		"protocol":   "https",
		"ip_address": clusterIP,
		"ip":         clusterIP,
		"port":       "9440",
		"creds_bag": map[string]interface{}{
			"username": username,
			"password": password,
		},
		"username": username,
		"password": password,
	}

	var lastErr error
	// Different NDB versions/builds accept different skip/query semantics for /clusters/{id}/json.
	for _, flags := range []struct {
		skipUpload  bool
		skipProfile bool
		updateJSON  bool
	}{
		{false, false, false},
		{false, true, true},
		{true, true, true},
	} {
		_, err := svc.UploadClusterWizardJSON(ctx, clusterID, payload, flags.skipUpload, flags.skipProfile, flags.updateJSON)
		if err == nil {
			return nil
		}
		lastErr = err
	}
	return lastErr
}

func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

func copyMap(src map[string]interface{}) map[string]interface{} {
	dst := make(map[string]interface{}, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func removeKeys(src map[string]interface{}, keys ...string) map[string]interface{} {
	for _, k := range keys {
		delete(src, k)
	}
	return src
}

func applyNDBConfigStep(ctx context.Context, svc era.Service, cfg onboardingResolvedConfig) error {
	dnsServers := cfg.DNSServers
	ntpServers := cfg.NTPServers
	timezone := cfg.Timezone

	// UI validates DNS+NTP first before sending split updates.
	validateBody := &era.OnboardingEraServerConfig{
		DNSServers: dnsServers,
		NTPServers: ntpServers,
	}
	if _, err := svc.ValidateEraServerConfig(ctx, validateBody, []string{"dns", "ntp"}); err != nil {
		return err
	}

	if len(dnsServers) > 0 {
		if _, err := svc.SetEraServerConfig(ctx, &era.OnboardingEraServerConfig{
			DNSServers: dnsServers,
		}, []string{"dns"}); err != nil {
			return err
		}
	}

	if len(ntpServers) > 0 {
		if _, err := svc.SetEraServerConfig(ctx, &era.OnboardingEraServerConfig{
			NTPServers: ntpServers,
		}, []string{"ntp"}); err != nil {
			return err
		}
	}

	if timezone != "" {
		if _, err := svc.SetEraServerConfig(ctx, &era.OnboardingEraServerConfig{
			Timezone: utils.StringPtr(timezone),
		}, []string{"timezone"}); err != nil {
			return err
		}
	}

	applySMTP := cfg.ApplySMTPEvenEmpty
	smtpServerIPPort := cfg.SMTPServerIPPort
	smtpUsername := cfg.SMTPUsername
	smtpPassword := cfg.SMTPPassword
	emailFromAddress := cfg.EmailFromAddress
	tlsSet := cfg.SMTPTLSEnabled != nil
	unsecuredSet := cfg.SMTPUnsecured != nil

	hasSMTPValues := smtpServerIPPort != "" || smtpUsername != "" || smtpPassword != "" || emailFromAddress != "" || tlsSet || unsecuredSet
	if applySMTP || hasSMTPValues {
		smtpCfg := &era.OnboardingEraSMTPConfig{}
		if smtpServerIPPort != "" {
			smtpCfg.SMTPServerIPPort = utils.StringPtr(smtpServerIPPort)
		}
		if smtpUsername != "" {
			smtpCfg.SMTPUsername = utils.StringPtr(smtpUsername)
		}
		if smtpPassword != "" {
			smtpCfg.SMTPPassword = utils.StringPtr(smtpPassword)
		}
		if emailFromAddress != "" {
			smtpCfg.EmailFromAddress = utils.StringPtr(emailFromAddress)
		}
		if tlsSet {
			smtpCfg.TLSEnabled = cfg.SMTPTLSEnabled
		}
		if unsecuredSet {
			smtpCfg.Unsecured = cfg.SMTPUnsecured
		}
		if _, err := svc.SetEraServerConfig(ctx, &era.OnboardingEraServerConfig{
			SMTPConfig: smtpCfg,
		}, []string{"smtp"}); err != nil {
			return err
		}
	}

	return nil
}

func applyStorageStep(ctx context.Context, svc era.Service, clusterID, storageContainer string, peInfo []interface{}) error {
	if storageContainer == "" {
		return nil
	}
	pe := peInfo[0].(map[string]interface{})
	var lastErr error
	waitSchedule := []time.Duration{0, 6 * time.Second, 10 * time.Second, 15 * time.Second, 20 * time.Second}
	for _, wait := range waitSchedule {
		if wait > 0 {
			time.Sleep(wait)
		}
		cluster, err := svc.GetCluster(ctx, clusterID, "")
		if err != nil {
			lastErr = err
			if isRetryableStep4Error(err) {
				continue
			}
			return err
		}
		name := ""
		description := ""
		cloudType := "NTNX"
		version := "v2"
		if cluster.Name != nil {
			name = *cluster.Name
		}
		if cluster.Description != nil {
			description = *cluster.Description
		}
		if cluster.Cloudtype != nil && *cluster.Cloudtype != "" {
			cloudType = *cluster.Cloudtype
		}
		if cluster.Version != nil && *cluster.Version != "" {
			version = *cluster.Version
		}
		ipAddresses := make([]string, 0, len(cluster.Ipaddresses))
		for _, ip := range cluster.Ipaddresses {
			if ip != nil && *ip != "" {
				ipAddresses = append(ipAddresses, *ip)
			}
		}
		body := map[string]interface{}{
			"name":        name,
			"description": description,
			"ipAddresses": ipAddresses,
			"username":    pe["username"].(string),
			"password":    pe["password"].(string),
			"status":      "UP",
			"version":     version,
			"cloudType":   cloudType,
			"properties": []map[string]interface{}{
				{
					"name":  "ERA_STORAGE_CONTAINER",
					"value": storageContainer,
				},
			},
		}
		_, err = svc.ReplaceClusterWizard(ctx, clusterID, body)
		if err == nil {
			return nil
		}
		lastErr = err
		// Fallback for NDB versions that reject the PUT shape during onboarding.
		_, fallbackErr := svc.UploadClusterWizardJSON(ctx, clusterID, map[string]interface{}{
			"storageContainer": storageContainer,
		}, true, true, true)
		if fallbackErr == nil {
			return nil
		}
		lastErr = err
		if !isRetryableStep4Error(err) && !isRetryableStep4Error(fallbackErr) {
			return err
		}
	}
	if lastErr != nil {
		return lastErr
	}
	return fmt.Errorf("storage update did not complete")
}

func applyNetworkStep(ctx context.Context, svc era.Service, clusterID string, net onboardingResolvedNetwork) error {
	cluster, err := svc.GetCluster(ctx, clusterID, "")
	if err != nil {
		return err
	}
	clusterIP := ""
	if len(cluster.Ipaddresses) > 0 && cluster.Ipaddresses[0] != nil {
		clusterIP = *cluster.Ipaddresses[0]
	}
	if clusterIP == "" {
		clusterIP = net.StaticIP
	}
	if net.Skip {
		_, err = svc.UploadClusterWizardJSON(ctx, clusterID, map[string]interface{}{
			"protocol":   "https",
			"port":       "9440",
			"ip_address": clusterIP,
			"clusterIp":  clusterIP,
		}, true, false, false)
		return err
	}
	ipAddress := clusterIP
	if net.StaticIP != "" {
		ipAddress = net.StaticIP
	}
	body := map[string]interface{}{
		"protocol":   "https",
		"port":       "9440",
		"ip_address": ipAddress,
		"clusterIp":  ipAddress,
	}
	if net.NetworkName != "" {
		body["vlanName"] = net.NetworkName
	}
	if net.Gateway != "" {
		body["gateway"] = net.Gateway
	}
	if net.SubnetMask != "" {
		body["subnetMask"] = net.SubnetMask
	}
	_, err = svc.UploadClusterWizardJSON(ctx, clusterID, body, true, false, false)
	return err
}

func applySetupStep(ctx context.Context, svc era.Service, clusterID string, setupRaw interface{}) (string, error) {
	trigger := true
	timeoutMinutes := 90
	if setupRaw != nil {
		sl := setupRaw.([]interface{})
		if len(sl) > 0 {
			setupCfg := sl[0].(map[string]interface{})
			trigger = setupCfg["trigger"].(bool)
			if v, ok := setupCfg["timeout_minutes"]; ok && v.(int) > 0 {
				timeoutMinutes = v.(int)
			}
		}
	}
	if !trigger {
		return "", nil
	}
	_, err := svc.UploadClusterWizardJSON(ctx, clusterID, map[string]interface{}{
		"action": "setup",
	}, true, true, true)
	if err != nil {
		return "", err
	}

	operationID, err := waitForSetupOperationID(ctx, svc, clusterID)
	if err != nil {
		return "", err
	}
	return operationID, waitForOperationCompletion(ctx, svc, operationID, time.Duration(timeoutMinutes)*time.Minute)
}

func expandStringList(items []interface{}) []string {
	out := make([]string, 0, len(items))
	for _, v := range items {
		if s, ok := v.(string); ok && s != "" {
			out = append(out, s)
		}
	}
	return out
}

func waitForSetupOperationID(ctx context.Context, svc era.Service, clusterID string) (string, error) {
	timeout := time.NewTimer(2 * time.Minute)
	ticker := time.NewTicker(3 * time.Second)
	defer timeout.Stop()
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-timeout.C:
			return "", fmt.Errorf("timed out waiting for setup operation id for cluster %s", clusterID)
		case <-ticker.C:
			shortInfo, err := svc.GetOperationsShortInfo(ctx)
			if err != nil {
				continue
			}
			for _, op := range shortInfo.Operations {
				if op.Type == nil || op.NxClusterID == nil || op.ID == nil {
					continue
				}
				if *op.Type == "configure_oob_software_profiles" && *op.NxClusterID == clusterID {
					return *op.ID, nil
				}
			}
		}
	}
}

func waitForOperationCompletion(ctx context.Context, svc era.Service, operationID string, timeoutDuration time.Duration) error {
	timeout := time.NewTimer(timeoutDuration)
	ticker := time.NewTicker(10 * time.Second)
	defer timeout.Stop()
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeout.C:
			return fmt.Errorf("timed out waiting for operation %s completion", operationID)
		case <-ticker.C:
			op, err := svc.GetOperation(era.GetOperationRequest{OperationID: operationID})
			if err != nil {
				continue
			}
			if op.Status == nil {
				continue
			}
			switch *op.Status {
			case "5":
				return nil
			case "3", "4", "6":
				return fmt.Errorf("setup operation %s failed with status %s", operationID, *op.Status)
			}
		}
	}
}

type onboardingDiscovery struct {
	StorageContainers []string
	DNSServers        []string
	NTPServers        []string
	NetworkNames      []string
}

type onboardingResolvedConfig struct {
	DNSServers         []string
	NTPServers         []string
	Timezone           string
	SMTPServerIPPort   string
	SMTPUsername       string
	SMTPPassword       string
	EmailFromAddress   string
	SMTPTLSEnabled     *bool
	SMTPUnsecured      *bool
	ApplySMTPEvenEmpty bool
}

type onboardingResolvedNetwork struct {
	Skip        bool
	NetworkName string
	StaticIP    string
	Gateway     string
	SubnetMask  string
}

func discoverOnboardingOptions(ctx context.Context, svc era.Service, clusterID string) (*onboardingDiscovery, error) {
	out := &onboardingDiscovery{}
	var storageErr error
	for _, wait := range []time.Duration{0, 5 * time.Second, 10 * time.Second, 15 * time.Second} {
		if wait > 0 {
			time.Sleep(wait)
		}
		storageResp, err := svc.GetClusterStorageContainers(ctx, clusterID)
		if err == nil {
			out.StorageContainers = parseStorageContainers(storageResp)
			storageErr = nil
			break
		}
		storageErr = err
		if !isRetryableStep4Error(err) {
			break
		}
	}
	if storageErr != nil && !isRetryableStep4Error(storageErr) {
		return nil, storageErr
	}
	serverCfg, err := svc.GetEraServerConfig(ctx)
	if err == nil {
		out.DNSServers = uniqStrings(serverCfg.DNSServers)
		out.NTPServers = uniqStrings(serverCfg.NTPServers)
	}
	networks, err := svc.ListNetwork(ctx)
	if err == nil {
		out.NetworkNames = parseNetworkNames(networks)
	}
	return out, nil
}

func parseStorageContainers(storageResp map[string]interface{}) []string {
	names := []string{}
	entities, ok := storageResp["entities"].([]interface{})
	if !ok {
		return names
	}
	for _, raw := range entities {
		entity, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		if n, ok := entity["name"].(string); ok && n != "" {
			names = append(names, n)
		}
		if n, ok := entity["containerName"].(string); ok && n != "" {
			names = append(names, n)
		}
		if vstores, ok := entity["vstore_name_list"].([]interface{}); ok {
			for _, v := range vstores {
				if s, ok := v.(string); ok && s != "" {
					names = append(names, s)
				}
			}
		}
	}
	return uniqStrings(names)
}

func parseNetworkNames(networks *era.ListNetworkResponse) []string {
	if networks == nil {
		return []string{}
	}
	out := make([]string, 0, len(*networks))
	for _, n := range *networks {
		if n == nil || n.Name == nil || *n.Name == "" {
			continue
		}
		out = append(out, *n.Name)
	}
	return uniqStrings(out)
}

func resolveNDBConfig(raw interface{}, discovery *onboardingDiscovery) onboardingResolvedConfig {
	resolved := onboardingResolvedConfig{
		DNSServers:         discovery.DNSServers,
		NTPServers:         discovery.NTPServers,
		Timezone:           "UTC",
		ApplySMTPEvenEmpty: true,
	}
	cfgList, ok := raw.([]interface{})
	if ok && len(cfgList) > 0 && cfgList[0] != nil {
		cfg := cfgList[0].(map[string]interface{})
		if v, ok := cfg["dns_servers"]; ok {
			dns := expandStringList(v.([]interface{}))
			if len(dns) > 0 {
				resolved.DNSServers = dns
			}
		}
		if v, ok := cfg["ntp_servers"]; ok {
			ntp := expandStringList(v.([]interface{}))
			if len(ntp) > 0 {
				resolved.NTPServers = ntp
			}
		}
		if v, ok := cfg["timezone"].(string); ok && v != "" {
			resolved.Timezone = v
		}
		if v, ok := cfg["smtp_server_ip_port"].(string); ok {
			resolved.SMTPServerIPPort = v
		}
		if v, ok := cfg["smtp_username"].(string); ok {
			resolved.SMTPUsername = v
		}
		if v, ok := cfg["smtp_password"].(string); ok {
			resolved.SMTPPassword = v
		}
		if v, ok := cfg["email_from_address"].(string); ok {
			resolved.EmailFromAddress = v
		}
		if v, ok := cfg["apply_smtp_even_empty"].(bool); ok {
			resolved.ApplySMTPEvenEmpty = v
		}
		if v, ok := cfg["smtp_tls_enabled"].(bool); ok {
			resolved.SMTPTLSEnabled = &v
		}
		if v, ok := cfg["smtp_unsecured"].(bool); ok {
			resolved.SMTPUnsecured = &v
		}
	}
	if len(resolved.DNSServers) == 0 {
		resolved.DNSServers = []string{"10.40.64.15", "10.40.64.16"}
	}
	if len(resolved.NTPServers) == 0 {
		resolved.NTPServers = []string{"pool.ntp.org"}
	}
	return resolved
}

func resolveStorageContainer(selectionMode string, raw interface{}, available []string) (string, error) {
	var userValue string
	if sl, ok := raw.([]interface{}); ok && len(sl) > 0 && sl[0] != nil {
		if s, ok := sl[0].(map[string]interface{}); ok {
			if v, ok := s["container_name"].(string); ok {
				userValue = strings.TrimSpace(v)
			}
		}
	}
	if userValue != "" {
		if selectionMode == "strict" && len(available) > 0 && !containsString(available, userValue) {
			return "", fmt.Errorf("storage container %q not found in discovered options: %s", userValue, strings.Join(available, ", "))
		}
		return userValue, nil
	}
	if len(available) > 0 {
		return available[0], nil
	}
	if selectionMode == "strict" {
		return "", fmt.Errorf("no storage container provided and no discovered options available")
	}
	return "", fmt.Errorf("no storage container resolved from user input or discovery")
}

func resolveNetworkDetails(raw interface{}, discovery *onboardingDiscovery) onboardingResolvedNetwork {
	resolved := onboardingResolvedNetwork{Skip: true}
	netList, ok := raw.([]interface{})
	if !ok || len(netList) == 0 || netList[0] == nil {
		return resolved
	}
	n := netList[0].(map[string]interface{})
	if v, ok := n["skip"].(bool); ok {
		resolved.Skip = v
	}
	existingName := ""
	if v, ok := n["existing_network_name"].(string); ok {
		existingName = strings.TrimSpace(v)
	}
	if existingName == "" {
		if v, ok := n["vlan_name"].(string); ok {
			existingName = strings.TrimSpace(v)
		}
	}
	if existingName == "" && !resolved.Skip && len(discovery.NetworkNames) > 0 {
		existingName = discovery.NetworkNames[0]
	}
	resolved.NetworkName = existingName
	if v, ok := n["static_ip"].(string); ok {
		resolved.StaticIP = strings.TrimSpace(v)
	}
	if v, ok := n["gateway"].(string); ok {
		resolved.Gateway = strings.TrimSpace(v)
	}
	if v, ok := n["subnet_mask"].(string); ok {
		resolved.SubnetMask = strings.TrimSpace(v)
	}
	return resolved
}

func containsString(values []string, target string) bool {
	for _, v := range values {
		if v == target {
			return true
		}
	}
	return false
}

func uniqStrings(input []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(input))
	for _, s := range input {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}

func validationStringInSlice(valid []string) schema.SchemaValidateFunc {
	return func(v interface{}, _ string) (ws []string, errors []error) {
		s, ok := v.(string)
		if !ok {
			return nil, []error{fmt.Errorf("must be a string")}
		}
		for _, allowed := range valid {
			if s == allowed {
				return nil, nil
			}
		}
		return nil, []error{fmt.Errorf("must be one of: %s", strings.Join(valid, ", "))}
	}
}

func isRetryableStep4Error(err error) bool {
	if err == nil {
		return false
	}
	errText := strings.ToLower(err.Error())
	return strings.Contains(errText, "could not get prism rest caller for cloudid") ||
		strings.Contains(errText, "era-0000000") ||
		strings.Contains(errText, "era-sql-0000001") ||
		strings.Contains(errText, "an internal error has occurred")
}
