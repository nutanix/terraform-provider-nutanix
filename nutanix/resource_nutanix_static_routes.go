package nutanix

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const (
	StaticRouteDelayTime  = 2 * time.Second
	StaticRouteMinTimeout = 5 * time.Second
)

func resourceNutanixStaticRoute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixStaticRouteCreate,
		ReadContext:   resourceNutanixStaticRouteRead,
		UpdateContext: resourceNutanixStaticRouteUpdate,
		DeleteContext: resourceNutanixStaticRouteDelete,
		Schema: map[string]*schema.Schema{
			"vpc_uuid": {
				Type:     schema.TypeString,
				Required: true,
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
							Required: true,
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
	conn := meta.(*Client).API

	request := &v3.StaticRouteIntentInput{}
	spec := &v3.StaticRouteSpec{}
	res := &v3.StaticRouteResources{}
	metadata := &v3.Metadata{}

	// get the details for route table

	vpcUUID, ok := d.GetOk("vpc_uuid")
	if !ok {
		return diag.Errorf("vpc_uuid is required")
	}
	resp, err := conn.V3.GetStaticRoute(ctx, vpcUUID.(string))
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

	spec.Resources = res
	request.Spec = spec
	request.Metadata = metadata

	// make a request to update the VPC Static Route Table
	response, err := conn.V3.UpdateStaticRoute(ctx, vpcUUID.(string), request)
	if err != nil {
		return diag.FromErr(err)
	}

	taskUUID := response.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the Static Route to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      StaticRouteDelayTime,
		MinTimeout: StaticRouteMinTimeout,
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for vpc Static Route (%s) to create: %s", d.Id(), errWaitTask)
	}

	d.SetId(vpcUUID.(string))
	return resourceNutanixStaticRouteRead(ctx, d, meta)
}

func resourceNutanixStaticRouteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).API

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
	return nil
}

func resourceNutanixStaticRouteUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).API

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
		Pending:    []string{"PENDING", "RUNNING"},
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

			if v1, ok1 := v["destination"]; ok1 {
				stat.Destination = utils.StringPtr(v1.(string))
			}
			if v, ok := v["external_subnet_reference_uuid"]; ok {
				nexthop := &v3.NextHop{}

				nexthop.ExternalSubnetReference = buildReference(v.(string), "subnet")

				stat.NextHop = nexthop
			}
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
			stats["external_subnet_reference"] = v.NextHop.ExternalSubnetReference.UUID

			statList[k] = stats
		}
		return statList
	}
	return nil
}
