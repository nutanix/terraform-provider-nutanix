package nutanix

import (
	"fmt"
	"log"

	karbon "github.com/terraform-providers/terraform-provider-nutanix/client/karbon"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceNutanixKarbonPrivateRegistry() *schema.Resource {
	return &schema.Resource{
		Create: resourceNutanixKarbonPrivateRegistryCreate,
		Read:   resourceNutanixKarbonPrivateRegistryRead,
		Update: resourceNutanixKarbonPrivateRegistryUpdate,
		Delete: resourceNutanixKarbonPrivateRegistryDelete,
		Exists: resourceNutanixKarbonPrivateRegistryExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 1,
		Schema:        KarbonPrivateRegistryResourceMap(),
	}
}

func KarbonPrivateRegistryResourceMap() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"cert": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"url": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"port": {
			Type:     schema.TypeInt,
			Required: true,
			// ForceNew: true,
			ValidateFunc: validation.IntAtLeast(1),
		},
		"endpoint": {
			Type:     schema.TypeString,
			Computed: true,
			// ForceNew: true,
		},
	}
}

func resourceNutanixKarbonPrivateRegistryCreate(d *schema.ResourceData, meta interface{}) error {
	log.Print("[Debug] Entering resourceNutanixKarbonPrivateRegistryCreate")
	// Get client connection
	client := meta.(*Client)
	conn := client.KarbonAPI
	setTimeout(meta)
	// Prepare request
	karbonPrivateRegistry := &karbon.PrivateRegistryIntentInput{}
	if name, ok := d.GetOk("name"); ok {
		n := name.(string)
		karbonPrivateRegistry.Name = &n
	} else {
		return fmt.Errorf("error occurred during private registry creation:\n Name must be set")
	}
	if url, ok := d.GetOk("url"); ok {
		u := url.(string)
		karbonPrivateRegistry.URL = &u
	} else {
		return fmt.Errorf("Error occurred during private registry creation:\n URL must be set")
	}
	if port, ok := d.GetOk("port"); ok {
		p := int64(port.(int))
		karbonPrivateRegistry.Port = &p
	}

	if cert, ok := d.GetOk("cert"); ok {
		c := cert.(string)
		karbonPrivateRegistry.Cert = &c
	}
	utils.PrintToJSON(karbonPrivateRegistry, "[DEBUG karbonPrivateRegistry: ")
	createPrivateRegistryResponse, err := conn.PrivateRegistry.CreateKarbonPrivateRegistry(karbonPrivateRegistry)
	if err != nil {
		return fmt.Errorf("Error occurred during private registry creation:\n %s", err)
	}
	utils.PrintToJSON(createPrivateRegistryResponse, "[DEBUG createPrivateRegistryResponse: ")

	// Set terraform state id
	d.SetId(*createPrivateRegistryResponse.UUID)
	return resourceNutanixKarbonPrivateRegistryRead(d, meta)
}

func resourceNutanixKarbonPrivateRegistryRead(d *schema.ResourceData, meta interface{}) error {
	log.Print("[Debug] Entering resourceNutanixKarbonPrivateRegistryRead")
	// Get client connection
	conn := meta.(*Client).KarbonAPI
	setTimeout(meta)
	// Make request to the API
	var name interface{}
	var ok bool
	if name, ok = d.GetOk("name"); !ok {
		return fmt.Errorf("Cannot read Private Registry without name!")
	}
	resp, err := conn.PrivateRegistry.GetKarbonPrivateRegistry(name.(string))
	if err != nil {
		d.SetId("")
		return nil
	}
	if err := d.Set("name", *resp.Name); err != nil {
		return fmt.Errorf("error setting name for Karbon private registry %s: %s", d.Id(), err)
	}
	// log.Print(*resp.Endpoint)
	if err := d.Set("endpoint", *resp.Endpoint); err != nil {
		return fmt.Errorf("error setting endpoint for Karbon private registry %s: %s", d.Id(), err)
	}
	d.SetId(*resp.UUID)
	return nil
}

func resourceNutanixKarbonPrivateRegistryUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Print("[Debug] Entering resourceNutanixKarbonPrivateRegistryUpdate")
	return resourceNutanixKarbonPrivateRegistryRead(d, meta)
}

func resourceNutanixKarbonPrivateRegistryDelete(d *schema.ResourceData, meta interface{}) error {
	log.Print("[Debug] Entering resourceNutanixKarbonPrivateRegistryDelete")
	client := meta.(*Client)
	conn := client.KarbonAPI
	setTimeout(meta)
	karbonPrivateRegistryName := d.Get("name").(string)
	log.Printf("[DEBUG] Deleting Karbon cluster: %s, %s", karbonPrivateRegistryName, d.Id())

	_, err := conn.PrivateRegistry.DeleteKarbonPrivateRegistry(karbonPrivateRegistryName)
	if err != nil {
		return fmt.Errorf("error while deleting Karbon Private Registry UUID(%s): %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}

func resourceNutanixKarbonPrivateRegistryExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Print("[DEBUG] Entering resourceNutanixKarbonPrivateRegistryExists")
	conn := meta.(*Client).KarbonAPI
	setTimeout(meta)
	// Make request to the API
	var name interface{}
	var ok bool
	if name, ok = d.GetOk("name"); !ok {
		return false, fmt.Errorf("Cannot read Private Registry without name!")
	}
	resp, err := conn.PrivateRegistry.GetKarbonPrivateRegistry(name.(string))
	log.Print("error:")
	log.Print(err)
	utils.PrintToJSON(resp, "resourceNutanixKarbonPrivateRegistryExists resp: ")
	if err != nil {
		d.SetId("")
		return false, nil
	}
	return true, nil
}
