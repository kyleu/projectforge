-- {% func QueueCreateTable(table string) %}
create table "{%s table %}" (
  id text primary key,
  topic text not null,
  param text not null,
  retries integer not null default 0,
  max_age text not null default (strftime('%Y-%m-%dT%H:%M:%fZ')),
  created text not null default (strftime('%Y-%m-%dT%H:%M:%fZ')),
  updated text not null default (strftime('%Y-%m-%dT%H:%M:%fZ'))
) strict;

create trigger {%s table %}_updated_timestamp after update on {%s table %} begin
    update {%s table %} set updated = strftime('%Y-%m-%dT%H:%M:%fZ') where id = old.id;
end;

create index {%s table %}_topic_created_idx on {%s table %} (topic, created);
-- {% endfunc %}

-- {% func QueueCount(table string) %}
select count(*) as "x" from "{%s table %}" where "topic" = $1;
-- {% endfunc %}

-- {% func QueueRead(table string) %}
update {%s table %} set "max_age" = $1, "retries" = "retries" + 1 where "id" = (
  select "id" from "{%s table %}" where "topic" = $2 and "max_age" < $3 and "retries" < $4 order by "created" limit 1
) returning "id", "topic", "param", "retries";
-- {% endfunc %}

-- {% func QueueWrite(table string) %}
insert into "{%s table %}" ("id", "topic", "param") values ($1, $2, $3);
-- {% endfunc %}

-- {% func QueueDelete(table string) %}
delete from "{%s table %}" where "id" = $1;
-- {% endfunc %}
