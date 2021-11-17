# Table: ldap_user

A user is known as the customer or end-user.

**Important notes:**

This table supports optional quals. Queries with optional quals in `where` clause are optimised to use ldap filters.

Optional quals are supported for the following columns:

- `filter` - Allows use of explicit query. Refer [LDAP filter language](https://ldap.com/ldap-filters/)
- `changed`
- `cn`
- `created_on`
- `description`
- `disabled`
- `display_name`
- `given_name`
- `log_stream_name`
- `mail`
- `object_sid`
- `sam_account_name`
- `surname`
- `user_principal_name`

## Examples

### Basic info

```sql
select
  dn,
  cn,
  initials,
  created_on,
  mail,
  department,
  sam_account_name
from
  ldap_user limit 100;
```

### Get logon name, e-mail, and groups that each user is a member of

```sql
select
  user_principal_name,
  display_name,
  mail,
  jsonb_pretty(member_of) as user_groups
from
  ldap_user;
```

### List disabled users

```sql
select
  dn,
  sam_account_name,
  mail,
  object_sid
from
  ldap_user
where
  disabled;
```

### List users belonging to the 'Engineering' Department

```sql
select
  dn,
  sam_account_name,
  mail,
  department
from
  ldap_user
where
  department = 'Engineering';
```

### List users that have been created in the last '30' days

```sql
select
  dn,
  sam_account_name,
  mail,
  created_on
from
  ldap_user
where
  created_on > current_timestamp - interval '30 days';
```

### Get details of user 'Adelhard Frey' and the groups which he is a member of

```sql
select
  u.dn as userDn,
  u.mail as email,
  u.object_sid as userObjectSid,
  g.dn as groupDn,
  g.cn as groupName,
  g.object_sid as groupObjectSid
from
  ldap.ldap_user as u
cross join
  jsonb_array_elements_text(u.member_of) as groups
inner join
  ldap.ldap_group as g
on
  g.dn = groups
where
  u.cn = 'Adelhard Frey';
```

### List users whose name start with John using a filter

```sql
select
  dn,
  sam_account_name,
  mail,
  created_on
from
  ldap_user
where
  filter = '(cn=Adam*)';
```
