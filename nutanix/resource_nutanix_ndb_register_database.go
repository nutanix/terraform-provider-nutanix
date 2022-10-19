package nutanix

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	era "github.com/terraform-providers/terraform-provider-nutanix/client/era"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func resourceNutanixNDBRegisterDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixNDBRegisterDatabaseCreate,
		ReadContext:   resourceNutanixNDBRegisterDatabaseRead,
		UpdateContext: resourceNutanixNDBRegisterDatabaseUpdate,
		DeleteContext: resourceNutanixNDBRegisterDatabaseDelete,
		Schema: map[string]*schema.Schema{
			"nx_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"database_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"database_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"clustered": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  false,
			},
			"forced_install": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  true,
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "DEFAULT",
			},
			"vm_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vm_username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vm_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vm_sshkey": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vm_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"reset_description_in_nx_cluster": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  false,
			},
			"auto_tune_staging_drive": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  true,
			},
			"working_directory": {
				Type:     schema.TypeString,
				Required: true,
			},
			"time_machine": timeMachineInfoSchema(),
			"tags":         dataSourceEraDBInstanceTags(),
			"postgress_info": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: schema.Resource{
					Schema: map[string]*schema.Schema{
						"listener_port": {
							Type:     schema.TypeString,
							Required: true,
						},
						"db_user": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"switch_log": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"allow_multiple_databases": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"backup_policy": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "prefer_secondary",
						},
						"vm_ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"postgres_software_home": {
							Type:     schema.TypeString,
							Required: true,
						},
						"db_password": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"db_name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceNutanixNDBRegisterDatabaseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).Era

	log.Println("Creating the request!!!")
	req, err := buildReisterDbRequest(d)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, er := conn.Service.RegisterDatabase(ctx, req)
	if er != nil {
		return diag.FromErr(er)
	}
	log.Println(resp)
	d.SetId("ssss")
	return nil
}

func resourceNutanixNDBRegisterDatabaseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
func resourceNutanixNDBRegisterDatabaseUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
func resourceNutanixNDBRegisterDatabaseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func buildReisterDbRequest(d *schema.ResourceData) (*era.RegisterDBInputRequest, error) {
	return &era.RegisterDBInputRequest{
		NxClusterID:                 utils.StringPtr(d.Get("nx_cluster_id").(string)),
		DatabaseType:                utils.StringPtr(d.Get("database_type").(string)),
		DatabaseName:                utils.StringPtr(d.Get("database_name").(string)),
		Description:                 utils.StringPtr(d.Get("description").(string)),
		Clustered:                   d.Get("clustered").(bool),
		ForcedInstall:               d.Get("forced_install").(bool),
		Category:                    utils.StringPtr(d.Get("category").(string)),
		VMIP:                        utils.StringPtr(d.Get("vm_ip").(string)),
		VMUsername:                  utils.StringPtr(d.Get("vm_username").(string)),
		VMPassword:                  utils.StringPtr(d.Get("vm_password").(string)),
		VMSshkey:                    utils.StringPtr(d.Get("vm_sshkey").(string)),
		VMDescription:               utils.StringPtr(d.Get("vm_description").(string)),
		ResetDescriptionInNxCluster: d.Get("reset_description_in_nx_cluster").(bool),
		AutoTuneStagingDrive:        d.Get("auto_tune_staging_drive").(bool),
		WorkingDirectory:            utils.StringPtr(d.Get("working_directory").(string)),
		TimeMachineInfo:             buildTimeMachineFromResourceData(d.Get("time_machine").(*schema.Set)),
		Actionarguments:             expandRegisterDbActionArguments(d),
	}, nil
}

func expandRegisterDbActionArguments(d *schema.ResourceData) []*era.Actionarguments {
	args := []*era.Actionarguments{}
	if post, ok := d.GetOk("postgresql_info"); ok {
		brr := post.([]interface{})

		for _, arg := range brr {
			val := arg.(map[string]interface{})
			var values interface{}
			if plist, pok := val["listener_port"]; pok && len(plist.(string)) > 0 {
				values = plist

				args = append(args, &era.Actionarguments{
					Name:  "listener_port",
					Value: values,
				})
			}
			if plist, pok := val["db_user"]; pok && len(plist.(string)) > 0 {
				values = plist

				args = append(args, &era.Actionarguments{
					Name:  "db_user",
					Value: values,
				})
			}
			if plist, pok := val["switch_log"]; pok && plist.(bool) {
				values = plist

				args = append(args, &era.Actionarguments{
					Name:  "switch_log",
					Value: values,
				})
			}
			if plist, pok := val["allow_multiple_databases"]; pok && plist.(bool) {
				values = plist

				args = append(args, &era.Actionarguments{
					Name:  "allow_multiple_databases",
					Value: values,
				})
			}
			if plist, pok := val["backup_policy"]; pok && len(plist.(string)) > 0 {
				values = plist

				args = append(args, &era.Actionarguments{
					Name:  "backup_policy",
					Value: values,
				})
			}
			if plist, pok := val["vm_ip"]; pok && len(plist.(string)) > 0 {
				values = plist

				args = append(args, &era.Actionarguments{
					Name:  "vmIp",
					Value: values,
				})
			}
			if plist, pok := val["postgres_software_home"]; pok && len(plist.(string)) > 0 {
				values = plist

				args = append(args, &era.Actionarguments{
					Name:  "postgres_software_home",
					Value: values,
				})
			}
			if plist, pok := val["db_password"]; pok && len(plist.(string)) > 0 {
				values = plist

				args = append(args, &era.Actionarguments{
					Name:  "db_password",
					Value: values,
				})
			}
			if plist, pok := val["db_name"]; pok && len(plist.(string)) > 0 {
				values = plist

				args = append(args, &era.Actionarguments{
					Name:  "db_name",
					Value: values,
				})
			}
		}
	}

	resp := buildActionArgumentsFromResourceData(d.Get("actionarguments").(*schema.Set), args)
	return resp
}
