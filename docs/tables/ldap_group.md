# Table: ldap_group

A group is a collection of digital identities, e.g., users, groups.

**Note:**

- This table supports optional quals. Queries with optional quals in a `where` clause are optimised to use LDAP search filters.
- If `filter` is provided, other optional quals will not be used when searching.
- Optional quals are supported for the following columns:
  - `cn`
  - `description`
  - `filter` - Allows use of an explicit filter. Please refer to [LDAP filter language](https://ldap.com/ldap-filters/).
  - `object_sid`
  - `sam_account_name`
  - `when_changed`
  - `when_created`

## Examples

### Basic info

```sql
select
  dn,
  cn,
  ou,
  when_created,
  sam_account_name
from
  ldap_group;
```

### List all members for each group

```sql
select
  jsonb_pretty(attributes -> 'member') as members
from
  ldap_group;
```
### List groups that have been created in the last 30 days

```sql
select
  dn,
  sam_account_name,
  when_created
from
  ldap_group
where
  when_created > current_timestamp - interval '30 days';
```

### Get details for groups the group 'Database' is a member of

```sql
select
  g.dn as group_nd,
  g.ou as group_ou,
  g.object_sid as group_object_sid,
  mg.dn as parent_group_dn,
  mg.cn as parent_group_name,
  mg.object_sid as parent_group_object_sid
from
  ldap.ldap_group as g
cross join
  jsonb_array_elements_text(g.member_of) as groups
inner join
  ldap.ldap_group as mg
on
  mg.dn = groups
where
  g.cn = 'Database';
```

## Filter Examples

### List groups that "Bob Smith" is a member of

```sql
select
  dn,
  ou,
  description,
  when_created
from
  ldap_group
where
  filter = '(member=CN=Bob Smith,OU=Devs,OU=SP,DC=sp,DC=turbot,DC=com)';
```
