package file

import (
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

type Type struct {
	Key      string
	Suffixes []string
	Title    string
}

var dockerfileSuffixes = []string{"Dockerfile", "Dockerfile.debug", "Dockerfile.desktop", "Dockerfile.publish", "Dockerfile.release"}

var (
	TypeBatch        = Type{Key: "batch", Suffixes: []string{".bat"}, Title: "Windows batch file"}
	TypeCodeowners   = Type{Key: "codeowners", Suffixes: []string{"CODEOWNERS"}, Title: "Code owners source control file"}
	TypeConf         = Type{Key: "conf", Suffixes: []string{".conf"}, Title: "Configuration file"}
	TypeCSS          = Type{Key: "css", Suffixes: []string{util.ExtCSS}, Title: "Cascading Stylesheet"}
	TypeDocker       = Type{Key: "docker", Suffixes: dockerfileSuffixes, Title: "Docker build configuration"}
	TypeDockerIgnore = Type{Key: "dockerignore", Suffixes: []string{".dockerignore"}, Title: "Docker ignore configuration"}
	TypeEditorConfig = Type{Key: "editorconfig", Suffixes: []string{".editorconfig"}, Title: "IDE configuration and styling"}
	TypeEntitlements = Type{Key: "entitlements", Suffixes: []string{".entitlements"}, Title: "Apple entitlements file"}
	TypeEnv          = Type{Key: "env", Suffixes: []string{".env"}, Title: "Shell environment files"}
	TypeESLint       = Type{Key: "eslint", Suffixes: []string{".eslintrc"}, Title: "ESLint configuration file"}
	TypeGitIgnore    = Type{Key: "gitignore", Suffixes: []string{".gitignore"}, Title: "Git ignore file"}
	TypeGo           = Type{Key: "go", Suffixes: []string{util.ExtGo}, Title: "Go source"}
	TypeGoMod        = Type{Key: "gomod", Suffixes: []string{util.ExtMod}, Title: "Go module configuration"}
	TypeGradle       = Type{Key: "gradle", Suffixes: []string{".gradle"}, Title: "Gradle build or source file"}
	TypeGraphQL      = Type{Key: "graphql", Suffixes: []string{".graphql", ".graphqls"}, Title: "Query and schema files for GraphQL"}
	TypeHCL          = Type{Key: "hcl", Suffixes: []string{".hcl"}, Title: "HashiCorp Configuration Language"}
	TypeHTML         = Type{Key: "html", Suffixes: []string{util.ExtHTML, ".htm"}, Title: "HTML web page"}
	TypeIgnore       = Type{Key: "dockerignore", Suffixes: []string{".dockerignore"}, Title: "Docker ignore file"}
	TypeIcons        = Type{Key: "icons", Suffixes: []string{".icns"}, Title: "Apple icon collection"}
	TypeJavaScript   = Type{Key: "javascript", Suffixes: []string{util.ExtJS, ".javascript"}, Title: "JavaScript source"}
	TypeJSON         = Type{Key: "json", Suffixes: []string{util.ExtJSON}, Title: "JavaScript Object Notation"}
	TypeKotlin       = Type{Key: "kotlin", Suffixes: []string{".kt"}, Title: "Kotlin source file"}
	TypeMakefile     = Type{Key: "makefile", Suffixes: []string{"Makefile", "makefile"}, Title: "Build file"}
	TypeMarkdown     = Type{Key: "markdown", Suffixes: []string{util.ExtMarkdown, ".markdown"}, Title: "Markdown document"}
	TypePList        = Type{Key: "plist", Suffixes: []string{".plist"}, Title: "Apple configuration plist"}
	TypeProperties   = Type{Key: "properties", Suffixes: []string{".properties"}, Title: "Java properties file"}
	TypeProtobuf     = Type{Key: "protobuf", Suffixes: []string{".proto"}, Title: "Protobuf definition"}
	TypeShell        = Type{Key: "shell", Suffixes: []string{".sh"}, Title: "Shell script"}
	TypeSQL          = Type{Key: "sql", Suffixes: []string{util.ExtSQL}, Title: "SQL query file"}
	TypeSVG          = Type{Key: "svg", Suffixes: []string{util.ExtSVG}, Title: "Simple Vector Graphics file"}
	TypeSwift        = Type{Key: "swift", Suffixes: []string{".swift"}, Title: "Swift source file"}
	TypeTypeScript   = Type{Key: "typescript", Suffixes: []string{util.ExtTS}, Title: "TypeScript source"}
	TypeXML          = Type{Key: "xml", Suffixes: []string{".xml"}, Title: "XML document"}
	TypeYAML         = Type{Key: "yaml", Suffixes: []string{".yml", ".yaml"}, Title: "YAML configuration source"}
)

var AllTypes = []Type{
	TypeBatch, TypeCodeowners, TypeConf, TypeCSS, TypeDocker, TypeDockerIgnore, TypeEditorConfig, TypeEntitlements, TypeEnv,
	TypeESLint, TypeGitIgnore, TypeGo, TypeGoMod, TypeGradle, TypeGraphQL, TypeHCL, TypeHTML, TypeIcons, TypeIgnore,
	TypeJavaScript, TypeJSON, TypeKotlin, TypeMakefile, TypeMarkdown, TypePList, TypeProperties,
	TypeProtobuf, TypeShell, TypeSQL, TypeSVG, TypeSwift, TypeTypeScript, TypeXML, TypeYAML,
}

func TypeFromString(s string) Type {
	return lo.FindOrElse(AllTypes, errorType("No file type available with key ["+s+"]"), func(t Type) bool {
		return t.Key == s
	})
}

func errorType(msg string) Type {
	return Type{Key: util.KeyError, Title: "Error: " + msg}
}

func (t *Type) String() string {
	return t.Key
}

func (t *Type) MarshalJSON() ([]byte, error) {
	return util.ToJSONBytes(t.Key, false), nil
}

func (t *Type) UnmarshalJSON(data []byte) error {
	var s string
	if err := util.FromJSON(data, &s); err != nil {
		return err
	}
	x := TypeFromString(s)
	*t = x
	return nil
}

func getType(fn string) Type {
	for _, t := range AllTypes {
		for _, suf := range t.Suffixes {
			if strings.HasSuffix(fn, suf) {
				return t
			}
		}
	}
	return errorType("invalid: " + fn)
}
