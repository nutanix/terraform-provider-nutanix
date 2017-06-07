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

// HostReferenceStruct struct
type HostReferenceStruct struct {
	Kind string `json:"kind"`
	UUID string `json:"uuid"`
}

// StatusStruct has status of VM
type StatusStruct struct {
	State             string               `json:"state,omitempty"`
	Name              string               `json:"name,omitempty"`
	Resources         *st.ResourcesStruct  `json:"resources,omitempty"`
	HostReference     *HostReferenceStruct `json:"host_reference,omitempty"`
	HypervisorType    string               `json:"hypervisor_type",omitempty`
	NumVcpusPerSocket int                  `json:"num_vcpus_per_socket,omitempty"`
	NumSockets        int                  `json:"num_sockets,omitempty"`
	MemorySizeMb      int                  `json:"memory_size_mb,omitempty"`
	GpuList           []string             `json:"gpu_list,omitempty"`
	PowerState        string               `json:"power_state,omitempty"`
}

// VMResponse is struct returned by Post call for creating vm
type VMResponse struct {
	Status     *StatusStruct      `json:"status"`
	Spec       *st.SpecStruct     `json:"spec,omitempty"`
	APIVersion string             `json:"api_version",omitempty`
	Metadata   *st.MetaDataStruct `json:"metadata,omitempty"`
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
	url := "https://" + c.Endpoint + ":9440/api/nutanix/v3/vms/list"
	method := "POST"
	jsonResponse, err := requestutils.RequestHandler(url, method, jsonStr, c.Username, c.Password, c.Insecure)
	if err != nil {
		return err
	}

	var uuid string
	var vmlist vmList
	err = json.Unmarshal(jsonResponse, &vmlist)
	check(err)

	for _, vm := range vmlist.Entities {
		if vm.Spec.Name == m.Spec.Name {
			uuid = vm.Metadata.UUID
		}
	}

	url = "https://" + c.Endpoint + ":9440/api/nutanix/v3/vms/" + uuid
	method = "DELETE"
	jsonResponse, err = requestutils.RequestHandler(url, method, jsonStr, c.Username, c.Password, c.Insecure)
	if err != nil {
		return err
	}
	return nil
}

// UpdateMachine function updates the vm specifications using PUT api call
func (c *MyClient) UpdateMachine(m *Machine, name string) error {

	jsonStr := []byte(`{}`)
	url := "https://" + c.Endpoint + ":9440/api/nutanix/v3/vms/list"
	method := "POST"
	jsonResponse, err := requestutils.RequestHandler(url, method, jsonStr, c.Username, c.Password, c.Insecure)
	if err != nil {
		return err
	}

	var uuid string
	var vmlist vmList
	err = json.Unmarshal(jsonResponse, &vmlist)
	check(err)

	for _, vm := range vmlist.Entities {
		if vm.Spec.Name == name {
			uuid = vm.Metadata.UUID
		}
	}
	jsonStr, err = json.Marshal(m)
	check(err)

	url = "https://" + c.Endpoint + ":9440/api/nutanix/v3/vms/" + uuid
	method = "PUT"
	jsonResponse, err = requestutils.RequestHandler(url, method, jsonStr, c.Username, c.Password, c.Insecure)
	if err != nil {
		return err
	}
	return nil
}

// WaitForProcess waits till the nutanix gets to running
func (c *MyClient) WaitForProcess(vmresp1 *VMResponse) (bool, error) {
	uuid := vmresp1.Metadata.UUID
	url := "https://" + c.Endpoint + ":9440/api/nutanix/v3/vms/" + uuid
	method := "GET"
	var vmresp VMResponse
	var payload []byte
	for {
		resp, err := requestutils.RequestHandler(url, method, payload, c.Username, c.Password, c.Insecure)
		if err != nil {
			return false, err
		}
		json.Unmarshal(resp, &vmresp)

		if vmresp.Status.State == "COMPLETE" {
			return true, nil
		}
	}
	return false, nil
}

// WaitForIP function sets the ip address obtained by the GET request
func (c *MyClient) WaitForIP(vmresp *VMResponse, d *schema.ResourceData) error {
	uuid := vmresp.Metadata.UUID
	url := "https://" + c.Endpoint + ":9440/api/nutanix/v3/vms/" + uuid
	method := "GET"
	var payload []byte

	for {
		resp, err := requestutils.RequestHandler(url, method, payload, c.Username, c.Password, c.Insecure)
		if err != nil {
			return err
		}
		var vmresp VMResponse
		json.Unmarshal(resp, &vmresp)

		if len(vmresp.Status.Resources.NicList) != 0 {
			for _, nic := range vmresp.Status.Resources.NicList {
				if len(nic.IPEndpointList) != 0 {
					if ip := nic.IPEndpointList[0].Address; ip != "" {
						d.Set("ip_address", ip)
						return nil
					}
				}
			}
		}
	}
	return nil
}

// CreateMachine function creates the vm using POST api call
func (c *MyClient) CreateMachine(m *Machine) ([]byte, error) {

	payload, err := json.Marshal(m)
	check(err)

	method := "POST"
	url := "https://" + c.Endpoint + ":9440/api/nutanix/v3/vms"
	resp, err := requestutils.RequestHandler(url, method, payload, c.Username, c.Password, c.Insecure)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func resourceNutanixVirtualMachineCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*MyClient)
	machine := Machine(vm.SetMachineConfig(d))
	machine.Spec.Name = d.Get("name").(string)
	machine.Metadata.Name = d.Get("name").(string)
	log.Printf("[DEBUG] Creating Virtual Machine: %s", machine.ID())

	resp, err := client.CreateMachine(&machine)
	if err != nil {
		return err
	}
	var vmresp VMResponse
	json.Unmarshal(resp, &vmresp)

	status, err := client.WaitForProcess(&vmresp)
	if status != true {
		return err
	}

	err = client.WaitForIP(&vmresp, d)
	if err != nil {
		return err
	}

	d.SetId(machine.ID())
	return nil

}

func resourceNutanixVirtualMachineRead(d *schema.ResourceData, m interface{}) error {

	return nil
}

func resourceNutanixVirtualMachineUpdate(d *schema.ResourceData, meta interface{}) error {
	// Enable partial state mode
	d.Partial(true)
	if d.HasChange("spec") || d.HasChange("metadata") {
		client := meta.(*MyClient)
		machine := Machine(vm.SetMachineConfig(d))
		machine.Metadata.Name = d.Get("name").(string)
		machine.Spec.Name = d.Get("name").(string)
		log.Printf("[DEBUG] Updating Virtual Machine: %s", d.Id())

		name := d.Get("name")
		if d.HasChange("name") {
			name, _ = d.GetChange("name")
		}

		err := client.UpdateMachine(&machine, name.(string))
		if err != nil {
			return err
		}
		d.SetPartial("spec")
		d.SetPartial("metadata")
	}
	//Disabling partial state mode. This will cause terraform to save all fields again
	d.Partial(false)

	return nil
}

func resourceNutanixVirtualMachineDelete(d *schema.ResourceData, m interface{}) error {

	client := m.(*MyClient)
	log.Printf("[DEBUG] Deleting Virtual Machine: %s", d.Id())
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
			"ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
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
