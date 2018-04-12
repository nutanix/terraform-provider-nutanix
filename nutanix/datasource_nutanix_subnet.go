package nutanix

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func dataSourceNutanixSubnet() *schema.Resource {
	return &schema.Resource{
		Read: datasourceNutanixSubnetRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"kind": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"subnet_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vlan_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ip_config": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pool_range": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"subnet_ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"prefix_length": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"default_gateway_ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"dhcp_options": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dhcp_server_address_host": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"dhcp_server_address_port": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"domain_name_server_list": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"boot_file_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"domain_search_list": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"domain_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"tftp_server_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func datasourceNutanixSubnetRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Datasource Subnet Read : %s", d.Get("name").(string))

	client := meta.(*V3Client)
	SubnetAPIInstance := SubnetAPIInstance(client)

	uuid, err := client.SubnetExists(d.Get("name").(string))
	if err != nil {
		return err
	}

	if uuid == "" {
		return fmt.Errorf("Nic doesn't exists in given cluster.")
	}

	get_subnet, APIResponse, err := SubnetAPIInstance.SubnetsUuidGet(uuid)
	if err != nil {
		return err
	}

	err = checkAPIResponse(*APIResponse)
	if err != nil {
		return err
	}

	get_subnet_json, _ := json.Marshal(get_subnet)
	log.Printf("[DEBUG] Read Subnet %s", get_subnet_json)

	d.Set("name", get_subnet.Status.Name)
	d.Set("uuid", uuid)

	d.Set("default_gateway_ip", get_subnet.Status.Resources.IpConfig.DefaultGatewayIp)
	d.Set("prefix_length", get_subnet.Status.Resources.IpConfig.PrefixLength)
	d.Set("subnet_ip", get_subnet.Status.Resources.IpConfig.SubnetIp)
	d.SetId(uuid)
	return nil
}
