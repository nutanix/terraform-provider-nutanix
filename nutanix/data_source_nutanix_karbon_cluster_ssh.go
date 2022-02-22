package nutanix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/karbon"
)

func dataSourceNutanixKarbonClusterSSH() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceNutanixKarbonClusterSSHRead,
		SchemaVersion: 1,
		Schema:        KarbonClusterSSHConfigElementDataSourceMap(),
	}
}

func dataSourceNutanixKarbonClusterSSHRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection
	conn := meta.(*Client).KarbonAPI
	setTimeout(meta)
	// Make request to the API
	karbonClusterID, iok := d.GetOk("karbon_cluster_id")
	karbonClusterNameInput, nok := d.GetOk("karbon_cluster_name")
	if !iok && !nok {
		return diag.Errorf("please provide one of karbon_cluster_id or karbon_cluster_name attributes")
	}
	var err error
	var resp *karbon.ClusterSSHconfig
	var karbonClusterName string
	if iok {
		var c *karbon.ClusterIntentResponse
		c, err = conn.Cluster.GetKarbonCluster(karbonClusterID.(string))
		if err != nil {
			return diag.Errorf("unable to find cluster with id %s: %s", karbonClusterID, err)
		}
		karbonClusterName = *c.Name
	} else {
		karbonClusterName = karbonClusterNameInput.(string)
	}

	resp, err = conn.Cluster.GetSSHConfigForKarbonCluster(karbonClusterName)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	if err := d.Set("certificate", resp.Certificate); err != nil {
		return diag.Errorf("failed to set certificate output: %s", err)
	}
	if err := d.Set("expiry_time", resp.ExpiryTime); err != nil {
		return diag.Errorf("failed to set expiry_time output: %s", err)
	}
	if err := d.Set("private_key", resp.PrivateKey); err != nil {
		return diag.Errorf("failed to set private_key output: %s", err)
	}
	if err := d.Set("username", resp.Username); err != nil {
		return diag.Errorf("failed to set username output: %s", err)
	}
	d.SetId(resource.UniqueId())

	return nil
}

func KarbonClusterSSHConfigElementDataSourceMap() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"karbon_cluster_id": {
			Type:          schema.TypeString,
			Optional:      true,
			ConflictsWith: []string{"karbon_cluster_name"},
		},
		"karbon_cluster_name": {
			Type:          schema.TypeString,
			Optional:      true,
			ConflictsWith: []string{"karbon_cluster_id"},
		},
		"certificate": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"expiry_time": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"private_key": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"username": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
