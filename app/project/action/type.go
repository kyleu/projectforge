package action

import (
	"projectforge.dev/projectforge/app/util"
)

type Type struct {
	Key         string
	Title       string
	Icon        string
	Description string
	Hidden      bool
}

var (
	TypeAudit    = Type{Key: "audit", Title: "Audit", Icon: "scale", Description: "Audits the project files, detecting invalid files and empty folders"}
	TypeBuild    = Type{Key: "build", Title: "Build", Icon: "hammer", Description: "Builds the project, many options available"}
	TypeCreate   = Type{Key: "create", Title: "Create", Icon: "folder-plus", Description: "Creates a new project"}
	TypeDebug    = Type{Key: "debug", Title: "Debug", Icon: "bug", Description: "Dumps information about the project"}
	TypeDoctor   = Type{Key: "doctor", Title: "Doctor", Icon: "first-aid", Description: "Makes sure your machine has the required dependencies"}
	TypeGenerate = Type{Key: "generate", Title: "Generate", Icon: "forward", Description: "Applies pending changes to files as required"}
	TypePreview  = Type{Key: "preview", Title: "Preview", Icon: "play", Description: "Shows what would happen if you generate"}
	TypeRules    = Type{Key: "rules", Title: "Rules", Icon: "play", Description: "Apply export rules from json file located at ./rules.json"}
	TypeSVG      = Type{Key: "svg", Title: "SVG", Icon: "icons", Description: "Builds the project's SVG files"}
	TypeTest     = Type{Key: "test", Title: "Test", Icon: "wrench", Description: "Runs internal tests, you probably don't want this", Hidden: true}
)

var (
	AllTypes     = []Type{TypeAudit, TypeBuild, TypeCreate, TypeDebug, TypeDoctor, TypeGenerate, TypePreview, TypeRules, TypeSVG, TypeTest}
	ProjectTypes = []Type{TypePreview, TypeGenerate, TypeAudit, TypeBuild}
)

func TypeFromString(s string) Type {
	for _, t := range AllTypes {
		if t.Key == s {
			return t
		}
	}
	return Type{Key: s, Title: "Error", Icon: "star", Description: "No action type available with key [" + s + "]"}
}

func (t *Type) String() string {
	return t.Key
}

func (t *Type) Breadcrumb() string {
	return t.Title + "**" + t.Icon
}

func (t *Type) Matches(x Type) bool {
	return t.Key == x.Key
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

func (t *Type) Expensive(cfg util.ValueMap) bool {
	if t.Matches(TypeBuild) {
		switch cfg.GetStringOpt("phase") {
		case buildDeps.Key, buildLint.Key, buildFull.Key:
			return true
		}
	}
	return false
}
