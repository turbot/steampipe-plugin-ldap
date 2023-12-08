---
title: "Steampipe Table: ldap_user - Query LDAP Users using SQL"
description: "Allows users to query LDAP Users, specifically user attributes like distinguished name (DN), common name (CN), and user email, providing insights into user information and attributes."
---

# Table: ldap_user - Query LDAP Users using SQL

The Lightweight Directory Access Protocol (LDAP) is a protocol designed to manage and access distributed directory information services over an Internet Protocol (IP) network. LDAP is used to look up encryption certificates, pointers to printers and other services on a network, and provide 'single sign-on' where one password for a user is shared between many services. An LDAP user is an entry or record within the LDAP directory that represents a single user with attributes such as common name, distinguished name, and email address.

## Table Usage Guide

The `ldap_user` table provides insights into user entries within the LDAP directory. As an IT administrator, explore user-specific details through this table, including distinguished names, common names, and email addresses. Utilize it to uncover information about users, such as their specific attributes, the hierarchical structure of the directory, and the relationships between different entries.

## Examples

### Basic info
Determine the areas in which certain user details were created, including email and department specifics, to gain insights into your user base. This can help in assessing the elements within your organization for better user management.

```sql+postgres
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

```sql+sqlite
select
  dn,
  cn,
  when_created,
  mail,
  department,
  sam_account_name
from
  ldap_user limit 100;
```

### List disabled users
Identify instances where user accounts have been disabled to ensure proper access control and maintain the integrity of your system.

```sql+postgres
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

```sql+sqlite
select
  dn,
  sam_account_name,
  mail,
  object_sid
from
  ldap_user
where
  disabled = 1;
```

### List users in the 'Engineering' department
Identify individuals who are part of the Engineering team. This is useful for understanding team composition and for reaching out to specific team members.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that include users who were added recently. This can help in monitoring user growth and understanding the rate at which new users are being added to your system.

```sql+postgres
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

```sql+sqlite
select
  dn,
  sam_account_name,
  mail,
  when_created
from
  ldap_user
where
  when_created > datetime('now', '-30 days');
```

### Get details for groups the user 'Bob Smith' is a member of
Determine the groups that a specific user, such as 'Bob Smith', is associated with. This can be useful in managing user permissions and understanding user roles within a system.

```sql+postgres
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

```sql+sqlite
select
  u.dn as userDn,
  u.mail as email,
  u.object_sid as user_object_sid,
  g.dn as group_dn,
  g.cn as group_name,
  g.object_sid as group_object_sid
from
  ldap_user as u,
  json_each(u.member_of) as groups
inner join
  ldap_group as g
on
  g.dn = groups.value
where
  u.cn = 'Bob Smith';
```

## Filter Examples

### List users whose names start with "Adam"
Explore which users have names starting with 'Adam' to quickly locate their account details and creation dates, useful for user management and auditing purposes.

```sql+postgres
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

```sql+sqlite
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
Identify the members of a specific group in your network, including their display names and principal user names. This can help you understand who has access to certain resources and when they were added to the group.

```sql+postgres
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

```sql+sqlite
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