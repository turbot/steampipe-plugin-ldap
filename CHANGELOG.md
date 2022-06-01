## v0.1.1 [2022-06-01]

_Bug fixes_

- Fixed `ldap_group` and `ldap_user` table queries crashing on `ou` column if a DN doesn't contain an OU. ([#10](https://github.com/turbot/steampipe-plugin-ldap/pull/10))

## v0.1.0 [2022-04-27]

_Enhancements_

- Added support for native Linux ARM and Mac M1 builds. ([#7](https://github.com/turbot/steampipe-plugin-ldap/pull/7))
- Recompiled plugin with [steampipe-plugin-sdk v3.1.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v310--2022-03-30) and Go version `1.18`. ([#6](https://github.com/turbot/steampipe-plugin-ldap/pull/6))

## v0.0.1 [2022-02-01]

_What's new?_

- New tables added

  - [ldap_group](https://hub.steampipe.io/plugins/turbot/ldap/tables/ldap_group)
  - [ldap_organizational_unit](https://hub.steampipe.io/plugins/turbot/ldap/tables/ldap_organizational_unit)
  - [ldap_user](https://hub.steampipe.io/plugins/turbot/ldap/tables/ldap_user)
