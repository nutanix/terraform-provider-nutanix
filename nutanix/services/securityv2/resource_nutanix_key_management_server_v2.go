// Package securityv2 provides resources for managing security-related configurations in Nutanix.
package securityv2

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	commonCfg "github.com/nutanix/ntnx-api-golang-clients/security-go-client/v4/models/common/v1/config"
	securityPrism "github.com/nutanix/ntnx-api-golang-clients/security-go-client/v4/models/prism/v4/config"
	"github.com/nutanix/ntnx-api-golang-clients/security-go-client/v4/models/security/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"

	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixKeyManagementServerV2() *schema.Resource {
	endpointResource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ip_address": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1, //nolint:gomnd
				Elem:     common.SchemaForIPList(true),
			},
			"port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
		},
	}

	exactlyOneKMSAccess := []string{
		"access_information.0.azure_key_vault",
		"access_information.0.kmip_key_vault",
	}

	return &schema.Resource{
		CreateContext: ResourceNutanixKeyManagementServerV2Create,
		ReadContext:   ResourceNutanixKeyManagementServerV2Read,
		UpdateContext: ResourceNutanixKeyManagementServerV2Update,
		DeleteContext: ResourceNutanixKeyManagementServerV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(DEFAULTWAITTIMEOUT * time.Minute),
			Update: schema.DefaultTimeout(DEFAULTWAITTIMEOUT * time.Minute),
			Delete: schema.DefaultTimeout(DEFAULTWAITTIMEOUT * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"access_information": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1, //nolint:gomnd
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"azure_key_vault": {
							Type:         schema.TypeList,
							Optional:     true,
							MaxItems:     1, //nolint:gomnd
							ExactlyOneOf: exactlyOneKMSAccess,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"endpoint_url": {
										Type:     schema.TypeString,
										Required: true,
									},
									"key_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"tenant_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"client_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"client_secret": {
										Type:         schema.TypeString,
										Required:     true,
										Sensitive:    true,
										ValidateFunc: validation.StringLenBetween(8, 256),
									},
									"credential_expiry_date": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringMatch(
											regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`),
											"must be in YYYY-MM-DD format",
										),
									},
									"truncated_client_secret": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"kmip_key_vault": {
							Type:         schema.TypeList,
							Optional:     true,
							MaxItems:     1, //nolint:gomnd
							ExactlyOneOf: exactlyOneKMSAccess,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cert_pem": {
										Type:      schema.TypeString,
										Required:  true,
										Sensitive: true,
									},
									"private_key": {
										Type:      schema.TypeString,
										Required:  true,
										Sensitive: true,
									},
									"ca_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"ca_pem": {
										Type:      schema.TypeString,
										Required:  true,
										Sensitive: true,
									},
									"endpoint_url": {
										Type:     schema.TypeSet,
										Required: true,
										Elem:     endpointResource,
										Set:      schema.HashResource(endpointResource),
									},
								},
							},
						},
					},
				},
			},
			// Computed Attributes
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": common.LinksSchema(),
		},
	}
}

func ResourceNutanixKeyManagementServerV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).SecurityAPI

	body := config.NewKeyManagementServer()

	if name, ok := d.GetOk("name"); ok {
		body.Name = utils.StringPtr(name.(string))
	}
	if accessInfo, ok := d.GetOk("access_information"); ok {
		accessInfoVal, expandErr := expandAccessInformation(accessInfo.([]interface{}))
		if expandErr != nil {
			return diag.FromErr(expandErr)
		}
		if err := body.SetAccessInformation(accessInfoVal); err != nil {
			return diag.FromErr(err)
		}
	}

	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] key management server payload: %s", aJSON)

	resp, err := conn.KeyManagementServersAPIInstance.CreateKeyManagementServer(body)
	if err != nil {
		return diag.Errorf("error while creating Key Management Server : %v", err)
	}

	taskRefValue, ok := resp.Data.GetValue().(securityPrism.TaskReference)
	if !ok {
		return diag.Errorf("error: unexpected response type from create API, expected TaskReference")
	}
	TaskRef := taskRefValue
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task
	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the key management server to be created
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for key management server (%s) to be created: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching key management server create task (%s): %v", utils.StringValue(taskUUID), err)
	}
	taskDetailsValue, ok := taskResp.Data.GetValue().(prismConfig.Task)
	if !ok {
		return diag.Errorf("error: unexpected response type from task API, expected Task")
	}
	taskDetails := taskDetailsValue
	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Create Key Management Server Task Details: %s", string(aJSON))

	// Extract UUID from task using entity type constant
	kmsExtID, err := common.ExtractEntityUUIDFromTask(taskDetails, utils.RelEntityTypeKMS, "Key management server")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(utils.StringValue(kmsExtID))
	return ResourceNutanixKeyManagementServerV2Read(ctx, d, meta)
}

func ResourceNutanixKeyManagementServerV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).SecurityAPI

	resp, err := conn.KeyManagementServersAPIInstance.GetKeyManagementServerById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching key management server : %v", err)
	}

	getRespValue, ok := resp.Data.GetValue().(config.KeyManagementServer)
	if !ok {
		return diag.Errorf("error: unexpected response type from get API, expected KeyManagementServer")
	}
	getResp := getRespValue

	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] flattening access information")
	accessInfo, flattenErr := flattenAccessInformation(getResp.GetAccessInformation())
	if flattenErr != nil {
		return diag.FromErr(flattenErr)
	}
	// Preserve sensitive values from config/state (API may not return them).
	// This avoids perpetual diffs for required sensitive fields like client_secret.
	if v, ok := d.GetOk("access_information.0.azure_key_vault.0.client_secret"); ok && len(accessInfo) > 0 {
		if azureList, ok2 := accessInfo[0]["azure_key_vault"].([]map[string]interface{}); ok2 && len(azureList) > 0 {
			azureList[0]["client_secret"] = v.(string)
			accessInfo[0]["azure_key_vault"] = azureList
		}
	}
	if err := d.Set("access_information", accessInfo); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ResourceNutanixKeyManagementServerV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).SecurityAPI

	resp, err := conn.KeyManagementServersAPIInstance.GetKeyManagementServerById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching key management server : %v", err)
	}

	// Extract E-Tag Header
	etagValue := conn.KeyManagementServersAPIInstance.ApiClient.GetEtag(resp)
	args := make(map[string]interface{})
	args["If-Match"] = utils.StringPtr(etagValue)

	// Get current state to preserve values for fields that haven't changed
	updateSpec, ok := resp.Data.GetValue().(config.KeyManagementServer)
	if !ok {
		return diag.Errorf("error: unexpected response type from get API, expected KeyManagementServer")
	}

	// Update name if it has changed
	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			updateSpec.Name = utils.StringPtr(v.(string))
		}
	}

	// Update access_information if it has changed
	if d.HasChange("access_information") {
		if v, ok := d.GetOk("access_information"); ok {
			accessInfo, expandErr := expandAccessInformation(v.([]interface{}))
			if expandErr != nil {
				return diag.FromErr(expandErr)
			}

			// The API requires clientId and clientSecret to be present in update requests for Azure.
			// Get values from raw config as they may not be returned by API / may be sensitive.
			if _, isAzure := accessInfo.(config.AzureAccessInformation); isAzure {
				rawConfig := d.GetRawConfig()
				configMap := rawConfig.AsValueMap()

				var configClientSecret, configClientId string
				if accessInfoVal, exists := configMap["access_information"]; exists && !accessInfoVal.IsNull() && accessInfoVal.IsKnown() {
					if accessInfoVal.Type().IsListType() || accessInfoVal.Type().IsTupleType() {
						accessInfoList := accessInfoVal.AsValueSlice()
						if len(accessInfoList) > 0 {
							firstItem := accessInfoList[0]
							if firstItem.Type().IsObjectType() {
								itemMap := firstItem.AsValueMap()
								if azureVal, exists := itemMap["azure_key_vault"]; exists && !azureVal.IsNull() && azureVal.IsKnown() {
									azureList := azureVal.AsValueSlice()
									if len(azureList) > 0 && azureList[0].Type().IsObjectType() {
										azureMap := azureList[0].AsValueMap()
										if clientSecretVal, exists := azureMap["client_secret"]; exists && !clientSecretVal.IsNull() && clientSecretVal.IsKnown() {
											configClientSecret = clientSecretVal.AsString()
										}
										if clientIdVal, exists := azureMap["client_id"]; exists && !clientIdVal.IsNull() && clientIdVal.IsKnown() {
											configClientId = clientIdVal.AsString()
										}
									}
								}
							}
						}
					}
				}

				azure := accessInfo.(config.AzureAccessInformation)
				// If client_secret hasn't changed, use the value from raw config
				if !d.HasChange("access_information.0.azure_key_vault.0.client_secret") && configClientSecret != "" {
					azure.ClientSecret = utils.StringPtr(configClientSecret)
				} else if azure.ClientSecret != nil && utils.StringValue(azure.ClientSecret) == "" && configClientSecret != "" {
					azure.ClientSecret = utils.StringPtr(configClientSecret)
				}
				// client_id
				if !d.HasChange("access_information.0.azure_key_vault.0.client_id") && configClientId != "" {
					azure.ClientId = utils.StringPtr(configClientId)
				} else if azure.ClientId != nil && utils.StringValue(azure.ClientId) == "" && configClientId != "" {
					azure.ClientId = utils.StringPtr(configClientId)
				}
				accessInfo = azure
			}

			if err := updateSpec.SetAccessInformation(accessInfo); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	aJSON, _ := json.MarshalIndent(updateSpec, "", "  ")
	log.Printf("[DEBUG] update key management server payload: %s", aJSON)

	updateResp, err := conn.KeyManagementServersAPIInstance.UpdateKeyManagementServerById(utils.StringPtr(d.Id()), &updateSpec, args)
	if err != nil {
		return diag.Errorf("error while updating key management server : %v", err)
	}
	taskRefValue, ok := updateResp.Data.GetValue().(securityPrism.TaskReference)
	if !ok {
		return diag.Errorf("error: unexpected response type from update API, expected TaskReference")
	}
	TaskRef := taskRefValue
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task
	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the key management server to be updated
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for key management server (%s) to be updated: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get task details for logging
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching key management server update task (%s): %v", utils.StringValue(taskUUID), err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Update Key Management Server Task Details: %s", string(aJSON))

	return ResourceNutanixKeyManagementServerV2Read(ctx, d, meta)
}

func ResourceNutanixKeyManagementServerV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).SecurityAPI

	resp, err := conn.KeyManagementServersAPIInstance.DeleteKeyManagementServerById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while deleting key management server : %v", err)
	}
	taskRefValue, ok := resp.Data.GetValue().(securityPrism.TaskReference)
	if !ok {
		return diag.Errorf("error: unexpected response type from delete API, expected TaskReference")
	}
	TaskRef := taskRefValue
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task
	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the key management server to be deleted
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for key management server (%s) to be deleted: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get task details for logging
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching key management server delete task (%s): %v", utils.StringValue(taskUUID), err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Delete Key Management Server Task Details: %s", string(aJSON))

	return nil
}

func expandAccessInformation(accessInfo []interface{}) (interface{}, error) {
	if len(accessInfo) == 0 {
		return nil, fmt.Errorf("access information is required")
	}

	accessInfoVal, ok := accessInfo[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("access information must be a map")
	}

	azureList, _ := accessInfoVal["azure_key_vault"].([]interface{})
	kmipList, _ := accessInfoVal["kmip_key_vault"].([]interface{})

	if len(azureList) > 0 && len(kmipList) > 0 {
		return nil, fmt.Errorf("exactly one of azure_key_vault or kmip_key_vault must be specified")
	}
	if len(azureList) == 0 && len(kmipList) == 0 {
		return nil, fmt.Errorf("exactly one of azure_key_vault or kmip_key_vault must be specified")
	}

	if len(azureList) > 0 {
		azureMap, ok := azureList[0].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("azure_key_vault must be a map")
		}

		expiryStr := azureMap["credential_expiry_date"].(string)
		expiryTime, err := time.Parse("2006-01-02", expiryStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse credential_expiry_date %q: %w", expiryStr, err)
		}

		azureAccessInfo := config.NewAzureAccessInformation()

		azureAccessInfo.EndpointUrl = utils.StringPtr(azureMap["endpoint_url"].(string))
		azureAccessInfo.KeyId = utils.StringPtr(azureMap["key_id"].(string))
		azureAccessInfo.TenantId = utils.StringPtr(azureMap["tenant_id"].(string))
		azureAccessInfo.ClientId = utils.StringPtr(azureMap["client_id"].(string))
		azureAccessInfo.ClientSecret = utils.StringPtr(azureMap["client_secret"].(string))
		azureAccessInfo.CredentialExpiryDate = utils.Time(expiryTime)

		// The generated API client's OneOf setter expects the non-pointer model type.
		// Returning a pointer here can lead to: "*config.AzureAccessInformation(...) is not expected type".
		return *azureAccessInfo, nil
	}

	kmipMap, ok := kmipList[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("kmip_key_vault must be a map")
	}

	endpoints, err := expandKMIPEndpoints(kmipMap["endpoint_url"])
	if err != nil {
		return nil, err
	}

	kmipAccessInfo := config.NewKmipAccessInformation()
	kmipAccessInfo.CaName = utils.StringPtr(kmipMap["ca_name"].(string))
	kmipAccessInfo.CaPem = utils.StringPtr(kmipMap["ca_pem"].(string))
	kmipAccessInfo.CertPem = utils.StringPtr(kmipMap["cert_pem"].(string))
	kmipAccessInfo.PrivateKey = utils.StringPtr(kmipMap["private_key"].(string))
	kmipAccessInfo.Endpoints = endpoints

	// The generated API client's OneOf setter expects the non-pointer model type.
	return *kmipAccessInfo, nil
}

func expandKMIPEndpoints(raw interface{}) ([]config.EndpointInfo, error) {
	var rawList []interface{}
	switch v := raw.(type) {
	case *schema.Set:
		rawList = v.List()
	case []interface{}:
		rawList = v
	default:
		return nil, fmt.Errorf("endpoint_url must be a set")
	}

	if len(rawList) == 0 {
		return nil, fmt.Errorf("endpoint_url must not be empty")
	}

	endpoints := make([]config.EndpointInfo, 0, len(rawList))
	for _, item := range rawList {
		m, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("endpoint_url entry must be an object")
		}
		ipList, ok := m["ip_address"].([]interface{})
		if !ok || len(ipList) == 0 {
			return nil, fmt.Errorf("endpoint_url.ip_address is required")
		}
		ip := expandIPAddressOrFQDN(ipList)
		if ip == nil || !ip.IsValid() {
			return nil, fmt.Errorf("endpoint_url.ip_address must contain ipv4, ipv6, or fqdn")
		}

		port, ok := m["port"].(int)
		if !ok {
			return nil, fmt.Errorf("endpoint_url.port must be an integer")
		}

		endpoints = append(endpoints, config.EndpointInfo{
			IpAddress: ip,
			Port:      utils.IntPtr(port),
		})
	}
	return endpoints, nil
}

func expandIPAddressOrFQDN(pr []interface{}) *commonCfg.IPAddressOrFQDN {
	if len(pr) == 0 {
		return nil
	}
	val, ok := pr[0].(map[string]interface{})
	if !ok {
		return nil
	}
	ip := commonCfg.NewIPAddressOrFQDN()

	if ipv4, ok := val["ipv4"]; ok && len(ipv4.([]interface{})) > 0 {
		ip.Ipv4 = expandIPv4Address(ipv4)
	}
	if ipv6, ok := val["ipv6"]; ok && len(ipv6.([]interface{})) > 0 {
		ip.Ipv6 = expandIPv6Address(ipv6)
	}
	if fqdn, ok := val["fqdn"]; ok && len(fqdn.([]interface{})) > 0 {
		ip.Fqdn = expandFQDN(fqdn.([]interface{}))
	}

	return ip
}

func expandIPv4Address(pr interface{}) *commonCfg.IPv4Address {
	prSlice, ok := pr.([]interface{})
	if !ok || len(prSlice) == 0 {
		return nil
	}
	valMap, ok := prSlice[0].(map[string]interface{})
	if !ok {
		return nil
	}
	ipv4 := commonCfg.NewIPv4Address()
	if v, ok := valMap["value"].(string); ok {
		ipv4.Value = utils.StringPtr(v)
	}
	if p, ok := valMap["prefix_length"].(int); ok {
		ipv4.PrefixLength = utils.IntPtr(p)
	}
	return ipv4
}

func expandIPv6Address(pr interface{}) *commonCfg.IPv6Address {
	prSlice, ok := pr.([]interface{})
	if !ok || len(prSlice) == 0 {
		return nil
	}
	valMap, ok := prSlice[0].(map[string]interface{})
	if !ok {
		return nil
	}
	ipv6 := commonCfg.NewIPv6Address()
	if v, ok := valMap["value"].(string); ok {
		ipv6.Value = utils.StringPtr(v)
	}
	if p, ok := valMap["prefix_length"].(int); ok {
		ipv6.PrefixLength = utils.IntPtr(p)
	}
	return ipv6
}

func expandFQDN(pr []interface{}) *commonCfg.FQDN {
	if len(pr) == 0 {
		return nil
	}
	valMap, ok := pr[0].(map[string]interface{})
	if !ok {
		return nil
	}
	f := commonCfg.NewFQDN()
	if v, ok := valMap["value"].(string); ok {
		f.Value = utils.StringPtr(v)
	}
	return f
}
