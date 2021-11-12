package ldap

import (
	"context"
	"time"

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
	// Create Date
	CreatedOn *time.Time
	// Modified Date
	ModifiedOn *time.Time
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
	// Groups to which the group belongs
	MemberOf []string
	// All attributes that are configured to be returned
	Attributes map[string][]string
}

func tableLDAPGroup(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ldap_group",
		Description: "A group is a collection of digital identities i.e. users, groups etc.",
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
				{Name: "created_on", Operators: []string{">=", "=", "<="}, Require: plugin.Optional},
				{Name: "modified_on", Operators: []string{">=", "=", "<="}, Require: plugin.Optional},
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
				Name:        "created_on",
				Description: "Date & time the group was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "modified_on",
				Description: "Date & time the group was last modified.",
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
				Description: "The Base distinguished name on which the search was performed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "filter",
				Description: "Optional custom filter passed.",
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
	logger.Trace("getGroup")

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
			Attributes:     transformAttributes(entry.Attributes),
		}

		// Populate Time fields
		if !time.Time.IsZero(*convertToTimestamp(ctx, entry.GetAttributeValue("whenCreated"))) {
			row.CreatedOn = convertToTimestamp(ctx, entry.GetAttributeValue("whenCreated"))
		}
		if !time.Time.IsZero(*convertToTimestamp(ctx, entry.GetAttributeValue("whenChanged"))) {
			row.ModifiedOn = convertToTimestamp(ctx, entry.GetAttributeValue("whenChanged"))
		}

		return row, nil
	}

	return nil, nil
}

func listGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listGroups")

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

	keyQuals := d.KeyColumnQuals
	quals := d.Quals

	// default value for the group object filter if nothing is passed
	groupObjectFilter = "(objectClass=group)"

	filter := generateFilterString(logger, keyQuals, quals, groupObjectFilter)

	if d.QueryContext.Limit != nil {
		if uint32(*d.QueryContext.Limit) < pageSize {
			pageSize = uint32(*d.QueryContext.Limit)
		}
	}

	logger.Info("baseDN", baseDN)
	logger.Info("filter", filter)
	logger.Info("attributes", attributes)

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
			logger.Error("ldap_group.listGroups", "search_error", err)
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
				Attributes:     transformAttributes(entry.Attributes),
			}

			if keyQuals["filter"] != nil {
				row.Filter = keyQuals["filter"].GetStringValue()
			}

			// Populate Time fields
			if !time.Time.IsZero(*convertToTimestamp(ctx, entry.GetAttributeValue("whenCreated"))) {
				row.CreatedOn = convertToTimestamp(ctx, entry.GetAttributeValue("whenCreated"))
			}
			if !time.Time.IsZero(*convertToTimestamp(ctx, entry.GetAttributeValue("whenChanged"))) {
				row.ModifiedOn = convertToTimestamp(ctx, entry.GetAttributeValue("whenChanged"))
			}

			d.StreamListItem(ctx, row)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
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
