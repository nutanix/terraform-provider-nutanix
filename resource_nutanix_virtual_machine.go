package nutanix

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ideadevice/terraform-ahv-provider-plugin/requestutils"
	vmdefn "github.com/ideadevice/terraform-ahv-provider-plugin/virtualmachine"
	vm "github.com/ideadevice/terraform-ahv-provider-plugin/virtualmachineconfig"
	"log"
	"runtime/debug"
)

type vmStruct struct {
	Metadata vmdefn.MetaDataStruct `json:"metadata"`
	Status   interface{}           `json:"status"`
	Spec     vmdefn.SpecStruct     `json:"spec"`
}

type vmList struct {
	APIVersion string                `json:"api_version"`
	MetaData   vmdefn.MetaDataStruct `json:"metadata"`
	Entities   []vmStruct            `json:"entities"`
}

// HostReferenceStruct struct
type HostReferenceStruct struct {
	Kind string `json:"kind"`
	UUID string `json:"uuid"`
}

// StatusStruct has status of VM
type StatusStruct struct {
	State             string                  `json:"state,omitempty"`
	Name              string                  `json:"name,omitempty"`
	Resources         *vmdefn.ResourcesStruct `json:"resources,omitempty"`
	HostReference     *HostReferenceStruct    `json:"host_reference,omitempty"`
	HypervisorType    string                  `json:"hypervisor_type",omitempty`
	NumVcpusPerSocket int                     `json:"num_vcpus_per_socket,omitempty"`
	NumSockets        int                     `json:"num_sockets,omitempty"`
	MemorySizeMb      int                     `json:"memory_size_mb,omitempty"`
	GpuList           []string                `json:"gpu_list,omitempty"`
	PowerState        string                  `json:"power_state,omitempty"`
}

// VMResponse is struct returned by Post call for creating vm
type VMResponse struct {
	Status     *StatusStruct          `json:"status"`
	Spec       *vmdefn.SpecStruct     `json:"spec,omitempty"`
	APIVersion string                 `json:"api_version",omitempty`
	Metadata   *vmdefn.MetaDataStruct `json:"metadata,omitempty"`
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

// DeleteMachine function deletes the vm using DELETE api call
func (c *V3Client) DeleteMachine(m *vmdefn.VirtualMachine, d *schema.ResourceData) error {

	log.Printf("[DEBUG] Updating Virtual Machine: %s", m.Spec.Name)
	var jsonStr []byte
	url := c.URL + "/vms/" + d.Id()
	method := "DELETE"
	_, err := requestutils.RequestHandler(url, method, jsonStr, c.Username, c.Password, c.Insecure)
	if err != nil {
		return err
	}
	return nil
}

// UpdateMachine function updates the vm specifications using PUT api call
func (c *V3Client) UpdateMachine(m *vmdefn.VirtualMachine, d *schema.ResourceData) error {

	log.Printf("[DEBUG] Updating Virtual Machine: %s", m.Spec.Name)

	jsonStr, err := json.Marshal(m)
	check(err)

	url := c.URL + "/vms/" + d.Id()
	method := "PUT"
	_, err = requestutils.RequestHandler(url, method, jsonStr, c.Username, c.Password, c.Insecure)
	if err != nil {
		return err
	}
	return nil
}

// MachineExists function returns the uuid of the machine with given name
func (c *V3Client) MachineExists(name string) (string, error) {
	log.Printf("[DEBUG] Checking Virtual Machine Existance: %s", name)
	payload := []byte(`{}`)
	url := c.URL + "/vms/list"
	method := "POST"
	jsonResponse, err := requestutils.RequestHandler(url, method, payload, c.Username, c.Password, c.Insecure)
	if err != nil {
		return "", err
	}

	var uuid string
	var vmlist vmList
	err = json.Unmarshal(jsonResponse, &vmlist)
	check(err)

	for _, vm := range vmlist.Entities {
		if vm.Spec.Name == name {
			uuid = vm.Metadata.UUID
			return uuid, nil
		}
	}
	return "", nil
}

// ShutdownMachine function shut vm using PUT api call
func (c *V3Client) ShutdownMachine(m *vmdefn.VirtualMachine, d *schema.ResourceData) error {

	log.Printf("[DEBUG] Shutting Down Virtual Machine: %s", m.Metadata.Name)

	data := &vmdefn.VirtualMachine{
		Spec: &vmdefn.SpecStruct{
			Name: m.Spec.Name,
			Resources: &vmdefn.ResourcesStruct{
				PowerState: "POWERED_OFF",
			},
		},
		APIVersion: "3.0",
		Metadata: &vmdefn.MetaDataStruct{
			Name:        m.Spec.Name,
			Kind:        "vm",
			SpecVersion: 0,
		},
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	url := c.URL + "/vms/" + d.Id()
	method := "PUT"
	_, err = requestutils.RequestHandler(url, method, payload, c.Username, c.Password, c.Insecure)
	return err
}

// StartMachine function starts the vm using PUT api call
func (c *V3Client) StartMachine(m *vmdefn.VirtualMachine, d *schema.ResourceData) error {

	log.Printf("[DEBUG] Starting Virtual Machine: %s", m.Metadata.Name)

	m.Spec.Resources.PowerState = "POWERED_ON"
	payload, err := json.Marshal(m)
	if err != nil {
		return err
	}
	url := c.URL + "/vms/" + d.Id()
	method := "PUT"
	_, err = requestutils.RequestHandler(url, method, payload, c.Username, c.Password, c.Insecure)
	return err
}

// WaitForProcess waits till the nutanix gets to running
func (c *V3Client) WaitForProcess(vmresp1 *VMResponse) (bool, error) {
	uuid := vmresp1.Metadata.UUID
	url := c.URL + "/vms/" + uuid
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
func (c *V3Client) WaitForIP(vmresp *VMResponse, d *schema.ResourceData) error {
	uuid := vmresp.Metadata.UUID
	url := c.URL + "/vms/" + uuid
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
func (c *V3Client) CreateMachine(m *vmdefn.VirtualMachine) ([]byte, error) {

	payload, err := json.Marshal(m)
	check(err)

	method := "POST"
	resp, err := requestutils.RequestHandler(c.URL+"/vms", method, payload, c.Username, c.Password, c.Insecure)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func resourceNutanixVirtualMachineCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*V3Client)
	machine := vm.SetMachineConfig(d)
	machine.Spec.Name = d.Get("name").(string)
	machine.Metadata.Name = d.Get("name").(string)
	log.Printf("[DEBUG] Creating Virtual Machine: %s", d.Id())

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
	d.Set("ip_address", "")

	if machine.Spec.Resources.NicList != nil && machine.Spec.Resources.PowerState == "POWERED_ON" {
		err = client.WaitForIP(&vmresp, d)
	}
	if err != nil {
		return err
	}

	uuid, err := client.MachineExists(machine.Spec.Name)
	if err != nil {
		return err
	}
	if uuid == "" {
		fmt.Errorf("Machine doesn't exists")
	}

	d.SetId(uuid)
	return nil

}

func resourceNutanixVirtualMachineRead(d *schema.ResourceData, m interface{}) error {

	return nil
}

func resourceNutanixVirtualMachineUpdate(d *schema.ResourceData, meta interface{}) error {
	// Enable partial state mode
	d.Partial(true)
	client := meta.(*V3Client)
	machine := vm.SetMachineConfig(d)
	machine.Metadata.Name = d.Get("name").(string)
	machine.Spec.Name = d.Get("name").(string)

	if d.HasChange("name") || d.HasChange("spec") || d.HasChange("metadata") {

		err := client.UpdateMachine(&machine, d)
		if err != nil {
			return err
		}
		d.SetPartial("spec")
		d.SetPartial("metadata")
	}
	//Disabling partial state mode. This will cause terraform to save all fields again
	d.Partial(false)
	vmresp := VMResponse{Metadata: &vmdefn.MetaDataStruct{UUID: d.Id()}}
	status, err := client.WaitForProcess(&vmresp)
	if status != true {
		return err
	}
	d.Set("ip_address", "")
	if len(machine.Spec.Resources.NicList) > 0 && machine.Spec.Resources.PowerState == "POWERED_ON" {
		err := client.WaitForIP(&vmresp, d)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceNutanixVirtualMachineDelete(d *schema.ResourceData, m interface{}) error {

	client := m.(*V3Client)
	log.Printf("[DEBUG] Deleting Virtual Machine: %s", d.Id())
	machine := vm.SetMachineConfig(d)
	machine.Spec.Name = d.Get("name").(string)
	machine.Metadata.Name = d.Get("name").(string)

	err := client.DeleteMachine(&machine, d)
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
																	Schema: referenceSchema(),
																},
															},
															"availability_zone_reference": &schema.Schema{
																Type:     schema.TypeSet,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: referenceSchema(),
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
								Schema: referenceSchema(),
							},
						},
						"cluster_reference": &schema.Schema{
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: referenceSchema(),
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
											Schema: referenceSchema(),
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
														Schema: referenceSchema(),
													},
												},
												"subnet_reference": &schema.Schema{
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Resource{
														Schema: referenceSchema(),
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
														Schema: referenceSchema(),
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
								Schema: referenceSchema(),
							},
						},
					},
				},
			},
		},
	}
}

// Schema of Reference
func referenceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
	}
}
