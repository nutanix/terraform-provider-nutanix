package nutanix

import (
	"errors"
	"fmt"
	"os"
	"bufio"
	"time"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ideadevice/terraform-ahv-provider-plugin/flg"
	vmschema "github.com/ideadevice/terraform-ahv-provider-plugin/virtualmachineschema"
	vmconfig "github.com/ideadevice/terraform-ahv-provider-plugin/virtualmachineconfig"
	nutanixV3 "nutanixV3"
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
		time.Sleep(3000*time.Millisecond)
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
		time.Sleep(3000*time.Millisecond)
	}
	return nil
}

func resourceNutanixVirtualMachineCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*V3Client)
	machine := vmconfig.SetMachineConfig(d)
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

func resourceNutanixVirtualMachineRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*V3Client)
	APIInstance := setAPIInstance(client)
	//VMIntentResponse, APIResponse, err := APIInstance.VmsUuidGet(d.Id())
	_, APIResponse, err := APIInstance.VmsUuidGet(d.Id())
	if err != nil {
		return err
	}

	err = checkAPIResponse(*APIResponse)
	if err != nil {
		return err
	}
	return nil
}

func resourceNutanixVirtualMachineUpdate(d *schema.ResourceData, meta interface{}) error {
	// Enable partial state mode
	d.Partial(true)
	client := meta.(*V3Client)
	machine := vmconfig.SetMachineConfig(d)
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

		Schema: vmschema.VMSchema(),
	}
}
