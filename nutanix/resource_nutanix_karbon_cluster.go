package nutanix

import (
	"fmt"
	"log"
	"strings"
	"time"

	karbon "github.com/terraform-providers/terraform-provider-nutanix/client/karbon"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

const (
	DEFAULTMASTERNODEPOOLNAME = "master_node_pool"
	DEFAULTETCDNODEPOOLNAME   = "etcd_node_pool"
	DEFAULTWORKERNODEPOOLNAME = "worker_node_pool"
	DEFAULTPODIPV4CIDR        = "172.20.0.0/16"
	DEFAULTSERVICEIPV4CIDR    = "172.19.0.0/16"
	DEFAULTRECLAIMPOLICY      = "Delete"
	DEFAULTFILESYSTEM         = "ext4"
	DEFAULTSTORAGECLASSNAME   = "default-storageclass"
	DEFAULTNODECIDRMASKSIZE   = 24
	DEFAULTETCDNODECPU        = 4
	DEFAULTETCDNODEDISKMIB    = 40960
	DEFAULTETCDNODEEMORYMIB   = 8192
	DEFAULTMASTERNODECPU      = 2
	DEFAULTMASTERNODEDISKMIB  = 122880
	DEFAULTMASTERNODEEMORYMIB = 4096
	DEFAULTWORKERNODECPU      = 8
	DEFAULTWORKERNODEDISKMIB  = 122880
	DEFAULTWORKERNODEEMORYMIB = 8192
	MINDISKMIB                = 1024
	MINMEMORYMIB              = 1024
	MINCPU                    = 2
	MINNUMINSTANCES           = 1
	MAXMASTERNODES            = 5
	MINMASTERNODES            = 2
	CPUDIVISIONAMOUNT         = 2
	KARBONAPIVERSION          = "2.0.0"
	MINIMUMWAITTIMEOUT        = 1
	DEFAULTWAITTIMEOUT        = 60
	WAITDELAY                 = 10 * time.Second
	WAITMINTIMEOUT            = 10 * time.Second
)

// Known issues:
//  - Importing karbon clusters do not contain cni_configs and storage_class_configs
//  - Importing karbon clusters show an incorrect version

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
		"wait_timeout_minutes": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      DEFAULTWAITTIMEOUT,
			ValidateFunc: validation.IntAtLeast(MINIMUMWAITTIMEOUT),
		},
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
			ForceNew: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Optional: true,
						Default:  DEFAULTSTORAGECLASSNAME,
						ForceNew: true,
					},
					"reclaim_policy": {
						Type:         schema.TypeString,
						Optional:     true,
						Default:      DEFAULTRECLAIMPOLICY,
						ValidateFunc: validation.StringInSlice(getSupportedReclaimPolicies(), false),
					},
					"volumes_config": {
						Type:     schema.TypeList,
						Required: true,
						MaxItems: 1,
						ForceNew: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"file_system": {
									Type:         schema.TypeString,
									ValidateFunc: validation.StringInSlice(getSupportedFileSystems(), false),
									Optional:     true,
									Default:      DEFAULTFILESYSTEM,
								},
								"flash_mode": {
									Type:     schema.TypeBool,
									Optional: true,
									Default:  false,
								},
								"password": {
									Type:        schema.TypeString,
									Required:    true,
									Sensitive:   true,
									DefaultFunc: schema.EnvDefaultFunc("NUTANIX_PE_PASSWORD", nil),
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
									Type:        schema.TypeString,
									Required:    true,
									DefaultFunc: schema.EnvDefaultFunc("NUTANIX_PE_USERNAME", nil),
								},
							},
						},
					},
				},
			},
		},
		"single_master_config": {
			Type:          schema.TypeList,
			Optional:      true,
			ForceNew:      true,
			MaxItems:      1,
			ConflictsWith: []string{"external_lb_config", "active_passive_config"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{},
			},
		},
		"active_passive_config": {
			Type:          schema.TypeList,
			Optional:      true,
			ForceNew:      true,
			MaxItems:      1,
			ConflictsWith: []string{"external_lb_config", "single_master_config"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"external_ipv4_address": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.SingleIP(),
					},
				},
			},
		},
		"external_lb_config": {
			Type:          schema.TypeList,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"active_passive_config", "single_master_config"},
			MaxItems:      1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"external_ipv4_address": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.SingleIP(),
					},
					"master_nodes_config": {
						Type:     schema.TypeSet,
						Required: true,
						MaxItems: MAXMASTERNODES,
						MinItems: MINMASTERNODES,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"ipv4_address": {
									Type:         schema.TypeString,
									Required:     true,
									ValidateFunc: validation.SingleIP(),
								},
								"node_pool_name": {
									Type:     schema.TypeString,
									Optional: true,
									Default:  DEFAULTMASTERNODEPOOLNAME,
								},
							},
						},
					},
				},
			},
		},
		"private_registry": {
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
		"etcd_node_pool":   nodePoolSchema(DEFAULTETCDNODEPOOLNAME, true, DEFAULTETCDNODECPU, DEFAULTETCDNODEDISKMIB, DEFAULTETCDNODEEMORYMIB),
		"master_node_pool": nodePoolSchema(DEFAULTMASTERNODEPOOLNAME, true, DEFAULTMASTERNODECPU, DEFAULTMASTERNODEDISKMIB, DEFAULTMASTERNODEEMORYMIB),
		"worker_node_pool": nodePoolSchema(DEFAULTWORKERNODEPOOLNAME, false, DEFAULTWORKERNODECPU, DEFAULTWORKERNODEDISKMIB, DEFAULTWORKERNODEEMORYMIB),
		"cni_config":       CNISchema(),
	}
}

func CNISchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		ForceNew: true,
		// Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"node_cidr_mask_size": {
					Type: schema.TypeInt,
					// Required: true,
					Optional: true,
					Default:  DEFAULTNODECIDRMASKSIZE,
				},
				"pod_ipv4_cidr": {
					Type: schema.TypeString,
					// Required: true,
					Optional:     true,
					Default:      DEFAULTPODIPV4CIDR,
					ValidateFunc: validation.CIDRNetwork(0, 32),
				},
				"service_ipv4_cidr": {
					Type: schema.TypeString,
					// Required: true,
					Optional:     true,
					Default:      DEFAULTSERVICEIPV4CIDR,
					ValidateFunc: validation.CIDRNetwork(0, 32),
				},
				"flannel_config": {
					Type:     schema.TypeList,
					Optional: true,
					ForceNew: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{},
					},
				},
				"calico_config": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					ForceNew: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"ip_pool_config": {
								Type:     schema.TypeList,
								Optional: true,
								ForceNew: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"cidr": {
											Type:         schema.TypeString,
											Optional:     true,
											ForceNew:     true,
											Default:      DEFAULTPODIPV4CIDR,
											ValidateFunc: validation.CIDRNetwork(0, 32),
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

func nodePoolSchema(defaultNodepoolName string, forceNewNodes bool, cpuDefault int, diskMibDefault int, memoryMibDefault int) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  defaultNodepoolName,
					ForceNew: true,
				},
				"node_os_version": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
				"num_instances": {
					Type:         schema.TypeInt,
					Required:     true,
					ForceNew:     forceNewNodes,
					ValidateFunc: validation.IntAtLeast(MINNUMINSTANCES),
				},
				"ahv_config": {
					Type:     schema.TypeList,
					MaxItems: 1,
					Optional: true,
					// Computed: true,
					ForceNew: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"cpu": {
								Type:         schema.TypeInt,
								Optional:     true,
								Default:      cpuDefault,
								ValidateFunc: validation.IntAtLeast(MINCPU),
							},
							"disk_mib": {
								Type:         schema.TypeInt,
								Optional:     true,
								Default:      diskMibDefault,
								ValidateFunc: validation.IntAtLeast(MINDISKMIB),
							},
							"memory_mib": {
								Type:         schema.TypeInt,
								Optional:     true,
								Default:      memoryMibDefault,
								ValidateFunc: validation.IntAtLeast(MINMEMORYMIB),
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
	// Node pools
	var err error
	etcdNodePoolInput, okETCD := d.GetOk("etcd_node_pool")
	if !okETCD {
		return fmt.Errorf("unable to retrieve mandatory parameter etcd_node_pool")
	}
	workerNodePoolInput, okWorker := d.GetOk("worker_node_pool")
	if !okWorker {
		return fmt.Errorf("unable to retrieve mandatory parameter worker_node_pool")
	}
	masterNodePoolInput, okMaster := d.GetOk("master_node_pool")
	if !okMaster {
		return fmt.Errorf("unable to retrieve mandatory parameter master_node_pool")
	}
	storageClassConfigInput, okStorageClassConfig := d.GetOk("storage_class_config")
	if !okStorageClassConfig {
		return fmt.Errorf("unable to retrieve mandatory parameter storage_class_config")
	}
	cniInput, okCNI := d.GetOk("cni_config")
	if !okCNI {
		return fmt.Errorf("unable to retrieve mandatory parameter cni_config")
	}
	karbonClusterNameInput, okName := d.GetOk("name")
	if !okName {
		return fmt.Errorf("unable to retrieve mandatory parameter name")
	}
	versionInput, okVersion := d.GetOk("version")
	if !okVersion {
		return fmt.Errorf("unable to retrieve mandatory parameter version")
	}
	timeout, timeoutErr := getTimeout(d)
	if timeoutErr != nil {
		return timeoutErr
	}

	karbonClusterName := karbonClusterNameInput.(string)

	etcdNodePool, err := expandNodePool(etcdNodePoolInput.([]interface{}))
	if err != nil {
		return err
	}
	workerNodePool, err := expandNodePool(workerNodePoolInput.([]interface{}))
	if err != nil {
		return err
	}
	masterNodePool, err := expandNodePool(masterNodePoolInput.([]interface{}))
	if err != nil {
		return err
	}
	// storageclass
	storageClassConfig, err := expandStorageClassConfig(storageClassConfigInput.(*schema.Set).List())
	if err != nil {
		return err
	}
	// CNI
	cniConfig, err := expandCNI(cniInput.([]interface{}))
	if err != nil {
		return err
	}
	karbonCluster := &karbon.ClusterIntentInput{
		Name:      karbonClusterName,
		Version:   versionInput.(string),
		CNIConfig: *cniConfig,
		ETCDConfig: karbon.ClusterETCDConfigIntentInput{
			NodePools: etcdNodePool,
		},
		MastersConfig: karbon.ClusterMasterConfigIntentInput{
			NodePools: masterNodePool,
		},
		Metadata: karbon.ClusterMetadataIntentInput{
			APIVersion: KARBONAPIVERSION,
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
		err = addActivePassiveConfig(activePassiveConfig, karbonCluster)
		if err != nil {
			return err
		}
	}
	if elbcOk {
		// set active active config
		err = addExternalLBConfig(externalLbConfig, karbonCluster)
		if err != nil {
			return err
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
	err = WaitForKarbonCluster(client, timeout, createClusterResponse.TaskUUID)
	if err != nil {
		return err
	}

	fmt.Printf("Cluster uuid: %s", createClusterResponse.ClusterUUID)
	fmt.Printf("Task uuid: %s", createClusterResponse.TaskUUID)
	if privateRegistries, ok := d.GetOk("private_registry"); ok {
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
	if err = d.Set("kubeapi_server_ipv4_address", utils.StringValue(resp.KubeAPIServerIPv4Address)); err != nil {
		return fmt.Errorf("error setting kubeapi_server_ipv4_address for Karbon Cluster %s: %s", d.Id(), err)
	}
	if err = d.Set("deployment_type", resp.MasterConfig.DeploymentType); err != nil {
		return fmt.Errorf("error setting deployment_type for Karbon Cluster %s: %s", d.Id(), err)
	}
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
		return fmt.Errorf("error getting flat private_registry for Karbon Cluster %s: %s", d.Id(), err)
	}
	// utils.PrintToJSON(flattenedPrivateRegistries, "flattenedPrivateRegistries: ")
	if err = d.Set("private_registry", flattenedPrivateRegistries); err != nil {
		return fmt.Errorf("error setting private_registry for Karbon Cluster %s: %s", d.Id(), err)
	}
	flatCNIConfig := flattenCNIConfig(resp.CNIConfig)
	utils.PrintToJSON(flatCNIConfig, "flatCNIConfig: ")
	if err = d.Set("cni_config", flatCNIConfig); err != nil {
		return fmt.Errorf("error setting cni_config for Karbon Cluster %s: %s", d.Id(), err)
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
	if d.HasChange("worker_node_pool") {
		timeout, timeoutErr := getTimeout(d)
		if timeoutErr != nil {
			return timeoutErr
		}
		log.Printf("Change!")
		_, n := d.GetChange("worker_node_pool")
		newWorkerNodePool, err := expandNodePool(n.([]interface{}))
		if err != nil {
			return fmt.Errorf("Error occured while expanding new worker node pool: %s", err)
		}
		utils.PrintToJSON(newWorkerNodePool, "new_worker_node_pool: ")
		currentNodePool, err := GetNodePoolsForCluster(conn, karbonClusterName, resp.WorkerConfig.NodePools)
		if err != nil {
			return err
		}
		utils.PrintToJSON(currentNodePool, "current_node_pool: ")
		taskUUID, err := determineNodepoolsScaling(client, karbonClusterName, currentNodePool, newWorkerNodePool)
		if err != nil {
			return err
		}
		err = WaitForKarbonCluster(client, timeout, taskUUID)
	}
	if d.HasChange("private_registry") {
		_, p := d.GetChange("private_registry")
		// utils.PrintToJSON(p.(*schema.Set).List(), "p private_registry: ")
		newPrivateRegistries, err := expandPrivateRegistries(p.(*schema.Set).List())
		if err != nil {
			return err
		}
		// utils.PrintToJSON(newPrivateRegistries, "newPrivateRegistries: ")
		currentPrivateRegistriesList, err := conn.Cluster.ListPrivateRegistries(karbonClusterName)
		if err != nil {
			return err
		}
		// utils.PrintToJSON(currentPrivateRegistriesList, "currentPrivateRegistriesList: ")
		currentPrivateRegistries := convertKarbonPrivateRegistriesIntentInputToOperations(*currentPrivateRegistriesList)
		// utils.PrintToJSON(currentPrivateRegistries, "currentPrivateRegistries: ")
		toAdd := diffFlatPrivateRegistrySlices(*newPrivateRegistries, currentPrivateRegistries)
		// utils.PrintToJSON(toAdd, "toAdd: ")
		for _, a := range toAdd {
			log.Printf("adding private registry %s", *a.RegistryName)
			conn.Cluster.AddPrivateRegistry(karbonClusterName, a)
		}
		toRemove := diffFlatPrivateRegistrySlices(currentPrivateRegistries, *newPrivateRegistries)
		// utils.PrintToJSON(toRemove, "toRemove: ")
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
	timeout, timeoutErr := getTimeout(d)
	if timeoutErr != nil {
		return timeoutErr
	}
	karbonClusterNameInput, okName := d.GetOk("name")
	if !okName {
		return fmt.Errorf("unable to retrieve mandatory parameter name")
	}
	karbonClusterName := karbonClusterNameInput.(string)
	log.Printf("[DEBUG] Deleting Karbon cluster: %s, %s", karbonClusterName, d.Id())

	clusterDeleteResponse, err := conn.Cluster.DeleteKarbonCluster(karbonClusterName)
	if err != nil {
		return fmt.Errorf("error while deleting Karbon Cluster UUID(%s): %s", d.Id(), err)
	}
	err = WaitForKarbonCluster(client, timeout, clusterDeleteResponse.TaskUUID)
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
	karbonClusterName, ok := d.GetOk("name")
	var exists bool
	var err error
	// search by Name
	if ok {
		exists, err = checkNutanixKarbonClusterExistsByName(conn, karbonClusterName.(string))
	} else {
		//search by uuid
		exists, err = checkNutanixKarbonClusterExistsByUUID(conn, d.Id())
	}
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "cluster not found") {
			d.SetId("")
			return exists, nil
		}
		return exists, fmt.Errorf("error checking kubernetes cluster %s existence: %s", d.Id(), err)
	}
	return exists, nil
}

func getTimeout(d *schema.ResourceData) (int64, error) {
	timeoutInput, okTimeout := d.GetOk("wait_timeout_minutes")
	if !okTimeout {
		return 0, fmt.Errorf("unable to retrieve mandatory parameter wait_timeout_minutes")
	}
	return int64(timeoutInput.(int)), nil
}

func checkNutanixKarbonClusterExistsByUUID(conn *karbon.Client, uuid string) (bool, error) {
	// Make request to the API
	log.Print("[Debug] checkNutanixKarbonClusterExistsByUUID UUID:")
	log.Print(uuid)
	karbonClusters, err := conn.Cluster.ListKarbonClusters()
	utils.PrintToJSON(karbonClusters, "checkNutanixKarbonClusterExistsByUUID karbonClusters: ")
	log.Print("error:")
	log.Print(err)
	if err != nil {
		log.Print("checkNutanixKarbonClusterExistsByUUID returning false nil")
		return false, err
	}
	for _, k := range *karbonClusters {
		if *k.UUID == uuid {
			log.Print("checkNutanixKarbonClusterExistsByUUID returning true nil")
			return true, nil
		}
	}
	log.Print("checkNutanixKarbonClusterExistsByUUID returning false nil")
	return false, fmt.Errorf("k8s cluster not found")
}

//"cluster not found"
func checkNutanixKarbonClusterExistsByName(conn *karbon.Client, clusterName string) (bool, error) {
	// Make request to the API
	log.Print("[Debug] checkNutanixKarbonClusterExistsByName clusterName:")
	log.Print(clusterName)
	resp, err := conn.Cluster.GetKarbonCluster(clusterName)
	utils.PrintToJSON(resp, "checkNutanixKarbonClusterExistsByName resp: ")
	log.Print("error:")
	log.Print(err)
	if err != nil {
		log.Print("checkNutanixKarbonClusterExistsByName returning false nil")
		return false, err
	}
	log.Print("checkNutanixKarbonClusterExistsByName returning true nil")
	return true, nil
}

func addActivePassiveConfig(activePassiveConfig interface{}, karbonCluster *karbon.ClusterIntentInput) error {
	activePassiveConfigList := activePassiveConfig.([]interface{})
	if len(activePassiveConfigList) != 1 {
		return fmt.Errorf("cannot have more (or less) than one active_passive_config element")
	}
	externalIPV4Address, okExtAddr := activePassiveConfigList[0].(map[string]interface{})["external_ipv4_address"]
	if !okExtAddr {
		return fmt.Errorf("must set external_ipv4_address when using active_passive_config")
	}
	karbonCluster.MastersConfig.ActivePassiveConfig = &karbon.ClusterActivePassiveMasterConfigIntentInput{
		ExternalIPv4Address: externalIPV4Address.(string),
	}
	return nil
}

func addExternalLBConfig(externalLbConfig interface{}, karbonCluster *karbon.ClusterIntentInput) error {
	externalLbConfigList := externalLbConfig.([]interface{})
	if len(externalLbConfigList) != 1 {
		return fmt.Errorf("cannot have more (or less) than one external_lb_config element")
	}
	externalLbConfigElement := externalLbConfigList[0].(map[string]interface{})
	masterNodesConfig := make([]karbon.ClusterMasterNodeMasterConfigIntentInput, 0)
	if mnc, ok := externalLbConfigElement["master_nodes_config"]; ok && len(mnc.(*schema.Set).List()) > 0 {
		masterNodesConfigSlice := mnc.(*schema.Set).List()
		for _, mnce := range masterNodesConfigSlice {
			masterConf := karbon.ClusterMasterNodeMasterConfigIntentInput{}
			if val, ok := mnce.(map[string]interface{})["ipv4_address"]; ok {
				masterConf.IPv4Address = val.(string)
			} else {
				return fmt.Errorf("ipv4_address must be set when defining a master node in a external_lb_config element")
			}
			if val, ok := mnce.(map[string]interface{})["node_pool_name"]; ok {
				masterConf.NodePoolName = val.(string)
			} else {
				return fmt.Errorf("node_pool_name must be set when defining a master node in a external_lb_config element")
			}
			masterNodesConfig = append(masterNodesConfig, masterConf)
		}
	} else {
		return fmt.Errorf("master_nodes_config (>0) must be passed when configuring external_lb_config")
	}
	karbonCluster.MastersConfig.ExternalLBConfig = &karbon.ClusterExternalLBMasterConfigIntentInput{
		ExternalIPv4Address: externalLbConfigElement["external_ipv4_address"].(string),
		MasterNodesConfig:   masterNodesConfig,
	}
	return nil
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
			"registry_name": p.Name,
		})
	}
	return flatPrivReg, nil
}

func flattenCNIConfig(cniConfig karbon.ClusterCNIConfig) []map[string]interface{} {
	flatCNIConfigList := make([]map[string]interface{}, 0)
	flatCNIConfig := map[string]interface{}{
		"flannel_config":      cniConfig.FlannelConfig,
		"node_cidr_mask_size": cniConfig.NodeCIDRMaskSize,
		"pod_ipv4_cidr":       cniConfig.PodIPv4CIDR,
		"service_ipv4_cidr":   cniConfig.ServiceIPv4CIDR,
	}
	if cniConfig.CalicoConfig != nil {
		flatCNIConfig["calico_config"] = flattenCalicoConfig(cniConfig.CalicoConfig)
	}
	flatCNIConfigList = append(flatCNIConfigList, flatCNIConfig)
	return flatCNIConfigList
}

func flattenCalicoConfig(calicoConfig *karbon.ClusterCalicoConfig) []map[string]interface{} {
	flatCalicoConfigList := make([]map[string]interface{}, 0)
	if calicoConfig != nil {
		ipPoolList := make([]map[string]interface{}, 0)
		for _, ippc := range calicoConfig.IPPoolConfigs {
			ipPoolList = append(ipPoolList, map[string]interface{}{
				"cidr": ippc.CIDR,
			})
		}
		if len(ipPoolList) > 0 {
			flatCalicoConfigList = append(flatCalicoConfigList,
				map[string]interface{}{
					"ip_pool_config": ipPoolList,
				})
		}
	}
	return flatCalicoConfigList
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
		// } else {
		// 	return nil, fmt.Errorf("unable to get mandatory attribute %s", nodePoolKey)
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
	//API bug karbon
	diskMib := nodepool.AHVConfig.DiskMib
	networkUUID := nodepool.AHVConfig.NetworkUUID
	if userDefinedNodePools != nil {
		utils.PrintToJSON(userDefinedNodePools.AHVConfig, "userDefinedNodePools.AHVConfig: ")
		diskMib = userDefinedNodePools.AHVConfig.DiskMib
		log.Printf("using modified networkUUID %s", networkUUID)
		networkUUID = userDefinedNodePools.AHVConfig.NetworkUUID
	}
	flatNodepool["ahv_config"] = []map[string]interface{}{
		{
			"cpu": nodepool.AHVConfig.CPU,
			// karbon api bug 	GetKarbonClusterLegacy(uuid string) (*KarbonClusterLegacyIntentResponse, error)
			"disk_mib": diskMib,
			// "disk_mib":   nodepool.AHVConfig.DiskMib,
			"memory_mib": nodepool.AHVConfig.MemoryMib,
			//karbon api bug => network_uuid not set KRBN-3520
			// "network_uuid":               nodepool.AHVConfig.NetworkUUID,
			"network_uuid":               networkUUID,
			"prism_element_cluster_uuid": nodepool.AHVConfig.PrismElementClusterUUID,
		},
	}
	flatNodepool["name"] = nodepool.Name
	flatNodepool["num_instances"] = nodepool.NumInstances
	flatNodepool["node_os_version"] = nodepool.NodeOSVersion
	// utils.PrintToJSON(flatNodepool, "flatNodepool: ")
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

func WaitForKarbonCluster(client *Client, waitTimeoutMinutes int64, taskUUID string) error {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(client.API, taskUUID),
		Timeout:    time.Duration(waitTimeoutMinutes) * time.Minute,
		Delay:      WAITDELAY,
		MinTimeout: WAITMINTIMEOUT,
	}

	if _, errWaitTask := stateConf.WaitForState(); errWaitTask != nil {
		return fmt.Errorf("error waiting for karbon cluster to create: %s", errWaitTask)
	}
	return nil
	// log.Printf("Starting wait")
	// sleepTime := 30
	// var status string = "QUEUED"

	// for status == "QUEUED" || status == "RUNNING" {
	// 	time.Sleep(time.Duration(sleepTime) * time.Second)
	// 	v, err := client.API.V3.GetTask(taskUUID)

	// 	if err != nil {
	// 		if strings.Contains(fmt.Sprint(err), "INVALID_UUID") {
	// 			return fmt.Errorf("invalid uuid retrieved")
	// 		}
	// 		return err
	// 	}
	// 	status = *v.Status
	// 	log.Printf("Status: %s", status)
	// 	if status == "INVALID_UUID" || status == "FAILED" {
	// 		return fmt.Errorf("error_detail: %s, progress_message: %s", utils.StringValue(v.ErrorDetail), utils.StringValue(v.ProgressMessage))
	// 	}
	// }
	// if status == "SUCCEEDED" {
	// 	return nil
	// }
	// return fmt.Errorf("end state was not succeeded but was %s", status)
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
		// Name:                "default-storageclass",
		VolumesConfig: karbon.ClusterVolumesConfigIntentInput{},
	}
	if valName, okName := storageClassConfigInput["name"]; okName {
		storageClassConfig.Name = valName.(string)
	} else {
		return nil, fmt.Errorf("storage_class_config name was not set")
	}
	if val, ok := storageClassConfigInput["reclaim_policy"]; ok {
		storageClassConfig.ReclaimPolicy = val.(string)
	}
	if volumesConfigListRaw, ok3 := storageClassConfigInput["volumes_config"]; ok3 {
		volumesConfigList := volumesConfigListRaw.([]interface{})
		if len(volumesConfigList) != 1 {
			return nil, fmt.Errorf("at least one volume_config must be passed")
		}
		volumesConfig := volumesConfigList[0].(map[string]interface{})
		if valFileSystem, ok := volumesConfig["file_system"]; ok {
			storageClassConfig.VolumesConfig.FileSystem = valFileSystem.(string)
		}
		if valFlashMode, ok := volumesConfig["flash_mode"]; ok {
			storageClassConfig.VolumesConfig.FlashMode = valFlashMode.(bool)
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
		return nil, fmt.Errorf("cannot have more (or less) than one CNI configuration")
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
	if calicoConfig, cok := cniConfigMap["calico_config"]; cok && len(calicoConfig.([]interface{})) > 0 {
		if flannelConfig, fok := cniConfigMap["flannel_config"]; fok && len(flannelConfig.([]interface{})) > 0 {
			return nil, fmt.Errorf("cannot have both calico and flannel config")
		}
		calicoConfigList := calicoConfig.([]interface{})
		if len(calicoConfigList) != 1 {
			return nil, fmt.Errorf("cannot have more (or less) than one calico configuration")
		}
		calicoConfigMap := calicoConfigList[0].(map[string]interface{})
		ipPoolConfigs := make([]karbon.ClusterCalicoConfigIPPoolConfig, 0)
		var ipPoolConfigsFromMap []interface{}
		if ipcfm, ok := calicoConfigMap["ip_pool_configs"]; ok {
			ipPoolConfigsFromMap = ipcfm.([]interface{})
		}
		if ipcfm, ok := calicoConfigMap["ip_pool_config"]; ok {
			ipPoolConfigsFromMap = ipcfm.([]interface{})
		}

		for _, ipc := range ipPoolConfigsFromMap {
			mipc := ipc.(map[string]interface{})
			ipPoolConfigs = append(ipPoolConfigs, karbon.ClusterCalicoConfigIPPoolConfig{
				CIDR: mipc["cidr"].(string),
			})
		}
		cniConfig.CalicoConfig = &karbon.ClusterCalicoConfig{
			IPPoolConfigs: ipPoolConfigs,
		}
	} else {
		cniConfig.FlannelConfig = &karbon.ClusterFlannelConfig{}
	}
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
		if ahvConfigListRaw, ok3 := nodepoolInput["ahv_config"]; ok3 {
			ahvConfigList := ahvConfigListRaw.([]interface{})
			if len(ahvConfigList) != 1 {
				return nil, fmt.Errorf("ahv_config must have 1 element")
			}
			ahvConfig := ahvConfigList[0].(map[string]interface{})
			if valCPU, ok := ahvConfig["cpu"]; ok {
				i := int64(valCPU.(int))
				// Karbon CPU workaround
				modi := i % CPUDIVISIONAMOUNT
				if modi != 0 {
					return nil, fmt.Errorf("amount of CPU must be an even number")
				}
				divi := i / CPUDIVISIONAMOUNT
				nodepool.AHVConfig.CPU = divi
			}
			if valDiskMib, ok := ahvConfig["disk_mib"]; ok {
				// log.Print("[DEBUG] valDiskMib")
				// log.Print(valDiskMib)
				i := int64(valDiskMib.(int))
				// log.Print(i)
				nodepool.AHVConfig.DiskMib = i
			}
			if valMemoryMib, ok := ahvConfig["memory_mib"]; ok {
				// log.Print("[DEBUG] valMemoryMib")
				// log.Print(valMemoryMib)
				// i, _ := strconv.ParseInt(valMemoryMib.(string), 10, 64)
				i := int64(valMemoryMib.(int))
				// log.Print(i)
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

func determineNodepoolsScaling(client *Client, karbonClusterName string, currentNodepools []karbon.ClusterNodePool, newNodepools []karbon.ClusterNodePool) (string, error) {
	log.Printf("[DEBUG] entering determineNodepoolsScaling")
	var taskUUID string
	for _, cnp := range currentNodepools {
		log.Printf("cnp.Name: %s", *cnp.Name)
		for _, nnp := range newNodepools {
			log.Printf("nnp.Name: %s", *nnp.Name)
			if *cnp.Name == *nnp.Name {
				log.Print("cnp.Name == nnp.Name")
				if *cnp.NumInstances < *nnp.NumInstances {
					// scale up
					log.Print("scale up")
					amountOfNodes := *nnp.NumInstances - *cnp.NumInstances
					scaleUpRequest := &karbon.ClusterScaleUpIntentInput{
						Count: amountOfNodes,
					}
					// taskUUID, err = scaleUpNodepool(client, karbonClusterName, nnp, amountOfNodes)
					karbonClusterActionResponse, err := client.KarbonAPI.Cluster.ScaleUpKarbonCluster(
						karbonClusterName,
						*nnp.Name,
						scaleUpRequest,
					)
					if err != nil {
						return "", fmt.Errorf("Error occred while scaling up nodepool %s: %s", *nnp.Name, err)
					}
					taskUUID = karbonClusterActionResponse.TaskUUID
				}
				if *cnp.NumInstances > *nnp.NumInstances {
					log.Print("scale down")
					amountOfNodes := *cnp.NumInstances - *nnp.NumInstances
					scaleDownRequest := &karbon.ClusterScaleDownIntentInput{
						Count: amountOfNodes,
					}
					// taskUUID, err = scaleDownNodepool(client, karbonClusterName, nnp, amountOfNodes)
					karbonClusterActionResponse, err := client.KarbonAPI.Cluster.ScaleDownKarbonCluster(
						karbonClusterName,
						*nnp.Name,
						scaleDownRequest,
					)
					if err != nil {
						return "", fmt.Errorf("Error occred while scaling down nodepool %s: %s", *nnp.Name, err)
					}
					taskUUID = karbonClusterActionResponse.TaskUUID
				}
				log.Print("no match?")
			}
		}
	}
	return taskUUID, nil
}

func getSupportedFileSystems() []string {
	return []string{
		"ext4",
		"xfs",
	}
}

func getSupportedReclaimPolicies() []string {
	return []string{
		"Delete",
		"Retain",
	}
}
