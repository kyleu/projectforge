package build

import (
	"slices"
	"strings"

	"github.com/samber/lo"
)

const textArrow = " -> "

type Dependency struct {
	Key        string   `json:"key"`
	Version    string   `json:"version,omitzero"`
	Available  string   `json:"available,omitzero"`
	References []string `json:"references,omitempty"`
}

func ParseDependency(line string) *Dependency {
	if strings.HasPrefix(line, "\t") && strings.Contains(line, " ") {
		start := strings.Index(line, " v")
		if start == -1 {
			return nil
		}

		return &Dependency{Key: strings.TrimSpace(line[:start]), Version: strings.TrimSpace(line[start:])}
	}
	return nil
}

func (d *Dependency) AddRef(r string) {
	if lo.Contains(d.References, r) {
		return
	}
	d.References = append(d.References, r)
	slices.Sort(d.References)
}

func (d *Dependency) String() string {
	ret := d.Key
	if d.Version != "" {
		ret += ": " + d.Version
	}
	if d.Available != "" && d.Available != d.Version {
		ret += textArrow + d.Available
	}
	return ret
}

type Dependencies []*Dependency

func (d Dependencies) Get(k string) *Dependency {
	return lo.FindOrElse(d, nil, func(x *Dependency) bool {
		return x.Key == k
	})
}
