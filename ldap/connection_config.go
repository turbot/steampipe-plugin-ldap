package ldap

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/schema"
)

type ldapConfig struct {
	Attributes                     []string `cty:"attributes"`
	BaseDN                         *string  `cty:"base_dn"`
	Username                       *string  `cty:"username"`
	Password                       *string  `cty:"password"`
	Host                           *string  `cty:"host"`
	Port                           *string  `cty:"port"`
	TLSRequired                    *bool    `cty:"tls_required"`
	UserObjectFilter               *string  `cty:"user_object_filter"`
	GroupObjectFilter              *string  `cty:"group_object_filter"`
	OrganizationalUnitObjectFilter *string  `cty:"ou_object_filter"`
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
	"host": {
		Type: schema.TypeString,
	},
	"port": {
		Type: schema.TypeString,
	},
	"tls_required": {
		Type: schema.TypeBool,
	},
	"user_object_filter": {
		Type: schema.TypeString,
	},
	"group_object_filter": {
		Type: schema.TypeString,
	},
	"ou_object_filter": {
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
