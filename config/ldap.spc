connection "ldap" {
  plugin = "ldap"

  # The following set of properties are mandatory for the ldap plugin to make a connection to the server
  # Distinguished name of the user which will be used to bind to the server
  # username = "CN=Admin,OU=Users,DC=domain,DC=example,DC=com"

  # The corresponding password of the user defined above
  # password = "55j%@8RnFakePassword"

  # Host to which the plugin will connect to e.g. ad.example.com, ldap.example.com etc.
  # host = "domain.example.com"

  # Port on which the directory server is listening i.e. 389, 636 etc.
  # port = "389"

  # Set to true to enable TLS encryption
  # tls_required = false

  # Distinguished name of the base object on which queries will be executed
  # base_dn = "DC=domain,DC=example,DC=com"

  # Fixed set of attributes that will be requested for each LDAP query. This attribute list is shared across all tables.
  # If nothing is provided, Steampipe will request for all attributes
  # attributes = ["cn", "displayName", "uid"]

  # Optional user object filter to be used to filter objects. If not provided, defaults to - "(&(objectCategory=person)(objectClass=user))"
  # user_object_filter = "(&(objectCategory=person)(objectClass=user))"

  # Optional group object filter to be used to filter objects. If not provided, defaults to - "(objectClass=group)"
  # group_object_filter = "(objectClass=group)"

  # Optional organizational object filter to be used to filter objects. If not provided, defaults to - "(objectClass=organizationalUnit)"
  # ou_object_filter = "(objectClass=organizationalUnit)"
}
