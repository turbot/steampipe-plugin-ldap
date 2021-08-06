package ldap

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-ldap/ldap"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableLDAPUser(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ldap_user",
		Description: "LDAP users.",
		List: &plugin.ListConfig{
			Hydrate: listUsers,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "path", Require: plugin.Optional},
				{Name: "filter", Require: plugin.Optional},
				{Name: "scope", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "dn",
				Description: "The distinguished name (DN) for this resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DN"),
			},
			{
				Name:        "path",
				Description: "The path to search in.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("path"),
			},
			{
				Name:        "filter",
				Description: "The filter to search with.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("filter"),
			},
			{
				Name:        "scope",
				Description: "The scope to search in.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("scope"),
			},
			{
				Name:        "attributes",
				Description: "The attributes for this resource.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "raw",
				Description: "The attributes for this resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromValue(),
			},

			// Standard columns
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
		},
	}
}

func listUsers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listUsers")

	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ldap_user.listUsers", "connection_error", err)
		return nil, err
	}

	defer conn.Close()

	var baseDN, filter, scope string

	keyQuals := d.KeyColumnQuals
	if keyQuals["path"] != nil {
		baseDN = keyQuals["path"].GetStringValue()
	} else {
		baseDN = "OU=people,DC=planetexpress,DC=com"
	}

	// Filters must start and finish with ()!
	if keyQuals["filter"] != nil {
		filter = keyQuals["filter"].GetStringValue()
	} else {
		filter = fmt.Sprintf("(&(objectClass=person))")
	}

	//filter = ldap.EscapeFilter(filter)

	if keyQuals["scope"] != nil {
		scope = keyQuals["scope"].GetStringValue()
	} else {
		scope = "sub"
	}

	var scopeInt int
	switch scope {
	case "base":
		scopeInt = 0
	case "single":
		scopeInt = 1
	case "sub":
		scopeInt = 2
	default:
		scopeErr := errors.New("Scope must be base, single, or sub")
		plugin.Logger(ctx).Error("ldap_user.listUsers", "scope_error", scopeErr)
		return nil, scopeErr
	}

	logger.Warn("baseDN", baseDN)
	logger.Warn("filter", filter)
	logger.Warn("scope", scopeInt)

	//searchReq := ldap.NewSearchRequest(baseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, []string{"displayName", "description"}, []ldap.Control{})
	searchReq := ldap.NewSearchRequest(baseDN, scopeInt, 0, 0, 0, false, filter, []string{"cn", "displayName", "description", "dn", "objectClass", "employeeType"}, []ldap.Control{})

	result, err := conn.Search(searchReq)
	if err != nil {
		plugin.Logger(ctx).Error("ldap_user.listUsers", "search_error", err)
		return nil, err
	}

	for _, entry := range result.Entries {
		d.StreamListItem(ctx, entry)
	}

	return nil, nil
}
