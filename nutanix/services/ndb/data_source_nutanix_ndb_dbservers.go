package ndb

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DataSourceNutanixNDBDBServers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixNDBDBServersRead,
		Schema: map[string]*schema.Schema{
			"dbservers": {
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
						"date_created": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"date_modified": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"properties": {
							Type:        schema.TypeList,
							Description: "List of all the properties",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "",
									},

									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "",
									},
								},
							},
						},
						"tags": dataSourceEraDBInstanceTags(),
						"vm_cluster_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_addresses": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"fqdns": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mac_addresses": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"era_drive_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"era_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vm_timezone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vm_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"secure_info": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"info": dataSourceEraDatabaseInfo(),
									"deregister_info": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"message": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"operations": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"os_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"os_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"distribution": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"network_info": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"vlan_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"vlan_uuid": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"vlan_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"era_configured": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"gateway": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"subnet_mask": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"hostname": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"device_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"mac_address": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"flags": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"mtu": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"ip_addresses": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"default_gateway_device": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"access_info": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"access_type": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"destination_subnet": {
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
							},
						},
						"clustered": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_server_driven": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"protection_domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"query_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"database_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dbserver_invalid_ea_state": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"working_directory": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"valid_diagnostic_bundle_state": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"windows_db_server": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"associated_time_machine_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"access_key_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNutanixNDBDBServersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	resp, err := conn.Service.ListDBServerVM(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	if e := d.Set("dbservers", flattenDBServerVMResponse(resp)); e != nil {
		return diag.FromErr(e)
	}

	uuid, er := uuid.GenerateUUID()

	if er != nil {
		return diag.Errorf("Error generating UUID for era dbservers: %+v", er)
	}
	d.SetId(uuid)

	return nil
}

func flattenDBServerVMResponse(pr *era.ListDBServerVMResponse) []interface{} {
	if pr != nil {
		lst := []interface{}{}

		for _, v := range *pr {
			vms := map[string]interface{}{}

			vms["id"] = v.ID
			vms["name"] = v.Name
			vms["description"] = v.Description
			vms["date_created"] = v.DateCreated
			vms["date_modified"] = v.DateModified
			vms["access_level"] = v.AccessLevel
			if v.Properties != nil {
				props := []interface{}{}
				for _, prop := range v.Properties {
					props = append(props, map[string]interface{}{
						"name":  prop.Name,
						"value": prop.Value,
					})
				}
				vms["properties"] = props
			}
			vms["tags"] = flattenDBTags(v.Tags)
			vms["vm_cluster_uuid"] = v.VMClusterUUID
			vms["ip_addresses"] = utils.StringValueSlice(v.IPAddresses)
			vms["fqdns"] = v.Fqdns
			vms["mac_addresses"] = utils.StringValueSlice(v.MacAddresses)
			vms["type"] = v.Type
			vms["status"] = v.Status
			vms["client_id"] = v.ClientID
			vms["era_drive_id"] = v.EraDriveID
			vms["era_version"] = v.EraVersion
			vms["vm_timezone"] = v.VMTimeZone
			vms["vm_info"] = flattenDBServerVMInfo(v.VMInfo)
			vms["clustered"] = v.Clustered
			vms["is_server_driven"] = v.IsServerDriven
			vms["protection_domain_id"] = v.ProtectionDomainID
			vms["query_count"] = v.QueryCount
			vms["database_type"] = v.DatabaseType
			vms["dbserver_invalid_ea_state"] = v.DbserverInValidEaState
			vms["working_directory"] = v.WorkingDirectory
			vms["valid_diagnostic_bundle_state"] = v.ValidDiagnosticBundleState
			vms["windows_db_server"] = v.WindowsDBServer
			vms["associated_time_machine_ids"] = utils.StringValueSlice(v.AssociatedTimeMachineIds)
			vms["access_key_id"] = v.AccessKeyID

			lst = append(lst, vms)
		}
		return lst
	}
	return nil
}
