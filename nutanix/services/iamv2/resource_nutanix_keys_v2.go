package iamv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	import1 "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/iam/v4/authn"
)

func ResourceNutanixUserKeysV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixUserKeysV2Create,
		ReadContext:   resourceNutanixUserKeysV2Read,
		UpdateContext: resourceNutanixUserKeysV2Update,
		DeleteContext: resourceNutanixUserKeysV2Delete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"keyType": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{"API_KEY", "OBJECT_KEY"}, false),
			},
			"creationType": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{"PREDEFINED", "SERVICEDEFINED", "USERDEFINED"}, false),
			},
			"expiryTime": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"assignedTo": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceNutanixUserKeysV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).IamAPI
	spec := &import1.Key{}
	resp, err := conn.UsersAPIInstance.CreateUserKey(spec)
	if err != nil {
		return diag.Errorf("error while creating User Key: %v", err)
	}
	getResp := resp.Data.GetValue().(import1.CreateKeyApiResponse)

	d.SetId(*getResp.ExtId)
	return resourceNutanixUserKeysV2Read(ctx, d, meta)
}

func resourceNutanixUserKeysV2Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// conn := meta.(*conns.Client).IamAPI
	return nil
}

func resourceNutanixUserKeysV2Update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// conn := meta.(*conns.Client).IamAPI
	return nil
}

func resourceNutanixUserKeysV2Delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {	
	// conn := meta.(*conns.Client).IamAPI
	return nil
}
		