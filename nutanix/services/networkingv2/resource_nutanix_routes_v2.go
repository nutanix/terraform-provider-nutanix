package networkingv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	common "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/common/v1/config"
	"github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
	networkingPrism "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixRoutesV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixRoutesV2Create,
		ReadContext:   ResourceNutanixRoutesV2Read,
		UpdateContext: ResourceNutanixRoutesV2Update,
		DeleteContext: ResourceNutanixRoutesV2Delete,
		Schema: map[string]*schema.Schema{
			"route_table_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"metadata": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: DatasourceMetadataSchemaV4(),
				},
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"destination": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4": SchemaForValueRequiredPrefixLengthRequired(),
						"ipv6": SchemaForValueRequiredPrefixLengthRequired(),
					},
				},
			},
			"next_hop": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"next_hop_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"LOCAL_SUBNET", "DIRECT_CONNECT_VIF", "VPN_CONNECTION", "IP_ADDRESS", "EXTERNAL_SUBNET"}, false),
						},
						"next_hop_reference": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"next_hop_ip_address": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv4": SchemaForValueRequiredPrefixLength(),
									"ipv6": SchemaForValueRequiredPrefixLength(),
								},
							},
						},
						"next_hop_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"route_table_reference": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vpc_reference": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"external_routing_domain_reference": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"route_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"LOCAL", "STATIC", "DYNAMIC"}, false),
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rel": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"is_active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func SchemaForValueRequiredPrefixLength() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"value": {
					Type:     schema.TypeString,
					Required: true,
				},
				"prefix_length": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},
			},
		},
	}
}

func SchemaForValueRequiredPrefixLengthRequired() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ip": SchemaForValueRequiredPrefixLength(),
				"prefix_length": {
					Type:     schema.TypeInt,
					Computed: true,
					Optional: true,
				},
			},
		},
	}
}

func ResourceNutanixRoutesV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ResourceNutanixRoutesV2Create \n")
	conn := meta.(*conns.Client).NetworkingAPI

	routeTableExtID := d.Get("route_table_ext_id").(string)

	reqBody := &config.Route{}

	if metadata, ok := d.GetOk("metadata"); ok {
		reqBody.Metadata = expandMetadata(metadata.([]interface{}))
	}
	if name, ok := d.GetOk("name"); ok {
		reqBody.Name = utils.StringPtr(name.(string))
	}
	if description, ok := d.GetOk("description"); ok {
		reqBody.Description = utils.StringPtr(description.(string))
	}
	if destination, ok := d.GetOk("destination"); ok {
		reqBody.Destination = expandDestination(destination)
	}
	if nextHop, ok := d.GetOk("next_hop"); ok {
		reqBody.Nexthop = expandNextHop(nextHop)
	}
	if routeTableReference, ok := d.GetOk("route_table_reference"); ok {
		reqBody.RouteTableReference = utils.StringPtr(routeTableReference.(string))
	}
	if vpcReference, ok := d.GetOk("vpc_reference"); ok {
		reqBody.VpcReference = utils.StringPtr(vpcReference.(string))
	}
	if externalRoutingDomainReference, ok := d.GetOk("external_routing_domain_reference"); ok {
		reqBody.ExternalRoutingDomainReference = utils.StringPtr(externalRoutingDomainReference.(string))
	}
	if routeType, ok := d.GetOk("route_type"); ok {
		const two, three, four = 2, 3, 4
		routeTypeMap := map[string]interface{}{
			"DYNAMIC": two,
			"LOCAL":   three,
			"STATIC":  four,
		}
		pVal := routeTypeMap[routeType.(string)]
		p := config.RouteType(pVal.(int))
		reqBody.RouteType = &p
	}
	aJSON, _ := json.Marshal(reqBody)
	log.Printf("[DEBUG] Route Request Body: %v", string(aJSON))

	resp, err := conn.Routes.CreateRouteForRouteTable(&routeTableExtID, reqBody)
	if err != nil {
		return diag.Errorf("error while creating route for table : %v, error: %v", routeTableExtID, err)
	}

	taskRef := resp.Data.GetValue().(networkingPrism.TaskReference)
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the route table to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for route table (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching cretaet route task UUID : %v", err)
	}
	rUUID := resourceUUID.Data.GetValue().(prismConfig.Task)

	aJSON, _ = json.Marshal(rUUID)
	log.Printf("[DEBUG] create Route Task Details: %v", string(aJSON))

	uuid := rUUID.EntitiesAffected[0].ExtId
	d.SetId(*uuid)

	return ResourceNutanixRoutesV2Read(ctx, d, meta)
}

func ResourceNutanixRoutesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ResourceNutanixRoutesV2Read \n")
	conn := meta.(*conns.Client).NetworkingAPI

	routeTableExtID := d.Get("route_table_ext_id").(string)

	resp, err := conn.Routes.GetRouteForRouteTableById(utils.StringPtr(d.Id()), &routeTableExtID)
	if err != nil {
		return diag.Errorf("error while fetching route : %v", err)
	}

	getResp := resp.Data.GetValue().(config.Route)

	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("metadata", flattenMetadata(getResp.Metadata)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("destination", flattenDestination(getResp.Destination)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("next_hop", flattenNextHop(getResp.Nexthop)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("route_table_reference", getResp.RouteTableReference); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vpc_reference", getResp.VpcReference); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("external_routing_domain_reference", getResp.ExternalRoutingDomainReference); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("route_type", flattenRouteType(getResp.RouteType)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_active", getResp.IsActive); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("priority", getResp.Priority); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func ResourceNutanixRoutesV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ResourceNutanixRoutesV2Update \n")
	conn := meta.(*conns.Client).NetworkingAPI

	routeTableExtID := d.Get("route_table_ext_id").(string)

	// Get Etag
	routResp, err := conn.Routes.GetRouteForRouteTableById(utils.StringPtr(d.Id()), &routeTableExtID)
	if err != nil {
		return diag.Errorf("error while fetching route : %v", err)
	}

	// Extract E-Tag Header
	etagValue := conn.APIClientInstance.GetEtag(routResp)
	args := make(map[string]interface{})
	args["If-Match"] = utils.StringPtr(etagValue)

	getResp := routResp.Data.GetValue().(config.Route)

	updateSpec := &getResp

	if d.HasChange("metadata") {
		updateSpec.Metadata = expandMetadata(d.Get("metadata").([]interface{}))
	}
	if d.HasChange("name") {
		updateSpec.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("description") {
		updateSpec.Description = utils.StringPtr(d.Get("description").(string))
	}
	if d.HasChange("destination") {
		updateSpec.Destination = expandDestination(d.Get("destination"))
	}
	if d.HasChange("next_hop") {
		updateSpec.Nexthop = expandNextHop(d.Get("next_hop"))
	}
	if d.HasChange("route_table_reference") {
		updateSpec.RouteTableReference = utils.StringPtr(d.Get("route_table_reference").(string))
	}
	if d.HasChange("vpc_reference") {
		updateSpec.VpcReference = utils.StringPtr(d.Get("vpc_reference").(string))
	}
	if d.HasChange("external_routing_domain_reference") {
		updateSpec.ExternalRoutingDomainReference = utils.StringPtr(d.Get("external_routing_domain_reference").(string))
	}
	if d.HasChange("route_type") {
		const two, three, four = 2, 3, 4
		routeTypeMap := map[string]interface{}{
			"DYNAMIC": two,
			"LOCAL":   three,
			"STATIC":  four,
		}
		pVal := routeTypeMap[d.Get("route_type").(string)]
		p := config.RouteType(pVal.(int))
		updateSpec.RouteType = &p
	}

	aJSON, _ := json.Marshal(updateSpec)
	log.Printf("[DEBUG] Update Route Request Body: %v", string(aJSON))

	updateResp, err := conn.Routes.UpdateRouteForRouteTableById(utils.StringPtr(d.Id()), &routeTableExtID, updateSpec, args)
	if err != nil {
		return diag.Errorf("error while updating route : %v", err)
	}

	taskRef := updateResp.Data.GetValue().(networkingPrism.TaskReference)
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the route table to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for route table (%s) to perform: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching update route task UUID : %v", err)
	}

	rUUID := resourceUUID.Data.GetValue().(prismConfig.Task)

	aJSON, _ = json.Marshal(rUUID)
	log.Printf("[DEBUG] Update Route Task Details: %v", string(aJSON))

	return ResourceNutanixRoutesV2Read(ctx, d, meta)
}

func ResourceNutanixRoutesV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	routeTableExtID := d.Get("route_table_ext_id").(string)

	resp, err := conn.Routes.DeleteRouteForRouteTableById(utils.StringPtr(d.Id()), &routeTableExtID)
	if err != nil {
		return diag.Errorf("error while deleting route : %v", err)
	}

	taskRef := resp.Data.GetValue().(networkingPrism.TaskReference)
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the route table to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for route table (%s) to perform: %s", utils.StringValue(taskUUID), errWaitTask)
	}
	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching update route task UUID : %v", err)
	}

	rUUID := resourceUUID.Data.GetValue().(prismConfig.Task)

	aJSON, _ := json.Marshal(rUUID)
	log.Printf("[DEBUG] Update Route Task Details: %v", string(aJSON))

	return nil
}

func expandDestination(destination interface{}) *config.IPSubnet {
	if len(destination.([]interface{})) == 0 {
		log.Printf("[DEBUG] No destination found")
		return nil
	}
	destinationMap := destination.([]interface{})[0].(map[string]interface{})
	destinationObj := &config.IPSubnet{}
	aJSON, _ := json.Marshal(destinationMap)
	log.Printf("[DEBUG] Destination Map: %v", string(aJSON))
	if ipv4, ok := destinationMap["ipv4"]; ok && len(ipv4.([]interface{})) > 0 {
		destinationObj.Ipv4 = expandIPv4Subnet(ipv4)
	}
	if ipv6, ok := destinationMap["ipv6"]; ok && len(ipv6.([]interface{})) > 0 {
		destinationObj.Ipv6 = expandIPv6Subnet(ipv6)
	}
	aJSON, _ = json.Marshal(destinationObj)
	log.Printf("[DEBUG] Destination Object: %v", string(aJSON))
	return destinationObj
}

func expandNextHop(nextHop interface{}) *config.Nexthop {
	if len(nextHop.([]interface{})) == 0 {
		log.Printf("[DEBUG] No next hop found")
		return nil
	}
	nextHopMap := nextHop.([]interface{})[0].(map[string]interface{})
	nextHopObj := &config.Nexthop{}

	aJSON, _ := json.Marshal(nextHopMap)
	log.Printf("[DEBUG] Next Hop Map: %v", string(aJSON))

	if nextHopType, ok := nextHopMap["next_hop_type"]; ok {
		nextHopObj.NexthopType = expandNextHopType(nextHopType)
	}
	if nextHopReference, ok := nextHopMap["next_hop_reference"]; ok {
		nextHopObj.NexthopReference = utils.StringPtr(nextHopReference.(string))
	}
	if nextHopIPAddress, ok := nextHopMap["next_hop_ip_address"]; ok && len(nextHopIPAddress.([]interface{})) > 0 {
		nextHopObj.NexthopIpAddress = expandNextHopIPAddress(nextHopIPAddress)
	}
	log.Printf("[DEBUG] Next Hop Object: %v", nextHopObj)
	return nextHopObj
}

func expandNextHopIPAddress(address interface{}) *common.IPAddress {
	if len(address.([]interface{})) == 0 {
		log.Printf("[DEBUG] No next hop IP address found")
		return nil
	}
	addressMap := address.([]interface{})
	addressVal := addressMap[0].(map[string]interface{})
	addressObj := &common.IPAddress{}

	if ipv4, ok := addressVal["ipv4"]; ok && len(ipv4.([]interface{})) > 0 {
		addressObj.Ipv4 = expandIPv4Address(ipv4)
	}
	if ipv6, ok := addressVal["ipv6"]; ok && len(ipv6.([]interface{})) > 0 {
		addressObj.Ipv6 = expandIPv6Address(ipv6)
	}
	return addressObj
}

func expandNextHopType(hopType interface{}) *config.NexthopType {
	if hopType != nil {
		const two, three, four, five, six = 2, 3, 4, 5, 6
		nextHopTypeMap := map[string]interface{}{
			"IP_ADDRESS":         two,
			"DIRECT_CONNECT_VIF": three,
			"LOCAL_SUBNET":       four,
			"EXTERNAL_SUBNET":    five,
			"VPN_CONNECTION":     six,
		}
		pVal := nextHopTypeMap[hopType.(string)]
		p := config.NexthopType(pVal.(int))
		return &p
	}
	return nil
}

func expandMetadata(metadata []interface{}) *common.Metadata {
	if len(metadata) == 0 {
		log.Printf("[DEBUG] No metadata found")
		return nil
	}
	metadataMap := metadata[0].(map[string]interface{})
	metadataObj := &common.Metadata{}
	if ownerRefID, ok := metadataMap["owner_reference_id"]; ok {
		metadataObj.OwnerReferenceId = utils.StringPtr(ownerRefID.(string))
	}
	if ownerUserName, ok := metadataMap["owner_user_name"]; ok {
		metadataObj.OwnerUserName = utils.StringPtr(ownerUserName.(string))
	}
	if projRefID, ok := metadataMap["project_reference_id"]; ok {
		metadataObj.ProjectReferenceId = utils.StringPtr(projRefID.(string))
	}
	if projName, ok := metadataMap["project_name"]; ok {
		metadataObj.ProjectName = utils.StringPtr(projName.(string))
	}
	if categoryIDs, ok := metadataMap["category_ids"]; ok {
		categoryIDList := categoryIDs.([]interface{})
		categoryIDListStr := make([]string, len(categoryIDList))
		for i, v := range categoryIDList {
			categoryIDListStr[i] = v.(string)
		}
		metadataObj.CategoryIds = categoryIDListStr
	}
	aJSON, _ := json.Marshal(metadataObj)
	log.Printf("[DEBUG] Metadata Object: %v", string(aJSON))

	return metadataObj
}
