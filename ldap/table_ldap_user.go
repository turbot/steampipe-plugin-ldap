package ldap

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-ldap/ldap/v3"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

// TODO: Add missing LDAP config options
// TODO: Fix 'Error: LDAP Result Code 200 "Network Error": ldap: connection closed' after error queries
func tableLDAPUser(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ldap_user",
		Description: "LDAP users.",
		List: &plugin.ListConfig{
			Hydrate: listUsers,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "filter", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// TODO: How to avoid all transform.FromField calls?
			{
				Name:        "dn",
				Description: "The distinguished name (DN) for this resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("dn"),
			},
			{
				Name:        "base_dn",
				Description: "The base path to search in.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("base_dn"),
			},
			{
				Name:        "filter",
				Description: "The filter to search with.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("filter"),
			},
			{
				Name:        "cn",
				Description: "The user's common name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("cn"),
			},
			{
				Name:        "description",
				Description: "The user's description.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("description"),
			},
			{
				Name:        "display_name",
				Description: "The user's display name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("displayName"),
			},
			{
				Name:        "given_name",
				Description: "The user's given name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("givenName"),
			},
			{
				Name:        "initials",
				Description: "The user's initials.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("initials"),
			},
			{
				Name:        "mail",
				Description: "The user's email address.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("mail"),
			},
			{
				Name:        "object_class",
				Description: "The user's object classes.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("objectClass"),
			},
			{
				Name:        "ou",
				Description: "The user's organizational unit (OU).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ou"),
			},
			{
				Name:        "sn",
				Description: "The user's surname.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("sn"),
			},
			{
				Name:        "uid",
				Description: "The user's ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("uid"),
			},
			{
				Name:        "attributes",
				Description: "The attributes for this resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("attributes"),
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
				Transform:   transform.FromField("cn"),
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

	// TODO: Where to close connection?
	//defer conn.Close()

	var baseDN, filter string
	var attributes []string

	ldapConfig := GetConfig(d.Connection)
	if &ldapConfig != nil {
		if ldapConfig.BaseDN != nil {
			baseDN = *ldapConfig.BaseDN
		}
		if ldapConfig.Attributes != nil {
			attributes = ldapConfig.Attributes
		}
	}

	// Check for all required config args
	if baseDN == "" {
		return nil, errors.New("'base_dn' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}

	keyQuals := d.KeyColumnQuals

	// Filters must start and finish with ()!
	// TODO: Construct filters based on passed in quals
	if keyQuals["filter"] != nil {
		filter = keyQuals["filter"].GetStringValue()
	} else {
		// TODO: Why doesn't objectCategory work, data?
		//filter = fmt.Sprintf("(&(objectClass=person)(objectCategory=person))")
		filter = fmt.Sprintf("(&(objectClass=person))")
	}

	// TODO: Do we need to escape what the users pass in?
	//filter = ldap.EscapeFilter(filter)

	logger.Warn("baseDN", baseDN)
	logger.Warn("filter", filter)
	logger.Warn("attributes", attributes)

	var searchReq *ldap.SearchRequest
	// If no attributes are passed in, search request will get all of them
	if attributes != nil {
		searchReq = ldap.NewSearchRequest(baseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, attributes, []ldap.Control{})
	} else {
		searchReq = ldap.NewSearchRequest(baseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, []string{}, []ldap.Control{})
	}

	result, err := conn.Search(searchReq)
	if err != nil {
		plugin.Logger(ctx).Error("ldap_user.listUsers", "search_error", err)
		return nil, err
	}

	// TODO: Add standardizing for output casing
	for _, entry := range result.Entries {
		row := make(map[string]interface{})
		for _, attr := range entry.Attributes {
			// TODO: Handle null char \u0000 better to avoid 'Error: unsupported Unicode escape sequence'
			if attr.Name != "jpegPhoto" {
				// TODO: Better handle single/multiple values
				if len(attr.Values) == 1 {
					row[attr.Name] = entry.GetAttributeValue(attr.Name)
				} else if len(attr.Values) > 1 {
					row[attr.Name] = entry.GetAttributeValues(attr.Name)
				}
			}
		}
		row["base_dn"] = baseDN
		row["dn"] = entry.DN
		row["attributes"] = entry.Attributes
		d.StreamListItem(ctx, row)
	}

	return nil, nil
}
