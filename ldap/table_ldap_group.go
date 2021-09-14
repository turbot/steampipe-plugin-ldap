package ldap

import (
	"context"

	"github.com/go-ldap/ldap/v3"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

type groupRow struct {
	// Distinguished Name
	Dn string
	// Base Domain Name
	BaseDn string
	// Filter string (if passed as query clause)
	Filter string
	// Common Name
	Cn string
	// Description
	Description string
	// Object Class
	ObjectClass []string
	// Organizational Unit to which the group belongs
	Ou string
	// Object SID
	ObjectSid string
	// SAM Account Name
	SamAccountName string
	// Title
	Title string
	// All attributes that are configured to be returned
	Attributes []*ldap.EntryAttribute
	// Raw data from LDAP
	Raw []string
}

func tableLDAPGroup(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ldap_group",
		Description: "LDAP groups.",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("dn"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidVolume.NotFound", "InvalidParameterValue"}),
			Hydrate:           getGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listGroups,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "filter", Require: plugin.Optional},
				{Name: "cn", Require: plugin.Optional},
				{Name: "ou", Require: plugin.Optional},
				{Name: "object_sid", Require: plugin.Optional},
				{Name: "sam_account_name", Require: plugin.Optional},
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
				Description: "The group's common name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "description",
				Description: "The group's description.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "object_class",
				Description: "The group's object classes.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "ou",
				Description: "The group's organizational unit (OU).",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "object_sid",
				Description: "The Object SID of the group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sam_account_name",
				Description: "The SAM Account Name of the group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "attributes",
				Description: "The attributes of the group.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "raw",
				Description: "The attributes of the group.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "title",
				Description: "Title of the group.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Cn"),
			},
		},
	}
}

func getGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getGroup")

	groupDN := d.KeyColumnQuals["dn"].GetStringValue()

	conn, err := connect(ctx, d)
	if err != nil {
		logger.Error("ldap_group.getGroup", "connection_error", err)
		return nil, err
	}

	ldapConfig := GetConfig(d.Connection)

	var searchReq *ldap.SearchRequest

	if ldapConfig.Attributes != nil {
		searchReq = ldap.NewSearchRequest(groupDN, ldap.ScopeBaseObject, 0, 1, 0, false, "(&)", ldapConfig.Attributes, []ldap.Control{})
	} else {
		searchReq = ldap.NewSearchRequest(groupDN, ldap.ScopeBaseObject, 0, 1, 0, false, "(&)", []string{}, []ldap.Control{})
	}

	result, err := conn.Search(searchReq)
	if err != nil {
		logger.Error("ldap_group.getGroup", "search_error", err)
		return nil, err
	}

	if result.Entries != nil && len(result.Entries) > 0 {
		entry := result.Entries[0]
		row := groupRow{
			Dn:             entry.DN,
			BaseDn:         *ldapConfig.BaseDN,
			Cn:             entry.GetAttributeValue("cn"),
			Description:    entry.GetAttributeValue("description"),
			ObjectClass:    entry.GetAttributeValues("objectClass"),
			Ou:             getOrganizationUnit(entry.DN),
			Title:          entry.GetAttributeValue("title"),
			ObjectSid:      getObjectSid(entry),
			SamAccountName: entry.GetAttributeValue("sAMAccountName"),
			Attributes:     entry.Attributes,
		}
		return row, nil
	}

	return nil, nil
}

func listGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listGroups")

	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ldap_group.listGroups", "connection_error", err)
		return nil, err
	}

	// TODO: Where to close connection?
	//defer conn.Close()

	var baseDN, groupObjectFilter string
	var attributes []string
	var limit *int64
	var pageSize uint32 = PageSize

	ldapConfig := GetConfig(d.Connection)
	if &ldapConfig != nil {
		if ldapConfig.BaseDN != nil {
			baseDN = *ldapConfig.BaseDN
		}
		if ldapConfig.Attributes != nil {
			attributes = ldapConfig.Attributes
		}
	}

	keyQuals := d.KeyColumnQuals

	// default value for the group object filter if nothing is passed
	groupObjectFilter = "(objectClass=group)"

	filter := generateFilterString(keyQuals, groupObjectFilter)

	if d.QueryContext.Limit != nil {
		limit = d.QueryContext.Limit
		if uint32(*limit) < pageSize {
			pageSize = uint32(*limit)
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
			plugin.Logger(ctx).Error("ldap_group.listGroups", "search_error", err)
			return nil, err
		}

		for _, entry := range result.Entries {
			row := groupRow{
				Dn:             entry.DN,
				BaseDn:         baseDN,
				Cn:             entry.GetAttributeValue("cn"),
				Description:    entry.GetAttributeValue("description"),
				ObjectClass:    entry.GetAttributeValues("objectClass"),
				Ou:             getOrganizationUnit(entry.DN),
				Title:          entry.GetAttributeValue("title"),
				ObjectSid:      getObjectSid(entry),
				SamAccountName: entry.GetAttributeValue("sAMAccountName"),
				Attributes:     entry.Attributes,
			}

			if keyQuals["filter"] != nil {
				row.Filter = keyQuals["filter"].GetStringValue()
			}

			d.StreamListItem(ctx, row)

			// Decrement the limit and exit outer loop if all results have been streamed or in case of manual cancellation
			if limit != nil {
				*limit--
				if *limit == 0 || plugin.IsCancelled(ctx) {
					break out
				}
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
