package iamv2

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import1 "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/iam/v4/authn"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixUserKeyV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixUserKeyV2Create,
		ReadContext:   resourceNutanixUserKeyV2Read,
		UpdateContext: resourceNutanixUserKeyV2Update,
		DeleteContext: resourceNutanixUserKeyV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				const expectedPartsCount = 2
				parts := strings.Split(d.Id(), "/")
				if len(parts) != expectedPartsCount {
					return nil, fmt.Errorf("invalid import uuid (%q), expected user_ext_id/user_key_ext_id", d.Id())
				}
				d.Set("user_ext_id", parts[0])
				d.SetId(parts[1])
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"user_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": SchemaForLinks(),
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"key_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"API_KEY", "OBJECT_KEY"}, false),
			},
			"created_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"last_updated_by": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"creation_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"PREDEFINED", "SERVICEDEFINED", "USERDEFINED"}, false),
				Computed:     true,
			},
			"expiry_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"REVOKED", "VALID", "EXPIRED"}, false),
				Computed:     true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"last_updated_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"last_used_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"assigned_to": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"key_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_key_details": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"api_key": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"object_key_details": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"secret_key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"access_key": {
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
	}
}

func resourceNutanixUserKeyV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI
	spec := &import1.Key{}

	var creationType = map[string]import1.CreationType{
		"PREDEFINED":     import1.CREATIONTYPE_PREDEFINED,
		"USERDEFINED":    import1.CREATIONTYPE_USERDEFINED,
		"SERVICEDEFINED": import1.CREATIONTYPE_SERVICEDEFINED,
	}

	var KeyType = map[string]import1.KeyKind{
		"API_KEY":    import1.KEYKIND_API_KEY,
		"OBJECT_KEY": import1.KEYKIND_OBJECT_KEY,
	}

	var KeyStatus = map[string]import1.KeyStatus{
		"VALID":   import1.KEYSTATUS_VALID,
		"REVOKED": import1.KEYSTATUS_REVOKED,
		"EXPIRED": import1.KEYSTATUS_EXPIRED,
	}

	var userExtID *string
	if v, ok := d.GetOk("user_ext_id"); ok {
		userExtID = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("name"); ok {
		spec.Name = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("description"); ok {
		spec.Description = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("key_type"); ok {
		if strValue, isString := v.(string); isString {
			if enumValue, exists := KeyType[strValue]; exists {
				spec.KeyType = &enumValue
			}
		}
	}
	if v, ok := d.GetOk("creation_type"); ok {
		if strValue, isString := v.(string); isString {
			if enumValue, exists := creationType[strValue]; exists {
				spec.CreationType = &enumValue
			}
		}
	}
	if v, ok := d.GetOk("expiry_time"); ok {
		expiryStr := v.(string)
		expiryTime, err := time.Parse(time.RFC3339, expiryStr)
		if err == nil {
			spec.ExpiryTime = &expiryTime
		}
	}
	if v, ok := d.GetOk("status"); ok {
		if strValue, isString := v.(string); isString {
			if enumValue, exists := KeyStatus[strValue]; exists {
				spec.Status = &enumValue
			}
		}
	}
	if v, ok := d.GetOk("assigned_to"); ok {
		spec.AssignedTo = utils.StringPtr(v.(string))
	}

	resp, err := conn.UsersAPIInstance.CreateUserKey(userExtID, spec)
	if err != nil {
		return diag.Errorf("error while creating User Key: %v", err)
	}
	getResp := resp.Data.GetValue().(import1.Key)
	d.SetId(utils.StringValue(getResp.ExtId))
	return resourceNutanixUserKeyV2Read(ctx, d, meta)
}

func resourceNutanixUserKeyV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection
	conn := meta.(*conns.Client).IamAPI

	var userExtID *string
	if v, ok := d.GetOk("user_ext_id"); ok {
		userExtID = utils.StringPtr(v.(string))
	}

	resp, err := conn.UsersAPIInstance.GetUserKeyById(userExtID, utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching the user key: %v", err)
	}

	keyConfig := resp.Data.GetValue().(import1.Key)

	aJSON, _ := json.MarshalIndent(keyConfig, "", "  ")
	log.Printf("[DEBUG] Retrieved User Key: %s", aJSON)

	if err := d.Set("tenant_id", keyConfig.TenantId); err != nil {
		return diag.Errorf("error while setting tenant_id: %v", err)
	}
	if err := d.Set("ext_id", keyConfig.ExtId); err != nil {
		return diag.Errorf("error while setting ext_id: %v", err)
	}
	if err := d.Set("links", flattenLinks(keyConfig.Links)); err != nil {
		return diag.Errorf("error while setting links: %v", err)
	}
	if err := d.Set("name", keyConfig.Name); err != nil {
		return diag.Errorf("error while setting name: %v", err)
	}
	if err := d.Set("description", keyConfig.Description); err != nil {
		return diag.Errorf("error while setting description: %v", err)
	}
	if err := d.Set("key_type", keyConfig.KeyType.GetName()); err != nil {
		return diag.Errorf("error while setting key_type: %v", err)
	}
	if err := d.Set("created_time", flattenTime(keyConfig.CreatedTime)); err != nil {
		return diag.Errorf("error while setting created_time: %v", err)
	}
	if err := d.Set("last_updated_by", keyConfig.LastUpdatedBy); err != nil {
		return diag.Errorf("error while setting last_updated_by: %v", err)
	}
	if err := d.Set("creation_type", keyConfig.CreationType.GetName()); err != nil {
		return diag.Errorf("error while setting creation_type: %v", err)
	}
	if err := d.Set("expiry_time", flattenTime(keyConfig.ExpiryTime)); err != nil {
		return diag.Errorf("error while setting expiry_time: %v", err)
	}
	if err := d.Set("status", keyConfig.Status.GetName()); err != nil {
		return diag.Errorf("error while setting status: %v", err)
	}
	if err := d.Set("created_by", keyConfig.CreatedBy); err != nil {
		return diag.Errorf("error while setting created_by: %v", err)
	}
	if err := d.Set("last_updated_time", flattenTime(keyConfig.LastUpdatedTime)); err != nil {
		return diag.Errorf("error while setting last_updated_time: %v", err)
	}
	if err := d.Set("assigned_to", keyConfig.AssignedTo); err != nil {
		return diag.Errorf("error while setting assigned_to: %v", err)
	}
	if err := d.Set("last_used_time", flattenTime(keyConfig.LastUsedTime)); err != nil {
		return diag.Errorf("error while setting last_used_time: %v", err)
	}
	if err := d.Set("key_details", flattenKeyDetails(keyConfig.KeyDetails)); err != nil {
		return diag.Errorf("error while setting key_details: %v", err)
	}
	d.SetId(utils.StringValue(keyConfig.ExtId))
	return nil
}

func flattenKeyDetails(oneOfKeyKeyDetails *import1.OneOfKeyKeyDetails) interface{} {
	if oneOfKeyKeyDetails == nil {
		return nil
	}

	keyDetailsMap := make(map[string]interface{})

	// Determine which type is set using discriminator
	switch v := oneOfKeyKeyDetails.GetValue().(type) {
	case import1.ApiKeyDetails:
		apiKeyList := []map[string]interface{}{
			{
				"api_key": v.ApiKey,
			},
		}
		keyDetailsMap["api_key_details"] = apiKeyList

	case import1.ObjectKeyDetails:
		objectKeyList := []map[string]interface{}{
			{
				"secret_key": v.SecretKey,
				"access_key": v.AccessKey,
			},
		}
		keyDetailsMap["object_key_details"] = objectKeyList

	default:
		// If discriminator not set or unknown type
		return nil
	}

	return []map[string]interface{}{keyDetailsMap}
}

func resourceNutanixUserKeyV2Update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceNutanixUserKeyV2Create(ctx, d, m)
}

func resourceNutanixUserKeyV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI

	var userExtID *string
	if v, ok := d.GetOk("user_ext_id"); ok {
		userExtID = utils.StringPtr(v.(string))
	}

	resp, err := conn.UsersAPIInstance.GetUserKeyById(userExtID, utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching the user key: %v", err)
	}

	// Extract E-Tag Header
	args := make(map[string]interface{})
	etagValue := conn.UsersAPIInstance.ApiClient.GetEtag(resp)
	args["If-Match"] = utils.StringPtr(etagValue)

	_, delErr := conn.UsersAPIInstance.DeleteUserKeyById(userExtID, utils.StringPtr(d.Id()), args)
	if delErr != nil {
		return diag.Errorf("error while deleting the user key: %v", delErr)
	}
	d.SetId("")
	return nil
}
