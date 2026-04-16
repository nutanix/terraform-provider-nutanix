package ndb

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	era "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixNDBCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixNDBClusterCreate,
		ReadContext:   resourceNutanixNDBClusterRead,
		UpdateContext: resourceNutanixNDBClusterUpdate,
		DeleteContext: resourceNutanixNDBClusterDelete,
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
			"cluster_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"storage_container": {
				Type:     schema.TypeString,
				Required: true,
			},
			"agent_vm_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "EraAgent",
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  "9440",
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "https",
			},
			"cluster_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "NTNX",
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "v2",
			},
			"agent_network_info": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dns": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ntp": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"networks_info": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"network_info": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vlan_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"static_ip": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"gateway": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"subnet_mask": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"access_type": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			// computed
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"unique_name": {
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
			"nx_cluster_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloud_type": {
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
			"hypervisor_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hypervisor_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"properties": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ref_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"reference_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cloud_info": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_threshold_percentage": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"memory_threshold_percentage": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
					},
				},
			},
			"management_server_info": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"entity_counts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_servers": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"engine_counts": engineCountSchema(),
					},
				},
			},
			"healthy": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceNutanixNDBClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	req := &era.ClusterIntentInput{}

	if name, ok := d.GetOk("name"); ok {
		req.ClusterName = utils.StringPtr(name.(string))
	}
	if desc, ok := d.GetOk("description"); ok {
		req.ClusterDescription = utils.StringPtr(desc.(string))
	}
	if clsip, ok := d.GetOk("cluster_ip"); ok {
		req.ClusterIP = utils.StringPtr(clsip.(string))
	}
	if storageContainer, ok := d.GetOk("storage_container"); ok {
		req.StorageContainer = utils.StringPtr(storageContainer.(string))
	}
	if protocol, ok := d.GetOk("protocol"); ok {
		req.Protocol = utils.StringPtr(protocol.(string))
	}
	if agentPrefix, ok := d.GetOk("agent_vm_prefix"); ok {
		req.AgentVMPrefix = utils.StringPtr(agentPrefix.(string))
	}
	if port, ok := d.GetOk("port"); ok {
		req.Port = utils.IntPtr(port.(int))
	}
	if clsType, ok := d.GetOk("cluster_type"); ok {
		req.ClusterType = utils.StringPtr(clsType.(string))
	}
	if version, ok := d.GetOk("version"); ok {
		req.Version = utils.StringPtr(version.(string))
	}

	if username, ok := d.GetOk("username"); ok {
		creds := make([]*era.NameValueParams, 0)

		creds = append(creds, &era.NameValueParams{
			Name:  utils.StringPtr("username"),
			Value: utils.StringPtr(username.(string)),
		})
		creds = append(creds, &era.NameValueParams{
			Name:  utils.StringPtr("password"),
			Value: utils.StringPtr(d.Get("password").(string)),
		})

		req.CredentialsInfo = creds
	}
	if agentNetInfo, ok := d.GetOk("agent_network_info"); ok {
		req.AgentNetworkInfo = expandCredentialInfo(agentNetInfo.([]interface{}))
	}
	if netinfo, ok := d.GetOk("networks_info"); ok {
		req.NetworksInfo = expandNetworkInfo(netinfo.([]interface{}))
	}
	// api to create cluster
	resp, err := conn.Service.CreateCluster(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	// Get Operation ID from response of Cluster and poll for the operation to get completed.
	opID := resp.Operationid
	if opID == "" {
		return diag.Errorf("error: operation ID is an empty string")
	}
	opReq := era.GetOperationRequest{
		OperationID: opID,
	}

	log.Printf("polling for operation with id: %s\n", opID)

	// Poll for operation here - Operation GET Call
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED", "FAILED"},
		Refresh: eraRefresh(ctx, conn, opReq),
		Timeout: d.Timeout(schema.TimeoutCreate),
		Delay:   eraDelay,
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for cluster (%s) to register: %s", resp.Entityid, errWaitTask)
	}

	clsName := d.Get("name")
	// api to fetch clusters based on name
	getResp, er := conn.Service.GetCluster(ctx, "", clsName.(string))
	if er != nil {
		return diag.FromErr(er)
	}
	d.SetId(*getResp.ID)
	log.Printf("NDB cluster with %s id is registered successfully", d.Id())
	return resourceNutanixNDBClusterRead(ctx, d, meta)
}

func resourceNutanixNDBClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	// check if d.Id() is nil
	if d.Id() == "" {
		return diag.Errorf("id is required for read operation")
	}

	resp, err := conn.Service.GetCluster(ctx, d.Id(), "")
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", resp.ID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", resp.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("unique_name", resp.Uniquename); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ip_addresses", resp.Ipaddresses); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("fqdns", resp.Fqdns); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("nx_cluster_uuid", resp.Nxclusteruuid); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("description", resp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cloud_type", resp.Cloudtype); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("date_created", resp.Datecreated); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("date_modified", resp.Datemodified); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("owner_id", resp.Ownerid); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("status", resp.Status); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("version", resp.Version); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("hypervisor_type", resp.Hypervisortype); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("hypervisor_version", resp.Hypervisorversion); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("properties", flattenClusterProperties(resp.Properties)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("reference_count", resp.Referencecount); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("username", resp.Username); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cloud_info", resp.Cloudinfo); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("resource_config", flattenResourceConfig(resp.Resourceconfig)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("management_server_info", resp.Managementserverinfo); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("entity_counts", flattenEntityCounts(resp.EntityCounts)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("healthy", resp.Healthy); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceNutanixNDBClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	req := &era.ClusterUpdateInput{}

	resp, err := conn.Service.GetCluster(ctx, d.Id(), "")
	if err != nil {
		return diag.FromErr(err)
	}

	if resp != nil {
		req.Name = resp.Name
		req.Description = resp.Description
		req.IPAddresses = resp.Ipaddresses
	}

	if d.HasChange("name") {
		req.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("description") {
		req.Description = utils.StringPtr(d.Get("description").(string))
	}
	if d.HasChange("cluster_ip") {
		ips := make([]*string, 0)
		clsIPs := d.Get("cluster_ip").([]interface{})

		for _, v := range clsIPs {
			ips = append(ips, utils.StringPtr(v.(string)))
		}
		req.IPAddresses = ips
	}

	if d.HasChange("username") {
		req.Username = utils.StringPtr(d.Get("username").(string))
	}
	if d.HasChange("password") {
		req.Password = utils.StringPtr(d.Get("password").(string))
	}

	// call update cluster API
	_, er := conn.Service.UpdateCluster(ctx, req, d.Id())
	if er != nil {
		return diag.FromErr(er)
	}

	log.Printf("NDB cluster with %s id is updated successfully", d.Id())
	return resourceNutanixNDBClusterRead(ctx, d, meta)
}

func resourceNutanixNDBClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Era

	req := &era.DeleteClusterInput{
		DeleteRemoteSites: false,
	}

	resp, err := conn.Service.DeleteCluster(ctx, req, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("Operation to delete cluster with id %s has started, operation id: %s", d.Id(), resp.Operationid)
	opID := resp.Operationid
	if opID == "" {
		return diag.Errorf("error: operation ID is an empty string")
	}
	opReq := era.GetOperationRequest{
		OperationID: opID,
	}

	log.Printf("polling for operation with id: %s\n", opID)

	// Poll for operation here - Cluster GET Call
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED", "FAILED"},
		Refresh: eraRefresh(ctx, conn, opReq),
		Timeout: d.Timeout(schema.TimeoutCreate),
		Delay:   eraDelay,
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for cluster (%s) to delete: %s", resp.Entityid, errWaitTask)
	}
	log.Printf("NDB cluster with %s id is deleted successfully", d.Id())
	return nil
}

func expandCredentialInfo(pr []interface{}) []*era.NameValueParams {
	if len(pr) > 0 {
		creds := make([]*era.NameValueParams, 0)

		cred := pr[0].(map[string]interface{})
		for key, v := range cred {
			creds = append(creds, &era.NameValueParams{
				Name:  utils.StringPtr(key),
				Value: utils.StringPtr(v.(string)),
			})
		}
		return creds
	}
	return nil
}

func expandNetworkInfo(pr []interface{}) []*era.NetworksInfo {
	if len(pr) > 0 {
		networkInfo := make([]*era.NetworksInfo, 0)

		for _, v := range pr {
			val := v.(map[string]interface{})
			netInfo := &era.NetworksInfo{}
			if netType, ok := val["type"]; ok {
				netInfo.Type = utils.StringPtr(netType.(string))
			}
			if infos, ok := val["network_info"]; ok {
				netInfo.NetworkInfo = expandClusterNetworkInfo(infos.([]interface{}))
			}
			if accessType, ok := val["access_type"]; ok {
				accessList := accessType.([]interface{})
				res := make([]*string, 0)
				for _, v := range accessList {
					res = append(res, utils.StringPtr(v.(string)))
				}
				netInfo.AccessType = res
			}
			networkInfo = append(networkInfo, netInfo)
		}
		return networkInfo
	}
	return nil
}

func expandClusterNetworkInfo(pr []interface{}) []*era.NameValueParams {
	if len(pr) > 0 {
		networkInfos := make([]*era.NameValueParams, 0)

		for _, v := range pr {
			val := v.(map[string]interface{})

			if vlan, ok := val["vlan_name"]; ok {
				networkInfos = append(networkInfos, &era.NameValueParams{
					Name:  utils.StringPtr("vlanName"),
					Value: utils.StringPtr(vlan.(string)),
				})
			}

			if vlan, ok := val["static_ip"]; ok {
				networkInfos = append(networkInfos, &era.NameValueParams{
					Name:  utils.StringPtr("staticIP"),
					Value: utils.StringPtr(vlan.(string)),
				})
			}

			if vlan, ok := val["gateway"]; ok {
				networkInfos = append(networkInfos, &era.NameValueParams{
					Name:  utils.StringPtr("gateway"),
					Value: utils.StringPtr(vlan.(string)),
				})
			}

			if vlan, ok := val["subnet_mask"]; ok {
				networkInfos = append(networkInfos, &era.NameValueParams{
					Name:  utils.StringPtr("subnetMask"),
					Value: utils.StringPtr(vlan.(string)),
				})
			}
		}
		return networkInfos
	}
	return nil
}
