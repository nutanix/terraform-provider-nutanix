package objectstoresv2

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	objectsCommon "github.com/nutanix/ntnx-api-golang-clients/objects-go-client/v4/models/common/v1/config"
	objectsResponse "github.com/nutanix/ntnx-api-golang-clients/objects-go-client/v4/models/common/v1/response"
	"github.com/nutanix/ntnx-api-golang-clients/objects-go-client/v4/models/objects/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixObjectStoreV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixObjectStoreV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
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
			},
			"metadata": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: metadataSchema(),
				},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deployment_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"num_worker_nodes": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"cluster_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_network_reference": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_network_vip": {
				Type:     schema.TypeList,
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
				Computed: true,
			},
			"public_network_ips": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4": SchemaForValuePrefixLength(ipv4PrefixLengthDefaultValue),
						"ipv6": SchemaForValuePrefixLength(ipv6PrefixLengthDefaultValue),
					},
				},
			},
			"total_capacity_gib": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_ext_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func DatasourceNutanixObjectStoreV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ObjectStoreAPI

	objectStoreExtID := d.Get("ext_id").(string)

	readResp, err := conn.ObjectStoresAPIInstance.GetObjectstoreById(utils.StringPtr(objectStoreExtID))
	if err != nil || readResp.Data == nil {
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

	d.SetId(utils.StringValue(objectStore.ExtId))

	return nil
}

// flatteners
func flattenMetadata(metadata *objectsCommon.Metadata) []map[string]interface{} {
	if metadata != nil {
		metadataMapList := make([]map[string]interface{}, 0)

		metadataMap := make(map[string]interface{})

		metadataMap["owner_reference_id"] = metadata.OwnerReferenceId
		metadataMap["owner_user_name"] = metadata.OwnerUserName
		metadataMap["project_reference_id"] = metadata.ProjectReferenceId
		metadataMap["project_name"] = metadata.ProjectName
		metadataMap["category_ids"] = metadata.CategoryIds

		metadataMapList = append(metadataMapList, metadataMap)
		return metadataMapList
	}
	return nil
}

func flattenLinks(links []objectsResponse.ApiLink) []map[string]interface{} {
	if len(links) > 0 {
		linksMapList := make([]map[string]interface{}, len(links))

		for k, v := range links {
			linkMap := map[string]interface{}{}
			if v.Href != nil {
				linkMap["href"] = v.Href
			}
			if v.Rel != nil {
				linkMap["rel"] = v.Rel
			}

			linksMapList[k] = linkMap
		}
		return linksMapList
	}
	return nil
}

func flattenTime(inTime *time.Time) string {
	if inTime != nil {
		return inTime.UTC().Format(time.RFC3339)
	}
	return ""
}

func flattenIPAddress(ipAddresses []objectsCommon.IPAddress) []map[string]interface{} {
	if ipAddresses != nil || len(ipAddresses) > 0 {
		ipAddressesMapList := make([]map[string]interface{}, 0)

		for _, ipAddress := range ipAddresses {
			ipAddressMap := make(map[string]interface{})

			ipAddressMap["ipv4"] = flattenFloatingIPv4Address(ipAddress.Ipv4)
			ipAddressMap["ipv6"] = flattenFloatingIPv6Address(ipAddress.Ipv6)

			ipAddressesMapList = append(ipAddressesMapList, ipAddressMap)
		}

		return ipAddressesMapList
	}
	return nil
}

func flattenFloatingIPv4Address(pr *objectsCommon.IPv4Address) []map[string]interface{} {
	if pr != nil {
		ips := make([]map[string]interface{}, 0)

		ip := make(map[string]interface{})

		ip["prefix_length"] = pr.PrefixLength
		ip["value"] = pr.Value

		ips = append(ips, ip)

		return ips
	}
	return nil
}

func flattenFloatingIPv6Address(pr *objectsCommon.IPv6Address) []map[string]interface{} {
	if pr != nil {
		ips := make([]map[string]interface{}, 0)

		ip := make(map[string]interface{})

		ip["prefix_length"] = pr.PrefixLength
		ip["value"] = pr.Value

		ips = append(ips, ip)

		return ips
	}
	return nil
}
