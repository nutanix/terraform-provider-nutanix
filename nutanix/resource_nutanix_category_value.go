package nutanix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func resourceNutanixCategoryValue() *schema.Resource {
	return &schema.Resource{
		Create: resourceNutanixCategoryValueCreateOrUpdate,
		Read:   resourceNutanixCategoryValueRead,
		Update: resourceNutanixCategoryValueCreateOrUpdate,
		Delete: resourceNutanixCategoryValueDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: getCategoryValueSchema(),
	}
}

func resourceNutanixCategoryValueCreateOrUpdate(resourceData *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Creating CategoryValue: %s", resourceData.Get("value").(string))

	conn := meta.(*Client).API

	request := &v3.CategoryValue{}

	name, nameOK := resourceData.GetOk("name")

	value, valueOK := resourceData.GetOk("value")

	// Read Arguments and set request values
	if v, ok := resourceData.GetOk("api_version"); ok {
		request.APIVersion = utils.String(v.(string))
	}

	if desc, ok := resourceData.GetOk("description"); ok {
		request.Description = utils.String(desc.(string))
	}

	// validaste required fields
	if !nameOK || !valueOK {
		return fmt.Errorf("Please provide the required attributes name and value")
	}

	request.Value = utils.String(value.(string))

	//Make request to the API
	resp, err := conn.V3.CreateOrUpdateCategoryValue(name.(string), request)

	if err != nil {
		return err
	}

	v := *resp.Value

	// set terraform state
	resourceData.SetId(v)

	return resourceNutanixCategoryValueRead(resourceData, meta)
}

func resourceNutanixCategoryValueRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading CategoryValue: %s", d.Get("value").(string))

	name, nameOK := d.GetOk("name")

	if !nameOK {
		return fmt.Errorf("Please provide the required attributes name")
	}

	// Get client connection
	conn := meta.(*Client).API

	// Make request to the API
	resp, err := conn.V3.GetCategoryValue(name.(string), d.Id())

	if err != nil {
		return err
	}

	d.Set("api_version", utils.StringValue(resp.APIVersion))
	d.Set("name", utils.StringValue(resp.Name))
	d.Set("description", utils.StringValue(resp.Description))

	return d.Set("system_defined", utils.BoolValue(resp.SystemDefined))
}

func resourceNutanixCategoryValueDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API

	name, nameOK := d.GetOk("name")

	if !nameOK {
		return fmt.Errorf("Please provide the required attributes name")
	}

	log.Printf("Destroying the category with the name %s", d.Id())
	fmt.Printf("Destroying the category with the name %s", d.Id())

	if err := conn.V3.DeleteCategoryValue(name.(string), d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func getCategoryValueSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"value": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"system_defined": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"api_version": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}
