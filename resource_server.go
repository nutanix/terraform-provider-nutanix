package ahvproviderplugin

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ideadevice/terraform-ahv-provider-plugin/requestutils"
	"io/ioutil"
	"os"
	"strings"
)

type specStruct struct {
	Name      string      `json:"name"`
	Resources interface{} `json:"resources"`
}

type metaStruct struct {
	OwnerReference interface{} `json:"owner_reference"`
	SpecVersion    int64       `json:"spec_version"`
	UUID           string      `json:"uuid"`
	Kind           string      `json:"kind"`
	Categories     interface{} `json:"categories"`
}

type vmStruct struct {
	Metadata metaStruct  `json:"metadata"`
	Status   interface{} `json:"status"`
	Spec     specStruct  `json:"spec"`
}

type vmList struct {
	APIVersion string      `json:"api_version"`
	MetaData   interface{} `json:"metadata"`
	Entities   []vmStruct  `json:"entities"`
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

func resourceServerCreate(d *schema.ResourceData, m interface{}) error {

	f, err := os.Create("file1")
	check(err)

	w := bufio.NewWriter(f)

	// Opening the json template file
	filejson, err1 := os.Open("json_template")
	check(err1)
	defer filejson.Close()

	// Create a new scanner and read the file line by line
	scanner := bufio.NewScanner(filejson)
	strName := "\"name\": \"string\""
	flagName := true
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), strName) && flagName {
			name := d.Get("name").(string)
			strnew := "\"name\": \"" + name + "\""
			str := strings.Replace(scanner.Text(), strName, strnew, 1)
			_, err = fmt.Fprintf(w, "%v\n", str)
			check(err)
			flagName = true
		} else {
			_, err = fmt.Fprintf(w, "%v\n", scanner.Text())
			check(err)
		}
	}
	w.Flush()
	f.Close()

	json, err2 := ioutil.ReadFile("file1")
	check(err2)
	jsonStr := []byte(json)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	url := "https://10.5.68.6:9440/api/nutanix/v3/vms"
	method := "POST"
	requestutils.RequestHandler(url, method, jsonStr, username, password)

	address := d.Get("address").(string)
	d.SetId("MyID " + address)
	return nil
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	/*
		client := meta.(*MyClient)

		// Attempt to read from an upstream API
		obj, ok := client.Get(d.Id())

		// If resource does not exist, inform Terraform.
		// We want to return immediately return here to prevent further processing
		if !ok {
			d.SetId("")
			return nil
		}

		d.Set("address", obj.Address)
	*/
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	url := "http://www.example.com/customers/12345"
	var jsonStr = []byte(`yo`)
	method := "GET"
	requestutils.RequestHandler(url, method, jsonStr, username, password)
	return nil

}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
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

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {

	username := d.Get("username").(string)
	password := d.Get("password").(string)
	name := d.Get("name").(string)

	jsonStr := []byte(`{}`)

	url := "https://10.5.68.6:9440/api/nutanix/v3/vms/list"
	method := "POST"
	jsonResponse := requestutils.RequestHandler(url, method, jsonStr, username, password)

	var uuid string
	var vmlist vmList
	err := json.Unmarshal(jsonResponse, &vmlist)
	check(err)
	for _, vm := range vmlist.Entities {
		if vm.Spec.Name == name {
			uuid = vm.Metadata.UUID
		}
	}

	url = "https://10.5.68.6:9440/api/nutanix/v3/vms/" + uuid
	jsonStr = []byte(`yo`)
	method = "DELETE"
	requestutils.RequestHandler(url, method, jsonStr, username, password)

	d.SetId("")
	return nil
}

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}
