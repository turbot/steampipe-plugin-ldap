package ldap

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type ldapConfig struct {
	Attributes                     []string `hcl:"attributes"`
	BaseDN                         *string  `hcl:"base_dn"`
	Username                       *string  `hcl:"username"`
	Password                       *string  `hcl:"password"`
	Host                           *string  `hcl:"host"`
	Port                           *string  `hcl:"port"`
	TLSRequired                    *bool    `hcl:"tls_required"`
	UserObjectFilter               *string  `hcl:"user_object_filter"`
	GroupObjectFilter              *string  `hcl:"group_object_filter"`
	OrganizationalUnitObjectFilter *string  `hcl:"ou_object_filter"`
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
