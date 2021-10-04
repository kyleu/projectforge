package file

import (
	"strings"

	"github.com/kyleu/projectforge/app/util"
)

type Type struct {
	Key      string
	Suffixes []string
	Title    string
}

var (
	TypeConf       = Type{Key: "conf", Suffixes: []string{".conf"}, Title: "Configuration file"}
	TypeCSS        = Type{Key: "css", Suffixes: []string{".css"}, Title: "Cascading Stylesheet"}
	TypeDocker     = Type{Key: "docker", Suffixes: []string{"Dockerfile"}, Title: "Docker build configuration"}
	TypeGo         = Type{Key: "go", Suffixes: []string{".go"}, Title: "Go source"}
	TypeGoMod      = Type{Key: "gomod", Suffixes: []string{".mod"}, Title: "Go module configuration"}
	TypeHCL        = Type{Key: "hcl", Suffixes: []string{".hcl"}, Title: "HashiCorp Configuration Language"}
	TypeHTML       = Type{Key: "html", Suffixes: []string{".html", ".htm"}, Title: "HTML web page"}
	TypeIgnore     = Type{Key: "dockerignore", Suffixes: []string{".dockerignore"}, Title: "Docker ignore file"}
	TypeIcons      = Type{Key: "icons", Suffixes: []string{".icns"}, Title: "Apple icon collection"}
	TypeJavaScript = Type{Key: "javascript", Suffixes: []string{".js", ".javascript"}, Title: "JavaScript source"}
	TypeJSON       = Type{Key: "json", Suffixes: []string{".json"}, Title: "JavaScript Object Notation"}
	TypeMakefile   = Type{Key: "makefile", Suffixes: []string{"Makefile", "makefile"}, Title: "Build file"}
	TypeMarkdown   = Type{Key: "markdown", Suffixes: []string{".md", ".markdown"}, Title: "Markdown document"}
	TypePBXProject = Type{Key: "pbxproject", Suffixes: []string{".pbxproj"}, Title: "Apple project configuration"}
	TypePList      = Type{Key: "plist", Suffixes: []string{".plist"}, Title: "Apple configuration plist"}
	TypeShell      = Type{Key: "shell", Suffixes: []string{".sh"}, Title: "Shell script"}
	TypeSVG        = Type{Key: "svg", Suffixes: []string{".svg"}, Title: "Simple Vector Graphics file"}
	TypeSwift      = Type{Key: "swift", Suffixes: []string{".swift"}, Title: "Swift source code file"}
	TypeText       = Type{Key: "text", Suffixes: []string{".txt"}, Title: "Plaintext"}
	TypeTypeScript = Type{Key: "typescript", Suffixes: []string{".ts"}, Title: "TypeScript source"}
	TypeYAML       = Type{Key: "yaml", Suffixes: []string{".yml", ".yaml"}, Title: "YAML configuration source"}
)

var AllTypes = []Type{
	TypeConf, TypeCSS, TypeDocker, TypeGo, TypeGoMod, TypeHCL, TypeHTML, TypeIcons, TypeIgnore, TypeJavaScript, TypeJSON,
	TypeMakefile, TypeMarkdown, TypePBXProject, TypePList, TypeShell, TypeSVG, TypeSwift, TypeText, TypeTypeScript, TypeYAML,
}

func TypeFromString(s string) Type {
	for _, t := range AllTypes {
		if t.Key == s {
			return t
		}
	}
	return errorType("No file type available with key [" + s + "]")
}

func errorType(msg string) Type {
	return Type{Key: "error", Title: "Error: " + msg}
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
