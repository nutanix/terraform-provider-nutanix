package prismv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	prismCommon "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/common/v1/config"
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
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address": {
													Type:     schema.TypeList,
													Required: true,
													Elem:     schemaForIPAddressOrFqdn(),
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
													Elem:     schemaForIPAddressOrFqdn(),
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
			"links": schemaForLinks(),
			"config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemaForPcConfig(),
			},
			"is_registered_with_hosting_cluster": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"network": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemaForPcNetwork(),
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

func ResourceNutanixClusterPCRegistrationV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Create PC Registration\n")

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

	// check which remote cluster spec is provided
	if domainManagerRemoteClusterSpec, ok := remoteCluster["domain_manager_remote_cluster_spec"].([]interface{}); ok && len(domainManagerRemoteClusterSpec) > 0 {
		// domain manager remote cluster spec is provided
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
		// aos remote cluster spec is provided
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
		// cluster reference is provided
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

	body.RemoteCluster = remoteClusterBodySpec

	aJSON, _ := json.Marshal(body)
	log.Printf("[DEBUG] PC Registration Request Body: %s", string(aJSON))

	// pass nil for the new dyRun flag
	resp, err := conn.DomainManagerAPIInstance.Register(&pcExtID, body, nil, args)
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

func ResourceNutanixClusterPCRegistrationV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Read PC Registration Details\n")
	conn := meta.(*conns.Client).PrismAPI

	pcExtID := d.Id()

	readResp, err := conn.DomainManagerAPIInstance.GetDomainManagerById(&pcExtID)
	if err != nil {
		return diag.Errorf("error while fetching domain manager with id %s : %v", pcExtID, err)
	}

	pcBody := readResp.Data.GetValue().(prismConfig.DomainManager)
	aJSON, _ := json.MarshalIndent(pcBody, "", "  ")
	log.Printf("[DEBUG] PC Registration Read Response: %s", string(aJSON))

	// set attributes
	if err := d.Set("tenant_id", utils.StringValue(pcBody.TenantId)); err != nil {
		return diag.Errorf("error setting tenant_id: %s", err)
	}
	if err := d.Set("ext_id", utils.StringValue(pcBody.ExtId)); err != nil {
		return diag.Errorf("error setting ext_id: %s", err)
	}
	if err := d.Set("links", flattenLinks(pcBody.Links)); err != nil {
		return diag.Errorf("error setting links: %s", err)
	}
	if err := d.Set("config", flattenPCConfig(pcBody.Config)); err != nil {
		return diag.Errorf("error setting config: %s", err)
	}
	if err := d.Set("is_registered_with_hosting_cluster", utils.BoolValue(pcBody.IsRegisteredWithHostingCluster)); err != nil {
		return diag.Errorf("error setting is_registered_with_hosting_cluster: %s", err)
	}
	if err := d.Set("network", flattenPCNetwork(pcBody.Network)); err != nil {
		return diag.Errorf("error setting network: %s", err)
	}
	if err := d.Set("hosting_cluster_ext_id", utils.StringValue(pcBody.HostingClusterExtId)); err != nil {
		return diag.Errorf("error setting hosting_cluster_ext_id: %s", err)
	}
	if err := d.Set("should_enable_high_availability", utils.BoolValue(pcBody.ShouldEnableHighAvailability)); err != nil {
		return diag.Errorf("error setting should_enable_high_availability: %s", err)
	}
	if err := d.Set("node_ext_ids", pcBody.NodeExtIds); err != nil {
		return diag.Errorf("error setting node_ext_ids: %s", err)
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
			addressMap := address.([]interface{})[0].(map[string]interface{})
			ipAdd := expandIPAddressOrFqdn(addressMap)
			remoteClusterSpec.Address = &ipAdd
		}
		if credentials, ok := cluster["credentials"]; ok {
			remoteClusterSpec.Credentials = expandRegisterCredentials(credentials)
		}
		return remoteClusterSpec
	}
	return nil
}

func expandRegisterCredentials(credentials interface{}) *prismManagment.Credentials {
	if credentials != nil {
		credentialsMap := credentials.([]interface{})[0].(map[string]interface{})
		credentialsObj := &prismManagment.Credentials{}
		if authentication, ok := credentialsMap["authentication"]; ok {
			credentialsObj.Authentication = expandAuthentication(authentication)
		}
		return credentialsObj
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
		const OnpremCloud, NutanixHostedCloud = 2, 3
		subMap := map[string]interface{}{
			"ONPREM_CLOUD":         OnpremCloud,
			"NUTANIX_HOSTED_CLOUD": NutanixHostedCloud,
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
