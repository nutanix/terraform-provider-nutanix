package main

import (
	"github.com/hashicorp/terraform/plugin"
	"terraform-provider-nutanix"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: nutanix.Provider,
	})
}
