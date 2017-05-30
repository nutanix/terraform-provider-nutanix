package ahvproviderplugin

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	st "github.com/ideadevice/terraform-ahv-provider-plugin/jsonstruct"
	"github.com/ideadevice/terraform-ahv-provider-plugin/requestutils"
	"io/ioutil"
	"log"
	"reflect"
	"runtime/debug"
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

//RecoverFunc can be used to recover from panics. name is the name of the caller
func RecoverFunc(name string) {
	if err := recover(); err != nil {
		log.Printf("Recovered from error %s, %s", err, name)
		log.Printf("Stack Trace: %s", debug.Stack())
		panic(err)
	}
}

func setStructStringField(v interface{}, field string, path string, d *schema.ResourceData) {
	defer RecoverFunc("setStructStringField")
	a := reflect.ValueOf(v).Elem().FieldByName(field)
	valueCurr := fmt.Sprintf("%v", a)
	temp := d.Get(path)
	if a.IsValid() && temp != nil && temp.(string) != "" {
		a.SetString(temp.(string))
	} else if valueCurr == "string" {
		a.SetString("")
	}
}

func setStructIntField(v interface{}, field string, path string, d *schema.ResourceData) {
	a := reflect.ValueOf(v).Elem().FieldByName(field)
	temp := d.Get(path)
	if a.IsValid() && temp != nil && temp.(int) != 0 {
		var intVal int64
		intVal = int64(temp.(int))
		a.SetInt(intVal)
	}
}

// SetJSONFields is function for setting the JSON fields from config file
func SetJSONFields(myStruct interface{}, path string, d *schema.ResourceData) {
	v := reflect.ValueOf(myStruct).Elem()
	typeOfT := v.Type()

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		typeName := fmt.Sprintf("%s", f.Type())
		fieldName := fmt.Sprintf("%s", typeOfT.Field(i).Name)

		var pathNew string
		if path == "" {
			pathNew = fieldName
		} else {
			pathNew = path + "_" + fieldName
		}

		if typeName == "string" {
			setStructStringField(myStruct, fieldName, pathNew, d)
		} else if typeName == "int" {
			setStructIntField(myStruct, fieldName, pathNew, d)
		} else if typeName != strings.TrimSuffix(typeName, "Struct") && !(reflect.ValueOf(f.Interface()).IsNil()) {
			SetJSONFields(f.Interface(), pathNew, d)
		}
	}
}

// ID returns the id to be set
func (m *Machine) ID() string {
	return "ID-" + m.Name + "!!"
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
		if vm.Spec.Name == m.Name {
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

	var JSON st.JSONstruct

	Input, err := ioutil.ReadFile("json_template")
	check(err)
	InputPattern := []byte(Input)

	json.Unmarshal(InputPattern, &JSON)

	SetJSONFields(&JSON, "", d)
	JSON.Spec.Name = m.Name
	JSON.Metadata.Name = m.Name

	jsonStr, err1 := json.Marshal(JSON)
	check(err1)

	method := "POST"
	requestutils.RequestHandler(c.Endpoint, method, jsonStr, c.Username, c.Password)
	return nil
}

func resourceServerCreate(d *schema.ResourceData, m interface{}) error {

	client := m.(*MyClient)
	machine := Machine{
		Name: d.Get("name").(string),
		SpecResourcesNumVCPUsPerSocket: d.Get("Spec_Resources_NumVCPUsPerSocket").(int),
		SpecResourcesNumSockets:        d.Get("Spec_Resources_NumSockets").(int),
		SpecResourcesMemorySizeMib:     d.Get("Spec_Resources_MemorySizeMib").(int),
		SpecResourcesPowerState:        d.Get("Spec_Resources_PowerState").(string),
		APIversion:                     d.Get("APIversion").(string),
	}

	err := client.CreateMachine(&machine, d)
	if err != nil {
		return err
	}

	d.SetId(machine.ID())
	return nil

}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
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

	client := m.(*MyClient)
	machine := Machine{
		Name: d.Get("name").(string),
		SpecResourcesNumVCPUsPerSocket: d.Get("Spec_Resources_NumVCPUsPerSocket").(int),
		SpecResourcesNumSockets:        d.Get("Spec_Resources_NumSockets").(int),
		SpecResourcesMemorySizeMib:     d.Get("Spec_Resources_MemorySizeMib").(int),
		SpecResourcesPowerState:        d.Get("Spec_Resources_PowerState").(string),
		APIversion:                     d.Get("APIversion").(string),
	}

	err := client.DeleteMachine(&machine)
	if err != nil {
		return err
	}

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
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"Spec_Resources_NumVCPUsPerSocket": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"Spec_Resources_NumSockets": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"Spec_Resources_MemorySizeMib": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"Spec_Resources_PowerState": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"APIversion": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}
