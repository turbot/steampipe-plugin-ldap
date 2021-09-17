connection "ldap" {
  plugin     = "ldap"

  # The following set of properties are mandatory for the ldap plugin to make a connection to the server
  # Distinguished name of the user which will be used to bind to the server
  # username   = "CN=Admin,OU=Users,DC=example,DC=domain,DC=com"

  # The corresponding password of the user defined above
  # password   = "55j%@8Rn[Ct8#\Mm"

  # Url to which the plugin will connect to in the format -> <ip-address>:<port>
  # url        = "10.84.11.5:389"

  # Distinguished name of the base object on which queries will be executed
  # base_dn    = "DC=example,DC=domain,DC=com"

  # Fixed set of attributes that will be requested for each LDAP query. This attribute list is shared across all tables.
  # If nothing is provided, Steampipe will request for all attributes
  # attributes = ["cn", "displayName", "uid"]

  # Should be set to true if you want to secure communications via SSL/Transport Layer Security(TLS) technology
  # tls_required = false
  # Set to true for TLS accepting any certificate presented by the server and any host name in that certificate
  # tls_insecure_skip_verify = true

  # Optional user object filter to be used to filter objects. If not provided, defaults to - "(&(objectCategory=person)(objectClass=user))"
  # user_object_filter = "(&(objectCategory=person)(objectClass=user))"
}
