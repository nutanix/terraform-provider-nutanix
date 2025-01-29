package selfservice

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceNutanixCalmRunbook() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceNutanixCalmRunbookCreate,
		Schema: map[string]*schema.Schema{
			"app_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"length": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"offset": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"recovery_point_info_list": {
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
									"snapshot_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"kind": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"creation_time": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"recovery_point_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"expiration_time": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"location_agnostic_uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"service_references": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"config_spec_reference": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"api_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spec": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spec_version": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"total_matches": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNutanixCalmRunbookCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
