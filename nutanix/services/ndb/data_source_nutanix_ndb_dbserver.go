package ndb

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DataSourceNutanixNDBDBServer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixNDBDBServerRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vm_cluster_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vm_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dbserver_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"nx_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// computed
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
	}
}

func dataSourceNutanixNDBDBServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	filterReq := &era.DBServerFilterRequest{}

	if name, ok := d.GetOk("name"); ok {
		filterReq.Name = utils.StringPtr(name.(string))
	}
	if id, ok := d.GetOk("id"); ok {
		filterReq.ID = utils.StringPtr(id.(string))
	}
	if ip, ok := d.GetOk("ip"); ok {
		filterReq.IP = utils.StringPtr(ip.(string))
	}
	if vmClsName, ok := d.GetOk("vm_cluster_name"); ok {
		filterReq.VMClusterName = utils.StringPtr(vmClsName.(string))
	}
	if vmClsid, ok := d.GetOk("vm_cluster_id"); ok {
		filterReq.VMClusterID = utils.StringPtr(vmClsid.(string))
	}
	if nxclsID, ok := d.GetOk("nx_cluster_id"); ok {
		filterReq.NxClusterID = utils.StringPtr(nxclsID.(string))
	}
	if dbserver, ok := d.GetOk("dbserver_cluster_id"); ok {
		filterReq.DBServerClusterID = utils.StringPtr(dbserver.(string))
	}

	resp, err := conn.Service.GetDBServerVM(ctx, filterReq)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", resp.ID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", resp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", resp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("date_created", resp.DateCreated); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("date_modified", resp.DateModified); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("access_level", resp.AccessLevel); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("vm_cluster_uuid", resp.VMClusterUUID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ip_addresses", utils.StringValueSlice(resp.IPAddresses)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("fqdns", resp.Fqdns); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("mac_addresses", utils.StringValueSlice(resp.MacAddresses)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("type", resp.Type); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("status", resp.Status); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("client_id", resp.ClientID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("era_drive_id", resp.EraDriveID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("era_version", resp.EraVersion); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vm_timezone", resp.VMTimeZone); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("clustered", resp.Clustered); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_server_driven", resp.IsServerDriven); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("protection_domain_id", resp.ProtectionDomainID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("query_count", resp.QueryCount); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("database_type", resp.DatabaseType); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("dbserver_invalid_ea_state", resp.DbserverInValidEaState); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("working_directory", resp.WorkingDirectory); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("valid_diagnostic_bundle_state", resp.ValidDiagnosticBundleState); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("windows_db_server", resp.WindowsDBServer); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("associated_time_machine_ids", utils.StringValueSlice(resp.AssociatedTimeMachineIds)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("access_key_id", resp.AccessKeyID); err != nil {
		return diag.FromErr(err)
	}

	props := []interface{}{}
	for _, prop := range resp.Properties {
		props = append(props, map[string]interface{}{
			"name":  prop.Name,
			"value": prop.Value,
		})
	}
	if err := d.Set("properties", props); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tags", flattenDBTags(resp.Tags)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("vm_info", flattenDBServerVMInfo(resp.VMInfo)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*resp.ID)
	return nil
}

func flattenDBServerMetadata(pr *era.DBServerMetadata) []interface{} {
	if pr != nil {
		metaList := make([]interface{}, 0)

		meta := map[string]interface{}{}

		meta["secure_info"] = pr.Secureinfo
		meta["info"] = pr.Info
		meta["deregister_info"] = flattenDeRegiserInfo(pr.Deregisterinfo)
		meta["database_type"] = pr.Databasetype
		meta["physical_era_drive"] = pr.Physicaleradrive
		meta["clustered"] = pr.Clustered
		meta["single_instance"] = pr.Singleinstance
		meta["era_drive_initialized"] = pr.Eradriveinitialised
		meta["provision_operation_id"] = pr.Provisionoperationid
		meta["marked_for_deletion"] = pr.Markedfordeletion
		if pr.Associatedtimemachines != nil {
			meta["associated_time_machines"] = utils.StringValueSlice(pr.Associatedtimemachines)
		}
		meta["software_snaphot_interval"] = pr.Softwaresnaphotinterval

		metaList = append(metaList, meta)
		return metaList
	}
	return nil
}

func flattenDBServerVMInfo(pr *era.VMInfo) []interface{} {
	if pr != nil {
		infoList := make([]interface{}, 0)
		info := map[string]interface{}{}

		info["secure_info"] = pr.SecureInfo
		info["info"] = pr.Info
		info["deregister_info"] = flattenDeRegiserInfo(pr.DeregisterInfo)
		info["os_type"] = pr.OsType
		info["os_version"] = pr.OsVersion
		info["distribution"] = pr.Distribution
		info["network_info"] = flattenVMNetworkInfo(pr.NetworkInfo)

		infoList = append(infoList, info)
		return infoList
	}
	return nil
}

func flattenVMNetworkInfo(pr []*era.NetworkInfo) []interface{} {
	if len(pr) > 0 {
		netList := make([]interface{}, len(pr))

		for k, v := range pr {
			nwt := make(map[string]interface{})

			nwt["vlan_name"] = v.VlanName
			nwt["vlan_uuid"] = v.VlanUUID
			nwt["vlan_type"] = v.VlanType
			nwt["era_configured"] = v.EraConfigured
			nwt["gateway"] = v.Gateway
			nwt["subnet_mask"] = v.SubnetMask
			nwt["hostname"] = v.Hostname
			nwt["device_name"] = v.DeviceName
			nwt["mac_address"] = v.MacAddress
			nwt["flags"] = v.Flags
			nwt["mtu"] = v.Mtu
			nwt["ip_addresses"] = utils.StringValueSlice(v.IPAddresses)
			nwt["default_gateway_device"] = v.DefaultGatewayDevice
			nwt["access_info"] = flattenAccessInfo(v.AccessInfo)

			netList[k] = nwt
		}
		return netList
	}
	return nil
}

func flattenAccessInfo(pr []*era.AccessInfo) []interface{} {
	if len(pr) > 0 {
		accessList := make([]interface{}, len(pr))

		for k, v := range pr {
			access := make(map[string]interface{})

			access["access_type"] = v.AccessType
			access["destination_subnet"] = v.DestinationSubnet

			accessList[k] = access
		}
		return accessList
	}
	return nil
}
