package clustersv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/clustermgmt-go-client/v17/models/clustermgmt/v4/config"
	import2 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/clustermgmt-go-client/v17/models/clustermgmt/v4/request/clusters"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixClusterEntitiesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixClusterEntitiesV2Read,
		Schema: map[string]*schema.Schema{
			"page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"apply": {
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
			"cluster_entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixClusterEntityV2(),
			},
		},
	}
}

func DatasourceNutanixClusterEntitiesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI

	listClustersRequest := import2.ListClustersRequest{}

	if v, ok := d.GetOk("page"); ok {
		listClustersRequest.Page_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("limit"); ok {
		listClustersRequest.Limit_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("filter"); ok {
		listClustersRequest.Filter_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("order_by"); ok {
		listClustersRequest.Orderby_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("apply"); ok {
		listClustersRequest.Apply_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("expand"); ok {
		listClustersRequest.Expand_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("select"); ok {
		listClustersRequest.Select_ = utils.StringPtr(v.(string))
	}

	resp, err := conn.ClusterEntityAPI.ListClusters(ctx, &listClustersRequest)
	if err != nil {
		return diag.Errorf("error while fetching cluster entities : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("cluster_entities", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(resource.UniqueId())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No Clusters found",
			Detail:   "The API returned an empty list of clusters.",
		}}
	}
	getResp := resp.Data.GetValue().([]import1.Cluster)

	if err := d.Set("cluster_entities", flattenClusterEntities(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenClusterEntities(pr []import1.Cluster) []interface{} {
	if len(pr) > 0 {
		clsList := make([]interface{}, len(pr))

		for k, v := range pr {
			cls := make(map[string]interface{})

			cls["ext_id"] = v.ExtId
			cls["tenant_id"] = v.TenantId
			cls["links"] = common.FlattenLinks(v.Links)
			cls["name"] = v.Name
			cls["nodes"] = flattenNodeReference(v.Nodes)
			cls["network"] = flattenClusterNetworkReference(v.Network)
			cls["config"] = flattenClusterConfigReference(v.Config)
			cls["upgrade_status"] = flattenUpgradeStatus(v.UpgradeStatus)
			cls["vm_count"] = v.VmCount
			cls["inefficient_vm_count"] = v.InefficientVmCount
			cls["container_name"] = v.ContainerName
			cls["categories"] = v.Categories
			cls["cluster_profile_ext_id"] = v.ClusterProfileExtId
			cls["backup_eligibility_score"] = v.BackupEligibilityScore

			clsList[k] = cls
		}
		return clsList
	}
	return nil
}
