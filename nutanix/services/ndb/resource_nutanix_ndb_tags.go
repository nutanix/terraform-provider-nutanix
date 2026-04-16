package ndb

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixNDBTags() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixNDBTagsCreate,
		ReadContext:   resourceNutanixNDBTagsRead,
		UpdateContext: resourceNutanixNDBTagsUpdate,
		DeleteContext: resourceNutanixNDBTagsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"required": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"entity_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"DATABASE", "TIME_MACHINE",
					"CLONE", "DATABASE_SERVER",
				}, false),
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"ENABLED", "DEPRECATED"}, false),
			},
			//computed values

			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"values": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNutanixNDBTagsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	req := &era.CreateTagsInput{}
	tagName := ""
	entityType := ""

	if name, ok := d.GetOk("name"); ok {
		req.Name = utils.StringPtr(name.(string))
		tagName = name.(string)
	}

	if desc, ok := d.GetOk("description"); ok {
		req.Description = utils.StringPtr(desc.(string))
	}

	if require, ok := d.GetOk("required"); ok {
		req.Required = utils.BoolPtr(require.(bool))
	}

	if entity, ok := d.GetOk("entity_type"); ok {
		req.EntityType = utils.StringPtr(entity.(string))
		entityType = entity.(string)
	}

	_, err := conn.Service.CreateTags(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	uniqueID := ""
	// fetch all the tags
	tagsListResp, er := conn.Service.ListTags(ctx, entityType)
	if er != nil {
		return diag.FromErr(er)
	}

	for _, v := range *tagsListResp {
		if tagName == utils.StringValue(v.Name) && entityType == utils.StringValue(v.EntityType) {
			uniqueID = *v.ID
		}
	}
	d.SetId(uniqueID)
	log.Printf("NDB Tag with %s id is created successfully", uniqueID)
	return resourceNutanixNDBTagsRead(ctx, d, meta)
}

func resourceNutanixNDBTagsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	// check if d.Id() is nil
	if d.Id() == "" {
		return diag.Errorf("tag id is required for read operation")
	}
	resp, err := conn.Service.ReadTags(ctx, d.Id(), "")
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("name", resp.Name); err != nil {
		return diag.Errorf("error setting name for tag %s: %s", d.Id(), err)
	}

	if err = d.Set("description", resp.Description); err != nil {
		return diag.Errorf("error setting description for tag %s: %s", d.Id(), err)
	}
	if err = d.Set("date_created", resp.DateCreated); err != nil {
		return diag.Errorf("error setting date created for tag %s: %s", d.Id(), err)
	}

	if err = d.Set("date_modified", resp.DateModified); err != nil {
		return diag.Errorf("error setting date modified for tag %s: %s", d.Id(), err)
	}
	if err = d.Set("owner", resp.Owner); err != nil {
		return diag.Errorf("error setting owner id for tag %s: %s", d.Id(), err)
	}

	if err = d.Set("required", resp.Required); err != nil {
		return diag.Errorf("error setting required for tag %s: %s", d.Id(), err)
	}
	if err = d.Set("status", resp.Status); err != nil {
		return diag.Errorf("error setting status for tag %s: %s", d.Id(), err)
	}

	if err = d.Set("entity_type", resp.EntityType); err != nil {
		return diag.Errorf("error setting entity type for tag %s: %s", d.Id(), err)
	}
	if err = d.Set("values", resp.Values); err != nil {
		return diag.Errorf("error setting values for tag %s: %s", d.Id(), err)
	}
	return nil
}

func resourceNutanixNDBTagsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	updateReq := &era.GetTagsResponse{}

	// read the tag

	resp, err := conn.Service.ReadTags(ctx, d.Id(), "")
	if err != nil {
		return diag.FromErr(err)
	}

	if resp != nil {
		updateReq.Name = resp.Name
		updateReq.Description = resp.Description
		updateReq.DateCreated = resp.DateCreated
		updateReq.DateModified = resp.DateModified
		updateReq.Owner = resp.Owner
		updateReq.Required = resp.Required
		updateReq.Status = resp.Status
		updateReq.EntityType = resp.EntityType
		updateReq.Values = resp.Values
	}

	if d.HasChange("name") {
		updateReq.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("description") {
		updateReq.Description = utils.StringPtr(d.Get("description").(string))
	}
	if d.HasChange("required") {
		updateReq.Required = utils.BoolPtr(d.Get("required").(bool))
	}
	if d.HasChange("status") {
		updateReq.Status = utils.StringPtr(d.Get("status").(string))
	}

	updateResp, er := conn.Service.UpdateTags(ctx, updateReq, d.Id())
	if er != nil {
		return diag.FromErr(er)
	}
	log.Printf("NDB Tag with %s id updated successfully", *updateResp.ID)
	return resourceNutanixNDBTagsRead(ctx, d, meta)
}

func resourceNutanixNDBTagsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	resp, err := conn.Service.DeleteTags(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if resp == utils.StringPtr("Tag Successfully Deleted.") {
		log.Printf("NDB Tag with %s id is deleted successfully", d.Id())
		d.SetId("")
	}
	return nil
}
