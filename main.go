// Copyright 2020 BlueCat Networks. All rights reserved

package main

import (
	"terraform-provider-bluecat-gss/bluecat"

	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return bluecat.Provider()
		},
	})
}
