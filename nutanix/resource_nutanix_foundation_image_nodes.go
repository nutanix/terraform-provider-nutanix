package nutanix

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/fc"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func resourceNutanixFCImageCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixFCImageClusterCreate,
		ReadContext:   resourceNutanixFCImageClusterRead,
		UpdateContext: resourceNutanixFCImageClusterUpdate,
		DeleteContext: resourceNutanixFCImageClusterDelete,
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
				Type:     schema.TypeBool,
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
		},
	}
}

func resourceNutanixFCImageClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// conn := meta.(*Client).FC
	// resp, err := conn.GetImagedNode(d.Id())
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
		req.AosPackageUrl = utils.StringPtr(aosPackageURL.(string))

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

func expandCommonNetworkSettings(d *schema.ResourceData) *fc.CommonNetworkSettings {
	cns := fc.CommonNetworkSettings{}
	resourceData, ok := d.GetOk("common_network_settings")
	if !ok {
		return nil
	}
	settingsMap := resourceData.([]interface{})[0].(map[string]interface{})

	cns.CvmDnsServers = settingsMap["cvm_dns_servers"].([]string)
	cns.CvmNtpServers = settingsMap["cvm_ntp_servers"].([]string)
	cns.HypervisorDnsServers = settingsMap["hypervisor_dns_servers"].([]string)
	cns.HypervisorNtpServers = settingsMap["hypervisor_ntp_servers"].([]string)

	return &cns
}

func expandHyperVisorIsoDetails(d *schema.ResourceData) *fc.HypervisorIsoDetails {
	hid := fc.HypervisorIsoDetails{}
	resourceData, ok := d.GetOk("hypervisor_iso_details")
	if !ok {
		return nil
	}
	settingsMap := resourceData.([]interface{})[0].(map[string]interface{})

	hid.HypervSku = utils.StringPtr(settingsMap["hyperv_sku"].(string))
	hid.Url = utils.StringPtr(settingsMap["url"].(string))
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
		if cvmRamGb, ok := nodeSettings["cvm_ram_gb"]; ok {
			node.CvmRamGb = utils.IntPtr(cvmRamGb.(int))
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
			node.HardwareAttributesOverride = hardwareAttrs.(map[string]interface{})
		}
		nodeList = append(nodeList, &node)
	}

	return nodeList
}

func resourceNutanixFCImageClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Create FC Image: %s", d.Get("name").(string))

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
	req.AosPackageUrl = utils.StringPtr(aosPackageURL.(string))

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

	req.CommonNetworkSettings = expandCommonNetworkSettings(d)
	req.HypervisorIsoDetails = expandHyperVisorIsoDetails(d)
	req.NodesList = expandNodesList(d)

	//Make request to the API
	resp, err := conn.Service.CreateCluster(&req)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.ImagedClusterUUID == nil {
		return diag.Errorf("returned image cluster uuid is empty")
	}

	d.SetId(*resp.ImagedClusterUUID)

	// Poll for operation here

	return nil
}

func resourceNutanixFCImageClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixFCImageClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
