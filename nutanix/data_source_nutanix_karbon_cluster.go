package nutanix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/karbon"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func dataSourceNutanixKarbonCluster() *schema.Resource {
	return &schema.Resource{
		Read:          dataSourceNutanixKarbonClusterRead,
		SchemaVersion: 1,
		Schema:        KarbonClusterDataSourceMap(),
	}
}

func dataSourceNutanixKarbonClusterRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*Client).KarbonAPI
	setTimeout(meta)
	// Make request to the API
	karbonClusterID, iok := d.GetOk("karbon_cluster_id")
	karbonClusterName, nok := d.GetOk("karbon_cluster_name")
	if !iok && !nok {
		return fmt.Errorf("please provide one of karbon_cluster_id or karbon_cluster_name attributes")
	}
	var err error
	var resp *karbon.KarbonClusterIntentResponse

	if iok {
		resp, err = conn.Cluster.GetKarbonCluster(karbonClusterID.(string))
	} else {
		resp, err = conn.Cluster.GetKarbonCluster(karbonClusterName.(string))
	}

	if err != nil {
		d.SetId("")
		return err
	}

	karbon_cluster_name := *resp.Name
	flattenedEtcdNodepool, err := flattenNodePools(d, conn, "etcd_node_pool", karbon_cluster_name, resp.ETCDConfig.NodePools)
	if err != nil {
		return err
	}
	flattenedWorkerNodepool, err := flattenNodePools(d, conn, "worker_node_pool", karbon_cluster_name, resp.WorkerConfig.NodePools)
	if err != nil {
		return err
	}
	flattenedMasterNodepool, err := flattenNodePools(d, conn, "master_node_pool", karbon_cluster_name, resp.MasterConfig.NodePools)
	if err != nil {
		return err
	}
	d.Set("name", utils.StringValue(resp.Name))

	d.Set("status", utils.StringValue(resp.Status))

	//Must use legacy API because GA API reports different version
	log.Printf("Getting existing version: %s", d.Get("version").(string))
	d.Set("version", d.Get("version").(string))
	// d.Set("version", utils.StringValue(resp.Version))
	// d.Set("version", utils.StringValue(respLegacy.K8sConfig.Version))
	d.Set("kubeapi_server_ipv4_address", utils.StringValue(resp.KubeApiServerIPv4Address))
	d.Set("deployment_type", resp.MasterConfig.DeploymentType)
	d.Set("worker_node_pool", flattenedWorkerNodepool)

	d.Set("etcd_node_pool", flattenedEtcdNodepool)
	d.Set("master_node_pool", flattenedMasterNodepool)

	d.SetId(*resp.UUID)

	return nil
}

func KarbonClusterDataSourceMap() map[string]*schema.Schema {
	kcsm := KarbonClusterElementDataSourceMap()
	kcsm["karbon_cluster_id"] = &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		ConflictsWith: []string{"karbon_cluster_name"},
	}
	kcsm["karbon_cluster_name"] = &schema.Schema{
		Type:          schema.TypeString,
		Optional:      true,
		ConflictsWith: []string{"karbon_cluster_id"},
	}
	return kcsm
}

func KarbonClusterElementDataSourceMap() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"version": {
			Type:     schema.TypeString,
			Computed: true,
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
		"etcd_node_pool":   nodePoolDatasourceSchema(),
		"master_node_pool": nodePoolDatasourceSchema(),
		"worker_node_pool": nodePoolDatasourceSchema(),
	}
}

func nodePoolDatasourceSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"node_os_version": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"num_instances": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"ahv_config": {
					Type: schema.TypeMap,

					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"cpu": {
								Type:     schema.TypeInt,
								Computed: true,
							},
							"disk_mib": {
								Type:     schema.TypeInt,
								Computed: true,
							},
							"memory_mib": {
								Type:     schema.TypeInt,
								Computed: true,
							},
							"network_uuid": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"prism_element_cluster_uuid": {
								Type:     schema.TypeString,
								Computed: true,
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
