package ldap

import (
	"context"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

type organizationalUnitRow struct {
	// Distinguished name
	Dn string
	// Base domain name
	BaseDn string
	// Filter string
	Filter string
	// Organizational unit name
	Ou string
	// Description
	Description string
	// Creation Date
	WhenCreated *time.Time
	// Last modified Date
	WhenChanged *time.Time
	// Object class
	ObjectClass []string
	// Title
	Title string
	// Entity that manages the organizational unit
	ManagedBy string
	// All attributes that are configured to be returned
	Attributes map[string][]string
}

func tableLDAPOrganizationalUnit(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ldap_organizational_unit",
		Description: "LDAP Organizational Unit",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("dn"),
			Hydrate:    getOrganizationalUnit,
		},
		List: &plugin.ListConfig{
			Hydrate: listOrganizationalUnits,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "filter", Require: plugin.Optional},
				{Name: "ou", Require: plugin.Optional},
				{Name: "description", Require: plugin.Optional},
				{Name: "when_created", Operators: []string{">", ">=", "=", "<", "<="}, Require: plugin.Optional},
				{Name: "when_changed", Operators: []string{">", ">=", "=", "<", "<="}, Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Top Columns
			{
				Name:        "dn",
				Description: "Distinguished Name of the organizational unit.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ou",
				Description: "Name of the organizational unit",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "managed_by",
				Description: "The Person/Group that manages the organizational unit.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "when_created",
				Description: "Date & Time the organizational unit was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "when_changed",
				Description: "Date & Time the organizational unit was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
			},

			// Other Columns
			{
				Name:        "description",
				Description: "Description of the organizational unit.",
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

			// JSON Columns
			{
				Name:        "object_class",
				Description: "Object Classes of the organizational unit.",
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
				Description: "Title of the organizational unit.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Ou"),
			},
		},
	}
}

func getOrganizationalUnit(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("ldap_organizational_unit.getOrganizationalUnit")

	organizationalUnitDN := d.KeyColumnQuals["dn"].GetStringValue()

	ldapConfig := GetConfig(d.Connection)

	var searchReq *ldap.SearchRequest

	if ldapConfig.Attributes != nil {
		searchReq = ldap.NewSearchRequest(organizationalUnitDN, ldap.ScopeBaseObject, 0, 1, 0, false, "(&)", ldapConfig.Attributes, []ldap.Control{})
	} else {
		searchReq = ldap.NewSearchRequest(organizationalUnitDN, ldap.ScopeBaseObject, 0, 1, 0, false, "(&)", []string{}, []ldap.Control{})
	}

	result, err := search(ctx, d, searchReq)
	if err != nil {
		logger.Error("ldap_organizational_unit.getOrganizationalUnit", "search_error", err)
		return nil, err
	}

	if result.Entries != nil && len(result.Entries) > 0 {
		entry := result.Entries[0]
		row := organizationalUnitRow{
			Dn:          entry.DN,
			BaseDn:      *ldapConfig.BaseDN,
			Ou:          entry.GetAttributeValue("ou"),
			Description: entry.GetAttributeValue("description"),
			ObjectClass: entry.GetAttributeValues("objectClass"),
			ManagedBy:   entry.GetAttributeValue("managedBy"),
			Attributes:  transformAttributes(ctx, entry.Attributes),
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

func listOrganizationalUnits(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("ldap_organizational_unit.listOrganizationalUnits")

	var baseDN, organizationalUnitObjectFilter string
	var attributes []string
	var pageSize uint32 = PageSize

	ldapConfig := GetConfig(d.Connection)
	if ldapConfig.BaseDN != nil {
		baseDN = *ldapConfig.BaseDN
	}
	if ldapConfig.Attributes != nil {
		attributes = ldapConfig.Attributes
	}
	if ldapConfig.OrganizationalUnitObjectFilter != nil {
		organizationalUnitObjectFilter = *ldapConfig.OrganizationalUnitObjectFilter
	}

	keyQuals := d.KeyColumnQuals
	quals := d.Quals

	// default value for the organizational unit object filter if nothing is passed
	if organizationalUnitObjectFilter == "" {
		organizationalUnitObjectFilter = "(objectClass=organizationalUnit)"
	}

	filter := generateFilterString(keyQuals, quals, organizationalUnitObjectFilter)

	logger.Debug("ldap_organizational_unit.listOrganizationalUnits", "baseDN", baseDN)
	logger.Debug("ldap_organizational_unit.listOrganizationalUnits", "filter", filter)
	logger.Debug("ldap_organizational_unit.listOrganizationalUnits", "attributes", attributes)

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
			logger.Error("ldap_organizational_unit.listOrganizationalUnits", "search_error", err)
			return nil, err
		}

		for _, entry := range result.Entries {
			row := organizationalUnitRow{
				Dn:          entry.DN,
				BaseDn:      *ldapConfig.BaseDN,
				Ou:          entry.GetAttributeValue("ou"),
				Description: entry.GetAttributeValue("description"),
				ObjectClass: entry.GetAttributeValues("objectClass"),
				ManagedBy:   entry.GetAttributeValue("managedBy"),
				Attributes:  transformAttributes(ctx, entry.Attributes),
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
