package ldap

import (
	"context"

	"github.com/go-ldap/ldap/v3"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

type organizationalUnitRow struct {
	// Distinguished Name
	Dn string
	// Base Domain Name
	BaseDn string
	// Filter string (if passed as query clause)
	Filter string
	// Name of the Organizational Unit
	Ou string
	// Description
	Description string
	// Object Class
	ObjectClass []string
	// Title
	Title string
	// Entity that manages the Organizational Unit
	ManagedBy string
	// All attributes that are configured to be returned
	Attributes []*ldap.EntryAttribute
	// Raw data from LDAP
	Raw []string
}

func tableLDAPOrganizationalUnit(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "ldap_organizational_unit",
		Description: "LDAP Organizational Unit",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("dn"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidVolume.NotFound", "InvalidParameterValue"}),
			Hydrate:           getOrganizationalUnit,
		},
		List: &plugin.ListConfig{
			Hydrate: listOrganizationalUnits,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "filter", Require: plugin.Optional},
				{Name: "ou", Require: plugin.Optional},
				{Name: "description", Require: plugin.Optional},
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
			{
				Name:        "raw",
				Description: "All attributes along with their raw data values.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromValue(),
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
	logger.Trace("getOrganizationalUnit")

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
			Attributes:  entry.Attributes,
		}
		return row, nil
	}

	return nil, nil
}

func listOrganizationalUnits(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listOrganizationalUnits")

	var baseDN, organizationalUnitObjectFilter string
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
	}

	keyQuals := d.KeyColumnQuals

	// default value for the organizational unit object filter if nothing is passed
	organizationalUnitObjectFilter = "(objectClass=organizationalUnit)"

	filter := generateFilterString(keyQuals, organizationalUnitObjectFilter)

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
				Attributes:  entry.Attributes,
			}

			if keyQuals["filter"] != nil {
				row.Filter = keyQuals["filter"].GetStringValue()
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
