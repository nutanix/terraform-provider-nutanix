package nutanix

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func resourceNutanixSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceNutanixSubnetCreate,
		Read:   resourceNutanixSubnetRead,
		Update: resourceNutanixSubnetUpdate,
		Delete: resourceNutanixSubnetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: getSubnetSchema(),
	}
}

func resourceNutanixSubnetCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API

	request := &v3.SubnetIntentInput{}
	spec := &v3.Subnet{}
	metadata := &v3.Metadata{}
	subnet := &v3.SubnetResources{}

	n, nok := d.GetOk("name")
	azr, azrok := d.GetOk("availability_zone_reference")
	cr, crok := d.GetOk("cluster_reference")
	_, stok := d.GetOk("subnet_type")

	if !stok && !nok {
		return fmt.Errorf("please provide the required attributes name, subnet_type")
	}

	if v, ok := d.GetOk("api_version"); ok {
		request.APIVersion = utils.String(v.(string))
	}
	if !nok {
		return fmt.Errorf("please provide the required name attribute")
	}
	if err := getMetadataAttributes(d, metadata, "subnet"); err != nil {
		return err
	}

	if azrok {
		a := azr.(map[string]interface{})
		spec.AvailabilityZoneReference = validateRef(a)
	}
	if crok {
		a := cr.(map[string]interface{})
		spec.ClusterReference = validateRef(a)
	}

	if err := getSubnetResources(d, subnet); err != nil {
		return err
	}

	subnetUUID, err := resourceNutanixSubnetExists(conn, d.Get("name").(string))

	if err != nil {
		return err
	}

	if subnetUUID != nil {
		return fmt.Errorf("subnet already with name %s exists in the given cluster, UUID %s", d.Get("name").(string), *subnetUUID)
	}

	spec.Name = utils.String(n.(string))
	spec.Resources = subnet
	request.Metadata = metadata
	request.Spec = spec

	resp, err := conn.V3.CreateSubnet(request)
	if err != nil {
		return err
	}

	d.SetId(*resp.Metadata.UUID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING", "RUNNING"},
		Target:     []string{"COMPLETE"},
		Refresh:    subnetStateRefreshFunc(conn, d.Id()),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for subnet (%s) to create: %s", d.Id(), err)
	}

	return resourceNutanixSubnetRead(d, meta)
}

func resourceNutanixSubnetRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API

	resp, err := conn.V3.GetSubnet(d.Id())
	if err != nil {
		return err
	}

	m, c := setRSEntityMetadata(resp.Metadata)

	if err := d.Set("metadata", m); err != nil {
		return err
	}
	if err := d.Set("categories", c); err != nil {
		return err
	}
	if err := d.Set("project_reference", getReferenceValues(resp.Metadata.ProjectReference)); err != nil {
		return err
	}
	if err := d.Set("owner_reference", getReferenceValues(resp.Metadata.OwnerReference)); err != nil {
		return err
	}
	if err := d.Set("availability_zone_reference", getReferenceValues(resp.Status.AvailabilityZoneReference)); err != nil {
		return err
	}
	if err := d.Set("cluster_reference", getClusterReferenceValues(resp.Status.ClusterReference)); err != nil {
		return err
	}

	dgIP := ""
	sIP := ""
	pl := int64(0)
	port := int64(0)
	dhcpSA := make(map[string]interface{})
	dOptions := make(map[string]interface{})
	ipcpl := make([]string, 0)
	dnsList := make([]string, 0)
	dsList := make([]string, 0)

	if resp.Status.Resources.IPConfig != nil {
		dgIP = utils.StringValue(resp.Status.Resources.IPConfig.DefaultGatewayIP)
		pl = utils.Int64Value(resp.Status.Resources.IPConfig.PrefixLength)
		sIP = utils.StringValue(resp.Status.Resources.IPConfig.SubnetIP)

		if resp.Status.Resources.IPConfig.DHCPServerAddress != nil {
			dhcpSA["ip"] = utils.StringValue(resp.Status.Resources.IPConfig.DHCPServerAddress.IP)
			dhcpSA["fqdn"] = utils.StringValue(resp.Status.Resources.IPConfig.DHCPServerAddress.FQDN)
			dhcpSA["ipv6"] = utils.StringValue(resp.Status.Resources.IPConfig.DHCPServerAddress.IPV6)
			port = utils.Int64Value(resp.Status.Resources.IPConfig.DHCPServerAddress.Port)
		}

		if resp.Status.Resources.IPConfig.PoolList != nil {
			pl := resp.Status.Resources.IPConfig.PoolList
			poolList := make([]string, len(pl))
			for k, v := range pl {
				poolList[k] = utils.StringValue(v.Range)
			}
			ipcpl = poolList
		}
		if resp.Status.Resources.IPConfig.DHCPOptions != nil {
			dOptions["boot_file_name"] = utils.StringValue(resp.Status.Resources.IPConfig.DHCPOptions.BootFileName)
			dOptions["domain_name"] = utils.StringValue(resp.Status.Resources.IPConfig.DHCPOptions.DomainName)
			dOptions["tftp_server_name"] = utils.StringValue(resp.Status.Resources.IPConfig.DHCPOptions.TFTPServerName)

			if resp.Status.Resources.IPConfig.DHCPOptions.DomainNameServerList != nil {
				dnsList = utils.StringValueSlice(resp.Status.Resources.IPConfig.DHCPOptions.DomainNameServerList)
			}
			if resp.Status.Resources.IPConfig.DHCPOptions.DomainSearchList != nil {
				dsList = utils.StringValueSlice(resp.Status.Resources.IPConfig.DHCPOptions.DomainSearchList)
			}
		}
	}

	if err := d.Set("dhcp_server_address", dhcpSA); err != nil {
		return nil
	}
	if err := d.Set("ip_config_pool_list_ranges", ipcpl); err != nil {
		return nil
	}
	if err := d.Set("dhcp_options", dOptions); err != nil {
		return nil
	}
	if err := d.Set("dhcp_domain_name_server_list", dnsList); err != nil {
		return nil
	}
	if err := d.Set("dhcp_domain_search_list", dsList); err != nil {
		return nil
	}

	d.Set("cluster_reference_name", utils.StringValue(resp.Status.ClusterReference.Name))
	d.Set("api_version", utils.StringValue(resp.APIVersion))
	d.Set("name", utils.StringValue(resp.Status.Name))
	d.Set("description", utils.StringValue(resp.Status.Description))
	d.Set("state", utils.StringValue(resp.Status.State))
	d.Set("vswitch_name", utils.StringValue(resp.Status.Resources.VswitchName))
	d.Set("subnet_type", utils.StringValue(resp.Status.Resources.SubnetType))
	d.Set("default_gateway_ip", dgIP)
	d.Set("prefix_length", pl)
	d.Set("subnet_ip", sIP)
	d.Set("dhcp_server_address_port", port)
	d.Set("vlan_id", utils.Int64Value(resp.Status.Resources.VlanID))
	d.Set("network_function_chain_reference", getReferenceValues(resp.Status.Resources.NetworkFunctionChainReference))

	d.SetId(*resp.Metadata.UUID)

	return nil
}

func resourceNutanixSubnetUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API

	request := &v3.SubnetIntentInput{}
	metadata := &v3.Metadata{}
	res := &v3.SubnetResources{}
	ipcfg := &v3.IPConfig{}
	dhcpO := &v3.DHCPOptions{}
	spec := &v3.Subnet{}

	response, err := conn.V3.GetSubnet(d.Id())

	if err != nil {
		return err
	}

	if response.Metadata != nil {
		metadata = response.Metadata
	}

	if response.Spec != nil {
		spec = response.Spec

		if response.Spec.Resources != nil {
			res = response.Spec.Resources
			ipcfg = res.IPConfig
			if ipcfg != nil {
				dhcpO = ipcfg.DHCPOptions
			}
		}
	}

	if d.HasChange("categories") {
		catl := d.Get("categories").([]interface{})

		if len(catl) > 0 {
			cl := make(map[string]string)
			for _, v := range catl {
				item := v.(map[string]interface{})

				if i, ok := item["name"]; ok && i.(string) != "" {
					if k, kok := item["value"]; kok && k.(string) != "" {
						cl[i.(string)] = k.(string)
					}
				}
			}
			metadata.Categories = cl
		} else {
			metadata.Categories = nil
		}
	}
	if d.HasChange("owner_reference") {
		or := d.Get("owner_reference").(map[string]interface{})
		metadata.OwnerReference = validateRef(or)
	}
	if d.HasChange("project_reference") {
		pr := d.Get("project_reference").(map[string]interface{})
		metadata.ProjectReference = validateRef(pr)
	}
	if d.HasChange("name") {
		spec.Name = utils.String(d.Get("name").(string))
	}
	if d.HasChange("availability_zone_reference") {
		a := d.Get("availability_zone_reference").(map[string]interface{})
		spec.AvailabilityZoneReference = validateRef(a)
	}
	if d.HasChange("cluster_reference") {
		a := d.Get("cluster_reference").(map[string]interface{})
		spec.ClusterReference = validateRef(a)
	}
	if d.HasChange("dhcp_domain_name_server_list") {
		dd := d.Get("dhcp_domain_name_server_list").([]interface{})
		ddn := make([]*string, len(dd))
		for k, v := range dd {
			ddn[k] = utils.String(v.(string))
		}
		dhcpO.DomainNameServerList = ddn
	}
	if d.HasChange("dhcp_domain_search_list") {
		dd := d.Get("dhcp_domain_search_list").([]interface{})
		ddn := make([]*string, len(dd))
		for k, v := range dd {
			ddn[k] = utils.String(v.(string))
		}
		dhcpO.DomainSearchList = ddn
	}
	if d.HasChange("ip_config_pool_list_ranges") {
		dd := d.Get("ip_config_pool_list_ranges").([]interface{})
		ddn := make([]*v3.IPPool, len(dd))
		for k, v := range dd {
			i := &v3.IPPool{}
			i.Range = utils.String(v.(string))
			ddn[k] = i
		}
		ipcfg.PoolList = ddn
	}
	if d.HasChange("dhcp_options") {
		dOptions := d.Get("dhcp_options").(map[string]interface{})

		dhcpO.BootFileName = validateMapStringValue(dOptions, "boot_file_name")
		dhcpO.DomainName = validateMapStringValue(dOptions, "domain_name")
		dhcpO.TFTPServerName = validateMapStringValue(dOptions, "tftp_server_name")
	}
	if d.HasChange("network_function_chain_reference") {
		a := d.Get("network_function_chain_reference").(map[string]interface{})
		res.NetworkFunctionChainReference = validateRef(a)
	}
	if d.HasChange("vswitch_name") {
		res.VswitchName = utils.String(d.Get("vswitch_name").(string))
	}
	if d.HasChange("subnet_type") {
		res.SubnetType = utils.String(d.Get("subnet_type").(string))
	}
	if d.HasChange("default_gateway_ip") {
		ipcfg.DefaultGatewayIP = utils.String(d.Get("default_gateway_ip").(string))
	}
	if d.HasChange("prefix_length") {
		ipcfg.PrefixLength = utils.Int64(int64(d.Get("prefix_length").(int)))
	}
	if d.HasChange("subnet_ip") {
		ipcfg.SubnetIP = utils.String(d.Get("subnet_ip").(string))
	}
	if d.HasChange("dhcp_server_address") {
		dh := d.Get("dhcp_server_address").(map[string]interface{})

		ipcfg.DHCPServerAddress = &v3.Address{
			IP:   validateMapStringValue(dh, "ip"),
			IPV6: validateMapStringValue(dh, "ipv6"),
			FQDN: validateMapStringValue(dh, "fqdn"),
		}
	}
	if d.HasChange("dhcp_server_address_port") {
		ipcfg.DHCPServerAddress.Port = utils.Int64(int64(d.Get("dhcp_server_address_port").(int)))
	}
	if d.HasChange("vlan_id") {
		res.VlanID = utils.Int64(int64(d.Get("vlan_id").(int)))
	}

	ipcfg.DHCPOptions = dhcpO
	res.IPConfig = ipcfg
	spec.Resources = res
	request.Metadata = metadata
	request.Spec = spec

	if _, errUpdate := conn.V3.UpdateSubnet(d.Id(), request); errUpdate != nil {
		return errUpdate
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING", "RUNNING"},
		Target:     []string{"COMPLETE"},
		Refresh:    subnetStateRefreshFunc(conn, d.Id()),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"error waiting for subnet (%s) to update: %s", d.Id(), err)
	}
	return resourceNutanixSubnetRead(d, meta)
}

func resourceNutanixSubnetDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API

	if err := conn.V3.DeleteSubnet(d.Id()); err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING", "RUNNING", "DELETE_IN_PROGRESS", "COMPLETE"},
		Target:     []string{"DELETED"},
		Refresh:    subnetStateRefreshFunc(conn, d.Id()),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"error waiting for subnet (%s) to delete: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func resourceNutanixSubnetExists(conn *v3.Client, name string) (*string, error) {
	subnetEntities := &v3.DSMetadata{}
	var subnetUUID *string

	subnetList, err := conn.V3.ListSubnet(subnetEntities)

	if err != nil {
		return nil, err
	}

	for _, subnet := range subnetList.Entities {
		if subnet.Status.Name == utils.String(name) {
			subnetUUID = subnet.Metadata.UUID
		}
	}
	return subnetUUID, nil
}

func getSubnetResources(d *schema.ResourceData, subnet *v3.SubnetResources) error {

	ip := &v3.IPConfig{}
	dhcpo := &v3.DHCPOptions{}

	if v, ok := d.GetOk("vswitch_name"); ok {
		subnet.VswitchName = utils.String(v.(string))
	}
	if st, ok := d.GetOk("subnet_type"); ok {
		subnet.SubnetType = utils.String(st.(string))
	}
	if v, ok := d.GetOk("default_gateway_ip"); ok {
		ip.DefaultGatewayIP = utils.String(v.(string))
	}
	if v, ok := d.GetOk("prefix_length"); ok {
		ip.PrefixLength = utils.Int64(int64(v.(int)))
	}
	if v, ok := d.GetOk("subnet_ip"); ok {
		ip.SubnetIP = utils.String(v.(string))
	}
	if v, ok := d.GetOk("dhcp_server_address"); ok {
		dhcpa := v.(map[string]interface{})
		address := &v3.Address{}

		if ip, ok := dhcpa["ip"]; ok {
			address.IP = utils.String(ip.(string))
		}
		if fqdn, ok := dhcpa["fqdn"]; ok {
			address.FQDN = utils.String(fqdn.(string))
		}
		if v, ok := d.GetOk("dhcp_server_address_port"); ok {
			address.Port = utils.Int64(int64(v.(int)))
		}
		if ipv6, ok := dhcpa["ipv6"]; ok {
			address.IPV6 = utils.String(ipv6.(string))
		}

		ip.DHCPServerAddress = address
	}
	if v, ok := d.GetOk("ip_config_pool_list_ranges"); ok {
		p := v.([]interface{})
		pool := make([]*v3.IPPool, len(p))

		for k, v := range p {
			pItem := &v3.IPPool{}
			pItem.Range = utils.String(v.(string))
			pool[k] = pItem
		}

		ip.PoolList = pool
	}
	if v, ok := d.GetOk("dhcp_options"); ok {
		dop := v.(map[string]interface{})

		if boot, ok := dop["boot_file_name"]; ok {
			dhcpo.BootFileName = utils.String(boot.(string))
		}

		if dn, ok := dop["domain_name"]; ok {
			dhcpo.DomainName = utils.String(dn.(string))
		}

		if tsn, ok := dop["tftp_server_name"]; ok {
			dhcpo.TFTPServerName = utils.String(tsn.(string))
		}
	}

	if v, ok := d.GetOk("dhcp_domain_name_server_list"); ok {
		p := v.([]interface{})
		pool := make([]*string, len(p))

		for k, v := range p {
			pool[k] = utils.String(v.(string))
		}

		dhcpo.DomainNameServerList = pool
	}
	if v, ok := d.GetOk("dhcp_domain_search_list"); ok {
		p := v.([]interface{})
		pool := make([]*string, len(p))

		for k, v := range p {
			pool[k] = utils.String(v.(string))
		}

		dhcpo.DomainSearchList = pool
	}

	v, ok := d.GetOk("vlan_id")
	if v.(int) == 0 || ok {
		subnet.VlanID = utils.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("network_function_chain_reference"); ok {
		subnet.NetworkFunctionChainReference = validateRef(v.(map[string]interface{}))
	}

	ip.DHCPOptions = dhcpo

	subnet.IPConfig = ip

	return nil
}

func subnetStateRefreshFunc(client *v3.Client, uuid string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := client.V3.GetSubnet(uuid)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
				return v, DELETED, nil
			}
			return nil, "", err
		}

		return v, *v.Status.State, nil
	}
}

func getSubnetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"api_version": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"metadata": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"last_update_time": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"kind": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"creation_time": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"spec_version": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"spec_hash": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"categories": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"value": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"owner_reference": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"name": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"project_reference": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"name": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"availability_zone_reference": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Required: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Required: true,
					},
					"name": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"cluster_reference": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Required: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"cluster_name": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"vswitch_name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"subnet_type": {
			Type:     schema.TypeString,
			Required: true,
		},
		"default_gateway_ip": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"prefix_length": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"subnet_ip": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"dhcp_server_address": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"ip": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"fqdn": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"ipv6": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"dhcp_server_address_port": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"ip_config_pool_list_ranges": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"dhcp_options": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"boot_file_name": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
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
		"dhcp_domain_name_server_list": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"dhcp_domain_search_list": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"vlan_id": {
			Type:     schema.TypeInt,
			Optional: true,
			ForceNew: true,
			Computed: true,
		},
		"network_function_chain_reference": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Required: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Required: true,
					},
					"name": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
	}
}
