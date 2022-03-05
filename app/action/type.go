package action

import (
	"github.com/kyleu/projectforge/app/util"
)

type Type struct {
	Key         string
	Title       string
	Icon        string
	Description string
	Hidden      bool
}

var (
	TypeAudit   = Type{Key: "audit", Title: "Audit", Icon: "star", Description: "Audits the project files"}
	TypeBuild   = Type{Key: "build", Title: "Build", Icon: "star", Description: "Builds the project"}
	TypeCreate  = Type{Key: "create", Title: "Create", Icon: "star", Description: "Creates a new project"}
	TypeDebug   = Type{Key: "debug", Title: "Debug", Icon: "star", Description: "Dumps a ton of information about the project"}
	TypeDoctor  = Type{Key: "doctor", Title: "Doctor", Icon: "star", Description: "Makes sure your machine has the required dependencies"}
	TypeMerge   = Type{Key: "merge", Title: "Merge", Icon: "star", Description: "Merges changed files as required"}
	TypePreview = Type{Key: "preview", Title: "Preview", Icon: "star", Description: "Show what would happen if you did merge"}
	TypeSlam    = Type{Key: "slam", Title: "Slam", Icon: "star", Description: "Slams all files to the target, ignoring changes"}
	TypeSVG     = Type{Key: "svg", Title: "SVG", Icon: "star", Description: "Builds the project's SVG files"}
	TypeTest    = Type{Key: "test", Title: "Test", Icon: "star", Description: "Runs internal tests, you probably don't want this", Hidden: true}
)

var (
	AllTypes     = []Type{TypeAudit, TypeBuild, TypeCreate, TypeDebug, TypeDoctor, TypeMerge, TypePreview, TypeSlam, TypeSVG, TypeTest}
	ProjectTypes = []Type{TypePreview, TypeMerge, TypeSlam, TypeAudit}
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
