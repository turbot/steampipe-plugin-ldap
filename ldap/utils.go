package ldap

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-ldap/ldap"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func connect(_ context.Context, d *plugin.QueryData) (*ldap.Conn, error) {

	// Load connection from cache
	cacheKey := "ldap"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*ldap.Conn), nil
	}

	var username, password, url string

	ldapConfig := GetConfig(d.Connection)
	if &ldapConfig != nil {
		if ldapConfig.Username != nil {
			username = *ldapConfig.Username
			password = *ldapConfig.Password
			url = *ldapConfig.URL
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

	ldapURL := fmt.Sprintf("ldap://%s", url)
	conn, err := ldap.DialURL(ldapURL)
	if err != nil {
		return nil, errors.New("Failed to connect to LDAP server")
	}

	if err := conn.Bind(username, password); err != nil {
		return nil, errors.New("Failed to bind")
	}

	// Save to cache
	d.ConnectionManager.Cache.Set(cacheKey, conn)

	return conn, nil
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
