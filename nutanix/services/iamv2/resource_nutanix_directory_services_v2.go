package iamv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/iam/v4/authn"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixDirectoryServicesV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixDirectoryServicesV2Create,
		ReadContext:   ResourceNutanixDirectoryServicesV2Read,
		UpdateContext: ResourceNutanixDirectoryServicesV2Update,
		DeleteContext: ResourceNutanixDirectoryServicesV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"secondary_urls": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"directory_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_account": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Required: true,
						},
						"password": {
							Type:      schema.TypeString,
							Sensitive: true,
							Required:  true,
						},
					},
				},
			},
			"open_ldap_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_configuration": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"user_object_class": {
										Type:     schema.TypeString,
										Required: true,
									},
									"user_search_base": {
										Type:     schema.TypeString,
										Required: true,
									},
									"username_attribute": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"user_group_configuration": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_object_class": {
										Type:     schema.TypeString,
										Required: true,
									},
									"group_search_base": {
										Type:     schema.TypeString,
										Required: true,
									},
									"group_member_attribute": {
										Type:     schema.TypeString,
										Required: true,
									},
									"group_member_attribute_value": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"group_search_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"white_listed_groups": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_updated_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func ResourceNutanixDirectoryServicesV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	input := &import1.DirectoryService{}

	if name, ok := d.GetOk("name"); ok {
		input.Name = utils.StringPtr(name.(string))
	}
	if url, ok := d.GetOk("url"); ok {
		input.Url = utils.StringPtr(url.(string))
	}
	if secUrls, ok := d.GetOk("secondary_urls"); ok {
		secondaryUrlsList := secUrls.([]interface{})
		secondaryUrlsListStr := make([]string, len(secondaryUrlsList))
		for i, v := range secondaryUrlsList {
			secondaryUrlsListStr[i] = v.(string)
		}
		input.SecondaryUrls = secondaryUrlsListStr
	}
	if domainName, ok := d.GetOk("domain_name"); ok {
		input.DomainName = utils.StringPtr(domainName.(string))
	}
	if dType, ok := d.GetOk("directory_type"); ok {
		const two, three = 2, 3
		subMap := map[string]interface{}{
			"ACTIVE_DIRECTORY": two,
			"OPEN_LDAP":        three,
		}
		pInt := subMap[dType.(string)]
		p := import1.DirectoryType(pInt.(int))

		input.DirectoryType = &p
	}
	if serviceAcc, ok := d.GetOk("service_account"); ok {
		input.ServiceAccount = expandDsServiceAccount(serviceAcc)
	}
	if ldap, ok := d.GetOk("open_ldap_configuration"); ok {
		input.OpenLdapConfiguration = expandOpenLdapConfig(ldap)
	}
	if grpSearchType, ok := d.GetOk("group_search_type"); ok {
		const two, three = 2, 3
		subMap := map[string]interface{}{
			"NON_RECURSIVE": two,
			"RECURSIVE":     three,
		}
		pInt := subMap[grpSearchType.(string)]
		p := import1.GroupSearchType(pInt.(int))
		input.GroupSearchType = &p
	}
	if whitelistedGrp, ok := d.GetOk("white_listed_groups"); ok {
		whitelistedGrpList := whitelistedGrp.([]interface{})
		whitelistedGrpListStr := make([]string, len(whitelistedGrpList))
		for i, v := range whitelistedGrpList {
			whitelistedGrpListStr[i] = v.(string)
		}
		input.WhiteListedGroups = whitelistedGrpListStr
	}

	aJSON, _ := json.MarshalIndent(input, "", " ")
	log.Println("[DEBUG] Directory Service create payload: ", string(aJSON))

	resp, err := conn.DirectoryServiceAPIInstance.CreateDirectoryService(input)
	if err != nil {
		return diag.Errorf("error while creating directory services : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.DirectoryService)

	d.SetId(*getResp.ExtId)
	return ResourceNutanixDirectoryServicesV2Read(ctx, d, meta)
}

func ResourceNutanixDirectoryServicesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	resp, err := conn.DirectoryServiceAPIInstance.GetDirectoryServiceById(utils.StringPtr(d.Id()))
	if err != nil {
		var errordata map[string]interface{}
		e := json.Unmarshal([]byte(err.Error()), &errordata)
		if e != nil {
			return diag.FromErr(e)
		}
		data := errordata["data"].(map[string]interface{})
		errorList := data["error"].([]interface{})
		errorMessage := errorList[0].(map[string]interface{})
		return diag.Errorf("error while fetching directory services: %v", errorMessage["message"])
	}

	getResp := resp.Data.GetValue().(import1.DirectoryService)
	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("url", getResp.Url); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("secondary_urls", getResp.SecondaryUrls); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("domain_name", getResp.DomainName); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("directory_type", flattenDirectoryType(getResp.DirectoryType)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("service_account", flattenDsServiceAccountForResource(getResp.ServiceAccount)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("open_ldap_configuration", flattenOpenLdapConfig(getResp.OpenLdapConfiguration)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("group_search_type", flattenGroupSearchType(getResp.GroupSearchType)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("white_listed_groups", getResp.WhiteListedGroups); err != nil {
		return diag.FromErr(err)
	}

	if getResp.CreatedTime != nil {
		t := getResp.CreatedTime
		if err := d.Set("created_time", t.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if getResp.LastUpdatedTime != nil {
		t := getResp.LastUpdatedTime
		if err := d.Set("last_updated_time", t.String()); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("created_by", getResp.CreatedBy); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func ResourceNutanixDirectoryServicesV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI
	updatedSpec := import1.DirectoryService{}

	readResp, err := conn.DirectoryServiceAPIInstance.GetDirectoryServiceById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching Directory service : %v", err)
	}
	// get etag value from read response to pass in update request If-Match header, Required for update request
	etagValue := conn.SamlIdentityAPIInstance.ApiClient.GetEtag(readResp)
	headers := make(map[string]interface{})
	headers["If-Match"] = utils.StringPtr(etagValue)

	updatedSpec = readResp.Data.GetValue().(import1.DirectoryService)

	// service account required for update
	// expand service account from row config as it required for update and
	// password is saved as sensitive value
	rawConfig := d.GetRawConfig()
	configMap := rawConfig.AsValueMap()

	val, exists := configMap["service_account"]
	if !exists || val.IsNull() || !val.IsKnown() {
		return diag.Errorf("service_account is missing or unknown")
	}

	if !val.Type().IsTupleType() && !val.Type().IsListType() {
		return diag.Errorf("service_account is not a list")
	}

	serviceAccounts := make([]interface{}, 0)

	// Extract each element from the list
	for _, item := range val.AsValueSlice() {
		if item.IsNull() || !item.IsKnown() || !item.Type().IsObjectType() {
			continue
		}

		account := make(map[string]interface{})
		for key, value := range item.AsValueMap() {
			if value.IsNull() || !value.IsKnown() {
				continue
			}
			account[key] = value.AsString()
		}

		serviceAccounts = append(serviceAccounts, account)
	}

	updatedSpec.ServiceAccount = expandDsServiceAccount(serviceAccounts)

	if d.HasChange("name") {
		updatedSpec.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("url") {
		updatedSpec.Url = utils.StringPtr(d.Get("url").(string))
	}
	if d.HasChange("secondary_urls") {
		secUrls := d.Get("secondary_urls")
		secondaryUrlsList := secUrls.([]interface{})
		secondaryUrlsListStr := make([]string, len(secondaryUrlsList))
		for i, v := range secondaryUrlsList {
			secondaryUrlsListStr[i] = v.(string)
		}
		updatedSpec.SecondaryUrls = secondaryUrlsListStr
	}
	if d.HasChange("domain_name") {
		updatedSpec.DomainName = utils.StringPtr(d.Get("domain_name").(string))
	}
	if d.HasChange("directory_type") {
		const two, three = 2, 3
		subMap := map[string]interface{}{
			"ACTIVE_DIRECTORY": two,
			"OPEN_LDAP":        three,
		}
		pInt := subMap[d.Get("directory_type").(string)]
		p := import1.DirectoryType(pInt.(int))

		updatedSpec.DirectoryType = &p
	}
	if d.HasChange("open_ldap_configuration") {
		updatedSpec.OpenLdapConfiguration = expandOpenLdapConfig(d.Get("open_ldap_configuration"))
	}
	if d.HasChange("group_search_type") {
		const two, three = 2, 3
		subMap := map[string]interface{}{
			"NON_RECURSIVE": two,
			"RECURSIVE":     three,
		}
		pInt := subMap[d.Get("group_search_type").(string)]
		p := import1.GroupSearchType(pInt.(int))
		updatedSpec.GroupSearchType = &p
	}
	if d.HasChange("white_listed_groups") {
		whitelistedGrp := d.Get("white_listed_groups")
		whitelistedGrpList := whitelistedGrp.([]interface{})
		whitelistedGrpListStr := make([]string, len(whitelistedGrpList))
		for i, v := range whitelistedGrpList {
			whitelistedGrpListStr[i] = v.(string)
		}
		updatedSpec.WhiteListedGroups = whitelistedGrpListStr
	}

	aJSON, _ := json.MarshalIndent(updatedSpec, "", " ")
	log.Println("[DEBUG] Directory Service update payload: ", string(aJSON))

	updatedResp, err := conn.DirectoryServiceAPIInstance.UpdateDirectoryServiceById(utils.StringPtr(d.Id()), &updatedSpec, headers)
	if err != nil {
		return diag.Errorf("error while updating directory services: %v", err)
	}

	updatedResponse := updatedResp.Data.GetValue().(import1.DirectoryService)

	if updatedResponse.ExtId != nil {
		log.Println("[DEBUG] updated the directory services")
	}
	return ResourceNutanixDirectoryServicesV2Read(ctx, d, meta)
}

func ResourceNutanixDirectoryServicesV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	readResp, err := conn.DirectoryServiceAPIInstance.GetDirectoryServiceById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching Directory service : %v", err)
	}
	// get etag value from read response to pass in update request If-Match header, Required for update request
	etagValue := conn.SamlIdentityAPIInstance.ApiClient.GetEtag(readResp)
	headers := make(map[string]interface{})
	headers["If-Match"] = utils.StringPtr(etagValue)

	resp, err := conn.DirectoryServiceAPIInstance.DeleteDirectoryServiceById(utils.StringPtr(d.Id()), headers)
	if err != nil {
		return diag.Errorf("error while deleting directory services : %v", err)
	}

	if resp == nil {
		log.Println("[DEBUG] Directory Services deleted successfully.")
	}
	return nil
}

func expandDsServiceAccount(pr interface{}) *import1.DsServiceAccount {
	if pr != nil && len(pr.([]interface{})) > 0 {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		ds := &import1.DsServiceAccount{}

		if pass, ok := val["password"]; ok {
			ds.Password = utils.StringPtr(pass.(string))
		}
		if user, ok := val["username"]; ok {
			ds.Username = utils.StringPtr(user.(string))
		}
		return ds
	}
	return nil
}

func expandOpenLdapConfig(pr interface{}) *import1.OpenLdapConfig {
	if pr != nil {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		ldap := &import1.OpenLdapConfig{}

		if userCfg, ok := val["user_configuration"]; ok {
			ldap.UserConfiguration = expandUserConfiguration(userCfg)
		}
		if userGroupCfg, ok := val["user_group_configuration"]; ok {
			ldap.UserGroupConfiguration = expandUserGroupConfiguration(userGroupCfg)
		}
		return ldap
	}
	return nil
}

func expandUserConfiguration(pr interface{}) *import1.UserConfiguration {
	if pr != nil {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		usrcfg := &import1.UserConfiguration{}

		if usrObjClass, ok := val["user_object_class"]; ok {
			usrcfg.UserObjectClass = utils.StringPtr(usrObjClass.(string))
		}
		if usrSearchbase, ok := val["user_search_base"]; ok {
			usrcfg.UserSearchBase = utils.StringPtr(usrSearchbase.(string))
		}
		if usernameAttr, ok := val["username_attribute"]; ok {
			usrcfg.UsernameAttribute = utils.StringPtr(usernameAttr.(string))
		}

		return usrcfg
	}
	return nil
}

func expandUserGroupConfiguration(pr interface{}) *import1.UserGroupConfiguration {
	if pr != nil {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		usrGrp := &import1.UserGroupConfiguration{}

		if grpObjClass, ok := val["group_object_class"]; ok {
			usrGrp.GroupObjectClass = utils.StringPtr(grpObjClass.(string))
		}
		if grpSearchbase, ok := val["group_search_base"]; ok {
			usrGrp.GroupSearchBase = utils.StringPtr(grpSearchbase.(string))
		}
		if grpMemberAttr, ok := val["group_member_attribute"]; ok {
			usrGrp.GroupMemberAttribute = utils.StringPtr(grpMemberAttr.(string))
		}
		if grpAttrVal, ok := val["group_member_attribute_value"]; ok {
			usrGrp.GroupMemberAttributeValue = utils.StringPtr(grpAttrVal.(string))
		}
		return usrGrp
	}
	return nil
}

func flattenDsServiceAccountForResource(pr *import1.DsServiceAccount) []map[string]interface{} {
	if pr != nil {
		accs := make([]map[string]interface{}, 0)
		acc := make(map[string]interface{})

		// password is not returned in the response
		// so we are using the password from the configuration as it required for update
		acc["username"] = pr.Username
		acc["password"] = pr.Password

		accs = append(accs, acc)
		return accs
	}
	return nil
}
