package vmmv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	vmmPrismConfig "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/prism/v4/config"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixVmGuestCustomizationProfileV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixVmGuestCustomizationProfileV2Create,
		ReadContext:   ResourceNutanixVmGuestCustomizationProfileV2Read,
		UpdateContext: ResourceNutanixVmGuestCustomizationProfileV2Update,
		DeleteContext: ResourceNutanixVmGuestCustomizationProfileV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sysprep_config": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"customization": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"sysprep_params": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: resourceVmGcProfileSysprepParamsSchema(),
													},
												},
												"answer_file": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unattend_xml": {
																Type:     schema.TypeString,
																Required: true,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"updated_by": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rel": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceVmGcProfileSysprepParamsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"first_logon_commands": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"general_settings": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"administrator_password": {
						Type:      schema.TypeString,
						Optional:  true,
						Sensitive: true,
					},
					"auto_logon_settings": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"logon_count": {
									Type:     schema.TypeInt,
									Required: true,
								},
							},
						},
					},
					"computer_name": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"use_vm_name": {
									Type:     schema.TypeBool,
									Optional: true,
								},
								"must_provide_during_deployment": {
									Type:     schema.TypeBool,
									Optional: true,
								},
							},
						},
					},
					"registered_organization": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"registered_owner": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"timezone": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"windows_product_key": {
						Type:      schema.TypeString,
						Optional:  true,
						Sensitive: true,
					},
				},
			},
		},
		"locale_settings": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"system_locale": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"ui_language": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"user_locale": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"network_settings": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"nic_config_list": {
						Type:     schema.TypeList,
						Required: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"dns_config": {
									Type:     schema.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"alternate_dns_server_addresses": {
												Type:     schema.TypeList,
												Optional: true,
												Elem: &schema.Schema{
													Type: schema.TypeString,
												},
											},
											"preferred_dns_server_address": {
												Type:     schema.TypeString,
												Required: true,
											},
										},
									},
								},
								"ipv4_config": {
									Type:     schema.TypeList,
									Required: true,
									MaxItems: 1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"use_dhcp": {
												Type:     schema.TypeBool,
												Optional: true,
											},
											"must_provide_during_deployment": {
												Type:     schema.TypeBool,
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"workgroup_or_domain_info": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"workgroup": {
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:     schema.TypeString,
									Required: true,
								},
							},
						},
					},
					"domain_settings": {
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"credentials": {
									Type:     schema.TypeList,
									Required: true,
									MaxItems: 1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"domain_name": {
												Type:     schema.TypeString,
												Required: true,
											},
											"password": {
												Type:      schema.TypeString,
												Required:  true,
												Sensitive: true,
											},
											"username": {
												Type:     schema.TypeString,
												Required: true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func ResourceNutanixVmGuestCustomizationProfileV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	body := config.VmGuestCustomizationProfile{}

	if name, ok := d.GetOk("name"); ok {
		body.Name = utils.StringPtr(name.(string))
	}
	if desc, ok := d.GetOk("description"); ok {
		body.Description = utils.StringPtr(desc.(string))
	}
	if cfg, ok := d.GetOk("config"); ok {
		configList := cfg.([]interface{})
		if len(configList) > 0 && configList[0] != nil {
			expandedConfig := expandVmGcProfileConfig(configList[0].(map[string]interface{}))
			if expandedConfig != nil {
				oneOfConfig := config.NewOneOfVmGuestCustomizationProfileConfig()
				oneOfConfig.SetValue(expandedConfig)
				body.Config = oneOfConfig
			}
		}
	}

	aJSON, _ := json.MarshalIndent(body, "", " ")
	log.Printf("[DEBUG] VM Guest Customization Profile Create Request Body: %s", string(aJSON))

	resp, err := conn.VmGuestCustomizationProfilesAPIInstance.CreateVmGuestCustomizationProfile(&body)
	if err != nil {
		return diag.Errorf("error while creating VM Guest Customization Profile: %v", err)
	}

	TaskRef := resp.Data.GetValue().(vmmPrismConfig.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for VM Guest Customization Profile (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching VM Guest Customization Profile task: %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)

	uuid, err := common.ExtractEntityUUIDFromTask(taskDetails, utils.RelEntityTypeVmGuestCustomizationProfile, "VM Guest Customization Profile")
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(utils.StringValue(uuid))
	return ResourceNutanixVmGuestCustomizationProfileV2Read(ctx, d, meta)
}

func ResourceNutanixVmGuestCustomizationProfileV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	resp, err := conn.VmGuestCustomizationProfilesAPIInstance.GetVmGuestCustomizationProfileById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching VM Guest Customization Profile: %v", err)
	}

	getResp := resp.Data.GetValue().(config.VmGuestCustomizationProfile)

	flattenedProfile := flattenVmGuestCustomizationProfileEntity(getResp)

	for k, v := range flattenedProfile {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func ResourceNutanixVmGuestCustomizationProfileV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	resp, err := conn.VmGuestCustomizationProfilesAPIInstance.GetVmGuestCustomizationProfileById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching VM Guest Customization Profile: %v", err)
	}

	updateSpec := resp.Data.GetValue().(config.VmGuestCustomizationProfile)
	clearProfileDiscriminators(&updateSpec)

	if d.HasChange("name") {
		updateSpec.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("description") {
		updateSpec.Description = utils.StringPtr(d.Get("description").(string))
	}
	if d.HasChange("config") {
		if cfg, ok := d.GetOk("config"); ok {
			configList := cfg.([]interface{})
			if len(configList) > 0 && configList[0] != nil {
				expandedConfig := expandVmGcProfileConfig(configList[0].(map[string]interface{}))
				if expandedConfig != nil {
					oneOfConfig := config.NewOneOfVmGuestCustomizationProfileConfig()
					oneOfConfig.SetValue(expandedConfig)
					updateSpec.Config = oneOfConfig
				}
			}
		}
	}

	aJSON, _ := json.MarshalIndent(updateSpec, "", " ")
	log.Printf("[DEBUG] VM Guest Customization Profile Update Request Body: %s", string(aJSON))

	etagValue := conn.VmGuestCustomizationProfilesAPIInstance.ApiClient.GetEtag(resp)
	args := make(map[string]interface{})
	args["If-Match"] = utils.StringPtr(etagValue)

	updateResp, err := conn.VmGuestCustomizationProfilesAPIInstance.UpdateVmGuestCustomizationProfileById(utils.StringPtr(d.Id()), &updateSpec, args)
	if err != nil {
		return diag.Errorf("error while updating VM Guest Customization Profile: %v", err)
	}

	TaskRef := updateResp.Data.GetValue().(vmmPrismConfig.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for VM Guest Customization Profile (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	return ResourceNutanixVmGuestCustomizationProfileV2Read(ctx, d, meta)
}

func ResourceNutanixVmGuestCustomizationProfileV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	readResp, err := conn.VmGuestCustomizationProfilesAPIInstance.GetVmGuestCustomizationProfileById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching VM Guest Customization Profile: %v", err)
	}

	etagValue := conn.VmGuestCustomizationProfilesAPIInstance.ApiClient.GetEtag(readResp)
	args := make(map[string]interface{})
	args["If-Match"] = utils.StringPtr(etagValue)

	resp, err := conn.VmGuestCustomizationProfilesAPIInstance.DeleteVmGuestCustomizationProfileById(utils.StringPtr(d.Id()), args)
	if err != nil {
		return diag.Errorf("error while deleting VM Guest Customization Profile: %v", err)
	}

	TaskRef := resp.Data.GetValue().(vmmPrismConfig.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for VM Guest Customization Profile (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	return nil
}

func expandVmGcProfileConfig(cfgMap map[string]interface{}) interface{} {
	if sysprepList, ok := cfgMap["sysprep_config"].([]interface{}); ok && len(sysprepList) > 0 && sysprepList[0] != nil {
		sysprepMap := sysprepList[0].(map[string]interface{})
		sysprepConfig := config.NewVmGcProfileSysprepConfig()

		if custList, ok := sysprepMap["customization"].([]interface{}); ok && len(custList) > 0 && custList[0] != nil {
			custMap := custList[0].(map[string]interface{})
			oneOfCust := config.NewOneOfVmGcProfileSysprepConfigCustomization()

			if spList, ok := custMap["sysprep_params"].([]interface{}); ok && len(spList) > 0 && spList[0] != nil {
				params := expandSysprepParams(spList[0].(map[string]interface{}))
				oneOfCust.SetValue(*params)
			} else if afList, ok := custMap["answer_file"].([]interface{}); ok && len(afList) > 0 && afList[0] != nil {
				afMap := afList[0].(map[string]interface{})
				af := config.NewVmGcProfileAnswerFile()
				if v, ok := afMap["unattend_xml"].(string); ok && v != "" {
					af.UnattendXml = utils.StringPtr(v)
				}
				oneOfCust.SetValue(*af)
			}
			sysprepConfig.Customization = oneOfCust
		}
		return *sysprepConfig
	}

	return nil
}

func expandSysprepParams(spMap map[string]interface{}) *config.VmGcProfileSysprepParams {
	params := config.NewVmGcProfileSysprepParams()

	if cmds, ok := spMap["first_logon_commands"].([]interface{}); ok && len(cmds) > 0 {
		commands := make([]string, len(cmds))
		for i, c := range cmds {
			commands[i] = c.(string)
		}
		params.FirstLogonCommands = commands
	}

	if gsList, ok := spMap["general_settings"].([]interface{}); ok && len(gsList) > 0 && gsList[0] != nil {
		gsMap := gsList[0].(map[string]interface{})
		gs := config.NewVmGcProfileGeneralSettings()

		if v, ok := gsMap["administrator_password"].(string); ok && v != "" {
			gs.AdministratorPassword = utils.StringPtr(v)
		}
		if alsList, ok := gsMap["auto_logon_settings"].([]interface{}); ok && len(alsList) > 0 && alsList[0] != nil {
			alsMap := alsList[0].(map[string]interface{})
			als := config.NewVmGcProfileAutoLogonSettings()
			if v, ok := alsMap["logon_count"].(int); ok {
				als.LogonCount = utils.IntPtr(v)
			}
			gs.AutoLogonSettings = als
		}
		if cnList, ok := gsMap["computer_name"].([]interface{}); ok && len(cnList) > 0 && cnList[0] != nil {
			cnMap := cnList[0].(map[string]interface{})
			oneOfCN := config.NewOneOfVmGcProfileGeneralSettingsComputerName()
			if useVM, ok := cnMap["use_vm_name"].(bool); ok && useVM {
				oneOfCN.SetValue(*config.NewVmGcProfileUseVmName())
			} else if mustProvide, ok := cnMap["must_provide_during_deployment"].(bool); ok && mustProvide {
				oneOfCN.SetValue(*config.NewVmGcProfileMustProvideDuringDeployment())
			}
			gs.ComputerName = oneOfCN
		}
		if v, ok := gsMap["registered_organization"].(string); ok && v != "" {
			gs.RegisteredOrganization = utils.StringPtr(v)
		}
		if v, ok := gsMap["registered_owner"].(string); ok && v != "" {
			gs.RegisteredOwner = utils.StringPtr(v)
		}
		if v, ok := gsMap["timezone"].(string); ok && v != "" {
			gs.Timezone = utils.StringPtr(v)
		}
		if v, ok := gsMap["windows_product_key"].(string); ok && v != "" {
			gs.WindowsProductKey = utils.StringPtr(v)
		}
		params.GeneralSettings = gs
	}

	if lsList, ok := spMap["locale_settings"].([]interface{}); ok && len(lsList) > 0 && lsList[0] != nil {
		lsMap := lsList[0].(map[string]interface{})
		ls := config.NewVmGcProfileLocaleSettings()
		if v, ok := lsMap["system_locale"].(string); ok && v != "" {
			ls.SystemLocale = utils.StringPtr(v)
		}
		if v, ok := lsMap["ui_language"].(string); ok && v != "" {
			ls.UiLanguage = utils.StringPtr(v)
		}
		if v, ok := lsMap["user_locale"].(string); ok && v != "" {
			ls.UserLocale = utils.StringPtr(v)
		}
		params.LocaleSettings = ls
	}

	if nsList, ok := spMap["network_settings"].([]interface{}); ok && len(nsList) > 0 && nsList[0] != nil {
		nsMap := nsList[0].(map[string]interface{})
		ns := config.NewVmGcProfileNetworkSettings()
		if nicList, ok := nsMap["nic_config_list"].([]interface{}); ok && len(nicList) > 0 {
			nics := make([]config.VmGcProfileNicConfig, len(nicList))
			for i, nicRaw := range nicList {
				nicMap := nicRaw.(map[string]interface{})
				nic := *config.NewVmGcProfileNicConfig()
				if dnsList, ok := nicMap["dns_config"].([]interface{}); ok && len(dnsList) > 0 && dnsList[0] != nil {
					dnsMap := dnsList[0].(map[string]interface{})
					dns := config.NewVmGcProfileDnsConfig()
					if v, ok := dnsMap["preferred_dns_server_address"].(string); ok && v != "" {
						dns.PreferredDnsServerAddress = utils.StringPtr(v)
					}
					if altDns, ok := dnsMap["alternate_dns_server_addresses"].([]interface{}); ok && len(altDns) > 0 {
						addrs := make([]string, len(altDns))
						for j, a := range altDns {
							addrs[j] = a.(string)
						}
						dns.AlternateDnsServerAddresses = addrs
					}
					nic.DnsConfig = dns
				}
				if ipv4List, ok := nicMap["ipv4_config"].([]interface{}); ok && len(ipv4List) > 0 && ipv4List[0] != nil {
					ipv4Map := ipv4List[0].(map[string]interface{})
					oneOfIPv4 := config.NewOneOfVmGcProfileNicConfigIpv4Config()
					if useDhcp, ok := ipv4Map["use_dhcp"].(bool); ok && useDhcp {
						oneOfIPv4.SetValue(*config.NewVmGcProfileUseDhcp())
					} else if mustProvide, ok := ipv4Map["must_provide_during_deployment"].(bool); ok && mustProvide {
						oneOfIPv4.SetValue(*config.NewVmGcProfileMustProvideDuringDeployment())
					}
					nic.Ipv4Config = oneOfIPv4
				}
				nics[i] = nic
			}
			ns.NicConfigList = nics
		}
		params.NetworkSettings = ns
	}

	if wdList, ok := spMap["workgroup_or_domain_info"].([]interface{}); ok && len(wdList) > 0 && wdList[0] != nil {
		wdMap := wdList[0].(map[string]interface{})
		oneOfWD := config.NewOneOfVmGcProfileSysprepParamsWorkgroupOrDomainInfo()
		if wgList, ok := wdMap["workgroup"].([]interface{}); ok && len(wgList) > 0 && wgList[0] != nil {
			wgMap := wgList[0].(map[string]interface{})
			wg := config.NewVmGcProfileWorkgroup()
			if v, ok := wgMap["name"].(string); ok && v != "" {
				wg.Name = utils.StringPtr(v)
			}
			oneOfWD.SetValue(*wg)
		} else if dsList, ok := wdMap["domain_settings"].([]interface{}); ok && len(dsList) > 0 && dsList[0] != nil {
			dsMap := dsList[0].(map[string]interface{})
			ds := config.NewVmGcProfileDomainSettings()
			if credsList, ok := dsMap["credentials"].([]interface{}); ok && len(credsList) > 0 && credsList[0] != nil {
				credsMap := credsList[0].(map[string]interface{})
				creds := config.NewVmGcProfileDomainCredentials()
				if v, ok := credsMap["domain_name"].(string); ok && v != "" {
					creds.DomainName = utils.StringPtr(v)
				}
				if v, ok := credsMap["password"].(string); ok && v != "" {
					creds.Password = utils.StringPtr(v)
				}
				if v, ok := credsMap["username"].(string); ok && v != "" {
					creds.Username = utils.StringPtr(v)
				}
				ds.Credentials = creds
			}
			oneOfWD.SetValue(*ds)
		}
		params.WorkgroupOrDomainInfo = oneOfWD
	}

	return params
}

// clearProfileDiscriminators removes SDK-populated item discriminator fields
// that are rejected by the API server during create/update requests.
func clearProfileDiscriminators(profile *config.VmGuestCustomizationProfile) {
	profile.ConfigItemDiscriminator_ = nil
	if profile.Config == nil {
		return
	}
	val := profile.Config.GetValue()
	if sysprep, ok := val.(config.VmGcProfileSysprepConfig); ok {
		sysprep.CustomizationItemDiscriminator_ = nil
		if sysprep.Customization != nil {
			custVal := sysprep.Customization.GetValue()
			if sp, ok := custVal.(config.VmGcProfileSysprepParams); ok {
				sp.WorkgroupOrDomainInfoItemDiscriminator_ = nil
				if sp.GeneralSettings != nil {
					sp.GeneralSettings.ComputerNameItemDiscriminator_ = nil
				}
				if sp.NetworkSettings != nil {
					for i := range sp.NetworkSettings.NicConfigList {
						sp.NetworkSettings.NicConfigList[i].Ipv4ConfigItemDiscriminator_ = nil
					}
				}
				oneOfCust := config.NewOneOfVmGcProfileSysprepConfigCustomization()
				oneOfCust.SetValue(sp)
				sysprep.Customization = oneOfCust
			}
		}
		oneOfConfig := config.NewOneOfVmGuestCustomizationProfileConfig()
		oneOfConfig.SetValue(sysprep)
		profile.Config = oneOfConfig
	}
}
