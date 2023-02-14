-- $PF_GENERATE_ONCE$
-- {% func CreateDatabase() %}
create role "{{{ .Key }}}" with login password '{{{ .Key }}}';

create database "{{{ .Key }}}";
alter database "{{{ .Key }}}" set timezone to 'utc';
grant all privileges on database "{{{ .Key }}}" to "{{{ .Key }}}";
-- {% endfunc %}
