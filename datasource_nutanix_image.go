package nutanix

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNutanixImage() *schema.Resource {
	// TODO: return &schema.Resource for nutanix_image
	return &schema.Resource{
		Schema: map[string]*schema.Schema{},
	}
}
