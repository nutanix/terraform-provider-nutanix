package ndb

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

var SoftwareVersionProfileTimeout = 15 * time.Minute

func ResourceNutanixNDBSoftwareVersionProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixNDBSoftwareVersionProfileCreate,
		ReadContext:   resourceNutanixNDBSoftwareVersionProfileRead,
		UpdateContext: resourceNutanixNDBSoftwareVersionProfileUpdate,
		DeleteContext: resourceNutanixNDBSoftwareVersionProfileDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(SoftwareVersionProfileTimeout),
		},
		Schema: map[string]*schema.Schema{
			"profile_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"engine_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"deprecated", "published", "unpublished"}, false),
			},
			"postgres_database": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_dbserver_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"os_notes": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"db_software_notes": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"available_cluster_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			// computed arguments
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"topology": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"system_profile": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"published": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"deprecated": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"properties": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secure": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"properties_map": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"version_cluster_association": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nx_cluster_id": {
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
						"owner_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"profile_version_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"properties": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"secure": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"optimized_for_provisioning": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceNutanixNDBSoftwareVersionProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	req := &era.ProfileRequest{}
	profileID := ""
	// pre-filled requests

	req.DBVersion = utils.StringPtr("ALL")
	req.SystemProfile = false
	req.Type = utils.StringPtr("Software")

	if pID, ok := d.GetOk("profile_id"); ok {
		profileID = pID.(string)
	}

	if name, ok := d.GetOk("name"); ok {
		req.Name = utils.StringPtr(name.(string))
	}

	if desc, ok := d.GetOk("description"); ok {
		req.Description = utils.StringPtr(desc.(string))
	}

	if engType, ok := d.GetOk("engine_type"); ok {
		req.EngineType = utils.StringPtr(engType.(string))
	}

	if ps, ok := d.GetOk("postgres_database"); ok {
		req.Properties = expandSoftwareProfileProp(ps.([]interface{}))
	}

	if ac, ok1 := d.GetOk("available_cluster_ids"); ok1 {
		st := ac.([]interface{})
		sublist := make([]*string, len(st))

		for a := range st {
			sublist[a] = utils.StringPtr(st[a].(string))
		}
		req.AvailableClusterIds = sublist
	}

	// API to create software versions

	resp, err := conn.Service.CreateSoftwareProfileVersion(ctx, profileID, req)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get Operation ID from response of SoftwareProfileVersion  and poll for the operation to get completed.
	opID := resp.OperationID
	if opID == utils.StringPtr("") {
		return diag.Errorf("error: operation ID is an empty string")
	}
	opReq := era.GetOperationRequest{
		OperationID: utils.StringValue(opID),
	}

	log.Printf("polling for operation with id: %s\n", *opID)

	// Poll for operation here - Operation GET Call
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED", "FAILED"},
		Refresh: eraRefresh(ctx, conn, opReq),
		Timeout: d.Timeout(schema.TimeoutCreate),
		Delay:   eraDelay,
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for software profile version (%s) to create: %s", *resp.EntityID, errWaitTask)
	}
	d.SetId(*resp.EntityID)

	// spec for update profile request

	updateSpec := &era.ProfileRequest{}

	// getting name & description
	updateSpec.Name = req.Name
	updateSpec.Description = req.Description

	// now call the Update Profile API if publish params given
	if status, ok := d.GetOk("status"); ok {
		statusValue := status.(string)

		switch {
		case statusValue == "published":
			updateSpec.Published = true
			updateSpec.Deprecated = false
		case statusValue == "unpublished":
			updateSpec.Published = false
			updateSpec.Deprecated = false
		default:
			updateSpec.Published = false
			updateSpec.Deprecated = true
		}
	}

	//update for software profile version
	_, er := conn.Service.UpdateProfileVersion(ctx, updateSpec, profileID, d.Id())
	if er != nil {
		return diag.FromErr(er)
	}
	log.Printf("NDB Software Version Profile with %s id is created successfully", d.Id())
	return resourceNutanixNDBSoftwareVersionProfileRead(ctx, d, meta)
}

func resourceNutanixNDBSoftwareVersionProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	// Get Profile Version API
	profileVersionID := d.Get("profile_id")
	// check if d.Id() is nil
	if d.Id() == "" {
		return diag.Errorf("profile version id is required for read operation")
	}
	resp, err := conn.Service.GetSoftwareProfileVersion(ctx, profileVersionID.(string), d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", resp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", resp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("engine_type", resp.Enginetype); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("owner", resp.Owner); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("db_version", resp.Dbversion); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("topology", resp.Topology); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("system_profile", resp.Systemprofile); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("version", resp.Version); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("published", resp.Published); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("deprecated", resp.Deprecated); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("properties", flattenProperties(resp.Properties)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("properties_map", utils.ConvertMapString(resp.Propertiesmap)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("version_cluster_association", flattenClusterAssociation(resp.VersionClusterAssociation)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNutanixNDBSoftwareVersionProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	req := &era.ProfileRequest{}

	profileID := d.Get("profile_id")
	// get the software profile version

	oldResp, err := conn.Service.GetSoftwareProfileVersion(ctx, profileID.(string), d.Id())
	if err != nil {
		diag.FromErr(err)
	}

	if oldResp != nil {
		req.Name = oldResp.Name
		req.Description = oldResp.Description
	}

	if d.HasChange("name") {
		req.Name = utils.StringPtr(d.Get("name").(string))
	}

	if d.HasChange("description") {
		req.Description = utils.StringPtr(d.Get("description").(string))
	}

	if status, ok := d.GetOk("status"); ok {
		statusValue := status.(string)
		switch {
		case statusValue == "published":
			req.Published = true
			req.Deprecated = false
		case statusValue == "unpublished":
			req.Published = false
			req.Deprecated = false
		default:
			req.Published = false
			req.Deprecated = true
		}
	}

	//update for software profile version
	_, er := conn.Service.UpdateProfileVersion(ctx, req, profileID.(string), d.Id())
	if er != nil {
		return diag.FromErr(er)
	}
	log.Printf("NDB Software Version Profile with %s id is updated successfully", d.Id())
	return resourceNutanixNDBSoftwareVersionProfileRead(ctx, d, meta)
}

func resourceNutanixNDBSoftwareVersionProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era
	profileID := d.Get("profile_id")
	resp, err := conn.Service.DeleteProfileVersion(ctx, profileID.(string), d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if resp == utils.StringPtr("Profile Successfully Deleted.") {
		log.Printf("NDB Software Version Profile with %s id is deleted successfully", d.Id())
		d.SetId("")
	}
	return nil
}
