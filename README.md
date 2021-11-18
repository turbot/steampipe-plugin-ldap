![image](https://hub.steampipe.io/images/plugins/turbot/ldap-social-graphic.png)

# LDAP Plugin for Steampipe

Use SQL to query infrastructure including users, groups, organizational units and more from LDAP.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/ldap)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/ldap/tables)
- Community: [Slack Channel](https://steampipe.io/community/join)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-ldap/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install ldap
```

Run a query:

```sql
select dn, mail, department from ldap_user
```

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/turbot/steampipe-plugin-ldap.git
cd steampipe-plugin-ldap
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```
make
```

Configure the plugin:

```
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/ldap.spc
```

Try it!

```
steampipe query
> .inspect ldap
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-ldap/blob/main/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [ldap Plugin](https://github.com/turbot/steampipe-plugin-ldap/labels/help%20wanted)
