package helper

import (
	"strings"

	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/util"
)

func LoadRef(col *model.Column, models model.Models, extraTypes model.Models) (*types.Reference, *model.Model, error) {
	ret, err := model.AsRef(col.Type)
	if err != nil {
		return nil, nil, err
	}
	mdl := models.Get(ret.K)
	if mdl == nil {
		mdl = extraTypes.Get(ret.K)
	}
	if mdl == nil {
		deref := util.StringToCamel(strings.TrimPrefix(ret.K, "*"))
		mdl = models.Get(deref)
		if mdl == nil {
			mdl = extraTypes.Get(deref)
		}
		ss := util.StringToSingular(deref)
		if mdl == nil && ss != deref {
			mdl = models.Get(ss)
			if mdl == nil {
				mdl = extraTypes.Get(ss)
			}
		}
	}
	if mdl != nil {
		pkg := mdl.Group
		pkg = append(pkg, mdl.Package)
		return &types.Reference{Pkg: pkg, K: ret.K}, mdl, nil
	}
	return ret, nil, nil
}
