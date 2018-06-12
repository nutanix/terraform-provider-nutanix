package nutanix

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func resourceNutanixCategoryKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceNutanixCategoryKeyCreateOrUpdate,
		Read:   resourceNutanixCategoryKeyRead,
		Update: resourceNutanixCategoryKeyCreateOrUpdate,
		Delete: resourceNutanixCategoryKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: getCategoryKeySchema(),
	}
}

func resourceNutanixCategoryKeyCreateOrUpdate(resourceData *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Creating CategoryKey: %s", resourceData.Get("name").(string))

	conn := meta.(*Client).API

	request := &v3.CategoryKey{}

	name, nameOK := resourceData.GetOk("name")

	// Read Arguments and set request values
	if v, ok := resourceData.GetOk("api_version"); ok {
		request.APIVersion = utils.String(v.(string))
	}

	if desc, ok := resourceData.GetOk("description"); ok {
		request.Description = utils.String(desc.(string))
	}

	// validaste required fields
	if !nameOK {
		return fmt.Errorf("please provide the required attribute name")
	}

	request.Name = utils.String(name.(string))

	//Make request to the API
	resp, err := conn.V3.CreateOrUpdateCategoryKey(request)

	if err != nil {
		return err
	}

	n := *resp.Name

	// set terraform state
	resourceData.SetId(n)

	return resourceNutanixCategoryKeyRead(resourceData, meta)
}

func resourceNutanixCategoryKeyRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading CategoryKey: %s", d.Get("name").(string))

	// Get client connection
	conn := meta.(*Client).API

	// Make request to the API
	resp, err := conn.V3.GetCategoryKey(d.Id())

	if err != nil {
		return err
	}

	d.Set("api_version", utils.StringValue(resp.APIVersion))
	d.Set("name", utils.StringValue(resp.Name))
	d.Set("description", utils.StringValue(resp.Description))

	return d.Set("system_defined", utils.BoolValue(resp.SystemDefined))
}

func resourceNutanixCategoryKeyDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API

	log.Printf("Destroying the category with the name %s", d.Id())
	fmt.Printf("Destroying the category with the name %s", d.Id())

	if err := conn.V3.DeleteCategoryKey(d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func getCategoryKeySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
