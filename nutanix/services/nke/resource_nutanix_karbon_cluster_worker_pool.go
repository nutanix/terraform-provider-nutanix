package nke

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/karbon"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixKarbonWorkerNodePool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixKarbonWorkerNodePoolCreate,
		ReadContext:   resourceNutanixKarbonWorkerNodePoolRead,
		UpdateContext: resourceNutanixKarbonWorkerNodePoolUpdate,
		DeleteContext: resourceNutanixKarbonWorkerNodePoolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(DEFAULTWAITTIMEOUT * time.Minute),
			Update: schema.DefaultTimeout(DEFAULTWAITTIMEOUT * time.Minute),
			Delete: schema.DefaultTimeout(DEFAULTWAITTIMEOUT * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_os_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"num_instances": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validation.IntAtLeast(MINNUMINSTANCES),
			},
			"ahv_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cpu": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      "8",
							ValidateFunc: validation.IntAtLeast(MINCPU),
						},
						"disk_mib": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      "122880",
							ValidateFunc: validation.IntAtLeast(DEFAULTWORKERNODEDISKMIB),
						},
						"memory_mib": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      "8192",
							ValidateFunc: validation.IntAtLeast(DEFAULTWORKERNODEEMORYMIB),
						},
						"network_uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"prism_element_cluster_uuid": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"iscsi_network_uuid": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
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
	}
}

func resourceNutanixKarbonWorkerNodePoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*conns.Client)
	conn := client.KarbonAPI

	addworkerRequest := &karbon.ClusterNodePool{}
	nkeName := ""
	if karbonNodeName, ok := d.GetOk("cluster_name"); ok && len(karbonNodeName.(string)) > 0 {
		nkeName = karbonNodeName.(string)
	} else {
		return diag.Errorf("cluster_name is a required field")
	}

	if workerName, ok := d.GetOk("name"); ok {
		addworkerRequest.Name = utils.StringPtr(workerName.(string))
	}

	if numOfInst, ok := d.GetOk("num_instances"); ok {
		numInstances := int64(numOfInst.(int))
		addworkerRequest.NumInstances = &numInstances
	}

	if ahvConfig, ok := d.GetOk("ahv_config"); ok {
		ahvConfigList := ahvConfig.([]interface{})
		nodepool := &karbon.ClusterNodePool{
			AHVConfig: &karbon.ClusterNodePoolAHVConfig{},
		}
		if len(ahvConfigList) != 1 {
			return diag.Errorf("ahv_config must have 1 element")
		}
		ahvConfig := ahvConfigList[0].(map[string]interface{})
		if valCPU, ok := ahvConfig["cpu"]; ok {
			i := int64(valCPU.(int))
			nodepool.AHVConfig.CPU = i
		}
		if valDiskMib, ok := ahvConfig["disk_mib"]; ok {
			i := int64(valDiskMib.(int))
			nodepool.AHVConfig.DiskMib = i
		}
		if valMemoryMib, ok := ahvConfig["memory_mib"]; ok {
			i := int64(valMemoryMib.(int))
			nodepool.AHVConfig.MemoryMib = i
		}
		if valNetworkUUID, ok := ahvConfig["network_uuid"]; ok && len(valNetworkUUID.(string)) > 0 {
			nodepool.AHVConfig.NetworkUUID = valNetworkUUID.(string)
		}
		if valPrismElementClusterUUID, ok := ahvConfig["prism_element_cluster_uuid"]; ok && len(valPrismElementClusterUUID.(string)) > 0 {
			nodepool.AHVConfig.PrismElementClusterUUID = valPrismElementClusterUUID.(string)
		}
		if valICSUUUID, ok := ahvConfig["iscsi_network_uuid"]; ok && len(valICSUUUID.(string)) > 0 {
			nodepool.AHVConfig.IscsiNetworkUUID = valICSUUUID.(string)
		}
		addworkerRequest.AHVConfig = nodepool.AHVConfig
	}
	if label, ok := d.GetOk("labels"); ok && label.(map[string]interface{}) != nil {
		addworkerRequest.Labels = utils.ConvertMapString(label.(map[string]interface{}))
	}
	karbonClusterActionResponse, err := conn.Cluster.AddWorkerNodePool(
		nkeName,
		addworkerRequest,
	)
	if err != nil {
		return diag.FromErr(err)
	}
	err = WaitForKarbonCluster(ctx, client, 0, karbonClusterActionResponse.TaskUUID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(karbonClusterActionResponse.TaskUUID)
	return resourceNutanixKarbonWorkerNodePoolRead(ctx, d, meta)
}

func resourceNutanixKarbonWorkerNodePoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).KarbonAPI
	setTimeout(meta)
	// Make request to the API
	var err error
	karbonClsName := d.Get("cluster_name")
	resp, err := conn.Cluster.GetKarbonCluster(karbonClsName.(string))
	if err != nil {
		d.SetId("")
		return nil
	}
	karbonClusterName := *resp.Name
	workerName := d.Get("name")
	nodepool, err := conn.Cluster.GetKarbonClusterNodePool(karbonClusterName, workerName.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	nodes := make([]map[string]interface{}, 0)
	for _, npn := range *nodepool.Nodes {
		nodes = append(nodes, map[string]interface{}{
			"hostname":     npn.Hostname,
			"ipv4_address": npn.IPv4Address,
		})
	}
	if err = d.Set("name", nodepool.Name); err != nil {
		return diag.Errorf("error setting name for nke Worker Node Pool %s: %s", d.Id(), err)
	}
	if err = d.Set("node_os_version", nodepool.NodeOSVersion); err != nil {
		return diag.Errorf("error setting node_os_version for nke Worker Node Pool %s: %s", d.Id(), err)
	}
	if err = d.Set("num_instances", nodepool.NumInstances); err != nil {
		return diag.Errorf("error setting num_instances for nke Worker Node Pool %s: %s", d.Id(), err)
	}
	if err = d.Set("nodes", nodes); err != nil {
		return diag.Errorf("error setting nodes for nke Worker Node Pool %s: %s", d.Id(), err)
	}
	if err = d.Set("labels", nodepool.Labels); err != nil {
		return diag.Errorf("error setting labels for nke Worker Node Pool %s: %s", d.Id(), err)
	}
	if err = d.Set("ahv_config", flattenAHVNodePoolConfig(nodepool.AHVConfig)); err != nil {
		return diag.Errorf("error setting ahv_config for nke Worker Node Pool %s: %s", d.Id(), err)
	}
	return nil
}

func resourceNutanixKarbonWorkerNodePoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*conns.Client)
	conn := client.KarbonAPI

	karbonClsName := d.Get("cluster_name")
	workerName := d.Get("name")
	resp, err := conn.Cluster.GetKarbonCluster(karbonClsName.(string))
	if err != nil {
		d.SetId("")
		return nil
	}
	nodepool, err := conn.Cluster.GetKarbonClusterNodePool(*resp.Name, workerName.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("num_instances") {
		old, new := d.GetChange("num_instances")

		if old.(int) > new.(int) {
			amountOfNodes := old.(int) - new.(int)
			scaleDownRequest := &karbon.ClusterScaleDownIntentInput{
				Count: int64(amountOfNodes),
			}
			karbonClusterActionResponse, err := client.KarbonAPI.Cluster.ScaleDownKarbonCluster(
				*resp.Name,
				*nodepool.Name,
				scaleDownRequest,
			)
			if err != nil {
				return diag.FromErr(err)
			}
			err = WaitForKarbonCluster(ctx, client, 0, karbonClusterActionResponse.TaskUUID, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			amountOfNodes := new.(int) - old.(int)
			scaleUpRequest := &karbon.ClusterScaleUpIntentInput{
				Count: int64(amountOfNodes),
			}
			karbonClusterActionResponse, err := client.KarbonAPI.Cluster.ScaleUpKarbonCluster(
				*resp.Name,
				*nodepool.Name,
				scaleUpRequest,
			)
			if err != nil {
				return diag.FromErr(err)
			}
			err = WaitForKarbonCluster(ctx, client, 0, karbonClusterActionResponse.TaskUUID, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if d.HasChange("labels") {
		old, new := d.GetChange("labels")
		updateLabelRequest := &karbon.UpdateWorkerNodeLabels{}

		newMap := new.(map[string]interface{})
		oldMap := old.(map[string]interface{})
		addLabelMap := map[string]string{}
		removeLabel := []string{}

		// check any new is label is added.
		for key := range newMap {
			if _, ok := oldMap[key]; ok {
				continue
			} else {
				addLabelMap[key] = (newMap[key]).(string)
			}
		}
		// check any label is removed
		for key := range oldMap {
			if _, ok := newMap[key]; ok {
				continue
			} else {
				removeLabel = append(removeLabel, key)
			}
		}

		updateLabelRequest.AddLabel = addLabelMap
		updateLabelRequest.RemoveLabel = removeLabel

		nodeLabelActionResponse, err := client.KarbonAPI.Cluster.UpdateWorkerNodeLables(
			*resp.Name,
			*nodepool.Name,
			updateLabelRequest,
		)
		if err != nil {
			return diag.FromErr(err)
		}
		err = WaitForKarbonCluster(ctx, client, 0, nodeLabelActionResponse.TaskUUID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceNutanixKarbonWorkerNodePoolRead(ctx, d, meta)
}

func resourceNutanixKarbonWorkerNodePoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*conns.Client)
	conn := client.KarbonAPI

	var err error
	karbonClsName := d.Get("cluster_name")
	resp, err := conn.Cluster.GetKarbonCluster(karbonClsName.(string))
	if err != nil {
		d.SetId("")
		return nil
	}
	karbonClusterName := *resp.Name
	workerName := d.Get("name")
	nodepool, err := conn.Cluster.GetKarbonClusterNodePool(karbonClusterName, workerName.(string))
	if err != nil {
		return diag.FromErr(err)
	}
	nodes := []*string{}
	for _, v := range *nodepool.Nodes {
		nodes = append(nodes, v.Hostname)
	}
	removeWorkerRequest := &karbon.RemoveWorkerNodeRequest{
		NodeList: nodes,
	}
	workerPoolActionResponse, err := conn.Cluster.RemoveWorkerNodePool(
		karbonClusterName,
		workerName.(string),
		removeWorkerRequest,
	)
	if err != nil {
		return diag.FromErr(err)
	}
	err = WaitForKarbonCluster(ctx, client, 0, workerPoolActionResponse.TaskUUID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	workerNodeDeleteResponse, er := client.KarbonAPI.Cluster.DeleteWorkerNodePool(
		karbonClusterName,
		workerName.(string),
	)
	if er != nil {
		return diag.FromErr(er)
	}
	err = WaitForKarbonCluster(ctx, client, 0, workerNodeDeleteResponse.TaskUUID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func flattenAHVNodePoolConfig(ahv *karbon.ClusterNodePoolAHVConfig) []map[string]interface{} {
	if ahv != nil {
		ahvConfig := make([]map[string]interface{}, 0)

		config := map[string]interface{}{}

		config["cpu"] = ahv.CPU
		config["disk_mib"] = ahv.DiskMib
		config["memory_mib"] = ahv.MemoryMib
		config["network_uuid"] = ahv.NetworkUUID
		config["prism_element_cluster_uuid"] = ahv.PrismElementClusterUUID
		if ahv.IscsiNetworkUUID != "" {
			config["iscsi_network_uuid"] = ahv.IscsiNetworkUUID
		}

		ahvConfig = append(ahvConfig, config)
		return ahvConfig
	}
	return nil
}
