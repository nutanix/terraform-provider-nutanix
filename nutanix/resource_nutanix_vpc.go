package nutanix

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func resourceNutanixVPC() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixVPCCreate,
		ReadContext:   resourceNutanixVPCRead,
		UpdateContext: resourceNutanixVPCUpdate,
		DeleteContext: resourceNutanixVPCDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"api_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"categories": categoriesSchema(),
			"external_subnet_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external_subnet_reference": {
							Type:     schema.TypeMap,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						// "gateway_node_uuid_list": {
						// 	Type:     schema.TypeList,
						// 	Optional: true,
						// 	Computed: true,
						// 	Elem: &schema.Schema{
						// 		Type: schema.TypeString,
						// 	},
						// },
						// "active_gateway_node": {
						// 	Type:     schema.TypeList,
						// 	Optional: true,
						// 	Computed: true,
						// 	MaxItems: 1,
						// 	Elem: &schema.Resource{
						// 		Schema: map[string]*schema.Schema{
						// 			"host_reference": {
						// 				Type:     schema.TypeMap,
						// 				Required: true,
						// 				Elem: &schema.Schema{
						// 					Type: schema.TypeString,
						// 				},
						// 			},
						// 			"ip_address": {
						// 				Type:     schema.TypeString,
						// 				Optional: true,
						// 				Computed: true,
						// 			},
						// 		},
						// 	},
						// },
					},
				},
			},
			"externally_routable_prefix_list": {
				Type:     schema.TypeList,
				Optional: true,
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
			"common_domain_name_server_ip_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ipv6": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"fqdn": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"is_backup": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"external_subnet_list_status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external_subnet_reference": {
							Type:     schema.TypeMap,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"external_ip_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"active_gateway_node": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host_reference": {
										Type:     schema.TypeMap,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"ip_address": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceNutanixVPCCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).API

	request := &v3.VPCIntentInput{}
	spec := &v3.VPC{}
	metadata := &v3.Metadata{}
	res := &v3.VpcResources{}

	if n, nok := d.GetOk("name"); nok {
		spec.Name = utils.StringPtr(n.(string))
	}

	if err := getMetadataAttributes(d, metadata, "vpc"); err != nil {
		return diag.Errorf("error reading metadata for VPC %s", err)
	}

	if err := getVPCResources(d, res); err != nil {
		return diag.FromErr(err)
	}

	spec.Resources = res
	request.Metadata = metadata
	request.Spec = spec

	// Make request to the API
	resp, err := conn.V3.CreateVPC(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}
	uuid := *resp.Metadata.UUID
	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the VPC to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      2 * time.Second,
		MinTimeout: 5 * 60,
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for vpc (%s) to create: %s", uuid, errWaitTask)
	}

	d.SetId(uuid)
	return resourceNutanixVPCRead(ctx, d, meta)
}

func resourceNutanixVPCRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).API
	// Make request to the API
	resp, err := conn.V3.GetVPC(ctx, d.Id())

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return nil
		}
		return diag.Errorf("error reading VPC %s: %s", d.Id(), err)
	}

	m, c := setRSEntityMetadata(resp.Metadata)

	if err = d.Set("metadata", m); err != nil {
		return diag.Errorf("error setting metadata for VPC %s: %s", d.Id(), err)
	}
	if err := d.Set("categories", c); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("external_subnet_list", flattenExtSubnetList(resp.Spec.Resources.ExternalSubnetList)); err != nil {
		return diag.Errorf("error setting external_subnet_list for VPC %s: %s", d.Id(), err)
	}

	if err = d.Set("externally_routable_prefix_list", flattenExtRoutableList(resp.Spec.Resources.ExternallyRoutablePrefixList)); err != nil {
		return diag.Errorf("error setting externally_routable_prefix_list for VPC %s: %s", d.Id(), err)
	}

	if err = d.Set("common_domain_name_server_ip_list", flattenCommonDNSIPList(resp.Spec.Resources.CommonDomainNameServerIPList)); err != nil {
		return diag.Errorf("error setting externally_routable_prefix_list for VPC %s: %s", d.Id(), err)
	}

	if err = d.Set("external_subnet_list_status", flattenExtSubnetListStatus(resp.Status.Resources.ExternalSubnetList)); err != nil {
		return diag.Errorf("error setting external_subnet_list_status for VPC %s: %s", d.Id(), err)
	}

	d.Set("name", resp.Spec.Name)
	d.Set("api_version", resp.APIVersion)

	// d.SetId(*resp.Metadata.UUID)
	return nil
}
func resourceNutanixVPCUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
func resourceNutanixVPCDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func getVPCResources(d *schema.ResourceData, vpc *v3.VpcResources) error {

	if az, azok := d.GetOk("externally_routable_prefix_list"); azok {
		vpc.ExternallyRoutablePrefixList = expandExternallyRoutablePL(az)
	}

	if ext, extok := d.GetOk("external_subnet_list"); extok {
		vpc.ExternalSubnetList = expandExternalSubnet(ext)
	}

	if cmn, cmnok := d.GetOk("common_domain_name_server_ip_list"); cmnok {
		vpc.CommonDomainNameServerIPList = expandCommonDNSIPList(cmn)
	}
	return nil
}

func expandExternalSubnet(exs interface{}) []*v3.ExternalSubnetList {
	ex := exs.([]interface{})

	if len(ex) > 0 {
		dls := make([]*v3.ExternalSubnetList, len(ex))

		for k, val := range ex {
			v := val.(map[string]interface{})
			dl := &v3.ExternalSubnetList{}

			if v1, ok1 := v["external_subnet_reference"]; ok1 {
				dl.ExternalSubnetReference = validateShortRef(v1.(map[string]interface{}))
			}

			// if v2, ok1 := v["gateway_node_uuid_list"]; ok1 {
			// 	vs2 := v2.([]interface{})
			// 	dl.GatewayNodeUUIDList = expandStringList(vs2)
			// }

			// if v3, ok1 := v["active_gateway_node"]; ok1 {
			// 	dl.ActiveGatewayNode = expandActiveGateway(v3.([]interface{}))
			// }

			dls[k] = (dl)
		}
		return dls
	}
	return nil
}

func flattenExtSubnetListStatus(ext []*v3.ExternalSubnetList) []map[string]interface{} {
	if len(ext) > 0 {
		extSubStatusList := make([]map[string]interface{}, len(ext))

		for k, v := range ext {
			extSub := make(map[string]interface{})

			extSub["external_subnet_reference"] = flattenReferenceValues(v.ExternalSubnetReference)
			extSub["active_gateway_node"] = flattenActiveGatewayNode(v.ActiveGatewayNode)
			extSub["external_ip_list"] = v.ExternalIPList

			extSubStatusList[k] = extSub
		}
		return extSubStatusList
	}
	return nil
}

func flattenActiveGatewayNode(act *v3.ActiveGatewayNode) []interface{} {
	actNodeList := make([]interface{}, 0)

	if act != nil {
		actNode := make(map[string]interface{})
		actNode["host_reference"] = flattenReferenceValues(act.HostReference)
		actNode["ip_address"] = (act.IPAddress)

		actNodeList = append(actNodeList, actNode)
	}
	return actNodeList
}

func expandExternallyRoutablePL(ext interface{}) []*v3.ExternallyRoutablePrefixList {
	extP := ext.([]interface{})

	if len(extP) > 0 {
		extPL := make([]*v3.ExternallyRoutablePrefixList, len(extP))

		for k, val := range extP {
			v := val.(map[string]interface{})
			epl := &v3.ExternallyRoutablePrefixList{}

			if v1, ok1 := v["ip"]; ok1 {
				epl.IP = utils.StringPtr(v1.(string))
			}

			if v2, ok2 := v["prefix_length"]; ok2 {
				epl.PrefixLength = utils.IntPtr(v2.(int))
			}

			extPL[k] = epl
		}
		return extPL
	}
	return nil
}

func expandCommonDNSIPList(cms interface{}) []*v3.CommonDomainNameServerIPList {
	cmnDNS := cms.([]interface{})

	if len(cmnDNS) > 0 {
		cmnDNSIPList := make([]*v3.CommonDomainNameServerIPList, len(cmnDNS))

		for k, val := range cmnDNS {
			cmnDNSIP := &v3.CommonDomainNameServerIPList{}
			v := val.(map[string]interface{})

			if v1, ok1 := v["ip"]; ok1 {
				cmnDNSIP.IP = utils.StringPtr(v1.(string))
			}

			cmnDNSIPList[k] = cmnDNSIP
		}
		return cmnDNSIPList
	}
	return nil
}

func flattenExtSubnetList(ext []*v3.ExternalSubnetList) []map[string]interface{} {

	if len(ext) > 0 {
		extSub := make([]map[string]interface{}, len(ext))

		for k, v := range ext {
			exts := make(map[string]interface{})

			exts["external_subnet_reference"] = flattenReferenceValues(v.ExternalSubnetReference)
			extSub[k] = exts
		}
		return extSub
	}
	return nil
}

func flattenExtRoutableList(ext []*v3.ExternallyRoutablePrefixList) []map[string]interface{} {
	if len(ext) > 0 {
		extRout := make([]map[string]interface{}, len(ext))

		for k, v := range ext {
			exts := make(map[string]interface{})

			exts["ip"] = v.IP
			exts["prefix_length"] = v.PrefixLength

			extRout[k] = exts
		}
		return extRout
	}
	return nil
}

func flattenCommonDNSIPList(cmn []*v3.CommonDomainNameServerIPList) []map[string]interface{} {
	if len(cmn) > 0 {
		cmnDNSList := make([]map[string]interface{}, len(cmn))

		for k, v := range cmn {
			cmnDNS := make(map[string]interface{})

			cmnDNS["ip"] = v.IP

			cmnDNSList[k] = cmnDNS
		}
		return cmnDNSList
	}
	return nil
}
