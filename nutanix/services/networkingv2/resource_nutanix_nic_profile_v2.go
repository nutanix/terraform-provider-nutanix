package networkingv2

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/request/nicprofiles"
	import3 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixNicProfileV2() *schema.Resource {
	resourceSchema := nicProfileSchema(false)

	return &schema.Resource{
		CreateContext: ResourceNutanixNicProfileV2Create,
		ReadContext:   ResourceNutanixNicProfileV2Read,
		UpdateContext: ResourceNutanixNicProfileV2Update,
		DeleteContext: ResourceNutanixNicProfileV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: resourceSchema,
	}
}

func ResourceNutanixNicProfileV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	inputSpec := expandNicProfile(map[string]interface{}{
		"name":              d.Get("name"),
		"description":       d.Get("description"),
		"capability_config": d.Get("capability_config"),
		"metadata":          d.Get("metadata"),
		"nic_family":        d.Get("nic_family"),
		"operating_mode":    d.Get("operating_mode"),
		"owner_type":        d.Get("owner_type"),
		"project_ext_id":    d.Get("project_ext_id"),
	})

	createRequest := import2.CreateNicProfileRequest{
		Body: &inputSpec,
	}

	aJSON, _ := json.MarshalIndent(createRequest, "", "  ")
	log.Printf("[DEBUG] Create NIC Profile Payload: %s", string(aJSON))

	resp, err := conn.NicProfilesAPI.CreateNicProfile(ctx, &createRequest)
	if err != nil {
		return diag.Errorf("error while creating NIC profile: %v", err)
	}

	taskVal := resp.Data.GetValue()
	aJSON, _ = json.MarshalIndent(taskVal, "", "  ")
	log.Printf("[DEBUG] Create NIC Profile Task Details: %s", string(aJSON))

	taskRef, err := extractNicProfileTaskReference(taskVal, "create nic profile")
	if err != nil {
		return diag.FromErr(err)
	}
	if err := waitForNicProfileTask(ctx, meta, utils.StringValue(taskRef.ExtId), d.Timeout(schema.TimeoutCreate), fmt.Sprintf("NIC profile (%s) creation", utils.StringValue(taskRef.ExtId))); err != nil {
		return diag.FromErr(err)
	}

	filter := fmt.Sprintf("name eq '%s'", d.Get("name").(string))
	listRequest := import2.ListNicProfilesRequest{
		Filter_: &filter,
	}
	listResp, err := conn.NicProfilesAPI.ListNicProfiles(ctx, &listRequest)
	if err != nil || listResp == nil || listResp.Data == nil {
		if err != nil {
			return diag.Errorf("NIC profile created but lookup by name failed: %v", err)
		}
		return diag.Errorf("NIC profile created but lookup by name returned no data")
	}

	switch v := listResp.Data.GetValue().(type) {
	case []import1.NicProfile:
		if len(v) > 0 && v[0].ExtId != nil {
			d.SetId(utils.StringValue(v[0].ExtId))
			if diags := reconcileNicProfileHostNicAssociations(ctx, meta, d.Id(), nil, expandStringSet(d.Get("host_nic_ext_ids")), d.Timeout(schema.TimeoutCreate)); diags != nil {
				return diags
			}
			return ResourceNutanixNicProfileV2Read(ctx, d, meta)
		}
	case []import1.NicProfileProjection:
		if len(v) > 0 && v[0].ExtId != nil {
			d.SetId(utils.StringValue(v[0].ExtId))
			if diags := reconcileNicProfileHostNicAssociations(ctx, meta, d.Id(), nil, expandStringSet(d.Get("host_nic_ext_ids")), d.Timeout(schema.TimeoutCreate)); diags != nil {
				return diags
			}
			return ResourceNutanixNicProfileV2Read(ctx, d, meta)
		}
	}

	return diag.Errorf("NIC profile created but ext_id could not be determined")
}

func ResourceNutanixNicProfileV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	getRequest := import2.GetNicProfileByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}
	resp, err := conn.NicProfilesAPI.GetNicProfileById(ctx, &getRequest)
	if err != nil {
		return diag.Errorf("error while fetching NIC profile: %v", err)
	}

	raw := resp.Data.GetValue()
	var getResp import1.NicProfile
	switch v := raw.(type) {
	case import1.NicProfile:
		getResp = v
	case *import1.NicProfile:
		if v == nil {
			return diag.Errorf("NIC profile response was nil")
		}
		getResp = *v
	default:
		return diag.Errorf("unexpected NIC profile response type: %T", raw)
	}

	if err := d.Set("ext_id", utils.StringValue(getResp.ExtId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", utils.StringValue(getResp.Name)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", utils.StringValue(getResp.Description)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("capability_config", flattenNicProfileCapabilityConfig(getResp.CapabilityConfig)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("host_nic_ext_ids", flattenNicProfileHostNicExtIDs(getResp.HostNicReferences)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("host_nic_references", flattenNicProfileHostNicReferences(getResp.HostNicReferences)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("metadata", flattenMetadata(getResp.Metadata)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("nic_family", utils.StringValue(getResp.NicFamily)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("operating_mode", common.FlattenPtrEnum(getResp.OperatingMode)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("owner_type", common.FlattenPtrEnum(getResp.OwnerType)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("project_ext_id", utils.StringValue(getResp.ProjectExtId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", utils.StringValue(getResp.TenantId)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ResourceNutanixNicProfileV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChange("capability_config") {
		return diag.Errorf("error while updating capability_config: Update of capability_config is not supported")
	}
	if d.HasChange("nic_family") {
		return diag.Errorf("error while updating nic_family: Update of nic_family is not supported")
	}
	if d.HasChange("operating_mode") {
		return diag.Errorf("error while updating operating_mode: Update of operating_mode is not supported")
	}
	if d.HasChange("owner_type") {
		return diag.Errorf("error while updating owner_type: Update of owner_type is not supported")
	}
	if d.HasChange("project_ext_id") {
		return diag.Errorf("error while updating project_ext_id: Update of project_ext_id is not supported")
	}

	if d.HasChange("name") || d.HasChange("description") || d.HasChange("metadata") {
		conn := meta.(*conns.Client).NetworkingAPI
		getRequest := import2.GetNicProfileByIdRequest{
			ExtId: utils.StringPtr(d.Id()),
		}
		resp, err := conn.NicProfilesAPI.GetNicProfileById(ctx, &getRequest)
		if err != nil {
			return diag.Errorf("error while fetching NIC profile: %v", err)
		}

		raw := resp.Data.GetValue()
		var updateSpec import1.NicProfile
		switch v := raw.(type) {
		case import1.NicProfile:
			updateSpec = v
		case *import1.NicProfile:
			if v == nil {
				return diag.Errorf("NIC profile response was nil")
			}
			updateSpec = *v
		default:
			return diag.Errorf("unexpected NIC profile response type: %T", raw)
		}

		if d.HasChange("name") {
			updateSpec.Name = utils.StringPtr(d.Get("name").(string))
		}
		if d.HasChange("description") {
			if v, ok := d.GetOk("description"); ok {
				updateSpec.Description = utils.StringPtr(v.(string))
			} else {
				updateSpec.Description = nil
			}
		}
		if d.HasChange("metadata") {
			updateSpec.Metadata = expandMetadata(d.Get("metadata").([]interface{}))
		}

		updateRequest := import2.UpdateNicProfileByIdRequest{
			ExtId: utils.StringPtr(d.Id()),
			Body:  &updateSpec,
		}

		aJSON, _ := json.MarshalIndent(updateRequest, "", "  ")
		log.Printf("[DEBUG] Update NIC Profile Payload: %s", string(aJSON))

		updateResp, err := conn.NicProfilesAPI.UpdateNicProfileById(ctx, &updateRequest)
		if err != nil {
			return diag.Errorf("error while updating NIC profile: %v", err)
		}

		taskVal := updateResp.Data.GetValue()
		aJSON, _ = json.MarshalIndent(taskVal, "", "  ")
		log.Printf("[DEBUG] Update NIC Profile Task Details: %s", string(aJSON))

		taskRef, err := extractNicProfileTaskReference(taskVal, "update nic profile")
		if err != nil {
			return diag.FromErr(err)
		}
		if err := waitForNicProfileTask(ctx, meta, utils.StringValue(taskRef.ExtId), d.Timeout(schema.TimeoutUpdate), fmt.Sprintf("NIC profile (%s) update", utils.StringValue(taskRef.ExtId))); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("host_nic_ext_ids") {
		oldValue, newValue := d.GetChange("host_nic_ext_ids")
		if diags := reconcileNicProfileHostNicAssociations(ctx, meta, d.Id(), expandStringSet(oldValue), expandStringSet(newValue), d.Timeout(schema.TimeoutUpdate)); diags != nil {
			return diags
		}
	}

	return ResourceNutanixNicProfileV2Read(ctx, d, meta)
}

func ResourceNutanixNicProfileV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	deleteRequest := import2.DeleteNicProfileByIdRequest{
		ExtId: utils.StringPtr(d.Id()),
	}
	resp, err := conn.NicProfilesAPI.DeleteNicProfileById(ctx, &deleteRequest)
	if err != nil {
		return diag.Errorf("error while deleting NIC profile: %v", err)
	}

	taskVal := resp.Data.GetValue()
	aJSON, _ := json.MarshalIndent(taskVal, "", "  ")
	log.Printf("[DEBUG] Delete NIC Profile Task Details: %s", string(aJSON))

	taskRef, err := extractNicProfileTaskReference(taskVal, "delete nic profile")
	if err != nil {
		return diag.FromErr(err)
	}
	if err := waitForNicProfileTask(ctx, meta, utils.StringValue(taskRef.ExtId), d.Timeout(schema.TimeoutDelete), fmt.Sprintf("NIC profile (%s) deletion", utils.StringValue(taskRef.ExtId))); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func reconcileNicProfileHostNicAssociations(ctx context.Context, meta interface{}, nicProfileExtID string, currentIDs []string, desiredIDs []string, timeout time.Duration) diag.Diagnostics {
	current := make(map[string]struct{}, len(currentIDs))
	for _, id := range currentIDs {
		current[id] = struct{}{}
	}

	desired := make(map[string]struct{}, len(desiredIDs))
	for _, id := range desiredIDs {
		desired[id] = struct{}{}
	}

	disassociateIDs := make([]string, 0)
	for _, id := range currentIDs {
		if _, ok := desired[id]; !ok {
			disassociateIDs = append(disassociateIDs, id)
		}
	}

	associateIDs := make([]string, 0)
	for _, id := range desiredIDs {
		if _, ok := current[id]; !ok {
			associateIDs = append(associateIDs, id)
		}
	}

	sort.Strings(disassociateIDs)
	sort.Strings(associateIDs)

	for _, hostNicExtID := range disassociateIDs {
		if _, err := disassociateHostNicFromNicProfile(ctx, meta, nicProfileExtID, hostNicExtID, timeout); err != nil {
			return diag.FromErr(err)
		}
	}

	for _, hostNicExtID := range associateIDs {
		if _, err := associateHostNicToNicProfile(ctx, meta, nicProfileExtID, hostNicExtID, timeout); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func associateHostNicToNicProfile(ctx context.Context, meta interface{}, nicProfileExtID string, hostNicExtID string, timeout time.Duration) (string, error) {
	conn := meta.(*conns.Client).NetworkingAPI

	request := import2.AssociateHostNicToNicProfileRequest{
		ExtId: utils.StringPtr(nicProfileExtID),
		Body: &import1.HostNic{
			HostNicExtId: utils.StringPtr(hostNicExtID),
		},
	}

	aJSON, _ := json.MarshalIndent(request, "", "  ")
	log.Printf("[DEBUG] Associate Host NIC to NIC Profile Payload: %s", string(aJSON))

	resp, err := conn.NicProfilesAPI.AssociateHostNicToNicProfile(ctx, &request)
	if err != nil {
		return "", fmt.Errorf("error while associating host NIC %s to NIC profile %s: %w", hostNicExtID, nicProfileExtID, err)
	}

	taskRef, err := extractNicProfileTaskReference(resp.Data.GetValue(), "associate host NIC to nic profile")
	if err != nil {
		return "", err
	}

	taskExtID := utils.StringValue(taskRef.ExtId)
	if err := waitForNicProfileTask(ctx, meta, taskExtID, timeout, fmt.Sprintf("host NIC %s association", hostNicExtID)); err != nil {
		return "", err
	}
	return taskExtID, nil
}

func disassociateHostNicFromNicProfile(ctx context.Context, meta interface{}, nicProfileExtID string, hostNicExtID string, timeout time.Duration) (string, error) {
	conn := meta.(*conns.Client).NetworkingAPI

	request := import2.DisassociateHostNicFromNicProfileRequest{
		ExtId: utils.StringPtr(nicProfileExtID),
		Body: &import1.HostNic{
			HostNicExtId: utils.StringPtr(hostNicExtID),
		},
	}

	aJSON, _ := json.MarshalIndent(request, "", "  ")
	log.Printf("[DEBUG] Disassociate Host NIC from NIC Profile Payload: %s", string(aJSON))

	resp, err := conn.NicProfilesAPI.DisassociateHostNicFromNicProfile(ctx, &request)
	if err != nil {
		return "", fmt.Errorf("error while disassociating host NIC %s from NIC profile %s: %w", hostNicExtID, nicProfileExtID, err)
	}

	taskRef, err := extractNicProfileTaskReference(resp.Data.GetValue(), "disassociate host NIC from nic profile")
	if err != nil {
		return "", err
	}

	taskExtID := utils.StringValue(taskRef.ExtId)
	if err := waitForNicProfileTask(ctx, meta, taskExtID, timeout, fmt.Sprintf("host NIC %s disassociation", hostNicExtID)); err != nil {
		return "", err
	}
	return taskExtID, nil
}

func extractNicProfileTaskReference(taskVal interface{}, operation string) (import3.TaskReference, error) {
	taskRef, ok := taskVal.(import3.TaskReference)
	if !ok {
		return import3.TaskReference{}, fmt.Errorf("unexpected %s task type: %T", operation, taskVal)
	}
	return taskRef, nil
}

func waitForNicProfileTask(ctx context.Context, meta interface{}, taskExtID string, timeout time.Duration, operation string) error {
	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, taskExtID),
		Timeout: timeout,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for %s task (%s): %w", operation, taskExtID, err)
	}
	return nil
}
