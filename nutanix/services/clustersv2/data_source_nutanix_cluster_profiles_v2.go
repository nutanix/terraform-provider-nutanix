package clustersv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/clustermgmt-go-client/v17/models/clustermgmt/v4/config"
	import2 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/clustermgmt-go-client/v17/models/clustermgmt/v4/request/clusterprofiles"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixClusterProfilesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixClusterProfilesV2Read,
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
			"select": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_profiles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixClusterProfileV2(),
			},
		},
	}
}

func DatasourceNutanixClusterProfilesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI

	listClusterProfilesRequest := import2.ListClusterProfilesRequest{}

	if v, ok := d.GetOk("page"); ok {
		listClusterProfilesRequest.Page_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("limit"); ok {
		listClusterProfilesRequest.Limit_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("filter"); ok {
		listClusterProfilesRequest.Filter_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("order_by"); ok {
		listClusterProfilesRequest.Orderby_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("select"); ok {
		listClusterProfilesRequest.Select_ = utils.StringPtr(v.(string))
	}

	resp, err := conn.ClusterProfilesAPI.ListClusterProfiles(ctx, &listClusterProfilesRequest)
	if err != nil {
		return diag.Errorf("error while fetching cluster profiles : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("cluster_profiles", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(resource.UniqueId())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No Cluster Profiles found",
			Detail:   "The API returned an empty list of cluster profiles.",
		}}
	}

	clusterProfiles := resp.Data.GetValue().([]import1.ClusterProfile)

	if err := d.Set("cluster_profiles", flattenClusterProfiles(clusterProfiles)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenClusterProfiles(clusterProfiles []import1.ClusterProfile) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(clusterProfiles))

	for _, clusterProfile := range clusterProfiles {
		clusterMap := map[string]interface{}{
			"tenant_id":             clusterProfile.TenantId,
			"links":                 common.FlattenLinks(clusterProfile.Links),
			"ext_id":                utils.StringValue(clusterProfile.ExtId),
			"name":                  utils.StringValue(clusterProfile.Name),
			"description":           utils.StringValue(clusterProfile.Description),
			"allowed_overrides":     common.FlattenEnumValueList(clusterProfile.AllowedOverrides),
			"name_server_ip_list":   flattenIPAddressList(clusterProfile.NameServerIpList),
			"ntp_server_ip_list":    flattenIPAddressOrFQDN(clusterProfile.NtpServerIpList),
			"smtp_server":           flattenSMTPServerRef(clusterProfile.SmtpServer),
			"nfs_subnet_white_list": clusterProfile.NfsSubnetWhitelist,
			"snmp_config":           flattenSnmpConfig(clusterProfile.SnmpConfig),
			"rsyslog_server_list":   flattenRsyslogServerList(clusterProfile.RsyslogServerList),
			"pulse_status":          flattenPulseStatus(clusterProfile.PulseStatus),
		}

		result = append(result, clusterMap)
	}

	return result
}
