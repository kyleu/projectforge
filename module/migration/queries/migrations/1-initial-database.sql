-- $PF_GENERATE_ONCE$
{{{ if .HasExport }}}-- {% import "{{{ .Package }}}/queries/ddl" %}
-- {% import "{{{ .Package }}}/queries/seeddata" %}

-- {% func Migration1InitialDatabase(debug bool) %}

-- {%= ddl.CreateAll() %}
-- {%= seeddata.SeedDataAll() %}

-- {% endfunc %}{{{ else }}}
-- {% func Migration1InitialDatabase(debug bool) %}

select 1;

-- {% endfunc %}{{{ end }}}
