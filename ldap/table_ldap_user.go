package ldap

import (
	"context"

	"github.com/go-ldap/ldap/v3"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

type userRow struct {
	// Distinguished Name
	Dn string
	// Base Domain Name
	BaseDn string
	// Filter string (if passed as query clause)
	Filter string
	// Common Name / Full Name
	Cn string
	// Description
	Description string
	// Display Name / Full Name
	DisplayName string
	// First name
	GivenName string
	// Middle Initials
	Initials string
	// Email id
	Mail string
	// Object Class
	ObjectClass []string
	// Organizational Unit to which the user belongs
	Ou string
	// Last Name
	Sn string
	// User ID
	Uid string
	// Department
	Department string
	// Object SID
	ObjectSid string
	// SAM Account Name
	SamAccountName string
	// User Principal Name
	UserPrincipalName string
	// Job Title
	Title string
	// All attributes that are configured to be returned
	Attributes []*ldap.EntryAttribute
	// Raw data from LDAP
	Raw []string
}

func tableLDAPUser(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ldap_user",
		Description: "LDAP users.",
		List: &plugin.ListConfig{
			Hydrate: listUsers,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "filter", Require: plugin.Optional},
				{Name: "cn", Require: plugin.Optional},
				{Name: "sn", Require: plugin.Optional},
				{Name: "dn", Require: plugin.Optional},
				{Name: "mail", Require: plugin.Optional},
				{Name: "ou", Require: plugin.Optional},
				{Name: "uid", Require: plugin.Optional},
				{Name: "display_name", Require: plugin.Optional},
				{Name: "given_name", Require: plugin.Optional},
				{Name: "department", Require: plugin.Optional},
				{Name: "object_sid", Require: plugin.Optional},
				{Name: "sam_account_name", Require: plugin.Optional},
				{Name: "user_principal_name", Require: plugin.Optional},
				{Name: "description", Require: plugin.Optional},
			},
		},
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
				Name:        "department",
				Description: "The department to which the user belongs to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "object_sid",
				Description: "The Object SID of the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sam_account_name",
				Description: "The SAM Account Name of the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "user_principal_name",
				Description: "The User Principal Name of the user.",
				Type:        proto.ColumnType_STRING,
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
			{
				Name:        "title",
				Description: "Job Title of this resource.",
				Type:        proto.ColumnType_STRING,
				// Transform:   transform.FromField("Cn"),
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
	limit = 1000
	pageSize = 500

	ldapConfig := GetConfig(d.Connection)
	if &ldapConfig != nil {
		if ldapConfig.Attributes != nil {
			attributes = ldapConfig.Attributes
		}
		if ldapConfig.UserObjectFilter != nil {
			userObjectFilter = *ldapConfig.UserObjectFilter
		}
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
	paging := ldap.NewControlPaging(pageSize)

	// label for outer for loop
out:
	for {
		// If no attributes are passed in, search request will get all of them
		if attributes != nil {
			searchReq = ldap.NewSearchRequest(baseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, attributes, []ldap.Control{paging})
		} else {
			searchReq = ldap.NewSearchRequest(baseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, []string{}, []ldap.Control{paging})
		}

		result, err := conn.Search(searchReq)
		if err != nil {
			plugin.Logger(ctx).Error("ldap_user.listUsers", "search_error", err)
			return nil, err
		}

		for _, entry := range result.Entries {
			row := userRow{
				Dn:                entry.DN,
				BaseDn:            baseDN,
				Cn:                entry.GetAttributeValue("cn"),
				Description:       entry.GetAttributeValue("description"),
				DisplayName:       entry.GetAttributeValue("displayName"),
				GivenName:         entry.GetAttributeValue("givenName"),
				Initials:          entry.GetAttributeValue("initials"),
				Mail:              entry.GetAttributeValue("mail"),
				ObjectClass:       entry.GetAttributeValues("objectClass"),
				Ou:                entry.GetAttributeValue("ou"),
				Sn:                entry.GetAttributeValue("sn"),
				Uid:               entry.GetAttributeValue("uid"),
				Title:             entry.GetAttributeValue("title"),
				Department:        entry.GetAttributeValue("department"),
				ObjectSid:         entry.GetAttributeValue("objectSid"),
				SamAccountName:    entry.GetAttributeValue("sAMAccountName"),
				UserPrincipalName: entry.GetAttributeValue("userPrincipalName"),
				Attributes:        entry.Attributes,
			}

			if keyQuals["filter"] != nil {
				row.Filter = keyQuals["filter"].GetStringValue()
			}

			d.StreamListItem(ctx, row)

			// Decrement the limit and exit outer loop if all results have been streamed or in case of manual cancellation
			limit--
			if limit == 0 || plugin.IsCancelled(ctx) {
				break out
			}
		}

		// If the result control does not have paging or if the paging control does not
		// have a next page cookie we exit from the loop
		resultCtrl := ldap.FindControl(result.Controls, paging.GetControlType())
		if resultCtrl == nil {
			break
		}
		if pagingCtrl, ok := resultCtrl.(*ldap.ControlPaging); ok {
			if len(pagingCtrl.Cookie) == 0 {
				break
			}
			paging.SetCookie(pagingCtrl.Cookie)
		}
	}

	return nil, nil
}
