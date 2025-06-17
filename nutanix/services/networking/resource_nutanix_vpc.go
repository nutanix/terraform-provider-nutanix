package networking

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const (
	VpcDelayTime     = 2 * time.Second
	VpcMinTimeout    = 5 * time.Second
	VpcDeleteTimeout = 1 * time.Minute
)

func ResourceNutanixVPC() *schema.Resource {
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
			"external_subnet_reference_uuid": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ConflictsWith: []string{"external_subnet_reference_name"},
			},
			"external_subnet_reference_name": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ConflictsWith: []string{"external_subnet_reference_uuid"},
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
	conn := meta.(*conns.Client).API

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

	if extname, extnameok := d.GetOk("external_subnet_reference_name"); extnameok {
		extnames := extname.([]interface{})
		subsList := make([]*v3.ExternalSubnetList, len(extnames))
		for k, v := range extnames {
			subs := &v3.ExternalSubnetList{}
			subResp, err := findSubnetByName(conn, v.(string), nil)
			if err != nil {
				return diag.FromErr(err)
			}
			subs.ExternalSubnetReference = buildReference(*subResp.Metadata.UUID, "subnet")
			subsList[k] = subs
		}
		res.ExternalSubnetList = subsList
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
		Pending:    []string{"PENDING", "RUNNING", "QUEUED"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      VpcDelayTime,
		MinTimeout: VpcMinTimeout,
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for vpc (%s) to create: %s", uuid, errWaitTask)
	}

	d.SetId(uuid)
	return resourceNutanixVPCRead(ctx, d, meta)
}

func resourceNutanixVPCRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API
	// Make request to the API
	resp, err := conn.V3.GetVPC(ctx, d.Id())
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return nil
		}
		return diag.Errorf("error reading VPC %s: %s", d.Id(), err)
	}

	m, _ := setRSEntityMetadata(resp.Metadata)

	if err = d.Set("metadata", m); err != nil {
		return diag.Errorf("error setting metadata for VPC %s: %s", d.Id(), err)
	}

	subList := make([]string, 0, len(resp.Status.Resources.ExternalSubnetList))
	for _, v := range resp.Status.Resources.ExternalSubnetList {
		subList = append(subList, *v.ExternalSubnetReference.UUID)
	}

	if err = d.Set("external_subnet_reference_uuid", subList); err != nil {
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

	return nil
}

func resourceNutanixVPCUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API

	request := &v3.VPCIntentInput{}
	spec := &v3.VPC{}
	metadata := &v3.Metadata{}
	res := &v3.VpcResources{}

	response, err := conn.V3.GetVPC(ctx, d.Id())
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

	if d.HasChange("external_subnet_reference_uuid") {
		newSubs := d.Get("external_subnet_reference_uuid")
		exts := newSubs.([]interface{})
		subsList := make([]*v3.ExternalSubnetList, len(exts))
		for k, v := range exts {
			subs := &v3.ExternalSubnetList{}
			subs.ExternalSubnetReference = buildReference(v.(string), "subnet")
			subsList[k] = subs
		}
		res.ExternalSubnetList = subsList
	}

	if d.HasChange("common_domain_name_server_ip_list") {
		res.CommonDomainNameServerIPList = expandCommonDNSIPList(d.Get("common_domain_name_server_ip_list"))
	}

	if d.HasChange("externally_routable_prefix_list") {
		res.ExternallyRoutablePrefixList = expandExternallyRoutablePL(d.Get("externally_routable_prefix_list"))
	}

	if d.HasChange("external_subnet_reference_name") {
		extname := d.Get("external_subnet_reference_name")
		extnames := extname.([]interface{})
		subsList := make([]*v3.ExternalSubnetList, len(extnames))
		for k, v := range extnames {
			subs := &v3.ExternalSubnetList{}
			subResp, er := findSubnetByName(conn, v.(string), nil)
			if er != nil {
				return diag.FromErr(er)
			}
			subs.ExternalSubnetReference = buildReference(*subResp.Metadata.UUID, "subnet")
			subsList[k] = subs
		}
		res.ExternalSubnetList = subsList
	}

	spec.Resources = res
	request.Metadata = metadata
	request.Spec = spec

	// Make request to the API
	resp, err := conn.V3.UpdateVPC(ctx, d.Id(), request)
	if err != nil {
		return diag.FromErr(err)
	}

	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the VPC to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING", "RUNNING", "QUEUED"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      VpcDelayTime,
		MinTimeout: VpcMinTimeout,
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for vpc (%s) to create: %s", d.Id(), errWaitTask)
	}

	return resourceNutanixVPCRead(ctx, d, meta)
}

func resourceNutanixVPCDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API
	log.Printf("[DEBUG] Deleting VPC: %s, %s", d.Get("name").(string), d.Id())
	resp, err := conn.V3.DeleteVPC(ctx, d.Id())
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return nil
		}
		return diag.Errorf("error while deleting VPC UUID(%s): %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETE_PENDING", "RUNNING", "QUEUED"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, resp.Status.ExecutionContext.TaskUUID.(string)),
		Timeout:    VpcDeleteTimeout,
		Delay:      VpcDelayTime,
		MinTimeout: VpcMinTimeout,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf(
			"error waiting for vpc (%s) to delete: %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}

func getVPCResources(d *schema.ResourceData, vpc *v3.VpcResources) error {
	if az, azok := d.GetOk("externally_routable_prefix_list"); azok {
		vpc.ExternallyRoutablePrefixList = expandExternallyRoutablePL(az)
	}

	if cmn, cmnok := d.GetOk("common_domain_name_server_ip_list"); cmnok {
		vpc.CommonDomainNameServerIPList = expandCommonDNSIPList(cmn)
	}

	if ext, extok := d.GetOk("external_subnet_reference_uuid"); extok {
		exts := ext.([]interface{})
		subsList := make([]*v3.ExternalSubnetList, len(exts))
		for k, v := range exts {
			subs := &v3.ExternalSubnetList{}
			subs.ExternalSubnetReference = buildReference(v.(string), "subnet")
			subsList[k] = subs
		}
		vpc.ExternalSubnetList = subsList
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
