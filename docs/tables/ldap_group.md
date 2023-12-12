---
title: "Steampipe Table: ldap_group - Query LDAP Groups using SQL"
description: "Allows users to query LDAP Groups, specifically the group DN, CN, and member details, providing insights into group configurations and membership."
---

# Table: ldap_group - Query LDAP Groups using SQL

Lightweight Directory Access Protocol (LDAP) is a protocol used to access directory listings within Active Directory (AD), OpenLDAP, and other directory systems. It allows users to access and manage a variety of information, including user profiles, groups, and network information. The LDAP service in provides a way to connect to and manage your LDAP directories.

## Table Usage Guide

The `ldap_group` table provides insights into LDAP groups within LDAP service. As a Systems Administrator, explore group-specific details through this table, including distinguished names (DN), common names (CN), and member details. Utilize it to uncover information about groups, such as those with specific members, the hierarchical relationships between groups, and the verification of group configurations.

**Important Notes**

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
Explore which groups exist within your LDAP directory by identifying their distinguished names, common names, and organizational units. This can be particularly useful for auditing purposes, helping to ensure that all groups are accounted for and appropriately organized.

```sql+postgres
select
  dn,
  cn,
  ou,
  when_created,
  sam_account_name
from
  ldap_group;
```

```sql+sqlite
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
Explore which members belong to each group in your LDAP directory. This helps in managing access controls and user permissions effectively.

```sql+postgres
select
  jsonb_pretty(attributes -> 'member') as members
from
  ldap_group;
```

```sql+sqlite
select
  attributes as members
from
  ldap_group;
```
### List groups that have been created in the last 30 days
Discover the groups that have been established in the recent 30 days. This query is useful for monitoring the creation of new groups and maintaining an up-to-date overview of your system's group structure.

```sql+postgres
select
  dn,
  sam_account_name,
  when_created
from
  ldap_group
where
  when_created > current_timestamp - interval '30 days';
```

```sql+sqlite
select
  dn,
  sam_account_name,
  when_created
from
  ldap_group
where
  when_created > datetime('now', '-30 days');
```

### Get details for groups the group 'Database' is a member of
Explore the hierarchical relationships within your group structures, specifically identifying the parent groups to which your 'Database' group belongs. This is useful for managing access controls and understanding how permissions are inherited within your organization.

```sql+postgres
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

```sql+sqlite
select
  g.dn as group_nd,
  g.ou as group_ou,
  g.object_sid as group_object_sid,
  mg.dn as parent_group_dn,
  mg.cn as parent_group_name,
  mg.object_sid as parent_group_object_sid
from
  ldap.ldap_group as g,
  json_each(g.member_of) as groups
inner join
  ldap.ldap_group as mg
on
  mg.dn = groups.value
where
  g.cn = 'Database';
```

## Filter Examples

### List groups that "Bob Smith" is a member of
Discover the groups that a specific user, such as Bob Smith, belongs to, providing insights into the user's roles and permissions within the organization. This can be useful for auditing user access and ensuring appropriate security measures.

```sql+postgres
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

```sql+sqlite
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