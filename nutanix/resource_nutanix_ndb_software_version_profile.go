package nutanix

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/terraform-providers/terraform-provider-nutanix/client/era"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func resourceNutanixNDBSoftwareVersionProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixNDBSoftwareVersionProfileCreate,
		ReadContext:   resourceNutanixNDBSoftwareVersionProfileRead,
		UpdateContext: resourceNutanixNDBSoftwareVersionProfileUpdate,
		DeleteContext: resourceNutanixNDBSoftwareVersionProfileDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"profile_id": {
				Type:     schema.TypeString,
				Required: true,
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
			"latest_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"latest_version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"topology": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_version": {
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
						"profile_id": {
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
				},
			},
			"nx_cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"assoc_databases": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cluster_availability": {
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
						"profile_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceNutanixNDBSoftwareVersionProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).Era

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
		return diag.Errorf("error waiting for software profile	version (%s) to create: %s", *resp.EntityID, errWaitTask)
	}
	d.SetId(*resp.EntityID)

	// spec for update profile request

	updateSpec := &era.ProfileRequest{}

	// getting name & description
	updateSpec.Name = req.Name
	updateSpec.Description = req.Description

	// now call the Update Profile API if publish params given
	if status, ok := d.GetOk("status"); ok {
		if status.(string) == "published" {
			updateSpec.Published = true
			updateSpec.Deprecated = false
		} else if status.(string) == "unpublished" {
			updateSpec.Published = false
			updateSpec.Deprecated = false
		} else {
			updateSpec.Published = false
			updateSpec.Deprecated = true
		}
	}

	//update for software profile version
	_, er := conn.Service.UpdateProfileVersion(ctx, updateSpec, profileID, d.Id())
	if er != nil {
		return diag.FromErr(er)
	}

	return nil
}

func resourceNutanixNDBSoftwareVersionProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixNDBSoftwareVersionProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).Era

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
		if status.(string) == "published" {
			req.Published = true
			req.Deprecated = false
		} else if status.(string) == "unpublished" {
			req.Published = false
			req.Deprecated = false
		} else {
			req.Published = false
			req.Deprecated = true
		}
	}

	//update for software profile version
	_, er := conn.Service.UpdateProfileVersion(ctx, req, profileID.(string), d.Id())
	if er != nil {
		return diag.FromErr(er)
	}

	return nil
}

func resourceNutanixNDBSoftwareVersionProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).Era
	profileID := d.Get("profile_id")
	resp, err := conn.Service.DeleteProfileVersion(ctx, profileID.(string), d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if resp == utils.StringPtr("Profile Successfully Deleted.") {
		d.SetId("")
	}
	return nil
}
