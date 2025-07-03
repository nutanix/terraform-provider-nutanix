package lcmv2

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	lcmCommonConfig "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/common/v1/config"
	"github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/common"
	lcmOperations "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/operations"
	taskRef "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/prism/v4/config"
	lcmSecurityConfig "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/security/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const (
	ipV4PrefixLengthDefault = 32
	ipV6PrefixLengthDefault = 128
)

func ResourceNutanixLcmPerformInventoryV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixLcmPerformInventoryV2Create,
		ReadContext:   ResourceNutanixLcmPerformInventoryV2Read,
		UpdateContext: ResourceNutanixLcmPerformInventoryV2Update,
		DeleteContext: ResourceNutanixLcmPerformInventoryV2Delete,
		Schema: map[string]*schema.Schema{
			"dryrun": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"x_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"credentials": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 4, //nolint:gomnd
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"credential_detail": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1, //nolint:gomnd
							Elem:     credentialDetailSchema(),
						},
					},
				},
			},
		},
	}
}

func ResourceNutanixLcmPerformInventoryV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LcmAPI
	clusterExtID := d.Get("x_cluster_id").(string)

	var dryRun *bool
	var clusterID *string

	if dryRunVar, ok := d.GetOk("dryrun"); ok {
		dryRun = utils.BoolPtr(dryRunVar.(bool))
	} else {
		dryRun = utils.BoolPtr(false)
	}

	if clusterExtID != "" {
		clusterID = utils.StringPtr(clusterExtID)
	} else {
		clusterID = nil
	}

	credentials := &lcmOperations.InventorySpec{}

	if credList, ok := d.GetOk("credentials"); ok {
		credentials.Credentials = expandCredentials(credList.([]interface{}))
	}

	resp, err := conn.LcmInventoryAPIInstance.PerformInventory(credentials, clusterID, dryRun)
	if err != nil {
		return diag.Errorf("error while performing the inventory: %v", err)
	}

	TaskRef := resp.Data.GetValue().(taskRef.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task
	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the inventory to be successful
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroup(taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("Perform inventory task failed: %s", errWaitTask)
	}

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching the Lcm inventory task : %v", err)
	}

	task := resourceUUID.Data.GetValue().(prismConfig.Task)
	aJSON, _ := json.MarshalIndent(task, "", "  ")
	log.Printf("[DEBUG] Perform Inventory Task Response: %s", string(aJSON))

	// randomly generating the id
	d.SetId(utils.GenUUID())
	return nil
}

func ResourceNutanixLcmPerformInventoryV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixLcmPerformInventoryV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return ResourceNutanixLcmPerformInventoryV2Create(ctx, d, meta)
}

func ResourceNutanixLcmPerformInventoryV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func taskStateRefreshPrismTaskGroup(client *prism.Client, taskUUID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		// data := base64.StdEncoding.EncodeToString([]byte("ergon"))
		// encodeUUID := data + ":" + taskUUID
		vresp, err := client.TaskRefAPI.GetTaskById(utils.StringPtr(taskUUID), nil)
		if err != nil {
			return "", "", (fmt.Errorf("error while polling prism task: %v", err))
		}

		// get the group results

		v := vresp.Data.GetValue().(prismConfig.Task)

		if getTaskStatus(v.Status) == "CANCELED" || getTaskStatus(v.Status) == "FAILED" {
			return v, getTaskStatus(v.Status),
				fmt.Errorf("error_detail: %s, progress_message: %d", utils.StringValue(v.ErrorMessages[0].Message), utils.IntValue(v.ProgressPercentage))
		}
		return v, getTaskStatus(v.Status), nil
	}
}

func getTaskStatus(pr *prismConfig.TaskStatus) string {
	const two, three, five, six, seven = 2, 3, 5, 6, 7
	if pr != nil {
		if *pr == prismConfig.TaskStatus(six) {
			return "FAILED"
		}
		if *pr == prismConfig.TaskStatus(seven) {
			return "CANCELED"
		}
		if *pr == prismConfig.TaskStatus(two) {
			return "QUEUED"
		}
		if *pr == prismConfig.TaskStatus(three) {
			return "RUNNING"
		}
		if *pr == prismConfig.TaskStatus(five) {
			return "SUCCEEDED"
		}
	}
	return "UNKNOWN"
}

func expandCredentials(credentialsRawData []interface{}) []common.Credential {
	credentials := make([]common.Credential, 0, len(credentialsRawData))

	for _, raw := range credentialsRawData {
		credentialsRawMap := raw.(map[string]interface{})

		credentialDetailList := credentialsRawMap["credential_detail"].([]interface{})

		if len(credentialDetailList) == 0 {
			continue
		}

		credentialDetailMap := credentialDetailList[0].(map[string]interface{})

		credentialDetailObj := common.NewOneOfCredentialCredentialDetail()
		// dispatch on which block is set:
		if credentialRef, ok := credentialDetailMap["credential_reference"].([]interface{}); ok && len(credentialRef) > 0 {
			ref := expandCredentialReference(credentialRef[0].(map[string]interface{}))
			if err := credentialDetailObj.SetValue(*ref); err != nil {
				log.Printf("[ERROR] error setting credential reference: %s", err)
				continue
			}
		} else if vendorManagementCred, ok := credentialDetailMap["vendor_management_credential"].([]interface{}); ok && len(vendorManagementCred) > 0 {
			vendorManagementMap := vendorManagementCred[0].(map[string]interface{})
			credentialSpecList := vendorManagementMap["credential_spec"].([]interface{})

			if len(credentialSpecList) == 0 {
				log.Printf("[ERROR] credential_spec is not set in vendor_management_credential")
				continue
			}

			credentialSpecMap := credentialSpecList[0].(map[string]interface{})
			credentialSpecObj := common.NewOneOfVendorManagementCredentialCredentialSpec()

			if intersightConnection, ok := credentialSpecMap["intersight_connection"].([]interface{}); ok && len(intersightConnection) > 0 {
				intersightConnMap := intersightConnection[0].(map[string]interface{})

				intersightConnObj := extractIntersightConnection(intersightConnMap)

				if err := credentialSpecObj.SetValue(intersightConnObj); err != nil {
					log.Printf("[ERROR] error setting intersight connection: %s", err)
					continue
				}
			} else if vcenterCredential, ok := credentialSpecMap["vcenter_credential"].([]interface{}); ok && len(vcenterCredential) > 0 {
				vcenterCredMap := vcenterCredential[0].(map[string]interface{})

				vcenterCredObj := extractVCenterCredential(vcenterCredMap)

				if err := credentialSpecObj.SetValue(vcenterCredObj); err != nil {
					log.Printf("[ERROR] error setting vCenter credential: %s", err)
					continue
				}
			} else {
				log.Printf("[ERROR] neither intersight_connection nor vcenter_credential is set in credential_spec")
				continue
			}

			// set the credential spec in the credential detail object
			err := credentialDetailObj.SetValue(*credentialSpecObj)
			if err != nil {
				log.Printf("[ERROR] error setting vendor management credential spec: %s", err)
				continue
			}
		} else {
			log.Printf("[ERROR] neither credential_reference nor vendor_management_credential is set in credential_detail")
			continue
		}
	}
	return credentials
}

func expandCredentialReference(refMap map[string]interface{}) *common.CredentialReference {
	credentialRef := common.NewCredentialReference()
	if extID, ok := refMap["credential_ext_id"].(string); ok && extID != "" {
		credentialRef.CredentialExtId = utils.StringPtr(extID)
	} else {
		log.Printf("[WARN] credential_ext_id is not set or empty in credential_reference")
	}
	return credentialRef
}

func extractIntersightConnection(intersightConnMap map[string]interface{}) *lcmSecurityConfig.IntersightCredential {
	intersightCred := lcmSecurityConfig.NewIntersightCredential()

	if credentialList, ok := intersightConnMap["credential"].([]interface{}); ok && len(credentialList) > 0 {
		credentialMap := credentialList[0].(map[string]interface{})

		credential := lcmSecurityConfig.NewKeyBasedAuth()

		if apiKey, ok := credentialMap["api_key"].(string); ok && apiKey != "" {
			credential.ApiKey = utils.StringPtr(apiKey)
		} else {
			log.Printf("[ERROR] api_key is not set or empty in intersight_connection")
			return nil
		}

		if secretKey, ok := credentialMap["secret_key"].(string); ok && secretKey != "" {
			credential.SecretKey = utils.StringPtr(secretKey)
		} else {
			log.Printf("[ERROR] secret_key is not set or empty in intersight_connection")
			return nil
		}

		intersightCred.Credential = credential
	}

	if url, ok := intersightConnMap["url"].(string); ok && url != "" {
		intersightCred.Url = utils.StringPtr(url)
	}

	if deploymentType, ok := intersightConnMap["deployment_type"].(string); ok && deploymentType != "" {
		if deploymentType == "INTERSIGHT_VIRTUAL_APPLIANCE" {
			intersightCred.DeploymentType = lcmSecurityConfig.INTERSIGHTCONNECTIONTYPE_INTERSIGHT_VIRTUAL_APPLIANCE.Ref()
		} else if deploymentType == "INTERSIGHT_SAAS" {
			intersightCred.DeploymentType = lcmSecurityConfig.INTERSIGHTCONNECTIONTYPE_INTERSIGHT_SAAS.Ref()
		} else {
			log.Printf("[ERROR] invalid deployment_type: %s", deploymentType)
			// return nil
		}
	}

	return intersightCred
}

func extractVCenterCredential(vcenterCredMap map[string]interface{}) *lcmSecurityConfig.VcenterCredential {
	vcenterCred := lcmSecurityConfig.NewVcenterCredential()
	if credentialList, ok := vcenterCredMap["credential"].([]interface{}); ok && len(credentialList) > 0 {
		credentialMap := credentialList[0].(map[string]interface{})
		credential := lcmCommonConfig.NewBasicAuth()
		if username, ok := credentialMap["username"].(string); ok && username != "" {
			credential.Username = utils.StringPtr(username)
		}

		if password, ok := credentialMap["password"].(string); ok && password != "" {
			credential.Password = utils.StringPtr(password)
		}

		vcenterCred.Credential = credential
	}

	if addressList, ok := vcenterCredMap["address"].([]interface{}); ok && len(addressList) > 0 {
		addressMap := addressList[0].(map[string]interface{})
		address := lcmCommonConfig.NewIPAddressOrFQDN()

		if ipv4List, ok := addressMap["ipv4"].([]interface{}); ok && len(ipv4List) > 0 {
			address.Ipv4 = expandIPv4Address(ipv4List[0].(map[string]interface{}))
		}

		if ipv6List, ok := addressMap["ipv6"].([]interface{}); ok && len(ipv6List) > 0 {
			address.Ipv6 = expandIPv6Address(ipv6List[0].(map[string]interface{}))
		}

		if fqdnList, ok := addressMap["fqdn"].([]interface{}); ok && len(fqdnList) > 0 {
			fqdnMap := fqdnList[0].(map[string]interface{})
			fqdn := lcmCommonConfig.NewFQDN()
			if fqdnValue, ok := fqdnMap["value"].(string); ok &&

				fqdnValue != "" {
				fqdn.Value = utils.StringPtr(fqdnValue)
			} else {
				log.Printf("[ERROR] fqdn value is not set or empty in vcenter_credential")
				return nil
			}
			address.Fqdn = fqdn
		} else {
			log.Printf("[ERROR] neither ipv4 nor ipv6 nor fqdn is set in vcenter_credential address")

			return nil
		}
	}

	return vcenterCred
}

func expandIPv4Address(ipv4 map[string]interface{}) *lcmCommonConfig.IPv4Address {
	ipv4Address := lcmCommonConfig.NewIPv4Address()

	if value, ok := ipv4["value"].(string); ok && value != "" {
		ipv4Address.Value = utils.StringPtr(value)
	}

	if prefixLength, ok := ipv4["prefix_length"].(int); ok && prefixLength > 0 {
		ipv4Address.PrefixLength = utils.IntPtr(prefixLength)
	}

	return ipv4Address
}

func expandIPv6Address(ipv6 map[string]interface{}) *lcmCommonConfig.IPv6Address {
	ipv6Address := lcmCommonConfig.NewIPv6Address()

	if value, ok := ipv6["value"].(string); ok && value != "" {
		ipv6Address.Value = utils.StringPtr(value)
	}

	if prefixLength, ok := ipv6["prefix_length"].(int); ok && prefixLength > 0 {
		ipv6Address.PrefixLength = utils.IntPtr(prefixLength)
	}

	return ipv6Address
}

func credentialDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"credential_reference": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"credential_ext_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				// exactly one of reference OR vendor_management must be set
				// ConflictsWith: []string{"vendor_management_credential"},
				// ExactlyOneOf:  []string{"credential_reference", "vendor_management_credential"},
			},
			"vendor_management_credential": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     vendorManagementCredentialSchema(),
				// exactly one of reference OR vendor_management must be set
				// ConflictsWith: []string{"credential_reference"},
				// ExactlyOneOf:  []string{"credential_reference", "vendor_management_credential"},
			},
		},
	}
}

func vendorManagementCredentialSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"credential_spec": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"intersight_connection": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1, //nolint:gomnd
							Elem:     interSightConnectionSchema(),
						},
						"vcenter_credential": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1, //nolint:gomnd
							Elem:     vcenterCredentialSchema(),
						},
					},
				},
			},
		},

		// enforce exactly one of intersight vs vcenter inside the spec
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			specs := d.Get("credential_spec").([]interface{})
			if len(specs) == 0 {
				return fmt.Errorf("credential_spec must be set")
			}
			spec := specs[0].(map[string]interface{})
			_, hasIS := spec["intersight_connection"]
			_, hasVC := spec["vcenter_credential"]

			if hasIS == hasVC {
				return fmt.Errorf("exactly one of 'intersight_connection' or 'vcenter_credential' must be set")
			}
			return nil
		},
	}
}

func interSightConnectionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"credential": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1, //nolint:gomnd
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_key": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"secret_key": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
					},
				},
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"deployment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"INTERSIGHT_VIRTUAL_APPLIANCE", "INTERSIGHT_SAAS"}, false),
			},
		},
	}
}

func vcenterCredentialSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"credential": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1, //nolint:gomnd
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Required: true,
						},
						"password": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
					},
				},
			},
			"address": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1, //nolint:gomnd
				Elem:     ipAddressOrFqdnSchema(),
			},
		},
	}
}

func ipAddressOrFqdnSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ipv4": {
				Type:         schema.TypeList,
				Optional:     true,
				Computed:     true,
				MaxItems:     1, //nolint:gomnd
				// ExactlyOneOf: []string{"ipv6", "fqdn"},
				Elem:         ipSchema(ipV4PrefixLengthDefault),
			},
			"ipv6": {
				Type:         schema.TypeList,
				Optional:     true,
				Computed:     true,
				MaxItems:     1, //nolint:gomnd
				// ExactlyOneOf: []string{"ipv4", "fqdn"},
				Elem:         ipSchema(ipV6PrefixLengthDefault),
			},
			"fqdn": {
				Type:         schema.TypeList,
				Optional:     true,
				Computed:     true,
				MaxItems:     1, //nolint:gomnd
				// ExactlyOneOf: []string{"ipv4", "fqdn"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func ipSchema(defaultPrefix int) *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"prefix_length": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      defaultPrefix,
				ValidateFunc: validation.IntBetween(0, defaultPrefix),
			},
		},
	}
}
