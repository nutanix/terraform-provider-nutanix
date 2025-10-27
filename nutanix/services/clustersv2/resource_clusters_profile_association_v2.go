package clustersv2

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func ResourceNutanixClustersProfileAssociationV2() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"profile_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}
