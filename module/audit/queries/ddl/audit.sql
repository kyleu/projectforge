-- {% func AuditDrop() %}
drop table if exists "audit_record";
drop table if exists "audit";
-- {% endfunc %}

-- {% func AuditCreate() %}
create table if not exists "audit" (
  "id" uuid not null,
  "app" text not null,
  "act" text not null,
  "client" text not null,
  "server" text not null,
  "user" text not null,
  "metadata" jsonb not null,
  "message" text not null,
  "started" timestamp not null default now(),
  "completed" timestamp not null default now(),
  primary key ("id")
);

create index if not exists "audit__act" on "audit" using btree ("act" asc nulls last);
create index if not exists "audit__app" on "audit" using btree ("app" asc nulls last);
create index if not exists "audit__client" on "audit" using btree ("client" asc nulls last);
create index if not exists "audit__server" on "audit" using btree ("server" asc nulls last);
create index if not exists "audit__user_id" on "audit" using btree ("user" asc nulls last);

create table if not exists "audit_record" (
  "id" uuid not null,
  "audit_id" uuid not null,
  "t" text not null,
  "pk" text not null,
  "changes" jsonb not null,
  "occurred" timestamp not null default now(),
  foreign key ("audit_id") references "audit" ("id"),
  primary key ("id")
);

create index if not exists "audit_record__t" on "audit_record" using btree ("t" asc nulls last);
create index if not exists "audit_record__pk" on "audit_record" using btree ("pk" asc nulls last);
create index if not exists "audit_record__changes" on "audit_record" using gin ("changes");

create index if not exists "audit_record__audit_id_idx" on "audit_record" ("audit_id");
-- {% endfunc %}
