# Table: ldap_organizational_unit

An organizational unit contains users, computers, groups etc.

**Note:** This table supports an optional `filter` column to query results based on the LDAP [filter](https://ldap.com/ldap-filters/) language.

## Examples

### Basic info

```sql
select
  dn,
  ou,
  created
from
  ldap_organizational_unit;
```

### Get name and the person/group managing the organizational unit

```sql
select
  ou,
  created,
  managed_by
from
  ldap_organizational_unit;
```

### List organizational units that have been created in the last '30' days

```sql
select
  dn,
  ou,
  created
from
  ldap_organizational_unit
where
  created > current_timestamp - interval '30 days';
```

### List all organizational units that are critical system objects and cannot be deleted without replacement

```sql
select
  dn,
  ou,
  created,
  managed_by
from
  ldap_organizational_unit
where
  attributes->'isCriticalSystemObject' ? 'TRUE';
```
