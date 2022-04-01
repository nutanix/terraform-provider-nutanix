package nutanix

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/fc"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const (
	aggregatePercentComplete = 100
	ImageMinTimeout          = 2 * time.Hour
	DelayTime                = 15 * time.Minute
	NodePollTimeout          = 30 * time.Minute
)

func resourceNutanixFCImageCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixFCImageClusterCreate,
		ReadContext:   resourceNutanixFCImageClusterRead,
		UpdateContext: resourceNutanixFCImageClusterUpdate,
		DeleteContext: resourceNutanixFCImageClusterDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(ImageMinTimeout),
		},
		Schema: map[string]*schema.Schema{
			"cluster_external_ip": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"common_network_settings": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cvm_dns_servers": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"hypervisor_dns_servers": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"cvm_ntp_servers": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"hypervisor_ntp_servers": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"hypervisor_iso_details": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hyperv_sku": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"hyperv_product_key": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"sha256sum": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"storage_node_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"redundancy_factor": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"aos_package_url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"aos_package_sha256sum": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"image_cluster_uuid": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"node_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cvm_gateway": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ipmi_netmask": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"rdma_passthrough": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"imaged_node_uuid": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"cvm_vlan_id": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"hypervisor_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"image_now": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"hypervisor_hostname": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"hypervisor_netmask": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"cvm_netmask": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ipmi_ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"hypervisor_gateway": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"hardware_attributes_override": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
						},
						"cvm_ram_gb": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"cvm_ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"hypervisor_ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ipmi_gateway": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"use_existing_network_settings": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"skip_cluster_creation": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceNutanixFCImageClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// conn := meta.(*Client).FC
	// resp, err := conn.GetImagedNode(ctx, d.Id())
	// if err != nil {
	// 	diag.FromErr(err)
	// }

	/*
		clusterExternalIP, ok := d.GetOk("cluster_external_ip")
		if !ok {
			log.Println("cluster_external_ip is not set")
		}
		req.ClusterExternalIP = utils.StringPtr(clusterExternalIP.(string))

		storageCount, ok := d.GetOk("storage_node_count")
		if !ok {
			log.Println("storage_node_count is not set")
		}
		req.StorageNodeCount = utils.IntPtr(storageCount.(int))

		redundancyFactor, ok := d.GetOk("redundancy_factor")
		if !ok {
			log.Println("redundancy_factor is not set")
		}
		req.RedundancyFactor = utils.IntPtr(redundancyFactor.(int))

		clusterName, ok := d.GetOk("cluster_name")
		if !ok {
			log.Println("cluster_name is not set")
		}
		req.ClusterName = utils.StringPtr(clusterName.(string))

		aosPackageURL, ok := d.GetOk("aos_package_url")
		if !ok {
			log.Println("aos_package_url is not set")
		}
		req.AosPackageURL = utils.StringPtr(aosPackageURL.(string))

		aosPackageSha, ok := d.GetOk("aos_package_sha256sum")
		if !ok {
			log.Println("aos_package_url is not set")
		}
		req.AosPackageSha256sum = utils.StringPtr(aosPackageSha.(string))

		clusterSize, ok := d.GetOk("cluster_size")
		if !ok {
			log.Println("cluster_size is not set")
		}
		req.ClusterSize = utils.IntPtr(clusterSize.(int))

		timezone, ok := d.GetOk("timezone")
		if !ok {
			log.Println("timezone is not set")
		}
		req.Timezone = utils.StringPtr(timezone.(string))
	*/

	// req := fc.CreateClusterInput{}
	return nil
}

func foundationCentralClusterRefresh(ctx context.Context, conn *fc.Client, imageUUID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := conn.Service.GetImagedCluster(ctx, imageUUID)

		if err != nil {
			return nil, "FAILED", err
		}
		if *v.ClusterStatus.ImagingStopped {
			return v, "COMPLETED", nil
		}
		return v, "PENDING", nil
	}
}

func expandCommonNetworkSettings(d *schema.ResourceData) *fc.CommonNetworkSettings {
	cns := fc.CommonNetworkSettings{}
	resourceData, ok := d.GetOk("common_network_settings")
	if !ok {
		return nil
	}
	settingsMap := resourceData.([]interface{})[0].(map[string]interface{})

	cns.CvmDNSServers = expandNetworklist(settingsMap["cvm_dns_servers"])
	cns.CvmNtpServers = expandNetworklist(settingsMap["cvm_ntp_servers"])
	cns.HypervisorDNSServers = expandNetworklist(settingsMap["hypervisor_dns_servers"])
	cns.HypervisorNtpServers = expandNetworklist(settingsMap["hypervisor_ntp_servers"])

	return &cns
}

func expandNetworklist(pr interface{}) []string {
	c := make([]string, 0)
	prList := pr.([]interface{})

	for _, k := range prList {
		c = append(c, k.(string))
	}
	return c
}

func expandHyperVisorIsoDetails(d *schema.ResourceData) *fc.HypervisorIsoDetails {
	hid := fc.HypervisorIsoDetails{}
	resourceData, ok := d.GetOk("hypervisor_iso_details")
	if !ok {
		return nil
	}
	settingsMap := resourceData.([]interface{})[0].(map[string]interface{})

	hid.HypervSku = utils.StringPtr(settingsMap["hyperv_sku"].(string))
	hid.URL = utils.StringPtr(settingsMap["url"].(string))
	hid.HypervProductKey = utils.StringPtr(settingsMap["hyperv_product_key"].(string))
	hid.Sha256sum = utils.StringPtr(settingsMap["sha256sum"].(string))

	return &hid
}

func expandNodesList(d *schema.ResourceData) []*fc.Node {
	nodeList := []*fc.Node{}
	resourceData, ok := d.GetOk("node_list")
	if !ok {
		return nil
	}
	nodesConfig := resourceData.([]interface{})

	for _, nodeConfig := range nodesConfig {
		nodeSettings := nodeConfig.(map[string]interface{})
		node := fc.Node{}
		if cvmGateway, ok := nodeSettings["cvm_gateway"]; ok {
			node.CvmGateway = utils.StringPtr(cvmGateway.(string))
		}
		if ipmiGateway, ok := nodeSettings["ipmi_gateway"]; ok {
			node.IpmiGateway = utils.StringPtr(ipmiGateway.(string))
		}
		if ipmiNetmask, ok := nodeSettings["ipmi_netmask"]; ok {
			node.IpmiNetmask = utils.StringPtr(ipmiNetmask.(string))
		}
		if ipmiIP, ok := nodeSettings["ipmi_ip"]; ok {
			node.IpmiIP = utils.StringPtr(ipmiIP.(string))
		}
		if hypGateway, ok := nodeSettings["hypervisor_gateway"]; ok {
			node.HypervisorGateway = utils.StringPtr(hypGateway.(string))
		}
		if imageNodeUUID, ok := nodeSettings["imaged_node_uuid"]; ok {
			node.ImagedNodeUUID = utils.StringPtr(imageNodeUUID.(string))
		}
		if hypervisorType, ok := nodeSettings["hypervisor_type"]; ok {
			node.HypervisorType = utils.StringPtr(hypervisorType.(string))
		}
		if hypervisorHostname, ok := nodeSettings["hypervisor_hostname"]; ok {
			node.HypervisorHostname = utils.StringPtr(hypervisorHostname.(string))
		}
		if hypervisorNetmask, ok := nodeSettings["hypervisor_netmask"]; ok {
			node.HypervisorNetmask = utils.StringPtr(hypervisorNetmask.(string))
		}
		if cvmNetmask, ok := nodeSettings["cvm_netmask"]; ok {
			node.CvmNetmask = utils.StringPtr(cvmNetmask.(string))
		}
		if cvmIP, ok := nodeSettings["cvm_ip"]; ok {
			node.CvmIP = utils.StringPtr(cvmIP.(string))
		}
		if hypervisorIP, ok := nodeSettings["hypervisor_ip"]; ok {
			node.HypervisorIP = utils.StringPtr(hypervisorIP.(string))
		}

		if cvmVlanID, ok := nodeSettings["cvm_vlan_id"]; ok {
			node.CvmVlanID = utils.IntPtr(cvmVlanID.(int))
		}
		if cvmRAMGb, ok := nodeSettings["cvm_ram_gb"]; ok {
			node.CvmRAMGb = utils.IntPtr(cvmRAMGb.(int))
		}

		if rdmaPassthrough, ok := nodeSettings["rdma_passthrough"]; ok {
			node.RdmaPassthrough = utils.BoolPtr(rdmaPassthrough.(bool))
		}
		if imageNow, ok := nodeSettings["image_now"]; ok {
			node.ImageNow = utils.BoolPtr(imageNow.(bool))
		}
		if useExistingNetworkSettings, ok := nodeSettings["use_existing_network_settings"]; ok {
			node.UseExistingNetworkSettings = utils.BoolPtr(useExistingNetworkSettings.(bool))
		}

		if hardwareAttrs, ok := nodeSettings["hardware_attributes_override"]; ok {
			// Convert map to json string
			jsonStr, err := json.Marshal(hardwareAttrs)
			if err != nil {
				fmt.Println(err)
			}
			// Convert json string to map[string]interface{}
			var mapData map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &mapData); err != nil {
				fmt.Println(err)
			}
			node.HardwareAttributesOverride = mapData
		}
		nodeList = append(nodeList, &node)
	}

	return nodeList
}

func resourceNutanixFCImageClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// Get client connection
	conn := meta.(*Client).FC
	req := fc.CreateClusterInput{}

	clusterExternalIP, ok := d.GetOk("cluster_external_ip")
	if !ok {
		log.Println("cluster_external_ip is not set")
	}
	req.ClusterExternalIP = utils.StringPtr(clusterExternalIP.(string))

	storageCount, ok := d.GetOk("storage_node_count")
	if !ok {
		log.Println("storage_node_count is not set")
	}
	req.StorageNodeCount = utils.IntPtr(storageCount.(int))

	redundancyFactor, ok := d.GetOk("redundancy_factor")
	if !ok {
		log.Println("redundancy_factor is not set")
	}
	req.RedundancyFactor = utils.IntPtr(redundancyFactor.(int))

	clusterName, ok := d.GetOk("cluster_name")
	if !ok {
		log.Println("cluster_name is not set")
	}
	req.ClusterName = utils.StringPtr(clusterName.(string))

	aosPackageURL, ok := d.GetOk("aos_package_url")
	if !ok {
		log.Println("aos_package_url is not set")
	}
	req.AosPackageURL = utils.StringPtr(aosPackageURL.(string))

	aosPackageSha, ok := d.GetOk("aos_package_sha256sum")
	if !ok {
		log.Println("aos_package_url is not set")
	}
	req.AosPackageSha256sum = utils.StringPtr(aosPackageSha.(string))

	clusterSize, ok := d.GetOk("cluster_size")
	if !ok {
		log.Println("cluster_size is not set")
	}
	req.ClusterSize = utils.IntPtr(clusterSize.(int))

	timezone, ok := d.GetOk("timezone")
	if !ok {
		log.Println("timezone is not set")
	}
	req.Timezone = utils.StringPtr(timezone.(string))

	skipClusterCreation, ok := d.GetOk("skip_cluster_creation")
	if !ok {
		req.SkipClusterCreation = false
	}
	req.SkipClusterCreation = skipClusterCreation.(bool)

	req.CommonNetworkSettings = expandCommonNetworkSettings(d)
	req.HypervisorIsoDetails = expandHyperVisorIsoDetails(d)
	req.NodesList = expandNodesList(d)

	// Poll for operation here - Node Detail GET Call
	for _, vv := range req.NodesList {
		stateConfig := &resource.StateChangeConf{
			Pending: []string{"STATE_DISCOVERING", "STATE_UNAVAILABLE"},
			Target:  []string{"STATE_AVAILABLE"},
			Refresh: foundationCentralPollingNode(ctx, conn, *vv.ImagedNodeUUID),
			Timeout: NodePollTimeout,
			Delay:   10 * time.Second,
		}
		infos, err := stateConfig.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error waiting for node (%s) to be available: %v", *vv.ImagedNodeUUID, err)
		}
		if progress, ok := infos.(*fc.ImagedNodeDetails); ok {
			if !(*progress.Available) {
				return diag.Errorf("Node is not available to image or alraedy be a part of cluster")
			}
		}
	}
	//Make request to the API
	resp, err := conn.Service.CreateCluster(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.ImagedClusterUUID == nil {
		return diag.Errorf("returned image cluster uuid is empty")
	}

	d.SetId(*resp.ImagedClusterUUID)

	// Poll for operation here - Cluster GET Call
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED", "FAILED"},
		Refresh: foundationCentralClusterRefresh(ctx, conn, *resp.ImagedClusterUUID),
		Timeout: d.Timeout(schema.TimeoutCreate),
		Delay:   DelayTime,
	}
	info, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for image (%s) to be ready: %v", *resp.ImagedClusterUUID, err)
	}

	if progress, ok := info.(*fc.ImagedClusterDetails); ok {
		if utils.Float64Value(progress.ClusterStatus.AggregatePercentComplete) < float64(aggregatePercentComplete) {
			return collectIndividualErrorDiagnosticsFC(progress)
		}
	}

	return nil
}

func resourceNutanixFCImageClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixFCImageClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).FC
	log.Printf("[DEBUG] Deleting Cluster: %s, %s", d.Get("cluster_name").(string), d.Id())
	err := conn.DeleteCluster(ctx, d.Id())
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return nil
		}
		return diag.Errorf("error while Deleting Cluster: UUID(%s): %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}

func collectIndividualErrorDiagnosticsFC(progress *fc.ImagedClusterDetails) diag.Diagnostics {
	// create empty diagnostics
	var diags diag.Diagnostics

	// append errors for failed node imaging
	for _, v := range progress.ClusterStatus.NodeProgressDetails {
		if utils.Float64Value(v.PercentComplete) < float64(aggregatePercentComplete) {
			message := ""
			for _, v1 := range v.MessageList {
				message += *v1
			}
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Node imaging for imaged_node_uuid IP: %s failed with error:  %s.", *v.ImagedNodeUUID, *v.Status),
				Detail:   message,
			})
		}
	}

	// append errors for failed cluster creation
	for _, v := range progress.ClusterStatus.ClusterProgressDetails {
		if utils.Float64Value(v.PercentComplete) < float64(aggregatePercentComplete) {
			message := ""
			for _, v1 := range v.MessageList {
				message += *v1
			}
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Cluster creation for Cluster : %s failed with error:  %s.", *v.ClusterName, *v.Status),
				Detail:   message,
			})
		}
	}
	return diags
}

func foundationCentralPollingNode(ctx context.Context, conn *fc.Client, imageUUID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Polling Node to be Available: %s", imageUUID)
		v, err := conn.Service.GetImagedNode(ctx, imageUUID)
		if err != nil {
			return nil, *v.NodeState, err
		}

		if *v.NodeState == "STATE_UNAVAILABLE" || *v.NodeState == "STATE_DISCOVERING" {
			return v, *v.NodeState, nil

		}
		return v, "STATE_AVAILABLE", nil
	}
}
