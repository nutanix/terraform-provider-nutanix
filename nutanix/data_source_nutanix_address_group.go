package nutanix

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func dataSourceNutanixAddressGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNutanixAddressGroupRead,
		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_address_block_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"prefix_length": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"address_group_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNutanixAddressGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API

	if uuid, uuidOk := d.GetOk("uuid"); uuidOk {
		group, reqErr := conn.V3.GetAddressGroup(uuid.(string))

		if reqErr != nil {
			if strings.Contains(fmt.Sprint(reqErr), "ENTITY_NOT_FOUND") {
				d.SetId("")
			}
			return fmt.Errorf("error reading user with error %s", reqErr)
		}

		if err := d.Set("name", utils.StringValue(group.AddressGroup.Name)); err != nil {
			return err
		}

		if err := d.Set("description", utils.StringValue(group.AddressGroup.Description)); err != nil {
			return err
		}

		if err := d.Set("address_group_string", utils.StringValue(group.AddressGroup.AddressGroupString)); err != nil {
			return err
		}

		if err := d.Set("ip_address_block_list", flattenAddressEntry(group.AddressGroup.BlockList)); err != nil {
			return err
		}

		d.SetId(uuid.(string))
	} else {
		return fmt.Errorf("please provide `uuid`")
	}
	return nil
}
