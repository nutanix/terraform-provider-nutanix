package clustersv2

import (
	"context"
	"encoding/json"
	"log"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/common/v1/config"
	import3 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/prism/v4/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixClusterProfileV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixClusterProfileV2Create,
		ReadContext:   ResourceNutanixClusterProfileV2Read,
		UpdateContext: ResourceNutanixClusterProfileV2Update,
		DeleteContext: ResourceNutanixClusterProfileV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"allowed_overrides": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(AllowedOverridesStrings, false),
				},
			},
			"name_server_ip_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Set:      common.HashIPItem,
				Elem:     common.SchemaForIPList(false), // do not include FQDN
			},
			"ntp_server_ip_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Set:      common.HashIPItem,
				Elem:     common.SchemaForIPList(true), // include FQDN
			},
			"smtp_server": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email_address": {
							Type:     schema.TypeString,
							Required: true,
						},
						"server": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_address": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem:     common.SchemaForIPList(true), // include FQDN
									},
									"port": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"username": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"password": {
										Type:      schema.TypeString,
										Sensitive: true,
										Optional:  true,
									},
								},
							},
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice(SMTPTypeStrings, false),
						},
					},
				},
			},
			"nfs_subnet_white_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringMatch(regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}/(?:\d{1,3}\.){3}\d{1,3}\b`),
						"Must be in CIDR-like format x.x.x.x/y.y.y.y"),
				},
			},
			"snmp_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"links": common.LinksSchema(),
						"is_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"users": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ext_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"tenant_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"links": common.LinksSchema(),
									"username": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringLenBetween(1, 64),
									},
									"auth_type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(SnmpAuthTypeStrings, false),
									},
									"auth_key": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[^']+$`), "cannot contain single quotes"),
									},
									"priv_type": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice(SnmpPrivTypeStrings, false),
									},
									"priv_key": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[^']+$`), "cannot contain single quotes"),
									},
								},
							},
						},
						"transports": {
							Type:     schema.TypeList,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(SnmpProtocolStrings, false),
									},
									"port": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						"traps": {
							Type:     schema.TypeList,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ext_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"tenant_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"links": common.LinksSchema(),
									"address": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem:     common.SchemaForIPList(false), // do not include FQDN
									},
									"username": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringLenBetween(1, 64),
									},
									"protocol": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice(SnmpProtocolStrings, false),
									},
									"port": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"should_inform": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"engine_id": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringMatch(regexp.MustCompile(`^(?:0[xX])?[0-9a-fA-F]+$`), "must be a valid hex string"),
									},
									"version": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(SnmpTrapVersionStrings, false),
									},
									"receiver_name": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringLenBetween(1, 64),
									},
									"community_string": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"rsyslog_server_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"links": common.LinksSchema(),
						"server_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 64),
						},

						"ip_address": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem:     common.SchemaForIPList(false), // do not include FQDN
						},
						"port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"network_protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(RsyslogNetworkProtocolStrings, false),
						},
						"modules": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(RsyslogModuleNameStrings, false),
									},
									"log_severity_level": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(RsyslogLogSeverityLevelStrings, false),
									},
									"should_log_monitor_files": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
								},
							},
						},
					},
				},
			},
			"pulse_status": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"pii_scrubbing_level": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice(PIIScrubbingLevelStrings, false),
						},
					},
				},
			},
			"dryrun": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			// computed fields
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": common.LinksSchema(),
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"drifted_cluster_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_compliant": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"last_synced_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"config_drifts": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func ResourceNutanixClusterProfileV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI
	body := expandClusterProfile(d)

	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] Cluster Profile Create Payload: %s", string(aJSON))

	createResp, createErr := conn.ClusterProfilesAPI.CreateClusterProfile(body)
	if createErr != nil {
		return diag.FromErr(createErr)
	}

	TaskRef := createResp.Data.GetValue().(import3.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the cluster profile to be created
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}
	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for cluster profile (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	// Get UUID from TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching cluster profile task: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(import2.Task)
	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Create Cluster Profile Task Details: %s", string(aJSON))

	// Get UUID from TASK API using entity type constant
	uuid, err := common.ExtractEntityUUIDFromTask(taskDetails, utils.RelEntityTypeClusterProfile, "Cluster profile")
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(uuid))

	return ResourceNutanixClusterProfileV2Read(ctx, d, meta)
}

func ResourceNutanixClusterProfileV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI

	// Fetch the Cluster Profile by UUID
	clusterProfileResp, err := conn.ClusterProfilesAPI.GetClusterProfileById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.FromErr(err)
	}

	// Check if Data is nil or empty
	if clusterProfileResp.Data == nil || clusterProfileResp.Data.GetValue() == nil {
		return diag.Errorf("ClusterProfile API returned empty data for ID %s", d.Id())
	}

	// Safe type assertion
	clusterProfile, ok := clusterProfileResp.Data.GetValue().(config.ClusterProfile)
	if !ok {
		return diag.Errorf("ClusterProfile API returned unexpected type for ID %s", d.Id())
	}

	// Set the resource data from the API response
	if err := d.Set("tenant_id", clusterProfile.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ext_id", clusterProfile.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", common.FlattenLinks(clusterProfile.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", clusterProfile.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("create_time", utils.TimeStringValue(clusterProfile.CreateTime)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("last_update_time", utils.TimeStringValue(clusterProfile.LastUpdateTime)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("created_by", utils.StringValue(clusterProfile.CreatedBy)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("last_updated_by", utils.StringValue(clusterProfile.LastUpdatedBy)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_count", utils.IntValue(clusterProfile.ClusterCount)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("drifted_cluster_count", utils.IntValue(clusterProfile.DriftedClusterCount)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("clusters", flattenClusterProfileClusterList(clusterProfile.Clusters)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", clusterProfile.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("allowed_overrides", common.FlattenEnumValueList(clusterProfile.AllowedOverrides)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name_server_ip_list", flattenIPAddressList(clusterProfile.NameServerIpList)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ntp_server_ip_list", flattenIPAddressOrFQDN(clusterProfile.NtpServerIpList)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("smtp_server", flattenSMTPServerRef(clusterProfile.SmtpServer)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("nfs_subnet_white_list", clusterProfile.NfsSubnetWhitelist); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("snmp_config", flattenSnmpConfig(clusterProfile.SnmpConfig)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("rsyslog_server_list", flattenRsyslogServerList(clusterProfile.RsyslogServerList)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("pulse_status", flattenPulseStatus(clusterProfile.PulseStatus)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ResourceNutanixClusterProfileV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI

	// Fetch the Cluster Profile by UUID
	clusterProfileResp, err := conn.ClusterProfilesAPI.GetClusterProfileById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.FromErr(err)
	}

	// Check if Data is nil or empty
	if clusterProfileResp.Data == nil || clusterProfileResp.Data.GetValue() == nil {
		return diag.Errorf("ClusterProfile API returned empty data for ID %s", d.Id())
	}

	// Safe type assertion
	clusterProfile, ok := clusterProfileResp.Data.GetValue().(config.ClusterProfile)
	if !ok {
		return diag.Errorf("ClusterProfile API returned unexpected type for ID %s", d.Id())
	}

	body := config.NewClusterProfile()

	if d.HasChange("name") {
		body.Name = utils.StringPtr(d.Get("name").(string))
	} else {
		// name is required in update payload
		body.Name = clusterProfile.Name
	}
	if d.HasChange("description") {
		body.Description = utils.StringPtr(d.Get("description").(string))
	}
	if d.HasChange("allowed_overrides") {
		aoList, _ := d.GetOk("allowed_overrides")
		body.AllowedOverrides = common.ExpandEnumList(aoList, AllowedOverridesMap, "allowed_override")
	}
	if d.HasChange("name_server_ip_list") {
		nameServerIPRaw, _ := d.GetOk("name_server_ip_list")
		nameServerIPList := common.InterfaceToSlice(nameServerIPRaw)
		result := make([]import1.IPAddress, 0)
		for _, ip := range nameServerIPList {
			result = append(result, *expandIPAddress(common.InterfaceToSlice(ip)))
		}
		body.NameServerIpList = result
	}
	if d.HasChange("ntp_server_ip_list") {
		ntpServerIPRaw, _ := d.GetOk("ntp_server_ip_list")
		body.NtpServerIpList = expandIPAddressOrFQDN(common.InterfaceToSlice(ntpServerIPRaw))
	}
	if d.HasChange("smtp_server") {
		smtpConfigRaw, _ := d.GetOk("smtp_server")
		smtpConfigList := common.InterfaceToSlice(smtpConfigRaw)
		body.SmtpServer = expandSMTPServerRef(smtpConfigList)
	}
	if d.HasChange("nfs_subnet_white_list") {
		nfsWhiteListRaw, _ := d.GetOk("nfs_subnet_white_list")
		nfsWhiteList := common.InterfaceToSlice(nfsWhiteListRaw)
		body.NfsSubnetWhitelist = common.ExpandListOfString(nfsWhiteList)
	}
	if d.HasChange("snmp_config") {
		snmpConfigRaw, _ := d.GetOk("snmp_config")
		snmpConfigList := common.InterfaceToSlice(snmpConfigRaw)
		body.SnmpConfig = expandSNMPConfig(snmpConfigList)
	}
	if d.HasChange("rsyslog_server_list") {
		rsyslogServerListRaw, _ := d.GetOk("rsyslog_server_list")
		rsyslogServerList := common.InterfaceToSlice(rsyslogServerListRaw)
		body.RsyslogServerList = expandRsyslogServerList(rsyslogServerList)
	}
	if d.HasChange("pulse_status") {
		pulseStatusRaw, _ := d.GetOk("pulse_status")
		pulseStatusList := common.InterfaceToSlice(pulseStatusRaw)
		body.PulseStatus = expandPulseStatus(pulseStatusList)
	}

	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] Cluster Profile Update Payload: %s", string(aJSON))

	var dryRun *bool
	if v, ok := d.GetOk("dryrun"); ok {
		dryRun = utils.BoolPtr(v.(bool))
	}

	// Extract E-tag from read response
	getClusterProfileResp, getErr := conn.ClusterProfilesAPI.GetClusterProfileById(utils.StringPtr(d.Id()))
	if getErr != nil {
		return diag.Errorf("error fetching cluster profile: %v", getErr)
	}

	// Extract E-Tag Header
	etagValue := conn.ClusterProfilesAPI.ApiClient.GetEtag(getClusterProfileResp)
	args := make(map[string]interface{})
	args["If-Match"] = utils.StringPtr(etagValue)

	// Update Cluster Profile
	updateResp, updateErr := conn.ClusterProfilesAPI.UpdateClusterProfileById(utils.StringPtr(d.Id()), body, dryRun, args)
	if updateErr != nil {
		return diag.FromErr(updateErr)
	}

	TaskRef := updateResp.Data.GetValue().(import3.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the cluster profile to be updated
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}
	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for cluster profile (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	// Get UUID from TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching cluster profile task: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(import2.Task)
	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Update Cluster Profile Task Details: %s", string(aJSON))

	return ResourceNutanixClusterProfileV2Read(ctx, d, meta)
}

func ResourceNutanixClusterProfileV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI

	deleteResp, deleteErr := conn.ClusterProfilesAPI.DeleteClusterProfileById(utils.StringPtr(d.Id()))
	if deleteErr != nil {
		return diag.FromErr(deleteErr)
	}

	TaskRef := deleteResp.Data.GetValue().(import3.TaskReference)
	taskUUID := TaskRef.ExtId
	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the cluster profile to be deleted
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}
	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for cluster profile (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	// Get UUID from TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching cluster profile delete task: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(import2.Task)
	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Delete Cluster Profile Task Details: %s", string(aJSON))

	return nil
}

// ###########################################
// ###### Cluster Profiles Expanders #########
// ###########################################

// expandClusterProfile expands the Cluster Profile resource data into the API model
func expandClusterProfile(d *schema.ResourceData) *config.ClusterProfile {
	body := config.NewClusterProfile()

	// Simple string fields
	if name, ok := d.GetOk("name"); ok {
		body.Name = utils.StringPtr(name.(string))
	}
	if description, ok := d.GetOk("description"); ok {
		body.Description = utils.StringPtr(description.(string))
	}

	// Enum list
	if aoList, ok := d.GetOk("allowed_overrides"); ok {
		body.AllowedOverrides = common.ExpandEnumList(aoList, AllowedOverridesMap, "allowed_override")
	}

	// Name server IP list
	if nameServerIPRaw, ok := d.GetOk("name_server_ip_list"); ok {
		nameServerIPList := common.InterfaceToSlice(nameServerIPRaw)
		result := make([]import1.IPAddress, 0)
		for _, ip := range nameServerIPList {
			result = append(result, *expandIPAddress(common.InterfaceToSlice(ip)))
		}
		body.NameServerIpList = result
	}

	// NTP server IPs
	if ntpServerIPRaw, ok := d.GetOk("ntp_server_ip_list"); ok {
		body.NtpServerIpList = expandIPAddressOrFQDN(common.InterfaceToSlice(ntpServerIPRaw))
	}

	// SMTP server
	if smtpConfigRaw, ok := d.GetOk("smtp_server"); ok {
		smtpConfigList := common.InterfaceToSlice(smtpConfigRaw)
		body.SmtpServer = expandSMTPServerRef(smtpConfigList)
	}

	// NFS subnet whitelist
	if nfsWhiteListRaw, ok := d.GetOk("nfs_subnet_white_list"); ok {
		nfsWhiteList := common.InterfaceToSlice(nfsWhiteListRaw)
		body.NfsSubnetWhitelist = common.ExpandListOfString(nfsWhiteList)
	}

	// SNMP config
	if snmpConfigRaw, ok := d.GetOk("snmp_config"); ok {
		snmpConfigList := common.InterfaceToSlice(snmpConfigRaw)
		body.SnmpConfig = expandSNMPConfig(snmpConfigList)
	}

	// Rsyslog servers
	if rsyslogServerListRaw, ok := d.GetOk("rsyslog_server_list"); ok {
		rsyslogServerList := common.InterfaceToSlice(rsyslogServerListRaw)
		body.RsyslogServerList = expandRsyslogServerList(rsyslogServerList)
	}

	// Pulse status
	if pulseStatusRaw, ok := d.GetOk("pulse_status"); ok {
		pulseStatusList := common.InterfaceToSlice(pulseStatusRaw)
		body.PulseStatus = expandPulseStatus(pulseStatusList)
	}

	return body
}

// expandSNMPConfig expands SNMP configuration from the resource data
func expandSNMPConfig(snmpConfigList []interface{}) *config.SnmpConfig {
	if len(snmpConfigList) == 0 || snmpConfigList[0] == nil {
		return nil
	}

	raw := snmpConfigList[0].(map[string]interface{})
	snmp := &config.SnmpConfig{}

	if v, ok := raw["is_enabled"]; ok {
		snmp.IsEnabled = utils.BoolPtr(v.(bool))
	}

	// Users
	if usersRaw, ok := raw["users"]; ok {
		users := common.InterfaceToSlice(usersRaw)
		snmp.Users = make([]config.SnmpUser, 0, len(users))
		for _, u := range users {
			userMap := u.(map[string]interface{})
			user := config.SnmpUser{
				Username: utils.StringPtr(userMap["username"].(string)),
				AuthType: common.ExpandEnum(userMap["auth_type"].(string), SnmpAuthTypeMap, "auth_type"),
			}

			// Only set AuthKey if it exists and is non-empty
			if v, ok := userMap["auth_key"].(string); ok && v != "" {
				user.AuthKey = utils.StringPtr(v)
			}

			// Only set PrivType if it exists and is non-empty
			if v, ok := userMap["priv_type"].(string); ok && v != "" {
				user.PrivType = common.ExpandEnum(v, SnmpPrivTypeMap, "priv_type")
			}

			// Only set PrivKey if it exists and is non-empty
			if v, ok := userMap["priv_key"].(string); ok && v != "" {
				user.PrivKey = utils.StringPtr(v)
			}

			snmp.Users = append(snmp.Users, user)
		}
	}

	// Transports
	if transportsRaw, ok := raw["transports"]; ok {
		transports := common.InterfaceToSlice(transportsRaw)
		snmp.Transports = make([]config.SnmpTransport, 0, len(transports))
		for _, t := range transports {
			tMap := t.(map[string]interface{})
			transport := config.SnmpTransport{
				Protocol: common.ExpandEnum(tMap["protocol"].(string), SnmpProtocolMap, "protocol"),
				Port:     utils.IntPtr(int(tMap["port"].(int))),
			}
			snmp.Transports = append(snmp.Transports, transport)
		}
	}

	// Traps
	if trapsRaw, ok := raw["traps"]; ok {
		traps := common.InterfaceToSlice(trapsRaw)
		snmp.Traps = make([]config.SnmpTrap, 0, len(traps))
		for _, tr := range traps {
			trMap := tr.(map[string]interface{})
			trap := config.SnmpTrap{
				Address:         expandIPAddress(trMap["address"]),
				Username:        utils.StringPtr(trMap["username"].(string)),
				Protocol:        common.ExpandEnum(trMap["protocol"].(string), SnmpProtocolMap, "protocol"),
				Port:            utils.IntPtr(int(trMap["port"].(int))),
				ShouldInform:    utils.BoolPtr(trMap["should_inform"].(bool)),
				EngineId:        utils.StringPtr(trMap["engine_id"].(string)),
				Version:         common.ExpandEnum(trMap["version"].(string), SnmpTrapVersionMap, "version"),
				RecieverName:    utils.StringPtr(trMap["receiver_name"].(string)),
				CommunityString: utils.StringPtr(trMap["community_string"].(string)),
			}
			snmp.Traps = append(snmp.Traps, trap)
		}
	}

	return snmp
}

// expandRsyslogServerList expands the Rsyslog server list from the resource data
func expandRsyslogServerList(rsyslogServerList []interface{}) []config.RsyslogServer {
	if len(rsyslogServerList) == 0 {
		return nil
	}

	result := make([]config.RsyslogServer, 0, len(rsyslogServerList))

	for _, item := range rsyslogServerList {
		serverMap := item.(map[string]interface{})
		server := config.RsyslogServer{
			ServerName:      utils.StringPtr(serverMap["server_name"].(string)),
			Port:            utils.IntPtr(int(serverMap["port"].(int))),
			NetworkProtocol: common.ExpandEnum(serverMap["network_protocol"].(string), RsyslogNetworkProtocolMap, "network_protocol"),
		}

		// IP Address
		if ipRaw, ok := serverMap["ip_address"]; ok {
			ipList := common.InterfaceToSlice(ipRaw)
			if len(ipList) > 0 && ipList[0] != nil {
				ipMap := ipList[0].(map[string]interface{})
				server.IpAddress = expandIPAddress(common.InterfaceToSlice(ipMap))
			}
		}

		// Modules
		if modulesRaw, ok := serverMap["modules"]; ok {
			modules := common.InterfaceToSlice(modulesRaw)
			server.Modules = make([]config.RsyslogModuleItem, 0, len(modules))
			for _, m := range modules {
				modMap := m.(map[string]interface{})
				module := config.RsyslogModuleItem{
					Name:                  common.ExpandEnum(modMap["name"].(string), RsyslogModuleNameMap, "name"),
					LogSeverityLevel:      common.ExpandEnum(modMap["log_severity_level"].(string), RsyslogLogSeverityLevelMap, "log_severity_level"),
					ShouldLogMonitorFiles: utils.BoolPtr(modMap["should_log_monitor_files"].(bool)),
				}
				server.Modules = append(server.Modules, module)
			}
		}

		result = append(result, server)
	}

	return result
}

// ############################################
// ###### Cluster Profiles Flatteners #########
// ############################################

// flattenSnmpConfig flattens the SNMP configuration into the resource data format
func flattenSnmpConfig(snmpConfig *config.SnmpConfig) interface{} {
	if snmpConfig == nil {
		return nil
	}

	m := map[string]interface{}{}

	// Flatten ext_id
	if snmpConfig.ExtId != nil {
		m["ext_id"] = utils.StringValue(snmpConfig.ExtId)
	}

	// Flatten tenant_id
	if snmpConfig.TenantId != nil {
		m["tenant_id"] = utils.StringValue(snmpConfig.TenantId)
	}

	// Flatten links
	if snmpConfig.Links != nil {
		links := common.FlattenLinks(snmpConfig.Links)
		m["links"] = links
	}
	// Flatten is_enabled
	if snmpConfig.IsEnabled != nil {
		m["is_enabled"] = utils.BoolValue(snmpConfig.IsEnabled)
	}

	// Flatten users
	if len(snmpConfig.Users) > 0 {
		users := make([]interface{}, 0, len(snmpConfig.Users))
		for _, u := range snmpConfig.Users {
			user := map[string]interface{}{
				"ext_id":    utils.StringValue(u.ExtId),
				"tenant_id": utils.StringValue(u.TenantId),
				"links":     common.FlattenLinks(u.Links),
				"username":  u.Username,
				"auth_type": common.FlattenPtrEnum(u.AuthType),
				"auth_key":  utils.StringValue(u.AuthKey),
				"priv_type": common.FlattenPtrEnum(u.PrivType),
				"priv_key":  utils.StringValue(u.PrivKey),
			}
			users = append(users, user)
		}
		m["users"] = users
	}

	// Flatten transports
	if len(snmpConfig.Transports) > 0 {
		transports := make([]interface{}, 0, len(snmpConfig.Transports))
		for _, t := range snmpConfig.Transports {
			trans := map[string]interface{}{
				"protocol": common.FlattenPtrEnum(t.Protocol),
				"port":     utils.IntValue(t.Port),
			}
			transports = append(transports, trans)
		}
		m["transports"] = transports
	}

	// Flatten traps
	if len(snmpConfig.Traps) > 0 {
		traps := make([]interface{}, 0, len(snmpConfig.Traps))
		for _, tr := range snmpConfig.Traps {
			trap := map[string]interface{}{
				"ext_id":           utils.StringValue(tr.ExtId),
				"tenant_id":        utils.StringValue(tr.TenantId),
				"links":            common.FlattenLinks(tr.Links),
				"username":         tr.Username,
				"protocol":         common.FlattenPtrEnum(tr.Protocol),
				"port":             utils.IntValue(tr.Port),
				"should_inform":    utils.BoolValue(tr.ShouldInform),
				"engine_id":        utils.StringValue(tr.EngineId),
				"version":          common.FlattenPtrEnum(tr.Version),
				"receiver_name":    tr.RecieverName,
				"community_string": utils.StringValue(tr.CommunityString),
			}

			// Flatten IP address if available
			if tr.Address != nil {
				ip := flattenIPAddress(tr.Address)
				trap["address"] = ip
			}

			traps = append(traps, trap)
		}
		m["traps"] = traps
	}

	return []interface{}{m}
}

// flattenRsyslogServerList flattens the Rsyslog server list into the resource data format
func flattenRsyslogServerList(rsyslogServers []config.RsyslogServer) interface{} {
	if len(rsyslogServers) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(rsyslogServers))
	for _, srv := range rsyslogServers {
		s := map[string]interface{}{
			"ext_id":           utils.StringValue(srv.ExtId),
			"tenant_id":        utils.StringValue(srv.TenantId),
			"links":            common.FlattenLinks(srv.Links),
			"server_name":      utils.StringValue(srv.ServerName),
			"port":             utils.IntValue(srv.Port),
			"network_protocol": common.FlattenPtrEnum(srv.NetworkProtocol),
		}

		// Flatten IP if present
		if srv.IpAddress != nil {
			ip := flattenIPAddress(srv.IpAddress)
			s["ip_address"] = ip
		}

		// Flatten modules
		if len(srv.Modules) > 0 {
			modules := make([]interface{}, 0, len(srv.Modules))
			for _, mod := range srv.Modules {
				m := map[string]interface{}{
					"name":                     common.FlattenPtrEnum(mod.Name),
					"log_severity_level":       common.FlattenPtrEnum(mod.LogSeverityLevel),
					"should_log_monitor_files": utils.BoolValue(mod.ShouldLogMonitorFiles),
				}
				modules = append(modules, m)
			}
			s["modules"] = modules
		}

		result = append(result, s)
	}

	return result
}
