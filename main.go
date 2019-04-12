package main

import (
	"github.com/mavidser/terraform-provider-linux/linux"

	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return linux.Provider()
		},
	})
}
