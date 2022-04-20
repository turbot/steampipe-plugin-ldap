package ldap

import (
	"context"
	"strconv"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

type userRow struct {
	// Distinguished name
	Dn string
	// Base domain name
	BaseDn string
	// Filter string
	Filter string
	// Common name / Full name
	Cn string
	// Description
	Description string
	// Display name / Full name
	DisplayName string
	// First name
	GivenName string
	// Middle initials
	Initials string
	// Email address
	Mail string
	// Creation date
	WhenCreated *time.Time
	// Last modified date
	WhenChanged *time.Time
	// Object class
	ObjectClass []string
	// Organizational unit the user belongs to
	Ou string
	// Last name
	Surname string
	// Department
	Department string
	// Object SID
	ObjectSid string
	// SAM account name
	SamAccountName string
	// User principal name
	UserPrincipalName string
	// Title
	Title string
	// Job title
	JobTitle string
	// Groups the user belongs to
	MemberOf []string
	// Whether the user account is disabled
	Disabled *bool
	// All attributes that are configured to be returned
	Attributes map[string][]string
}

func tableLDAPUser(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ldap_user",
		Description: "A user is known as the customer or end-user.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("dn"),
			Hydrate:    getUser,
		},
		List: &plugin.ListConfig{
			Hydrate: listUsers,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "cn", Require: plugin.Optional},
				{Name: "department", Require: plugin.Optional},
				{Name: "description", Require: plugin.Optional},
				{Name: "disabled", Operators: []string{"<>", "="}, Require: plugin.Optional},
				{Name: "display_name", Require: plugin.Optional},
				{Name: "filter", Require: plugin.Optional},
				{Name: "given_name", Require: plugin.Optional},
				{Name: "mail", Require: plugin.Optional},
				{Name: "object_sid", Require: plugin.Optional},
				{Name: "sam_account_name", Require: plugin.Optional},
				{Name: "surname", Require: plugin.Optional},
				{Name: "user_principal_name", Require: plugin.Optional},
				{Name: "when_changed", Operators: []string{">", ">=", "=", "<", "<="}, Require: plugin.Optional},
				{Name: "when_created", Operators: []string{">", ">=", "=", "<", "<="}, Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Top Columns
			{
				Name:        "dn",
				Description: "Distinguished name of the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "cn",
				Description: "Full name of the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "display_name",
				Description: "Display name of the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "object_sid",
				Description: "The security identifier (SID) of the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "given_name",
				Description: "Given name of the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "surname",
				Description: "Family name of the user.",
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
				Name:        "when_created",
				Description: "Date when the user was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "when_changed",
				Description: "Date when the user was last changed.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "sam_account_name",
				Description: "Logon name (pre-Windows 2000) of the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "user_principal_name",
				Description: "Login name of the user, usually mapped to the user email name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "job_title",
				Description: "Job title of the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ou",
				Description: "Organizational unit to which the user belongs to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "disabled",
				Description: "Whether the user account is disabled.",
				Type:        proto.ColumnType_BOOL,
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
				Description: "Optional search filter.",
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
				Description: "Object classes of the user.",
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
	logger.Trace("ldap_user.getUser")

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
			Attributes:        transformAttributes(ctx, entry.Attributes),
			Disabled:          verifyUserDisabled(ctx, entry),
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

func listUsers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("ldap_user.listUsers")

	var baseDN, userObjectFilter string
	var attributes []string
	var pageSize uint32 = PageSize

	ldapConfig := GetConfig(d.Connection)

	if ldapConfig.BaseDN != nil {
		baseDN = *ldapConfig.BaseDN
	}
	if ldapConfig.Attributes != nil {
		attributes = ldapConfig.Attributes
	}
	if ldapConfig.UserObjectFilter != nil {
		userObjectFilter = *ldapConfig.UserObjectFilter
	}

	keyQuals := d.KeyColumnQuals
	quals := d.Quals

	// default value for the user object filter if nothing is passed
	if userObjectFilter == "" {
		userObjectFilter = "(&(objectCategory=person)(objectClass=user))"
	}

	filter := generateFilterString(keyQuals, quals, userObjectFilter)

	logger.Debug("ldap_user.listUsers", "baseDN", baseDN)
	logger.Debug("ldap_user.listUsers", "filter", filter)
	logger.Debug("ldap_user.listUsers", "attributes", attributes)

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
				Attributes:        transformAttributes(ctx, entry.Attributes),
				Disabled:          verifyUserDisabled(ctx, entry),
			}

			if keyQuals["filter"] != nil {
				row.Filter = keyQuals["filter"].GetStringValue()
			}

			// Populate Time fields
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

func verifyUserDisabled(ctx context.Context, entry *ldap.Entry) *bool {
	var disabled bool
	userAccountControl := entry.GetAttributeValue("userAccountControl")
	if userAccountControl == "" {
		return nil
	}
	control, err := strconv.Atoi(userAccountControl)
	if err != nil {
		plugin.Logger(ctx).Error("ldap_user.verifyUserDisabled", "Error while converting userAccountControl to integer", err)
		return nil
	}
	// If the masking of the second bit returns 2, it means that the account is disabled
	// Refer - http://www.selfadsi.org/ads-attributes/user-userAccountControl.htm
	if control&2 == 2 {
		disabled = true
	}
	if control&2 != 2 {
		disabled = false
	}
	return &disabled
}
