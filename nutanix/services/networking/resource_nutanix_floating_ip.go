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
	FipDelayTime  = 2 * time.Second
	FipMinTimeout = 5 * time.Second
)

func ResourceNutanixFloatingIP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixFloatingIPCreate,
		ReadContext:   resourceNutanixFloatingIPRead,
		UpdateContext: resourceNutanixFloatingIPUpdate,
		DeleteContext: resourceNutanixFloatingIPDelete,
		Schema: map[string]*schema.Schema{
			"external_subnet_reference_uuid": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"external_subnet_reference_name"},
			},
			"external_subnet_reference_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"external_subnet_reference_uuid"},
			},
			"vm_nic_reference_uuid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_reference_uuid": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				RequiredWith:  []string{"private_ip"},
				ConflictsWith: []string{"vpc_reference_name"},
			},
			"vpc_reference_name": {
				Type:          schema.TypeString,
				Optional:      true,
				RequiredWith:  []string{"private_ip"},
				ConflictsWith: []string{"vpc_reference_uuid"},
			},
			"private_ip": {
				Type:     schema.TypeString,
				Optional: true,
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

func resourceNutanixFloatingIPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API

	request := &v3.FIPIntentInput{}
	spec := &v3.FloatingIPSpec{}
	res := &v3.FIPResource{}
	metadata := &v3.Metadata{}

	if err := getMetadataAttributes(d, metadata, "floating_ip"); err != nil {
		return diag.Errorf("error reading metadata for floating_ip %s", err)
	}

	if extSub, eok := d.GetOk("external_subnet_reference_uuid"); eok {
		res.ExternalSubnetReference = buildReference(extSub.(string), "subnet")
	}

	if extSubName, eok := d.GetOk("external_subnet_reference_name"); eok {
		extResp, er := findSubnetByName(conn, extSubName.(string), nil)
		if er != nil {
			return diag.FromErr(er)
		}

		res.ExternalSubnetReference = buildReference(*extResp.Metadata.UUID, "subnet")
	}

	if vmNic, vok := d.GetOk("vm_nic_reference_uuid"); vok {
		res.VMNICReference = buildReference(vmNic.(string), "vm_nic")
	}

	if vpc, ok := d.GetOk("vpc_reference_uuid"); ok {
		res.VPCReference = buildReference(vpc.(string), "vpc")
	}

	if vpcName, ok := d.GetOk("vpc_reference_name"); ok {
		vpcResp, ver := findVPCByName(ctx, conn, vpcName.(string))
		if ver != nil {
			return diag.FromErr(ver)
		}
		res.VPCReference = buildReference(*vpcResp.Metadata.UUID, "vpc")
	}

	if pri, pok := d.GetOk("private_ip"); pok {
		res.PrivateIP = utils.StringPtr(pri.(string))
	}

	spec.Resource = res
	request.Metadata = metadata
	request.Spec = spec

	resp, err := conn.V3.CreateFloatingIPs(ctx, request)
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
		Delay:      FipDelayTime,
		MinTimeout: FipMinTimeout,
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for vpc (%s) to create: %s", uuid, errWaitTask)
	}

	d.SetId(uuid)
	return resourceNutanixFloatingIPRead(ctx, d, meta)
}

func resourceNutanixFloatingIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API

	resp, err := conn.V3.GetFloatingIPs(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	m, _ := setRSEntityMetadata(resp.Metadata)

	if err = d.Set("metadata", m); err != nil {
		return diag.Errorf("error setting metadata for Floating IP %s: %s", d.Id(), err)
	}

	if err = d.Set("api_version", resp.APIVersion); err != nil {
		return diag.FromErr(err)
	}
	if resp.Status.Resource.ExternalSubnetReference != nil {
		d.Set("external_subnet_reference_uuid", resp.Status.Resource.ExternalSubnetReference.UUID)
	}

	if resp.Status.Resource.VMNICReference != nil {
		d.Set("vm_nic_reference_uuid", resp.Status.Resource.VMNICReference.UUID)
	}

	if resp.Status.Resource.VPCReference != nil {
		d.Set("vpc_reference_uuid", resp.Status.Resource.VPCReference.UUID)
	}

	if err := d.Set("private_ip", resp.Status.Resource.PrivateIP); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*resp.Metadata.UUID)
	return nil
}

func resourceNutanixFloatingIPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API

	request := &v3.FIPIntentInput{}
	spec := &v3.FloatingIPSpec{}
	res := &v3.FIPResource{}
	metadata := &v3.Metadata{}

	resp, err := conn.V3.GetFloatingIPs(ctx, d.Id())
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return nil
		}
		return diag.Errorf("error reading Floating IP %s: %s", d.Id(), err)
	}

	if resp.Metadata != nil {
		metadata = resp.Metadata
	}

	if resp.Spec != nil {
		spec = resp.Spec

		if resp.Spec.Resource != nil {
			res = resp.Spec.Resource
		}
	}

	if d.HasChange("external_subnet_reference_uuid") {
		_, n := d.GetChange("external_subnet_reference_uuid")
		if len(n.(string)) > 0 {
			res.ExternalSubnetReference = buildReference(n.(string), "subnet")
		} else {
			res.ExternalSubnetReference = nil
		}
	}

	if d.HasChange("external_subnet_reference_name") {
		_, n := d.GetChange("external_subnet_reference_name")

		if len(n.(string)) > 0 {
			extResp, er := findSubnetByName(conn, n.(string), nil)
			if er != nil {
				return diag.FromErr(er)
			}
			res.ExternalSubnetReference = buildReference(*extResp.Metadata.UUID, "subnet")
		} else {
			res.ExternalSubnetReference = nil
		}
	}

	if d.HasChange("vm_nic_reference_uuid") {
		_, vm := d.GetChange("vm_nic_reference_uuid")
		res.VMNICReference = expandUpdateReference(vm.(string), "vm_nic")
	}

	if d.HasChange("vpc_reference_uuid") {
		_, vpc := d.GetChange("vpc_reference_uuid")

		if len(vpc.(string)) > 0 {
			res.VPCReference = expandUpdateReference(vpc.(string), "vpc")
		} else {
			res.VPCReference = nil
		}
	}

	if d.HasChange("vpc_reference_name") {
		_, vpc := d.GetChange("vpc_reference_name")

		if len(vpc.(string)) > 0 {
			vpcResp, ver := findVPCByName(ctx, conn, vpc.(string))
			if ver != nil {
				return diag.FromErr(ver)
			}
			res.VPCReference = expandUpdateReference(*vpcResp.Metadata.UUID, "vpc")
		} else {
			res.VPCReference = nil
		}
	}

	if d.HasChange("private_ip") {
		if len(d.Get("private_ip").(string)) > 0 {
			res.PrivateIP = utils.StringPtr(d.Get("private_ip").(string))
		} else {
			res.PrivateIP = nil
		}
	}

	request.Metadata = metadata
	spec.Resource = res
	request.Spec = spec

	// request to update Floating IP
	response, err := conn.V3.UpdateFloatingIP(ctx, d.Id(), request)
	if err != nil {
		return diag.FromErr(err)
	}

	taskUUID := response.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the Floating IP to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING", "RUNNING", "QUEUED"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      FipDelayTime,
		MinTimeout: FipMinTimeout,
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for Floating IP (%s) to update: %s", d.Id(), errWaitTask)
	}

	return resourceNutanixFloatingIPRead(ctx, d, meta)
}

func resourceNutanixFloatingIPDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API

	log.Printf("[DEBUG] Deleting Floating IP: %s", d.Id())
	resp, err := conn.V3.DeleteFloatingIP(ctx, d.Id())
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return nil
		}
		return diag.Errorf("error while deleting Floating IP UUID(%s): %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETE_PENDING", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, resp.Status.ExecutionContext.TaskUUID.(string)),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      FipDelayTime,
		MinTimeout: FipMinTimeout,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf(
			"error waiting for floating ip (%s) to delete: %s", d.Id(), err.Error())
	}
	d.SetId("")
	return nil
}

func expandUpdateReference(uuid string, kind string) *v3.Reference {
	if len(uuid) > 0 {
		return &v3.Reference{
			Kind: utils.StringPtr(kind),
			UUID: utils.StringPtr(uuid),
		}
	}
	return nil
}
