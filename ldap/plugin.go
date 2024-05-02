package ldap

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             "steampipe-plugin-ldap",
		DefaultTransform: transform.FromGo(),
		ConnectionKeyColumns: []plugin.ConnectionKeyColumn{
			{
				Name:    "host_name",
				Hydrate: getHostName,
			},
		},
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
		},
		TableMap: map[string]*plugin.Table{
			"ldap_group":               tableLDAPGroup(ctx),
			"ldap_organizational_unit": tableLDAPOrganizationalUnit(ctx),
			"ldap_user":                tableLDAPUser(ctx),
		},
	}
	return p
}
