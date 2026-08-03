package objectstoresv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	objectsCommon "github.com/nutanix-core/ntnx-api-golang-sdk-internal/objects-go-client/v17/models/common/v1/config"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/objects-go-client/v17/models/objects/v4/config"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/objects-go-client/v17/models/objects/v4/request/objectstores"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixObjectStoresV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixObjectStoresV2Read,
		Schema: map[string]*schema.Schema{
			"page": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0, //nolint:gomnd
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  50, //nolint:gomnd
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"expand": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"select": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"object_stores": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixObjectStoreV2(),
			},
		},
	}
}

func DatasourceNutanixObjectStoresV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ObjectStoreAPI

	listObjectstoresRequest := import1.ListObjectstoresRequest{}

	if common.IsExplicitlySet(d, "page") {
		listObjectstoresRequest.Page_ = utils.IntPtr(d.Get("page").(int))
	} else {
		listObjectstoresRequest.Page_ = utils.IntPtr(0)
	}
	if v, ok := d.GetOk("limit"); ok {
		listObjectstoresRequest.Limit_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("filter"); ok {
		listObjectstoresRequest.Filter_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("order_by"); ok {
		listObjectstoresRequest.Orderby_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("expand"); ok {
		listObjectstoresRequest.Expand_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("select"); ok {
		listObjectstoresRequest.Select_ = utils.StringPtr(v.(string))
	}

	resp, err := conn.ObjectStoresAPIInstance.ListObjectstores(ctx, &listObjectstoresRequest)
	if err != nil {
		return diag.Errorf("error while fetching object stores : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("object_stores", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No Objects store found",
			Detail:   "The API returned an empty list of Objects stores.",
		}}
	}

	objectStoreList := resp.Data.GetValue().([]config.ObjectStore)

	if err := d.Set("object_stores", flattenObjectStoreEntities(objectStoreList)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.GenUUID())

	return nil
}

func flattenObjectStoreEntities(objectStoresList []config.ObjectStore) []map[string]interface{} {
	if len(objectStoresList) == 0 {
		return []map[string]interface{}{}
	}

	result := make([]map[string]interface{}, 0)

	for _, objectStore := range objectStoresList {
		objectStoreMap := map[string]interface{}{
			"ext_id":                    objectStore.ExtId,
			"tenant_id":                 objectStore.TenantId,
			"links":                     flattenLinks(objectStore.Links),
			"metadata":                  flattenMetadata(objectStore.Metadata),
			"name":                      utils.StringValue(objectStore.Name),
			"creation_time":             flattenTime(objectStore.CreationTime),
			"last_update_time":          flattenTime(objectStore.LastUpdateTime),
			"description":               utils.StringValue(objectStore.Description),
			"deployment_version":        utils.StringValue(objectStore.DeploymentVersion),
			"domain":                    utils.StringValue(objectStore.Domain),
			"region":                    utils.StringValue(objectStore.Region),
			"num_worker_nodes":          utils.Int64Value(objectStore.NumWorkerNodes),
			"cluster_ext_id":            utils.StringValue(objectStore.ClusterExtId),
			"storage_network_reference": utils.StringValue(objectStore.StorageNetworkReference),
			"storage_network_vip":       flattenIPAddress([]objectsCommon.IPAddress{*objectStore.StorageNetworkVip}),
			"storage_network_dns_ip":    flattenIPAddress([]objectsCommon.IPAddress{*objectStore.StorageNetworkDnsIp}),
			"public_network_reference":  utils.StringValue(objectStore.PublicNetworkReference),
			"public_network_ips":        flattenIPAddress(objectStore.PublicNetworkIps),
			"total_capacity_gib":        utils.Int64Value(objectStore.TotalCapacityGiB),
			"state":                     objectStore.State.GetName(),
			"certificate_ext_ids":       objectStore.CertificateExtIds,
		}
		result = append(result, objectStoreMap)
	}
	return result
}
