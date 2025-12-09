package storagecontainersv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clustermgmt "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixStorageContainerV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixStorageContainerV2Read,
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
			"container_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_pool_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_marked_for_removal": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"max_capacity_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"logical_explicit_reserved_capacity_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"logical_implicit_reserved_capacity_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"logical_advertised_capacity_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"replication_factor": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"nfs_whitelist_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4": SchemaForValuePrefixLength(),
						"ipv6": SchemaForValuePrefixLength(),
						"fqdn": SchemaForFqdnValue(),
					},
				},
			},
			"is_nfs_whitelist_inherited": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"erasure_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_inline_ec_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"has_higher_ec_fault_domain_preference": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"erasure_code_delay_secs": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cache_deduplication": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"on_disk_dedup": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_compression_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"compression_delay_secs": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"is_internal": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_software_encryption_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_encrypted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"affinity_host_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DatasourceNutanixStorageContainerV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI

	extID := d.Get("ext_id")
	resp, err := conn.StorageContainersAPI.GetStorageContainerById(utils.StringPtr(extID.(string)))
	if err != nil {
		return diag.Errorf("error while fetching Storage Container : %v", err)
	}

	getResp := resp.Data.GetValue().(clustermgmt.StorageContainer)

	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("container_ext_id", getResp.ContainerExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("owner_ext_id", getResp.OwnerExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_ext_id", getResp.ClusterExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("storage_pool_ext_id", getResp.StoragePoolExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_marked_for_removal", getResp.IsMarkedForRemoval); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("max_capacity_bytes", getResp.MaxCapacityBytes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("logical_explicit_reserved_capacity_bytes", getResp.LogicalExplicitReservedCapacityBytes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("logical_implicit_reserved_capacity_bytes", getResp.LogicalImplicitReservedCapacityBytes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("logical_advertised_capacity_bytes", getResp.LogicalAdvertisedCapacityBytes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("replication_factor", getResp.ReplicationFactor); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("nfs_whitelist_addresses", flattenNfsWhitelistAddresses(getResp.NfsWhitelistAddress)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_nfs_whitelist_inherited", getResp.IsNfsWhitelistInherited); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("erasure_code", flattenErasureCodeStatus(getResp.ErasureCode)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_inline_ec_enabled", getResp.IsInlineEcEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("has_higher_ec_fault_domain_preference", getResp.HasHigherEcFaultDomainPreference); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("erasure_code_delay_secs", getResp.ErasureCodeDelaySecs); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cache_deduplication", flattenCacheDeduplication(getResp.CacheDeduplication)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("on_disk_dedup", flattenOnDiskDedup(getResp.OnDiskDedup)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_compression_enabled", getResp.IsCompressionEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("compression_delay_secs", getResp.CompressionDelaySecs); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_internal", getResp.IsInternal); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_software_encryption_enabled", getResp.IsSoftwareEncryptionEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_encrypted", getResp.IsEncrypted); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("affinity_host_ext_id", getResp.AffinityHostExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_name", getResp.ClusterName); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(getResp.ContainerExtId))
	return nil
}

func SchemaForFqdnValue() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"value": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func SchemaForValuePrefixLength() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"value": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"prefix_length": {
					Type:     schema.TypeInt,
					Computed: true,
				},
			},
		},
	}
}
