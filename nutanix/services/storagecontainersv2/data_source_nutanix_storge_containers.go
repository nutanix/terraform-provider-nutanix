package storagecontainersv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clustermgmt "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	clsConfig "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/common/v1/config"
	clsResponse "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/common/v1/response"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixStorageContainersV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixStorageContainersV2Read,
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
			"select": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"storage_containers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixStorageContainerV2(),
			},
		},
	}
}

func DatasourceNutanixStorageContainersV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI

	// initialize query params
	var filter, orderBy, selectQ *string
	var page, limit *int

	if pagef, ok := d.GetOk("page"); ok {
		page = utils.IntPtr(pagef.(int))
	} else {
		page = nil
	}
	if limitf, ok := d.GetOk("limit"); ok {
		limit = utils.IntPtr(limitf.(int))
	} else {
		limit = nil
	}
	if filterf, ok := d.GetOk("filter"); ok {
		filter = utils.StringPtr(filterf.(string))
	} else {
		filter = nil
	}
	if order, ok := d.GetOk("order_by"); ok {
		orderBy = utils.StringPtr(order.(string))
	} else {
		orderBy = nil
	}
	if selectQy, ok := d.GetOk("apply"); ok {
		selectQ = utils.StringPtr(selectQy.(string))
	} else {
		selectQ = nil
	}

	resp, err := conn.StorageContainersAPI.ListStorageContainers(page, limit, filter, orderBy, selectQ)
	if err != nil {
		return diag.Errorf("error while fetching Storage Containers : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("storage_containers", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of storage containers.",
		}}
	}

	getResp := resp.Data.GetValue().([]clustermgmt.StorageContainer)

	if err := d.Set("storage_containers", flattenStorageContainers(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenStorageContainers(storageContainers []clustermgmt.StorageContainer) []interface{} {
	if len(storageContainers) > 0 {
		storageContainersList := make([]interface{}, len(storageContainers))

		for k, v := range storageContainers {
			storageContainer := make(map[string]interface{})

			storageContainer["ext_id"] = v.ContainerExtId
			storageContainer["tenant_id"] = v.TenantId
			storageContainer["links"] = flattenLinks(v.Links)
			storageContainer["container_ext_id"] = v.ContainerExtId
			storageContainer["owner_ext_id"] = v.OwnerExtId
			storageContainer["name"] = v.Name
			storageContainer["cluster_ext_id"] = v.ClusterExtId
			storageContainer["storage_pool_ext_id"] = v.StoragePoolExtId
			storageContainer["is_marked_for_removal"] = v.IsMarkedForRemoval
			// storageContainer["marked_for_removal"] = v.M
			storageContainer["max_capacity_bytes"] = v.MaxCapacityBytes
			storageContainer["logical_explicit_reserved_capacity_bytes"] = v.LogicalExplicitReservedCapacityBytes
			storageContainer["logical_implicit_reserved_capacity_bytes"] = v.LogicalImplicitReservedCapacityBytes
			storageContainer["logical_advertised_capacity_bytes"] = v.LogicalAdvertisedCapacityBytes
			storageContainer["replication_factor"] = v.ReplicationFactor
			storageContainer["nfs_whitelist_addresses"] = flattenNfsWhitelistAddresses(v.NfsWhitelistAddress)
			storageContainer["is_nfs_whitelist_inherited"] = v.IsNfsWhitelistInherited
			storageContainer["erasure_code"] = flattenErasureCodeStatus(v.ErasureCode)
			storageContainer["is_inline_ec_enabled"] = v.IsInlineEcEnabled
			storageContainer["has_higher_ec_fault_domain_preference"] = v.HasHigherEcFaultDomainPreference
			storageContainer["erasure_code_delay_secs"] = v.ErasureCodeDelaySecs
			storageContainer["cache_deduplication"] = flattenCacheDeduplication(v.CacheDeduplication)
			storageContainer["on_disk_dedup"] = flattenOnDiskDedup(v.OnDiskDedup)
			storageContainer["is_compression_enabled"] = v.IsCompressionEnabled
			storageContainer["compression_delay_secs"] = v.CompressionDelaySecs
			storageContainer["is_internal"] = v.IsInternal
			storageContainer["is_software_encryption_enabled"] = v.IsSoftwareEncryptionEnabled
			storageContainer["is_encrypted"] = v.IsEncrypted
			storageContainer["cluster_name"] = v.ClusterName

			storageContainersList[k] = storageContainer
		}
		return storageContainersList
	}
	return nil
}

func flattenNfsWhitelistAddresses(pr []clsConfig.IPAddressOrFQDN) []map[string]interface{} {
	if len(pr) > 0 {
		ips := make([]map[string]interface{}, len(pr))

		for k, v := range pr {
			ip := make(map[string]interface{})

			if v.Ipv4 != nil {
				ip["ipv4"] = flattenIPv4Address(v.Ipv4)
			}
			if v.Ipv6 != nil {
				ip["ipv6"] = flattenIPv6Address(v.Ipv6)
			}
			if v.Fqdn != nil {
				ip["fqdn"] = flattenFQDN(v.Fqdn)
			}
			ips[k] = ip
		}
		return ips
	}
	return nil
}

func flattenCacheDeduplication(pr *clustermgmt.CacheDeduplication) string {
	if pr != nil {
		const one, two, three, four = 1, 2, 3, 4
		if *pr == clustermgmt.CacheDeduplication(one) {
			return "REDACTED"
		}
		if *pr == clustermgmt.CacheDeduplication(two) {
			return "NONE"
		}
		if *pr == clustermgmt.CacheDeduplication(three) {
			return "OFF"
		}
		if *pr == clustermgmt.CacheDeduplication(four) {
			return "ON"
		}
	}
	return "UNKNOWN"
}

func flattenErasureCodeStatus(pr *clustermgmt.ErasureCodeStatus) string {
	if pr != nil {
		const one, two, three, four = 1, 2, 3, 4
		if *pr == clustermgmt.ErasureCodeStatus(one) {
			return "REDACTED"
		}
		if *pr == clustermgmt.ErasureCodeStatus(two) {
			return "NONE"
		}
		if *pr == clustermgmt.ErasureCodeStatus(three) {
			return "OFF"
		}
		if *pr == clustermgmt.ErasureCodeStatus(four) {
			return "ON"
		}
	}
	return "UNKNOWN"
}

func flattenOnDiskDedup(pr *clustermgmt.OnDiskDedup) string {
	if pr != nil {
		const one, two, three, four = 1, 2, 3, 4
		if *pr == clustermgmt.OnDiskDedup(one) {
			return "REDACTED"
		}
		if *pr == clustermgmt.OnDiskDedup(two) {
			return "NONE"
		}
		if *pr == clustermgmt.OnDiskDedup(three) {
			return "OFF"
		}
		if *pr == clustermgmt.OnDiskDedup(four) {
			return "POST_PROCESS"
		}
	}
	return "UNKNOWN"
}

func flattenLinks(pr []clsResponse.ApiLink) []map[string]interface{} {
	if len(pr) > 0 {
		linkList := make([]map[string]interface{}, len(pr))

		for k, v := range pr {
			links := map[string]interface{}{}
			if v.Href != nil {
				links["href"] = v.Href
			}
			if v.Rel != nil {
				links["rel"] = v.Rel
			}

			linkList[k] = links
		}
		return linkList
	}
	return nil
}

func flattenIPv4Address(pr *clsConfig.IPv4Address) []interface{} {
	if pr != nil {
		ipv4 := make([]interface{}, 0)

		ip := make(map[string]interface{})

		ip["value"] = pr.Value
		ip["prefix_length"] = pr.PrefixLength

		ipv4 = append(ipv4, ip)

		return ipv4
	}
	return nil
}

func flattenIPv6Address(pr *clsConfig.IPv6Address) []interface{} {
	if pr != nil {
		ipv6 := make([]interface{}, 0)

		ip := make(map[string]interface{})

		ip["value"] = pr.Value
		ip["prefix_length"] = pr.PrefixLength

		ipv6 = append(ipv6, ip)

		return ipv6
	}
	return nil
}

func flattenFQDN(pr *clsConfig.FQDN) []interface{} {
	if pr != nil {
		fqdn := make([]interface{}, 0)

		f := make(map[string]interface{})

		f["value"] = pr.Value

		fqdn = append(fqdn, f)

		return fqdn
	}
	return nil
}
