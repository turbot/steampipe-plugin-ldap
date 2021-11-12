# Table: ldap_group

A group is a collection of digital identities i.e. users, groups etc.

**Note:** This table supports an optional `filter` column to query results based on the LDAP [filter](https://ldap.com/ldap-filters/) language.

## Examples

### Basic info

```sql
select
  dn,
  cn,
  ou,
  created_on,
  sam_account_name
from
  ldap_group;
```

### Get logon name, organizational unit, and groups that the group is a member of

```sql
select
  sam_account_name,
  ou,
  jsonb_pretty(member_of) as group_groups
from
  ldap_group;
```

### List groups that have been created in the last '30' days

```sql
select
  dn,
  sam_account_name,
  created_on
from
  ldap_group
where
  created_on > current_timestamp - interval '30 days';
```

### Get details of group 'Database' and the groups which it is a member of

```sql
select
  g.dn as groupDn,
  g.ou as groupOu,
  g.object_sid as groupObjectSid,
  mg.dn as memberOfGroupDn,
  mg.cn as memberOfGroupName,
  mg.object_sid as memberOfGroupObjectSid
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

### Get all members of a particular group

```sql
select
  jsonb_pretty(attributes->'member') as members
from
  ldap_group
where
  cn = 'Database';
```