package build

import (
	"sort"

	"golang.org/x/exp/slices"
)

type Dependency struct {
	Key        string   `json:"key"`
	Version    string   `json:"version"`
	Available  string   `json:"available"`
	References []string `json:"references"`
}

func (d *Dependency) AddRef(r string) {
	if slices.Contains(d.References, r) {
		return
	}
	d.References = append(d.References, r)
	sort.Strings(d.References)
}

func (d *Dependency) String() string {
	ret := d.Key
	if d.Version != "" {
		ret += ": " + d.Version
	}
	if d.Available != "" && d.Available != d.Version {
		ret += " -> " + d.Available
	}
	return ret
}

type Dependencies []*Dependency

func (d Dependencies) Get(k string) *Dependency {
	for _, x := range d {
		if x.Key == k {
			return x
		}
	}
	return nil
}
