package clustersv2

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sort"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	import4 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/common/v1/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/prism/v4/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/clusters"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

//nolint:misspell // British English spelling is intentional
const (
	CANCELED = "CANCELLED"
)

func ResourceNutanixClusterV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixClusterV2Create,
		ReadContext:   ResourceNutanixClusterV2Read,
		UpdateContext: ResourceNutanixClusterV2Update,
		DeleteContext: ResourceNutanixClusterV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dryrun": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"expand": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"nodes": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"number_of_nodes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"node_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"controller_vm_ip": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ipv4": SchemaForValuePrefixLengthResource(),
												"ipv6": SchemaForValuePrefixLengthResource(),
											},
										},
									},
									"node_uuid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"host_ip": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ipv4": SchemaForValuePrefixLengthResource(),
												"ipv6": SchemaForValuePrefixLengthResource(),
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"network": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external_address": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv4": SchemaForValuePrefixLengthResource(),
									"ipv6": SchemaForValuePrefixLengthResource(),
								},
							},
						},
						"external_data_services_ip": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv4": SchemaForValuePrefixLengthResource(),
									"ipv6": SchemaForValuePrefixLengthResource(),
								},
							},
						},
						"external_subnet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internal_subnet": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nfs_subnet_white_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"name_server_ip_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv4": SchemaForValuePrefixLengthResource(),
									"ipv6": SchemaForValuePrefixLengthResource(),
									"fqdn": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"value": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"ntp_server_ip_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv4": SchemaForValuePrefixLengthResource(),
									"ipv6": SchemaForValuePrefixLengthResource(),
									"fqdn": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"value": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"smtp_server": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"email_address": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"server": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip_address": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ipv4": SchemaForValuePrefixLengthResource(),
															"ipv6": SchemaForValuePrefixLengthResource(),
															"fqdn": {
																Type:     schema.TypeList,
																Optional: true,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"value": {
																			Type:     schema.TypeString,
																			Computed: true,
																		},
																	},
																},
															},
														},
													},
												},
												"port": {
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
												"username": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"password": {
													Type:      schema.TypeString,
													Sensitive: true,
													Optional:  true,
												},
											},
										},
									},
									"type": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"PLAIN", "STARTTLS", "SSL"}, false),
									},
								},
							},
						},
						"masquerading_ip": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv4": SchemaForValuePrefixLengthResource(),
									"ipv6": SchemaForValuePrefixLengthResource(),
								},
							},
						},
						"masquerading_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"management_server": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ipv4": SchemaForValuePrefixLengthResource(),
												"ipv6": SchemaForValuePrefixLengthResource(),
											},
										},
									},
									"type": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"VCENTER"}, false),
									},
									"is_drs_enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"is_registered": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"is_in_use": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"fqdn": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"key_management_server_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"backplane": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_segmentation_enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"vlan_tag": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"subnet":  SchemaForValuePrefixLengthResource(),
									"netmask": SchemaForValuePrefixLengthResource(),
								},
							},
						},
						"http_proxy_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_address": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ipv4": SchemaForValuePrefixLengthResource(),
												"ipv6": SchemaForValuePrefixLengthResource(),
											},
										},
									},
									"port": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"username": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"password": {
										Type:      schema.TypeString,
										Optional:  true,
										Sensitive: true,
									},
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"proxy_types": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTPS", "SOCKS"}, false),
										},
									},
								},
							},
						},
						"http_proxy_white_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target": {
										Type:     schema.TypeString,
										Required: true,
									},
									"target_type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"IPV6_ADDRESS", "HOST_NAME", "DOMAIN_NAME_SUFFIX", "IPV4_NETWORK_MASK", "IPV4_ADDRESS"}, false),
									},
								},
							},
						},
					},
				},
			},
			"config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"incarnation_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"build_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"build_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"full_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"commit_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"short_commit_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"hypervisor_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cluster_function": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{"AOS", "ONE_NODE", "TWO_NODE"}, false),
							},
						},
						"timezone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"authorized_public_key_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"key": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},

							DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
								log.Printf("[DEBUG] DiffSuppressFunc authorized_public_key_list k: %v, oldValue: %v, newValue: %v", k, oldValue, newValue)

								// Check if the list has changed
								if d.HasChange("config.0.authorized_public_key_list") {
									log.Printf("[DEBUG] authorized_public_key_list has changed \n")
									oldRaw, newRaw := d.GetChange("config.0.authorized_public_key_list")
									// Convert to lists of interfaces
									oldList := oldRaw.([]interface{})
									newList := newRaw.([]interface{})
									// Sort lists based on a unique field (e.g., "key") for comparison
									sort.SliceStable(oldList, func(i, j int) bool {
										return oldList[i].(map[string]interface{})["key"].(string) < oldList[j].(map[string]interface{})["key"].(string)
									})
									sort.SliceStable(newList, func(i, j int) bool {
										return newList[i].(map[string]interface{})["key"].(string) < newList[j].(map[string]interface{})["key"].(string)
									})
									aJSON, _ := json.Marshal(oldList)
									log.Printf("[DEBUG] authorized_public_key_list oldList: %v", string(aJSON))
									aJSON, _ = json.Marshal(newList)
									log.Printf("[DEBUG] authorized_public_key_list newList: %v", string(aJSON))
									// Check if lists are equal
									if reflect.DeepEqual(oldList, newList) {
										log.Printf("[DEBUG] authorized_public_key_list are  equal \n")
										return true
									}

									log.Printf("[DEBUG] authorized_public_key_list are not equal \n")
									return false
								}
								log.Printf("[DEBUG] authorized_public_key_list has not changed \n")
								return false
							},
						},
						"redundancy_factor": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"cluster_software_map": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"software_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"version": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"cluster_arch": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"fault_tolerance_state": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"current_max_fault_tolerance": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"desired_max_fault_tolerance": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"domain_awareness_level": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"RACK", "NODE", "BLOCK", "DISK"}, false),
									},
									"current_cluster_fault_tolerance": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"desired_cluster_fault_tolerance": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"CFT_1N_OR_1D", "CFT_2N_OR_2D", "CFT_1N_AND_1D", "CFT_0N_AND_0D"}, false),
									},
									"redundancy_status": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"is_cassandra_preparation_done": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"is_zookeeper_preparation_done": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"is_remote_support_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"operation_mode": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"is_lts": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_password_remote_login_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"encryption_in_transit_status": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"encryption_option": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"encryption_scope": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"pulse_status": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"pii_scrubbing_level": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice([]string{"ALL", "DEFAULT"}, false),
									},
								},
							},
						},
						"is_available": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"upgrade_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"container_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"categories": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"vm_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"inefficient_vm_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cluster_profile_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_eligibility_score": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			// Computed fields
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rel": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func SchemaForValuePrefixLengthResource() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"value": {
					Type:     schema.TypeString,
					Required: true,
				},
				"prefix_length": {
					Type:     schema.TypeInt,
					Optional: true,
					Default:  defaultValue,
				},
			},
		},
	}
}

func ResourceNutanixClusterV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI
	body := config.NewCluster()
	var dryRun *bool

	if dryRunVar, ok := d.GetOk("dryrun"); ok {
		dryRun = utils.BoolPtr(dryRunVar.(bool))
	} else {
		dryRun = utils.BoolPtr(false)
	}
	if name, ok := d.GetOk("name"); ok {
		body.Name = utils.StringPtr(name.(string))
	}
	if nodes, ok := d.GetOk("nodes"); ok {
		body.Nodes = expandNodeReference(nodes)
	}
	if network, ok := d.GetOk("network"); ok {
		body.Network = expandClusterNetworkReference(network)
	}

	if configVar, ok := d.GetOk("config"); ok {
		body.Config = expandClusterConfigReference(configVar, d)
	}
	if upgradeStatus, ok := d.GetOk("upgrade_status"); ok {
		body.UpgradeStatus = expandUpgradeStatus(upgradeStatus)
	}

	if containerName, ok := d.GetOk("container_name"); ok {
		body.ContainerName = utils.StringPtr(containerName.(string))
	}

	if categories, ok := d.GetOk("categories"); ok {
		categoriesList := categories.([]interface{})
		categoriesListStr := common.ExpandListOfString(categoriesList)
		log.Printf("[DEBUG] categories List : %v", categoriesListStr)
		body.Categories = categoriesListStr
	}

	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] Create Cluster Request Body: %s", string(aJSON))

	resp, err := conn.ClusterEntityAPI.CreateCluster(body, dryRun)
	if err != nil {
		return diag.Errorf("error while creating clusters : %v", err)
	}

	TaskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the cluster to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for cluster (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching cluster UUID : %v", err)
	}
	rUUID := resourceUUID.Data.GetValue().(import2.Task)
	aJSON, _ = json.MarshalIndent(rUUID, "", "  ")
	log.Printf("[DEBUG] Create Cluster Task Response Details: %s", string(aJSON))

	randomID := utils.GenUUID()

	d.SetId(randomID)

	return nil
}

func ResourceNutanixClusterV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI
	var expand *string

	if expandVar, ok := d.GetOk("expand"); ok {
		expand = utils.StringPtr(expandVar.(string))
	} else {
		expand = nil
	}

	if d.Get("ext_id").(string) == "" {
		log.Printf("[DEBUG] ResourceNutanixClusterV2Read : extID is empty")
		err := getClusterExtID(d, conn)
		if err != nil {
			log.Printf("[DEBUG] ResourceNutanixClusterV2Read : Cluster not found, err -> %v", err)
			return diag.Errorf("error while fetching cluster : %v", err)
		}
	}

	log.Printf("[DEBUG] ResourceNutanixClusterV2Read : Cluster found, extID : %s", d.Id())
	resp, err := conn.ClusterEntityAPI.GetClusterById(utils.StringPtr(d.Id()), expand)
	if err != nil {
		log.Printf("[DEBUG] ResourceNutanixClusterV2Read : Cluster %s not found", d.Id())
		return diag.Errorf("error while fetching cluster : %v", err)
	}

	getResp := resp.Data.GetValue().(config.Cluster)
	aJSON, _ := json.MarshalIndent(getResp, "", "  ")
	log.Printf("[DEBUG] Read Cluster Response Details: %s", string(aJSON))

	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("nodes", flattenNodeReference(getResp.Nodes)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("network", flattenClusterNetworkReference(getResp.Network)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("config", flattenClusterConfigReference(getResp.Config)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("upgrade_status", flattenUpgradeStatus(getResp.UpgradeStatus)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("categories", getResp.Categories); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vm_count", getResp.VmCount); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("inefficient_vm_count", getResp.InefficientVmCount); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("container_name", getResp.ContainerName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_profile_ext_id", getResp.ClusterProfileExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("backup_eligibility_score", getResp.BackupEligibilityScore); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ResourceNutanixClusterV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI
	var expand *string

	if expandVar, ok := d.GetOk("expand"); ok {
		expand = utils.StringPtr(expandVar.(string))
	} else {
		expand = nil
	}

	if d.Get("ext_id").(string) == "" {
		log.Printf("[DEBUG] ResourceNutanixClusterV2Update : Cluster not found, extID is empty")
		err := getClusterExtID(d, conn)
		if err != nil {
			log.Printf("[DEBUG] ResourceNutanixClusterV2Update : Cluster not found, err -> %v", err)
			return diag.Errorf("error while fetching cluster : %v", err)
		}
	}

	log.Printf("[DEBUG] ResourceNutanixClusterV2Update : Cluster found, extID : %s", d.Id())

	resp, err := conn.ClusterEntityAPI.GetClusterById(utils.StringPtr(d.Id()), expand)
	if err != nil {
		return diag.Errorf("error while fetching cluster : %v", err)
	}

	// get etag value from read response to pass in update request If-Match header, Required for update request
	args := getEtagHeader(resp, conn)

	updateSpec := config.Cluster{}

	if d.HasChange("name") {
		updateSpec.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("nodes") {
		updateSpec.Nodes = expandNodeReference(d.Get("nodes"))
	}
	if d.HasChange("network") {
		updateSpec.Network = expandClusterNetworkReference(d.Get("network"))
	}
	if d.HasChange("config") {
		updateSpec.Config = expandClusterConfigReference(d.Get("config"), d)
	}
	if d.HasChange("upgrade_status") {
		updateSpec.UpgradeStatus = expandUpgradeStatus(d.Get("upgrade_status"))
	}

	if d.HasChange("container_name") {
		updateSpec.ContainerName = utils.StringPtr(d.Get("container_name").(string))
	}
	if d.HasChange("categories") {
		categories := d.Get("categories")
		categoriesList := categories.([]interface{})
		categoriesListStr := make([]string, len(categoriesList))
		for i, v := range categoriesList {
			categoriesListStr[i] = v.(string)
		}
		log.Printf("[DEBUG] categories List update Spec: %v", categoriesListStr)
		updateSpec.Categories = categoriesListStr
	}

	aJSON, _ := json.MarshalIndent(updateSpec, "", "  ")
	log.Printf("[DEBUG] Update Cluster Request Body: %s", string(aJSON))

	updateResp, err := conn.ClusterEntityAPI.UpdateClusterById(utils.StringPtr(d.Id()), &updateSpec, args)
	if err != nil {
		return diag.Errorf("error while updating clusters : %v", err)
	}

	TaskRef := updateResp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the cluster to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while updating clusters : %v", err)
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for cluster (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	rUUID := resourceUUID.Data.GetValue().(import2.Task)
	aJSON, _ = json.MarshalIndent(rUUID, "", "  ")
	log.Printf("[DEBUG] Update Cluster Task Response Details: %s", string(aJSON))

	//delay 1 min to get the updated data
	time.Sleep(1 * time.Minute)

	return ResourceNutanixClusterV2Read(ctx, d, meta)
}

func ResourceNutanixClusterV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI
	var expand *string

	if expandVar, ok := d.GetOk("expand"); ok {
		expand = utils.StringPtr(expandVar.(string))
	} else {
		expand = nil
	}

	if d.Get("ext_id").(string) == "" {
		err := getClusterExtID(d, conn)
		if err != nil {
			log.Printf("[DEBUG] ResourceNutanixClusterV2Delete : error while fetching cluster : %v", err)
			return diag.Errorf("error while fetching cluster : %v", err)
		}
	}

	readResp, err := conn.ClusterEntityAPI.GetClusterById(utils.StringPtr(d.Id()), expand)
	if err != nil {
		return diag.Errorf("error while reading cluster : %v", err)
	}
	// Extract E-Tag Header
	args := getEtagHeader(readResp, conn)

	var dryRun *bool

	if dryRunVar, ok := d.GetOk("dryrun"); ok {
		dryRun = utils.BoolPtr(dryRunVar.(bool))
	} else {
		dryRun = utils.BoolPtr(false)
	}

	resp, err := conn.ClusterEntityAPI.DeleteClusterById(utils.StringPtr(d.Id()), dryRun, args)
	if err != nil {
		return diag.Errorf("error while deleting cluster : %v", err)
	}
	TaskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	// calling group API to poll for completion of task

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the cluster to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for cluster (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	return nil
}

func taskStateRefreshPrismTaskGroupFunc(ctx context.Context, client *prism.Client, taskUUID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		// data := base64.StdEncoding.EncodeToString([]byte("ergon"))
		// encodeUUID := data + ":" + taskUUID
		vresp, err := client.TaskRefAPI.GetTaskById(utils.StringPtr(taskUUID), nil)
		if err != nil {
			return "", "", (fmt.Errorf("error while polling prism task: %v", err))
		}

		// get the group results

		v := vresp.Data.GetValue().(import2.Task)

		if getTaskStatus(v.Status) == "CANCELED" || getTaskStatus(v.Status) == "FAILED" {
			return v, getTaskStatus(v.Status),
				fmt.Errorf("error_detail: %s, progress_message: %d", utils.StringValue(v.ErrorMessages[0].Message), utils.IntValue(v.ProgressPercentage))
		}
		return v, getTaskStatus(v.Status), nil
	}
}

func getTaskStatus(pr *import2.TaskStatus) string {
	const two, three, five, six, seven = 2, 3, 5, 6, 7
	if pr != nil {
		if *pr == import2.TaskStatus(six) {
			return "FAILED"
		}
		if *pr == import2.TaskStatus(seven) {
			return "CANCELED"
		}
		if *pr == import2.TaskStatus(two) {
			return "QUEUED"
		}
		if *pr == import2.TaskStatus(three) {
			return "RUNNING"
		}
		if *pr == import2.TaskStatus(five) {
			return "SUCCEEDED"
		}
	}
	return "UNKNOWN"
}

func getClusterExtID(d *schema.ResourceData, conn *clusters.Client) error {
	var filter string
	if d.HasChange("name") {
		// if name changed, get the old name to fetch the cluster, since the name will be updated after update request
		oldName, _ := d.GetChange("name")
		filter = fmt.Sprintf(`name eq '%s'`, oldName.(string))
	} else {
		filter = fmt.Sprintf(`name eq '%s'`, d.Get("name").(string))
	}

	log.Printf("[DEBUG] getClusterExtID filter : %s", filter)

	// get Cluster Ext Id
	listResp, err := conn.ClusterEntityAPI.ListClusters(nil, nil, utils.StringPtr(filter), nil, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("error while fetching cluster : %v", err)
	}

	if listResp.Data == nil {
		log.Printf("[DEBUG] getClusterExtID Cluster not found, clustersResponse.Data is nil")
		return fmt.Errorf("cluster not found : %v", err)
	}
	cls := listResp.Data.GetValue().([]config.Cluster)

	if len(cls) == 0 {
		log.Printf("[DEBUG] getClusterExtID Cluster not found, len(clusters) is 0")
		return fmt.Errorf("cluster not found : %v", err)
	}
	extID := utils.StringValue(cls[0].ExtId)
	if extID == "" {
		log.Printf("[DEBUG] getClusterExtID Cluster not found, extID is empty")
		return fmt.Errorf("cluster not found : %v", err)
	}
	log.Printf("[DEBUG] getClusterExtID Cluster found, extId : %s", extID)
	d.SetId(extID)
	err = d.Set("ext_id", extID)
	if err != nil {
		return err
	}
	return nil
}

func expandNodeReference(pr interface{}) *config.NodeReference {
	if pr != nil {
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		nodeRef := config.NewNodeReference()

		if nodeList, ok := val["node_list"]; ok {
			nodeRef.NodeList = expandNodeListItemReference(nodeList.([]interface{}))
		}

		return nodeRef
	}

	return nil
}

func expandUpgradeStatus(upgradeStatus interface{}) *config.UpgradeStatus {
	const two, three, four, five, six, seven, eight, nine, ten = 2, 3, 4, 5, 6, 7, 8, 9, 10
	subMap := map[string]interface{}{
		"PENDING":     two,
		"DOWNLOADING": three,
		"QUEUED":      four,
		"PREUPGRADE":  five,
		"UPGRADING":   six,
		"SUCCEEDED":   seven,
		"FAILED":      eight,
		CANCELED:      nine,
		"SCHEDULED":   ten,
	}
	if subMap[upgradeStatus.(string)] != nil {
		pVal := subMap[upgradeStatus.(string)]
		p := config.UpgradeStatus(pVal.(int))
		return &p
	}
	log.Printf("[INFO] upgrade_status is not provided")
	return nil
}

func expandNodeListItemReference(pr []interface{}) []config.NodeListItemReference {
	if len(pr) > 0 {
		nodeList := make([]config.NodeListItemReference, len(pr))

		for k, v := range pr {
			val := v.(map[string]interface{})
			node := config.NewNodeListItemReference()

			if controllerVMIP, ok := val["controller_vm_ip"]; ok {
				log.Printf("[DEBUG] controller_vm_ip")
				node.ControllerVmIp = expandIPAddress(controllerVMIP)
			}
			if hostIP, ok := val["host_ip"]; ok {
				log.Printf("[DEBUG] host_ip")
				node.HostIp = expandIPAddress(hostIP)
			}

			nodeList[k] = *node
		}
		return nodeList
	}
	return nil
}

func expandClusterNetworkReference(pr interface{}) *config.ClusterNetworkReference {
	if pr != nil {
		cls := config.NewClusterNetworkReference()
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if externalAddress, ok := val["external_address"]; ok {
			log.Printf("[DEBUG] external_address")
			cls.ExternalAddress = expandIPAddress(externalAddress)
		}
		if externalDataServiceIP, ok := val["external_data_services_ip"]; ok {
			log.Printf("[DEBUG] external_data_services_ip")
			cls.ExternalDataServiceIp = expandIPAddress(externalDataServiceIP)
		}
		if nfsSubnetWhite, ok := val["nfs_subnet_white_list"]; ok {
			nfsSubnetWhitelist := nfsSubnetWhite.([]interface{})
			nfsSubnetWhitelistStr := make([]string, len(nfsSubnetWhitelist))
			for i, v := range nfsSubnetWhitelist {
				nfsSubnetWhitelistStr[i] = v.(string)
			}
			log.Printf("[DEBUG] nfs_subnet_white_list: %v", nfsSubnetWhitelistStr)
			cls.NfsSubnetWhitelist = nfsSubnetWhitelistStr
		}
		if nameServerIPList, ok := val["name_server_ip_list"]; ok {
			cls.NameServerIpList = expandIPAddressOrFQDN(nameServerIPList.([]interface{}))
		}
		if ntpServerIPList, ok := val["ntp_server_ip_list"]; ok {
			log.Printf("[DEBUG] ntp_server_ip_list ")
			cls.NtpServerIpList = expandIPAddressOrFQDN(ntpServerIPList.([]interface{}))
		}
		if smtpServer, ok := val["smtp_server"]; ok {
			cls.SmtpServer = expandSMTPServerRef(smtpServer)
		}
		if masqueradingIP, ok := val["masquerading_ip"]; ok {
			log.Printf("[DEBUG] masquerading_ip ")
			cls.MasqueradingIp = expandIPAddress(masqueradingIP)
		}
		if managementServer, ok := val["management_server"]; ok {
			cls.ManagementServer = expandManagementServerRef(managementServer)
		}
		if fqdn, ok := val["fqdn"]; ok && fqdn != "" {
			log.Printf("[DEBUG] network/fqdn : %v", fqdn)
			cls.Fqdn = utils.StringPtr(fqdn.(string))
		}
		if keyManagementServerType, ok := val["key_management_server_type"]; ok {
			log.Printf("[DEBUG] key_management_server_type : %s", keyManagementServerType)
			const zero, one, two, three, four = 0, 1, 2, 3, 4
			subMap := map[string]interface{}{
				"UNKNOWN":       zero,
				"$REDACTED":     one,
				"LOCAL":         two,
				"PRISM_CENTRAL": three,
				"EXTERNAL":      four,
			}
			if subMap[keyManagementServerType.(string)] != nil {
				pVal := subMap[keyManagementServerType.(string)]
				p := config.KeyManagementServerType(pVal.(int))
				cls.KeyManagementServerType = &p
			}
		}
		if backplane, ok := val["backplane"]; ok {
			cls.Backplane = expandBackplaneNetworkParams(backplane)
		}

		if httpProxyList, ok := val["http_proxy_list"]; ok {
			cls.HttpProxyList = expandHTTPProxyList(httpProxyList.([]interface{}))
		}
		if httpProxyWhiteList, ok := val["http_proxy_white_list"]; ok {
			cls.HttpProxyWhiteList = expandHTTPProxyWhiteList(httpProxyWhiteList.([]interface{}))
		}

		return cls
	}
	return nil
}

func expandHTTPProxyWhiteList(proxyTypesWhiteList []interface{}) []config.HttpProxyWhiteListConfig {
	if len(proxyTypesWhiteList) > 0 {
		httpProxyWhiteList := make([]config.HttpProxyWhiteListConfig, len(proxyTypesWhiteList))

		for k, v := range proxyTypesWhiteList {
			val := v.(map[string]interface{})
			httpProxy := config.NewHttpProxyWhiteListConfig()

			if target, ok := val["target"]; ok {
				httpProxy.Target = utils.StringPtr(target.(string))
			}
			if targetType, ok := val["target_type"]; ok {
				const two, three, four, five, six = 2, 3, 4, 5, 6
				subMap := map[string]interface{}{
					"IPV4_ADDRESS":       two,
					"IPV6_ADDRESS":       three,
					"IPV4_NETWORK_MASK":  four,
					"DOMAIN_NAME_SUFFIX": five,
					"HOST_NAME":          six,
				}
				if subMap[targetType.(string)] != nil {
					pVal := subMap[targetType.(string)]
					p := config.HttpProxyWhiteListTargetType(pVal.(int))
					httpProxy.TargetType = &p
				}
			}
			httpProxyWhiteList[k] = *httpProxy
		}
		return httpProxyWhiteList
	}
	return nil
}

func expandHTTPProxyList(httpProxyList []interface{}) []config.HttpProxyConfig {
	if len(httpProxyList) > 0 {
		httpProxyConfig := make([]config.HttpProxyConfig, len(httpProxyList))

		for k, v := range httpProxyList {
			val := v.(map[string]interface{})
			httpProxy := config.NewHttpProxyConfig()

			if ipAddr, ok := val["ip_address"]; ok {
				httpProxy.IpAddress = expandIPAddress(ipAddr)
			}
			if port, ok := val["port"]; ok {
				httpProxy.Port = utils.IntPtr(port.(int))
			}
			if username, ok := val["username"]; ok {
				httpProxy.Username = utils.StringPtr(username.(string))
			}
			if password, ok := val["password"]; ok {
				httpProxy.Password = utils.StringPtr(password.(string))
			}
			if name, ok := val["name"]; ok {
				httpProxy.Name = utils.StringPtr(name.(string))
			}
			if proxyTypes, ok := val["proxy_types"]; ok {
				if proxyTypes == nil || len(proxyTypes.([]interface{})) == 0 {
					httpProxy.ProxyTypes = nil
				} else {
					proxyTypesList := make([]config.HttpProxyType, len(proxyTypes.([]interface{})))
					const two, three, four = 2, 3, 4
					subMap := map[string]interface{}{
						"HTTP":  two,
						"HTTPS": three,
						"SOCKS": four,
					}
					for i, val := range proxyTypes.([]interface{}) {
						if subMap[val.(string)] != nil {
							pVal := subMap[val.(string)]
							p := config.HttpProxyType(pVal.(int))
							proxyTypesList[i] = p
						}
					}
					httpProxy.ProxyTypes = proxyTypesList
				}
			}
			httpProxyConfig[k] = *httpProxy
		}
		return httpProxyConfig
	}
	return nil
}

func expandClusterConfigReference(pr interface{}, d *schema.ResourceData) *config.ClusterConfigReference {
	if pr != nil {
		clsConf := config.NewClusterConfigReference()
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if buildInfo, ok := val["build_info"]; ok && d.HasChange("config.0.build_info") {
			clsConf.BuildInfo = expandBuildReference(buildInfo)
		}
		if clusterFunction, ok := val["cluster_function"]; ok && d.HasChange("config.0.cluster_function") {
			cfLen := len(clusterFunction.([]interface{}))
			cfs := make([]config.ClusterFunctionRef, cfLen)
			const two, three, four, five, six, seven, eight = 2, 3, 4, 5, 6, 7, 8
			subMap := map[string]interface{}{
				"AOS":                two,
				"PRISM_CENTRAL":      three,
				"CLOUD_DATA_GATEWAY": four,
				"AFS":                five,
				"ONE_NODE":           six,
				"TWO_NODE":           seven,
				"ANALYTICS_PLATFORM": eight,
			}

			for k, v := range clusterFunction.([]interface{}) {
				if subMap[v.(string)] != nil {
					pVal := subMap[v.(string)]
					p := config.ClusterFunctionRef(pVal.(int))
					cfs[k] = p
				}
			}
			clsConf.ClusterFunction = cfs
		}
		if _, ok := val["authorized_public_key_list"]; ok && d.HasChange("config.0.authorized_public_key_list") {
			_, newObj := d.GetChange("config.0.authorized_public_key_list")
			clsConf.AuthorizedPublicKeyList = expandPublicKey(newObj.([]interface{}))
		}
		if redundancyFactor, ok := val["redundancy_factor"]; ok && d.HasChange("config.0.redundancy_factor") {
			clsConf.RedundancyFactor = utils.Int64Ptr(int64(redundancyFactor.(int)))
		}
		if clusterArch, ok := val["cluster_arch"]; ok && d.HasChange("config.0.cluster_arch") {
			const two, three = 2, 3
			subMap := map[string]interface{}{
				"X86_64":  two,
				"PPC64LE": three,
			}
			if subMap[clusterArch.(string)] != nil {
				pVal := subMap[clusterArch.(string)]
				p := config.ClusterArchReference(pVal.(int))
				clsConf.ClusterArch = &p
			}
		}
		if faultToleranceState, ok := val["fault_tolerance_state"]; ok && d.HasChange("config.0.fault_tolerance_state") {
			clsConf.FaultToleranceState = expandFaultToleranceState(faultToleranceState)
		}
		if operationMode, ok := val["operation_mode"]; ok && d.HasChange("config.0.operation_mode") {
			const two, three, four, five, six = 2, 3, 4, 5, 6
			subMap := map[string]interface{}{
				"NORMAL":             two,
				"READ_ONLY":          three,
				"STAND_ALONE":        four,
				"SWITCH_TO_TWO_NODE": five,
				"OVERRIDE":           six,
			}
			if subMap[operationMode.(string)] != nil {
				pVal := subMap[operationMode.(string)]
				p := config.OperationMode(pVal.(int))
				clsConf.OperationMode = &p
			}
		}
		if encryptionInTransitStatus, ok := val["encryption_in_transit_status"]; ok && d.HasChange("config.0.encryption_in_transit_status") {
			const two, three = 2, 3
			subMap := map[string]interface{}{
				"ENABLED":  two,
				"DISABLED": three,
			}

			if subMap[encryptionInTransitStatus.(string)] != nil {
				pVal := subMap[encryptionInTransitStatus.(string)]
				p := config.EncryptionStatus(pVal.(int))
				clsConf.EncryptionInTransitStatus = &p
			}
		}

		if pulseStatus, ok := val["pulse_status"]; ok && d.HasChange("config.0.pulse_status") {
			clsConf.PulseStatus = expandPulseStatus(pulseStatus)
		}

		return clsConf
	}
	return nil
}

func expandPulseStatus(status interface{}) *config.PulseStatus {
	if status == nil || len(status.([]interface{})) == 0 {
		log.Printf("[DEBUG] PulseStatus is nil")
		return nil
	}

	pulse := config.NewPulseStatus()
	prI := status.([]interface{})
	val := prI[0].(map[string]interface{})

	if isEnabled, ok := val["is_enabled"]; ok {
		pulse.IsEnabled = utils.BoolPtr(isEnabled.(bool))
	}
	if piiScrubbingLevel, ok := val["pii_scrubbing_level"]; ok {
		const two, three = 2, 3
		subMap := map[string]interface{}{
			"DEFAULT": two,
			"ALL":     three,
		}
		if subMap[piiScrubbingLevel.(string)] != nil {
			pVal := subMap[piiScrubbingLevel.(string)]
			p := config.PIIScrubbingLevel(pVal.(int))
			pulse.PiiScrubbingLevel = &p
		}
	}
	return pulse
}

func expandIPAddress(pr interface{}) *import4.IPAddress {
	if pr != nil {
		ipAddress := import4.NewIPAddress()
		prI := pr.([]interface{})
		if len(prI) == 0 {
			return nil
		}
		val := prI[0].(map[string]interface{})

		if ipv4, ok := val["ipv4"]; ok && len(ipv4.([]interface{})) > 0 {
			ipAddress.Ipv4 = expandIPv4Address(ipv4)
		}
		if ipv6, ok := val["ipv6"]; ok && len(ipv6.([]interface{})) > 0 {
			ipAddress.Ipv6 = expandIPv6Address(ipv6)
		}
		aJSON, _ := json.Marshal(ipAddress)
		log.Printf("[DEBUG] ipAddress : %v", string(aJSON))
		return ipAddress
	}
	return nil
}

func expandIPAddressOrFQDN(pr []interface{}) []import4.IPAddressOrFQDN {
	if len(pr) > 0 {
		ips := make([]import4.IPAddressOrFQDN, len(pr))

		for k, v := range pr {
			val := v.(map[string]interface{})
			ip := import4.NewIPAddressOrFQDN()

			if ipv4, ok := val["ipv4"]; ok && len(ipv4.([]interface{})) > 0 {
				ip.Ipv4 = expandIPv4Address(ipv4)
			}
			if ipv6, ok := val["ipv6"]; ok && len(ipv6.([]interface{})) > 0 {
				ip.Ipv6 = expandIPv6Address(ipv6)
			}
			if fqdn, ok := val["fqdn"]; ok && len(fqdn.([]interface{})) > 0 {
				ip.Fqdn = expandFQDN(fqdn.([]interface{}))
			}
			ips[k] = *ip
		}
		aJSON, _ := json.Marshal(ips)
		log.Printf("[DEBUG] ipAddressOrFQDN : %v", string(aJSON))
		return ips
	}
	return nil
}

func expandIPv4Address(pr interface{}) *import4.IPv4Address {
	if len(pr.([]interface{})) == 0 {
		return nil
	}
	if pr != nil {
		ipv4 := import4.NewIPv4Address()
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if value, ok := val["value"]; ok {
			ipv4.Value = utils.StringPtr(value.(string))
		}
		if prefix, ok := val["prefix_length"]; ok {
			ipv4.PrefixLength = utils.IntPtr(prefix.(int))
		}
		return ipv4
	}
	return nil
}

func expandIPv6Address(pr interface{}) *import4.IPv6Address {
	if len(pr.([]interface{})) == 0 {
		return nil
	}

	if pr != nil {
		ipv6 := import4.NewIPv6Address()
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if value, ok := val["value"]; ok {
			ipv6.Value = utils.StringPtr(value.(string))
		}
		if prefix, ok := val["prefix_length"]; ok {
			ipv6.PrefixLength = utils.IntPtr(prefix.(int))
		}
		return ipv6
	}
	return nil
}

func expandSMTPServerRef(pr interface{}) *config.SmtpServerRef {
	if len(pr.([]interface{})) == 0 {
		return nil
	}

	if pr != nil {
		smtp := config.NewSmtpServerRef()
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if emailAddress, ok := val["email_address"]; ok {
			smtp.EmailAddress = utils.StringPtr(emailAddress.(string))
		}
		if server, ok := val["server"]; ok {
			smtp.Server = expandSMTPNetwork(server.([]interface{}))
		}
		if smtpType, ok := val["type"]; ok {
			const two, three, four = 2, 3, 4
			subMap := map[string]interface{}{
				"PLAIN":    two,
				"STARTTLS": three,
				"SSL":      four,
			}
			if subMap[smtpType.(string)] != nil {
				pVal := subMap[smtpType.(string)]
				p := config.SmtpType(pVal.(int))
				smtp.Type = &p
			}
		}
		return smtp
	}
	return nil
}

func expandBackplaneNetworkParams(pr interface{}) *config.BackplaneNetworkParams {
	if len(pr.([]interface{})) == 0 {
		return nil
	}

	if pr != nil && len(pr.([]interface{})) > 0 {
		backPlane := config.NewBackplaneNetworkParams()
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if isSegmentationEnabled, ok := val["is_segmentation_enabled"]; ok {
			backPlane.IsSegmentationEnabled = utils.BoolPtr(isSegmentationEnabled.(bool))
		}
		if subnet, ok := val["subnet"]; ok {
			backPlane.Subnet = expandIPv4Address(subnet)
		}
		if netmask, ok := val["netmask"]; ok {
			backPlane.Netmask = expandIPv4Address(netmask)
		}
		if vlanTag, ok := val["vlan_tag"]; ok {
			backPlane.VlanTag = utils.Int64Ptr(int64(vlanTag.(int)))
		}

		return backPlane
	}
	return nil
}

func expandManagementServerRef(pr interface{}) *config.ManagementServerRef {
	if pr != nil && len(pr.([]interface{})) > 0 {
		mgm := config.NewManagementServerRef()
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if ip, ok := val["ip"]; ok {
			log.Printf("[DEBUG] management server ip")
			mgm.Ip = expandIPAddress(ip.([]interface{}))
		}
		if mgmType, ok := val["type"]; ok {
			const two = 2
			switch mgmType.(string) {
			case "VCENTER":
				p := config.ManagementServerType(two)
				mgm.Type = &p
				log.Printf("[DEBUG] mgmType : VCENTER case")
			default:
				log.Printf("[DEBUG] mgmType : default case")
				mgm.Type = nil
			}
			log.Printf("[DEBUG] mgmType : %v", mgmType.(string))
		}
		if drsEnabled, ok := val["is_drs_enabled"]; ok {
			mgm.IsDrsEnabled = utils.BoolPtr(drsEnabled.(bool))
		}
		if isRegistered, ok := val["is_registered"]; ok {
			mgm.IsRegistered = utils.BoolPtr(isRegistered.(bool))
		}
		if inUse, ok := val["is_in_use"]; ok {
			mgm.IsInUse = utils.BoolPtr(inUse.(bool))
		}

		return mgm
	}
	return nil
}

func expandSMTPNetwork(pr []interface{}) *config.SmtpNetwork {
	if len(pr) > 0 {
		smtp := config.NewSmtpNetwork()
		val := pr[0].(map[string]interface{})

		if ipAddress, ok := val["ip_address"]; ok {
			smtp.IpAddress = &(expandIPAddressOrFQDN(ipAddress.([]interface{})))[0]
		}
		if port, ok := val["port"]; ok {
			smtp.Port = utils.IntPtr(port.(int))
		}
		if username, ok := val["username"]; ok {
			smtp.Username = utils.StringPtr(username.(string))
		}
		if password, ok := val["password"]; ok {
			smtp.Password = utils.StringPtr(password.(string))
		}

		return smtp
	}
	return nil
}

func expandFQDN(pr []interface{}) *import4.FQDN {
	if len(pr) > 0 {
		fqdn := import4.FQDN{}
		val := pr[0].(map[string]interface{})
		if value, ok := val["value"]; ok {
			fqdn.Value = utils.StringPtr(value.(string))
		}

		return &fqdn
	}
	return nil
}

func expandBuildReference(buildInfo interface{}) *config.BuildReference {
	if buildInfo == nil || len(buildInfo.([]interface{})) == 0 {
		log.Printf("[DEBUG] buildInfo is nil")
		return nil
	}

	buildReference := config.NewBuildReference()
	buildInfoI := buildInfo.([]interface{})
	buildInfoVal := buildInfoI[0].(map[string]interface{})

	if buildType, ok := buildInfoVal["build_type"]; ok {
		buildReference.BuildType = utils.StringPtr(buildType.(string))
	}
	if version, ok := buildInfoVal["version"]; ok {
		buildReference.Version = utils.StringPtr(version.(string))
	}
	if fullVersion, ok := buildInfoVal["full_version"]; ok {
		buildReference.FullVersion = utils.StringPtr(fullVersion.(string))
	}
	if commitID, ok := buildInfoVal["commit_id"]; ok {
		buildReference.CommitId = utils.StringPtr(commitID.(string))
	}
	if shortCommitID, ok := buildInfoVal["short_commit_id"]; ok {
		buildReference.ShortCommitId = utils.StringPtr(shortCommitID.(string))
	}

	return buildReference
}

func expandPublicKey(pr []interface{}) []config.PublicKey {
	if len(pr) > 0 {
		pubKey := make([]config.PublicKey, len(pr))
		aJSON, _ := json.Marshal(pr)
		log.Printf("[DEBUG] PublicKey : %v", string(aJSON))

		for k, v := range pr {
			val := v.(map[string]interface{})
			pub := config.NewPublicKey()

			if key, ok := val["key"]; ok {
				pub.Key = utils.StringPtr(key.(string))
			}
			if name, ok := val["name"]; ok {
				pub.Name = utils.StringPtr(name.(string))
			}
			pubKey[k] = *pub
		}
		return pubKey
	}
	return nil
}

func expandFaultToleranceState(pr interface{}) *config.FaultToleranceState {
	if pr != nil && len(pr.([]interface{})) > 0 {
		fts := config.NewFaultToleranceState()
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if domainAwarenessLevel, ok := val["domain_awareness_level"]; ok {
			const two, three, four, five = 2, 3, 4, 5
			subMap := map[string]interface{}{
				"NODE":  two,
				"BLOCK": three,
				"RACK":  four,
				"DISK":  five,
			}
			if subMap[domainAwarenessLevel.(string)] != nil {
				pVal := subMap[domainAwarenessLevel.(string)]
				p := config.DomainAwarenessLevel(pVal.(int))
				fts.DomainAwarenessLevel = &p
			}
		}

		if currentClusterFaultTolerance, ok := val["current_cluster_fault_tolerance"]; ok {
			const two, three, four, five = 2, 3, 4, 5
			subMap := map[string]interface{}{
				"CFT_0N_AND_0D": two,
				"CFT_1N_OR_1D":  three,
				"CFT_2N_OR_2D":  four,
				"CFT_1N_AND_1D": five,
			}
			if subMap[currentClusterFaultTolerance.(string)] != nil {
				pVal := subMap[currentClusterFaultTolerance.(string)]
				p := config.ClusterFaultToleranceRef(pVal.(int))
				fts.CurrentClusterFaultTolerance = &p
			}
		}
		if desiredClusterFaultTolerance, ok := val["desired_cluster_fault_tolerance"]; ok {
			const two, three, four, five = 2, 3, 4, 5
			subMap := map[string]interface{}{
				"CFT_0N_AND_0D": two,
				"CFT_1N_OR_1D":  three,
				"CFT_2N_OR_2D":  four,
				"CFT_1N_AND_1D": five,
			}
			if subMap[desiredClusterFaultTolerance.(string)] != nil {
				pVal := subMap[desiredClusterFaultTolerance.(string)]
				p := config.ClusterFaultToleranceRef(pVal.(int))
				fts.DesiredClusterFaultTolerance = &p
			}
		}
		return fts
	}
	return nil
}
