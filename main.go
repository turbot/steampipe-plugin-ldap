package main

import (
	"github.com/turbot/steampipe-plugin-ldap/ldap"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		PluginFunc: ldap.Plugin})
}
