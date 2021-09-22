# Table: ldap_user

A user is known as the customer or end-user.

Note: This table supports an optional `filter` column to query results based on the LDAP [filter](https://ldap.com/ldap-filters/) language.

## Examples

### Basic info

```sql
select
  dn,
  cn,
  initials,
  created,
  mail,
  department,
  sam_account_name
from
  ldap_user;
```

### Get logon name, e-mail, and groups that each user is a member of

```sql
select
  sam_account_name,
  mail,
  jsonb_pretty(member_of) as user_groups
from
  ldap_user;
```

### Get details of users whose account is disabled

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
  created
from
  ldap_user
where
  created > current_timestamp - interval '30 days';
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
  created
from
  ldap_user
where
  filter = '(cn=Adam*)';
```
