package nutanix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func resourceNutanixNDBAuthorizeDBServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixNDBAuthorizeDBServerCreate,
		ReadContext:   resourceNutanixNDBAuthorizeDBServerRead,
		UpdateContext: resourceNutanixNDBAuthorizeDBServerUpdate,
		DeleteContext: resourceNutanixNDBAuthorizeDBServerDelete,
		Schema: map[string]*schema.Schema{
			"time_machine_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"time_machine_name"},
			},
			"time_machine_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"time_machine_id"},
			},
			"dbservers_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceNutanixNDBAuthorizeDBServerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).Era
	req := make([]*string, 0)

	tmsID, tok := d.GetOk("time_machine_id")
	tmsName, tnOk := d.GetOk("time_machine_name")

	if !tok && !tnOk {
		return diag.Errorf("Atleast one of time_machine_id or time_machine_name is required to perform clone")
	}

	if len(tmsName.(string)) > 0 {
		// call time machine API with value-type name
		res, er := conn.Service.GetTimeMachine(ctx, "", tmsName.(string))
		if er != nil {
			return diag.FromErr(er)
		}

		tmsID = *res.ID
	}

	if dbserversID, ok := d.GetOk("dbservers_id"); ok {
		dbser := dbserversID.([]interface{})

		for _, v := range dbser {
			req = append(req, utils.StringPtr(v.(string)))
		}
	}
	// call for Authorize API

	resp, err := conn.Service.AuthorizeDbServer(ctx, tmsID.(string), req)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.Status == utils.StringPtr("success") {
		d.SetId(tmsID.(string))
	}

	return nil
}

func resourceNutanixNDBAuthorizeDBServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixNDBAuthorizeDBServerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixNDBAuthorizeDBServerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
