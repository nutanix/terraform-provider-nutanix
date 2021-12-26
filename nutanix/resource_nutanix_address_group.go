package nutanix

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
	"log"
	"strings"
)

func resourceNutanixAddressGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceNutanixAddressGroupCreate,
		Read:   resourceNutanixAddressGroupRead,
		Delete: resourceNutanixAddressGroupDelete,
		Update: resourceNutanixAddressGroupUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_address_block_list": {
				Type:     schema.TypeList,
				Required: true,
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

func resourceNutanixAddressGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API
	id := d.Id()
	response, err := conn.V3.GetAddressGroup(id)

	request := &v3.AddressGroupInput{}

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
		}
		return fmt.Errorf("error retrieving for address group id (%s) :%+v", id, err)
	}

	group := response.AddressGroup

	if d.HasChange("name") {
		group.Name = utils.StringPtr(d.Get("name").(string))
	}

	if d.HasChange("description") {
		group.Description = utils.StringPtr(d.Get("description").(string))
	}

	if d.HasChange("ip_address_block_list") {
		blockList, err := expandAddressEntry(d)

		if err != nil {
			return err
		}

		group.BlockList = blockList
	}

	request.Name = group.Name
	request.Description = group.Description
	request.BlockList = group.BlockList

	errUpdate := conn.V3.UpdateAddressGroup(d.Id(), request)
	if errUpdate != nil {
		return fmt.Errorf("error updating address group id %s): %s", d.Id(), errUpdate)
	}

	return resourceNutanixAddressGroupRead(d, meta)

}

func resourceNutanixAddressGroupDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API

	log.Printf("[Debug] Destroying the address group with the ID %s", d.Id())

	if err := conn.V3.DeleteAddressGroup(d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceNutanixAddressGroupRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading AddressGroup: %s", d.Get("name").(string))

	// Get client connection
	conn := meta.(*Client).API

	// Make request to the API
	resp, err := conn.V3.GetAddressGroup(d.Id())

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return nil
		}
		return err
	}

	if err := d.Set("ip_address_block_list", flattenAddressEntry(resp.AddressGroup.BlockList)); err != nil {
		return err
	}

	d.Set("name", utils.StringValue(resp.AddressGroup.Name))

	return d.Set("description", utils.StringValue(resp.AddressGroup.Description))
}

func flattenAddressEntry(group []*v3.IPAddressBlock) []map[string]interface{} {
	groupList := make([]map[string]interface{}, 0)
	for _, v := range group {
		groupItem := make(map[string]interface{})
		groupItem["ip"] = v.IPAddress
		groupItem["prefix_length"] = v.PrefixLength
		groupList = append(groupList, groupItem)
	}

	return groupList
}

func resourceNutanixAddressGroupCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API

	request := &v3.AddressGroupInput{}
	request.BlockList = make([]*v3.IPAddressBlock, 0)

	// Read Arguments and set request values

	if name, ok := d.GetOk("name"); ok {
		request.Name = utils.StringPtr(name.(string))
	}

	if desc, ok := d.GetOk("description"); ok {
		request.Description = utils.StringPtr(desc.(string))
	}
	addressList, err := expandAddressEntry(d)

	if err != nil {
		return err
	}

	request.BlockList = addressList

	resp, err := conn.V3.CreateAddressGroup(request)

	if err != nil {
		return err
	}

	n := *resp.UUID

	// set terraform state
	d.SetId(n)

	return resourceNutanixAddressGroupRead(d, meta)
}

func expandAddressEntry(d *schema.ResourceData) ([]*v3.IPAddressBlock, error) {
	if groups, ok := d.GetOk("ip_address_block_list"); ok {
		set := groups.([]interface{})
		outbound := make([]*v3.IPAddressBlock, len(set))

		for k, v := range set {
			entry := v.(map[string]interface{})

			block := &v3.IPAddressBlock{}
			if ip, ipok := entry["ip"]; ipok {
				block.IPAddress = utils.StringPtr(ip.(string))
			} else {
				return nil, fmt.Errorf("error updating address group id %s): ip missing", d.Id())
			}

			if length, lengthok := entry["prefix_length"]; lengthok {
				block.PrefixLength = utils.Int64Ptr(int64(length.(int)))
			} else {
				return nil, fmt.Errorf("error updating address group id %s): prefix_length missing", d.Id())
			}

			outbound[k] = block
		}
		return outbound, nil
	}

	return nil, nil
}
