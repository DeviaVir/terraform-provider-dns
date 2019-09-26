package main

import (
	"github.com/DeviaVir/terraform-provider-dns/dns"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: dns.Provider})
}
