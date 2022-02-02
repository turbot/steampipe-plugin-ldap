connection "ldap" {
  plugin = "ldap"

  # Distinguished name of the user which will be used to bind to the server
  # username = "CN=Admin,OU=Users,DC=domain,DC=example,DC=com"

  # The password for the user defined above
  # password = "55j%@8RnFakePassword"

  # Host to connect to, e.g. ad.example.com, ldap.example.com
  # host = "domain.example.com"

  # Port on which the directory server is listening, e.g., 389, 636
  # port = "389"

  # If true, enable TLS encryption
  # tls_required = false

  # Distinguished name of the base object on which queries will be executed
  # base_dn = "DC=domain,DC=example,DC=com"

  # Fixed set of attributes that will be requested for each LDAP query. This attribute list is shared across all tables.
  # If nothing is specified, Steampipe will request all attributes
  # attributes = ["cn", "displayName", "uid"]

  # Optional user object filter to be used to filter objects. If not provided, defaults to "(&(objectCategory=person)(objectClass=user))"
  # user_object_filter = "(&(objectCategory=person)(objectClass=user))"

  # Optional group object filter to be used to filter objects. If not provided, defaults to "(objectClass=group)"
  # group_object_filter = "(objectClass=group)"

  # Optional organizational object filter to be used to filter objects. If not provided, defaults to "(objectClass=organizationalUnit)"
  # ou_object_filter = "(objectClass=organizationalUnit)"
}
