package helper

import (
	"strings"

	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/util"
)

func LoadRef(col *model.Column, models model.Models) (*types.Reference, error) {
	ret, err := model.AsRef(col.Type)
	if err != nil {
		return nil, err
	}
	if len(ret.Pkg) == 0 {
		deref := util.StringToLowerCamel(strings.TrimPrefix(ret.K, "*"))
		mdl := models.Get(deref)
		if mdl != nil {
			return &types.Reference{Pkg: append(mdl.Group, mdl.Package), K: ret.K}, nil
		}
		if ss := util.StringToSingular(deref); ss != deref {
			mdl = models.Get(ss)
			if mdl != nil {
				return &types.Reference{Pkg: append(mdl.Group, mdl.Package), K: ret.K}, nil
			}
		}
	}
	return ret, nil
}
