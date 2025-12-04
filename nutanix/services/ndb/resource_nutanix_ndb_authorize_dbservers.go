package ndb

import (
	"context"
	"log"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixNDBAuthorizeDBServer() *schema.Resource {
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
	conn := meta.(*conns.Client).Era
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

	resp, err := conn.Service.AuthorizeDBServer(ctx, tmsID.(string), req)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.Status == utils.StringPtr("success") {
		uuid, er := uuid.GenerateUUID()

		if er != nil {
			return diag.Errorf("Error generating UUID for era clusters: %+v", err)
		}
		d.SetId(uuid)
	}
	log.Printf("NDB Authorize dbservers with %s id created successfully", d.Id())
	return nil
}

func resourceNutanixNDBAuthorizeDBServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixNDBAuthorizeDBServerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixNDBAuthorizeDBServerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

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

	deauthorizeDBs := make([]*string, 0)

	if dbserversID, ok := d.GetOk("dbservers_id"); ok {
		dbser := dbserversID.([]interface{})

		for _, v := range dbser {
			deauthorizeDBs = append(deauthorizeDBs, utils.StringPtr(v.(string)))
		}
	}

	_, err := conn.Service.DeAuthorizeDBServer(ctx, tmsID.(string), deauthorizeDBs)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("NDB Authorize dbservers with %s id deleted successfully", d.Id())
	d.SetId("")
	return nil
}
