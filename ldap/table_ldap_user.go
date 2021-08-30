package ldap

import (
	"context"
	"errors"

	"github.com/go-ldap/ldap/v3"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

// TODO: Test this table with 100, 500, and 1000+ users
// TODO: Test this table against an AD server

// TODO: Add all columns here to allow for proper hydration
type userRow struct {
	Dn          string
	BaseDn      string
	Filter      string
	Cn          string
	Description string
	DisplayName string
	GivenName   string
	Initials    string
	Mail        string
	ObjectClass []string
	Ou          string
	Sn          string
	Uid         string
	Attributes  []*ldap.EntryAttribute
	Raw         []string
}

// TODO: Add missing LDAP config options
func tableLDAPUser(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ldap_user",
		Description: "LDAP users.",
		List: &plugin.ListConfig{
			Hydrate: listUsers,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "filter", Require: plugin.Optional},
				{Name: "cn", Require: plugin.Optional},
				{Name: "dn", Require: plugin.Optional},
				{Name: "mail", Require: plugin.Optional},
				{Name: "ou", Require: plugin.Optional},
				{Name: "uid", Require: plugin.Optional},
				{Name: "display_name", Require: plugin.Optional},
				{Name: "given_name", Require: plugin.Optional},
				{Name: "description", Require: plugin.Optional},
			},
		},
		// TODO: Add any missing columns that are useful in LDAP/AD
		Columns: []*plugin.Column{
			{
				Name:        "dn",
				Description: "The distinguished name (DN) for this resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "base_dn",
				Description: "The base path to search in.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "filter",
				Description: "The filter to search with.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cn",
				Description: "The user's common name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "The user's description.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "The user's display name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "given_name",
				Description: "The user's given name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "initials",
				Description: "The user's initials.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "mail",
				Description: "The user's email address.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "object_class",
				Description: "The user's object classes.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "ou",
				Description: "The user's organizational unit (OU).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sn",
				Description: "The user's surname.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "uid",
				Description: "The user's ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Uid"),
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
				Transform:   transform.FromField("Cn"),
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

	var baseDN, userObjectFilter string
	var attributes []string
	var limit int64
	var pageSize uint32
	// how do we maintain the default limit for queries? do we make it a configuration?
	limit = 100
	pageSize = 25

	ldapConfig := GetConfig(d.Connection)
	if &ldapConfig != nil {
		if ldapConfig.BaseDN != nil {
			baseDN = *ldapConfig.BaseDN
		}
		if ldapConfig.Attributes != nil {
			attributes = ldapConfig.Attributes
		}
		if ldapConfig.UserObjectFilter != nil {
			userObjectFilter = *ldapConfig.UserObjectFilter
		}
	}

	// Check for all required config args
	if baseDN == "" {
		return nil, errors.New("'base_dn' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}

	keyQuals := d.KeyColumnQuals

	// default value for the user object filter if nothing is passed
	if userObjectFilter == "" {
		userObjectFilter = "(&(objectCategory=person)(objectClass=user))"
	}

	filter := generateFilterString(keyQuals, userObjectFilter)

	if d.QueryContext.Limit != nil {
		if *d.QueryContext.Limit < limit {
			limit = *d.QueryContext.Limit
			if uint32(limit) < pageSize {
				pageSize = uint32(limit)
			}
		}
	}

	logger.Warn("baseDN", baseDN)
	logger.Warn("filter", filter)
	logger.Warn("attributes", attributes)

	var searchReq *ldap.SearchRequest
	// If no attributes are passed in, search request will get all of them
	if attributes != nil {
		searchReq = ldap.NewSearchRequest(baseDN, ldap.ScopeWholeSubtree, 0, int(limit), 0, false, filter, attributes, []ldap.Control{})
	} else {
		searchReq = ldap.NewSearchRequest(baseDN, ldap.ScopeWholeSubtree, 0, int(limit), 0, false, filter, []string{}, []ldap.Control{})
	}

	result, err := conn.SearchWithPaging(searchReq, pageSize)
	// result, err := conn.Search(searchReq)
	if err != nil {
		plugin.Logger(ctx).Error("ldap_user.listUsers", "search_error", err)
		return nil, err
	}

	for _, entry := range result.Entries {
		row := userRow{
			Dn:          entry.DN,
			BaseDn:      baseDN,
			Cn:          entry.GetAttributeValue("cn"),
			Description: entry.GetAttributeValue("description"),
			DisplayName: entry.GetAttributeValue("displayName"),
			GivenName:   entry.GetAttributeValue("givenName"),
			Initials:    entry.GetAttributeValue("initials"),
			Mail:        entry.GetAttributeValue("mail"),
			ObjectClass: entry.GetAttributeValues("objectClass"),
			Ou:          entry.GetAttributeValue("ou"),
			Sn:          entry.GetAttributeValue("sn"),
			Uid:         entry.GetAttributeValue("uid"),
			Attributes:  entry.Attributes,
		}

		if keyQuals["filter"] != nil {
			row.Filter = keyQuals["filter"].GetStringValue()
		}

		d.StreamListItem(ctx, row)
	}

	return nil, nil
}
