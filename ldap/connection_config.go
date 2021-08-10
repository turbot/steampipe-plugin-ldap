package ldap

import (
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/schema"
)

type ldapConfig struct {
	Attributes []string `cty:"attributes"`
	BaseDN     *string  `cty:"base_dn"`
	Username   *string  `cty:"username"`
	Password   *string  `cty:"password"`
	URL        *string  `cty:"url"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"attributes": {
		Type: schema.TypeList,
		Elem: &schema.Attribute{Type: schema.TypeString},
	},
	"base_dn": {
		Type: schema.TypeString,
	},
	"username": {
		Type: schema.TypeString,
	},
	"password": {
		Type: schema.TypeString,
	},
	"url": {
		Type: schema.TypeString,
	},
}

func ConfigInstance() interface{} {
	return &ldapConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) ldapConfig {
	if connection == nil || connection.Config == nil {
		return ldapConfig{}
	}
	config, _ := connection.Config.(ldapConfig)
	return config
}
