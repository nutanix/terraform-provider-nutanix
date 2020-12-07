package nutanix

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/karbon"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
	"gopkg.in/yaml.v2"
)

func dataSourceNutanixKarbonClusterKubeconfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNutanixKarbonClusterKubeconfigRead,

		Schema: map[string]*schema.Schema{
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
			"access_token": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"cluster_ca_certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNutanixKarbonClusterKubeconfigRead(d *schema.ResourceData, meta interface{}) error {
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
	var resp *karbon.KarbonClusterKubeconfig

	if iok {
		resp, err = GetKubeConfigForCluster(conn, karbonClusterID.(string))
	} else {
		resp, err = GetKubeConfigForCluster(conn, karbonClusterName.(string))
	}

	if err != nil {
		d.SetId("")
		return err
	}
	utils.PrintToJSON(resp, "resp: ")
	if len(resp.Clusters) != 1 {
		return fmt.Errorf("Incorrect amount of cluster information retrieved via Kubeconfig. Must be 1.")
	}
	if len(resp.Users) != 1 {
		return fmt.Errorf("Incorrect amount of user information retrieved via Kubeconfig. Must be 1.")
	}

	if err := d.Set("cluster_ca_certificate", resp.Clusters[0].Cluster.CertificateAuthorityData); err != nil {
		return fmt.Errorf("error setting `cluster_ca_certificate` for Karbon cluster (%s): %s", d.Id(), err)
	}
	if err := d.Set("cluster_url", resp.Clusters[0].Cluster.Server); err != nil {
		return fmt.Errorf("error setting `cluster_url` for Karbon cluster (%s): %s", d.Id(), err)
	}
	if err := d.Set("access_token", resp.Users[0].User.Token); err != nil {
		return fmt.Errorf("error setting `access_token` for Karbon cluster (%s): %s", d.Id(), err)
	}
	karbonClusterNameRetrieved := resp.Clusters[0].Name
	d.SetId(karbonClusterNameRetrieved)

	return nil
}

func GetKubeConfigForCluster(con *karbon.Client, karbonClusterName string) (*karbon.KarbonClusterKubeconfig, error) {
	kubeconfig, err := con.Cluster.GetKubeConfigForKarbonCluster(karbonClusterName)
	if err != nil {
		return nil, err
	}
	karbonClusterKubeconfig := karbon.KarbonClusterKubeconfig{}
	err = yaml.Unmarshal([]byte(kubeconfig.KubeConfig), &karbonClusterKubeconfig)
	if err != nil {
		return nil, err
	}
	utils.PrintToJSON(karbonClusterKubeconfig, "[karbonClusterKubeconfig]")
	return &karbonClusterKubeconfig, nil
}
