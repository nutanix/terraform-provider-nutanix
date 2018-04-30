package nutanix

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const (
	//SubnetKind represets the type of resource
	SubnetKind = "subnet"
)

func resourceNutanixSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceNutanixSubnetCreate,
		Read:   resourceNutanixSubnetRead,
		Update: resourceNutanixSubnetUpdate,
		Delete: resourceNutanixSubnetDelete,

		Schema: getSubnetSchema(),
	}
}

func resourceNutanixSubnetCreate(d *schema.ResourceData, meta interface{}) error {
	//Get client connection
	conn := meta.(*NutanixClient).API

	var version string
	if v, ok := d.GetOk("api_version"); ok {
		version = v.(string)
	} else {
		version = Version
	}

	// Prepare request
	request := &v3.SubnetIntentInput{
		APIVersion: utils.String(version),
		Spec: &v3.Subnet{
			Resources: &v3.SubnetResources{},
		},
	}

	//Read arguments and set request values
	m, mok := d.GetOk("metadata") //req
	n, nok := d.GetOk("name")     //req
	desc, descok := d.GetOk("description")
	azr, azrok := d.GetOk("availability_zone_reference")
	cr, crok := d.GetOk("cluster_reference")
	vswich, vok := d.GetOk("vswitch_name")
	subnet, stok := d.GetOk("subnet_type") //req
	ipconf, ipok := d.GetOk("ip_config")
	plr, plrok := d.GetOk("ip_config_pool_list_ranges")
	dhcp, dok := d.GetOk("dhcp_options")
	dnsl, dnslok := d.GetOk("dhcp_domain_name_server_list")
	dsl, dslok := d.GetOk("dhcp_domain_search_list")
	vlan, vlok := d.GetOk("vlan_id")
	nfcr, nfok := d.GetOk("network_function_chain_reference")
	pr, prok := d.GetOk("project_reference")
	or, orok := d.GetOk("owner_reference")

	if !mok && !stok && !nok {
		return fmt.Errorf("Please provide the required attributes metadata, name, subnet_type")
	}

	if azrok {
		a := azr.(map[string]interface{})
		r := &v3.Reference{
			Kind: utils.String(a["kind"].(string)),
			UUID: utils.String(a["uuid"].(string)),
		}
		if v, ok := a["name"]; ok {
			r.Name = utils.String(v.(string))
		}
		request.Spec.AvailabilityZoneReference = r
	}
	if descok {
		request.Spec.Description = utils.String(desc.(string))
	}
	if crok {
		a := cr.(map[string]interface{})
		r := &v3.Reference{
			Kind: utils.String(a["kind"].(string)),
			UUID: utils.String(a["uuid"].(string)),
		}
		if v, ok := a["name"]; ok {
			r.Name = utils.String(v.(string))
		}
		request.Spec.ClusterReference = r
	}

	if vok {
		request.Spec.Resources.VswitchName = utils.String(vswich.(string))
	}

	if stok {
		request.Spec.Resources.SubnetType = utils.String(subnet.(string))
	}

	//set ip_config
	request.Spec.Name = utils.String(n.(string))

	request.Metadata = setSubnetMetadata(m)

	if prok {

		pref := pr.(map[string]interface{})
		r := &v3.Reference{
			Kind: utils.String(pref["kind"].(string)),
			UUID: utils.String(pref["uuid"].(string)),
		}
		if v1, ok1 := pref["name"]; ok1 {
			r.Name = utils.String(v1.(string))
		}
		request.Metadata.ProjectReference = r
	}

	if orok {
		pr := or.(map[string]interface{})
		r := &v3.Reference{
			Kind: utils.String(pr["kind"].(string)),
			UUID: utils.String(pr["uuid"].(string)),
		}
		if v1, ok1 := pr["name"]; ok1 {
			r.Name = utils.String(v1.(string))
		}
		request.Metadata.OwnerReference = r
	}

	if ipok {
		ipConfig := setSubnetResourcesIPConfig(ipconf)

		if plrok {
			p := plr.([]interface{})

			pool := make([]*v3.IPPool, len(p))

			for k, v := range p {
				pItem := &v3.IPPool{}
				pItem.Range = utils.String(v.(string))
				pool[k] = pItem
			}

			ipConfig.PoolList = pool

		}

		if dok {
			ipConfig.DHCPOptions = setSubnetResourcesDHCPOptions(dhcp)

			//set domain_name_server_list
			if dnslok {
				dnslist := dnsl.([]interface{})

				domainNameServerList := make([]*string, len(dnslist))

				for k, v := range dnslist {
					domainNameServerList[k] = utils.String(v.(string))
				}

				ipConfig.DHCPOptions.DomainNameServerList = domainNameServerList
			}

			//set domain_name_search_list
			if dslok {
				dsList := dsl.([]interface{})

				domainSearchList := make([]*string, len(dsList))

				for k, v := range dsList {
					domainSearchList[k] = utils.String(v.(string))
				}

				ipConfig.DHCPOptions.DomainSearchList = domainSearchList
			}
		}
		request.Spec.Resources.IPConfig = ipConfig
	}

	//set vlan_id
	if vlok {
		request.Spec.Resources.VlanID = utils.Int64(int64(vlan.(int)))
	}

	// set network_function_chain_reference
	if nfok {
		ref := nfcr.(map[string]interface{})
		r := &v3.Reference{
			Kind: utils.String(ref["kind"].(string)),
			UUID: utils.String(ref["uuid"].(string)),
		}
		if v, ok := ref["name"]; ok {
			r.Name = utils.String(v.(string))
		}
		request.Spec.Resources.NetworkFunctionChainReference = r
	}

	subnetUUID, err := resourceNutanixSubnetExists(conn, d.Get("name").(string))

	if err != nil {
		return err
	}

	if subnetUUID != nil {
		return fmt.Errorf("Subnet already with name %s exists in the given cluster, UUID %s", d.Get("name").(string), *subnetUUID)
	}

	utils.PrintToJSON(request, "subnet request")

	//Make request to the API
	resp, err := conn.V3.CreateSubnet(request)
	if err != nil {
		return err
	}

	UUID := *resp.Metadata.UUID

	status, err := waitForSubnetProcess(conn, UUID)

	for status != true {
		return err
	}

	//set terraform state
	d.SetId(UUID)

	return resourceNutanixSubnetRead(d, meta)
}

func resourceNutanixSubnetRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading Subnet: %s", d.Get("name").(string))

	// Get client connection
	conn := meta.(*NutanixClient).API
	fmt.Printf("terraform-uuid %s", d.Id())
	// Make request to the API
	resp, err := conn.V3.GetSubnet(d.Id())

	if err != nil {
		return err
	}

	// Set subnet values
	// set availability zone reference values
	if resp.Status.AvailabilityZoneReference != nil {
		availabilityZoneReference := make(map[string]interface{})
		availabilityZoneReference["kind"] = utils.StringValue(resp.Status.AvailabilityZoneReference.Kind)
		availabilityZoneReference["name"] = utils.StringValue(resp.Status.AvailabilityZoneReference.Name)
		availabilityZoneReference["uuid"] = utils.StringValue(resp.Status.AvailabilityZoneReference.UUID)
		if err := d.Set("availability_zone_reference", availabilityZoneReference); err != nil {
			return err
		}
	}

	// set message list values
	if resp.Status.MessageList != nil {
		messages := make([]map[string]interface{}, len(resp.Status.MessageList))
		for k, v := range resp.Status.MessageList {
			message := make(map[string]interface{})

			message["message"] = *v.Message
			message["reason"] = *v.Reason
			message["details"] = v.Details

			messages[k] = message
		}
		if err := d.Set("message_list", messages); err != nil {
			return err
		}
	}

	// set cluster reference values
	if resp.Status.ClusterReference != nil {
		clusterReference := make(map[string]interface{})
		clusterReference["kind"] = *resp.Status.ClusterReference.Kind
		clusterReference["name"] = *resp.Status.ClusterReference.Name
		clusterReference["uuid"] = *resp.Status.ClusterReference.UUID

		if err := d.Set("cluster_reference", clusterReference); err != nil {
			return err
		}
	}

	// set metadata values
	metadata := make(map[string]interface{})
	metadata["last_update_time"] = resp.Metadata.LastUpdateTime.String()
	metadata["kind"] = utils.StringValue(resp.Metadata.Kind)
	metadata["uuid"] = utils.StringValue(resp.Metadata.UUID)
	metadata["creation_time"] = resp.Metadata.CreationTime.String()
	metadata["spec_version"] = strconv.FormatInt(utils.Int64Value(resp.Metadata.SpecVersion), 10) //convert to string
	metadata["spec_hash"] = utils.StringValue(resp.Metadata.SpecHash)
	//metadata["categories"] = resp.Metadata.Categories
	metadata["name"] = utils.StringValue(resp.Metadata.Name)

	if resp.Metadata.ProjectReference != nil {
		pr := make(map[string]interface{})
		pr["kind"] = utils.StringValue(resp.Metadata.ProjectReference.Kind)
		pr["name"] = utils.StringValue(resp.Metadata.ProjectReference.Name)
		pr["uuid"] = utils.StringValue(resp.Metadata.ProjectReference.UUID)

		if err := d.Set("project_reference", pr); err != nil {
			return err
		}
	}

	if resp.Metadata.OwnerReference != nil {
		or := make(map[string]interface{})
		or["kind"] = *resp.Metadata.OwnerReference.Kind
		or["name"] = *resp.Metadata.OwnerReference.Name
		or["uuid"] = *resp.Metadata.OwnerReference.UUID

		if err := d.Set("owner_reference", or); err != nil {
			return err
		}
	}

	// set ip_config
	ipConfig := make(map[string]interface{})

	ipConfig["default_gateway_ip"] = *resp.Status.Resources.IPConfig.DefaultGatewayIP

	//set ip_config.dhcp_server_address
	// address := make(map[string]interface{})

	// address["ip"] = utils.StringValue(resp.Status.Resources.IPConfig.DHCPServerAddress.IP)
	// address["fqdn"] = utils.StringValue(resp.Status.Resources.IPConfig.DHCPServerAddress.FQDN)
	// address["port"] = utils.Int64Value(resp.Status.Resources.IPConfig.DHCPServerAddress.Port)
	// address["ipv6"] = utils.StringValue(resp.Status.Resources.IPConfig.DHCPServerAddress.IPV6)

	// ipConfig["dhcp_server_address"] = address
	ipConfig["prefix_length"] = strconv.FormatInt(utils.Int64Value(resp.Status.Resources.IPConfig.PrefixLength), 10) //conver to string
	ipConfig["subnet_ip"] = utils.StringValue(resp.Status.Resources.IPConfig.SubnetIP)

	//set ip_config_pool_list_ranges
	pl := resp.Status.Resources.IPConfig.PoolList
	poolList := make([]string, len(pl))
	for k, v := range pl {
		poolList[k] = utils.StringValue(v.Range)
	}

	//set dhcp_options
	dOptions := make(map[string]interface{})

	dOptions["boot_file_name"] = utils.StringValue(resp.Status.Resources.IPConfig.DHCPOptions.BootFileName)
	dOptions["domain_name"] = utils.StringValue(resp.Status.Resources.IPConfig.DHCPOptions.DomainName)
	dOptions["tftp_server_name"] = utils.StringValue(resp.Status.Resources.IPConfig.DHCPOptions.TFTPServerName)
	ipConfig["default_gateway_ip"] = utils.StringValue(resp.Status.Resources.IPConfig.DefaultGatewayIP)

	//set dhcp_domain_name_server_list
	dnsl := resp.Status.Resources.IPConfig.DHCPOptions.DomainNameServerList
	dnsList := make([]string, len(dnsl))
	for k, v := range dnsl {
		dnsList[k] = utils.StringValue(v)
	}

	//set dhcp_domain_search_list
	dsl := resp.Status.Resources.IPConfig.DHCPOptions.DomainSearchList
	dsList := make([]string, len(dsl))
	for k, v := range dnsl {
		dsList[k] = utils.StringValue(v)
	}

	// set network_function_chain_reference
	if resp.Status.Resources.NetworkFunctionChainReference != nil {
		nfcr := make(map[string]interface{})
		nfcr["kind"] = utils.StringValue(resp.Status.Resources.NetworkFunctionChainReference.Kind)
		nfcr["name"] = utils.StringValue(resp.Status.Resources.NetworkFunctionChainReference.Name)
		nfcr["uuid"] = utils.StringValue(resp.Status.Resources.NetworkFunctionChainReference.UUID)

		if err := d.Set("network_function_chain_reference", nfcr); err != nil {
			return err
		}
	}
	if err := d.Set("api_version", resp.APIVersion); err != nil {
		return err
	}
	if err := d.Set("name", resp.Status.Name); err != nil {
		return err
	}
	if err := d.Set("state", resp.Status.State); err != nil {
		return err
	}
	if err := d.Set("description", resp.Status.Description); err != nil {
		return err
	}

	if err := d.Set("metadata", metadata); err != nil {
		return err
	}

	if err := d.Set("vswitch_name", resp.Status.Resources.VswitchName); err != nil {
		return err
	}

	if err := d.Set("subnet_type", resp.Status.Resources.SubnetType); err != nil {
		return err
	}

	if err := d.Set("ip_config", ipConfig); err != nil {
		return err
	}

	if err := d.Set("ip_config_pool_list_ranges", poolList); err != nil {
		return err
	}

	if err := d.Set("dhcp_options", dOptions); err != nil {
		return err
	}

	if err := d.Set("dhcp_domain_name_server_list", dnsList); err != nil {
		return err
	}

	if err := d.Set("dhcp_domain_search_list", dsList); err != nil {
		return err
	}

	if err := d.Set("vlan_id", resp.Status.Resources.VlanID); err != nil {
		return err
	}

	d.SetId(resource.UniqueId())

	return nil
}

func resourceNutanixSubnetUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NutanixClient).API

	// get state
	uuid := d.Id()
	m := d.Get("metadata").(map[string]interface{})
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	azr := d.Get("availability_zone_reference")
	cr := d.Get("cluster_reference")
	spec := d.Get("resources").(map[string]interface{})

	log.Printf("[DEBUG] Updating Subnet: %s, %s", name, uuid)

	subnetSpec, err := setSubnetResources(spec)

	if err != nil {
		return err
	}

	request := &v3.SubnetIntentInput{}

	if d.HasChange("metadata") {
		request.Metadata = setSubnetMetadata(m)
	}

	if d.HasChange("name") {
		request.Spec.Name = utils.String(name)
	}

	if d.HasChange("description") {
		request.Spec.Description = utils.String(description)
	}

	if d.HasChange("availability_zone_reference") {
		a := azr.(map[string]interface{})
		r := &v3.Reference{
			Kind: utils.String(a["kind"].(string)),
			UUID: utils.String(a["uuid"].(string)),
		}
		if v, ok := a["name"]; ok {
			r.Name = utils.String(v.(string))
		}

		request.Spec.AvailabilityZoneReference = r
	}

	if d.HasChange("cluster_reference") {
		a := cr.(map[string]interface{})
		r := &v3.Reference{
			Kind: utils.String(a["kind"].(string)),
			UUID: utils.String(a["uuid"].(string)),
		}
		if v, ok := a["name"]; ok {
			r.Name = utils.String(v.(string))
		}

		request.Spec.ClusterReference = r
	}

	d.Partial(true)

	if d.HasChange("resources") {
		request.Spec.Resources = subnetSpec
	}

	_, errUpdate := conn.V3.UpdateSubnet(uuid, request)
	if errUpdate != nil {
		return errUpdate
	}

	d.Partial(false)

	status, err := waitForSubnetProcess(conn, uuid)
	for status != true {
		return err
	}

	return resourceNutanixSubnetRead(d, meta)
}

func resourceNutanixSubnetDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NutanixClient).API
	UUID := d.Id()

	if err := conn.V3.DeleteSubnet(UUID); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceNutanixSubnetExists(conn *v3.Client, name string) (*string, error) {
	log.Printf("[DEBUG] Get Subnet Existance : %s", name)

	subnetEntities := &v3.SubnetListMetadata{}
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

func waitForSubnetProcess(conn *v3.Client, UUID string) (bool, error) {
	for {
		SubnetIntentResponse, err := conn.V3.GetSubnet(UUID)

		if err != nil {
			return false, err
		}

		if utils.StringValue(SubnetIntentResponse.Status.State) == "COMPLETE" {
			return true, nil
		} else if utils.StringValue(SubnetIntentResponse.Status.State) == "ERROR" {
			return false, fmt.Errorf("%s", utils.StringValue(SubnetIntentResponse.Status.MessageList[0].Message))
		}
		time.Sleep(3000 * time.Millisecond)
	}
}

func setSubnetResources(m interface{}) (*v3.SubnetResources, error) {

	subnet := &v3.SubnetResources{}

	resources := m.(map[string]interface{})

	if v, ok := resources["vswitch_name"]; ok {
		subnet.VswitchName = utils.String(v.(string))
	}

	st, stok := resources["subnet_type"]

	if !stok {
		return nil, fmt.Errorf("Plase provide required subnet_type attribute")
	}

	if vlan, ok := resources["vlan_id"]; ok {
		if n, err := strconv.Atoi(vlan.(string)); err == nil {
			subnet.VlanID = utils.Int64(int64(n))
		}

	}

	nfcr, nfcrok := resources["network_function_chain_reference"]

	if nfcrok {
		a := nfcr.(map[string]interface{})
		r := &v3.Reference{
			Kind: utils.String(a["kind"].(string)),
			UUID: utils.String(a["uuid"].(string)),
		}
		if v, ok := a["name"]; ok {
			r.Name = utils.String(v.(string))
		}
		subnet.NetworkFunctionChainReference = r
	}

	//ip config
	if ipcfg, ipcok := resources["ip_config"]; ipcok {
		ipc := ipcfg.([]interface{})
		if len(ipc) > 0 {
			cfg := ipc[0].(map[string]interface{})
			ipConf := &v3.IPConfig{}

			if d, ok := cfg["default_gateway_ip"]; ok {
				ipConf.DefaultGatewayIP = utils.String(d.(string))
			}

			if d, ok := cfg["prefix_length"]; ok {
				ipConf.PrefixLength = utils.Int64(d.(int64))
			}

			if d, ok := cfg["subnet_ip"]; ok {
				ipConf.SubnetIP = utils.String(d.(string))
			}

			if dhcp, dok := cfg["dhcp_server_address"]; dok {
				dhcpa := dhcp.(map[string]interface{})
				address := &v3.Address{}

				if ip, ok := dhcpa["ip"]; ok {
					address.IP = utils.String(ip.(string))
				}

				if fqdn, ok := dhcpa["fqdn"]; ok {
					address.FQDN = utils.String(fqdn.(string))
				}

				if port, ok := dhcpa["port"]; ok {
					address.Port = utils.Int64(int64(port.(int64)))
				}

				if ipv6, ok := dhcpa["ipv6"]; ok {
					address.IPV6 = utils.String(ipv6.(string))
				}

				ipConf.DHCPServerAddress = address
			}

			if pl, ok := cfg["pool_list"]; ok {
				p := pl.([]map[string]interface{})

				pool := make([]*v3.IPPool, len(p))

				for k, v := range p {
					pItem := &v3.IPPool{}
					if val, ok := v["range"]; ok {
						pItem.Range = utils.String(val.(string))
					}
					pool[k] = pItem
				}

				ipConf.PoolList = pool
			}

			if do, ok := cfg["dhcp_options"]; ok {
				dhcpo := &v3.DHCPOptions{}

				dop := do.(map[string]interface{})

				if dn, ok := dop["domain_name_server_list"]; ok {
					dnsl := dn.([]*string)

					domainNameServerList := make([]*string, len(dnsl))

					for k, v := range dnsl {
						domainNameServerList[k] = v
					}

					dhcpo.DomainNameServerList = domainNameServerList
				}

				if boot, ok := dop["boot_file_name"]; ok {
					dhcpo.BootFileName = utils.String(boot.(string))
				}

				if ds, ok := dop["domain_search_list"]; ok {
					dsl := ds.([]*string)

					domainSearchList := make([]*string, len(dsl))

					for k, v := range dsl {
						domainSearchList[k] = v
					}

					dhcpo.DomainSearchList = domainSearchList
				}

				if dn, ok := dop["domain_name"]; ok {
					dhcpo.DomainName = utils.String(dn.(string))
				}

				if tsn, ok := dop["tftp_server_name"]; ok {
					dhcpo.TFTPServerName = utils.String(tsn.(string))
				}

				ipConf.DHCPOptions = dhcpo
			}
			subnet.IPConfig = ipConf
		}

	}

	subnet.SubnetType = utils.String(st.(string))

	return subnet, nil
}

func setSubnetResourcesIPConfig(ic interface{}) *v3.IPConfig {
	cfg := ic.(map[string]interface{})

	ipConf := &v3.IPConfig{}

	if d, ok := cfg["default_gateway_ip"]; ok {
		ipConf.DefaultGatewayIP = utils.String(d.(string))
	}

	if d, ok := cfg["prefix_length"]; ok {
		if n, err := strconv.Atoi(d.(string)); err == nil {
			ipConf.PrefixLength = utils.Int64(int64(n))
		}

	}

	if d, ok := cfg["subnet_ip"]; ok {
		ipConf.SubnetIP = utils.String(d.(string))
	}

	if dhcp, dok := cfg["dhcp_server_address"]; dok {
		dhcpa := dhcp.(map[string]interface{})
		address := &v3.Address{}

		if ip, ok := dhcpa["ip"]; ok {
			address.IP = utils.String(ip.(string))
		}

		if fqdn, ok := dhcpa["fqdn"]; ok {
			address.FQDN = utils.String(fqdn.(string))
		}

		if port, ok := dhcpa["port"]; ok {
			address.Port = utils.Int64(int64(port.(int64)))
		}

		if ipv6, ok := dhcpa["ipv6"]; ok {
			address.IPV6 = utils.String(ipv6.(string))
		}

		ipConf.DHCPServerAddress = address
	}
	return ipConf
}

func setSubnetResourcesDHCPOptions(dhcp interface{}) *v3.DHCPOptions {
	dhcpo := &v3.DHCPOptions{}

	dop := dhcp.(map[string]interface{})

	if boot, ok := dop["boot_file_name"]; ok {
		dhcpo.BootFileName = utils.String(boot.(string))
	}

	if dn, ok := dop["domain_name"]; ok {
		dhcpo.DomainName = utils.String(dn.(string))
	}

	if tsn, ok := dop["tftp_server_name"]; ok {
		dhcpo.TFTPServerName = utils.String(tsn.(string))
	}

	return dhcpo

}

func setSubnetMetadata(m interface{}) *v3.SubnetMetadata {
	metad := m.(map[string]interface{})
	metadata := &v3.SubnetMetadata{
		Kind: utils.String(metad["kind"].(string)),
	}
	if v, ok := metad["uuid"]; ok {
		metadata.UUID = utils.String(v.(string))
	}
	if v, ok := metad["spec_version"]; ok {
		fmt.Println("TYPE")
		fmt.Println(reflect.TypeOf(v))
		if n, err := strconv.Atoi(v.(string)); err == nil {
			metadata.SpecVersion = utils.Int64(int64(n))
		}
	}
	if v, ok := metad["spec_hash"]; ok {
		metadata.SpecHash = utils.String(v.(string))
	}
	if v, ok := metad["name"]; ok {
		metadata.Name = utils.String(v.(string))
	}
	if v, ok := metad["categories"]; ok {
		metadata.Categories = v.(map[string]string)
	}

	return metadata
}

func getSubnetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"api_version": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"metadata": &schema.Schema{
			Type:     schema.TypeMap,
			Required: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"last_update_time": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
					"kind": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
					"uuid": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
					"creation_time": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
					"spec_version": &schema.Schema{
						Type:     schema.TypeInt,
						Optional: true,
					},
					"spec_hash": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
					"categories": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"owner_reference": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
					"uuid": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"project_reference": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
					"uuid": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"state": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"availability_zone_reference": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"uuid": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"message_list": &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"message": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"reason": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"details": &schema.Schema{
						Type:     schema.TypeMap,
						Computed: true,
					},
				},
			},
		},
		"cluster_reference": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"uuid": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"vswitch_name": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"subnet_type": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"ip_config": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"default_gateway_ip": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"dhcp_server_address": &schema.Schema{
						Type:     schema.TypeMap,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"ip": &schema.Schema{
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"fqdn": &schema.Schema{
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"port": &schema.Schema{
									Type:     schema.TypeInt,
									Optional: true,
									Computed: true,
								},
								"ipv6": &schema.Schema{
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
							},
						},
					},
					"prefix_length": &schema.Schema{
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
					"subnet_ip": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"ip_config_pool_list_ranges": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"dhcp_options": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"boot_file_name": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"domain_name": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"tftp_server_name": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"dhcp_domain_name_server_list": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"dhcp_domain_search_list": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"vlan_id": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"network_function_chain_reference": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"uuid": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
	}
}
