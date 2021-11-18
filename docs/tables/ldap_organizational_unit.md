# Table: ldap_organizational_unit

An organizational unit contains users, computers, groups etc.

**Important notes:**

This table supports optional quals. Queries with optional quals in `where` clause are optimised to use ldap filters.

Optional quals are supported for the following columns:

- `filter` - Allows use of explicit query. Refer [LDAP filter language](https://ldap.com/ldap-filters/)
- `ou`
- `description`
- `created_on`
- `modified_on`

**Note:** This table supports an optional `filter` column to query results based on the LDAP [filter](https://ldap.com/ldap-filters/) language.

## Examples

### Basic info

```sql
select
  dn,
  ou,
  created_on
from
  ldap_organizational_unit;
```

### Get name and the person/group managing the organizational unit

```sql
select
  ou,
  created_on,
  managed_by
from
  ldap_organizational_unit;
```

### List organizational units that have been created in the last '30' days

```sql
select
  dn,
  ou,
  created_on
from
  ldap_organizational_unit
where
  created_on > current_timestamp - interval '30 days';
```

### List all organizational units that are critical system objects and cannot be deleted without replacement

```sql
select
  dn,
  ou,
  created_on,
  managed_by
from
  ldap_organizational_unit
where
  attributes->'isCriticalSystemObject' ? 'TRUE';
```
