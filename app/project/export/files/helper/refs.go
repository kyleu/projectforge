package helper

import (
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/util"
)

func GoArgsWithRef(cols model.Columns, pkg string, args *metamodel.Args) (string, error) {
	return strings.Join(lo.Map(cols, func(c *model.Column, _ int) string {
		t := GoTypeWithRef(c, pkg, args)
		return c.Camel() + " " + t
	}), ", "), nil
}

func GoTypeWithRef(c *model.Column, pkg string, args *metamodel.Args) string {
	gt, err := c.ToGoType(pkg, args.Enums)
	if err != nil {
		return err.Error()
	}
	if ref, mdl, _ := LoadRef(c, args.Models, args.Events, args.ExtraTypes); ref != nil && !strings.Contains(gt, ".") {
		if mdl != nil && mdl.PackageName() != pkg {
			if strings.HasPrefix(gt, "*") {
				gt = "*" + mdl.PackageName() + "." + gt[1:]
			} else {
				gt = mdl.PackageName() + "." + gt
			}
		}
	}
	if gt == "*"+types.KeyAny {
		gt = types.KeyAny
	}
	return gt
}

func LoadRef(col *model.Column, models model.Models, events model.Events, extraTypes model.Models) (*types.Reference, model.StringProvider, error) {
	ret, err := model.AsRef(col.Type)
	if err != nil {
		return nil, nil, err
	}
	k := strings.TrimPrefix(ret.K, "*")
	get := func(key string) model.StringProvider {
		if ret := models.Get(key); ret != nil {
			return ret
		}
		if ret := events.Get(key); ret != nil {
			return ret
		}
		if ret := extraTypes.Get(key); ret != nil {
			return ret
		}
		return nil
	}
	sp := get(k)
	if sp == nil {
		deref := util.StringToCamel(k)
		sp = get(deref)
		ss := util.StringToSingular(deref)
		if sp == nil && ss != deref {
			sp = get(ss)
		}
	}
	if sp != nil {
		return &types.Reference{Pkg: sp.GroupAndPackage(), K: ret.K}, sp, nil
	}
	return ret, nil, nil
}
