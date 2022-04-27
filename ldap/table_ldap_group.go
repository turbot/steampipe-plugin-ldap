package ldap

import (
	"context"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

type groupRow struct {
	// Distinguished name
	Dn string
	// Base domain name
	BaseDn string
	// Filter string
	Filter string
	// Common name
	Cn string
	// Description
	Description string
	// Creation date
	WhenCreated *time.Time
	// Last modified date
	WhenChanged *time.Time
	// Object class
	ObjectClass []string
	// Organizational unit the group belongs to
	Ou string
	// Object SID
	ObjectSid string
	// SAM account name
	SamAccountName string
	// Title
	Title string
	// Groups the group belongs to
	MemberOf []string
	// All attributes that are configured to be returned
	Attributes map[string][]string
}

func tableLDAPGroup(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ldap_group",
		Description: "A group is a collection of digital identities, e.g., users, groups.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("dn"),
			Hydrate:    getGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: listGroups,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "cn", Require: plugin.Optional},
				{Name: "description", Require: plugin.Optional},
				{Name: "filter", Require: plugin.Optional},
				{Name: "object_sid", Require: plugin.Optional},
				{Name: "sam_account_name", Require: plugin.Optional},
				{Name: "when_changed", Operators: []string{">", ">=", "=", "<", "<="}, Require: plugin.Optional},
				{Name: "when_created", Operators: []string{">", ">=", "=", "<", "<="}, Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Top Columns
			{
				Name:        "dn",
				Description: "Distinguished name of the group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cn",
				Description: "Common/Full name of the group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "object_sid",
				Description: "The security identifier (SID) of the group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ou",
				Description: "Organizational unit to which the group belongs to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "sam_account_name",
				Description: "SAM Account name of the group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "when_created",
				Description: "Date when the group was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "when_changed",
				Description: "Date when the group was last changed.",
				Type:        proto.ColumnType_TIMESTAMP,
			},

			// Other Columns
			{
				Name:        "description",
				Description: "Description of the group.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "base_dn",
				Description: "The Base DN on which the search was performed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "filter",
				Description: "Optional search filter.",
				Type:        proto.ColumnType_STRING,
			},

			// JSON Columns
			{
				Name:        "member_of",
				Description: "Groups that the group is a member of.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "object_class",
				Description: "Object classes of the group.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "attributes",
				Description: "All attributes that have been returned from LDAP.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe Columns
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
	logger.Trace("ldap_group.getGroup")

	groupDN := d.KeyColumnQuals["dn"].GetStringValue()

	ldapConfig := GetConfig(d.Connection)

	var searchReq *ldap.SearchRequest

	if ldapConfig.Attributes != nil {
		searchReq = ldap.NewSearchRequest(groupDN, ldap.ScopeBaseObject, 0, 1, 0, false, "(&)", ldapConfig.Attributes, []ldap.Control{})
	} else {
		searchReq = ldap.NewSearchRequest(groupDN, ldap.ScopeBaseObject, 0, 1, 0, false, "(&)", []string{}, []ldap.Control{})
	}

	result, err := search(ctx, d, searchReq)
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
			MemberOf:       entry.GetAttributeValues("memberOf"),
			Attributes:     transformAttributes(ctx, entry.Attributes),
		}

		// Populate Time fields
		if !time.Time.IsZero(*convertToTimestamp(ctx, entry.GetAttributeValue("whenCreated"))) {
			row.WhenCreated = convertToTimestamp(ctx, entry.GetAttributeValue("whenCreated"))
		}
		if !time.Time.IsZero(*convertToTimestamp(ctx, entry.GetAttributeValue("whenChanged"))) {
			row.WhenChanged = convertToTimestamp(ctx, entry.GetAttributeValue("whenChanged"))
		}

		return row, nil
	}

	return nil, nil
}

func listGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("ldap_group.listGroups")

	var baseDN, groupObjectFilter string
	var attributes []string
	var pageSize uint32 = PageSize

	ldapConfig := GetConfig(d.Connection)
	if ldapConfig.BaseDN != nil {
		baseDN = *ldapConfig.BaseDN
	}
	if ldapConfig.Attributes != nil {
		attributes = ldapConfig.Attributes
	}

	if ldapConfig.GroupObjectFilter != nil {
		groupObjectFilter = *ldapConfig.GroupObjectFilter
	}

	keyQuals := d.KeyColumnQuals
	quals := d.Quals

	// default value for the group object filter if nothing is passed
	if groupObjectFilter == "" {
		groupObjectFilter = "(objectClass=group)"
	}

	filter := generateFilterString(keyQuals, quals, groupObjectFilter)

	logger.Debug("ldap_group.listGroups", "baseDN", baseDN)
	logger.Debug("ldap_group.listGroups", "filter", filter)
	logger.Debug("ldap_group.listGroups", "attributes", attributes)

	if d.QueryContext.Limit != nil {
		if uint32(*d.QueryContext.Limit) < pageSize {
			pageSize = uint32(*d.QueryContext.Limit)
		}
	}

	var searchReq *ldap.SearchRequest
	paging := ldap.NewControlPaging(pageSize)

	for {
		// If no attributes are passed in, search request will get all of them
		if attributes != nil {
			searchReq = ldap.NewSearchRequest(baseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, attributes, []ldap.Control{paging})
		} else {
			searchReq = ldap.NewSearchRequest(baseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, []string{}, []ldap.Control{paging})
		}

		result, err := search(ctx, d, searchReq)
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
				MemberOf:       entry.GetAttributeValues("memberOf"),
				Attributes:     transformAttributes(ctx, entry.Attributes),
			}

			if keyQuals["filter"] != nil {
				row.Filter = keyQuals["filter"].GetStringValue()
			}

			// Populate time fields
			if !time.Time.IsZero(*convertToTimestamp(ctx, entry.GetAttributeValue("whenCreated"))) {
				row.WhenCreated = convertToTimestamp(ctx, entry.GetAttributeValue("whenCreated"))
			}
			if !time.Time.IsZero(*convertToTimestamp(ctx, entry.GetAttributeValue("whenChanged"))) {
				row.WhenChanged = convertToTimestamp(ctx, entry.GetAttributeValue("whenChanged"))
			}

			d.StreamListItem(ctx, row)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		// If the result control does not have paging or if the paging control does not
		// have a next page cookie exit from the loop
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
