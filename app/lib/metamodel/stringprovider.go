package metamodel

import "projectforge.dev/projectforge/app/util"

type StringProvider interface {
	PackageName() string
	Camel() string
	CamelLower() string
	CamelPlural() string
	Proper() string
	ProperPlural() string
	Kebab() string
	FirstLetter() string
	PackageWithGroup(prefix string) string
	GroupLen() int
	GroupAndPackage() []string
	ConfigMap() util.ValueMap
}
