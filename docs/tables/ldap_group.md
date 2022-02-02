# Table: ldap_group

A group is a collection of digital identities i.e. users, groups etc.

**Important notes:**

This table supports optional quals. Queries with optional quals in `where` clause are optimised to use ldap filters.

Optional quals are supported for the following columns:

- `cn`
- `description`
- `filter` - Allows use of explicit query. Refer [LDAP filter language](https://ldap.com/ldap-filters/)
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
  when_created
from
  ldap_group
where
  when_created > current_timestamp - interval '30 days';
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

### List groups a user is member of using user `dn` in `filter`

```sql
select
  dn,
  ou,
  description,
  when_created
from
  ldap_group
where
  filter = '(member:1.2.840.113556.1.4.1941:=CN=Ljiljana Rausch,OU=Mods,OU=VASHI,DC=vashi,DC=turbot,DC=com)';
```
