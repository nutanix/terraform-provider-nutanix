package vmmv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import3 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/common/v1/response"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func vmGcMaxItems(isDataSource bool) int {
	if isDataSource {
		return 0
	}
	return 1
}

func schemaForVmGcProfileConfig(isDataSource bool) *schema.Schema {
	s := &schema.Schema{
		Type:     schema.TypeList,
		Optional: !isDataSource,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"sysprep_config": {
					Type:     schema.TypeList,
					Optional: !isDataSource,
					Computed: isDataSource,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"customization": {
								Type:     schema.TypeList,
								Optional: !isDataSource,
								Computed: isDataSource,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"sysprep_params": {
											Type:     schema.TypeList,
											Optional: !isDataSource,
											Computed: isDataSource,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"first_logon_commands": {
														Type:     schema.TypeList,
														Optional: !isDataSource,
														Computed: isDataSource,
														Elem: &schema.Schema{
															Type: schema.TypeString,
														},
													},
													"general_settings": schemaForVmGcProfileGeneralSettings(isDataSource),
													"locale_settings":  schemaForVmGcProfileLocaleSettings(isDataSource),
													"network_settings": schemaForVmGcProfileNetworkSettings(isDataSource),
													"workgroup_or_domain_info": {
														Type:     schema.TypeList,
														Optional: !isDataSource,
														Computed: isDataSource,
														Elem: &schema.Resource{
															Schema: map[string]*schema.Schema{
																"workgroup": vmGcNestedBlock(isDataSource, map[string]*schema.Schema{
																	"name": {
																		Type:     schema.TypeString,
																		Optional: !isDataSource,
																		Computed: isDataSource,
																	},
																}),
																"domain_settings": vmGcNestedBlock(isDataSource, map[string]*schema.Schema{
																	"credentials": vmGcNestedBlock(isDataSource, map[string]*schema.Schema{
																		"domain_name": {
																			Type:     schema.TypeString,
																			Optional: !isDataSource,
																			Computed: isDataSource,
																		},
																		"password": {
																			Type:      schema.TypeString,
																			Optional:  !isDataSource,
																			Computed:  isDataSource,
																			Sensitive: true,
																		},
																		"username": {
																			Type:     schema.TypeString,
																			Optional: !isDataSource,
																			Computed: isDataSource,
																		},
																	}),
																}),
															},
														},
													},
												},
											},
										},
										"answer_file": vmGcNestedBlock(isDataSource, map[string]*schema.Schema{
											"unattend_xml": {
												Type:     schema.TypeString,
												Optional: !isDataSource,
												Computed: isDataSource,
											},
										}),
									},
								},
							},
						},
					},
				},
			},
		},
	}
	if !isDataSource {
		s.MaxItems = 1
	}
	return s
}

func vmGcNestedBlock(isDataSource bool, innerSchema map[string]*schema.Schema) *schema.Schema {
	s := &schema.Schema{
		Type:     schema.TypeList,
		Optional: !isDataSource,
		Computed: isDataSource,
		Elem: &schema.Resource{
			Schema: innerSchema,
		},
	}
	if !isDataSource {
		s.MaxItems = 1
	}
	return s
}

func schemaForVmGcProfileGeneralSettings(isDataSource bool) *schema.Schema {
	s := &schema.Schema{
		Type:     schema.TypeList,
		Optional: !isDataSource,
		Computed: isDataSource,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"administrator_password": {
					Type:      schema.TypeString,
					Optional:  !isDataSource,
					Computed:  isDataSource,
					Sensitive: true,
				},
				"auto_logon_settings": vmGcNestedBlock(isDataSource, map[string]*schema.Schema{
					"logon_count": {
						Type:     schema.TypeInt,
						Optional: !isDataSource,
						Computed: isDataSource,
					},
				}),
				"computer_name": vmGcNestedBlock(isDataSource, map[string]*schema.Schema{
					"must_provide_during_deployment": {
						Type:     schema.TypeBool,
						Optional: !isDataSource,
						Computed: isDataSource,
					},
					"use_vm_name": {
						Type:     schema.TypeBool,
						Optional: !isDataSource,
						Computed: isDataSource,
					},
				}),
				"registered_organization": {
					Type:     schema.TypeString,
					Optional: !isDataSource,
					Computed: isDataSource,
				},
				"registered_owner": {
					Type:     schema.TypeString,
					Optional: !isDataSource,
					Computed: isDataSource,
				},
				"timezone": {
					Type:     schema.TypeString,
					Optional: !isDataSource,
					Computed: isDataSource,
				},
				"windows_product_key": {
					Type:      schema.TypeString,
					Optional:  !isDataSource,
					Computed:  isDataSource,
					Sensitive: true,
				},
			},
		},
	}
	if !isDataSource {
		s.MaxItems = 1
	}
	return s
}

func schemaForVmGcProfileLocaleSettings(isDataSource bool) *schema.Schema {
	s := &schema.Schema{
		Type:     schema.TypeList,
		Optional: !isDataSource,
		Computed: isDataSource,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"system_locale": {
					Type:     schema.TypeString,
					Optional: !isDataSource,
					Computed: isDataSource,
				},
				"ui_language": {
					Type:     schema.TypeString,
					Optional: !isDataSource,
					Computed: isDataSource,
				},
				"user_locale": {
					Type:     schema.TypeString,
					Optional: !isDataSource,
					Computed: isDataSource,
				},
			},
		},
	}
	if !isDataSource {
		s.MaxItems = 1
	}
	return s
}

func schemaForVmGcProfileNetworkSettings(isDataSource bool) *schema.Schema {
	s := &schema.Schema{
		Type:     schema.TypeList,
		Optional: !isDataSource,
		Computed: isDataSource,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"nic_config_list": {
					Type:     schema.TypeList,
					Optional: !isDataSource,
					Computed: isDataSource,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"dns_config": vmGcNestedBlock(isDataSource, map[string]*schema.Schema{
								"alternate_dns_server_addresses": {
									Type:     schema.TypeList,
									Optional: !isDataSource,
									Computed: isDataSource,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								"preferred_dns_server_address": {
									Type:     schema.TypeString,
									Optional: !isDataSource,
									Computed: isDataSource,
								},
							}),
							"ipv4_config": vmGcNestedBlock(isDataSource, map[string]*schema.Schema{
								"use_dhcp": {
									Type:     schema.TypeBool,
									Optional: !isDataSource,
									Computed: isDataSource,
								},
								"must_provide_during_deployment": {
									Type:     schema.TypeBool,
									Optional: !isDataSource,
									Computed: isDataSource,
								},
							}),
						},
					},
				},
			},
		},
	}
	if !isDataSource {
		s.MaxItems = 1
	}
	return s
}

func flattenVmGcProfileLinks(links []import3.ApiLink) []map[string]interface{} {
	if len(links) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(links))
	for i, link := range links {
		result[i] = map[string]interface{}{
			"href": utils.StringValue(link.Href),
			"rel":  utils.StringValue(link.Rel),
		}
	}
	return result
}

func flattenVmGcProfileUserReference(ref *config.UserReference) []map[string]interface{} {
	if ref == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"ext_id": utils.StringValue(ref.ExtId),
		},
	}
}

func flattenVmGcProfileConfig(cfg *config.OneOfVmGuestCustomizationProfileConfig) []map[string]interface{} {
	if cfg == nil || cfg.ObjectType_ == nil {
		return nil
	}

	configMap := make(map[string]interface{})
	val := cfg.GetValue()
	if val == nil {
		return nil
	}

	switch *cfg.ObjectType_ {
	case "vmm.v4.ahv.config.VmGcProfileSysprepConfig":
		sysprepConfig, ok := val.(config.VmGcProfileSysprepConfig)
		if !ok {
			return nil
		}
		configMap["sysprep_config"] = flattenVmGcProfileSysprepConfig(&sysprepConfig)
	}

	return []map[string]interface{}{configMap}
}

func flattenVmGcProfileSysprepConfig(sc *config.VmGcProfileSysprepConfig) []map[string]interface{} {
	if sc == nil {
		return nil
	}
	result := map[string]interface{}{
		"customization": flattenVmGcProfileSysprepCustomization(sc.Customization),
	}
	return []map[string]interface{}{result}
}

func flattenVmGcProfileSysprepCustomization(cust *config.OneOfVmGcProfileSysprepConfigCustomization) []map[string]interface{} {
	if cust == nil || cust.ObjectType_ == nil {
		return nil
	}

	custMap := make(map[string]interface{})
	val := cust.GetValue()
	if val == nil {
		return nil
	}

	switch *cust.ObjectType_ {
	case "vmm.v4.ahv.config.VmGcProfileSysprepParams":
		params, ok := val.(config.VmGcProfileSysprepParams)
		if !ok {
			return nil
		}
		custMap["sysprep_params"] = flattenVmGcProfileSysprepParams(&params)
		custMap["answer_file"] = make([]interface{}, 0)
	case "vmm.v4.ahv.config.VmGcProfileAnswerFile":
		answerFile, ok := val.(config.VmGcProfileAnswerFile)
		if !ok {
			return nil
		}
		custMap["answer_file"] = []map[string]interface{}{
			{
				"unattend_xml": utils.StringValue(answerFile.UnattendXml),
			},
		}
		custMap["sysprep_params"] = make([]interface{}, 0)
	}

	return []map[string]interface{}{custMap}
}

func flattenVmGcProfileSysprepParams(params *config.VmGcProfileSysprepParams) []map[string]interface{} {
	if params == nil {
		return nil
	}
	result := map[string]interface{}{
		"first_logon_commands":     params.FirstLogonCommands,
		"general_settings":         flattenVmGcProfileGeneralSettings(params.GeneralSettings),
		"locale_settings":          flattenVmGcProfileLocaleSettings(params.LocaleSettings),
		"network_settings":         flattenVmGcProfileNetworkSettings(params.NetworkSettings),
		"workgroup_or_domain_info": flattenVmGcProfileWorkgroupOrDomainInfo(params.WorkgroupOrDomainInfo),
	}
	return []map[string]interface{}{result}
}

func flattenVmGcProfileGeneralSettings(gs *config.VmGcProfileGeneralSettings) []map[string]interface{} {
	if gs == nil {
		return nil
	}
	result := map[string]interface{}{
		"administrator_password":  utils.StringValue(gs.AdministratorPassword),
		"auto_logon_settings":     flattenVmGcProfileAutoLogonSettings(gs.AutoLogonSettings),
		"computer_name":           flattenVmGcProfileComputerName(gs.ComputerName),
		"registered_organization": utils.StringValue(gs.RegisteredOrganization),
		"registered_owner":        utils.StringValue(gs.RegisteredOwner),
		"timezone":                utils.StringValue(gs.Timezone),
		"windows_product_key":     utils.StringValue(gs.WindowsProductKey),
	}
	return []map[string]interface{}{result}
}

func flattenVmGcProfileAutoLogonSettings(als *config.VmGcProfileAutoLogonSettings) []map[string]interface{} {
	if als == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"logon_count": utils.IntValue(als.LogonCount),
		},
	}
}

func flattenVmGcProfileComputerName(cn *config.OneOfVmGcProfileGeneralSettingsComputerName) []map[string]interface{} {
	if cn == nil || cn.ObjectType_ == nil {
		return nil
	}

	val := cn.GetValue()
	if val == nil {
		return nil
	}

	switch *cn.ObjectType_ {
	case "vmm.v4.ahv.config.VmGcProfileMustProvideDuringDeployment":
		return []map[string]interface{}{
			{
				"must_provide_during_deployment": true,
				"use_vm_name":                    false,
			},
		}
	case "vmm.v4.ahv.config.VmGcProfileUseVmName":
		return []map[string]interface{}{
			{
				"must_provide_during_deployment": false,
				"use_vm_name":                    true,
			},
		}
	}
	return nil
}

func flattenVmGcProfileLocaleSettings(ls *config.VmGcProfileLocaleSettings) []map[string]interface{} {
	if ls == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"system_locale": utils.StringValue(ls.SystemLocale),
			"ui_language":   utils.StringValue(ls.UiLanguage),
			"user_locale":   utils.StringValue(ls.UserLocale),
		},
	}
}

func flattenVmGcProfileNetworkSettings(ns *config.VmGcProfileNetworkSettings) []map[string]interface{} {
	if ns == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"nic_config_list": flattenVmGcProfileNicConfigList(ns.NicConfigList),
		},
	}
}

func flattenVmGcProfileNicConfigList(nics []config.VmGcProfileNicConfig) []map[string]interface{} {
	if len(nics) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(nics))
	for i, nic := range nics {
		result[i] = map[string]interface{}{
			"dns_config":  flattenVmGcProfileDnsConfig(nic.DnsConfig),
			"ipv4_config": flattenVmGcProfileIpv4Config(nic.Ipv4Config),
		}
	}
	return result
}

func flattenVmGcProfileDnsConfig(dc *config.VmGcProfileDnsConfig) []map[string]interface{} {
	if dc == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"alternate_dns_server_addresses": dc.AlternateDnsServerAddresses,
			"preferred_dns_server_address":   utils.StringValue(dc.PreferredDnsServerAddress),
		},
	}
}

func flattenVmGcProfileIpv4Config(ic *config.OneOfVmGcProfileNicConfigIpv4Config) []map[string]interface{} {
	if ic == nil || ic.ObjectType_ == nil {
		return nil
	}

	val := ic.GetValue()
	if val == nil {
		return nil
	}

	switch *ic.ObjectType_ {
	case "vmm.v4.ahv.config.VmGcProfileUseDhcp":
		return []map[string]interface{}{
			{
				"use_dhcp":                       true,
				"must_provide_during_deployment": false,
			},
		}
	case "vmm.v4.ahv.config.VmGcProfileMustProvideDuringDeployment":
		return []map[string]interface{}{
			{
				"use_dhcp":                       false,
				"must_provide_during_deployment": true,
			},
		}
	}
	return nil
}

func flattenVmGcProfileWorkgroupOrDomainInfo(wodi *config.OneOfVmGcProfileSysprepParamsWorkgroupOrDomainInfo) []map[string]interface{} {
	if wodi == nil || wodi.ObjectType_ == nil {
		return nil
	}

	val := wodi.GetValue()
	if val == nil {
		return nil
	}

	wodiMap := make(map[string]interface{})
	switch *wodi.ObjectType_ {
	case "vmm.v4.ahv.config.VmGcProfileWorkgroup":
		wg, ok := val.(config.VmGcProfileWorkgroup)
		if !ok {
			return nil
		}
		wodiMap["workgroup"] = []map[string]interface{}{
			{
				"name": utils.StringValue(wg.Name),
			},
		}
		wodiMap["domain_settings"] = make([]interface{}, 0)
	case "vmm.v4.ahv.config.VmGcProfileDomainSettings":
		ds, ok := val.(config.VmGcProfileDomainSettings)
		if !ok {
			return nil
		}
		wodiMap["domain_settings"] = flattenVmGcProfileDomainSettings(&ds)
		wodiMap["workgroup"] = make([]interface{}, 0)
	}

	return []map[string]interface{}{wodiMap}
}

func flattenVmGcProfileDomainSettings(ds *config.VmGcProfileDomainSettings) []map[string]interface{} {
	if ds == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"credentials": flattenVmGcProfileDomainCredentials(ds.Credentials),
		},
	}
}

func flattenVmGcProfileDomainCredentials(creds *config.VmGcProfileDomainCredentials) []map[string]interface{} {
	if creds == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"domain_name": utils.StringValue(creds.DomainName),
			"password":    utils.StringValue(creds.Password),
			"username":    utils.StringValue(creds.Username),
		},
	}
}

func expandVmGcProfileConfig(cfgList []interface{}) *config.OneOfVmGuestCustomizationProfileConfig {
	if len(cfgList) == 0 || cfgList[0] == nil {
		return nil
	}

	cfgMap := cfgList[0].(map[string]interface{})

	if sysprepConfigList, ok := cfgMap["sysprep_config"].([]interface{}); ok && len(sysprepConfigList) > 0 {
		sysprepConfig := expandVmGcProfileSysprepConfig(sysprepConfigList)
		if sysprepConfig != nil {
			oneOf := config.NewOneOfVmGuestCustomizationProfileConfig()
			if err := oneOf.SetValue(*sysprepConfig); err == nil {
				return oneOf
			}
		}
	}

	return nil
}

func expandVmGcProfileSysprepConfig(sysprepList []interface{}) *config.VmGcProfileSysprepConfig {
	if len(sysprepList) == 0 || sysprepList[0] == nil {
		return nil
	}

	sysprepMap := sysprepList[0].(map[string]interface{})
	sc := config.NewVmGcProfileSysprepConfig()

	if custList, ok := sysprepMap["customization"].([]interface{}); ok && len(custList) > 0 {
		sc.Customization = expandVmGcProfileSysprepCustomization(custList)
	}

	return sc
}

func expandVmGcProfileSysprepCustomization(custList []interface{}) *config.OneOfVmGcProfileSysprepConfigCustomization {
	if len(custList) == 0 || custList[0] == nil {
		return nil
	}

	custMap := custList[0].(map[string]interface{})
	oneOf := config.NewOneOfVmGcProfileSysprepConfigCustomization()

	if paramsList, ok := custMap["sysprep_params"].([]interface{}); ok && len(paramsList) > 0 && paramsList[0] != nil {
		params := expandVmGcProfileSysprepParams(paramsList)
		if params != nil {
			if err := oneOf.SetValue(*params); err == nil {
				return oneOf
			}
		}
	}

	if answerFileList, ok := custMap["answer_file"].([]interface{}); ok && len(answerFileList) > 0 && answerFileList[0] != nil {
		answerFile := expandVmGcProfileAnswerFile(answerFileList)
		if answerFile != nil {
			if err := oneOf.SetValue(*answerFile); err == nil {
				return oneOf
			}
		}
	}

	return nil
}

func expandVmGcProfileSysprepParams(paramsList []interface{}) *config.VmGcProfileSysprepParams {
	if len(paramsList) == 0 || paramsList[0] == nil {
		return nil
	}

	paramsMap := paramsList[0].(map[string]interface{})
	params := config.NewVmGcProfileSysprepParams()

	if v, ok := paramsMap["first_logon_commands"].([]interface{}); ok && len(v) > 0 {
		cmds := make([]string, len(v))
		for i, cmd := range v {
			cmds[i] = cmd.(string)
		}
		params.FirstLogonCommands = cmds
	}

	if gsList, ok := paramsMap["general_settings"].([]interface{}); ok && len(gsList) > 0 {
		params.GeneralSettings = expandVmGcProfileGeneralSettings(gsList)
	}

	if lsList, ok := paramsMap["locale_settings"].([]interface{}); ok && len(lsList) > 0 {
		params.LocaleSettings = expandVmGcProfileLocaleSettings(lsList)
	}

	if nsList, ok := paramsMap["network_settings"].([]interface{}); ok && len(nsList) > 0 {
		params.NetworkSettings = expandVmGcProfileNetworkSettings(nsList)
	}

	if wodiList, ok := paramsMap["workgroup_or_domain_info"].([]interface{}); ok && len(wodiList) > 0 {
		params.WorkgroupOrDomainInfo = expandVmGcProfileWorkgroupOrDomainInfo(wodiList)
	}

	return params
}

func expandVmGcProfileAnswerFile(afList []interface{}) *config.VmGcProfileAnswerFile {
	if len(afList) == 0 || afList[0] == nil {
		return nil
	}

	afMap := afList[0].(map[string]interface{})
	af := config.NewVmGcProfileAnswerFile()

	if v, ok := afMap["unattend_xml"].(string); ok && v != "" {
		af.UnattendXml = utils.StringPtr(v)
	}

	return af
}

func expandVmGcProfileGeneralSettings(gsList []interface{}) *config.VmGcProfileGeneralSettings {
	if len(gsList) == 0 || gsList[0] == nil {
		return nil
	}

	gsMap := gsList[0].(map[string]interface{})
	gs := config.NewVmGcProfileGeneralSettings()

	if v, ok := gsMap["administrator_password"].(string); ok && v != "" {
		gs.AdministratorPassword = utils.StringPtr(v)
	}

	if alsList, ok := gsMap["auto_logon_settings"].([]interface{}); ok && len(alsList) > 0 && alsList[0] != nil {
		alsMap := alsList[0].(map[string]interface{})
		gs.AutoLogonSettings = config.NewVmGcProfileAutoLogonSettings()
		if v, ok := alsMap["logon_count"].(int); ok {
			gs.AutoLogonSettings.LogonCount = utils.IntPtr(v)
		}
	}

	if cnList, ok := gsMap["computer_name"].([]interface{}); ok && len(cnList) > 0 && cnList[0] != nil {
		cnMap := cnList[0].(map[string]interface{})
		oneOf := config.NewOneOfVmGcProfileGeneralSettingsComputerName()
		if v, ok := cnMap["use_vm_name"].(bool); ok && v {
			useVmName := config.NewVmGcProfileUseVmName()
			if err := oneOf.SetValue(*useVmName); err == nil {
				gs.ComputerName = oneOf
			}
		} else if v, ok := cnMap["must_provide_during_deployment"].(bool); ok && v {
			mustProvide := config.NewVmGcProfileMustProvideDuringDeployment()
			if err := oneOf.SetValue(*mustProvide); err == nil {
				gs.ComputerName = oneOf
			}
		}
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

	return gs
}

func expandVmGcProfileLocaleSettings(lsList []interface{}) *config.VmGcProfileLocaleSettings {
	if len(lsList) == 0 || lsList[0] == nil {
		return nil
	}

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

	return ls
}

func expandVmGcProfileNetworkSettings(nsList []interface{}) *config.VmGcProfileNetworkSettings {
	if len(nsList) == 0 || nsList[0] == nil {
		return nil
	}

	nsMap := nsList[0].(map[string]interface{})
	ns := config.NewVmGcProfileNetworkSettings()

	if nicList, ok := nsMap["nic_config_list"].([]interface{}); ok && len(nicList) > 0 {
		nics := make([]config.VmGcProfileNicConfig, len(nicList))
		for i, nicRaw := range nicList {
			nicMap := nicRaw.(map[string]interface{})
			nic := *config.NewVmGcProfileNicConfig()

			if dnsList, ok := nicMap["dns_config"].([]interface{}); ok && len(dnsList) > 0 && dnsList[0] != nil {
				dnsMap := dnsList[0].(map[string]interface{})
				nic.DnsConfig = config.NewVmGcProfileDnsConfig()
				if v, ok := dnsMap["preferred_dns_server_address"].(string); ok && v != "" {
					nic.DnsConfig.PreferredDnsServerAddress = utils.StringPtr(v)
				}
				if v, ok := dnsMap["alternate_dns_server_addresses"].([]interface{}); ok && len(v) > 0 {
					addrs := make([]string, len(v))
					for j, addr := range v {
						addrs[j] = addr.(string)
					}
					nic.DnsConfig.AlternateDnsServerAddresses = addrs
				}
			}

			if ipv4List, ok := nicMap["ipv4_config"].([]interface{}); ok && len(ipv4List) > 0 && ipv4List[0] != nil {
				ipv4Map := ipv4List[0].(map[string]interface{})
				oneOf := config.NewOneOfVmGcProfileNicConfigIpv4Config()
				if v, ok := ipv4Map["use_dhcp"].(bool); ok && v {
					useDhcp := config.NewVmGcProfileUseDhcp()
					if err := oneOf.SetValue(*useDhcp); err == nil {
						nic.Ipv4Config = oneOf
					}
				} else if v, ok := ipv4Map["must_provide_during_deployment"].(bool); ok && v {
					mustProvide := config.NewVmGcProfileMustProvideDuringDeployment()
					if err := oneOf.SetValue(*mustProvide); err == nil {
						nic.Ipv4Config = oneOf
					}
				}
			}

			nics[i] = nic
		}
		ns.NicConfigList = nics
	}

	return ns
}

func expandVmGcProfileWorkgroupOrDomainInfo(wodiList []interface{}) *config.OneOfVmGcProfileSysprepParamsWorkgroupOrDomainInfo {
	if len(wodiList) == 0 || wodiList[0] == nil {
		return nil
	}

	wodiMap := wodiList[0].(map[string]interface{})
	oneOf := config.NewOneOfVmGcProfileSysprepParamsWorkgroupOrDomainInfo()

	if wgList, ok := wodiMap["workgroup"].([]interface{}); ok && len(wgList) > 0 && wgList[0] != nil {
		wgMap := wgList[0].(map[string]interface{})
		wg := config.NewVmGcProfileWorkgroup()
		if v, ok := wgMap["name"].(string); ok && v != "" {
			wg.Name = utils.StringPtr(v)
		}
		if err := oneOf.SetValue(*wg); err == nil {
			return oneOf
		}
	}

	if dsList, ok := wodiMap["domain_settings"].([]interface{}); ok && len(dsList) > 0 && dsList[0] != nil {
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
		if err := oneOf.SetValue(*ds); err == nil {
			return oneOf
		}
	}

	return nil
}
