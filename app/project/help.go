package project

import (
	"strings"
)

var Helpers = func() map[string][]string {
	ret := map[string][]string{}
	add := func(k string, v string) {
		ret[k] = []string{v}
	}

	// project
	add("key", "The key of your project, lowercase letters only")
	add("name", "The project name is displayed in many places")
	add("icon", "The key of the SVG image used for your project's favicon and app icon")
	add("exec", "Your application's executable name, defaults to key")
	add("version", "Semantic version of the project")
	add("package", "Full Golang package, like github.com/org/key")
	add("args", "Arguments for your project when launched")
	add("port", "The TCP port used by your project's HTTP server")
	add("modules", "Project Forge modules used by this project")
	add("ignore", "Files ignored by your project")
	add("tags", "Tags for your project, not really used anywhere yet")

	// info
	add("org", "The Github organization responsible for this project")
	add("authorID", "The GitHub handle of the author of this project")
	add("authorName", "The full name of the author of this project")
	add("authorEmail", "The email address of the author of this project")
	add("team", "The team that owns this project")
	add("license", "The software license used by this project")
	add("homepage", "The main web page for this project")
	add("sourcecode", "The URL of this project's source repository")
	add("summary", "A one-line summary of this project")
	add("description", "A full multi-line description of this project")
	add("ci", "When to run CI")
	add("homebrew", "Override the URL to download Homebrew assets from")
	add("bundle", "App bundle used in iOS and macOS applications")
	add("signingIdentity", "Keychain identity to use for signing binaries")
	add("notarizeEmail", "Email address used for Apple notarization, defaults to owner")
	add("databaseEngine", "Primary database engine, used for DDL and seed data")
	add("slack", "Slack webhook for notifying after successful releases")
	add("channels", "Various communication channels used by this project")
	add("javaPackage", "The Java package used by the Android application")
	add("goVersion", "The version used for Go, defaults to \""+DefaultGoVersion+"\"")
	add("goBinary", "The binary used for Go, defaults to \"go\"")
	add("configVars", "Configuration variables for code generation")
	add("extraFiles", "Additional files to include in the builds")
	add("deployments", "Files that track deployment of this app")
	add("envvars", "additional environment variables used by this project")
	add("docs", "documentation resources for this project")

	return ret
}()

func clean(key string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(key, "-", "_"), "/", "_"), ".", "_")
}
