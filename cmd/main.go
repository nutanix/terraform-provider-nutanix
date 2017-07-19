package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/ideadevice/terraform-ahv-provider-plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: nutanix.Provider,
	})
}
