package ldap

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/go-objectsid"
	"github.com/go-ldap/ldap/v3"
	"github.com/iancoleman/strcase"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

// Define the constant page size to be used by all ldap tables
const PageSize uint32 = 1000

// Define the time filter timestamo format
const FilterTimestampFormat = "20060102150405.000Z"

// Disabled User Filter
const DisabledUserFilter = "(userAccountControl:1.2.840.113556.1.4.803:=2)"

func connect(_ context.Context, d *plugin.QueryData) (*ldap.Conn, error) {

	// Load connection from cache
	cacheKey := "ldap"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*ldap.Conn), nil
	}

	var username, password, host, port, baseDN string
	tlsRequired := false
	tlsInsecureSkipVerify := false

	ldapConfig := GetConfig(d.Connection)
	if &ldapConfig != nil {
		if ldapConfig.Username != nil {
			username = *ldapConfig.Username
		}
		if ldapConfig.Password != nil {
			password = *ldapConfig.Password
		}
		if ldapConfig.Host != nil {
			host = *ldapConfig.Host
		}
		if ldapConfig.Port != nil {
			port = *ldapConfig.Port
		}
		if ldapConfig.TLSRequired != nil {
			tlsRequired = *ldapConfig.TLSRequired
		}
		if ldapConfig.TLSInsecureSkipVerify != nil {
			tlsInsecureSkipVerify = *ldapConfig.TLSInsecureSkipVerify
		}
		if ldapConfig.BaseDN != nil {
			baseDN = *ldapConfig.BaseDN
		}
	}

	// Check for all required config args
	if username == "" {
		return nil, errors.New("'username' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}
	if password == "" {
		return nil, errors.New("'password' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}
	if host == "" {
		return nil, errors.New("'host' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}
	if port == "" {
		return nil, errors.New("'port' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}
	if baseDN == "" {
		return nil, errors.New("'base_dn' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}

	var ldapConn *ldap.Conn
	var connErr error

	if tlsRequired {
		ldapURL := fmt.Sprintf("ldaps://%s:%s", host, port)
		ldapConn, connErr = ldap.DialURL(ldapURL, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: tlsInsecureSkipVerify}))
	} else {
		ldapURL := fmt.Sprintf("ldap://%s:%s", host, port)
		ldapConn, connErr = ldap.DialURL(ldapURL)
	}

	if connErr != nil {
		return nil, connErr
	}

	if err := ldapConn.Bind(username, password); err != nil {
		return nil, err
	}

	// Save to cache
	// TODO: Use SetWithTTL once we know what default timeout is
	d.ConnectionManager.Cache.Set(cacheKey, ldapConn)

	return ldapConn, nil
}

func reconnect(ctx context.Context, d *plugin.QueryData) (*ldap.Conn, error) {
	d.ConnectionManager.Cache.Delete("ldap")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ldap_utils.reconnect", "reconnect_error", err)
		return nil, err
	}
	return conn, nil
}

func search(ctx context.Context, d *plugin.QueryData, searchReq *ldap.SearchRequest) (*ldap.SearchResult, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ldap_utils.search", "connection_error", err)
		return nil, err
	}
	searchResult, e := conn.Search(searchReq)
	if e != nil && ldap.IsErrorWithCode(e, 200) {
		plugin.Logger(ctx).Info("LDAP Connection closed, trying to reconnect ...")
		conn, err := reconnect(ctx, d)
		if err != nil {
			return nil, err
		}
		searchResult, err := conn.Search(searchReq)
		if err != nil {
			return nil, err
		}
		return searchResult, nil
	}
	return searchResult, nil
}

func isNotFoundError(notFoundErrors []string) plugin.ErrorPredicate {
	return func(err error) bool {
		errMsg := err.Error()
		for _, msg := range notFoundErrors {
			if strings.Contains(errMsg, msg) {
				return true
			}
		}
		return false
	}
}

func generateFilterString(keyQuals map[string]*proto.QualValue, quals map[string]*plugin.KeyColumnQuals, objectFilter string) string {
	var andClauses strings.Builder

	if keyQuals["filter"] != nil {
		andClauses.WriteString(keyQuals["filter"].GetStringValue())
	} else {
		// Range over the key quals
		for key, value := range keyQuals {
			if key == "filter" {
				continue
			}
			var clause string
			if value.GetStringValue() != "" {
				clause = buildClause(key, value.GetStringValue(), "=")
			} else if value.GetListValue() != nil {
				clause = generateOrClause(key, value.GetListValue())
			}
			andClauses.WriteString(clause)
		}

		// Get individual quals
		if quals["created"] != nil {
			for _, q := range quals["created"].Quals {
				timeString := q.Value.GetTimestampValue().AsTime().Format(FilterTimestampFormat)
				clause := buildClause("whenCreated", timeString, q.Operator)
				andClauses.WriteString(clause)
			}
		}

		if quals["changed"] != nil {
			for _, q := range quals["changed"].Quals {
				timeString := q.Value.GetTimestampValue().AsTime().Format(FilterTimestampFormat)
				clause := buildClause("whenChanged", timeString, q.Operator)
				andClauses.WriteString(clause)
			}
		}

		if quals["disabled"] != nil {
			for _, q := range quals["disabled"].Quals {
				value := q.Value
				if value != nil {
					clause := DisabledUserFilter
					if q.Operator == "<>" {
						clause = "(!" + clause + ")"
					}
					andClauses.WriteString(clause)
				}
			}
		}
	}

	return "(&" + objectFilter + andClauses.String() + ")"
}

func generateOrClause(key string, orValues *proto.QualValueList) string {
	var clauses strings.Builder

	for _, value := range orValues.Values {
		clauses.WriteString(buildClause(key, value.GetStringValue(), "="))
	}

	return "(|" + clauses.String() + ")"
}

func buildClause(key string, value string, operator string) string {
	return "(" + strcase.ToLowerCamel(key) + operator + value + ")"
}

func getOrganizationUnit(dn string) string {
	return dn[strings.Index(strings.ToUpper(dn), "OU"):]
}

func getObjectSid(entry *ldap.Entry) string {
	rawObjectSid := entry.GetRawAttributeValue("objectSid")
	if len(rawObjectSid) > 0 {
		return objectsid.Decode(rawObjectSid).String()
	}
	return ""
}

func convertToTimestamp(ctx context.Context, str string) *time.Time {
	// If there is a blank string, return zero time
	if str == "" {
		return &time.Time{}
	}

	// Frame the layout according to the data available. The front part remains constant to '20060102150405'
	// The second part i.e. after '.' can have a variable number of 0's followed by Z
	layout := "20060102150405." + strings.Split(str, ".")[1]
	t, err := time.Parse(layout, str)
	if err != nil {
		plugin.Logger(ctx).Error("ldap_utils.convertToTimestamp", "conversion_error", err)
		// Return zero time in case of a conversion error
		return &time.Time{}
	}
	// Return the converted time if conversion is successful
	return &t
}

func transformAttributes(attributes []*ldap.EntryAttribute) map[string][]string {
	var data = make(map[string][]string)
	for _, attribute := range attributes {
		data[attribute.Name] = attribute.Values
	}
	return data
}
