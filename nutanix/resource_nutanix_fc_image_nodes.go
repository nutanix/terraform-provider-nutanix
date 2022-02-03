package nutanix

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceNutanixFCImageCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceNutanixFCImageClusterCreate,
		Read:   resourceNutanixFCImageClusterRead,
		Update: resourceNutanixFCImageClusterUpdate,
		Delete: resourceNutanixFCImageClusterDelete,
		Schema: map[string]*schema.Schema{
			"cluster_external_ip": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"common_network_settings": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cvm_dns_servers": {
							Type:     schema.TypeList,
							Computed: true,
						},
						"hypervisor_dns_servers": {
							Type:     schema.TypeList,
							Computed: true,
						},
						"cvm_ntp_servers": {
							Type:     schema.TypeList,
							Computed: true,
						},
						"hypervisor_ntp_servers": {
							Type:     schema.TypeList,
							Computed: true,
						},
					},
				},
			},
			"hypervisor_iso_details": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hyperv_sku": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"hyperv_product_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sha256sum": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"storage_node_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"redundancy_factor": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"aos_package_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"aos_package_sha256sum": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"timezone": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"node_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cvm_gateway": {
							Type:     schema.TypeString,
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
							Computed: true,
						},
						"cvm_vlan_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"hypervisor_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"image_now": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"hypervisor_hostname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hypervisor_netmask": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cvm_netmask": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ipmi_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hypervisor_gateway": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hardware_attributes_override": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"cvm_ram_gb": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"cvm_ip": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"hypervisor_ip": {
							Type:     schema.TypeString,
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

func resourceNutanixFCImageClusterRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

// func expandCommonNetworkSettings(commonNetworkSettings interface{}) (*fc.CommonNetworkSettings, error) {
// 	cns := fc.CommonNetworkSettings{}
// 	settingsMap := commonNetworkSettings.(map[string]interface{})
// 	cns.CvmDnsServers = settingsMap["cvm_dns_server"].([]interface{})
// 	cns.CvmNtpServers = settingsMap["hypervisor_dns_servers"].([]interface{})
// 	cns.HypervisorDnsServers = settingsMap[""].([]interface{})
// }

func resourceNutanixFCImageClusterCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Create FC Image: %s", d.Get("name").(string))

	// Get client connection
	// conn := meta.(*Client).FC
	// uuid := d.Id()
	// input := new(fc.CreateClusterInput)
	// cluster_external_ip
	// common_network_settings
	// hypervisor_iso_details
	// storage_node_count
	// redundancy_factor
	// cluster_name
	// aos_package_url
	// cluster_size
	// aos_package_sha256sum
	// timezone
	// nodes_list

	// clusterExternalIP, _ := d.GetOk("cluster_external_ip")
	// expandCommonNetworkSettings(d.Get("common_network_settings"))

	//Make request to the API
	// resp, err := conn.Service.CreateCluster(input)
	// if err != nil {
	// 	if strings.Contains(fmt.Sprint(err), "ENITY_NOT_FOUND") {
	// 		d.SetId("")
	// 		return nil
	// 	}
	// 	return fmt.Errorf("error reading image UUID (%s) with error %s", uuid, err)
	// }

	// if err = d.Set("node_type", utils.StringValue(resp.NodeType)); err != nil {
	// 	return fmt.Errorf("error setting owner_reference for image UUID(%s), %s", d.Id(), err)
	// }

	// if err = d.Set("hardware_attributes", flattenReferenceValues(resp.Status.HardwareAttributes)); err != nil {
	// 	return fmt.Errorf("error setting owner_reference for image UUID(%s), %s", d.Id(), err)
	// }

	// if err = d.Set("node_serial", resp.Status.NodeSerial); err != nil {
	// 	return fmt.Errorf("error setting state for image UUID(%s), %s", d.Id(), err)
	// }

	// if err = d.Set("block_serial", resp.Status.BlockSerial); err != nil {
	// 	return fmt.Errorf("error setting image_type for image UUID(%s), %s", d.Id(), err)
	// }

	// if err = d.Set("model", resp.Status.Model); err != nil {
	// 	return fmt.Errorf("error setting source_uri for image UUID(%s), %s", d.Id(), err)
	// }

	// checksum := make(map[string]string)
	// if resp.Status.Resources.Checksum != nil {
	// 	checksum["checksum_algorithm"] = utils.StringValue(resp.Status.Resources.Checksum.ChecksumAlgorithm)
	// 	checksum["checksum_value"] = utils.StringValue(resp.Status.Resources.Checksum.ChecksumValue)
	// }

	return nil
}

func resourceNutanixFCImageClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	// client := meta.(*Client)
	// conn := client.API
	// timeout := client.WaitTimeout

	// if client.WaitTimeout == 0 {
	// 	timeout = 10
	// }

	// // get state
	// request := &v3.ImageIntentInput{}
	// metadata := &v3.Metadata{}
	// spec := &v3.Image{}
	// res := &v3.ImageResources{}

	// response, err := conn.fc.GetImage(d.Id())

	// if err != nil {
	// 	if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
	// 		d.SetId("")
	// 	}
	// 	return err
	// }

	// if d.HasChange("node_types") {
	// 	metadata.Categories = expandCategories(d.Get("node_types"))
	// }

	// if d.HasChange("hardware_attributes") {
	// 	or := d.Get("hardware_attributes").(map[string]interface{})
	// 	metadata.OwnerReference = validateRef(or)
	// }

	// if d.HasChange("node_serial") {
	// 	pr := d.Get("node_serial").(map[string]interface{})
	// 	metadata.ProjectReference = validateRef(pr)
	// }

	// if d.HasChange("block_serial") {
	// 	spec.Name = utils.StringPtr(d.Get("block_serial").(string))
	// }

	// if d.HasChange("model") {
	// 	spec.Description = utils.StringPtr(d.Get("model").(string))
	// }

	// resp, errUpdate := conn.fc.UpdateImage(d.Id(), request)

	// if errUpdate != nil {
	// 	return fmt.Errorf("error updating image(%s) %s", d.Id(), errUpdate)
	// }

	// taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// // Wait for the Image to be available
	// stateConf := &resource.StateChangeConf{
	// 	Pending:    []string{"QUEUED", "RUNNING"},
	// 	Target:     []string{"SUCCEEDED"},
	// 	Refresh:    taskStateRefreshFunc(conn, taskUUID),
	// 	Timeout:    time.Duration(timeout) * time.Minute,
	// 	Delay:      imageDelay,
	// 	MinTimeout: imageMinTimeout,
	// }

	// if _, err := stateConf.WaitForState(); err != nil {
	// 	delErr := resourceNutanixFCImageClusterDelete(d, meta)
	// 	if delErr != nil {
	// 		return fmt.Errorf("error waiting for image (%s) to delete in update: %s", d.Id(), delErr)
	// 	}
	// 	uuid := d.Id()
	// 	d.SetId("")
	// 	return fmt.Errorf("error waiting for image (%s) to update: %s", uuid, err)
	// }

	// return resourceNutanixFCImageClusterRead(d, meta)
	return nil
}

func resourceNutanixFCImageClusterDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Deleting Image: %s", d.Get("name").(string))

	// client := meta.(*Client)
	// conn := client.API
	// timeout := client.WaitTimeout

	// if client.WaitTimeout == 0 {
	// 	timeout = 10
	// }

	// UUID := d.Id()

	// resp, err := conn.fc.DeleteImage(UUID)
	// if err != nil {
	// 	if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
	// 		d.SetId("")
	// 	}
	// 	return err
	// }

	// taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// // Wait for the Image to be available
	// stateConf := &resource.StateChangeConf{
	// 	Pending:    []string{"QUEUED", "RUNNING"},
	// 	Target:     []string{"SUCCEEDED"},
	// 	Refresh:    taskStateRefreshFunc(conn, taskUUID),
	// 	Timeout:    time.Duration(timeout) * time.Minute,
	// 	Delay:      imageDelay,
	// 	MinTimeout: imageMinTimeout,
	// }

	// if _, err := stateConf.WaitForState(); err != nil {
	// 	d.SetId("")
	// 	return fmt.Errorf("error waiting for image (%s) to delete: %s", d.Id(), err)
	// }

	// d.SetId("")
	return nil
}

// func resourceImageInstanceStateUpgradeV0(is map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
// 	log.Printf("[DEBUG] Entering resourceImageInstanceStateUpgradeV0")
// 	return resourceNutanixCategoriesMigrateState(is, meta)
// }

// func resourceNutanixImageInstanceResourceV0() *schema.Resource {
// 	return &schema.Resource{
//     Schema: map[string]*schema.Schema {
//       "cluster_external_ip": {
//         Type: schema.TypeString
//         Computed: true
//         Optional: true
//       },
// 			"common_network_settings": {
// 				Type:     schema.TypeMap,
// 				Computed: true,
// 				Elem: &schema.Resource{
// 					Schema: map[string]*schema.Schema{
// 						"cvm_dns_servers": {
// 							Type:     schema.TypeList,
// 							Computed: true,
// 						},
// 						"hypervisor_dns_servers": {
// 							Type:     schema.TypeList,
// 							Computed: true,
// 						},
// 						"cvm_ntp_servers": {
// 							Type:     schema.TypeList,
// 							Computed: true,
// 						},
// 						"hypervisor_ntp_servers": {
// 							Type:     schema.TypeList,
// 							Computed: true,
// 						},
// 					},
// 				},
// 			},
// 			"hypervisor_iso_details": {
// 				Type:     schema.TypeMap,
// 				Computed: true,
// 				Elem: &schema.Resource{
// 					Schema: map[string]*schema.Schema{
// 						"hyperv_sku": {
// 							Type:     schema.TypeString,
// 							Computed: true,
// 						},
// 						"url": {
// 							Type:     schema.TypeString,
// 							Computed: true,
//               Required: true,
// 						},
// 						"hyperv_product_key": {
// 							Type:     schema.TypeString,
// 							Computed: true,
// 						},
// 						"sha256sum": {
// 							Type:     schema.TypeString,
// 							Computed: true,
// 						},
// 					},
// 				},
// 			},
//       "storage_node_count": {
//         Type: schema.TypeInt
//         Computed: true
//       },
//       "redundancy_factor": {
//         Type: schema.TypeInt
//         Computed: true
//         Required: true
//       },
//       "cluster_name": {
//         Type: schema.TypeString
//         Computed: true
//       },
//       "aos_package_url": {
//         Type: schema.TypeString
//         Computed: true
//       },
//       "cluster_size": {
//         Type: schema.TypeInt
//         Computed: true
//       },
//       "aos_package_sha256sum": {
//         Type: schema.TypeString
//         Computed: true
//       },
//       "timezone": {
//         Type: schema.TypeBool
//         Computed: true
//       },
// 			"node_list": {
// 				Type:     schema.TypeMap,
// 				Computed: true,
// 				Elem: &schema.Resource{
// 					Schema: map[string]*schema.Schema{
// 						"cvm_gateway": {
// 							Type:     schema.TypeString,
// 							Computed: true,
// 						},
// 						"ipmi_netmask": {
// 							Type:     schema.TypeString,
// 							Computed: true,
//               Required: true,
// 						},
// 						"rdma_passthrough": {
// 							Type:     schema.TypeBool,
// 							Computed: true,
//               Default: false
// 						},
// 						"imaged_node_uuid": {
// 							Type:     schema.TypeString,
// 							Computed: true,
// 						},
// 						"cvm_vlan_id": {
// 							Type:     schema.TypeInt,
// 							Computed: true,
// 						},
// 						"hypervisor_type": {
// 							Type:     schema.TypeString,
// 							Computed: true,
//               Required: true,
// 						},
// 						"image_now": {
// 							Type:     schema.TypeBool,
// 							Computed: true,
// 						},
// 						"hypervisor_hostname": {
// 							Type:     schema.TypeString,
// 							Computed: true,
// 						},
// 						"hypervisor_netmask": {
// 							Type:     schema.TypeString,
// 							Computed: true,
// 						},
// 						"cvm_netmask": {
// 							Type:     schema.TypeString,
// 							Computed: true,
//               Required: true,
// 						},
// 						"ipmi_ip": {
// 							Type:     schema.TypeString,
// 							Computed: true,
// 						},
// 						"hypervisor_gateway": {
// 							Type:     schema.TypeString,
// 							Computed: true,
// 						},
// 						"hardware_attributes_override": {
// 							Type:     schema.TypeObject,
// 							Computed: true,
// 						},
// 						"cvm_ram_gb": {
// 							Type:     schema.TypeInt,
// 							Computed: true,
//               Required: true,
// 						},
// 						"cvm_ip": {
// 							Type:     schema.TypeInt,
// 							Computed: true,
// 						},
// 						"hypervisor_ip": {
// 							Type:     schema.TypeString,
// 							Computed: true,
// 						},
//             "use_existing_network_settings": {
//               Type: schema.TypeBool,
//               Computed: true,
//               Default: false
//             }
// 					},
// 				},
// 			},
//     }
//   }
// }
