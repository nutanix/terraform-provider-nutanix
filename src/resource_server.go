package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func updateAddress(d *schema.ResourceData) error{
	return nil
}

func resourceServerCreate(d *schema.ResourceData, m interface{} ) error{
	address := d.Get("address").(string)
	d.SetId("myID "+address)
	return nil
}

func resourceServerRead(d *schema.ResourceData, m interface{} ) error{
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
	return nil

}

func resourceServerUpdate(d *schema.ResourceData, m interface{} ) error{
	// Enable partial state mode
	d.Partial(true)
    // checking that address has changed or not
	if d.HasChange("address"){
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

func resourceServerDelete(d *schema.ResourceData, m interface{} ) error{
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
				Type: schema.TypeString,
				Required: true,
			},
		},
	}
}
