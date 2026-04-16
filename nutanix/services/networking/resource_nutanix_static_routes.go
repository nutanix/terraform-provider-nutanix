package networking

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const (
	StaticRouteDelayTime  = 2 * time.Second
	StaticRouteMinTimeout = 5 * time.Second
)

func ResourceNutanixStaticRoute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixStaticRouteCreate,
		ReadContext:   resourceNutanixStaticRouteRead,
		UpdateContext: resourceNutanixStaticRouteUpdate,
		DeleteContext: resourceNutanixStaticRouteDelete,
		Schema: map[string]*schema.Schema{
			"vpc_uuid": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"vpc_name"},
			},
			"vpc_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"vpc_uuid"},
			},
			"static_routes_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination": {
							Type:     schema.TypeString,
							Required: true,
						},
						"external_subnet_reference_uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"vpn_connection_reference_uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"default_route_nexthop": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external_subnet_reference_uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
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
		},
	}
}

func resourceNutanixStaticRouteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API

	request := &v3.StaticRouteIntentInput{}
	spec := &v3.StaticRouteSpec{}
	res := &v3.StaticRouteResources{}
	metadata := &v3.Metadata{}

	// get the details for route table

	var vpcUUID string

	if vpcID, ok := d.GetOk("vpc_uuid"); ok {
		vpcUUID = vpcID.(string)
	}

	if vpcName, vnok := d.GetOk("vpc_name"); vnok {
		vpcResp, er := findVPCByName(ctx, conn, vpcName.(string))
		if er != nil {
			return diag.FromErr(er)
		}
		vpcUUID = *vpcResp.Metadata.UUID
	}

	resp, err := conn.V3.GetStaticRoute(ctx, vpcUUID)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.Metadata != nil {
		metadata = resp.Metadata
	}

	// request.Metadata = resp.Metadata
	if resp.Spec != nil {
		spec = resp.Spec

		if resp.Spec.Resources != nil {
			res = resp.Spec.Resources
		}
	}

	if er := getMetadataAttributes(d, metadata, "vpc_route_table"); er != nil {
		return diag.Errorf("error reading metadata for VPC Route table %s", er)
	}

	if stat, ok := d.GetOk("static_routes_list"); ok {
		res.StaticRoutesList = expandStaticRouteList(stat)
	}

	if def, ok := d.GetOk("default_route_nexthop"); ok {
		res.DefaultRouteNexthop = expandDefaultRoute(def)
	}

	spec.Resources = res
	request.Spec = spec
	request.Metadata = metadata

	// make a request to update the VPC Static Route Table
	response, err := conn.V3.UpdateStaticRoute(ctx, vpcUUID, request)
	if err != nil {
		return diag.FromErr(err)
	}

	taskUUID := response.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the Static Route to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING", "RUNNING", "QUEUED"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      StaticRouteDelayTime,
		MinTimeout: StaticRouteMinTimeout,
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for vpc Static Route (%s) to create: %s", d.Id(), errWaitTask)
	}

	d.SetId(vpcUUID)
	return resourceNutanixStaticRouteRead(ctx, d, meta)
}

func resourceNutanixStaticRouteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API

	// Get API to read the static routes

	resp, err := conn.V3.GetStaticRoute(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	m, _ := setRSEntityMetadata(resp.Metadata)

	if err = d.Set("metadata", m); err != nil {
		return diag.Errorf("error setting metadata for VPC %s: %s", d.Id(), err)
	}
	d.Set("api_version", resp.APIVersion)
	d.Set("vpc_uuid", d.Id())
	d.Set("static_routes_list", flattenStaticRouteList(resp.Spec.Resources.StaticRoutesList))
	d.Set("default_route_nexthop", flattendefaultRouteList(resp.Spec.Resources.DefaultRouteNexthop))
	return nil
}

func resourceNutanixStaticRouteUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API

	request := &v3.StaticRouteIntentInput{}
	spec := &v3.StaticRouteSpec{}
	res := &v3.StaticRouteResources{}
	metadata := &v3.Metadata{}

	// get the details for route table

	resp, err := conn.V3.GetStaticRoute(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.Metadata != nil {
		metadata = resp.Metadata
	}

	if resp.Spec != nil {
		spec = resp.Spec

		if resp.Spec.Resources != nil {
			res = resp.Spec.Resources
		}
	}

	if er := getMetadataAttributes(d, metadata, "vpc_route_table"); er != nil {
		return diag.Errorf("error reading metadata for VPC Route table %s", er)
	}

	if d.HasChange("static_routes_list") {
		res.StaticRoutesList = expandStaticRouteList(d.Get("static_routes_list"))
	}

	if d.HasChange("default_route_nexthop") {
		res.DefaultRouteNexthop = expandDefaultRoute(d.Get("default_route_nexthop"))
	}

	spec.Resources = res
	request.Spec = spec
	request.Metadata = metadata

	// make a request to update the VPC Static Route Table

	response, err := conn.V3.UpdateStaticRoute(ctx, d.Id(), request)
	if err != nil {
		return diag.FromErr(err)
	}

	taskUUID := response.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the Static Route to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING", "RUNNING", "QUEUED"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      StaticRouteDelayTime,
		MinTimeout: StaticRouteMinTimeout,
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for vpc Static Route (%s) to create: %s", d.Id(), errWaitTask)
	}

	return resourceNutanixStaticRouteRead(ctx, d, meta)
}

func resourceNutanixStaticRouteDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func expandStaticRouteList(pr interface{}) []*v3.StaticRoutesList {
	prI := pr.([]interface{})
	if len(prI) > 0 {
		prList := make([]*v3.StaticRoutesList, len(prI))
		for k, val := range prI {
			stat := &v3.StaticRoutesList{}
			v := val.(map[string]interface{})
			nexthop := &v3.NextHop{}
			if v1, ok1 := v["destination"]; ok1 {
				stat.Destination = utils.StringPtr(v1.(string))
			}
			if v, ok := v["external_subnet_reference_uuid"]; ok && len(v.(string)) > 0 {
				nexthop.ExternalSubnetReference = buildReference(v.(string), "subnet")
			}
			if v, ok := v["vpn_connection_reference_uuid"]; ok && len(v.(string)) > 0 {
				nexthop.VpnConnectionReference = buildReference(v.(string), "vpn_connection")
			}

			stat.NextHop = nexthop
			prList[k] = stat
		}
		return prList
	}
	return nil
}

func flattenStaticRouteList(stat []*v3.StaticRoutesList) []map[string]interface{} {
	if len(stat) > 0 {
		statList := make([]map[string]interface{}, len(stat))

		for k, v := range stat {
			stats := make(map[string]interface{})

			stats["destination"] = v.Destination
			if v.NextHop.ExternalSubnetReference != nil {
				stats["external_subnet_reference_uuid"] = v.NextHop.ExternalSubnetReference.UUID
			}
			if v.NextHop.VpnConnectionReference != nil {
				stats["vpn_connection_reference_uuid"] = v.NextHop.VpnConnectionReference.UUID
			}

			statList[k] = stats
		}
		return statList
	}
	return nil
}

func expandDefaultRoute(pr interface{}) *v3.NextHop {
	prI := pr.([]interface{})
	if len(prI) > 0 {
		nexthop := &v3.NextHop{}
		v := prI[0].(map[string]interface{})
		if v, ok := v["external_subnet_reference_uuid"]; ok && len(v.(string)) > 0 {
			nexthop.ExternalSubnetReference = buildReference(v.(string), "subnet")
		}
		return nexthop
	}
	return nil
}

func flattendefaultRouteList(pr *v3.NextHop) []map[string]interface{} {
	defList := make([]map[string]interface{}, 0)
	if pr != nil {
		defRoute := make(map[string]interface{})
		if pr.ExternalSubnetReference != nil {
			defRoute["external_subnet_reference_uuid"] = pr.ExternalSubnetReference.UUID
		}
		if pr.VpnConnectionReference != nil {
			defRoute["vpn_connection_reference_uuid"] = pr.VpnConnectionReference.UUID
		}

		defList = append(defList, defRoute)
	}
	return defList
}
