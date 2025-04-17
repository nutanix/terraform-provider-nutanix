package iamv2

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	import1 "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/iam/v4/authn"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixUserKeysV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixUserKeyV2Create,
		ReadContext:   resourceNutanixUserKeyV2Read,
		UpdateContext: resourceNutanixUserKeyV2Update,
		DeleteContext: resourceNutanixUserKeyV2Delete,
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
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"key_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{"API_KEY", "OBJECT_KEY"}, false),
				Computed: true,
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
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{"PREDEFINED", "SERVICEDEFINED", "USERDEFINED"}, false),
				Computed: true,
			},
			"expiry_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{"REVOKED", "VALID", "EXPIRED"}, false),
				Computed: true,
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

	var userExtId *string
	if v, ok := d.GetOk("user_ext_id"); ok {
		userExtId = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("name"); ok {
		spec.Name = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("description"); ok {
		spec.Description = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("key_type"); ok {
		spec.KeyType = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("creation_type"); ok {
		spec.CreationType = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("expiry_time"); ok {
		spec.ExpiryTime = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("status"); ok {
		spec.Status = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("assigned_to"); ok {
		spec.AssignedTo = utils.StringPtr(v.(string))
	}
	
	resp, err := conn.UsersAPIInstance.CreateUserKey(userExtId, spec)
	if err != nil {
		return diag.Errorf("error while creating User Key: %v", err)
	}
	getResp := resp.Data.GetValue().(import1.Key)
	d.SetId(*getResp.ExtId)
	return resourceNutanixUserKeyV2Read(ctx, d, meta)
}

func resourceNutanixUserKeyV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
   // Get client connection
	conn := meta.(*conns.Client).IamAPI

	var userExtId *string
	if v, ok := d.GetOk("user_ext_id"); ok {
		userExtId = utils.StringPtr(v.(string))
	}
 
	resp, err := conn.UsersAPIInstance.GetUserKeyById(userExtId, utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching the user key: %v", err)
	}

	keyConfig := resp.Data.GetValue().(import1.Key)
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
	if err := d.Set("key_details", keyConfig.KeyDetails); err != nil {
		return diag.Errorf("error while setting key_details: %v", err)
	}
	d.SetId(*keyConfig.ExtId)
	return nil
}

func resourceNutanixUserKeyV2Update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixUserKeyV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {	
	conn := meta.(*conns.Client).IamAPI

	var userExtId *string
	if v, ok := d.GetOk("user_ext_id"); ok {
		userExtId = utils.StringPtr(v.(string))
	}

	resp, err := conn.UsersAPIInstance.GetUserKeyById(userExtId, utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching the user key: %v", err)
	}

	// Extract E-Tag Header
	args := make(map[string]interface{})
	etagValue := conn.UsersAPIInstance.ApiClient.GetEtag(resp)
	args["If-Match"] = utils.StringPtr(etagValue)
  
	_, del_err := conn.UsersAPIInstance.DeleteUserKeyById(userExtId, utils.StringPtr(d.Id()), args)
	if del_err != nil {
		return diag.Errorf("error while deleting the user key: %v", del_err)
	}
	d.SetId("")
	return nil
}