package ldap

import (
	"context"
	"time"

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
	// Create Date
	Created *time.Time
	// Modified Date
	Changed *time.Time
	// Object Class
	ObjectClass []string
	// Organizational Unit to which the user belongs
	Ou string
	// Last Name
	Surname string
	// Department
	Department string
	// Object SID
	ObjectSid string
	// SAM Account Name
	SamAccountName string
	// User Principal Name
	UserPrincipalName string
	// Title
	Title string
	// Job Title
	JobTitle string
	// Groups to which the user belongs
	MemberOf []string
	// All attributes that are configured to be returned
	Attributes map[string][]string
}

func tableLDAPUser(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ldap_user",
		Description: "LDAP User",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("dn"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidVolume.NotFound", "InvalidParameterValue"}),
			Hydrate:           getUser,
		},
		List: &plugin.ListConfig{
			Hydrate: listUsers,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "filter", Require: plugin.Optional},
				{Name: "cn", Require: plugin.Optional},
				{Name: "surname", Require: plugin.Optional},
				{Name: "mail", Require: plugin.Optional},
				{Name: "ou", Require: plugin.Optional},
				{Name: "display_name", Require: plugin.Optional},
				{Name: "given_name", Require: plugin.Optional},
				{Name: "department", Require: plugin.Optional},
				{Name: "object_sid", Require: plugin.Optional},
				{Name: "sam_account_name", Require: plugin.Optional},
				{Name: "user_principal_name", Require: plugin.Optional},
				{Name: "description", Require: plugin.Optional},
				{Name: "job_title", Require: plugin.Optional},
				{Name: "created", Operators: []string{">=", "=", "<="}, Require: plugin.Optional},
				{Name: "changed", Operators: []string{">=", "=", "<="}, Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Top Columns
			{
				Name:        "dn",
				Description: "Distinguished Name of the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cn",
				Description: "Full Name of the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "Display Name of the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "given_name",
				Description: "Given Name of the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "surname",
				Description: "Family Name of the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "initials",
				Description: "Initials of the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "mail",
				Description: "E-mail address of the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "department",
				Description: "Department to which the user belongs to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created",
				Description: "Date & Time the user was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "changed",
				Description: "Date & Time the user was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "sam_account_name",
				Description: "Logon Name (pre-Windows 2000) of the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "user_principal_name",
				Description: "Logon Name of the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "job_title",
				Description: "Job Title of the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ou",
				Description: "Organizational Unit to which the user belongs to.",
				Type:        proto.ColumnType_STRING,
			},

			// Other Columns
			{
				Name:        "description",
				Description: "Description of the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "base_dn",
				Description: "The Base DN on which the search was performed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "filter",
				Description: "Optional custom filter passed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "object_sid",
				Description: "Object SID of the user.",
				Type:        proto.ColumnType_STRING,
			},

			// JSON Columns
			{
				Name:        "member_of",
				Description: "Groups that the user is a member of.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "object_class",
				Description: "Object Classes of the user.",
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
				Description: "Title of the user.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Cn"),
			},
		},
	}
}

func getUser(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getUser")

	userDN := d.KeyColumnQuals["dn"].GetStringValue()

	ldapConfig := GetConfig(d.Connection)

	var searchReq *ldap.SearchRequest

	if ldapConfig.Attributes != nil {
		searchReq = ldap.NewSearchRequest(userDN, ldap.ScopeBaseObject, 0, 1, 0, false, "(&)", ldapConfig.Attributes, []ldap.Control{})
	} else {
		searchReq = ldap.NewSearchRequest(userDN, ldap.ScopeBaseObject, 0, 1, 0, false, "(&)", []string{}, []ldap.Control{})
	}

	result, err := search(ctx, d, searchReq)
	if err != nil {
		logger.Error("ldap_user.getUser", "search_error", err)
		return nil, err
	}

	if result.Entries != nil && len(result.Entries) > 0 {
		entry := result.Entries[0]
		row := userRow{
			Dn:                entry.DN,
			BaseDn:            *ldapConfig.BaseDN,
			Cn:                entry.GetAttributeValue("cn"),
			Description:       entry.GetAttributeValue("description"),
			DisplayName:       entry.GetAttributeValue("displayName"),
			GivenName:         entry.GetAttributeValue("givenName"),
			Initials:          entry.GetAttributeValue("initials"),
			Mail:              entry.GetAttributeValue("mail"),
			ObjectClass:       entry.GetAttributeValues("objectClass"),
			Ou:                getOrganizationUnit(entry.DN),
			Surname:           entry.GetAttributeValue("sn"),
			JobTitle:          entry.GetAttributeValue("title"),
			Department:        entry.GetAttributeValue("department"),
			ObjectSid:         getObjectSid(entry),
			SamAccountName:    entry.GetAttributeValue("sAMAccountName"),
			UserPrincipalName: entry.GetAttributeValue("userPrincipalName"),
			MemberOf:          entry.GetAttributeValues("memberOf"),
			Attributes:        transformAttributes(entry.Attributes),
		}

		// Populate Time fields
		if !time.Time.IsZero(*convertToTimestamp(ctx, entry.GetAttributeValue("whenCreated"))) {
			row.Created = convertToTimestamp(ctx, entry.GetAttributeValue("whenCreated"))
		}
		if !time.Time.IsZero(*convertToTimestamp(ctx, entry.GetAttributeValue("whenChanged"))) {
			row.Changed = convertToTimestamp(ctx, entry.GetAttributeValue("whenChanged"))
		}

		return row, nil
	}

	return nil, nil
}

func listUsers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listUsers")

	var baseDN, userObjectFilter string
	var attributes []string
	var pageSize uint32 = PageSize

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

	keyQuals := d.KeyColumnQuals
	quals := d.Quals

	// default value for the user object filter if nothing is passed
	if userObjectFilter == "" {
		userObjectFilter = "(&(objectCategory=person)(objectClass=user))"
	}

	filter := generateFilterString(keyQuals, quals, userObjectFilter)

	if d.QueryContext.Limit != nil {
		if uint32(*d.QueryContext.Limit) < pageSize {
			pageSize = uint32(*d.QueryContext.Limit)
		}
	}

	logger.Warn("baseDN", baseDN)
	logger.Warn("filter", filter)
	logger.Warn("attributes", attributes)

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
				Ou:                getOrganizationUnit(entry.DN),
				Surname:           entry.GetAttributeValue("sn"),
				JobTitle:          entry.GetAttributeValue("title"),
				Department:        entry.GetAttributeValue("department"),
				ObjectSid:         getObjectSid(entry),
				SamAccountName:    entry.GetAttributeValue("sAMAccountName"),
				UserPrincipalName: entry.GetAttributeValue("userPrincipalName"),
				MemberOf:          entry.GetAttributeValues("memberOf"),
				Attributes:        transformAttributes(entry.Attributes),
			}

			if keyQuals["filter"] != nil {
				row.Filter = keyQuals["filter"].GetStringValue()
			}

			// Populate Time fields
			if !time.Time.IsZero(*convertToTimestamp(ctx, entry.GetAttributeValue("whenCreated"))) {
				row.Created = convertToTimestamp(ctx, entry.GetAttributeValue("whenCreated"))
			}
			if !time.Time.IsZero(*convertToTimestamp(ctx, entry.GetAttributeValue("whenChanged"))) {
				row.Changed = convertToTimestamp(ctx, entry.GetAttributeValue("whenChanged"))
			}

			d.StreamListItem(ctx, row)

			// Stop stearming items if the limit has been hit or in case of manual cancellation
			if plugin.IsCancelled(ctx) {
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
