package objectstoresv2

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	objectsCommon "github.com/nutanix/ntnx-api-golang-clients/objects-go-client/v4/models/common/v1/config"
	"github.com/nutanix/ntnx-api-golang-clients/objects-go-client/v4/models/objects/v4/config"
	objectPrismConfig "github.com/nutanix/ntnx-api-golang-clients/objects-go-client/v4/models/prism/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const ipv4PrefixLengthDefaultValue = 32
const ipv6PrefixLengthDefaultValue = 128

func ResourceNutanixObjectStoresV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixObjectsV2Create,
		ReadContext:   ResourceNutanixObjectsV2Read,
		UpdateContext: ResourceNutanixObjectsV2Update,
		DeleteContext: ResourceNutanixObjectsV2Delete,
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(1 * time.Hour),
			Create:  schema.DefaultTimeout(1 * time.Hour),
			Update:  schema.DefaultTimeout(1 * time.Hour),
			Delete:  schema.DefaultTimeout(1 * time.Hour),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: customdiff.ComputedIf("state", func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
			if d.Id() == "" {
				return false
			}
			client := meta.(*conns.Client).ObjectStoreAPI
			resp, err := client.ObjectStoresAPIInstance.GetObjectstoreById(utils.StringPtr(d.Id()))
			if err != nil {
				return false
			}
			os := resp.Data.GetValue().(config.ObjectStore)
			// trigger a diff when deployment has failed
			return os.State.GetName() == "OBJECT_STORE_DEPLOYMENT_FAILED"
		}),

		Schema: map[string]*schema.Schema{
			"metadata": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: metadataSchema(),
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"deployment_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"num_worker_nodes": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"cluster_ext_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"storage_network_reference": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"storage_network_vip": {
				Type:     schema.TypeList,
				MaxItems: 1, //nolint:gomnd
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4": SchemaForValuePrefixLength(ipv4PrefixLengthDefaultValue),
						"ipv6": SchemaForValuePrefixLength(ipv6PrefixLengthDefaultValue),
					},
				},
			},
			"storage_network_dns_ip": {
				Type:     schema.TypeList,
				MaxItems: 1, //nolint:gomnd
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4": SchemaForValuePrefixLength(ipv4PrefixLengthDefaultValue),
						"ipv6": SchemaForValuePrefixLength(ipv6PrefixLengthDefaultValue),
					},
				},
			},
			"public_network_reference": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"public_network_ips": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4": SchemaForValuePrefixLength(ipv4PrefixLengthDefaultValue),
						"ipv6": SchemaForValuePrefixLength(ipv6PrefixLengthDefaultValue),
					},
				},
				Set: schema.HashResource(&schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4": SchemaForValuePrefixLength(ipv4PrefixLengthDefaultValue),
						"ipv6": SchemaForValuePrefixLength(ipv6PrefixLengthDefaultValue),
					},
				}),
			},
			"total_capacity_gib": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"DEPLOYING_OBJECT_STORE",
					"OBJECT_STORE_DEPLOYMENT_FAILED",
					"DELETING_OBJECT_STORE",
					"OBJECT_STORE_OPERATION_FAILED",
					"UNDEPLOYED_OBJECT_STORE",
					"OBJECT_STORE_OPERATION_PENDING",
					"OBJECT_STORE_AVAILABLE",
					"OBJECT_STORE_CERT_CREATION_FAILED",
					"CREATING_OBJECT_STORE_CERT",
					"OBJECT_STORE_DELETION_FAILED",
				}, false),
			},
			// computed attributes
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": linksSchema(),
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_ext_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func ResourceNutanixObjectsV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ObjectStoreAPI

	objectStorePayload := config.NewObjectStore()

	if v, ok := d.GetOk("metadata"); ok {
		objectStorePayload.Metadata = expandMetadata(v.([]interface{}))
	}
	if name, ok := d.GetOk("name"); ok {
		objectStorePayload.Name = utils.StringPtr(name.(string))
	}
	if description, ok := d.GetOk("description"); ok {
		objectStorePayload.Description = utils.StringPtr(description.(string))
	}
	if deploymentVersion, ok := d.GetOk("deployment_version"); ok {
		objectStorePayload.DeploymentVersion = utils.StringPtr(deploymentVersion.(string))
	}
	if domain, ok := d.GetOk("domain"); ok {
		objectStorePayload.Domain = utils.StringPtr(domain.(string))
	}
	if region, ok := d.GetOk("region"); ok {
		objectStorePayload.Region = utils.StringPtr(region.(string))
	}
	if numWorkerNodes, ok := d.GetOk("num_worker_nodes"); ok {
		objectStorePayload.NumWorkerNodes = utils.Int64Ptr(int64(numWorkerNodes.(int)))
	}
	if clusterExtID, ok := d.GetOk("cluster_ext_id"); ok {
		objectStorePayload.ClusterExtId = utils.StringPtr(clusterExtID.(string))
	}
	if storageNetworkRef, ok := d.GetOk("storage_network_reference"); ok {
		objectStorePayload.StorageNetworkReference = utils.StringPtr(storageNetworkRef.(string))
	}
	if storageNetworkVIP, ok := d.GetOk("storage_network_vip"); ok {
		objectStorePayload.StorageNetworkVip = &expandIPAddress(storageNetworkVIP.([]interface{}))[0]
	}
	if storageNetworkDNSIP, ok := d.GetOk("storage_network_dns_ip"); ok {
		objectStorePayload.StorageNetworkDnsIp = &expandIPAddress(storageNetworkDNSIP.([]interface{}))[0]
	}
	if publicNetworkRef, ok := d.GetOk("public_network_reference"); ok {
		objectStorePayload.PublicNetworkReference = utils.StringPtr(publicNetworkRef.(string))
	}
	if publicNetworkIPs, ok := d.GetOk("public_network_ips"); ok {
		objectStorePayload.PublicNetworkIps = expandIPAddress(publicNetworkIPs.(*schema.Set).List())
	}
	if totalCapacityGiB, ok := d.GetOk("total_capacity_gib"); ok {
		objectStorePayload.TotalCapacityGiB = utils.Int64Ptr(int64(totalCapacityGiB.(int)))
	}
	if state, ok := d.GetOk("state"); ok {
		objectStorePayload.State = expandState(state.(string))
	}

	aJSON, _ := json.MarshalIndent(objectStorePayload, "", "  ")
	log.Printf("[DEBUG] Object Store create payload: %s", string(aJSON))

	// change the timeout for the create operation
	d.Timeout(schema.TimeoutCreate)

	resp, err := conn.ObjectStoresAPIInstance.CreateObjectstore(objectStorePayload)
	if err != nil {
		return diag.Errorf("Error creating object store: %s", err)
	}

	TaskRef := resp.Data.GetValue().(objectPrismConfig.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the cluster to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		log.Printf("[DEBUG] deploy object store task error: %s", err)

		taskResp, taskErr := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
		if taskErr != nil {
			return diag.Errorf("error while fetch deploy object store task: %s", taskErr)
		}

		taskDetails := taskResp.Data.GetValue().(prismConfig.Task)

		// Get created object store extID from TASK API
		objectStoreExtID := taskDetails.EntitiesAffected[0].ExtId

		log.Printf("[DEBUG] object store extID: %s", utils.StringValue(objectStoreExtID))

		// Check if the object store is deployed or not
		// If not deployed, then return error
		// If deployed, save the object store details to the state
		// this code to maintain the state of the object store when its failed and the object store is present in the PC
		_, readErr := conn.ObjectStoresAPIInstance.GetObjectstoreById(objectStoreExtID)
		if readErr != nil {
			log.Printf("[DEBUG] object store not found")
			// If the object store is not found, object store is not deployed
			// and return the error
			return diag.Errorf("error waiting for object store to be deployed : %s", err)
		}
		// else, the object store instance exists in the system
	}

	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetch deploy object store task: %s", err)
	}

	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] deploy object store task details: %s", string(aJSON))

	// Get created object store extID from TASK API
	objectStoreExtID := taskDetails.EntitiesAffected[0].ExtId
	d.SetId(*objectStoreExtID)

	return ResourceNutanixObjectsV2Read(ctx, d, meta)
}

func ResourceNutanixObjectsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Reading object store %s", d.Id())
	conn := meta.(*conns.Client).ObjectStoreAPI

	readResp, err := conn.ObjectStoresAPIInstance.GetObjectstoreById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("Error reading object store: %s", err)
	}

	objectStore := readResp.Data.GetValue().(config.ObjectStore)

	if err := d.Set("tenant_id", objectStore.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ext_id", objectStore.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(objectStore.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("metadata", flattenMetadata(objectStore.Metadata)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", objectStore.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("creation_time", flattenTime(objectStore.CreationTime)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("last_update_time", flattenTime(objectStore.LastUpdateTime)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", objectStore.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("deployment_version", objectStore.DeploymentVersion); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("domain", objectStore.Domain); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("region", objectStore.Region); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("num_worker_nodes", objectStore.NumWorkerNodes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_ext_id", objectStore.ClusterExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("storage_network_reference", objectStore.StorageNetworkReference); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("storage_network_vip", flattenIPAddress([]objectsCommon.IPAddress{*objectStore.StorageNetworkVip})); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("storage_network_dns_ip", flattenIPAddress([]objectsCommon.IPAddress{*objectStore.StorageNetworkDnsIp})); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("public_network_reference", objectStore.PublicNetworkReference); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("public_network_ips", flattenIPAddress(objectStore.PublicNetworkIps)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("total_capacity_gib", objectStore.TotalCapacityGiB); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("state", objectStore.State.GetName()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("certificate_ext_ids", objectStore.CertificateExtIds); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ResourceNutanixObjectsV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ObjectStoreAPI

	readResp, err := conn.ObjectStoresAPIInstance.GetObjectstoreById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("Error reading object store: %s", err)
	}

	// Extract E-Tag Header
	args := make(map[string]interface{})
	etagValue := conn.ObjectStoresAPIInstance.ApiClient.GetEtag(readResp)
	args["If-Match"] = utils.StringPtr(etagValue)

	objectStoreUpdatePayload := readResp.Data.GetValue().(config.ObjectStore)

	// change the timeout for the update operation
	d.Timeout(schema.TimeoutUpdate)

	// resume the object store deployment if the state is OBJECT_STORE_DEPLOYMENT_FAILED
	resp, err := conn.ObjectStoresAPIInstance.UpdateObjectstoreById(utils.StringPtr(d.Id()), &objectStoreUpdatePayload, args)
	if err != nil {
		return diag.Errorf("Error updating object store: %s", err)
	}
	TaskRef := resp.Data.GetValue().(objectPrismConfig.TaskReference)
	taskUUID := TaskRef.ExtId
	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the cluster to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for object store to be updated : %s", err)
	}
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while update object store task: %s", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Object store Update task details: %s", string(aJSON))

	return ResourceNutanixObjectsV2Read(ctx, d, meta)
}

func ResourceNutanixObjectsV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ObjectStoreAPI

	readResp, err := conn.ObjectStoresAPIInstance.GetObjectstoreById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("Error reading object store: %s", err)
	}

	// Extract E-Tag Header
	args := make(map[string]interface{})
	etagValue := conn.ObjectStoresAPIInstance.ApiClient.GetEtag(readResp)
	args["If-Match"] = utils.StringPtr(etagValue)

	resp, err := conn.ObjectStoresAPIInstance.DeleteObjectstoreById(utils.StringPtr(d.Id()), args)
	if err != nil {
		return diag.Errorf("Error deleting object store: %s", err)
	}

	TaskRef := resp.Data.GetValue().(objectPrismConfig.TaskReference)
	taskUUID := TaskRef.ExtId
	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the object store to be deleted
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for object store to be deleted : %s", err)
	}
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching object store delete task: %s", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Object store delete task details: %s", string(aJSON))

	return nil
}

func metadataSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"owner_reference_id": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"owner_user_name": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"project_reference_id": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"project_name": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"category_ids": {
			Type:     schema.TypeList,
			Computed: true,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func linksSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"rel": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"href": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func SchemaForValuePrefixLength(defaultPrefixLength int) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1, //nolint:gomnd
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"value": {
					Type:     schema.TypeString,
					Required: true,
				},
				"prefix_length": {
					Type:     schema.TypeInt,
					Optional: true,
					Default:  defaultPrefixLength,
				},
			},
		},
	}
}

// expanders
func expandMetadata(metadata []interface{}) *objectsCommon.Metadata {
	if len(metadata) == 0 {
		log.Printf("[DEBUG] No metadata found")
		return nil
	}
	metadataMap := metadata[0].(map[string]interface{})
	metadataObj := &objectsCommon.Metadata{}
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
		metadataObj.CategoryIds = common.ExpandListOfString(categoryIDs.([]interface{}))
	}
	return metadataObj
}

func expandIPAddress(pr interface{}) []objectsCommon.IPAddress {
	if len(pr.([]interface{})) > 0 {
		ipAddressesList := make([]objectsCommon.IPAddress, len(pr.([]interface{})))

		for i, v := range pr.([]interface{}) {
			ipFilter := objectsCommon.IPAddress{}

			if v.(map[string]interface{})["ipv4"] != nil {
				ipFilter.Ipv4 = expandIPv4Address(v.(map[string]interface{})["ipv4"])
			}
			if v.(map[string]interface{})["ipv6"] != nil {
				ipFilter.Ipv6 = expandIPv6Address(v.(map[string]interface{})["ipv6"])
			}

			ipAddressesList[i] = ipFilter
		}
		return ipAddressesList
	}
	return nil
}

func expandIPv4Address(pr interface{}) *objectsCommon.IPv4Address {
	if len(pr.([]interface{})) == 0 {
		return nil
	}
	if pr != nil {
		ipv4 := objectsCommon.NewIPv4Address()
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if value, ok := val["value"]; ok {
			ipv4.Value = utils.StringPtr(value.(string))
		}
		if prefix, ok := val["prefix_length"]; ok {
			ipv4.PrefixLength = utils.IntPtr(prefix.(int))
		}
		return ipv4
	}
	return nil
}

func expandIPv6Address(pr interface{}) *objectsCommon.IPv6Address {
	if len(pr.([]interface{})) == 0 {
		return nil
	}

	if pr != nil {
		ipv6 := objectsCommon.NewIPv6Address()
		prI := pr.([]interface{})
		val := prI[0].(map[string]interface{})

		if value, ok := val["value"]; ok {
			ipv6.Value = utils.StringPtr(value.(string))
		}
		if prefix, ok := val["prefix_length"]; ok {
			ipv6.PrefixLength = utils.IntPtr(prefix.(int))
		}
		return ipv6
	}
	return nil
}

func expandState(state string) *config.State {
	var stateEnum config.State

	switch state {
	case "OBJECT_STORE_DEPLOYMENT_FAILED":
		stateEnum = config.STATE_OBJECT_STORE_DEPLOYMENT_FAILED
	case "OBJECT_STORE_OPERATION_FAILED":
		stateEnum = config.STATE_OBJECT_STORE_OPERATION_FAILED
	case "OBJECT_STORE_OPERATION_PENDING":
		stateEnum = config.STATE_OBJECT_STORE_OPERATION_PENDING
	case "OBJECT_STORE_AVAILABLE":
		stateEnum = config.STATE_OBJECT_STORE_AVAILABLE
	case "OBJECT_STORE_CERT_CREATION_FAILED":
		stateEnum = config.STATE_OBJECT_STORE_CERT_CREATION_FAILED
	case "CREATING_OBJECT_STORE_CERT":
		stateEnum = config.STATE_CREATING_OBJECT_STORE_CERT
	case "OBJECT_STORE_DELETION_FAILED":
		stateEnum = config.STATE_OBJECT_STORE_DELETION_FAILED
	case "DEPLOYING_OBJECT_STORE":
		stateEnum = config.STATE_DEPLOYING_OBJECT_STORE
	case "DELETING_OBJECT_STORE":
		stateEnum = config.STATE_DELETING_OBJECT_STORE
	case "UNDEPLOYED_OBJECT_STORE":
		stateEnum = config.STATE_UNDEPLOYED_OBJECT_STORE
	default:
		stateEnum = config.STATE_UNKNOWN
	}

	return &stateEnum
}

// func to check pc task status, and return the task status or error message
func taskStateRefreshPrismTaskGroupFunc(client *prism.Client, taskUUID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		taskResp, err := client.TaskRefAPI.GetTaskById(utils.StringPtr(taskUUID), nil)

		if err != nil {
			return "", "", fmt.Errorf("error while polling prism task: %v", err)
		}

		// get the group results
		v := taskResp.Data.GetValue().(prismConfig.Task)

		if getTaskStatus(v.Status) == "CANCELED" || getTaskStatus(v.Status) == "FAILED" {
			return v, getTaskStatus(v.Status),
				fmt.Errorf("error_detail: %s, progress_message: %d", utils.StringValue(v.ErrorMessages[0].Message), utils.IntValue(v.ProgressPercentage))
		}
		return v, getTaskStatus(v.Status), nil
	}
}

// func to flatten the task status to string
func getTaskStatus(pr *prismConfig.TaskStatus) string {
	if pr != nil {
		const QUEUED, RUNNING, SUCCEEDED, FAILED, CANCELED = 2, 3, 5, 6, 7
		if *pr == prismConfig.TaskStatus(FAILED) {
			return "FAILED"
		}
		if *pr == prismConfig.TaskStatus(CANCELED) {
			return "CANCELED"
		}
		if *pr == prismConfig.TaskStatus(QUEUED) {
			return "QUEUED"
		}
		if *pr == prismConfig.TaskStatus(RUNNING) {
			return "RUNNING"
		}
		if *pr == prismConfig.TaskStatus(SUCCEEDED) {
			return "SUCCEEDED"
		}
	}
	return "UNKNOWN"
}
