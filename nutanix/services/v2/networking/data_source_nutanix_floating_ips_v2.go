package networking

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixFloatingIPsV4() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNutanixFloatingIPsV4Read,
		Schema: map[string]*schema.Schema{
			"page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"expand": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"floating_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixFloatingIPV4(),
			},
		},
	}
}

func datasourceNutanixFloatingIPsV4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	// initialize query params
	var filter, orderBy, expand *string
	var page, limit *int

	if pagef, ok := d.GetOk("page"); ok {
		page = utils.IntPtr(pagef.(int))
	} else {
		page = nil
	}
	if limitf, ok := d.GetOk("limit"); ok {
		limit = utils.IntPtr(limitf.(int))
	} else {
		limit = nil
	}
	if filterf, ok := d.GetOk("filter"); ok {
		filter = utils.StringPtr(filterf.(string))
	} else {
		filter = nil
	}
	if order, ok := d.GetOk("order_by"); ok {
		orderBy = utils.StringPtr(order.(string))
	} else {
		orderBy = nil
	}
	if expandf, ok := d.GetOk("expand"); ok {
		expand = utils.StringPtr(expandf.(string))
	} else {
		expand = nil
	}

	resp, err := conn.FloatingIPAPIInstance.ListFloatingIps(page, limit, filter, orderBy, expand)
	if err != nil {
		var errordata map[string]interface{}
		e := json.Unmarshal([]byte(err.Error()), &errordata)
		if e != nil {
			return diag.FromErr(e)
		}
		data := errordata["data"].(map[string]interface{})
		errorList := data["error"].([]interface{})
		errorMessage := errorList[0].(map[string]interface{})
		return diag.Errorf("error while fetching floating_ips : %v", errorMessage["message"])
	}

	getResp := resp.Data.GetValue().([]import1.FloatingIp)

	if err := d.Set("floating_ips", flattenFloatingIPsEntities(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenFloatingIPsEntities(pr []import1.FloatingIp) []map[string]interface{} {
	if len(pr) > 0 {
		fips := make([]map[string]interface{}, len(pr))

		for k, v := range pr {
			fip := make(map[string]interface{})

			fip["ext_id"] = v.ExtId
			fip["name"] = v.Name
			fip["description"] = v.Description
			fip["association"] = flattenAssociation(v.Association)
			fip["floating_ip"] = flattenFloatingIP(v.FloatingIp)
			fip["external_subnet_reference"] = v.ExternalSubnetReference
			fip["external_subnet"] = flattenExternalSubnet(v.ExternalSubnet)
			fip["private_ip"] = v.PrivateIp
			fip["floating_ip_value"] = v.FloatingIpValue
			fip["association_status"] = v.AssociationStatus
			fip["vpc_reference"] = v.VpcReference
			fip["vm_nic_reference"] = v.VmNicReference
			fip["vpc"] = flattenVpc(v.Vpc)
			fip["vm_nic"] = flattenVMNic(v.VmNic)
			fip["links"] = flattenLinks(v.Links)
			fip["tenant_id"] = v.TenantId
			fip["metadata"] = flattenMetadata(v.Metadata)

			fips[k] = fip
		}
		return fips
	}
	return nil
}
