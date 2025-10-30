package iamv2

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import1 "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/iam/v4/authn"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixSamlIdpV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixSamlIdpV2Create,
		ReadContext:   ResourceNutanixSamlIdpV2Read,
		UpdateContext: ResourceNutanixSamlIdpV2Update,
		DeleteContext: ResourceNutanixSamlIdpV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"idp_metadata_url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"idp_metadata": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entity_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"login_url": {
							Type:     schema.TypeString,
							Required: true,
						},
						"logout_url": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"error_url": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"certificate": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name_id_policy_format": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"emailAddress", "encrypted", "unspecified", "transient",
								"WindowsDomainQualifiedName", "X509SubjectName", "kerberos", "persistent", "entity",
							}, false),
						},
					},
				},
			},
			"idp_metadata_xml": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username_attribute": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"email_attribute": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"groups_attribute": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"groups_delim": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"custom_attributes": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"entity_issuer": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_signed_authn_req_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
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
		},
	}
}

func ResourceNutanixSamlIdpV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	input := &import1.SamlIdentityProvider{}
	if idpMetadataurl, ok := d.GetOk("idp_metadata_url"); ok {
		input.IdpMetadataUrl = utils.StringPtr(idpMetadataurl.(string))
	}
	if idpMetadata, ok := d.GetOk("idp_metadata"); ok {
		log.Printf("idp metadata: %v", idpMetadata)
		input.IdpMetadata = expandIdpMetadata(idpMetadata)
	}
	if idpMetaXML, ok := d.GetOk("idp_metadata_xml"); ok {
		input.IdpMetadataXml = utils.StringPtr(idpMetaXML.(string))
	}
	if name, ok := d.GetOk("name"); ok {
		input.Name = utils.StringPtr(name.(string))
	}
	if uName, ok := d.GetOk("username_attribute"); ok {
		input.UsernameAttribute = utils.StringPtr(uName.(string))
	}
	if emailAttr, ok := d.GetOk("email_attribute"); ok {
		input.EmailAttribute = utils.StringPtr(emailAttr.(string))
	}
	if grpAttr, ok := d.GetOk("groups_attribute"); ok {
		input.GroupsAttribute = utils.StringPtr(grpAttr.(string))
	}
	if grpDelim, ok := d.GetOk("groups_delim"); ok {
		input.GroupsDelim = utils.StringPtr(grpDelim.(string))
	}
	if customAttributes, ok := d.GetOk("custom_attributes"); ok {
		customAttributesList := customAttributes.([]interface{})
		customAttributesListStr := make([]string, len(customAttributesList))
		for i, v := range customAttributesList {
			customAttributesListStr[i] = v.(string)
		}
		input.CustomAttributes = customAttributesListStr
	}
	if entity, ok := d.GetOk("entity_issuer"); ok {
		input.EntityIssuer = utils.StringPtr(entity.(string))
	}
	if isSigned, ok := d.GetOk("is_signed_authn_req_enabled"); ok {
		input.IsSignedAuthnReqEnabled = utils.BoolPtr(isSigned.(bool))
	}

	resp, err := conn.SamlIdentityAPIInstance.CreateSamlIdentityProvider(input)
	if err != nil {
		return diag.Errorf("error while creating saml identity providers: %v", err)
	}

	getResp := resp.Data.GetValue().(import1.SamlIdentityProvider)

	d.SetId(*getResp.ExtId)
	return ResourceNutanixSamlIdpV2Read(ctx, d, meta)
}

func ResourceNutanixSamlIdpV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	resp, err := conn.SamlIdentityAPIInstance.GetSamlIdentityProviderById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching saml identity providers: %v", err)
	}

	getResp := resp.Data.GetValue().(import1.SamlIdentityProvider)

	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("idp_metadata_url", getResp.IdpMetadataUrl); err != nil {
		return diag.FromErr(err)
	}
	// if err := d.Set("idp_metadata_xml", getResp.IdpMetadataXml); err != nil {
	//	return diag.FromErr(err)
	//}
	if err := d.Set("idp_metadata", flattenIdpMetadata(getResp.IdpMetadata)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("username_attribute", getResp.UsernameAttribute); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("email_attribute", getResp.EmailAttribute); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("groups_attribute", getResp.GroupsAttribute); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("groups_delim", getResp.GroupsDelim); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("custom_attributes", getResp.CustomAttributes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("entity_issuer", getResp.EntityIssuer); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_signed_authn_req_enabled", getResp.IsSignedAuthnReqEnabled); err != nil {
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

func ResourceNutanixSamlIdpV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI
	updatedInput := import1.SamlIdentityProvider{}
	resp, err := conn.SamlIdentityAPIInstance.GetSamlIdentityProviderById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching saml identity providers: %v", err)
	}

	// get etag value from read response to pass in update request If-Match header, Required for update request
	etagValue := conn.SamlIdentityAPIInstance.ApiClient.GetEtag(resp)
	headers := make(map[string]interface{})
	headers["If-Match"] = utils.StringPtr(etagValue)

	updatedInput = resp.Data.GetValue().(import1.SamlIdentityProvider)

	if d.HasChange("name") {
		updatedInput.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("idp_metadata_url") {
		updatedInput.IdpMetadataUrl = utils.StringPtr(d.Get("idp_metadata_url").(string))
	}
	if d.HasChange("idp_metadata_xml") {
		updatedInput.IdpMetadataXml = utils.StringPtr(d.Get("idp_metadata_xml").(string))
	}
	if d.HasChange("idp_metadata") {
		updatedInput.IdpMetadata = expandIdpMetadata(d.Get("idp_metadata"))
	}
	if d.HasChange("username_attribute") {
		updatedInput.UsernameAttribute = utils.StringPtr(d.Get("username_attribute").(string))
	}
	if d.HasChange("email_attribute") {
		updatedInput.EmailAttribute = utils.StringPtr(d.Get("email_attribute").(string))
	}
	if d.HasChange("groups_attribute") {
		updatedInput.GroupsAttribute = utils.StringPtr(d.Get("groups_attribute").(string))
	}
	if d.HasChange("groups_delim") {
		updatedInput.GroupsDelim = utils.StringPtr(d.Get("groups_delim").(string))
	}
	if d.HasChange("custom_attributes") {
		customAttributes := d.Get("custom_attributes")
		customAttributesList := customAttributes.([]interface{})
		customAttributesListStr := make([]string, len(customAttributesList))
		for i, v := range customAttributesList {
			customAttributesListStr[i] = v.(string)
		}
		updatedInput.CustomAttributes = customAttributesListStr
	}
	if d.HasChange("entity_issuer") {
		updatedInput.EntityIssuer = utils.StringPtr(d.Get("entity_issuer").(string))
	}
	if d.HasChange("is_signed_authn_req_enabled") {
		updatedInput.IsSignedAuthnReqEnabled = utils.BoolPtr(d.Get("is_signed_authn_req_enabled").(bool))
	}

	updateResp, err := conn.SamlIdentityAPIInstance.UpdateSamlIdentityProviderById(utils.StringPtr(d.Id()), &updatedInput, headers)
	if err != nil {
		return diag.Errorf("error while updating saml identity providers: %v", err)
	}

	updateTaskResp := updateResp.Data.GetValue().(import1.SamlIdentityProvider)

	if updateTaskResp.ExtId != nil {
		log.Println("[DEBUG] Saml Identity provider updated successfully")
	}
	return nil
}

func ResourceNutanixSamlIdpV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	readResp, err := conn.SamlIdentityAPIInstance.GetSamlIdentityProviderById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching saml identity providers: %v", err)
	}
	// get etag value from read response to pass in update request If-Match header, Required for update request
	etagValue := conn.SamlIdentityAPIInstance.ApiClient.GetEtag(readResp)
	headers := make(map[string]interface{})
	headers["If-Match"] = utils.StringPtr(etagValue)

	resp, err := conn.SamlIdentityAPIInstance.DeleteSamlIdentityProviderById(utils.StringPtr(d.Id()), headers)
	if err != nil {
		return diag.Errorf("error while deleting saml idp : %v", err)
	}

	if resp == nil {
		log.Println("[DEBUG] Saml IDP deleted successfully.")
	}
	return nil
}

func expandIdpMetadata(pr interface{}) *import1.IdpMetadata {
	if pr != nil {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})
		idp := &import1.IdpMetadata{}

		if entityID, ok := val["entity_id"]; ok {
			idp.EntityId = utils.StringPtr(entityID.(string))
		}
		if loginURL, ok := val["login_url"]; ok {
			idp.LoginUrl = utils.StringPtr(loginURL.(string))
		}
		if logoutURL, ok := val["logout_url"]; ok {
			idp.LogoutUrl = utils.StringPtr(logoutURL.(string))
		}
		if errorURL, ok := val["error_url"]; ok {
			log.Printf("error url: %v", errorURL)
			if errorURL != "" {
				log.Printf("idp error url: %v", idp.ErrorUrl)
				idp.ErrorUrl = utils.StringPtr(errorURL.(string))
			} else {
				idp.ErrorUrl = nil
			}
		}
		if certi, ok := val["certificate"]; ok {
			idp.Certificate = utils.StringPtr(certi.(string))
		}
		if policyFormat, ok := val["name_id_policy_format"]; ok {
			const two, three, four, five, six, seven, eight, nine, ten = 2, 3, 4, 5, 6, 7, 8, 9, 10
			subMap := map[string]interface{}{
				"emailAddress":               two,
				"unspecified":                three,
				"X509SubjectName":            four,
				"WindowsDomainQualifiedName": five,
				"encrypted":                  six,
				"entity":                     seven,
				"kerberos":                   eight,
				"persistent":                 nine,
				"transient":                  ten,
			}
			pInt := subMap[policyFormat.(string)]
			p := import1.NameIdPolicyFormat(pInt.(int))
			idp.NameIdPolicyFormat = &p
		}
		log.Printf("idp: %v", idp)
		return idp
	}
	return nil
}
