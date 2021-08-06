package main

import (
	"fmt"
	"log"

	"github.com/go-ldap/ldap"
)

func main2() {
	ldapURL := "ldap://localhost:10389"
	l, err := ldap.DialURL(ldapURL)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	// connect code comes here
	user := "Hubert J. Farnsworth"
	baseDN := "DC=planetexpress,DC=com"
	filter := fmt.Sprintf("(CN=%s)", ldap.EscapeFilter(user))

	// Filters must start and finish with ()!
	searchReq := ldap.NewSearchRequest(baseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, []string{"sAMAccountName"}, []ldap.Control{})

	result, err := l.Search(searchReq)
	if err != nil {
		log.Fatal(err)
		//return fmt.Errorf("failed to query LDAP: %w", err)
	}

	log.Println("Got", len(result.Entries), "search results")
	result.PrettyPrint(2)

	err = l.Bind("cn=Hubert J. Farnsworth,ou=people,dc=planetexpress,dc=com", "professor")
	if err != nil {
		log.Fatal(err)
	}
}
