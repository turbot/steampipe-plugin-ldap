package ldap

import (
	"context"
	"fmt"
	"log"
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

	var username string

	ldapConfig := GetConfig(d.Connection)
	if &ldapConfig != nil {
		if ldapConfig.Username != nil {
			username = *ldapConfig.Username
		}
	}

	// Error if missing config
	if username == "" {
		return nil, fmt.Errorf("Partial credentials found in connection config, missing: username")
	}

	ldapURL := "ldap://localhost:10389"
	conn, err := ldap.DialURL(ldapURL)
	if err != nil {
		log.Fatal(err)
		return nil, err
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
