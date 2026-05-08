package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixVmGuestCustomizationProfileV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixVmGuestCustomizationProfileV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sysprep_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"customization": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"sysprep_params": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: vmGcProfileSysprepParamsSchema(true),
													},
												},
												"answer_file": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unattend_xml": {
																Type:     schema.TypeString,
																Computed: true,
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

func DatasourceNutanixVmGuestCustomizationProfileV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	extID := d.Get("ext_id").(string)

	resp, err := conn.VmGuestCustomizationProfilesAPIInstance.GetVmGuestCustomizationProfileById(utils.StringPtr(extID))
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

	d.SetId(flattenedProfile["ext_id"].(string))

	return nil
}

func vmGcProfileSysprepParamsSchema(computed bool) map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"first_logon_commands": {
			Type:     schema.TypeList,
			Computed: computed,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"general_settings": {
			Type:     schema.TypeList,
			Computed: computed,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"administrator_password": {
						Type:      schema.TypeString,
						Computed:  computed,
						Sensitive: true,
					},
					"auto_logon_settings": {
						Type:     schema.TypeList,
						Computed: computed,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"logon_count": {
									Type:     schema.TypeInt,
									Computed: computed,
								},
							},
						},
					},
					"computer_name": {
						Type:     schema.TypeList,
						Computed: computed,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"use_vm_name": {
									Type:     schema.TypeBool,
									Computed: computed,
								},
								"must_provide_during_deployment": {
									Type:     schema.TypeBool,
									Computed: computed,
								},
							},
						},
					},
					"registered_organization": {
						Type:     schema.TypeString,
						Computed: computed,
					},
					"registered_owner": {
						Type:     schema.TypeString,
						Computed: computed,
					},
					"timezone": {
						Type:     schema.TypeString,
						Computed: computed,
					},
					"windows_product_key": {
						Type:      schema.TypeString,
						Computed:  computed,
						Sensitive: true,
					},
				},
			},
		},
		"locale_settings": {
			Type:     schema.TypeList,
			Computed: computed,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"system_locale": {
						Type:     schema.TypeString,
						Computed: computed,
					},
					"ui_language": {
						Type:     schema.TypeString,
						Computed: computed,
					},
					"user_locale": {
						Type:     schema.TypeString,
						Computed: computed,
					},
				},
			},
		},
		"network_settings": {
			Type:     schema.TypeList,
			Computed: computed,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"nic_config_list": {
						Type:     schema.TypeList,
						Computed: computed,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"dns_config": {
									Type:     schema.TypeList,
									Computed: computed,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"alternate_dns_server_addresses": {
												Type:     schema.TypeList,
												Computed: computed,
												Elem: &schema.Schema{
													Type: schema.TypeString,
												},
											},
											"preferred_dns_server_address": {
												Type:     schema.TypeString,
												Computed: computed,
											},
										},
									},
								},
								"ipv4_config": {
									Type:     schema.TypeList,
									Computed: computed,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"use_dhcp": {
												Type:     schema.TypeBool,
												Computed: computed,
											},
											"must_provide_during_deployment": {
												Type:     schema.TypeBool,
												Computed: computed,
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
			Computed: computed,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"workgroup": {
						Type:     schema.TypeList,
						Computed: computed,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:     schema.TypeString,
									Computed: computed,
								},
							},
						},
					},
					"domain_settings": {
						Type:     schema.TypeList,
						Computed: computed,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"credentials": {
									Type:     schema.TypeList,
									Computed: computed,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"domain_name": {
												Type:     schema.TypeString,
												Computed: computed,
											},
											"password": {
												Type:      schema.TypeString,
												Computed:  computed,
												Sensitive: true,
											},
											"username": {
												Type:     schema.TypeString,
												Computed: computed,
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
	if !computed {
		s["first_logon_commands"].Optional = true
		s["general_settings"].Optional = true
		for _, v := range s["general_settings"].Elem.(*schema.Resource).Schema {
			v.Optional = true
			v.Computed = false
			if v.Elem != nil {
				if r, ok := v.Elem.(*schema.Resource); ok {
					for _, vv := range r.Schema {
						vv.Optional = true
						vv.Computed = false
					}
				}
			}
		}
		s["locale_settings"].Optional = true
		for _, v := range s["locale_settings"].Elem.(*schema.Resource).Schema {
			v.Optional = true
			v.Computed = false
		}
		s["network_settings"].Optional = true
		s["workgroup_or_domain_info"].Optional = true
	}
	return s
}

func flattenVmGuestCustomizationProfileEntity(profile config.VmGuestCustomizationProfile) map[string]interface{} {
	result := make(map[string]interface{})
	if profile.ExtId != nil {
		result["ext_id"] = *profile.ExtId
	}
	if profile.Name != nil {
		result["name"] = *profile.Name
	}
	if profile.Description != nil {
		result["description"] = *profile.Description
	}
	if profile.CreateTime != nil {
		result["create_time"] = utils.TimeStringValue(profile.CreateTime)
	}
	if profile.UpdateTime != nil {
		result["update_time"] = utils.TimeStringValue(profile.UpdateTime)
	}
	if profile.CreatedBy != nil && profile.CreatedBy.ExtId != nil {
		result["created_by"] = map[string]string{"ext_id": *profile.CreatedBy.ExtId}
	}
	if profile.UpdatedBy != nil && profile.UpdatedBy.ExtId != nil {
		result["updated_by"] = map[string]string{"ext_id": *profile.UpdatedBy.ExtId}
	}
	if profile.TenantId != nil {
		result["tenant_id"] = *profile.TenantId
	}
	linksList := make([]map[string]interface{}, 0)
	if profile.Links != nil {
		linksList = make([]map[string]interface{}, len(profile.Links))
		for i, l := range profile.Links {
			lMap := map[string]interface{}{}
			if l.Href != nil {
				lMap["href"] = *l.Href
			}
			if l.Rel != nil {
				lMap["rel"] = *l.Rel
			}
			linksList[i] = lMap
		}
	}
	result["links"] = linksList
	if profile.Config != nil {
		result["config"] = flattenVmGcProfileConfig(profile.Config)
	}
	return result
}

func flattenVmGcProfileConfig(oneOf *config.OneOfVmGuestCustomizationProfileConfig) []map[string]interface{} {
	if oneOf == nil {
		return nil
	}
	configMap := make(map[string]interface{})
	val := oneOf.GetValue()
	switch v := val.(type) {
	case config.VmGcProfileSysprepConfig:
		configMap["sysprep_config"] = flattenSysprepConfig(&v)
	}
	return []map[string]interface{}{configMap}
}

func flattenSysprepConfig(sc *config.VmGcProfileSysprepConfig) []map[string]interface{} {
	if sc == nil {
		return nil
	}
	result := map[string]interface{}{}
	if sc.Customization != nil {
		result["customization"] = flattenSysprepCustomization(sc.Customization)
	}
	return []map[string]interface{}{result}
}

func flattenSysprepCustomization(oneOf *config.OneOfVmGcProfileSysprepConfigCustomization) []map[string]interface{} {
	if oneOf == nil {
		return nil
	}
	custMap := make(map[string]interface{})
	val := oneOf.GetValue()
	switch v := val.(type) {
	case config.VmGcProfileSysprepParams:
		custMap["sysprep_params"] = flattenSysprepParams(&v)
	case config.VmGcProfileAnswerFile:
		custMap["answer_file"] = flattenAnswerFile(&v)
	}
	return []map[string]interface{}{custMap}
}

func flattenSysprepParams(sp *config.VmGcProfileSysprepParams) []map[string]interface{} {
	if sp == nil {
		return nil
	}
	result := map[string]interface{}{}
	if sp.FirstLogonCommands != nil {
		result["first_logon_commands"] = sp.FirstLogonCommands
	}
	if sp.GeneralSettings != nil {
		result["general_settings"] = flattenGeneralSettings(sp.GeneralSettings)
	}
	if sp.LocaleSettings != nil {
		result["locale_settings"] = flattenLocaleSettings(sp.LocaleSettings)
	}
	if sp.NetworkSettings != nil {
		result["network_settings"] = flattenNetworkSettings(sp.NetworkSettings)
	}
	if sp.WorkgroupOrDomainInfo != nil {
		result["workgroup_or_domain_info"] = flattenWorkgroupOrDomainInfo(sp.WorkgroupOrDomainInfo)
	}
	return []map[string]interface{}{result}
}

func flattenGeneralSettings(gs *config.VmGcProfileGeneralSettings) []map[string]interface{} {
	if gs == nil {
		return nil
	}
	result := map[string]interface{}{}
	if gs.AdministratorPassword != nil {
		result["administrator_password"] = *gs.AdministratorPassword
	}
	if gs.AutoLogonSettings != nil {
		als := map[string]interface{}{}
		if gs.AutoLogonSettings.LogonCount != nil {
			als["logon_count"] = *gs.AutoLogonSettings.LogonCount
		}
		result["auto_logon_settings"] = []map[string]interface{}{als}
	}
	if gs.ComputerName != nil {
		result["computer_name"] = flattenComputerName(gs.ComputerName)
	}
	if gs.RegisteredOrganization != nil {
		result["registered_organization"] = *gs.RegisteredOrganization
	}
	if gs.RegisteredOwner != nil {
		result["registered_owner"] = *gs.RegisteredOwner
	}
	if gs.Timezone != nil {
		result["timezone"] = *gs.Timezone
	}
	if gs.WindowsProductKey != nil {
		result["windows_product_key"] = *gs.WindowsProductKey
	}
	return []map[string]interface{}{result}
}

func flattenComputerName(oneOf *config.OneOfVmGcProfileGeneralSettingsComputerName) []map[string]interface{} {
	if oneOf == nil {
		return nil
	}
	result := map[string]interface{}{}
	val := oneOf.GetValue()
	switch val.(type) {
	case config.VmGcProfileUseVmName:
		result["use_vm_name"] = true
		result["must_provide_during_deployment"] = false
	case config.VmGcProfileMustProvideDuringDeployment:
		result["use_vm_name"] = false
		result["must_provide_during_deployment"] = true
	}
	return []map[string]interface{}{result}
}

func flattenLocaleSettings(ls *config.VmGcProfileLocaleSettings) []map[string]interface{} {
	if ls == nil {
		return nil
	}
	result := map[string]interface{}{}
	if ls.SystemLocale != nil {
		result["system_locale"] = *ls.SystemLocale
	}
	if ls.UiLanguage != nil {
		result["ui_language"] = *ls.UiLanguage
	}
	if ls.UserLocale != nil {
		result["user_locale"] = *ls.UserLocale
	}
	return []map[string]interface{}{result}
}

func flattenNetworkSettings(ns *config.VmGcProfileNetworkSettings) []map[string]interface{} {
	if ns == nil {
		return nil
	}
	nicList := make([]map[string]interface{}, len(ns.NicConfigList))
	for i, nic := range ns.NicConfigList {
		nicMap := map[string]interface{}{}
		if nic.DnsConfig != nil {
			dnsMap := map[string]interface{}{}
			if nic.DnsConfig.PreferredDnsServerAddress != nil {
				dnsMap["preferred_dns_server_address"] = *nic.DnsConfig.PreferredDnsServerAddress
			}
			if nic.DnsConfig.AlternateDnsServerAddresses != nil {
				dnsMap["alternate_dns_server_addresses"] = nic.DnsConfig.AlternateDnsServerAddresses
			}
			nicMap["dns_config"] = []map[string]interface{}{dnsMap}
		}
		if nic.Ipv4Config != nil {
			ipv4Map := map[string]interface{}{}
			val := nic.Ipv4Config.GetValue()
			switch val.(type) {
			case config.VmGcProfileUseDhcp:
				ipv4Map["use_dhcp"] = true
				ipv4Map["must_provide_during_deployment"] = false
			case config.VmGcProfileMustProvideDuringDeployment:
				ipv4Map["use_dhcp"] = false
				ipv4Map["must_provide_during_deployment"] = true
			}
			nicMap["ipv4_config"] = []map[string]interface{}{ipv4Map}
		}
		nicList[i] = nicMap
	}
	return []map[string]interface{}{{
		"nic_config_list": nicList,
	}}
}

func flattenWorkgroupOrDomainInfo(oneOf *config.OneOfVmGcProfileSysprepParamsWorkgroupOrDomainInfo) []map[string]interface{} {
	if oneOf == nil {
		return nil
	}
	result := map[string]interface{}{}
	val := oneOf.GetValue()
	switch v := val.(type) {
	case config.VmGcProfileWorkgroup:
		wg := map[string]interface{}{}
		if v.Name != nil {
			wg["name"] = *v.Name
		}
		result["workgroup"] = []map[string]interface{}{wg}
	case config.VmGcProfileDomainSettings:
		ds := map[string]interface{}{}
		if v.Credentials != nil {
			creds := map[string]interface{}{}
			if v.Credentials.DomainName != nil {
				creds["domain_name"] = *v.Credentials.DomainName
			}
			if v.Credentials.Password != nil {
				creds["password"] = *v.Credentials.Password
			}
			if v.Credentials.Username != nil {
				creds["username"] = *v.Credentials.Username
			}
			ds["credentials"] = []map[string]interface{}{creds}
		}
		result["domain_settings"] = []map[string]interface{}{ds}
	}
	return []map[string]interface{}{result}
}

func flattenAnswerFile(af *config.VmGcProfileAnswerFile) []map[string]interface{} {
	if af == nil {
		return nil
	}
	result := map[string]interface{}{}
	if af.UnattendXml != nil {
		result["unattend_xml"] = *af.UnattendXml
	}
	return []map[string]interface{}{result}
}
