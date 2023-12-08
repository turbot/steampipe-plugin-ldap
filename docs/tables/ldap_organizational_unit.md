---
title: "Steampipe Table: ldap_organizational_unit - Query LDAP Organizational Units using SQL"
description: "Allows users to query LDAP Organizational Units, providing insights into the structure and hierarchy of LDAP directories."
---

# Table: ldap_organizational_unit - Query LDAP Organizational Units using SQL

An LDAP Organizational Unit represents a container that can hold users, groups, and other organizational units within an LDAP directory. It is a crucial component in creating an organized structure in LDAP directories, allowing for efficient user and resource management. Organizational units can be nested within each other, creating a hierarchical structure that reflects the organization's structure.

## Table Usage Guide

The `ldap_organizational_unit` table provides insights into the structure and hierarchy of LDAP directories. As a system administrator, explore details about organizational units through this table, including their distinguished names, attributes, and associated metadata. Utilize it to understand the organization's structure, manage resources efficiently, and implement access control effectively.

## Examples

### Basic info
Explore the organizational units within your network, including when they were created and who manages them. This can help you understand your network's structure and identify areas for potential reorganization or management changes.

```sql+postgres
select
  dn,
  ou,
  when_created,
  managed_by
from
  ldap_organizational_unit;
```

```sql+sqlite
select
  dn,
  ou,
  when_created,
  managed_by
from
  ldap_organizational_unit;
```

### List organizational units that have been created in the last 30 days
Discover the segments that have been recently added to your organization within the past month. This can help keep track of organizational growth and changes.

```sql+postgres
select
  dn,
  ou,
  when_created
from
  ldap_organizational_unit
where
  when_created > current_timestamp - interval '30 days';
```

```sql+sqlite
select
  dn,
  ou,
  when_created
from
  ldap_organizational_unit
where
  when_created > datetime('now','-30 days');
```

## Filter Examples

### List organizational units that are critical system objects
Determine the areas in which organizational units are deemed as critical system objects. This can be used to identify key system components that require special attention or stricter security measures.

```sql+postgres
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

```sql+sqlite
select
  dn,
  ou,
  managed_by,
  json_extract(attributes, '$.isCriticalSystemObject') as is_critical_system_object
from
  ldap_organizational_unit
where
  json_extract(attributes, '$.isCriticalSystemObject') = 'TRUE';
```