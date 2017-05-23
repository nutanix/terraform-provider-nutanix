package ahvproviderplugin

import (
	//"encoding/json"
	"bufio"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	//"io/ioutil"
	"github.com/ideadevice/terraform-ahv-provider-plugin/requestutils"
	"os"
	"strings"
)

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
	defer f.Close()

	//	_, err = f.WriteString("executing create function\n")
	//	check(err)
	w := bufio.NewWriter(f)

	// Opening the json file
	filejson, err1 := os.Open("json_template")
	check(err1)
	defer filejson.Close()

	// Create a new scanner and read the file line by line
	scanner := bufio.NewScanner(filejson)
	for scanner.Scan() {
		strName := "\"name\": \"string\""
		if strings.Contains(scanner.Text(), strName) {
			name := d.Get("name").(string)
			strnew := "\"name\": \"" + name + "\""
			str := strings.Replace(scanner.Text(), strName, strnew, 1)
			_, err = fmt.Fprintf(w, "%v\n", str)
			check(err)
		} else {
			_, err = fmt.Fprintf(w, "%v\n", scanner.Text())
			check(err)
		}
	}
	w.Flush()
	// Parsing JSON-encoded data into data
	//var data interface{}
	//err = json.Unmarshal(plan, &data)
	//check(err)

	//	fmt.Println(data)
	url := "https://private-anon-466a7b0395-restapi3.apiary-mock.com/notes"
	var jsonStr = []byte(`{"title":" Yo POST requesting is working :D ."}`)
	method := "POST"
	requestutils.RequestHandler(url, method, jsonStr)

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
	url := "http://www.example.com/customers/12345"
	var jsonStr = []byte(`yo`)
	method := "GET"
	requestutils.RequestHandler(url, method, jsonStr)
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

	url := "http://www.example.com/customers/12345"
	var jsonStr = []byte(`yo`)
	method := "DELETE"
	requestutils.RequestHandler(url, method, jsonStr)

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
