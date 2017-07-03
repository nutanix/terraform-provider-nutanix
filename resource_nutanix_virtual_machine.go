package nutanix

import (
	nutanixV3 "nutanixV3"
	"errors"
	"fmt"
	"os"
	"bufio"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ideadevice/terraform-ahv-provider-plugin/flg"
	vmconfig "github.com/ideadevice/terraform-ahv-provider-plugin/virtualmachineconfig"
	"log"
	"runtime/debug"
)

var statusCodeFilter map[int]bool
func init() {
	statusMap := map[int]bool{
		200: true,
		201: true,
		202: true,
		203: true,
		204: true,
		205: true,
		206: true,
		207: true,
		208: true,
	}
	statusCodeFilter = statusMap
}

// Function checks if there is an error
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkAPIResponse(resp nutanixV3.APIResponse) error {
	response := fmt.Sprintf("Response ==> %+v\n Response Message ==> %+v\n Request ==> %+v\n Request Body==> %+v", resp.Response, resp.Message, resp.Response.Request, resp.Response.Request.Body)
	if flg.HTTPLog != "" {
		file, err := os.Create(flg.HTTPLog)
		if err != nil {
			return err
		}
		w := bufio.NewWriter(file)
		defer file.Close()
		defer w.Flush()
		fmt.Fprintf(w, "%v", response)
	}
	if !statusCodeFilter[resp.StatusCode] {
		errormsg := errors.New(response)
		return errormsg
	}
	return nil
}

// RecoverFunc can be used to recover from panics. name is the name of the caller
func RecoverFunc(name string) {
	if err := recover(); err != nil {
		log.Printf("Recovered from error %s, %s", err, name)
		log.Printf("Stack Trace: %s", debug.Stack())
		panic(err)
	}
}

// setAPIInstance sets the nutanixV3.VmApi from the V3Client
func setAPIInstance(c *V3Client) *(nutanixV3.VmApi) {
	APIInstance := nutanixV3.NewVmApi()
	APIInstance.Configuration.Username = c.Username
	APIInstance.Configuration.Password = c.Password
	APIInstance.Configuration.BasePath = c.URL
	APIInstance.Configuration.APIClient.Insecure = c.Insecure
	return APIInstance
}

// WaitForProcess waits till the nutanix gets to running
func (c *V3Client) WaitForProcess(uuid string) (bool, error) {
	APIInstance := setAPIInstance(c)
	for {
		VMIntentResponse, APIresponse, err := APIInstance.VmsUuidGet(uuid)
		if err != nil {
			return false, err
		}
		err = checkAPIResponse(*APIresponse)
		if err != nil {
			return false, err
		}
		if VMIntentResponse.Status.State == "COMPLETE" {
			return true, nil
		}
	}
	return false, nil
}

// WaitForIP function sets the ip address obtained by the GET request
func (c *V3Client) WaitForIP(uuid string, d *schema.ResourceData) error {
	APIInstance := setAPIInstance(c)
	for {
		VMIntentResponse, APIresponse, err := APIInstance.VmsUuidGet(uuid)
		if err != nil {
			return err
		}
		err = checkAPIResponse(*APIresponse)
		if err != nil {
			return  err
		}
		if len(VMIntentResponse.Status.Resources.NicList) != 0 {
			for i := range VMIntentResponse.Status.Resources.NicList {
				if len(VMIntentResponse.Status.Resources.NicList[i].IpEndpointList) != 0 {
					if ip := VMIntentResponse.Status.Resources.NicList[i].IpEndpointList[0].Address; ip != "" {
						d.Set("ip_address", ip)
						return nil
					}
				}
			}
		}
	}
	return nil
}

func resourceNutanixVirtualMachineCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*V3Client)
	machine := vmconfig.SetMachineConfig(d)
	machine.Spec.Name = d.Get("name").(string)
	machine.Metadata.Name = d.Get("name").(string)
	log.Printf("[DEBUG] Creating Virtual Machine: %s", d.Id())
	APIInstance := setAPIInstance(client)
	VMIntentResponse, APIResponse, err := APIInstance.VmsPost(machine)
	if err != nil {
		return err
	}
	err = checkAPIResponse(*APIResponse)
	if err != nil {
		return  err
	}

	uuid := VMIntentResponse.Metadata.Uuid
	status, err := client.WaitForProcess(uuid)
	for status != true {
			return err
	}
	d.Set("ip_address", "")

	if machine.Spec.Resources.NicList != nil && machine.Spec.Resources.PowerState == "POWERED_ON" {
		err = client.WaitForIP(uuid, d)
	}
	if err != nil {
		return err
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
	machine := vmconfig.SetMachineConfig(d)
	machine.Metadata.Name = d.Get("name").(string)
	machine.Spec.Name = d.Get("name").(string)
	APIInstance := setAPIInstance(client)
	uuid := d.Id()
	log.Printf("[DEBUG] Updating Virtual Machine: %s, %s", machine.Spec.Name, d.Id())

	if d.HasChange("name") || d.HasChange("spec") || d.HasChange("metadata") {

		_, APIResponse, err := APIInstance.VmsUuidPut(uuid, machine)
		if err != nil {
			return err
		}
		err = checkAPIResponse(*APIResponse)
		if err != nil {
			return  err
		}
		d.SetPartial("spec")
		d.SetPartial("metadata")
	}
	//Disabling partial state mode. This will cause terraform to save all fields again
	d.Partial(false)
	status, err := client.WaitForProcess(uuid)
	if status != true {
		return err
	}
	d.Set("ip_address", "")
	if len(machine.Spec.Resources.NicList) > 0 && machine.Spec.Resources.PowerState == "POWERED_ON" {
		err := client.WaitForIP(uuid, d)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceNutanixVirtualMachineDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*V3Client)
	log.Printf("[DEBUG] Deleting Virtual Machine: %s", d.Id())
	machine := vmconfig.SetMachineConfig(d)
	machine.Spec.Name = d.Get("name").(string)
	machine.Metadata.Name = d.Get("name").(string)
	APIInstance := setAPIInstance(client)
	uuid := d.Id()

	APIResponse, err := APIInstance.VmsUuidDelete(uuid)
	if err != nil {
		return err
	}
	err = checkAPIResponse(*APIResponse)
	if err != nil {
		return  err
	}

	d.SetId("")
	return nil
}

// MachineExists function returns the uuid of the machine with given name
func resourceNutanixVirtualMachineExists(d *schema.ResourceData, m interface{}) (bool, error) {
	log.Printf("[DEBUG] Checking Virtual Machine Existance: %s", d.Id())
	client := m.(*V3Client)
	APIInstance := setAPIInstance(client)

	getEntitiesRequest := nutanixV3.VmListMetadata{} // VmListMetadata
	VMListIntentResponse, APIResponse, err := APIInstance.VmsListPost(getEntitiesRequest)
	if err != nil {
		return false, err
	}
	err = checkAPIResponse(*APIResponse)
	if err != nil {
		return  false,err
	}

	for i := range VMListIntentResponse.Entities {
		if VMListIntentResponse.Entities[i].Metadata.Uuid == d.Id() {
			return true, nil
		}
	}
	return false, nil
}


func resourceNutanixVirtualMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceNutanixVirtualMachineCreate,
		Read:   resourceNutanixVirtualMachineRead,
		Update: resourceNutanixVirtualMachineUpdate,
		Delete: resourceNutanixVirtualMachineDelete,
		Exists: resourceNutanixVirtualMachineExists,

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
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
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
			"metadata": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"last_update_time": &schema.Schema{
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
						"creation_time": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
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
