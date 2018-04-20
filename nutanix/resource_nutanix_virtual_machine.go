package nutanix

import (
	"github.com/hashicorp/terraform/helper/schema"
)

// var statusCodeFilter map[int]bool
// var statusMap map[int]bool
// var version int64
// var powerON = vmconfig.PowerON
// var powerOFF = vmconfig.PowerOFF

// func init() {
// 	statusMap = map[int]bool{
// 		200: true,
// 		201: true,
// 		202: true,
// 		203: true,
// 		204: true,
// 		205: true,
// 		206: true,
// 		207: true,
// 		208: true,
// 	}
// 	statusCodeFilter = statusMap
// }

// // Function checks if there is an error
// func check(e error) {
// 	if e != nil {
// 		panic(e)
// 	}
// }

// func checkAPIResponse(resp nutanixV3.APIResponse) error {
// 	response := fmt.Sprintf("Response ==> %+v\n Response Message ==> %+v\n Request ==> %+v\n Request Body==> %+v", resp.Response, resp.Message, resp.Response.Request, resp.Response.Request.Body)
// 	if flg.HTTPLog != "" {
// 		file, err := os.Create(flg.HTTPLog)
// 		if err != nil {
// 			return err
// 		}
// 		w := bufio.NewWriter(file)
// 		defer file.Close()
// 		defer w.Flush()
// 		fmt.Fprintf(w, "%v", response)
// 	}
// 	if !statusCodeFilter[resp.StatusCode] {
// 		errormsg := errors.New(response)
// 		return errormsg
// 	}
// 	return nil
// }

// // RecoverFunc can be used to recover from panics. name is the name of the caller
// func RecoverFunc(name string) {
// 	if err := recover(); err != nil {
// 		log.Printf("Recovered from error %s, %s", err, name)
// 		log.Printf("Stack Trace: %s", debug.Stack())
// 		panic(err)
// 	}
// }

// // setAPIInstance sets the nutanixV3.VmApi from the V3Client
// func setAPIInstance(c *V3Client) *(nutanixV3.VmsApi) {
// 	APIInstance := nutanixV3.NewVmsApi()
// 	APIInstance.Configuration.Username = c.Username
// 	APIInstance.Configuration.Password = c.Password
// 	APIInstance.Configuration.BasePath = c.URL
// 	APIInstance.Configuration.APIClient.Insecure = c.Insecure
// 	return APIInstance
// }

// // WaitForProcess waits till the nutanix gets to running
// func (c *V3Client) WaitForProcess(uuid string) (bool, error) {
// 	APIInstance := setAPIInstance(c)
// 	for {
// 		VMIntentResponse, APIresponse, err := APIInstance.VmsUuidGet(uuid)
// 		if err != nil {
// 			return false, err
// 		}
// 		err = checkAPIResponse(*APIresponse)
// 		if err != nil {
// 			return false, err
// 		}
// 		if VMIntentResponse.Status.State == "COMPLETE" {
// 			return true, nil
// 		} else if VMIntentResponse.Status.State == "ERROR" {
// 			return false, fmt.Errorf("Error while waiting for resource to be up")
// 		}
// 		time.Sleep(3000 * time.Millisecond)
// 	}
// 	return false, nil
// }

// // WaitForIP function sets the ip address obtained by the GET request
// func (c *V3Client) WaitForIP(uuid string, d *schema.ResourceData) error {
// 	APIInstance := setAPIInstance(c)
// 	for {
// 		VMIntentResponse, APIresponse, err := APIInstance.VmsUuidGet(uuid)
// 		if err != nil {
// 			return err
// 		}
// 		err = checkAPIResponse(*APIresponse)
// 		if err != nil {
// 			return err
// 		}
// 		if len(VMIntentResponse.Status.Resources.NicList) != 0 {
// 			for i := range VMIntentResponse.Status.Resources.NicList {
// 				if len(VMIntentResponse.Status.Resources.NicList[i].IpEndpointList) != 0 {
// 					if ip := VMIntentResponse.Status.Resources.NicList[i].IpEndpointList[0].Ip; ip != "" {
// 						d.Set("ip_address", ip)
// 						return nil
// 					}
// 				}
// 			}
// 		}
// 		time.Sleep(3000 * time.Millisecond)
// 	}
// 	return nil
// }

// func resourceNutanixVirtualMachineCreate(d *schema.ResourceData, meta interface{}) error {
// 	client := meta.(*V3Client)
// 	machine := vmconfig.SetMachineConfig(d)
// 	log.Printf("[DEBUG] Creating Virtual Machine: %s", d.Id())
// 	APIInstance := setAPIInstance(client)
// 	VMIntentResponse, APIResponse, err := APIInstance.VmsPost(machine)
// 	if err != nil {
// 		return err
// 	}

// 	err = checkAPIResponse(*APIResponse)
// 	if err != nil {
// 		return err
// 	}
// 	uuid := VMIntentResponse.Metadata.Uuid
// 	status, err := client.WaitForProcess(uuid)
// 	for status != true {
// 		return err
// 	}
// 	d.Set("ip_address", "")

// 	if machine.Spec.Resources.NicList != nil && machine.Spec.Resources.PowerState == powerON {
// 		log.Printf("[DEBUG] Polling for IP\n")
// 		err = client.WaitForIP(uuid, d)
// 	}
// 	if err != nil {
// 		return err
// 	}

// 	d.SetId(uuid)
// 	return nil

// }

// func resourceNutanixVirtualMachineRead(d *schema.ResourceData, meta interface{}) error {
// 	client := meta.(*V3Client)
// 	APIInstance := setAPIInstance(client)
// 	VMIntentResponse, APIResponse, err := APIInstance.VmsUuidGet(d.Id())
// 	log.Printf("[DEBUG] Syncing the remote Virtual Machine instance with local instance: %s, %s", VMIntentResponse.Spec.Name, d.Id())
// 	if err != nil {
// 		return err
// 	}
// 	machine := vmconfig.SetMachineConfig(d)

// 	err = checkAPIResponse(*APIResponse)
// 	if err != nil {
// 		return err
// 	}

// 	VMIntentResponse.Spec.Resources = vmconfig.GetVMResources(VMIntentResponse.Status.Resources)

// 	machineTemp := nutanixV3.VmIntentInput{
// 		ApiVersion: "3.0",
// 		Spec:       VMIntentResponse.Spec,
// 		Metadata:   VMIntentResponse.Metadata,
// 	}

// 	if len(machineTemp.Spec.Resources.DiskList) == len(machine.Spec.Resources.DiskList) {
// 		machineTemp.Spec.Resources.DiskList = machine.Spec.Resources.DiskList
// 	}
// 	if len(machineTemp.Spec.Resources.NicList) == len(machine.Spec.Resources.NicList) {
// 		machineTemp.Spec.Resources.NicList = machine.Spec.Resources.NicList
// 	}
// 	machineTemp.Metadata.OwnerReference = machine.Metadata.OwnerReference
// 	machineTemp.Metadata.Uuid = machine.Metadata.Uuid
// 	machineTemp.Metadata.Name = machine.Metadata.Name

// 	if !reflect.DeepEqual(machineTemp, machine) {
// 		err = vmconfig.UpdateTerraformState(d, VMIntentResponse.Metadata, VMIntentResponse.Spec)
// 		if err != nil {
// 			return err
// 		}
// 		d.Set("ip_address", "")
// 		if len(VMIntentResponse.Spec.Resources.NicList) > 0 && VMIntentResponse.Spec.Resources.PowerState == powerON {
// 			err = client.WaitForIP(d.Id(), d)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		version = VMIntentResponse.Metadata.SpecVersion

// 	}

// 	return nil
// }

// func resourceNutanixVirtualMachineUpdate(d *schema.ResourceData, meta interface{}) error {
// 	// Enable partial state mode
// 	d.Partial(true)
// 	client := meta.(*V3Client)
// 	machine := vmconfig.SetMachineConfig(d)
// 	machine.Metadata.SpecVersion = version

// 	APIInstance := setAPIInstance(client)
// 	uuid := d.Id()
// 	log.Printf("[DEBUG] Updating Virtual Machine: %s, %s", machine.Spec.Name, d.Id())

// 	if d.HasChange("name") || d.HasChange("spec") || d.HasChange("metadata") {
// 		_, APIResponse, err := APIInstance.VmsUuidPut(uuid, machine)
// 		if err != nil {
// 			return err
// 		}
// 		err = checkAPIResponse(*APIResponse)
// 		if err != nil {
// 			return err
// 		}
// 		d.SetPartial("spec")
// 		d.SetPartial("metadata")
// 	}
// 	//Disabling partial state mode. This will cause terraform to save all fields again
// 	d.Partial(false)
// 	status, err := client.WaitForProcess(uuid)
// 	if status != true {
// 		return err
// 	}
// 	d.Set("ip_address", "")
// 	if len(machine.Spec.Resources.NicList) > 0 && machine.Spec.Resources.PowerState == powerON {
// 		log.Printf("[DEBUG] Polling for IP\n")
// 		err := client.WaitForIP(uuid, d)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func resourceNutanixVirtualMachineDelete(d *schema.ResourceData, m interface{}) error {
// 	client := m.(*V3Client)
// 	log.Printf("[DEBUG] Deleting Virtual Machine: %s", d.Id())
// 	APIInstance := setAPIInstance(client)
// 	uuid := d.Id()

// 	APIResponse, err := APIInstance.VmsUuidDelete(uuid)
// 	if err != nil {
// 		return err
// 	}
// 	err = checkAPIResponse(*APIResponse)
// 	if err != nil {
// 		return err
// 	}

// 	d.SetId("")
// 	return nil
// }

// // MachineExists function returns the uuid of the machine with given name
// func resourceNutanixVirtualMachineExists(d *schema.ResourceData, m interface{}) (bool, error) {
// 	log.Printf("[DEBUG] Checking Virtual Machine Existance: %s", d.Id())
// 	client := m.(*V3Client)
// 	APIInstance := setAPIInstance(client)

// 	getEntitiesRequest := nutanixV3.VmListMetadata{} // VmListMetadata
// 	VMListIntentResponse, APIResponse, err := APIInstance.VmsListPost(getEntitiesRequest)
// 	if err != nil {
// 		return false, err
// 	}
// 	err = checkAPIResponse(*APIResponse)
// 	if err != nil {
// 		return false, err
// 	}

// 	for i := range VMListIntentResponse.Entities {
// 		if VMListIntentResponse.Entities[i].Metadata.Uuid == d.Id() {
// 			return true, nil
// 		}
// 	}
// 	return false, nil
// }
// }

func resourceNutanixVirtualMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceNutanixVirtualMachineCreate,
		Read:   resourceNutanixVirtualMachineRead,
		Update: resourceNutanixVirtualMachineUpdate,
		Delete: resourceNutanixVirtualMachineDelete,
		Exists: resourceNutanixVirtualMachineExists,

		Schema: map[string]*schema.Schema{
			"status": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"state": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"availability_zone_reference": &schema.Schema{
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"uuid": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"name": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"message_list": &schema.Schema{
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							"message": &schema.Schema{
								Type:     schema.TypeString,
								Computed: true,
							},
							"reason": &schema.Schema{
								Type:     schema.TypeString,
								Computed: true,
							},
							"details": &schema.Schema{
								Type:     schema.TypeMap,
								Computed: true,
							},
						},
					},
					"cluster_reference": &schema.Schema{
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"uuid": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"name": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"resources": &schema.Schema{
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"vnuma_config": &schema.Schema{
									Type:     schema.TypeMap,
									Computed: true,
									Elem: &schema.Resource{
										Schema: &schema.Schema{
											"num_vnuma_nodes": &schema.Schema{
												Type:     schema.TypeInt,
												Computed: true,
											},
										},
									},
								},
								"nic_list": &schema.Schema{
									Type:     schema.TypeList,
									Computed: true,
									Elem: &schema.Resource{
										Schema: &schema.Schema{
											"nic_type": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"uuid": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"ip_endpoint_list": &schema.Schema{
												Type:     schema.TypeList,
												Computed: true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"ip": &schema.Schema{
															Type:     schema.TypeString,
															Computed: true,
														},
														"type": &schema.Schema{
															Type:     schema.TypeString,
															Computed: true,
														},
													},
												},
											},
											"network_function_chain_reference": &schema.Schema{
												Type:     schema.TypeMap,
												Computed: true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"kind": &schema.Schema{
															Type:     schema.TypeString,
															Computed: true,
														},
														"name": &schema.Schema{
															Type:     schema.TypeString,
															Computed: true,
														},
														"uuid": &schema.Schema{
															Type:     schema.TypeString,
															Computed: true,
														},
													},
												},
											},
											"floating_ip": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"network_function_nic_type": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"mac_address": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"subnet_reference": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"kind": &schema.Schema{
															Type:     schema.TypeString,
															Computed: true,
														},
														"name": &schema.Schema{
															Type:     schema.TypeString,
															Computed: true,
														},
														"uuid": &schema.Schema{
															Type:     schema.TypeString,
															Computed: true,
														},
													},
												},
											},
											"model": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
								},
								"host_reference": &schema.Schema{
									Type:     schema.TypeMap,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"kind": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"name": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"uuid": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
								},
								"guest_os_id": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"power_state": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"guest_tools": &schema.Schema{
									Type:     schema.TypeMap,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"nutanix_guest_tools": &schema.Schema{
												Type:     schema.TypeMap,
												Computed: true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"available_version": &schema.Schema{
															Type:     schema.TypeString,
															Computed: true,
														},
														"iso_mount_state": &schema.Schema{
															Type:     schema.TypeString,
															Computed: true,
														},
														"state": &schema.Schema{
															Type:     schema.TypeString,
															Computed: true,
														},
														"version": &schema.Schema{
															Type:     schema.TypeString,
															Computed: true,
														},
														"guest_os_version": &schema.Schema{
															Type:     schema.TypeString,
															Computed: true,
														},
														"enabled_capability_list": &schema.Schema{
															Type:     schema.TypeList,
															Computed: true,
															Elem:     &schema.Schema{Type: schema.TypeString},
														},
														"vss_snapshot_capable": &schema.Schema{
															Type:     schema.TypeBool,
															Computed: true,
														},
														"is_reachable": &schema.Schema{
															Type:     schema.TypeBool,
															Computed: true,
														},
														"vm_mobility_drivers_installed": &schema.Schema{
															Type:     schema.TypeBool,
															Computed: true,
														},
													},
												},
											},
										},
									},
								},
								"hypervisor_type": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"num_vcpus_per_socket": &schema.Schema{
									Type:     schema.TypeInt,
									Computed: true,
								},
								"num_sockets": &schema.Schema{
									Type:     schema.TypeInt,
									Computed: true,
								},
								"gpu_list": &schema.Schema{
									Type:     schema.TypeList,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"frame_buffer_size_mib": &schema.Schema{
												Type:     schema.TypeInt,
												Computed: true,
											},
											"vendor": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"uuid": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"name": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"pci_address": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"fraction": &schema.Schema{
												Type:     schema.TypeInt,
												Computed: true,
											},
											"mode": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"num_virtual_display_heads": &schema.Schema{
												Type:     schema.TypeInt,
												Computed: true,
											},
											"guest_driver_version": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"device_id": &schema.Schema{
												Type:     schema.TypeInt,
												Computed: true,
											},
										},
									},
								},
								"parent_reference": &schema.Schema{
									Type:     schema.TypeMap,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"kind": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"name": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"uuid": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
								},
								"memory_size_mib": &schema.Schema{
									Type:     schema.TypeInt,
									Computed: true,
								},
								"boot_config": &schema.Schema{
									Type:     schema.TypeMap,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"boot_device_order_list": &schema.Schema{
												Type:     schema.TypeList,
												Computed: true,
												Elem:     &schema.Schema{Type: schema.TypeString},
											},
											"boot_device": &schema.Schema{
												Type:     schema.TypeMap,
												Computed: true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"disk_address": &schema.Schema{
															Type:     schema.TypeMap,
															Computed: true,
															Elem: &schema.Resource{
																Schema: map[string]*schema.Schema{
																	"device_index": &schema.Schema{
																		Type:     schema.TypeInt,
																		Computed: true,
																	},
																	"adapter_type": &schema.Schema{
																		Type:     schema.TypeString,
																		Computed: true,
																	},
																},
															},
														},
														"mac_address": &schema.Schema{
															Type:     schema.TypeString,
															Computed: true,
														},
													},
												},
											},
										},
									},
								},
								"hardware_clock_timezone": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"guest_customization": &schema.Schema{
									Type:     schema.TypeMap,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"cloud_init": &schema.Schema{
												Type:     schema.TypeMap,
												Computed: true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"meta_data": &schema.Schema{
															Type:     schema.TypeString,
															Computed: true,
														},
														"user_data": &schema.Schema{
															Type:     schema.TypeString,
															Computed: true,
														},
														"custom_key_values": &schema.Schema{
															Type:     schema.TypeMap,
															Computed: true,
														},
													},
												},
											},
											"is_overridable": &schema.Schema{
												Type:     schema.TypeBool,
												Computed: true,
											},
											"sysprep": &schema.Schema{
												Type:     schema.TypeMap,
												Computed: true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"install_type": &schema.Schema{
															Type:     schema.TypeString,
															Computed: true,
														},
														"unattend_xml": &schema.Schema{
															Type:     schema.TypeString,
															Computed: true,
														},
														"custom_key_values": &schema.Schema{
															Type:     schema.TypeMap,
															Computed: true,
														},
													},
												},
											},
										},
									},
								},
								"power_state_mechanism": &schema.Schema{
									Type:     schema.TypeMap,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"guest_transition_config": &schema.Schema{
												Type:     schema.TypeMap,
												Computed: true,
												Elem: &schema.Resource{
													"should_fail_on_script_failure": &schema.Schema{
														Type:     schema.TypeBool,
														Computed: true,
													},
													"enable_script_exec": &schema.Schema{
														Type:     schema.TypeBool,
														Computed: true,
													},
												},
											},
											"mechanism": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
										},
									},
								},
								"vga_console_enabled": &schema.Schema{
									Type:     schema.TypeBool,
									Computed: true,
								},
								"disk_list": &schema.Schema{
									Type:     schema.TypeList,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"uuid": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"disk_size_bytes": &schema.Schema{
												Type:     schema.TypeInt,
												Computed: true,
											},
											"device_properties": &schema.Schema{
												Type:     schema.TypeMap,
												Computed: true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"device_type": &schema.Schema{
															Type:     schema.TypeString,
															Computed: true,
														},
														"disk_address": &schema.Schema{
															Type:     schema.TypeMap,
															Computed: true,
															Elem: &schema.Resource{
																"device_index": &schema.Schema{
																	Type:     schema.TypeInt,
																	Computed: true,
																},
																"adapter_type": &schema.Schema{
																	Type:     schema.TypeString,
																	Computed: true,
																},
															},
														},
													},
												},
											},
											"data_source_reference": &schema.Schema{
												Type:     schema.TypeMap,
												Computed: true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"kind": &schema.Schema{
															Type:     schema.TypeString,
															Computed: true,
														},
														"name": &schema.Schema{
															Type:     schema.TypeString,
															Computed: true,
														},
														"uuid": &schema.Schema{
															Type:     schema.TypeString,
															Computed: true,
														},
													},
												},
											},
											"disk_size_mib": &schema.Schema{
												Type:     schema.TypeInt,
												Computed: true,
											},
											"volume_group_reference": &schema.Schema{
												Type:     schema.TypeMap,
												Computed: true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"kind": &schema.Schema{
															Type:     schema.TypeString,
															Computed: true,
														},
														"name": &schema.Schema{
															Type:     schema.TypeString,
															Computed: true,
														},
														"uuid": &schema.Schema{
															Type:     schema.TypeString,
															Computed: true,
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					"description": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
			"spec": &schema.Schema{
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vm": &schema.Schema{
							Type:     schema.TypeMap,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"availability_zone_reference": &schema.Schema{
										Type:     schema.TypeMap,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"kind": &schema.Schema{
													Type:     schema.TypeString,
													Required: true,
												},
												"name": &schema.Schema{
													Type:     schema.TypeString,
													Optional: true,
												},
												"uuid": &schema.Schema{
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
									"description": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"resources": &schema.Schema{
										Type:     schema.TypeMap,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"vnuma_config": &schema.Schema{
													Type:     schema.TypeMap,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"num_vnuma_nodes": &schema.Schema{
																Type:     schema.TypeInt,
																Optional: true,
															},
														},
													},
												},
												"nic_list": &schema.Schema{
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"nic_type": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
															"uuid": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
															"ip_endpoint_list": &schema.Schema{
																Type:     schema.TypeList,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"ip": &schema.Schema{
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																		"type": &schema.Schema{
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																	},
																},
															},
															"network_function_chain_reference": &schema.Schema{
																Type:     schema.TypeMap,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"kind": &schema.Schema{
																			Type:     schema.TypeString,
																			Required: true,
																		},
																		"name": &schema.Schema{
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																		"uuid": &schema.Schema{
																			Type:     schema.TypeString,
																			Required: true,
																		},
																	},
																},
															},
															"network_function_nic_type": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
															"mac_address": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
															"subnet_reference": &schema.Schema{
																Type:     schema.TypeMap,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"kind": &schema.Schema{
																			Type:     schema.TypeString,
																			Required: true,
																		},
																		"name": &schema.Schema{
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																		"uuid": &schema.Schema{
																			Type:     schema.TypeString,
																			Required: true,
																		},
																	},
																},
															},
															"model": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
												"guest_os_id": &schema.Schema{
													Type:     schema.TypeString,
													Optional: true,
												},
												"power_state": &schema.Schema{
													Type:     schema.TypeString,
													Optional: true,
												},
												"guest_tools": &schema.Schema{
													Type:     schema.TypeMap,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"nutanix_guest_tools": &schema.Schema{
																Type:     schema.TypeMap,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"iso_mount_state": &schema.Schema{
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																		"state": &schema.Schema{
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																		"enabled_capability_list": &schema.Schema{
																			Type:     schema.TypeList,
																			Optional: true,
																			Elem:     &schema.Schema{Type: schema.TypeString},
																		},
																	},
																},
															},
														},
													},
												},
												"num_vcpus_per_socket": &schema.Schema{
													Type:     schema.TypeInt,
													Optional: true,
												},
												"num_sockets": &schema.Schema{
													Type:     schema.TypeInt,
													Optional: true,
												},
												"gpu_list": &schema.Schema{
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"vendor": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
															"mode": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
															"device_id": &schema.Schema{
																Type:     schema.TypeInt,
																Optional: true,
															},
														},
													},
												},
												"parent_reference": &schema.Schema{
													Type:     schema.TypeMap,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"kind": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
															"uuid": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
															"name": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
												"memory_size_mib": &schema.Schema{
													Type:     schema.TypeInt,
													Optional: true,
												},
												"boot_config": &schema.Schema{
													Type:     schema.TypeMap,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"boot_device_order_list": &schema.Schema{
																Type:     schema.TypeList,
																Optional: true,
																Elem:     &schema.Schema{Type: schema.TypeString},
															},
															"boot_device": &schema.Schema{
																Type:     schema.TypeMap,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"disk_address": &schema.Schema{
																			Type:     schema.TypeMap,
																			Optional: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"device_index": &schema.Schema{
																						Type:     schema.TypeInt,
																						Optional: true,
																					},
																					"adapter_type": &schema.Schema{
																						Type:     schema.TypeString,
																						Optional: true,
																					},
																				},
																			},
																		},
																		"mac_address": &schema.Schema{
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																	},
																},
															},
														},
													},
												},
												"hardware_clock_timezone": &schema.Schema{
													Type:     schema.TypeString,
													Optional: true,
												},
												"guest_customization": &schema.Schema{
													Type:     schema.TypeMap,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cloud_init": &schema.Schema{
																Type:     schema.TypeMap,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"meta_data": &schema.Schema{
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																		"user_data": &schema.Schema{
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																		"custom_key_values": &schema.Schema{
																			Type:     schema.TypeMap,
																			Optional: true,
																		},
																	},
																},
															},
															"is_overridable": &schema.Schema{
																Type:     schema.TypeBool,
																Optional: true,
															},
															"sysprep": &schema.Schema{
																Type:     schema.TypeMap,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"install_type": &schema.Schema{
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																		"unattend_xml": &schema.Schema{
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																		"custom_key_values": &schema.Schema{
																			Type:     schema.TypeMap,
																			Optional: true,
																		},
																	},
																},
															},
														},
													},
												},
												"power_state_mechanism": &schema.Schema{
													Type:     schema.TypeMap,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"guest_transition_config": &schema.Schema{
																Type:     schema.TypeMap,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"should_fail_on_script_failure": &schema.Schema{
																			Type:     schema.TypeBool,
																			Optional: true,
																		},
																		"enable_script_exec": &schema.Schema{
																			Type:     schema.TypeBool,
																			Optional: true,
																		},
																	},
																},
															},
															"mechanism": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
												"vga_console_enabled": &schema.Schema{
													Type:     schema.TypeBool,
													Optional: true,
												},
												"disk_list": &schema.Schema{
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"uuid": &schema.Schema{
																Type:     schema.TypeString,
																Optional: true,
															},
															"disk_size_bytes": &schema.Schema{
																Type:     schema.TypeInt,
																Optional: true,
															},
															"device_properties": &schema.Schema{
																Type:     schema.TypeMap,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"device_type": &schema.Schema{
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																		"disk_address": &schema.Schema{
																			Type:     schema.TypeMap,
																			Optional: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"device_index": &schema.Schema{
																						Type:     schema.TypeInt,
																						Optional: true,
																					},
																					"adapter_type": &schema.Schema{
																						Type:     schema.TypeString,
																						Optional: true,
																					},
																				},
																			},
																		},
																	},
																},
															},
															"data_source_reference": &schema.Schema{
																Type:     schema.TypeMap,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"kind": &schema.Schema{
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																		"uuid": &schema.Schema{
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																		"name": &schema.Schema{
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																	},
																},
															},
															"disk_size_mib": &schema.Schema{
																Type:     schema.TypeInt,
																Optional: true,
															},
															"volume_group_reference": &schema.Schema{
																Type:     schema.TypeMap,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"kind": &schema.Schema{
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																		"uuid": &schema.Schema{
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																		"name": &schema.Schema{
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
									},
									"cluster_reference": &schema.Schema{
										Type:     schema.TypeMap,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"kind": &schema.Schema{
													Type:     schema.TypeString,
													Optional: true,
												},
												"uuid": &schema.Schema{
													Type:     schema.TypeString,
													Optional: true,
												},
												"name": &schema.Schema{
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
			},
			"api_version": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"metadata": &schema.Schema{
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"last_update_time": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"kind": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"uuid": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"project_reference": &schema.Schema{
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"uuid": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"name": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"creation_time": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"spec_version": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"spec_hash": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"owner_reference": &schema.Schema{
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"uuid": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"name": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"categories": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}
