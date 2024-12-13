package clustersv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	prismClusterMang "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/clustermgmt/v4/config"
	prismCommon "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/common/v1/config"
	prismResponse "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/common/v1/response"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	prismManagment "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/management"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const (
	DomainManagerRemoteClusterSpec = "prism.v4.management.DomainManagerRemoteClusterSpec"
	AOSRemoteClusterSpec           = "prism.v4.management.AOSRemoteClusterSpec"
	ClusterReference               = "prism.v4.management.ClusterReference"
)

var exactlyOneOfRemoteClusterSpec = []string{ // Exactly one of the following fields must be set
	"remote_cluster.0.domain_manager_remote_cluster_spec",
	"remote_cluster.0.aos_remote_cluster_spec",
	"remote_cluster.0.cluster_reference",
}

func ResourceNutanixClusterPCRegistrationV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixClusterPCRegistrationV2Create,
		ReadContext:   ResourceNutanixClusterPCRegistrationV2Read,
		UpdateContext: ResourceNutanixClusterPCRegistrationV2Update,
		DeleteContext: ResourceNutanixClusterPCRegistrationV2Delete,
		Schema: map[string]*schema.Schema{
			"pc_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"remote_cluster": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_manager_remote_cluster_spec": {
							Type:         schema.TypeList,
							MaxItems:     1,
							Optional:     true,
							ExactlyOneOf: exactlyOneOfRemoteClusterSpec,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"remote_cluster": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address": {
													Type:     schema.TypeList,
													Required: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ipv4": SchemaForValuePrefixLengthResource(),
															"ipv6": SchemaForValuePrefixLengthResource(),
															"fqdn": schemaForFQDNValueResource(),
														},
													},
												},
												"credentials": {
													Type:     schema.TypeList,
													Required: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"authentication": {
																Type:     schema.TypeList,
																Required: true,
																MaxItems: 1,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"username": {
																			Type:     schema.TypeString,
																			Required: true,
																		},
																		"password": {
																			Type:     schema.TypeString,
																			Required: true,
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									}, // remote_cluster
									"cloud_type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"NUTANIX_HOSTED_CLOUD", "ONPREM_CLOUD"}, false),
									},
								},
							},
						},
						"aos_remote_cluster_spec": {
							Type:         schema.TypeList,
							MaxItems:     1,
							Optional:     true,
							ExactlyOneOf: exactlyOneOfRemoteClusterSpec,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"remote_cluster": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address": {
													Type:     schema.TypeList,
													Required: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ipv4": SchemaForValuePrefixLengthResource(),
															"ipv6": SchemaForValuePrefixLengthResource(),
															"fqdn": schemaForFQDNValueResource(),
														},
													},
												},
												"credentials": {
													Type:     schema.TypeList,
													Required: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"authentication": {
																Type:     schema.TypeList,
																Required: true,
																MaxItems: 1,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"username": {
																			Type:     schema.TypeString,
																			Required: true,
																		},
																		"password": {
																			Type:     schema.TypeString,
																			Required: true,
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									}, // remote_cluster
								},
							},
						},
						"cluster_reference": {
							Type:         schema.TypeList,
							MaxItems:     1,
							Optional:     true,
							ExactlyOneOf: exactlyOneOfRemoteClusterSpec,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ext_id": {
										Description: "Cluster UUID of a remote cluster.",
										Type:        schema.TypeString,
										Required:    true,
									},
								},
							},
						},
					},
				},
			}, // remote_cluster
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ext_id": {
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
			}, // links
			"config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"should_enable_lockdown_mode": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"build_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"version": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						}, // build_info
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bootstrap_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"environment_info": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"provider_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"provisioning_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						}, // bootstrap_config
						"resource_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"num_vcpus": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"memory_size_bytes": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"data_disk_size_bytes": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"container_ext_ids": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						}, // resource_config
					},
				},
			}, // config
			"is_registered_with_hosting_cluster": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"network": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external_address": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv4": SchemaForValuePrefixLength(),
									"ipv6": SchemaForValuePrefixLength(),
								},
							},
						},
						"name_servers": schemaForIPv4IPv6FQDNResource(),
						"ntp_servers":  schemaForIPv4IPv6FQDNResource(),
						"fqdn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"external_networks": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"default_gateway": schemaForIPv4IPv6FQDNResource(),
									"subnet_mask":     schemaForIPv4IPv6FQDNResource(),
									"ip_ranges": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"begin": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ipv4": SchemaForValuePrefixLength(),
															"ipv6": SchemaForValuePrefixLength(),
														},
													},
												}, // begin
												"end": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ipv4": SchemaForValuePrefixLength(),
															"ipv6": SchemaForValuePrefixLength(),
														},
													},
												}, // end
											},
										},
									}, // ip_ranges
									"network_ext_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						}, // external_networks
					},
				},
			}, // network
			"hosting_cluster_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"should_enable_high_availability": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"node_ext_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			}, // node_ext_ids
		},
	}
}

func schemaForIPv4IPv6FQDNResource() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ipv4": SchemaForValuePrefixLength(),
				"ipv6": SchemaForValuePrefixLength(),
				"fqdn": schemaForFQDNValueResource(),
			},
		},
	}
}

func schemaForFQDNValueResource() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"value": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func ResourceNutanixClusterPCRegistrationV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Create PC Registration\n")
	// validate attributes based on object_type
	// if err := validateAttributes(d); err != nil {
	//	return err
	//}

	conn := meta.(*conns.Client).PrismAPI

	pcExtID := d.Get("pc_ext_id").(string)

	readResp, err := conn.DomainManagerAPIInstance.GetDomainManagerById(&pcExtID)
	if err != nil {
		return diag.Errorf("error while fetching domain manager with id %s : %v", pcExtID, err)
	}

	// Extract etag value from the response
	etagValue := conn.DomainManagerAPIInstance.ApiClient.GetEtag(readResp)
	args := make(map[string]interface{})
	args["If-Match"] = utils.StringPtr(etagValue)

	// get remote cluster object
	remoteClusterObj, _ := d.GetOk("remote_cluster")

	// get remote cluster data
	remoteCluster := remoteClusterObj.([]interface{})[0].(map[string]interface{})
	log.Printf("[DEBUG] Remote Cluster: %v\n", remoteCluster)

	// create body spec based
	body := &prismManagment.ClusterRegistrationSpec{}
	remoteClusterBodySpec := &prismManagment.OneOfClusterRegistrationSpecRemoteCluster{}

	if domainManagerRemoteClusterSpec, ok := remoteCluster["domain_manager_remote_cluster_spec"].([]interface{}); ok && len(domainManagerRemoteClusterSpec) > 0 {
		domainManagerRemoteClusterData := domainManagerRemoteClusterSpec[0].(map[string]interface{})
		log.Printf("[DEBUG] %v is selected\n", domainManagerRemoteClusterData)
		domainManagerRemoteClusterObj := prismManagment.NewDomainManagerRemoteClusterSpec()
		if remoteClusterData := domainManagerRemoteClusterData["remote_cluster"].([]interface{})[0].(map[string]interface{}); remoteClusterData != nil {
			domainManagerRemoteClusterObj.RemoteCluster = expandDomainManagerRemoteCluster(remoteClusterData)
		}
		if cloudType, ok := domainManagerRemoteClusterData["cloud_type"].(string); ok {
			domainManagerRemoteClusterObj.CloudType = expandDomainManagerCloudType(cloudType)
		}
		err = remoteClusterBodySpec.SetValue(*domainManagerRemoteClusterObj)
		if err != nil {
			return diag.Errorf("error while setting Body Spec for %v: %v", DomainManagerRemoteClusterSpec, err)
		}
		aJSON, _ := json.Marshal(domainManagerRemoteClusterObj)
		log.Printf("[DEBUG] DomainManagerRemoteClusterSpec Body: %s\n", string(aJSON))
		body.RemoteCluster = remoteClusterBodySpec
	} else if aosRemoteClusterSpec, ok := remoteCluster["aos_remote_cluster_spec"].([]interface{}); ok && len(aosRemoteClusterSpec) > 0 {
		aosRemoteClusterData := aosRemoteClusterSpec[0].(map[string]interface{})
		log.Printf("[DEBUG] %v is selected", aosRemoteClusterData)
		aosRemoteClusterObj := prismManagment.NewAOSRemoteClusterSpec()
		if remoteClusterData := aosRemoteClusterData["remote_cluster"].([]interface{})[0].(map[string]interface{}); remoteClusterData != nil {
			aosRemoteClusterObj.RemoteCluster = expandDomainManagerRemoteCluster(remoteClusterData)
		}
		err = remoteClusterBodySpec.SetValue(*aosRemoteClusterObj)
		if err != nil {
			return diag.Errorf("error while setting Body Spec for %v: %v", AOSRemoteClusterSpec, err)
		}
		aJSON, _ := json.Marshal(aosRemoteClusterObj)
		log.Printf("[DEBUG] AOSRemoteClusterSpec Body: %s\n", string(aJSON))
		body.RemoteCluster = remoteClusterBodySpec
	} else if clusterReferenceSpec, ok := remoteCluster["cluster_reference"].([]interface{}); ok && len(clusterReferenceSpec) > 0 {
		clusterReferenceData := clusterReferenceSpec[0].(map[string]interface{})
		log.Printf("[DEBUG] %v is selected", clusterReferenceData)
		clusterReference := prismManagment.NewClusterReference()
		if extID, ok := clusterReferenceData["ext_id"].(string); ok {
			clusterReference.ExtId = utils.StringPtr(extID)
		}
		err = remoteClusterBodySpec.SetValue(*clusterReference)
		if err != nil {
			return diag.Errorf("error while setting Body Spec for %v: %v", ClusterReference, err)
		}
		aJSON, _ := json.Marshal(clusterReference)
		log.Printf("[DEBUG] ClusterReference Body: %s\n", string(aJSON))
		body.RemoteCluster = remoteClusterBodySpec
	} else {
		return diag.Errorf("non of [%v, %v, %v] is provided",
			DomainManagerRemoteClusterSpec, AOSRemoteClusterSpec, ClusterReference)
	}

	// set remote cluster body spec based on object_type
	// switch objectType := remoteCluster["object_type"].(string); objectType {
	//case DomainManagerRemoteClusterSpec:
	//	log.Printf("[DEBUG] %v is selected\n", DomainManagerRemoteClusterSpec)
	//	domainManagerRemoteClusterSpec := prismManagment.NewDomainManagerRemoteClusterSpec()
	//	if domainManagerRemoteCluster := remoteCluster["remote_cluster"].([]interface{})[0].(map[string]interface{}); domainManagerRemoteCluster != nil {
	//		domainManagerRemoteClusterSpec.RemoteCluster = expandDomainManagerRemoteCluster(domainManagerRemoteCluster)
	//	}
	//	if cloudType, ok := remoteCluster["cloud_type"].(string); ok {
	//		domainManagerRemoteClusterSpec.CloudType = expandDomainManagerCloudType(cloudType)
	//	}
	//	err = remoteClusterBodySpec.SetValue(*domainManagerRemoteClusterSpec)
	//	if err != nil {
	//		return diag.Errorf("error while setting Body Spec for %v: %v", DomainManagerRemoteClusterSpec, err)
	//	}
	//	aJSON, _ := json.Marshal(domainManagerRemoteClusterSpec)
	//	log.Printf("[DEBUG] DomainManagerRemoteClusterSpec Body: %s\n", string(aJSON))
	//	break
	//case AOSRemoteClusterSpec:
	//	log.Printf("[DEBUG] %v is selected", AOSRemoteClusterSpec)
	//	aosRemoteClusterSpec := prismManagment.NewAOSRemoteClusterSpec()
	//	if aosRemoteCluster := remoteCluster["remote_cluster"].([]interface{})[0].(map[string]interface{}); aosRemoteCluster != nil {
	//		aosRemoteClusterSpec.RemoteCluster = expandDomainManagerRemoteCluster(aosRemoteCluster)
	//	}
	//	err = remoteClusterBodySpec.SetValue(*aosRemoteClusterSpec)
	//	if err != nil {
	//		return diag.Errorf("error while setting Body Spec for %v: %v", DomainManagerRemoteClusterSpec, err)
	//	}
	//	aJSON, _ := json.Marshal(aosRemoteClusterSpec)
	//	log.Printf("[DEBUG] AOSRemoteClusterSpec Body: %s\n", string(aJSON))
	//	break
	//case ClusterReference:
	//	log.Printf("[DEBUG] %v is selected", ClusterReference)
	//	clusterReference := prismManagment.NewClusterReference()
	//	if extID, ok := remoteCluster["ext_id"].(string); ok {
	//		clusterReference.ExtId = utils.StringPtr(extID)
	//	}
	//	err = remoteClusterBodySpec.SetValue(*clusterReference)
	//	if err != nil {
	//		return diag.Errorf("error while setting Body Spec for %v: %v", DomainManagerRemoteClusterSpec, err)
	//	}
	//	aJSON, _ := json.Marshal(clusterReference)
	//	log.Printf("[DEBUG] ClusterReference Body: %s\n", string(aJSON))
	//	break
	//default:
	//
	//}

	body.RemoteCluster = remoteClusterBodySpec

	aJSON, _ := json.Marshal(body)
	log.Printf("[DEBUG] PC Registration Request Body: %s", string(aJSON))

	resp, err := conn.DomainManagerAPIInstance.Register(&pcExtID, body, args)
	if err != nil {
		return diag.Errorf("error while registering remote cluster with id %s : %v", pcExtID, err)
	}

	TaskRef := resp.Data.GetValue().(prismConfig.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the cluster to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for PC registration to complete: %v", err)
	}

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching PC Register task with id %s : %v", *taskUUID, err)
	}

	rUUID := resourceUUID.Data.GetValue().(prismConfig.Task)

	aJSON, _ = json.Marshal(rUUID)
	log.Printf("[DEBUG] PC Registration Task Details: %s", string(aJSON))

	d.SetId(pcExtID)
	return ResourceNutanixClusterPCRegistrationV2Read(ctx, d, meta)
}

func validateAttributes(d *schema.ResourceData) diag.Diagnostics {
	log.Printf("[DEBUG] validateAttributes\n")

	// handle validation for attributes based on object_type
	remoteCluster := d.Get("remote_cluster").([]interface{})[0].(map[string]interface{})

	objectType := remoteCluster["object_type"].(string)
	log.Printf("[DEBUG] object_type: %s\n", objectType)

	switch objectType {
	case DomainManagerRemoteClusterSpec:
		if remoteCluster["remote_cluster"] == nil || len(remoteCluster["remote_cluster"].([]interface{})) == 0 {
			return diag.Errorf("Missing required argument remote_cluster for object_type %s", objectType)
		}
		if remoteCluster["cloud_type"] == nil || remoteCluster["cloud_type"] == "" {
			return diag.Errorf("Missing required argument cloud_type for object_type %s", objectType)
		}
		log.Printf("[DEBUG] %v is validated\n", DomainManagerRemoteClusterSpec)
	case AOSRemoteClusterSpec:
		log.Printf("[DEBUG] remoteCluster: %v\n", remoteCluster)
		if remoteCluster["remote_cluster"] == nil || len(remoteCluster["remote_cluster"].([]interface{})) == 0 {
			return diag.Errorf("Missing required argument remote_cluster for object_type %s", objectType)
		}
		log.Printf("[DEBUG] %v is validated\n", AOSRemoteClusterSpec)
	case ClusterReference:
		if remoteCluster["ext_id"] == nil || remoteCluster["ext_id"] == "" {
			return diag.Errorf("Missing required argument cluster_ext_id for object_type %s", objectType)
		}
		log.Printf("[DEBUG] %v is validated\n", ClusterReference)
	default:
		return diag.Errorf(`Unsupported object_type: %s, it should be one of ["%v", "%v", "%v"]`,
			objectType, DomainManagerRemoteClusterSpec, AOSRemoteClusterSpec, ClusterReference)
	}
	log.Printf("[DEBUG] validateAttributes success\n")
	return nil
}

func ResourceNutanixClusterPCRegistrationV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Read PC Registration Details\n")
	conn := meta.(*conns.Client).PrismAPI

	pcExtID := d.Id()

	readResp, err := conn.DomainManagerAPIInstance.GetDomainManagerById(&pcExtID)
	if err != nil {
		return diag.Errorf("error while fetching domain manager with id %s : %v", pcExtID, err)
	}

	getResp := readResp.Data.GetValue().(prismConfig.DomainManager)
	aJSON, _ := json.Marshal(getResp)
	log.Printf("[DEBUG] PC Registration Read Response: %s", string(aJSON))

	// set attributes
	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.Errorf("error setting ext_id: %v", err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.Errorf("error setting tenant_id: %v", err)
	}
	if err := d.Set("links", flattenPrismLinks(getResp.Links)); err != nil {
		return diag.Errorf("error setting links: %v", err)
	}
	if err := d.Set("config", flattenPCConfig(getResp.Config)); err != nil {
		return diag.Errorf("error setting config: %v", err)
	}
	if err := d.Set("is_registered_with_hosting_cluster", getResp.IsRegisteredWithHostingCluster); err != nil {
		return diag.Errorf("error setting is_registered_with_hosting_cluster: %v", err)
	}
	if err := d.Set("network", flattenPCNetwork(getResp.Network)); err != nil {
		return diag.Errorf("error setting network: %v", err)
	}
	if err := d.Set("hosting_cluster_ext_id", getResp.HostingClusterExtId); err != nil {
		return diag.Errorf("error setting hosting_cluster_ext_id: %v", err)
	}
	if err := d.Set("should_enable_high_availability", getResp.ShouldEnableHighAvailability); err != nil {
		return diag.Errorf("error setting should_enable_high_availability: %v", err)
	}
	if err := d.Set("node_ext_ids", getResp.NodeExtIds); err != nil {
		return diag.Errorf("error setting node_ext_ids: %v", err)
	}

	return nil
}

func ResourceNutanixClusterPCRegistrationV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixClusterPCRegistrationV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func expandDomainManagerRemoteCluster(cluster map[string]interface{}) *prismManagment.RemoteClusterSpec {
	if cluster != nil {
		remoteClusterSpec := &prismManagment.RemoteClusterSpec{}
		if address, ok := cluster["address"]; ok {
			remoteClusterSpec.Address = expandDomainManagerRemoteClusterAddress(address)
		}
		if credentials, ok := cluster["credentials"]; ok {
			remoteClusterSpec.Credentials = expandCredentials(credentials)
		}
		return remoteClusterSpec
	}
	return nil
}

func expandDomainManagerRemoteClusterAddress(address interface{}) *prismCommon.IPAddressOrFQDN {
	if address != nil {
		addressMap := address.([]interface{})[0].(map[string]interface{})
		ip := &prismCommon.IPAddressOrFQDN{}
		if ipv4, ok := addressMap["ipv4"]; ok && len(ipv4.([]interface{})) > 0 {
			ip.Ipv4 = expandRemoteClusterIPv4Address(ipv4)
		}
		if ipv6, ok := addressMap["ipv6"]; ok && len(ipv6.([]interface{})) > 0 {
			ip.Ipv6 = expandRemoteClusterIPv6Address(ipv6)
		}
		if fqdn, ok := addressMap["fqdn"]; ok && len(fqdn.([]interface{})) > 0 {
			ip.Fqdn = expandRemoteClusterFQDN(fqdn)
		}
		return ip
	}
	return nil
}

func expandRemoteClusterIPv4Address(ipv4 interface{}) *prismCommon.IPv4Address {
	if ipv4 != nil {
		ipAddress := ipv4.([]interface{})[0].(map[string]interface{})
		ip := &prismCommon.IPv4Address{}
		if value, ok := ipAddress["value"].(string); ok {
			ip.Value = utils.StringPtr(value)
		}
		if prefixLength, ok := ipAddress["prefix_length"].(int); ok {
			ip.PrefixLength = utils.IntPtr(prefixLength)
		}
		return ip
	}
	return nil
}

func expandRemoteClusterIPv6Address(ipv6 interface{}) *prismCommon.IPv6Address {
	if ipv6 != nil {
		ipAddress := ipv6.([]interface{})[0].(map[string]interface{})
		ip := &prismCommon.IPv6Address{}
		if value, ok := ipAddress["value"].(string); ok {
			ip.Value = utils.StringPtr(value)
		}
		if prefixLength, ok := ipAddress["prefix_length"].(int); ok {
			ip.PrefixLength = utils.IntPtr(prefixLength)
		}
		return ip
	}
	return nil
}

func expandRemoteClusterFQDN(fqdn interface{}) *prismCommon.FQDN {
	if fqdn != nil {
		fqdnMap := fqdn.([]interface{})[0].(map[string]interface{})
		f := &prismCommon.FQDN{}
		if value, ok := fqdnMap["value"].(string); ok {
			f.Value = utils.StringPtr(value)
		}
		return f
	}
	return nil
}

func expandCredentials(credentials interface{}) *prismManagment.Credentials {
	if credentials != nil {
		credentialsMap := credentials.([]interface{})[0].(map[string]interface{})
		creds := &prismManagment.Credentials{}
		if authentication, ok := credentialsMap["authentication"]; ok {
			creds.Authentication = expandAuthentication(authentication)
		}
		return creds
	}
	return nil
}

func expandAuthentication(authentication interface{}) *prismCommon.BasicAuth {
	if authentication != nil {
		authMap := authentication.([]interface{})[0].(map[string]interface{})
		auth := &prismCommon.BasicAuth{}
		if username, ok := authMap["username"]; ok {
			auth.Username = utils.StringPtr(username.(string))
		}
		if password, ok := authMap["password"]; ok {
			auth.Password = utils.StringPtr(password.(string))
		}
		return auth
	}
	return nil
}

func expandDomainManagerCloudType(cloudType interface{}) *prismManagment.DomainManagerCloudType {
	if cloudType != nil && cloudType != "" {
		const two, three = 2, 3
		subMap := map[string]interface{}{
			"ONPREM_CLOUD":         two,
			"NUTANIX_HOSTED_CLOUD": three,
		}
		pVal := subMap[cloudType.(string)]
		if pVal == nil {
			return nil
		}
		p := prismManagment.DomainManagerCloudType(pVal.(int))
		return &p
	}
	return nil
}

// flattenLinks flattens the links field
func flattenPrismLinks(pr []prismResponse.ApiLink) []map[string]interface{} {
	if len(pr) > 0 {
		linkList := make([]map[string]interface{}, len(pr))

		for k, v := range pr {
			links := map[string]interface{}{}
			if v.Href != nil {
				links["href"] = v.Href
			}
			if v.Rel != nil {
				links["rel"] = v.Rel
			}

			linkList[k] = links
		}
		return linkList
	}
	return nil
}

func flattenPCNetwork(pcNetwork *prismConfig.DomainManagerNetwork) []map[string]interface{} {
	if pcNetwork != nil {
		network := make(map[string]interface{})
		if pcNetwork.ExternalAddress != nil {
			network["external_address"] = flattenPrismIPAddress(pcNetwork.ExternalAddress)
		}
		if pcNetwork.NameServers != nil {
			network["name_servers"] = flattenPrismIPAddressOrFQDN(pcNetwork.NameServers)
		}
		if pcNetwork.NtpServers != nil {
			network["ntp_servers"] = flattenPrismIPAddressOrFQDN(pcNetwork.NtpServers)
		}
		if pcNetwork.Fqdn != nil {
			network["fqdn"] = *pcNetwork.Fqdn
		}
		if pcNetwork.ExternalNetworks != nil {
			network["external_networks"] = flattenExternalNetworks(pcNetwork.ExternalNetworks)
		}
		return []map[string]interface{}{network}
	}
	return nil
}

func flattenExternalNetworks(externalNetworks []prismConfig.ExternalNetwork) []map[string]interface{} {
	if len(externalNetworks) > 0 {
		networks := make([]map[string]interface{}, len(externalNetworks))

		for k, v := range externalNetworks {
			network := make(map[string]interface{})
			ipAddressList := make([]prismCommon.IPAddressOrFQDN, 1)
			ipAddressList[0] = *v.DefaultGateway
			network["default_gateway"] = flattenPrismIPAddressOrFQDN(ipAddressList)
			ipAddressList[0] = *v.SubnetMask
			network["subnet_mask"] = flattenPrismIPAddressOrFQDN(ipAddressList)
			network["ip_ranges"] = flattenIPRanges(v.IpRanges)
			network["network_ext_id"] = v.NetworkExtId

			networks[k] = network
		}
		return networks
	}
	return nil
}

func flattenIPRanges(ipRanges []prismCommon.IpRange) []map[string]interface{} {
	if len(ipRanges) > 0 {
		ranges := make([]map[string]interface{}, len(ipRanges))

		for k, v := range ipRanges {
			ipRange := make(map[string]interface{})
			ipRange["begin"] = flattenPrismIPAddress(v.Begin)
			ipRange["end"] = flattenPrismIPAddress(v.End)

			ranges[k] = ipRange
		}
		return ranges
	}
	return nil
}

func flattenPrismIPAddress(ipAddress *prismCommon.IPAddress) []map[string]interface{} {
	if ipAddress != nil {
		ip := make(map[string]interface{})

		ip["ipv4"] = flattenPrismIPv4Address(ipAddress.Ipv4)
		ip["ipv6"] = flattenPrismIPv6Address(ipAddress.Ipv6)

		return []map[string]interface{}{ip}
	}
	return nil
}

func flattenPrismIPAddressOrFQDN(pr []prismCommon.IPAddressOrFQDN) interface{} {
	if len(pr) > 0 {
		ips := make([]map[string]interface{}, len(pr))

		for k, v := range pr {
			ip := make(map[string]interface{})

			ip["ipv4"] = flattenPrismIPv4Address(v.Ipv4)
			ip["ipv6"] = flattenPrismIPv6Address(v.Ipv6)
			ip["fqdn"] = flattenPrismFQDN(v.Fqdn)

			ips[k] = ip
		}
		return ips
	}
	return nil
}

func flattenPrismIPv4Address(pr *prismCommon.IPv4Address) []interface{} {
	if pr != nil {
		ipv4 := make([]interface{}, 0)

		ip := make(map[string]interface{})

		ip["value"] = pr.Value
		ip["prefix_length"] = pr.PrefixLength

		ipv4 = append(ipv4, ip)

		return ipv4
	}
	return nil
}

func flattenPrismIPv6Address(pr *prismCommon.IPv6Address) []interface{} {
	if pr != nil {
		ipv6 := make([]interface{}, 0)

		ip := make(map[string]interface{})

		ip["value"] = pr.Value
		ip["prefix_length"] = pr.PrefixLength

		ipv6 = append(ipv6, ip)

		return ipv6
	}
	return nil
}

func flattenPrismFQDN(pr *prismCommon.FQDN) []interface{} {
	if pr != nil {
		fqdn := make([]interface{}, 0)

		f := make(map[string]interface{})

		f["value"] = pr.Value

		fqdn = append(fqdn, f)

		return fqdn
	}
	return nil
}

func flattenPCConfig(pcConfig *prismConfig.DomainManagerClusterConfig) []map[string]interface{} {
	if pcConfig != nil {
		config := make(map[string]interface{})
		if pcConfig.ShouldEnableLockdownMode != nil {
			config["should_enable_lockdown_mode"] = *pcConfig.ShouldEnableLockdownMode
		}
		if pcConfig.BuildInfo != nil {
			config["build_info"] = flattenBuildInfo(pcConfig.BuildInfo)
		}
		if pcConfig.Name != nil {
			config["name"] = *pcConfig.Name
		}
		if pcConfig.Size != nil {
			config["size"] = flattenClusterSize(*pcConfig.Size)
		}
		if pcConfig.BootstrapConfig != nil {
			config["bootstrap_config"] = flattenBootstrapConfig(pcConfig.BootstrapConfig)
		}
		if pcConfig.ResourceConfig != nil {
			config["resource_config"] = flattenResourceConfig(pcConfig.ResourceConfig)
		}
		return []map[string]interface{}{config}
	}
	return nil
}

func flattenBuildInfo(buildInfo *prismClusterMang.BuildInfo) []map[string]interface{} {
	if buildInfo != nil {
		info := make(map[string]interface{})
		if buildInfo.Version != nil {
			info["version"] = *buildInfo.Version
		}
		return []map[string]interface{}{info}
	}
	return nil
}

func flattenClusterSize(size prismConfig.Size) string {
	const STARTER, SMALL, LARGE, EXTRALARGE = 2, 3, 4, 5

	switch size {
	case STARTER:
		return "STARTER"
	case SMALL:
		return "SMALL"
	case LARGE:
		return "LARGE"
	case EXTRALARGE:
		return "EXTRALARGE"
	default:
		return "UNKNOWN"
	}
}

func flattenBootstrapConfig(bootstrapConfig *prismConfig.BootstrapConfig) []map[string]interface{} {
	if bootstrapConfig != nil {
		config := make(map[string]interface{})
		if bootstrapConfig.EnvironmentInfo != nil {
			config["environment_info"] = flattenEnvironmentInfo(bootstrapConfig.EnvironmentInfo)
		}
		return []map[string]interface{}{config}
	}
	return nil
}

func flattenEnvironmentInfo(environmentInfo *prismConfig.EnvironmentInfo) []map[string]interface{} {
	if environmentInfo != nil {
		info := make(map[string]interface{})
		if environmentInfo.Type != nil {
			info["type"] = flattenEnvironmentType(*environmentInfo.Type)
		}
		if environmentInfo.ProviderType != nil {
			info["provider_type"] = flattenEnvironmentProviderType(*environmentInfo.ProviderType)
		}
		if environmentInfo.ProvisioningType != nil {
			info["provisioning_type"] = flattenEnvironmentProvisioningType(*environmentInfo.ProvisioningType)
		}
		return []map[string]interface{}{info}
	}
	return nil
}

func flattenEnvironmentType(environmentType prismConfig.EnvironmentType) string {
	const onprem, ntnxCloud = 2, 3

	switch environmentType {
	case onprem:
		return "ONPREM"
	case ntnxCloud:
		return "NTNX_CLOUD"
	default:
		return "UNKNOWN"
	}
}

func flattenEnvironmentProviderType(providerType prismConfig.ProviderType) string {
	const ntnx, azure, aws, gcp, vsphere = 2, 3, 4, 5, 6
	switch providerType {
	case ntnx:
		return "NTNX"
	case azure:
		return "AZURE"
	case aws:
		return "AWS"
	case gcp:
		return "GCP"
	case vsphere:
		return "VSPHERE"
	default:
		return "UNKNOWN"
	}
}

func flattenEnvironmentProvisioningType(provisioningType prismConfig.ProvisioningType) string {
	const ntnx, native = 2, 3
	switch provisioningType {
	case ntnx:
		return "NTNX"
	case native:
		return "NATIVE"
	default:
		return "UNKNOWN"
	}
}

func flattenResourceConfig(resourceConfig *prismConfig.DomainManagerResourceConfig) []map[string]interface{} {
	if resourceConfig != nil {
		config := make(map[string]interface{})
		if resourceConfig.NumVcpus != nil {
			config["num_vcpus"] = utils.IntValue(resourceConfig.NumVcpus)
		}
		if resourceConfig.MemorySizeBytes != nil {
			config["memory_size_bytes"] = utils.Int64Value(resourceConfig.MemorySizeBytes)
		}
		if resourceConfig.DataDiskSizeBytes != nil {
			config["data_disk_size_bytes"] = utils.Int64Value(resourceConfig.DataDiskSizeBytes)
		}
		if resourceConfig.ContainerExtIds != nil {
			config["container_ext_ids"] = resourceConfig.ContainerExtIds
		}
		return []map[string]interface{}{config}
	}
	return nil
}
