package iamv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import3 "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/iam/v4/authn"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixUserKeyV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixUserKeyV2Create,
		Schema: map[string]*schema.Schema{
			"user_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": SchemaForLinks(),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"key_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expiry_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_updated_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"assigned_to": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_used_time": {
				Type:     schema.TypeString,
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

func dataSourceNutanixUserKeyV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection
	conn := meta.(*conns.Client).IamAPI

	var userExtID *string
	if v, ok := d.GetOk("user_ext_id"); ok {
		userExtID = utils.StringPtr(v.(string))
	}

	var ExtID *string
	if v, ok := d.GetOk("ext_id"); ok {
		ExtID = utils.StringPtr(v.(string))
	}

	resp, err := conn.UsersAPIInstance.GetUserKeyById(userExtID, ExtID)
	if err != nil {
		return diag.Errorf("error while fetching the user key: %v", err)
	}

	keyConfig := resp.Data.GetValue().(import3.Key)
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
