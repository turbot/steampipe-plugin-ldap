## v1.1.1 [2025-04-18]

_Bug fixes_

- Fixed Linux AMD64 plugin build failures for `Postgres 14 FDW`, `Postgres 15 FDW`, and `SQLite Extension` by upgrading GitHub Actions runners from `ubuntu-20.04` to `ubuntu-22.04`.

## v1.1.0 [2025-04-17]

_Dependencies_

- Recompiled plugin with Go version `1.23.1`. ([#41](https://github.com/turbot/steampipe-plugin-ldap/pull/41))
- Recompiled plugin with [steampipe-plugin-sdk v5.11.5](https://github.com/turbot/steampipe-plugin-sdk/blob/v5.11.5/CHANGELOG.md#v5115-2025-03-31) that addresses critical and high vulnerabilities in dependent packages. ([#41](https://github.com/turbot/steampipe-plugin-ldap/pull/41))

## v1.0.0 [2024-10-22]

There are no significant changes in this plugin version; it has been released to align with [Steampipe's v1.0.0](https://steampipe.io/changelog/steampipe-cli-v1-0-0) release. This plugin adheres to [semantic versioning](https://semver.org/#semantic-versioning-specification-semver), ensuring backward compatibility within each major version.

_Dependencies_

- Recompiled plugin with Go version `1.22`. ([#39](https://github.com/turbot/steampipe-plugin-ldap/pull/39))
- Recompiled plugin with [steampipe-plugin-sdk v5.10.4](https://github.com/turbot/steampipe-plugin-sdk/blob/develop/CHANGELOG.md#v5104-2024-08-29) that fixes logging in the plugin export tool. ([#39](https://github.com/turbot/steampipe-plugin-ldap/pull/39))

## v0.5.1 [2023-12-12]

_Bug fixes_

- Fixed the missing optional tag on the `attributes` connection config parameter.

## v0.5.0 [2023-12-12]

_What's new?_

- The plugin can now be downloaded and used with the [Steampipe CLI](https://steampipe.io/docs), as a [Postgres FDW](https://steampipe.io/docs/steampipe_postgres/overview), as a [SQLite extension](https://steampipe.io/docs//steampipe_sqlite/overview) and as a standalone [exporter](https://steampipe.io/docs/steampipe_export/overview). ([#34](https://github.com/turbot/steampipe-plugin-ldap/pull/34))
- The table docs have been updated to provide corresponding example queries for Postgres FDW and SQLite extension. ([#34](https://github.com/turbot/steampipe-plugin-ldap/pull/34))
- Docs license updated to match Steampipe [CC BY-NC-ND license](https://github.com/turbot/steampipe-plugin-ldap/blob/main/docs/LICENSE). ([#34](https://github.com/turbot/steampipe-plugin-ldap/pull/34))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.8.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v580-2023-12-11) that includes plugin server encapsulation for in-process and GRPC usage, adding Steampipe Plugin SDK version to `_ctx` column, and fixing connection and potential divide-by-zero bugs. ([#33](https://github.com/turbot/steampipe-plugin-ldap/pull/33))

## v0.4.1 [2023-10-05]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.6.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v562-2023-10-03) which prevents nil pointer reference errors for implicit hydrate configs. ([#26](https://github.com/turbot/steampipe-plugin-ldap/pull/26))

## v0.4.0 [2023-10-02]

_Dependencies_

- Upgraded to [steampipe-plugin-sdk v5.6.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v561-2023-09-29) with support for rate limiters. ([#24](https://github.com/turbot/steampipe-plugin-ldap/pull/24))
- Recompiled plugin with Go version `1.21`. ([#24](https://github.com/turbot/steampipe-plugin-ldap/pull/24))

## v0.3.0 [2023-04-11]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.3.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v530-2023-03-16) which includes fixes for query cache pending item mechanism and aggregator connections not working for dynamic tables. ([#16](https://github.com/turbot/steampipe-plugin-ldap/pull/16))

## v0.2.0 [2022-09-26]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.7](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v417-2022-09-08) which includes several caching and memory management improvements. ([#13](https://github.com/turbot/steampipe-plugin-ldap/pull/13))
- Recompiled plugin with Go version `1.19`. ([#13](https://github.com/turbot/steampipe-plugin-ldap/pull/13))

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
