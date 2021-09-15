package ldap

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/go-objectsid"
	"github.com/go-ldap/ldap/v3"
	"github.com/iancoleman/strcase"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

// Define the constant page size to be used by all ldap tables
const PageSize uint32 = 1000

func connect(_ context.Context, d *plugin.QueryData) (*ldap.Conn, error) {

	// Load connection from cache
	cacheKey := "ldap"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*ldap.Conn), nil
	}

	var username, password, url, baseDN string
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
		if ldapConfig.URL != nil {
			url = *ldapConfig.URL
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
	if url == "" {
		return nil, errors.New("'url' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}
	if baseDN == "" {
		return nil, errors.New("'base_dn' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}

	var ldapConn *ldap.Conn
	var connErr error

	if tlsRequired {
		ldapURL := fmt.Sprintf("ldaps://%s", url)
		ldapConn, connErr = ldap.DialURL(ldapURL, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: tlsInsecureSkipVerify}))
	} else {
		ldapURL := fmt.Sprintf("ldap://%s", url)
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

func generateFilterString(keyQuals map[string]*proto.QualValue, objectFilter string) string {
	var andClauses strings.Builder

	if keyQuals["filter"] != nil {
		andClauses.WriteString(keyQuals["filter"].GetStringValue())
	} else {
		for key, value := range keyQuals {
			if key == "filter" {
				continue
			}
			var clause string
			if value.GetStringValue() != "" {
				clause = buildClause(key, value.GetStringValue())
			} else if value.GetListValue() != nil {
				clause = generateOrClause(key, value.GetListValue())
			}
			andClauses.WriteString(clause)
		}
	}

	return "(&" + objectFilter + andClauses.String() + ")"
}

func generateOrClause(key string, orValues *proto.QualValueList) string {
	var clauses strings.Builder

	for _, value := range orValues.Values {
		clauses.WriteString(buildClause(key, value.GetStringValue()))
	}

	return "(|" + clauses.String() + ")"
}

func buildClause(key string, value string) string {
	return "(" + strcase.ToLowerCamel(key) + "=" + value + ")"
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
