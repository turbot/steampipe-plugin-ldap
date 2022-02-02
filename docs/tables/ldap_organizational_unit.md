# Table: ldap_organizational_unit

An organizational unit contains users, computers, groups, and other objects.

**Note:**

- This table supports optional quals. Queries with optional quals in a `where` clause are optimised to use LDAP search filters.
- If `filter` is provided, other optional quals will not be used when searching.
- Optional quals are supported for the following columns:
  - `description`
  - `filter` - Allows use of an explicit filter. Please refer to [LDAP filter language](https://ldap.com/ldap-filters/).
  - `ou`
  - `when_changed`
  - `when_created`

## Examples

### Basic info

```sql
select
  dn,
  ou,
  when_created,
  managed_by
from
  ldap_organizational_unit;
```

### List organizational units that have been created in the last 30 days

```sql
select
  dn,
  ou,
  when_created
from
  ldap_organizational_unit
where
  when_created > current_timestamp - interval '30 days';
```

## Filter Examples

### List organizational units that are critical system objects

```sql
select
  dn,
  ou,
  managed_by,
  attributes -> 'isCriticalSystemObject' as is_critical_system_object
from
  ldap_organizational_unit
where
  filter = '(isCriticalSystemObject=TRUE)';
```
