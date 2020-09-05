package nutanix

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/karbon"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func dataSourceNutanixKarbonClusterSSH() *schema.Resource {
	return &schema.Resource{
		Read:          dataSourceNutanixKarbonClusterSSHRead,
		SchemaVersion: 1,
		Schema:        KarbonClusterSSHConfigElementDataSourceMap(),
	}
}

func dataSourceNutanixKarbonClusterSSHRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*Client).KarbonAPI
	setTimeout(meta)
	// Make request to the API
	karbonClusterID, iok := d.GetOk("karbon_cluster_id")
	karbonClusterName, nok := d.GetOk("karbon_cluster_name")
	if !iok && !nok {
		return fmt.Errorf("please provide one of karbon_cluster_id or karbon_cluster_name attributes")
	}
	var err error
	var resp *karbon.KarbonClusterSSHconfig

	if iok {
		c, err := conn.Cluster.GetKarbonCluster(karbonClusterID.(string))
		if err != nil {
			return fmt.Errorf("Unable to find cluster with id %s: %s", karbonClusterID, err)
		}
		resp, err = conn.Cluster.GetSSHConfigForKarbonCluster(*c.Name)
	} else {
		resp, err = conn.Cluster.GetSSHConfigForKarbonCluster(karbonClusterName.(string))
	}
	utils.PrintToJSON(resp, "resp: ")
	if err != nil {
		d.SetId("")
		return err
	}

	if err := d.Set("certificate", resp.Certificate); err != nil {
		return fmt.Errorf("Failed to set certificate output: %s", err)
	}
	if err := d.Set("expiry_time", resp.ExpiryTime); err != nil {
		return fmt.Errorf("Failed to set expiry_time output: %s", err)
	}
	if err := d.Set("private_key", resp.PrivateKey); err != nil {
		return fmt.Errorf("Failed to set private_key output: %s", err)
	}
	if err := d.Set("username", resp.Username); err != nil {
		return fmt.Errorf("Failed to set username output: %s", err)
	}
	d.SetId(resource.UniqueId())

	return nil
}

func KarbonClusterSSHConfigElementDataSourceMap() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"karbon_cluster_id": &schema.Schema{
			Type:          schema.TypeString,
			Optional:      true,
			ConflictsWith: []string{"karbon_cluster_name"},
		},
		"karbon_cluster_name": &schema.Schema{
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
