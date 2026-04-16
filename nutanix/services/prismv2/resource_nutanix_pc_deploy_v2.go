package prismv2

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	clustermgmtConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/clustermgmt/v4/config"
	commonConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/common/v1/config"
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	vmmConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/vmm/v4/ahv/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const (
	ipV4PrefixLengthDefault = 32
	ipV6PrefixLengthDefault = 128
)

func ResourceNutanixDeployPcV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixDeployPcV2Create,
		ReadContext:   ResourceNutanixDeployPcV2Read,
		UpdateContext: ResourceNutanixDeployPcV2Update,
		DeleteContext: ResourceNutanixDeployPcV2Delete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1*time.Hour + 30*time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"config": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     schemaForPcConfig(),
			},
			"network": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     schemaForPcNetwork(),
			},
			"should_enable_high_availability": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func ResourceNutanixDeployPcV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).PrismAPI

	deployPcBody := config.NewDomainManager()

	if configData, ok := d.GetOk("config"); ok {
		deployPcBody.Config = expandPCConfig(configData)
	}
	if networkData, ok := d.GetOk("network"); ok {
		deployPcBody.Network = expandPCNetwork(networkData)
	}
	if shouldEnableHighAvailability, ok := d.GetOk("should_enable_high_availability"); ok {
		deployPcBody.ShouldEnableHighAvailability = utils.BoolPtr(shouldEnableHighAvailability.(bool))
	}

	aJSON, _ := json.MarshalIndent(deployPcBody, "", "  ")
	log.Printf("[DEBUG] Payload to deploy PC: %s", string(aJSON))

	resp, err := conn.DomainManagerAPIInstance.CreateDomainManager(deployPcBody)
	if err != nil {
		return diag.Errorf("error while deploying PC: %s", err)
	}

	TaskRef := resp.Data.GetValue().(config.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the PC to be deployed
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for PC (%s) to be deployed: %s", utils.StringValue(taskUUID), err)
	}

	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching PC deploy task (%s): %s", utils.StringValue(taskUUID), err)
	}

	taskDetails := taskResp.Data.GetValue().(config.Task)
	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Deploy PC task details: %s", string(aJSON))

	d.SetId(utils.GenUUID())

	return nil
}

func ResourceNutanixDeployPcV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixDeployPcV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return ResourceNutanixDeployPcV2Create(ctx, d, meta)
}

func ResourceNutanixDeployPcV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

// schema Functions

func schemaForPcConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"should_enable_lockdown_mode": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"build_info": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"size": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"SMALL", "LARGE", "EXTRALARGE", "STARTER"}, false),
			},
			"bootstrap_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_init_config": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"datasource_type": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "CONFIG_DRIVE_V2",
										ValidateFunc: validation.StringInSlice([]string{"CONFIG_DRIVE_V2"}, false),
									},
									"metadata": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"cloud_init_script": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"user_data": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Optional: true,
																Computed: true,
															},
														},
													},
												},
												"custom_key_values": schemaForCustomKeyValuePairs(),
											},
										},
									},
								},
							},
						},
						"environment_info": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"NTNX_CLOUD", "ONPREM"}, false),
									},
									"provider_type": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"VSPHERE", "AZURE", "NTNX", "GCP", "AWS"}, false),
									},
									"provisioning_type": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"NATIVE", "NTNX"}, false),
									},
								},
							},
						},
					},
				},
			},
			"credentials": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				MaxItems: 5, //nolint:gomnd
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
			"resource_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"container_ext_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MinItems: 1,
							MaxItems: 3, //nolint:gomnd
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"num_vcpus": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"memory_size_bytes": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"data_disk_size_bytes": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func schemaForCustomKeyValuePairs() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key_value_pairs": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					MinItems: 0,
					MaxItems: 32, //nolint:gomnd
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name": {
								Type:     schema.TypeString,
								Optional: true,
								Computed: true,
							},
							"value": schemaForValue(),
						},
					},
				},
			},
		},
	}
}

func schemaForValue() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"string": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"integer": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},
				"boolean": {
					Type:     schema.TypeBool,
					Optional: true,
					Computed: true,
				},
				"string_list": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"object": {
					Type:     schema.TypeMap,
					Optional: true,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"map_of_strings": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"map": {
								Type:     schema.TypeMap,
								Optional: true,
								Computed: true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
						},
					},
				},
				"integer_list": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeInt,
					},
				},
			},
		},
	}
}

func schemaForPcNetwork() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"external_address": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     schemaForIPAddress(),
			},
			"name_servers": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1024, //nolint:gomnd
				Elem:     schemaForIPAddressOrFqdn(),
			},
			"ntp_servers": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1024, //nolint:gomnd
				Elem:     schemaForIPAddressOrFqdn(),
			},
			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"internal_networks": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default_gateway": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem:     schemaForIPAddressOrFqdn(),
						},
						"subnet_mask": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem:     schemaForIPAddressOrFqdn(),
						},
						"ip_ranges": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 15, //nolint:gomnd
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"begin": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem:     schemaForIPAddress(),
									},
									"end": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem:     schemaForIPAddress(),
									},
								},
							},
						},
					},
				},
			},
			"external_networks": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default_gateway": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem:     schemaForIPAddressOrFqdn(),
						},
						"subnet_mask": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem:     schemaForIPAddressOrFqdn(),
						},
						"ip_ranges": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 15, //nolint:gomnd
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"begin": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem:     schemaForIPAddress(),
									},
									"end": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem:     schemaForIPAddress(),
									},
								},
							},
						},
						"network_ext_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func schemaForIPAddress() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ipv4": SchemaForValuePrefixLengthResource(ipV4PrefixLengthDefault),
			"ipv6": SchemaForValuePrefixLengthResource(ipV6PrefixLengthDefault),
		},
	}
}

func schemaForIPAddressOrFqdn() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ipv4": SchemaForValuePrefixLengthResource(ipV4PrefixLengthDefault),
			"ipv6": SchemaForValuePrefixLengthResource(ipV6PrefixLengthDefault),
			"fqdn": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
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

func SchemaForValuePrefixLengthResource(defaultPrefixLength int) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"value": {
					Type:     schema.TypeString,
					Required: true,
				},
				"prefix_length": {
					Type:         schema.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(0, defaultPrefixLength),
					Default:      defaultPrefixLength,
				},
			},
		},
	}
}

// Expanders
// Pc Config Expanders
func expandPCConfig(configData interface{}) *config.DomainManagerClusterConfig {
	if len(configData.([]interface{})) == 0 {
		return nil
	}
	configDataInterface := configData.([]interface{})
	configDataMap := configDataInterface[0].(map[string]interface{})

	domainManagerClusterConfig := config.NewDomainManagerClusterConfig()

	if shouldEnableLockdownMode, ok := configDataMap["should_enable_lockdown_mode"]; ok {
		domainManagerClusterConfig.ShouldEnableLockdownMode = utils.BoolPtr(shouldEnableLockdownMode.(bool))
	}
	if buildInfo, ok := configDataMap["build_info"]; ok {
		buildInfoData := buildInfo.([]interface{})[0].(map[string]interface{})
		buildInfoObj := clustermgmtConfig.NewBuildInfo()
		if version, ok := buildInfoData["version"]; ok {
			buildInfoObj.Version = utils.StringPtr(version.(string))
		}
		domainManagerClusterConfig.BuildInfo = buildInfoObj
	}
	if name, ok := configDataMap["name"]; ok {
		domainManagerClusterConfig.Name = utils.StringPtr(name.(string))
	}
	if size, ok := configDataMap["size"]; ok {
		domainManagerClusterConfig.Size = expandClusterSize(size.(string))
	}
	if bootstrapConfig, ok := configDataMap["bootstrap_config"]; ok {
		domainManagerClusterConfig.BootstrapConfig = expandBootstrapConfig(bootstrapConfig)
	}
	if credentials, ok := configDataMap["credentials"]; ok {
		domainManagerClusterConfig.Credentials = expandCredentials(credentials.([]interface{}))
	}
	if resourceConfig, ok := configDataMap["resource_config"]; ok {
		domainManagerClusterConfig.ResourceConfig = expandResourceConfig(resourceConfig)
	}

	return domainManagerClusterConfig
}

func expandClusterSize(size string) *config.Size {
	const STARTER, SMALL, LARGE, EXTRALARGE = 2, 3, 4, 5
	switch size {
	case "STARTER":
		sizeVal := config.Size(STARTER)
		return &sizeVal
	case "SMALL":
		sizeVal := config.Size(SMALL)
		return &sizeVal
	case "LARGE":
		sizeVal := config.Size(LARGE)
		return &sizeVal
	case "EXTRALARGE":
		sizeVal := config.Size(EXTRALARGE)
		return &sizeVal
	default:
		return nil
	}
}

func expandBootstrapConfig(bootStrapConfigData interface{}) *config.BootstrapConfig {
	if len(bootStrapConfigData.([]interface{})) == 0 {
		return nil
	}
	bootStrapConfigDataInterface := bootStrapConfigData.([]interface{})
	bootStrapConfigDataMap := bootStrapConfigDataInterface[0].(map[string]interface{})

	bootstrapConfig := config.NewBootstrapConfig()

	if cloudInitConfig, ok := bootStrapConfigDataMap["cloud_init_config"]; ok {
		cloudInitConfigData := cloudInitConfig.([]interface{})
		cloudInitConfigList := make([]vmmConfig.CloudInit, 0)
		for _, cloudInitData := range cloudInitConfigData {
			cloudInitDataMap := cloudInitData.(map[string]interface{})

			cloudInitObj := vmmConfig.NewCloudInit()
			if datasourceType, ok := cloudInitDataMap["datasource_type"]; ok {
				if datasourceType != nil && datasourceType != "" {
					const ConfigDriveV2 = 2
					subMap := map[string]interface{}{
						"CONFIG_DRIVE_V2": ConfigDriveV2,
					}
					pVal := subMap[datasourceType.(string)]
					if pVal == nil {
						cloudInitObj.DatasourceType = nil
					} else {
						p := vmmConfig.CloudInitDataSourceType(pVal.(int))
						cloudInitObj.DatasourceType = &p
					}
				}
			}
			if metadata, ok := cloudInitDataMap["metadata"]; ok {
				cloudInitObj.Metadata = utils.StringPtr(metadata.(string))
			}
			if cloudInitScript, ok := cloudInitDataMap["cloud_init_script"]; ok {
				cloudInitScriptData := cloudInitScript.([]interface{})[0].(map[string]interface{})
				cloudInitScriptObj := vmmConfig.NewOneOfCloudInitCloudInitScript()

				if userdata := cloudInitScriptData["user_data"]; userdata != nil && len(userdata.([]interface{})) > 0 {
					user := vmmConfig.NewUserdata()
					userVal := userdata.([]interface{})[0].(map[string]interface{})

					if value, ok := userVal["value"]; ok {
						user.Value = utils.StringPtr(value.(string))
					}

					err := cloudInitScriptObj.SetValue(*user)
					if err != nil {
						log.Printf("[ERROR] cloudInitScript : Error setting value for userdata: %v", err)
						return nil
					}
					cloudInitObj.CloudInitScript = cloudInitScriptObj
				} else if customKeyValues, ok := cloudInitScriptData["custom_key_values"]; ok && len(customKeyValues.([]interface{})) > 0 {
					customKeyValuesObj := expandCustomKeyValuesPairs(customKeyValues)
					err := cloudInitScriptObj.SetValue(*customKeyValuesObj)
					if err != nil {
						log.Printf("[ERROR] cloudInitScript: Error setting value for custom key values: %v", err)
						return nil
					}
					cloudInitObj.CloudInitScript = cloudInitScriptObj
				}
			}
			cloudInitConfigList = append(cloudInitConfigList, *cloudInitObj)
		}
		bootstrapConfig.CloudInitConfig = cloudInitConfigList
	}

	if environmentInfo, ok := bootStrapConfigDataMap["environment_info"]; ok {
		environmentInfoData := environmentInfo.([]interface{})[0].(map[string]interface{})
		environmentInfoObj := config.NewEnvironmentInfo()

		if providerType, ok := environmentInfoData["provider_type"]; ok {
			if providerType != nil && providerType != "" {
				const NTNX, AZURE, AWS, GCP, VSPHERE = 2, 3, 4, 5, 6
				subMap := map[string]interface{}{
					"NTNX":    NTNX,
					"AZURE":   AZURE,
					"AWS":     AWS,
					"GCP":     GCP,
					"VSPHERE": VSPHERE,
				}
				pVal := subMap[providerType.(string)]
				if pVal == nil {
					environmentInfoObj.ProviderType = nil
				} else {
					p := config.ProviderType(pVal.(int))
					environmentInfoObj.ProviderType = &p
				}
			}
		}
		if provisioningType, ok := environmentInfoData["provisioning_type"]; ok {
			if provisioningType != nil && provisioningType != "" {
				const NTNX, NATIVE = 2, 3
				subMap := map[string]interface{}{
					"NTNX":   NTNX,
					"NATIVE": NATIVE,
				}
				pVal := subMap[provisioningType.(string)]
				if pVal == nil {
					environmentInfoObj.ProvisioningType = nil
				} else {
					p := config.ProvisioningType(pVal.(int))
					environmentInfoObj.ProvisioningType = &p
				}
			}
		}
		if environmentType, ok := environmentInfoData["type"]; ok {
			if environmentType != nil && environmentType != "" {
				const ONPREM, NtnxCloud = 2, 3
				subMap := map[string]interface{}{
					"ONPREM":     ONPREM,
					"NTNX_CLOUD": NtnxCloud,
				}
				pVal := subMap[environmentType.(string)]
				if pVal == nil {
					environmentInfoObj.Type = nil
				} else {
					p := config.EnvironmentType(pVal.(int))
					environmentInfoObj.Type = &p
				}
			}
		}
		bootstrapConfig.EnvironmentInfo = environmentInfoObj
	}

	return bootstrapConfig
}

func expandCustomKeyValuesPairs(customKeyValues interface{}) *vmmConfig.CustomKeyValues {
	if customKeyValues != nil {
		customKeyValuesObj := vmmConfig.NewCustomKeyValues()
		customKeyValuesListInterface := customKeyValues.([]interface{})
		customKeyValuesListValue := customKeyValuesListInterface[0].(map[string]interface{})
		customKeyValuesList := customKeyValuesListValue["key_value_pairs"].([]interface{})
		if len(customKeyValuesList) > 0 {
			kvpList := make([]commonConfig.KVPair, 0)
			for _, customKeyValuesData := range customKeyValuesList {
				if keyValue := customKeyValuesData.(map[string]interface{}); keyValue != nil {
					kvpList = append(kvpList, expandKVPair(keyValue))
				}
			}
			customKeyValuesObj.KeyValuePairs = kvpList
		}
		return customKeyValuesObj
	}
	return nil
}

func expandKVPair(attribute map[string]interface{}) commonConfig.KVPair {
	var kv commonConfig.KVPair
	if attribute["name"] != nil && attribute["value"] != nil &&
		attribute["name"] != "" && attribute["value"] != "" {
		kv.Name = utils.StringPtr(attribute["name"].(string))
		kv.Value = expandValue(attribute["value"])
	}
	aJSON, _ := json.MarshalIndent(kv, "", "  ")
	log.Printf("[DEBUG] KVPair: %v", string(aJSON))
	return kv
}

func expandValue(kvPairValue interface{}) *commonConfig.OneOfKVPairValue {
	valueObj := commonConfig.NewOneOfKVPairValue()
	if kvPairValue != nil {
		valueData := kvPairValue.([]interface{})[0].(map[string]interface{})
		log.Printf("[DEBUG] kvPair valueData: %v", valueData)

		//nolint:gocritic // Keeping if-else for clarity in this specific case
		if valueData["string_list"] != nil && len(valueData["string_list"].([]interface{})) > 0 {
			log.Printf("[DEBUG] valueData of type string_list")
			stringList := valueData["string_list"].([]interface{})
			stringsListStr := make([]string, len(stringList))
			for i, v := range stringList {
				stringsListStr[i] = v.(string)
			}
			log.Printf("[DEBUG] stringsListStr: %v", stringsListStr)
			err := valueObj.SetValue(stringsListStr)
			if err != nil {
				log.Printf("[ERROR] Error setting value for string_list: %s", err)
				diag.Errorf("Error setting value for string_list: %s", err)
				return nil
			}
		} else if valueData["integer_list"] != nil && len(valueData["integer_list"].([]interface{})) > 0 {
			log.Printf("[DEBUG] valueData of type integer_list")
			integerList := valueData["integer_list"].([]interface{})
			integersListInt := make([]int, len(integerList))
			for i, v := range integerList {
				integersListInt[i] = v.(int)
			}
			err := valueObj.SetValue(integersListInt)
			if err != nil {
				log.Printf("[ERROR] Error setting value for integer_list: %s", err)
				diag.Errorf("Error setting value for integer_list: %s", err)
				return nil
			}
		} else if valueData["map_of_strings"] != nil && len(valueData["map_of_strings"].([]interface{})) > 0 {
			log.Printf("[DEBUG] valueData of type map_of_strings")
			mapOfStrings := make([]commonConfig.MapOfStringWrapper, 0)

			for _, mapOfStringsData := range valueData["map_of_strings"].([]interface{}) {
				mapOfStringsDataMap := mapOfStringsData.(map[string]interface{})
				mapOfStringsObj := commonConfig.NewMapOfStringWrapper()
				mapOfStringsObj.Map = make(map[string]string)
				for k, v := range mapOfStringsDataMap["map"].(map[string]interface{}) {
					mapOfStringsObj.Map[k] = v.(string)
				}
				mapOfStrings = append(mapOfStrings, *mapOfStringsObj)
			}
			aJSON, _ := json.Marshal(mapOfStrings)
			log.Printf("[DEBUG] mapOfStrings: %v", string(aJSON))
			log.Printf("[DEBUG] mapOfStrings type: %T", mapOfStrings)
			err := valueObj.SetValue(mapOfStrings)
			if err != nil {
				log.Printf("[ERROR] Error setting value for map: %s", err)
				diag.Errorf("Error setting value for map: %s", err)
				return nil
			}
		} else if valueData["string"] != nil && valueData["string"] != "" {
			log.Printf("[DEBUG] valueData of type string")
			err := valueObj.SetValue(valueData["string"].(string))
			if err != nil {
				log.Printf("[ERROR] Error setting value for string: %s", err)
				diag.Errorf("Error setting value for string: %s", err)
				return nil
			}
		} else if valueData["object"] != nil && len(valueData["object"].(map[string]interface{})) > 0 {
			log.Printf("[DEBUG] valueData of type object")
			object := make(map[string]string)
			for k, v := range valueData["object"].(map[string]interface{}) {
				object[k] = v.(string)
			}
			err := valueObj.SetValue(object)
			if err != nil {
				log.Printf("[ERROR] Error setting value for object: %s", err)
				diag.Errorf("Error setting value for object: %s", err)
				return nil
			}
		} else if valueData["integer"] != nil && valueData["integer"] != 0 {
			log.Printf("[DEBUG] valueData of type integer")
			err := valueObj.SetValue(valueData["integer"].(int))
			if err != nil {
				log.Printf("[ERROR] Error setting value for integer: %s", err)
				diag.Errorf("Error setting value for integer: %s", err)
				return nil
			}
		} else if valueData["boolean"] != nil {
			log.Printf("[DEBUG] valueData of type boolean")
			err := valueObj.SetValue(valueData["boolean"].(bool))
			if err != nil {
				log.Printf("[ERROR] Error setting value for boolean: %s", err)
				diag.Errorf("Error setting value for boolean: %s", err)
				return nil
			}
		} else {
			log.Printf("[ERROR] invalid value type")
			return nil
		}
	}
	return valueObj
}

func expandCredentials(credentials []interface{}) []commonConfig.BasicAuth {
	if len(credentials) == 0 {
		return nil
	}

	credentialsList := make([]commonConfig.BasicAuth, 0)
	for _, credential := range credentials {
		credentialData := credential.(map[string]interface{})
		credentialObj := commonConfig.NewBasicAuth()
		if username, ok := credentialData["username"]; ok {
			credentialObj.Username = utils.StringPtr(username.(string))
		}
		if password, ok := credentialData["password"]; ok {
			credentialObj.Password = utils.StringPtr(password.(string))
		}
		credentialsList = append(credentialsList, *credentialObj)
	}
	return credentialsList
}

func expandResourceConfig(resourceConfig interface{}) *config.DomainManagerResourceConfig {
	if len(resourceConfig.([]interface{})) == 0 {
		return nil
	}
	resourceConfigI := resourceConfig.([]interface{})
	resourceConfigData := resourceConfigI[0].(map[string]interface{})

	resourceConfigObj := config.NewDomainManagerResourceConfig()

	if containerExtIds, ok := resourceConfigData["container_ext_ids"]; ok {
		containerExtIdsList := containerExtIds.([]interface{})
		containerExtIdsListObj := make([]string, 0)
		for _, containerExtID := range containerExtIdsList {
			containerExtIdsListObj = append(containerExtIdsListObj, containerExtID.(string))
		}
		resourceConfigObj.ContainerExtIds = containerExtIdsListObj
	}
	if dataDiskSizeBytes, ok := resourceConfigData["data_disk_size_bytes"]; ok {
		resourceConfigObj.DataDiskSizeBytes = utils.Int64Ptr(int64(dataDiskSizeBytes.(int)))
	}
	if memorySizeBytes, ok := resourceConfigData["memory_size_bytes"]; ok {
		resourceConfigObj.MemorySizeBytes = utils.Int64Ptr(int64(memorySizeBytes.(int)))
	}
	if numVcpus, ok := resourceConfigData["num_vcpus"]; ok {
		resourceConfigObj.NumVcpus = utils.IntPtr(numVcpus.(int))
	}

	return resourceConfigObj
}

// network expanders
func expandPCNetwork(pcNetwork interface{}) *config.DomainManagerNetwork {
	if len(pcNetwork.([]interface{})) == 0 {
		return nil
	}

	pcNetworkInterface := pcNetwork.([]interface{})
	pcNetworkData := pcNetworkInterface[0].(map[string]interface{})

	pcNetworkObj := config.NewDomainManagerNetwork()

	if externalAddress, ok := pcNetworkData["external_address"]; ok && len(externalAddress.([]interface{})) > 0 {
		externalAddressData := externalAddress.([]interface{})[0].(map[string]interface{})
		pcNetworkObj.ExternalAddress = expandIPAddress(externalAddressData)
	}
	if nameServers, ok := pcNetworkData["name_servers"]; ok && len(nameServers.([]interface{})) > 0 {
		nameServersData := nameServers.([]interface{})
		nameServersObj := make([]commonConfig.IPAddressOrFQDN, 0)
		for _, nameServerData := range nameServersData {
			nameServerObj := expandIPAddressOrFqdn(nameServerData.(map[string]interface{}))
			nameServersObj = append(nameServersObj, nameServerObj)
		}
		pcNetworkObj.NameServers = nameServersObj
	}
	if ntpServers, ok := pcNetworkData["ntp_servers"]; ok && len(ntpServers.([]interface{})) > 0 {
		ntpServersData := ntpServers.([]interface{})
		ntpServersObj := make([]commonConfig.IPAddressOrFQDN, 0)
		for _, ntpServerData := range ntpServersData {
			ntpServerObj := expandIPAddressOrFqdn(ntpServerData.(map[string]interface{}))
			ntpServersObj = append(ntpServersObj, ntpServerObj)
		}
		pcNetworkObj.NtpServers = ntpServersObj
	}
	if internalNetworks, ok := pcNetworkData["internal_networks"]; ok && len(internalNetworks.([]interface{})) > 0 {
		internalNetworksData := internalNetworks.([]interface{})
		internalNetworksList := make([]config.BaseNetwork, 0)
		for _, internalNetworkData := range internalNetworksData {
			internalNetworkObj := expandInternalNetwork(internalNetworkData.(map[string]interface{}))
			internalNetworksList = append(internalNetworksList, internalNetworkObj)
		}
		pcNetworkObj.InternalNetworks = internalNetworksList
	}
	if externalNetworks, ok := pcNetworkData["external_networks"]; ok && len(externalNetworks.([]interface{})) > 0 {
		externalNetworksData := externalNetworks.([]interface{})
		externalNetworksList := make([]config.ExternalNetwork, 0)
		for _, externalNetworkData := range externalNetworksData {
			externalNetworkObj := expandExternalNetwork(externalNetworkData.(map[string]interface{}))
			externalNetworksList = append(externalNetworksList, externalNetworkObj)
		}
		pcNetworkObj.ExternalNetworks = externalNetworksList
	}

	return pcNetworkObj
}

func expandExternalNetwork(externalNetwork map[string]interface{}) config.ExternalNetwork {
	externalNetworkObj := config.NewExternalNetwork()

	if defaultGateway, ok := externalNetwork["default_gateway"]; ok {
		defaultGatewayData := defaultGateway.([]interface{})[0].(map[string]interface{})
		defaultGatewayObj := expandIPAddressOrFqdn(defaultGatewayData)
		externalNetworkObj.DefaultGateway = &defaultGatewayObj
	}
	if subnetMask, ok := externalNetwork["subnet_mask"]; ok {
		subnetMaskData := subnetMask.([]interface{})[0].(map[string]interface{})
		subnetMaskObj := expandIPAddressOrFqdn(subnetMaskData)
		externalNetworkObj.SubnetMask = &subnetMaskObj
	}
	if ipRanges, ok := externalNetwork["ip_ranges"]; ok {
		ipRangesData := ipRanges.([]interface{})
		ipRangesList := make([]commonConfig.IpRange, 0)
		for _, ipRangeData := range ipRangesData {
			ipRangeObj := expandIPRange(ipRangeData.(map[string]interface{}))
			ipRangesList = append(ipRangesList, ipRangeObj)
		}
		externalNetworkObj.IpRanges = ipRangesList
	}
	if networkExtID, ok := externalNetwork["network_ext_id"]; ok {
		externalNetworkObj.NetworkExtId = utils.StringPtr(networkExtID.(string))
	}

	return *externalNetworkObj
}

func expandInternalNetwork(internalNetwork map[string]interface{}) config.BaseNetwork {
	internalNetworkObj := config.NewBaseNetwork()

	if defaultGateway, ok := internalNetwork["default_gateway"]; ok {
		defaultGatewayData := defaultGateway.([]interface{})[0].(map[string]interface{})
		defaultGatewayObj := expandIPAddressOrFqdn(defaultGatewayData)
		internalNetworkObj.DefaultGateway = &defaultGatewayObj
	}
	if subnetMask, ok := internalNetwork["subnet_mask"]; ok {
		subnetMaskData := subnetMask.([]interface{})[0].(map[string]interface{})
		subnetMaskObj := expandIPAddressOrFqdn(subnetMaskData)
		internalNetworkObj.SubnetMask = &subnetMaskObj
	}
	if ipRanges, ok := internalNetwork["ip_ranges"]; ok {
		ipRangesData := ipRanges.([]interface{})
		ipRangesObj := make([]commonConfig.IpRange, 0)
		for _, ipRangeData := range ipRangesData {
			ipRangeObj := expandIPRange(ipRangeData.(map[string]interface{}))
			ipRangesObj = append(ipRangesObj, ipRangeObj)
		}
		internalNetworkObj.IpRanges = ipRangesObj
	}

	return *internalNetworkObj
}

func expandIPRange(ipRange map[string]interface{}) commonConfig.IpRange {
	ipRangeObj := commonConfig.NewIpRange()

	if begin, ok := ipRange["begin"]; ok {
		beginData := begin.([]interface{})[0].(map[string]interface{})
		beginObj := expandIPAddress(beginData)
		ipRangeObj.Begin = beginObj
	}
	if end, ok := ipRange["end"]; ok {
		endData := end.([]interface{})[0].(map[string]interface{})
		endObj := expandIPAddress(endData)
		ipRangeObj.End = endObj
	}

	return *ipRangeObj
}

// ip address expanders
func expandIPAddress(ipAddress map[string]interface{}) *commonConfig.IPAddress {
	ipAddressObj := commonConfig.NewIPAddress()

	if ipv4, ok := ipAddress["ipv4"]; ok && len(ipv4.([]interface{})) > 0 {
		ipAddressObj.Ipv4 = expandIPv4Address(ipv4)
	}
	if ipv6, ok := ipAddress["ipv6"]; ok && len(ipv6.([]interface{})) > 0 {
		ipAddressObj.Ipv6 = expandIPv6Address(ipv6)
	}

	return ipAddressObj
}

func expandIPAddressOrFqdn(ipAddressOrFQDN map[string]interface{}) commonConfig.IPAddressOrFQDN {
	ipAddressOrFQDNObj := *commonConfig.NewIPAddressOrFQDN()

	if ipv4, ok := ipAddressOrFQDN["ipv4"]; ok && len(ipv4.([]interface{})) > 0 {
		ipAddressOrFQDNObj.Ipv4 = expandIPv4Address(ipv4)
	}
	if ipv6, ok := ipAddressOrFQDN["ipv6"]; ok && len(ipv6.([]interface{})) > 0 {
		ipAddressOrFQDNObj.Ipv6 = expandIPv6Address(ipv6)
	}
	if fqdn, ok := ipAddressOrFQDN["fqdn"]; ok && len(fqdn.([]interface{})) > 0 {
		ipAddressOrFQDNObj.Fqdn = expandFQDN(fqdn)
	}

	return ipAddressOrFQDNObj
}

func expandIPv4Address(ipv4 interface{}) *commonConfig.IPv4Address {
	if ipv4 != nil {
		ipAddress := ipv4.([]interface{})[0].(map[string]interface{})
		ip := &commonConfig.IPv4Address{}
		if value, ok := ipAddress["value"].(string); ok {
			ip.Value = utils.StringPtr(value)
		}
		if prefixLength, ok := ipAddress["prefix_length"].(int); ok {
			ip.PrefixLength = utils.IntPtr(prefixLength)
		}
		return ip
	}
	return nil
}

func expandIPv6Address(ipv6 interface{}) *commonConfig.IPv6Address {
	if ipv6 != nil {
		ipAddress := ipv6.([]interface{})[0].(map[string]interface{})
		ip := &commonConfig.IPv6Address{}
		if value, ok := ipAddress["value"].(string); ok {
			ip.Value = utils.StringPtr(value)
		}
		if prefixLength, ok := ipAddress["prefix_length"].(int); ok {
			ip.PrefixLength = utils.IntPtr(prefixLength)
		}
		return ip
	}
	return nil
}

func expandFQDN(fqdn interface{}) *commonConfig.FQDN {
	if fqdn != nil {
		fqdnMap := fqdn.([]interface{})[0].(map[string]interface{})
		f := &commonConfig.FQDN{}
		if value, ok := fqdnMap["value"].(string); ok {
			f.Value = utils.StringPtr(value)
		}
		return f
	}
	return nil
}
