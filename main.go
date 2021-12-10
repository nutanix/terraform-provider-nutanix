package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/nutanix/terraform-provider-nutanix/nutanix"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: nutanix.Provider,
	})
}
