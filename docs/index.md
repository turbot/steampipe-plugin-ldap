---
organization: Turbot
category: ["software development"]
icon_url: "/images/plugins/turbot/ldap.svg"
brand_color: "#CC2025"
display_name: "LDAP"
short_name: "ldap"
description: "Steampipe plugin for querying users, groups, organizational units and more from LDAP."
og_description: "Query LDAP with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/ldap-social-graphic.png"
---

# LDAP + Steampipe

[LDAP](https://ldap.com/) is a mature, flexible, and well supported standards-based mechanism for interacting with directory servers. It’s often used for authentication and storing information about users, groups, and applications, but an LDAP directory server is a fairly general-purpose data store and can be used in a wide variety of applications.

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

For example:

```sql
select
  dn,
  created,
  mail,
  department
from
  ldap_user;
```

```
+---------------------------------------------------------------+---------------------+---------------------------------+-------------+
| dn                                                            | created             | mail                            | department  |
+---------------------------------------------------------------+---------------------+---------------------------------+-------------+
| CN=Emine Braun,OU=Users,DC=example,DC=domain,DC=com           | 2021-08-30 11:21:05 | Emine.Braun@example.com         | IT          |
| CN=Richardis Lamprecht,OU=Users,DC=example,DC=domain,DC=com   | 2021-08-30 11:21:05 | Richardis.Lamprecht@example.com | Engineering |
| CN=Michl Gehring,OU=Users,DC=example,DC=domain,DC=com         | 2021-08-30 11:21:05 | Michl.Gehring@example.com       | Sales       |
| CN=Ottobert Giesen,OU=Users,DC=example,DC=domain,DC=com       | 2021-08-30 11:21:05 | Ottobert.Giesen@example.com     | Marketing   |
| CN=Mirjam Merker,OU=Users,DC=example,DC=domain,DC=com         | 2021-08-30 11:21:05 | Mirjam.Merker@example.com       | Engineering |
+---------------------------------------------------------------+---------------------+---------------------------------+-------------+
```

## Documentation

- **[Table definitions & examples →](/plugins/turbot/ldap/tables)**

## Get started

### Install

Download and install the latest LDAP plugin:

```bash
steampipe plugin install ldap
```

### Configuration

Installing the latest ldap plugin will create a config file (`~/.steampipe/config/ldap.spc`) with a single connection named `ldap`:

```hcl
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
```

## Get Involved

- Open source: https://github.com/turbot/steampipe-plugin-ldap
- Community: [Slack Channel](https://steampipe.io/community/join)
