-- {% func SizeInfo(dbType string) %}
-- {% switch dbType %}{{{ if .PostgreSQL }}}
-- {% case "postgres" %}
with recursive
  pg_inherit(inhrelid, inhparent) as (
    select inhrelid, inhparent
    from pg_inherits
    union
    select child.inhrelid, parent.inhparent
    from pg_inherit child, pg_inherits parent
    where child.inhparent = parent.inhrelid
  ),
  pg_inherit_short as (
    select *
    from pg_inherit
    where inhparent not in (select inhrelid from pg_inherit)
  )

select
  table_schema,
  table_name,
  row_estimate,
  total_bytes as "total",
  pg_size_pretty(total_bytes) as "total_pretty",
  case when index_bytes is null then 0 else index_bytes end as "index",
  pg_size_pretty(case when index_bytes is null then 0 else index_bytes end) as "index_pretty",
  case when toast_bytes is null then 0 else toast_bytes end as "toast",
  pg_size_pretty(case when toast_bytes is null then 0 else toast_bytes end) as "toast_pretty",
  case when table_bytes is null then 0 else table_bytes end as "table",
  pg_size_pretty(case when table_bytes is null then 0 else table_bytes end) as "table_pretty"
from (
  select *, total_bytes - index_bytes - coalesce(toast_bytes, 0) as table_bytes
  from (
    select c.oid, nspname as table_schema, relname as table_name,
      sum(c.reltuples) over (partition by parent) as row_estimate, sum(pg_total_relation_size(c.oid)) over (partition by parent) as total_bytes,
      sum(pg_indexes_size(c.oid)) over (partition by parent) as index_bytes,
      sum(pg_total_relation_size(reltoastrelid)) over (partition by parent) as toast_bytes, parent
    from (
      select pg_class.oid, reltuples, relname, relnamespace, pg_class.reltoastrelid, coalesce(inhparent, pg_class.oid) parent
      from pg_class left join pg_inherit_short on inhrelid = oid
      where relkind in ('r', 'p')
    ) c left join pg_namespace n on n.oid = c.relnamespace
    where nspname != 'pg_catalog' and nspname != 'information_schema'
  ) a
  where oid = parent
) a
order by total_bytes desc;{{{ end }}}{{{ if .SQLite }}}
-- {% case "sqlite" %}
select
  'default' as "table_schema",
  "name" as "table_name",
  0 as "row_estimate",
  0 as "total",
  '' as "total_pretty",
  0 as "index",
  '' as "index_pretty",
  0 as "toast",
  '' as "toast_pretty",
  0 as "table",
  '' as "table_pretty"
from "sqlite_master"
where "type" = 'table'
order by "table_name";{{{ end }}}{{{ if .SQLServer }}}
-- {% case "sqlserver" %}
with sizeinfo as (
  select
    s.name as "table_schema",
    t.name as "table_name",
    coalesce(sum(case when i.index_id < 2 then p.rows else 0 end), 0) as "row_estimate",
    coalesce(sum(a.used_pages) * 8 * 1024, 0) as "total",
    coalesce(sum(case when i.index_id >= 2 then a.used_pages else 0 end) * 8 * 1024, 0) as "index",
    0 as "toast",
    coalesce(sum(case when i.index_id < 2 then a.used_pages else 0 end) * 8 * 1024, 0) as "table"
  from sys.tables t
    join sys.schemas s on t.schema_id = s.schema_id
    join sys.indexes i on t.object_id = i.object_id
    join sys.partitions p on i.object_id = p.object_id and i.index_id = p.index_id
    join sys.allocation_units a on a.container_id = p.hobt_id or a.container_id = p.partition_id
  where t.is_ms_shipped = 0
  group by s.name, t.name
)
select
  "table_schema",
  "table_name",
  cast("row_estimate" as varchar(32)) as "row_estimate",
  "total",
  case
    when "total" >= 1099511627776 then concat(cast("total" / 1099511627776.0 as decimal(18,2)), ' TB')
    when "total" >= 1073741824 then concat(cast("total" / 1073741824.0 as decimal(18,2)), ' GB')
    when "total" >= 1048576 then concat(cast("total" / 1048576.0 as decimal(18,2)), ' MB')
    when "total" >= 1024 then concat(cast("total" / 1024.0 as decimal(18,2)), ' KB')
    else concat("total", ' B')
  end as "total_pretty",
  "index",
  case
    when "index" >= 1099511627776 then concat(cast("index" / 1099511627776.0 as decimal(18,2)), ' TB')
    when "index" >= 1073741824 then concat(cast("index" / 1073741824.0 as decimal(18,2)), ' GB')
    when "index" >= 1048576 then concat(cast("index" / 1048576.0 as decimal(18,2)), ' MB')
    when "index" >= 1024 then concat(cast("index" / 1024.0 as decimal(18,2)), ' KB')
    else concat("index", ' B')
  end as "index_pretty",
  "toast",
  '' as "toast_pretty",
  "table",
  case
    when "table" >= 1099511627776 then concat(cast("table" / 1099511627776.0 as decimal(18,2)), ' TB')
    when "table" >= 1073741824 then concat(cast("table" / 1073741824.0 as decimal(18,2)), ' GB')
    when "table" >= 1048576 then concat(cast("table" / 1048576.0 as decimal(18,2)), ' MB')
    when "table" >= 1024 then concat(cast("table" / 1024.0 as decimal(18,2)), ' KB')
    else concat("table", ' B')
  end as "table_pretty"
from sizeinfo
order by "total" desc;{{{ end }}}
-- {% default %}
select 'unhandled database type [{%s dbType %}]';
-- {% endswitch %}
-- {% endfunc %}
