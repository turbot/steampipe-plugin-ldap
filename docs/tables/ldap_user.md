# Table: ldap_user

A user is known as the customer or end-user.

**Notes:**

- This table supports optional quals. Queries with optional quals in a `where` clause are optimised to use LDAP search filters.
- If `filter` is provided, other optional quals will not be used when searching.
- Optional quals are supported for the following columns:
  - `cn`
  - `department`
  - `description`
  - `disabled`
  - `display_name`
  - `filter` - Allows use of an explicit filter. Please refer to [LDAP filter language](https://ldap.com/ldap-filters/).
  - `given_name`
  - `mail`
  - `object_sid`
  - `sam_account_name`
  - `surname`
  - `user_principal_name`
  - `when_created`
  - `when_changed`

## Examples

### Basic info

```sql
select
  dn,
  cn,
  when_created,
  mail,
  department,
  sam_account_name,
from
  ldap_user limit 100;
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

### List users in the 'Engineering' department

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

### List users that have been created in the last 30 days

```sql
select
  dn,
  sam_account_name,
  mail,
  when_created
from
  ldap_user
where
  when_created > current_timestamp - interval '30 days';
```

### Get details for groups the user 'Bob Smith' is a member of

```sql
select
  u.dn as userDn,
  u.mail as email,
  u.object_sid as user_object_sid,
  g.dn as group_dn,
  g.cn as group_name,
  g.object_sid as group_object_sid
from
  ldap.ldap_user as u
cross join
  jsonb_array_elements_text(u.member_of) as groups
inner join
  ldap.ldap_group as g
on
  g.dn = groups
where
  u.cn = 'Bob Smith';
```

## Filter Examples

### List users whose names start with "Adam"

```sql
select
  dn,
  sam_account_name,
  mail,
  when_created
from
  ldap_user
where
  filter = '(cn=Adam*)';
```

### List members of a group filtering by the group's DN

```sql
select
  cn,
  display_name,
  when_created,
  user_principal_name,
  ou,
  given_name
from
  ldap_user
where
  filter = '(memberof=CN=Devs,OU=Steampipe,OU=SP,DC=sp,DC=turbot,DC=com)';
```
