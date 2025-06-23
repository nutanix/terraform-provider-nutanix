package networking

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const (
	PbrDelayTime  = 2 * time.Second
	PbrMinTimeout = 5 * time.Second
)

func ResourceNutanixPbr() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixPbrCreate,
		ReadContext:   resourceNutanixPbrRead,
		UpdateContext: resourceNutanixPbrUpdate,
		DeleteContext: resourceNutanixPbrDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"source": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"INTERNET", "ALL"}, false),
						},
						"subnet_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"prefix_length": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"destination": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"INTERNET", "ALL"}, false),
						},
						"subnet_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"prefix_length": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"protocol_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP", "ICMP", "PROTOCOL_NUMBER", "ALL"}, false),
			},
			"protocol_parameters": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"udp": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source_port_range_list":      portRangeSchema(),
									"destination_port_range_list": portRangeSchema(),
								},
							},
						},
						"tcp": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source_port_range_list":      portRangeSchema(),
									"destination_port_range_list": portRangeSchema(),
								},
							},
						},
						"icmp": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"icmp_type": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"icmp_code": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						"protocol_number": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"action": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"DENY", "PERMIT", "REROUTE"}, false),
			},
			"service_ip_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"vpc_reference_uuid": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"vpc_name"},
			},
			"vpc_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"vpc_reference_uuid"},
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
			"is_bidirectional": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceNutanixPbrCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API

	request := &v3.PbrIntentInput{}
	res := &v3.PbrResources{}
	spec := &v3.PbrSpec{}
	metadata := &v3.Metadata{}

	n, nok := d.GetOk("name")
	if nok {
		spec.Name = utils.StringPtr(n.(string))
	}
	spec.Name = utils.StringPtr(n.(string))

	if err := getMetadataAttributes(d, metadata, "routing_policy"); err != nil {
		return diag.Errorf("error reading metadata for PBR %s", err)
	}

	if err := getPbrResources(d, res); err != nil {
		return diag.FromErr(err)
	}

	if vpcName, vnok := d.GetOk("vpc_name"); vnok {
		vpcResp, er := findVPCByName(ctx, conn, vpcName.(string))
		if er != nil {
			return diag.FromErr(er)
		}

		res.VpcReference = buildReference(*vpcResp.Metadata.UUID, "vpc")
	}

	spec.Resources = res
	request.Metadata = metadata
	request.Spec = spec

	//make call to create pbr

	resp, err := conn.V3.CreatePBR(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid := *resp.Metadata.UUID
	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the PBR to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING", "RUNNING", "QUEUED"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      PbrDelayTime,
		MinTimeout: PbrMinTimeout,
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for pbr (%s) to create: %s", uuid, errWaitTask)
	}

	d.SetId(*resp.Metadata.UUID)
	return resourceNutanixPbrRead(ctx, d, meta)
}

func resourceNutanixPbrRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API

	// Get the PBR
	resp, err := conn.V3.GetPBR(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	m, _ := setRSEntityMetadata(resp.Metadata)

	if err = d.Set("metadata", m); err != nil {
		return diag.Errorf("error setting metadata for VPC %s: %s", d.Id(), err)
	}

	if err = d.Set("source", flattenSourceDest(resp.Spec.Resources.Source)); err != nil {
		return diag.Errorf("error setting source for PBR %s: %s", d.Id(), err)
	}

	if err = d.Set("destination", flattenSourceDest(resp.Spec.Resources.Destination)); err != nil {
		return diag.Errorf("error setting destination for PBR %s: %s", d.Id(), err)
	}

	if err = d.Set("protocol_parameters", flattenProtocolParams(resp.Spec.Resources.ProtocolParameters)); err != nil {
		return diag.Errorf("error setting protocol parameter for PBR %s: %s", d.Id(), err)
	}

	if err = d.Set("action", resp.Spec.Resources.Action.Action); err != nil {
		return diag.Errorf("error setting action for PBR %s: %s", d.Id(), err)
	}

	if err = d.Set("priority", resp.Spec.Resources.Priority); err != nil {
		return diag.Errorf("error setting priority for PBR %s: %s", d.Id(), err)
	}

	if err = d.Set("is_bidirectional", resp.Spec.Resources.IsBidirectional); err != nil {
		return diag.Errorf("error setting is_bidirectional for PBR %s: %s", d.Id(), err)
	}

	if err = d.Set("vpc_reference_uuid", resp.Spec.Resources.VpcReference.UUID); err != nil {
		return diag.Errorf("error setting source for vpc_reference_uuid %s: %s", d.Id(), err)
	}

	if err = d.Set("protocol_type", resp.Spec.Resources.ProtocolType); err != nil {
		return diag.Errorf("error setting protocol_type for PBR %s: %s", d.Id(), err)
	}

	if err = d.Set("service_ip_list", utils.StringSlice(resp.Spec.Resources.Action.ServiceIPList)); err != nil {
		return diag.Errorf("error setting service_ip_list for PBR %s: %s", d.Id(), err)
	}

	d.Set("name", resp.Spec.Name)
	d.Set("api_version", resp.APIVersion)

	return nil
}

func resourceNutanixPbrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API

	spec := &v3.PbrSpec{}
	request := &v3.PbrIntentInput{}
	res := &v3.PbrResources{}
	metadata := &v3.Metadata{}

	// Get call

	response, err := conn.V3.GetPBR(ctx, d.Id())
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return nil
		}
		return diag.Errorf("error reading VPC %s: %s", d.Id(), err)
	}
	if response.Metadata != nil {
		metadata = response.Metadata
	}

	if response.Spec != nil {
		spec = response.Spec

		if response.Spec.Resources != nil {
			res = response.Spec.Resources
		}
	}

	if d.HasChange("name") {
		spec.Name = utils.StringPtr(d.Get("name").(string))
	}

	if d.HasChange("priority") {
		res.Priority = utils.IntPtr(d.Get("priority").(int))
	}

	if d.HasChange("protocol_type") {
		res.ProtocolType = utils.StringPtr(d.Get("protocol_type").(string))
	}

	if d.HasChange("source") {
		_, n := d.GetChange("source")
		res.Source = expandSourDest(n.([]interface{}))
	}

	if d.HasChange("destination") {
		_, n := d.GetChange("destination")
		res.Destination = expandSourDest(n.([]interface{}))
	}

	if d.HasChange("protocol_parameters") {
		_, n := d.GetChange("protocol_parameters")
		res.ProtocolParameters = expandProtocolParameters(n.([]interface{}))
	}

	if d.HasChange("action") {
		actUpdated := &v3.PbrAction{}
		action := d.Get("action")
		actUpdated.Action = utils.StringPtr(action.(string))

		if sip, sok := d.GetOk("service_ip_list"); sok {
			subips := sip.([]interface{})
			sublist := make([]string, len(subips))
			for a := range subips {
				sublist[a] = subips[a].(string)
			}
			actUpdated.ServiceIPList = sublist
		}
		res.Action = actUpdated
	}

	if d.HasChange("vpc_reference_uuid") {
		res.VpcReference = buildReference(d.Get("vpc_reference_uuid").(string), "vpc")
	}

	if d.HasChange("vpc_name") {
		vpcName := d.Get("vpc_name")
		vpcResp, er := findVPCByName(ctx, conn, vpcName.(string))
		if er != nil {
			return diag.FromErr(er)
		}
		res.VpcReference = buildReference(*vpcResp.Metadata.UUID, "vpc")
	}

	if d.HasChange("is_bidirectional") {
		res.IsBidirectional = utils.BoolPtr(d.Get("is_bidirectional").(bool))
	}

	spec.Resources = res
	request.Metadata = metadata
	request.Spec = spec

	// Make request to the API
	resp, err := conn.V3.UpdatePBR(ctx, d.Id(), request)
	if err != nil {
		return diag.FromErr(err)
	}

	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the PBR to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING", "RUNNING", "QUEUED"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      PbrDelayTime,
		MinTimeout: PbrMinTimeout,
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for PBR (%s) to update: %s", d.Id(), errWaitTask)
	}

	return resourceNutanixPbrRead(ctx, d, meta)
}

func resourceNutanixPbrDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API

	log.Printf("[DEBUG] Deleting PBR: %s, %s", d.Get("name").(string), d.Id())
	resp, err := conn.V3.DeletePBR(ctx, d.Id())
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return nil
		}
		return diag.Errorf("error while deleting PBR UUID(%s): %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETE_PENDING", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, resp.Status.ExecutionContext.TaskUUID.(string)),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      PbrDelayTime,
		MinTimeout: PbrMinTimeout,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf(
			"error waiting for vpc (%s) to delete: %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}

func getPbrResources(d *schema.ResourceData, pbr *v3.PbrResources) error {
	if sour, sok := d.GetOk("source"); sok {
		pbr.Source = expandSourDest(sour.([]interface{}))
	}

	if dest, dok := d.GetOk("destination"); dok {
		pbr.Destination = expandSourDest(dest.([]interface{}))
	}

	if proto, pok := d.GetOk("protocol_type"); pok {
		pbr.ProtocolType = utils.StringPtr(proto.(string))
	}

	if priorty, prok := d.GetOk("priority"); prok {
		pbr.Priority = utils.IntPtr(priorty.(int))
	}

	if bidirec, ok := d.GetOk("is_bidirectional"); ok {
		pbr.IsBidirectional = utils.BoolPtr(bidirec.(bool))
	}

	if protoParam, pok := d.GetOk("protocol_parameters"); pok {
		pbr.ProtocolParameters = expandProtocolParameters(protoParam.([]interface{}))
	}

	if action, aok := d.GetOk("action"); aok {
		act := &v3.PbrAction{}

		act.Action = utils.StringPtr(action.(string))

		if sip, sok := d.GetOk("service_ip_list"); sok {
			subips := sip.([]interface{})
			sublist := make([]string, len(subips))
			for a := range subips {
				sublist[a] = subips[a].(string)
			}
			act.ServiceIPList = sublist
		}
		pbr.Action = act
	}

	if vpc, vok := d.GetOk("vpc_reference_uuid"); vok {
		pbr.VpcReference = buildReference(vpc.(string), "vpc")
	}

	return nil
}

func expandSourDest(prs []interface{}) *v3.PbrSourDest {
	if len(prs) > 0 {
		res := &v3.PbrSourDest{}
		entry := prs[0].(map[string]interface{})
		if v1, ok1 := entry["address_type"]; ok1 && len(v1.(string)) > 0 {
			res.AddressType = utils.StringPtr(v1.(string))
		}

		if v2, ok2 := entry["subnet_ip"]; ok2 && len(v2.(string)) > 0 {
			subIP := &v3.PbrIPSubnet{}
			subIP.IP = utils.StringPtr(v2.(string))

			if v3, ok3 := entry["prefix_length"]; ok3 {
				subIP.PrefixLength = utils.IntPtr(v3.(int))
			}
			res.IPSubnet = subIP
		}
		return res
	}
	return nil
}

func expandProtocolParameters(prs []interface{}) *v3.PbrProtocolParams {
	if len(prs) > 0 {
		res := &v3.PbrProtocolParams{}
		ent := prs[0].(map[string]interface{})

		if pnum, pk := ent["protocol_number"]; pk && len(pnum.(string)) > 0 {
			val, _ := strconv.Atoi(pnum.(string))
			res.ProtocolNumber = utils.IntPtr(val)
		}

		if icmp, ok := ent["icmp"]; ok && len(icmp.([]interface{})) > 0 {
			icmpVal := &v3.PbrIcmp{}
			icmps := (icmp.([]interface{}))[0].(map[string]interface{})
			if code, cok := icmps["icmp_code"]; cok {
				icmpVal.IcmpCode = utils.IntPtr(code.(int))
			}

			if itype, tok := icmps["icmp_type"]; tok {
				icmpVal.IcmpType = utils.IntPtr(itype.(int))
			}
			res.Icmp = icmpVal
		}

		if udp, uok := ent["udp"]; uok && len(udp.([]interface{})) > 0 {
			res.UDP = expandRouteProtocol(udp.([]interface{}))
		}

		if tcp, tok := ent["tcp"]; tok && len(tcp.([]interface{})) > 0 {
			res.TCP = expandRouteProtocol(tcp.([]interface{}))
		}
		return res
	}
	return nil
}

func expandRouteProtocol(pr []interface{}) *v3.PortRangeList {
	if len(pr) > 0 {
		res := &v3.PortRangeList{}
		entry := pr[0].(map[string]interface{})

		if sp, sok := entry["source_port_range_list"]; sok {
			res.SourcePortRangeList = expandPortRangeList(sp)
		}

		if dp, dok := entry["destination_port_range_list"]; dok {
			res.DestinationPortRangeList = expandPortRangeList(dp)
		}
		return res
	}
	return nil
}

func flattenSourceDest(pr *v3.PbrSourDest) []interface{} {
	res := make([]interface{}, 0)
	if pr != nil {
		subDest := make(map[string]interface{})

		if pr.AddressType != nil {
			subDest["address_type"] = pr.AddressType
		}
		if pr.IPSubnet != nil {
			subDest["subnet_ip"] = pr.IPSubnet.IP
			subDest["prefix_length"] = pr.IPSubnet.PrefixLength
		}
		res = append(res, subDest)
	}
	return res
}

func flattenProtocolParams(pr *v3.PbrProtocolParams) []map[string]interface{} {
	if pr != nil {
		res := make([]map[string]interface{}, 0)
		proto := make(map[string]interface{})

		if pr.ProtocolNumber != nil {
			proto["protocol_number"] = strconv.Itoa(*pr.ProtocolNumber)
		}
		if pr.Icmp != nil {
			proto["icmp"] = flattenIcmpPbr(pr.Icmp)
		}
		if pr.TCP != nil {
			proto["tcp"] = flattenRouteSpec(pr.TCP)
		}
		if pr.UDP != nil {
			proto["udp"] = flattenRouteSpec(pr.UDP)
		}

		res = append(res, proto)
		return res
	}
	return nil
}

func flattenRouteSpec(pr *v3.PortRangeList) []interface{} {
	if pr != nil {
		res := make([]interface{}, 0)
		port := make(map[string]interface{})

		port["source_port_range_list"] = flattenPortRangeList(pr.SourcePortRangeList)
		port["destination_port_range_list"] = flattenPortRangeList(pr.DestinationPortRangeList)

		res = append(res, port)
		return res
	}
	return nil
}

func flattenIcmpPbr(pr *v3.PbrIcmp) []interface{} {
	if pr != nil {
		res := make([]interface{}, 0)
		icmp := make(map[string]interface{})

		icmp["icmp_code"] = pr.IcmpCode
		icmp["icmp_type"] = pr.IcmpType

		res = append(res, icmp)
		return res
	}
	return nil
}

func flattenPortRangeList(pr []*v3.PortRange) []map[string]interface{} {
	if len(pr) > 0 {
		res := make([]map[string]interface{}, len(pr))
		for i, v := range pr {
			item := make(map[string]interface{})
			item["end_port"] = v.EndPort
			item["start_port"] = v.StartPort
			res[i] = item
		}
		return res
	}
	return nil
}
