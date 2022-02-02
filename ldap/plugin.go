package ldap

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             "steampipe-plugin-ldap",
		DefaultTransform: transform.FromGo(),
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		TableMap: map[string]*plugin.Table{
			"ldap_group":               tableLDAPGroup(ctx),
			"ldap_organizational_unit": tableLDAPOrganizationalUnit(ctx),
			"ldap_user":                tableLDAPUser(ctx),
		},
	}
	return p
}
