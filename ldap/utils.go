package ldap

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func connect(_ context.Context, d *plugin.QueryData) (*ldap.Conn, error) {

	// Load connection from cache
	cacheKey := "ldap"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*ldap.Conn), nil
	}

	var username, password, url string
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
