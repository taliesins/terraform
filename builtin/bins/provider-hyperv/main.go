package main

import (
	"github.com/hashicorp/terraform/builtin/providers/hyperv"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: hyperv.Provider,
	})
}
