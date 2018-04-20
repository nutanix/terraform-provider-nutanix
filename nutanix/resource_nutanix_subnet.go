package nutanix

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

const (
	SUBNET_KIND = "subnet"
)

func resourceNutanixSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceNutanixSubnetCreate,
		Read:   resourceNutanixSubnetRead,
		Update: resourceNutanixSubnetUpdate,
		Delete: resourceNutanixSubnetDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
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
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ip_config": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pool_range": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"subnet_ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"prefix_length": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"default_gateway_ip": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"dhcp_options": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dhcp_server_address_host": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dhcp_server_address_port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"domain_name_server_list": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"boot_file_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"domain_search_list": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"domain_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"tftp_server_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceNutanixSubnetCreate(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] Creating Subnet: %s", d.Get("name").(string))

	client := meta.(*V3Client)
	SubnetAPIInstance := SubnetAPIInstance(client)
	vlan_id := int64(d.Get("vlan_id").(int))

	subnet := nutanixV3.SubnetIntentInput{
		ApiVersion: API_VERSION,
		Metadata: nutanixV3.SubnetMetadata{
			Name: d.Get("name").(string),
			Kind: SUBNET_KIND,
		},
		Spec: nutanixV3.Subnet{
			Description: d.Get("description").(string),
			Name:        d.Get("name").(string),
			Resources: nutanixV3.SubnetResources{
				VlanId: &vlan_id,
			},
		},
	}

	if _, ok := d.GetOk("subnet_type"); ok {
		subnet.Spec.Resources.SubnetType = strings.ToUpper(d.Get("subnet_type").(string))
	}

	if ipConfigData, ok := d.GetOk("ip_config"); ok {
		var ipconfig nutanixV3.IpConfig
		params := ipConfigData.(*schema.Set).List()
		for k := range params {
			data := params[k].(map[string]interface{})
			var pool_range []nutanixV3.IpPool
			for _, pool := range data["pool_range"].(*schema.Set).List() {
				ip_pool := nutanixV3.IpPool{Range_: pool.(string)}
				pool_range = append(pool_range, ip_pool)
			}

			ipconfig = nutanixV3.IpConfig{
				DefaultGatewayIp: data["default_gateway_ip"].(string),
				PrefixLength:     int64(data["prefix_length"].(int)),
				SubnetIp:         data["subnet_ip"].(string),
				PoolList:         pool_range,
			}
		}
		subnet.Spec.Resources.IpConfig = ipconfig
	}

	if dhcpOptionsData, ok := d.GetOk("dhcp_options"); ok {
		if _, ok := d.GetOk("ip_config"); !ok {
			return fmt.Errorf("Invalid or empty ip_config subnet specification")
		}
		var dhcpOptions nutanixV3.DhcpOptions
		var dhcpServerAddress nutanixV3.Address
		params := dhcpOptionsData.(*schema.Set).List()
		for k := range params {
			data := params[k].(map[string]interface{})

			dhcpServerAddress.Ip = data["dhcp_server_address_host"].(string)
			dhcpServerAddress.Port = int64(data["dhcp_server_address_port"].(int))

			dhcpOptions.BootFileName = data["boot_file_name"].(string)
			dhcpOptions.DomainName = data["domain_name"].(string)
			dhcpOptions.TftpServerName = data["tftp_server_name"].(string)
			var dns_server_list []string
			for _, server := range data["domain_name_server_list"].(*schema.Set).List() {
				dns_server_list = append(dns_server_list, server.(string))
			}
			var dns_namerserver_list []string
			for _, name_server := range data["domain_search_list"].(*schema.Set).List() {
				dns_namerserver_list = append(dns_namerserver_list, name_server.(string))
			}
			dhcpOptions.DomainNameServerList = dns_server_list
			dhcpOptions.DomainSearchList = dns_namerserver_list
		}
		subnet.Spec.Resources.IpConfig.DhcpServerAddress = dhcpServerAddress
		subnet.Spec.Resources.IpConfig.DhcpOptions = dhcpOptions
	}

	subnet_uuid, err := client.SubnetExists(d.Get("name").(string))
	if err != nil {
		return err
	}

	if subnet_uuid != "" {
		return fmt.Errorf("Subnet already with name %s exists in the given cluster, UUID %s", d.Get("name").(string), subnet_uuid)
	}
	/*
		if client.Categories != nil {
			Categories := make(map[string]string)
			for key, value := range client.Categories {
				Categories[key] = value.(string)
			}
			subnet.Metadata.Categories = Categories
		}
	*/
	subnet_json, _ := json.Marshal(subnet)
	log.Printf("[DEBUG] Subnet JSON :%s", subnet_json)

	SubnetIntentResponse, APIResponse, err := SubnetAPIInstance.SubnetsPost(subnet)

	if err != nil {
		return err
	}

	err = checkAPIResponse(*APIResponse)
	if err != nil {
		return err
	}

	uuid := SubnetIntentResponse.Metadata.Uuid
	status, err := client.WaitForSubnetProcess(uuid)
	for status != true {
		return err
	}
	d.SetId(uuid)
	return resourceNutanixSubnetRead(d, meta)
}

func resourceNutanixSubnetRead(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] Reading Subnet: %s", d.Get("name").(string))

	client := meta.(*V3Client)
	SubnetAPIInstance := SubnetAPIInstance(client)
	uuid := d.Id()

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

	d.Set("default_gateway_ip", get_subnet.Status.Resources.IpConfig.DefaultGatewayIp)
	d.Set("prefix_length", get_subnet.Status.Resources.IpConfig.PrefixLength)
	d.Set("subnet_ip", get_subnet.Status.Resources.IpConfig.SubnetIp)
	return nil
}

func resourceNutanixSubnetUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating Subnet:%s", d.Id())

	d.Partial(true)
	client := meta.(*V3Client)
	SubnetAPIInstance := SubnetAPIInstance(client)

	uuid := d.Id()
	get_subnet, APIResponse, err := SubnetAPIInstance.SubnetsUuidGet(uuid)
	if err != nil {
		return err
	}

	var subnet nutanixV3.SubnetIntentInput

	subnet.Spec = get_subnet.Spec
	subnet.Metadata = get_subnet.Metadata
	subnet.ApiVersion = API_VERSION

	if d.HasChange("name") {
		subnet.Spec.Name = d.Get("name").(string)
		subnet.Metadata.Name = d.Get("name").(string)
	}
	if d.HasChange("vlan_id") {
		vlan_id := int64(d.Get("vlan_id").(int))
		subnet.Spec.Resources.VlanId = &vlan_id
	}
	if d.HasChange("description") {
		subnet.Spec.Description = d.Get("description").(string)
	}
	if d.HasChange("ip_config") {
		var ipconfig nutanixV3.IpConfig
		params := d.Get("ip_config").(*schema.Set).List()
		for k := range params {
			data := params[k].(map[string]interface{})
			var pool_range []nutanixV3.IpPool
			for _, pool := range data["pool_range"].(*schema.Set).List() {
				ip_pool := nutanixV3.IpPool{Range_: pool.(string)}
				pool_range = append(pool_range, ip_pool)
			}

			ipconfig = nutanixV3.IpConfig{
				DefaultGatewayIp: data["default_gateway_ip"].(string),
				PrefixLength:     int64(data["prefix_length"].(int)),
				SubnetIp:         data["subnet_ip"].(string),
				PoolList:         pool_range,
			}
		}
		subnet.Spec.Resources.IpConfig = ipconfig
	}
	if d.HasChange("dhcp_options") {
		if _, ok := d.GetOk("ip_config"); !ok {
			return fmt.Errorf("Invalid or empty ip_config subnet specification")
		}
		var dhcpOptions nutanixV3.DhcpOptions
		var dhcpServerAddress nutanixV3.Address
		params := d.Get("dhcp_options").(*schema.Set).List()
		for k := range params {
			data := params[k].(map[string]interface{})

			dhcpServerAddress.Ip = data["dhcp_server_address_host"].(string)
			dhcpServerAddress.Port = int64(data["dhcp_server_address_port"].(int))

			dhcpOptions.BootFileName = data["boot_file_name"].(string)
			dhcpOptions.DomainName = data["domain_name"].(string)
			dhcpOptions.TftpServerName = data["tftp_server_name"].(string)
			var dns_server_list []string
			for _, server := range data["domain_name_server_list"].(*schema.Set).List() {
				dns_server_list = append(dns_server_list, server.(string))
			}
			var dns_namerserver_list []string
			for _, name_server := range data["domain_search_list"].(*schema.Set).List() {
				dns_namerserver_list = append(dns_namerserver_list, name_server.(string))
			}
			dhcpOptions.DomainNameServerList = dns_server_list
			dhcpOptions.DomainSearchList = dns_namerserver_list
		}
		subnet.Spec.Resources.IpConfig.DhcpServerAddress = dhcpServerAddress
		subnet.Spec.Resources.IpConfig.DhcpOptions = dhcpOptions
	}

	subnet_json, _ := json.Marshal(subnet)
	log.Printf("[DEBUG] Subnet JSON :%s", subnet_json)

	SubnetIntentResponse, APIResponse, err := SubnetAPIInstance.SubnetsUuidPut(uuid, subnet)
	if err != nil {
		return err
	}
	err = checkAPIResponse(*APIResponse)
	if err != nil {
		return err
	}
	d.Partial(false)
	d.SetId(SubnetIntentResponse.Metadata.Uuid)
	return resourceNutanixSubnetRead(d, meta)
}

func resourceNutanixSubnetDelete(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] Deleting Subnet: %s", d.Get("name").(string))

	client := meta.(*V3Client)
	SubnetAPIInstance := SubnetAPIInstance(client)
	uuid := d.Id()

	APIResponse, err := SubnetAPIInstance.SubnetsUuidDelete(uuid)
	if err != nil {
		return err
	}
	err = checkAPIResponse(*APIResponse)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
