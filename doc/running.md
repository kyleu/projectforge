<!--- Content managed by Project Forge, see [projectforge.md] for details. -->
# Environment Variables

All variable names can be set in uppercase or lowercase

| Name                          | Type   | Description                                                                          |
|-------------------------------|--------|--------------------------------------------------------------------------------------|
| `app_display_name`            | string | replaces the UI title                                                                |
| `app_display_name_append`     | string | added to the end of the UI title                                                     |
| `app_nav_color_dark`          | string | sets the navigation color for users with dark mode, defaults to theme color          |
| `app_nav_color_light`         | string | sets the navigation color for users with light mode, defaults to theme color         |
| `logging_format`              | string | When set to `json`, forces the logging format                                        |
| `logging_level`               | string | minimum logging level to display, one of [`debug`, `info`, `warn`, `error`]          |
| `telemetry_disabled`          | bool   | when set, disables all telemetry                                                     |
| `telemetry_endpoint`          | string | address of OpenTelemetry collector (when enabled), defaults to `localhost:55681`     |
| `projectforge_encryption_key` | string | encryption key for web sessions, defaults to `projectforge_secret`, warns if missing |
