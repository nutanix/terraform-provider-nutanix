package ndb

import (
	"context"
	"log"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixNDBTmsCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixNDBTmsClusterCreate,
		ReadContext:   resourceNutanixNDBTmsClusterRead,
		UpdateContext: resourceNutanixNDBTmsClusterUpdate,
		DeleteContext: resourceNutanixNDBTmsClusterDelete,
		Schema: map[string]*schema.Schema{
			"time_machine_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"nx_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sla_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "OTHER",
			},
			// computed
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schedule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"log_drive_status": {
				Type:     schema.TypeString,
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
			"log_drive_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceNutanixNDBTmsClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	req := &era.TmsClusterIntentInput{}

	tmsID := d.Get("time_machine_id")

	if nxcls, ok := d.GetOk("nx_cluster_id"); ok {
		req.NxClusterID = utils.StringPtr(nxcls.(string))
	}

	if slaid, ok := d.GetOk("sla_id"); ok {
		req.SLAID = utils.StringPtr(slaid.(string))
	}

	if clsType, ok := d.GetOk("type"); ok {
		req.Type = utils.StringPtr(clsType.(string))
	}

	_, err := conn.Service.CreateTimeMachineCluster(ctx, tmsID.(string), req)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, er := uuid.GenerateUUID()

	if er != nil {
		return diag.Errorf("Error generating UUID for era clusters: %+v", err)
	}
	d.SetId(uuid)
	log.Printf("NDB Time Machine Cluster with %s id is created successfully", d.Id())
	return resourceNutanixNDBTmsClusterRead(ctx, d, meta)
}

func resourceNutanixNDBTmsClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	tmsID := d.Get("time_machine_id")
	clsID := d.Get("nx_cluster_id")
	resp, err := conn.Service.ReadTimeMachineCluster(ctx, tmsID.(string), clsID.(string))
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("log_drive_id", resp.LogDrive); err != nil {
		return diag.Errorf("error occurred while setting log_drive_id for time machine cluster with id: %s : %s", d.Id(), err)
	}

	if err = d.Set("log_drive_status", resp.LogDriveStatus); err != nil {
		return diag.Errorf("error occurred while setting log_drive_status for time machine cluster with id: %s: %s", d.Id(), err)
	}

	if err = d.Set("type", resp.Type); err != nil {
		return diag.Errorf("error occurred while setting type for time machine cluster with id: %s: %s", d.Id(), err)
	}

	if err = d.Set("description", resp.Description); err != nil {
		return diag.Errorf("error occurred while setting description for time machine cluster with id: %s: %s", d.Id(), err)
	}

	if err = d.Set("status", resp.Status); err != nil {
		return diag.Errorf("error occurred while setting status for time machine cluster with id: %s: %s", d.Id(), err)
	}

	if err = d.Set("schedule_id", resp.ScheduleID); err != nil {
		return diag.Errorf("error occurred while setting schedule_id for time machine cluster with id: %s: %s", d.Id(), err)
	}

	if err = d.Set("owner_id", resp.OwnerID); err != nil {
		return diag.Errorf("error occurred while setting owner_id for time machine cluster with id: %s: %s", d.Id(), err)
	}

	if err = d.Set("date_created", resp.DateCreated); err != nil {
		return diag.Errorf("error occurred while setting date_created for time machine cluster with id: %s: %s", d.Id(), err)
	}

	if err = d.Set("source", resp.Source); err != nil {
		return diag.Errorf("error occurred while setting source for time machine cluster with id: %s: %s", d.Id(), err)
	}

	if err = d.Set("date_modified", resp.DateModified); err != nil {
		return diag.Errorf("error occurred while setting date_modified for time machine cluster with id: %s: %s", d.Id(), err)
	}
	if resp.SourceClusters != nil {
		sourceCls := make([]*string, 0)
		sourceCls = append(sourceCls, resp.SourceClusters...)

		d.Set("source_clusters", sourceCls)
	}
	return nil
}

func resourceNutanixNDBTmsClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	updateReq := &era.TmsClusterIntentInput{}

	tmsID := d.Get("time_machine_id")
	clsID := d.Get("nx_cluster_id")
	resp, err := conn.Service.ReadTimeMachineCluster(ctx, tmsID.(string), clsID.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if resp != nil {
		updateReq.Type = resp.Type
		updateReq.NxClusterID = resp.NxClusterID
	}

	if d.HasChange("sla_id") {
		updateReq.SLAID = utils.StringPtr(d.Get("sla_id").(string))
		updateReq.ResetSLAID = utils.BoolPtr(true)
	}

	if d.HasChange("nx_cluster_id") {
		updateReq.NxClusterID = utils.StringPtr(d.Get("nx_cluster_id").(string))
	}

	// update Call for time machine cluster

	_, er := conn.Service.UpdateTimeMachineCluster(ctx, tmsID.(string), clsID.(string), updateReq)
	if er != nil {
		return diag.FromErr(er)
	}
	log.Printf("NDB Time Machine Cluster with %s id is updated successfully", d.Id())
	return resourceNutanixNDBTmsClusterRead(ctx, d, meta)
}

func resourceNutanixNDBTmsClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	req := &era.DeleteTmsClusterInput{
		DeleteReplicatedSnapshots:         utils.BoolPtr(true),
		DeleteReplicatedProtectionDomains: utils.BoolPtr(true),
	}

	tmsID := d.Get("time_machine_id")
	clsID := d.Get("nx_cluster_id")

	resp, er := conn.Service.DeleteTimeMachineCluster(ctx, tmsID.(string), clsID.(string), req)
	if er != nil {
		return diag.FromErr(er)
	}

	if resp.Status == "" {
		d.SetId("")
		log.Printf("NDB Time Machine Cluster with %s id is deleted successfully", d.Id())
	}
	return nil
}
