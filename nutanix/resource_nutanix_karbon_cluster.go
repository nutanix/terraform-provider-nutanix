package nutanix

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	karbon "github.com/terraform-providers/terraform-provider-nutanix/client/karbon"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"

	"github.com/hashicorp/terraform/helper/schema"
)

const (
	maxMasterNodes    = 5
	minMasterNodes    = 2
	cpuDivisionAmount = 2
)

func resourceNutanixKarbonCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceNutanixKarbonClusterCreate,
		Read:   resourceNutanixKarbonClusterRead,
		Update: resourceNutanixKarbonClusterUpdate,
		Delete: resourceNutanixKarbonClusterDelete,
		Exists: resourceNutanixKarbonClusterExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 1,
		Schema:        KarbonClusterResourceMap(),
	}
}

func KarbonClusterResourceMap() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"version": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"deployment_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"kubeapi_server_ipv4_address": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"storage_class_config": {
			Type:     schema.TypeSet,
			Required: true,
			// ForceNew: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"reclaim_policy": {
						Type:     schema.TypeString,
						Required: true,
					},
					"volumes_config": {
						Type:     schema.TypeMap,
						Required: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"file_system": {
									Type:     schema.TypeString,
									Required: true,
								},
								"flash_mode": {
									Type:     schema.TypeBool,
									Required: true,
								},
								"password": {
									Type:      schema.TypeString,
									Required:  true,
									Sensitive: true,
								},
								"prism_element_cluster_uuid": {
									Type:     schema.TypeString,
									Required: true,
								},
								"storage_container": {
									Type:     schema.TypeString,
									Required: true,
								},
								"username": {
									Type:     schema.TypeString,
									Required: true,
								},
							},
						},
					},
				},
			},
		},
		"active_passive_config": {
			Type:          schema.TypeSet,
			Optional:      true,
			ForceNew:      true,
			MaxItems:      1,
			ConflictsWith: []string{"external_lb_config"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"external_ipv4_address": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"external_lb_config": {
			Type:          schema.TypeSet,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"active_passive_config"},
			MaxItems:      1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"external_ipv4_address": {
						Type:     schema.TypeString,
						Required: true,
					},
					"master_nodes_config": {
						Type:     schema.TypeSet,
						Required: true,
						MaxItems: maxMasterNodes,
						MinItems: minMasterNodes,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"ipv4_address": {
									Type:     schema.TypeString,
									Required: true,
								},
								"node_pool_name": {
									Type:     schema.TypeString,
									Required: true,
								},
							},
						},
					},
				},
			},
		},
		"private_registries": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"registry_name": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"etcd_node_pool":   nodePoolSchema(true),
		"master_node_pool": nodePoolSchema(true),
		"worker_node_pool": nodePoolSchema(true),
		"cni_config":       CNISchema(),
	}
}

func CNISchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Required: true,
		MaxItems: 1,
		ForceNew: true,
		// Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"node_cidr_mask_size": {
					Type:     schema.TypeInt,
					Required: true,
				},
				"pod_ipv4_cidr": {
					Type:     schema.TypeString,
					Required: true,
				},
				"service_ipv4_cidr": {
					Type:     schema.TypeString,
					Required: true,
				},
				"flannel_config": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{},
					},
				},
				"calico_config": {
					Type:     schema.TypeSet,
					Optional: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"ip_pool_configs": {
								Type:     schema.TypeList,
								Optional: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"cidr": {
											Type:     schema.TypeString,
											Optional: true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func nodePoolSchema(forceNewNodes bool) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"node_os_version": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
				"num_instances": {
					Type:     schema.TypeInt,
					Required: true,
					ForceNew: forceNewNodes,
				},
				"ahv_config": {
					Type:     schema.TypeMap,
					Optional: true,
					// Computed: true,
					// ForceNew: forceNewNodes,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"cpu": {
								Type:     schema.TypeInt,
								Optional: true,
							},
							"disk_mib": {
								Type:     schema.TypeInt,
								Optional: true,
							},
							"memory_mib": {
								Type:     schema.TypeInt,
								Optional: true,
							},
							"network_uuid": {
								Type:     schema.TypeString,
								Required: true,
							},
							"prism_element_cluster_uuid": {
								Type:     schema.TypeString,
								Required: true,
							},
						},
					},
				},
				"nodes": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"hostname": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"ipv4_address": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

func resourceNutanixKarbonClusterCreate(d *schema.ResourceData, meta interface{}) error {
	log.Print("[Debug] Entering resourceNutanixKarbonClusterCreate")
	// Get client connection
	client := meta.(*Client)
	conn := client.KarbonAPI
	setTimeout(meta)
	// Prepare request
	// Node pools
	etcdNodePool, err := expandNodePool(d.Get("etcd_node_pool").([]interface{}))
	if err != nil {
		return err
	}
	workerNodePool, err := expandNodePool(d.Get("worker_node_pool").([]interface{}))
	if err != nil {
		return err
	}
	masterNodePool, err := expandNodePool(d.Get("master_node_pool").([]interface{}))
	if err != nil {
		return err
	}
	// storageclass
	storageClassConfig, err := expandStorageClassConfig(d.Get("storage_class_config").(*schema.Set).List())
	if err != nil {
		return err
	}
	// CNI
	// todo modify these unchecked GETs
	cniConfig, err := expandCNI(d.Get("cni_config").(*schema.Set).List())
	if err != nil {
		return err
	}
	karbonClusterName := d.Get("name").(string)

	karbonCluster := &karbon.ClusterIntentInput{
		Name:      karbonClusterName,
		Version:   d.Get("version").(string),
		CNIConfig: *cniConfig,
		ETCDConfig: karbon.ClusterETCDConfigIntentInput{
			NodePools: etcdNodePool,
		},
		MastersConfig: karbon.ClusterMasterConfigIntentInput{
			NodePools: masterNodePool,
		},
		Metadata: karbon.ClusterMetadataIntentInput{
			APIVersion: "2.0.0",
		},
		StorageClassConfig: *storageClassConfig,
		WorkersConfig: karbon.ClusterWorkerConfigIntentInput{
			NodePools: workerNodePool,
		},
	}
	activePassiveConfig, apcOk := d.GetOk("active_passive_config")
	externalLbConfig, elbcOk := d.GetOk("external_lb_config")
	if apcOk && elbcOk {
		return fmt.Errorf("cannot pass both active_passive_config and external_lb_config")
	}
	if !apcOk && !elbcOk {
		karbonCluster.MastersConfig.SingleMasterConfig = &karbon.ClusterSingleMasterConfigIntentInput{}
	}
	// set active passive config
	if apcOk {
		activePassiveConfigList := activePassiveConfig.(*schema.Set).List()
		karbonCluster.MastersConfig.ActivePassiveConfig = &karbon.ClusterActivePassiveMasterConfigIntentInput{
			ExternalIPv4Address: activePassiveConfigList[0].(map[string]interface{})["external_ipv4_address"].(string),
		}
		// set active active config
	}
	if elbcOk {
		externalLbConfigList := externalLbConfig.(*schema.Set).List()
		externalLbConfigElement := externalLbConfigList[0].(map[string]interface{})
		masterNodesConfig := make([]karbon.ClusterMasterNodeMasterConfigIntentInput, 0)
		if mnc, ok := externalLbConfigElement["master_nodes_config"]; ok {
			masterNodesConfigSlice := mnc.(*schema.Set).List()
			for _, mnce := range masterNodesConfigSlice {
				masterConf := karbon.ClusterMasterNodeMasterConfigIntentInput{}
				if val, ok := mnce.(map[string]interface{})["ipv4_address"]; ok {
					masterConf.IPv4Address = val.(string)
				}
				if val, ok := mnce.(map[string]interface{})["node_pool_name"]; ok {
					masterConf.NodePoolName = val.(string)
				}
				masterNodesConfig = append(masterNodesConfig, masterConf)
			}
		} else {
			return fmt.Errorf("master_nodes_config must be passed when configuring external_lb_config")
		}
		karbonCluster.MastersConfig.ExternalLBConfig = &karbon.ClusterExternalLBMasterConfigIntentInput{
			ExternalIPv4Address: externalLbConfigElement["external_ipv4_address"].(string),
			MasterNodesConfig:   masterNodesConfig,
		}
	}

	utils.PrintToJSON(karbonCluster, "[DEBUG karbonCluster: ")
	createClusterResponse, err := conn.Cluster.CreateKarbonCluster(karbonCluster)
	if err != nil {
		return fmt.Errorf("error occurred during cluster creation:\n %s", err)
	}
	utils.PrintToJSON(createClusterResponse, "[DEBUG createClusterResponse: ")
	if createClusterResponse.TaskUUID == "" {
		return fmt.Errorf("did not retrieve task uuid")
	}
	if createClusterResponse.ClusterUUID == "" {
		return fmt.Errorf("did not retrieve cluster uuid")
	}
	err = WaitForKarbonCluster(client, createClusterResponse.TaskUUID)
	if err != nil {
		return err
	}

	fmt.Printf("Cluster uuid: %s", createClusterResponse.ClusterUUID)
	fmt.Printf("Task uuid: %s", createClusterResponse.TaskUUID)
	if privateRegistries, ok := d.GetOk("private_registries"); ok {
		newPrivateRegistries, err := expandPrivateRegistries(privateRegistries.(*schema.Set).List())
		if err != nil {
			return err
		}
		utils.PrintToJSON(newPrivateRegistries, "newPrivateRegistries: ")
		for _, newP := range *newPrivateRegistries {
			log.Printf("adding private registry %s", *newP.RegistryName)
			conn.Cluster.AddPrivateRegistry(karbonClusterName, newP)
		}
	}
	// Set terraform state id
	d.SetId(createClusterResponse.ClusterUUID)
	return resourceNutanixKarbonClusterRead(d, meta)
}

func resourceNutanixKarbonClusterRead(d *schema.ResourceData, meta interface{}) error {
	log.Print("[Debug] Entering resourceNutanixKarbonClusterRead")
	// Get client connection
	conn := meta.(*Client).KarbonAPI
	setTimeout(meta)
	// Make request to the API
	var err error
	resp, err := conn.Cluster.GetKarbonCluster(d.Id())
	if err != nil {
		d.SetId("")
		return nil
	}
	karbonClusterName := *resp.Name
	flattenedEtcdNodepool, err := flattenNodePools(d, conn, "etcd_node_pool", karbonClusterName, resp.ETCDConfig.NodePools)
	if err != nil {
		return err
	}
	flattenedWorkerNodepool, err := flattenNodePools(d, conn, "worker_node_pool", karbonClusterName, resp.WorkerConfig.NodePools)
	if err != nil {
		return err
	}
	flattenedMasterNodepool, err := flattenNodePools(d, conn, "master_node_pool", karbonClusterName, resp.MasterConfig.NodePools)
	if err != nil {
		return err
	}

	utils.PrintToJSON(flattenedWorkerNodepool, "pre set flattenedWorkerNodepool: ")
	// log.Printf("d.Get(master_node_pool)")
	// log.Print(d.Get("master_node_pool").([]interface{}))
	d.Set("name", utils.StringValue(resp.Name))

	if err = d.Set("status", utils.StringValue(resp.Status)); err != nil {
		return fmt.Errorf("error setting status for Karbon Cluster %s: %s", d.Id(), err)
	}

	// Must use know version because GA API reports different version
	var versionSet string
	log.Printf("Getting existing version: %s", d.Get("version").(string))
	if version, ok := d.GetOk("version"); ok {
		versionSet = version.(string)
	} else {
		versionSet = utils.StringValue(resp.Version)
	}
	log.Printf("using version: %s", versionSet)
	if err = d.Set("version", versionSet); err != nil {
		return fmt.Errorf("error setting version for Karbon Cluster %s: %s", d.Id(), err)
	}

	d.Set("kubeapi_server_ipv4_address", utils.StringValue(resp.KubeAPIServerIPv4Address))
	d.Set("deployment_type", resp.MasterConfig.DeploymentType)
	if err = d.Set("worker_node_pool", flattenedWorkerNodepool); err != nil {
		return fmt.Errorf("error setting worker_node_pool for Karbon Cluster %s: %s", d.Id(), err)
	}
	if err = d.Set("etcd_node_pool", flattenedEtcdNodepool); err != nil {
		return fmt.Errorf("error setting etcd_node_pool for Karbon Cluster %s: %s", d.Id(), err)
	}
	if err = d.Set("master_node_pool", flattenedMasterNodepool); err != nil {
		return fmt.Errorf("error setting worker_node_pool for Karbon Cluster %s: %s", d.Id(), err)
	}
	flattenedPrivateRegistries, err := flattenPrivateRegisties(conn, karbonClusterName)
	if err != nil {
		return fmt.Errorf("error getting flat private_registries for Karbon Cluster %s: %s", d.Id(), err)
	}
	utils.PrintToJSON(flattenedPrivateRegistries, "flattenedPrivateRegistries: ")
	if err = d.Set("private_registries", flattenedPrivateRegistries); err != nil {
		return fmt.Errorf("error setting private_registries for Karbon Cluster %s: %s", d.Id(), err)
	}
	d.SetId(*resp.UUID)
	return nil
}

func resourceNutanixKarbonClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Print("[Debug] Entering resourceNutanixKarbonClusterUpdate")
	// Get client connection
	client := meta.(*Client)
	conn := client.KarbonAPI
	setTimeout(meta)
	// Make request to the API
	resp, err := conn.Cluster.GetKarbonCluster(d.Id())
	if err != nil {
		return err
	}
	karbonClusterName := *resp.Name
	if d.HasChange("private_registries") {
		_, p := d.GetChange("private_registries")
		utils.PrintToJSON(p.(*schema.Set).List(), "p private_registries: ")
		newPrivateRegistries, err := expandPrivateRegistries(p.(*schema.Set).List())
		if err != nil {
			return err
		}
		utils.PrintToJSON(newPrivateRegistries, "newPrivateRegistries: ")
		currentPrivateRegistriesList, err := conn.Cluster.ListPrivateRegistries(karbonClusterName)
		if err != nil {
			return err
		}
		utils.PrintToJSON(currentPrivateRegistriesList, "currentPrivateRegistriesList: ")
		currentPrivateRegistries := convertKarbonPrivateRegistriesIntentInputToOperations(*currentPrivateRegistriesList)
		utils.PrintToJSON(currentPrivateRegistries, "currentPrivateRegistries: ")
		toAdd := diffFlatPrivateRegistrySlices(*newPrivateRegistries, currentPrivateRegistries)
		utils.PrintToJSON(toAdd, "toAdd: ")
		for _, a := range toAdd {
			log.Printf("adding private registry %s", *a.RegistryName)
			conn.Cluster.AddPrivateRegistry(karbonClusterName, a)
		}
		toRemove := diffFlatPrivateRegistrySlices(currentPrivateRegistries, *newPrivateRegistries)
		utils.PrintToJSON(toRemove, "toRemove: ")
		for _, r := range toRemove {
			log.Printf("removing private registry %s", *r.RegistryName)
			conn.Cluster.DeletePrivateRegistry(karbonClusterName, *r.RegistryName)
		}
	}
	return resourceNutanixKarbonClusterRead(d, meta)
}

func resourceNutanixKarbonClusterDelete(d *schema.ResourceData, meta interface{}) error {
	log.Print("[Debug] Entering resourceNutanixKarbonClusterDelete")
	client := meta.(*Client)
	conn := client.KarbonAPI
	setTimeout(meta)
	karbonClusterName := d.Get("name").(string)
	log.Printf("[DEBUG] Deleting Karbon cluster: %s, %s", karbonClusterName, d.Id())

	clusterDeleteResponse, err := conn.Cluster.DeleteKarbonCluster(karbonClusterName)
	if err != nil {
		return fmt.Errorf("error while deleting Karbon Cluster UUID(%s): %s", d.Id(), err)
	}
	err = WaitForKarbonCluster(client, clusterDeleteResponse.TaskUUID)
	if err != nil {
		return fmt.Errorf("error while waiting for Karbon Cluster deletion with UUID(%s): %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}

func resourceNutanixKarbonClusterExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Print("[DEBUG] Entering resourceNutanixKarbonClusterExists")
	conn := meta.(*Client).KarbonAPI
	setTimeout(meta)
	// Make request to the API
	resp, err := conn.Cluster.GetKarbonCluster(d.Id())
	log.Print("error:")
	log.Print(err)
	utils.PrintToJSON(resp, "resourceNutanixKarbonClusterExists resp: ")
	if err != nil {
		d.SetId("")
		return false, nil
	}
	return true, nil
}

func diffFlatPrivateRegistrySlices(prSlice1 []karbon.PrivateRegistryOperationIntentInput, prSlice2 []karbon.PrivateRegistryOperationIntentInput) []karbon.PrivateRegistryOperationIntentInput {
	prSliceResult := make([]karbon.PrivateRegistryOperationIntentInput, 0)
	for _, e1 := range prSlice1 {
		found := false
		for _, e2 := range prSlice2 {
			if *e1.RegistryName == *e2.RegistryName {
				found = true
			}
		}
		if !found {
			prSliceResult = append(prSliceResult, e1)
		}
	}
	return prSliceResult
}

func convertKarbonPrivateRegistriesIntentInputToOperations(privateRegistryResponses karbon.PrivateRegistryListResponse) []karbon.PrivateRegistryOperationIntentInput {
	s := make([]karbon.PrivateRegistryOperationIntentInput, 0)
	for _, p := range privateRegistryResponses {
		s = append(s, convertKarbonPrivateRegistryIntentInputToOperation(p))
	}
	return s
}

func convertKarbonPrivateRegistryIntentInputToOperation(privateRegistryResponse karbon.PrivateRegistryResponse) karbon.PrivateRegistryOperationIntentInput {
	return karbon.PrivateRegistryOperationIntentInput{
		RegistryName: privateRegistryResponse.Name,
	}
}

func expandPrivateRegistries(privateRegistries []interface{}) (*[]karbon.PrivateRegistryOperationIntentInput, error) {
	prSlice := make([]karbon.PrivateRegistryOperationIntentInput, 0)
	for _, p := range privateRegistries {
		fp, err := expandPrivateRegistry(p.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		prSlice = append(prSlice, *fp)
	}
	return &prSlice, nil
}

func expandPrivateRegistry(privateRegistry map[string]interface{}) (*karbon.PrivateRegistryOperationIntentInput, error) {
	if rn, ok := privateRegistry["registry_name"]; ok {
		rns := rn.(string)
		return &karbon.PrivateRegistryOperationIntentInput{
			RegistryName: &rns,
		}, nil
	}
	return nil, fmt.Errorf("failed to retrieve registry_name for private registry")
}

func flattenPrivateRegisties(conn *karbon.Client, karbonClusterName string) ([]map[string]interface{}, error) {
	flatPrivReg := make([]map[string]interface{}, 0)
	privRegList, err := conn.Cluster.ListPrivateRegistries(karbonClusterName)
	utils.PrintToJSON(privRegList, "privRegList: ")
	if err != nil {
		return nil, err
	}
	for _, p := range *privRegList {
		flatPrivReg = append(flatPrivReg, map[string]interface{}{
			// "endpoint": p.Endpoint,
			"registry_name": p.Name,
			// "UUID":     p.UUID,
		})
	}
	return flatPrivReg, nil
}

func flattenNodePools(d *schema.ResourceData, conn *karbon.Client, nodePoolKey string, karbonClusterName string, nodepools []string) ([]map[string]interface{}, error) {
	flatNodepools := make([]map[string]interface{}, 0)
	// start workaround for disk_mib bug GA API
	expandedUserDefinedNodePools := make([]karbon.ClusterNodePool, 0)
	var err error
	if nodepoolInterface, ok := d.GetOk(nodePoolKey); ok {
		expandedUserDefinedNodePools, err = expandNodePool(nodepoolInterface.([]interface{}))
		if err != nil {
			return nil, fmt.Errorf("unable to expand node pool during flattening: %s", err)
		}
	}
	// end workaround for disk_mib bug GA API
	for _, np := range nodepools {
		nodepool, err := conn.Cluster.GetKarbonClusterNodePool(karbonClusterName, np)
		if err != nil {
			return nil, err
		}
		var flattenedNodepool map[string]interface{}
		if len(expandedUserDefinedNodePools) == 0 {
			flattenedNodepool = flattenNodePool(nil, nodepool)
		} else {
			for _, udnp := range expandedUserDefinedNodePools {
				expandedUserDefinedNodePool := udnp
				if *expandedUserDefinedNodePool.Name == *nodepool.Name {
					flattenedNodepool = flattenNodePool(&expandedUserDefinedNodePool, nodepool)

					break
				}
			}
		}
		flatNodepools = append(flatNodepools, flattenedNodepool)
	}
	return flatNodepools, nil
}

func flattenNodePool(userDefinedNodePools *karbon.ClusterNodePool, nodepool *karbon.ClusterNodePool) map[string]interface{} {
	flatNodepool := map[string]interface{}{}
	// Nodes
	nodes := make([]map[string]interface{}, 0)
	for _, npn := range *nodepool.Nodes {
		nodes = append(nodes, map[string]interface{}{
			"hostname":     npn.Hostname,
			"ipv4_address": npn.IPv4Address,
		})
	}
	flatNodepool["nodes"] = nodes
	// AHV config
	// disk_mib, ok := d.GetOk("etcd_node_pool")
	diskMib := strconv.FormatInt(nodepool.AHVConfig.DiskMib, 10)
	if userDefinedNodePools != nil {
		utils.PrintToJSON(userDefinedNodePools, "userDefinedNodePools: ")
		log.Print(userDefinedNodePools.AHVConfig.DiskMib)
		diskMib = strconv.FormatInt(userDefinedNodePools.AHVConfig.DiskMib, 10)
	}
	flatNodepool["ahv_config"] = map[string]interface{}{
		"cpu": strconv.FormatInt(nodepool.AHVConfig.CPU, 10),
		// karbon api bug 	GetKarbonClusterLegacy(uuid string) (*KarbonClusterLegacyIntentResponse, error)
		"disk_mib": diskMib,
		// must check with legacy nodepool because GA API reports wrong disk space
		// "disk_mib":                   strconv.FormatInt(*legacyNodepool.ResourceConfig.DiskMib, 10),
		"memory_mib":                 strconv.FormatInt(nodepool.AHVConfig.MemoryMib, 10),
		"network_uuid":               nodepool.AHVConfig.NetworkUUID,
		"prism_element_cluster_uuid": nodepool.AHVConfig.PrismElementClusterUUID,
	}
	flatNodepool["name"] = nodepool.Name
	flatNodepool["num_instances"] = nodepool.NumInstances
	flatNodepool["node_os_version"] = nodepool.NodeOSVersion
	utils.PrintToJSON(flatNodepool, "flatNodepool: ")
	return flatNodepool
}

func GetNodePoolsForCluster(conn *karbon.Client, karbonClusterName string, nodepools []string) ([]karbon.ClusterNodePool, error) {
	nodepoolStructs := make([]karbon.ClusterNodePool, 0)
	for _, np := range nodepools {
		nodepool, err := conn.Cluster.GetKarbonClusterNodePool(karbonClusterName, np)
		if err != nil {
			return nil, err
		}
		nodepoolStructs = append(nodepoolStructs, *nodepool)
	}
	return nodepoolStructs, nil
}

func WaitForKarbonCluster(client *Client, taskUUID string) error {
	log.Printf("Starting wait")
	sleepTime := 30
	var status string = "QUEUED"

	for status == "QUEUED" || status == "RUNNING" {
		time.Sleep(time.Duration(sleepTime) * time.Second)
		v, err := client.API.V3.GetTask(taskUUID)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "INVALID_UUID") {
				return fmt.Errorf("invalid uuid retrieved")
			}
			return err
		}
		status = *v.Status
		log.Printf("Status: %s", status)
		if status == "INVALID_UUID" || status == "FAILED" {
			return fmt.Errorf("error_detail: %s, progress_message: %s", utils.StringValue(v.ErrorDetail), utils.StringValue(v.ProgressMessage))
		}
	}
	if status == "SUCCEEDED" {
		return nil
	}
	return fmt.Errorf("end state was not succeeded but was %s", status)
}

func setTimeout(meta interface{}) {
	client := meta.(*Client)
	if client.WaitTimeout != 0 {
		vmTimeout = time.Duration(client.WaitTimeout) * time.Minute
	}
}

func expandStorageClassConfig(storageClassConfigsInput []interface{}) (*karbon.ClusterStorageClassConfigIntentInput, error) {
	log.Print("[DEBUG] entering expandStorageClassConfig")
	if len(storageClassConfigsInput) != 1 {
		return nil, fmt.Errorf("more than one storage class input passed")
	}
	storageClassConfigInput := storageClassConfigsInput[0].(map[string]interface{})
	storageClassConfig := &karbon.ClusterStorageClassConfigIntentInput{
		DefaultStorageClass: true,
		Name:                "default-storageclass",
		VolumesConfig:       karbon.ClusterVolumesConfigIntentInput{},
	}
	if val, ok := storageClassConfigInput["reclaim_policy"]; ok {
		storageClassConfig.ReclaimPolicy = val.(string)
	}
	if volumesConfig, ok3 := storageClassConfigInput["volumes_config"]; ok3 {
		volumesConfig := volumesConfig.(map[string]interface{})
		if valFileSystem, ok := volumesConfig["file_system"]; ok {
			storageClassConfig.VolumesConfig.FileSystem = valFileSystem.(string)
		}
		if valFlashMode, ok := volumesConfig["flash_mode"]; ok {
			b, _ := strconv.ParseBool(valFlashMode.(string))
			storageClassConfig.VolumesConfig.FlashMode = b
		}
		if valPassword, ok := volumesConfig["password"]; ok {
			storageClassConfig.VolumesConfig.Password = valPassword.(string)
		}
		if valPrismElementClusterUUID, ok := volumesConfig["prism_element_cluster_uuid"]; ok {
			storageClassConfig.VolumesConfig.PrismElementClusterUUID = valPrismElementClusterUUID.(string)
		}
		if valStorageContainer, ok := volumesConfig["storage_container"]; ok {
			storageClassConfig.VolumesConfig.StorageContainer = valStorageContainer.(string)
		}
		if valUsername, ok := volumesConfig["username"]; ok {
			storageClassConfig.VolumesConfig.Username = valUsername.(string)
		}
	}
	return storageClassConfig, nil
}

func expandCNI(cniConfigInput []interface{}) (*karbon.ClusterCNIConfigIntentInput, error) {
	if len(cniConfigInput) != 1 {
		return nil, fmt.Errorf("cannot have more than one CNI configuration")
	}
	cniConfig := &karbon.ClusterCNIConfigIntentInput{}
	cniConfigMap := cniConfigInput[0].(map[string]interface{})
	if value, ok := cniConfigMap["node_cidr_mask_size"]; ok {
		cniConfig.NodeCIDRMaskSize = int64(value.(int))
	}
	if value, ok := cniConfigMap["pod_ipv4_cidr"]; ok && value.(string) != "" {
		cniConfig.PodIPv4CIDR = value.(string)
	}
	if value, ok := cniConfigMap["service_ipv4_cidr"]; ok && value.(string) != "" {
		cniConfig.ServiceIPv4CIDR = value.(string)
	}
	// todo ugly code
	if calicoConfig, cok := cniConfigMap["calico_config"]; cok && len(calicoConfig.(*schema.Set).List()) > 0 {
		utils.PrintToJSON(calicoConfig, "calicoConfig: ")
		if flannelConfig, fok := cniConfigMap["flannel_config"]; fok && len(flannelConfig.(*schema.Set).List()) > 0 {
			utils.PrintToJSON(flannelConfig, "flannelConfig: ")
			return nil, fmt.Errorf("cannot have both calico and flannel config")
		}
		calicoConfigMap := calicoConfig.(*schema.Set).List()[0].(map[string]interface{})
		ipPoolConfigs := make([]karbon.ClusterCalicoConfigIPPoolConfigIntentInput, 0)
		for _, ipc := range calicoConfigMap["ip_pool_configs"].([]interface{}) {
			mipc := ipc.(map[string]interface{})
			ipPoolConfigs = append(ipPoolConfigs, karbon.ClusterCalicoConfigIPPoolConfigIntentInput{
				CIDR: mipc["cidr"].(string),
			})
		}
		cniConfig.CalicoConfig = &karbon.ClusterCalicoConfigIntentInput{
			IPPoolConfigs: ipPoolConfigs,
		}
	} else {
		cniConfig.FlannelConfig = &karbon.ClusterFlannelConfigIntentInput{}
	}
	utils.PrintToJSON(cniConfig, "cniConfig: ")
	return cniConfig, nil
}

func expandNodePool(nodepoolsInput []interface{}) ([]karbon.ClusterNodePool, error) {
	nodepools := make([]karbon.ClusterNodePool, 0)
	for _, npi := range nodepoolsInput {
		nodepoolInput := npi.(map[string]interface{})
		nodepool := &karbon.ClusterNodePool{
			AHVConfig: &karbon.ClusterNodePoolAHVConfig{},
		}
		if nameVal, nameOk := nodepoolInput["name"]; nameOk && nameVal.(string) != "" {
			npName := nameVal.(string)
			nodepool.Name = &npName
		} else {
			return nil, fmt.Errorf("nodepool name must be passed")
		}
		if val, ok := nodepoolInput["node_os_version"]; ok {
			nodeOsVersion := val.(string)
			nodepool.NodeOSVersion = &nodeOsVersion
		}
		if val2, ok2 := nodepoolInput["num_instances"]; ok2 {
			numInstances := int64(val2.(int))
			nodepool.NumInstances = &numInstances
		}
		if ahvConfig, ok3 := nodepoolInput["ahv_config"]; ok3 {
			ahvConfig := ahvConfig.(map[string]interface{})
			if valCPU, ok := ahvConfig["cpu"]; ok {
				i, _ := strconv.ParseInt(valCPU.(string), 10, 64)
				// Karbon CPU workaround
				modi := i % cpuDivisionAmount
				if modi != 0 {
					return nil, fmt.Errorf("amount of CPU must be an even number")
				}
				divi := i / cpuDivisionAmount
				nodepool.AHVConfig.CPU = divi
			}
			if valDiskMib, ok := ahvConfig["disk_mib"]; ok {
				log.Print("[DEBUG] valDiskMib")
				log.Print(valDiskMib)
				i, _ := strconv.ParseInt(valDiskMib.(string), 10, 64)
				log.Print(i)
				nodepool.AHVConfig.DiskMib = i
			}
			if valMemoryMib, ok := ahvConfig["memory_mib"]; ok {
				log.Print("[DEBUG] valMemoryMib")
				log.Print(valMemoryMib)
				i, _ := strconv.ParseInt(valMemoryMib.(string), 10, 64)
				log.Print(i)
				nodepool.AHVConfig.MemoryMib = i
			}
			if valNetworkUUID, ok := ahvConfig["network_uuid"]; ok {
				nodepool.AHVConfig.NetworkUUID = valNetworkUUID.(string)
			}
			if valPrismElementClusterUUID, ok := ahvConfig["prism_element_cluster_uuid"]; ok {
				nodepool.AHVConfig.PrismElementClusterUUID = valPrismElementClusterUUID.(string)
			}
		}
		if nodes, ok4 := nodepoolInput["nodes"]; ok4 {
			nodesSlice := make([]karbon.ClusterNodeIntentResponse, 0)
			for _, n := range nodes.([]interface{}) {
				nmap := n.(map[string]interface{})
				node := karbon.ClusterNodeIntentResponse{}
				if nHostname, ok := nmap["hostname"]; ok && nHostname != "" {
					nh := nHostname.(string)
					node.Hostname = &nh
				}
				if nIP, ok := nmap["ipv4_address"]; ok && nIP != "" {
					ni := nIP.(string)
					node.IPv4Address = &ni
				}
				nodesSlice = append(nodesSlice, node)
			}
			nodepool.Nodes = &nodesSlice
		}
		nodepools = append(nodepools, *nodepool)
	}
	return nodepools, nil
}
