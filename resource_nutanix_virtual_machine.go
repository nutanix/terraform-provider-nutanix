package nutanix

import (
	"encoding/json"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ideadevice/terraform-ahv-provider-plugin/requestutils"
	vm "github.com/ideadevice/terraform-ahv-provider-plugin/virtualmachineconfig"
	st "github.com/ideadevice/terraform-ahv-provider-plugin/virtualmachinestruct"
	"log"
	"runtime/debug"
)

type vmStruct struct {
	Metadata st.MetaDataStruct `json:"metadata"`
	Status   interface{}       `json:"status"`
	Spec     st.SpecStruct     `json:"spec"`
}

type vmList struct {
	APIVersion string            `json:"api_version"`
	MetaData   st.MetaDataStruct `json:"metadata"`
	Entities   []vmStruct        `json:"entities"`
}

func updateAddress(d *schema.ResourceData) error {
	return nil
}

// Function checks if there is an error
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// RecoverFunc can be used to recover from panics. name is the name of the caller
func RecoverFunc(name string) {
	if err := recover(); err != nil {
		log.Printf("Recovered from error %s, %s", err, name)
		log.Printf("Stack Trace: %s", debug.Stack())
		panic(err)
	}
}

// ID returns the id to be set
func (m *Machine) ID() string {
	return "ID-" + m.Spec.Name + "!!"
}

// DeleteMachine function deletes the vm using DELETE api call
func (c *MyClient) DeleteMachine(m *Machine) error {

	jsonStr := []byte(`{}`)
	url := c.Endpoint + "/list"
	method := "POST"
	jsonResponse := requestutils.RequestHandler(url, method, jsonStr, c.Username, c.Password)

	var uuid string
	var vmlist vmList
	err := json.Unmarshal(jsonResponse, &vmlist)
	check(err)

	for _, vm := range vmlist.Entities {
		if vm.Spec.Name == m.Spec.Name {
			uuid = vm.Metadata.UUID
		}
	}

	url = c.Endpoint + "/" + uuid
	method = "DELETE"
	requestutils.RequestHandler(url, method, jsonStr, c.Username, c.Password)
	return nil
}

// CreateMachine function creates the vm using POST api call
func (c *MyClient) CreateMachine(m *Machine, d *schema.ResourceData) error {

	jsonStr, err1 := json.Marshal(m)
	check(err1)

	method := "POST"
	requestutils.RequestHandler(c.Endpoint, method, jsonStr, c.Username, c.Password)
	return nil
}

func resourceNutanixVirtualMachineCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*MyClient)
	machine := Machine(vm.SetMachineConfig(d))
	machine.Spec.Name = d.Get("name").(string)
	machine.Metadata.Name = d.Get("name").(string)

	err := client.CreateMachine(&machine, d)
	if err != nil {
		return err
	}

	d.SetId(machine.ID())
	return nil

}

func resourceNutanixVirtualMachineRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceNutanixVirtualMachineUpdate(d *schema.ResourceData, m interface{}) error {
	// Enable partial state mode
	d.Partial(true)
	// checking that address has changed or not
	if d.HasChange("address") {
		//Try updating the address
		if err := updateAddress(d); err != nil {
			return err
		}
		// After updating address
		d.SetPartial("address")
	}
	// If we were to return here, before disabling patial mode below, then only "address" field would be saved

	//Disabling partial state mode. This will cause terraform to save all fields again
	d.Partial(false)

	return nil
}

func resourceNutanixVirtualMachineDelete(d *schema.ResourceData, m interface{}) error {

	client := m.(*MyClient)
	machine := Machine(vm.SetMachineConfig(d))
	machine.Spec.Name = d.Get("name").(string)
	machine.Metadata.Name = d.Get("name").(string)

	err := client.DeleteMachine(&machine)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceNutanixVirtualMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceNutanixVirtualMachineCreate,
		Read:   resourceNutanixVirtualMachineRead,
		Update: resourceNutanixVirtualMachineUpdate,
		Delete: resourceNutanixVirtualMachineDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"spec": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"backup_policy": &schema.Schema{
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"consistency_group_identifier": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"default_snapshot_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"snapshot_policy_list": &schema.Schema{
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"snapshot_schedule_list": &schema.Schema{
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"local_retention_quantity": {
																Type:     schema.TypeInt,
																Optional: true,
															},
															"remote_retention_quantity": {
																Type:     schema.TypeInt,
																Optional: true,
															},
															"snapshot_type": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"schedule": &schema.Schema{
																Type:     schema.TypeSet,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"is_suspended": {
																			Type:     schema.TypeBool,
																			Optional: true,
																		},
																		"start_time": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																		"end_time": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																		"interval_type": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																		"duration_secs": {
																			Type:     schema.TypeInt,
																			Optional: true,
																		},
																		"interval_multiple": {
																			Type:     schema.TypeInt,
																			Optional: true,
																		},
																	},
																},
															},
														},
													},
												},
												"replication_target": &schema.Schema{
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cluster_reference": &schema.Schema{
																Type:     schema.TypeSet,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"name": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																		"uuid": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																		"kind": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																	},
																},
															},
															"availability_zone_reference": &schema.Schema{
																Type:     schema.TypeSet,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"name": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																		"uuid": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																		"kind": {
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
								},
							},
						},
						"availability_zone_reference": &schema.Schema{
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"uuid": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"kind": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"cluster_reference": &schema.Schema{
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"uuid": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"kind": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"resources": &schema.Schema{
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"num_vcpus_per_socket": &schema.Schema{
										Type:     schema.TypeInt,
										Required: true,
									},
									"num_sockets": &schema.Schema{
										Type:     schema.TypeInt,
										Required: true,
									},
									"memory_size_mb": &schema.Schema{
										Type:     schema.TypeInt,
										Required: true,
									},
									"hard_clock_timezone": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"guest_os_id": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"power_state": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"parent_reference": &schema.Schema{
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"uuid": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"kind": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"guest_tools": &schema.Schema{
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"nutanix_guest_tools": &schema.Schema{
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"iso_mount_state": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"state": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"enabled_capability_list": {
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
									"guest_customization": &schema.Schema{
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"sysprep": &schema.Schema{
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"install_type": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"unattend_xml": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
												"cloud_init": &schema.Schema{
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"meta_data": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"user_data": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
											},
										},
									},
									"boot_config": &schema.Schema{
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"mac_address": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"disk_address": &schema.Schema{
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"device_index": {
																Type:     schema.TypeInt,
																Optional: true,
															},
															"adapter": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
											},
										},
									},
									"nic_list": &schema.Schema{
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip_endpoint_list": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ip": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"type": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
												"nic_type": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"mac_address": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"network_function_nic_type": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"network_function_chain_reference": &schema.Schema{
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"kind": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"name": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"uuid": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
												"subnet_reference": &schema.Schema{
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"kind": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"name": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"uuid": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
											},
										},
									},
									"gpu_list": &schema.Schema{
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"vendor": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"mode": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"device_id": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
									},
									"disk_list": &schema.Schema{
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"uuid": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"disk_size_mib": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"device_properties": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"device_type": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"disk_address": {
																Type:     schema.TypeSet,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"adapter_type": {
																			Type:     schema.TypeString,
																			Optional: true,
																		},
																		"device_index": {
																			Type:     schema.TypeInt,
																			Optional: true,
																		},
																	},
																},
															},
														},
													},
												},
												"data_source_reference": &schema.Schema{
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"kind": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"name": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"uuid": {
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
					},
				},
			},
			"api_version": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"metadata": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"last_update_time": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"kind": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"uuid": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"creation_time": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"spec_version": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"entity_version": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"categories": &schema.Schema{
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"owner_reference": &schema.Schema{
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"name": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"uuid": &schema.Schema{
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
	}
}
