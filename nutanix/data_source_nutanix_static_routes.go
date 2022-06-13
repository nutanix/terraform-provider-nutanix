package nutanix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNutanixStaticRoute() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixStaticRouteRead,
		Schema:      map[string]*schema.Schema{},
	}
}

func dataSourceNutanixStaticRouteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
