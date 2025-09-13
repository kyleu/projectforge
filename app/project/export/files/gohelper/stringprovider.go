package gohelper

type StringProvider interface {
	PackageName() string
	Camel() string
	CamelPlural() string
	Proper() string
	ProperPlural() string
	FirstLetter() string
	PackageWithGroup(prefix string) string
	RelativePath(rGroup []string, extra ...string) string
	GroupLen() int
}
