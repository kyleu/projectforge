package controller

import (
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

const defaultPrefix = "controller."

func Controller(m *model.Model, args *metamodel.Args, linebreak string) (*file.File, error) {
	fn := m.Package
	if len(m.Group) > 0 {
		fn = m.GroupString("c", "") + "/" + fn
	}
	g := golang.NewFile(m.LastGroup("c", "controller"), []string{"app", "controller"}, fn)
	lo.ForEach(helper.ImportsForTypes("parse", "", m.PKs().Types()...), func(imp *model.Import, _ int) {
		g.AddImport(imp)
	})
	g.AddImport(m.Imports.Supporting("controller")...)
	if len(m.Group) > 0 {
		g.AddImport(helper.ImpAppController)
	}
	g.AddImport(helper.ImpFmt, helper.ImpErrors, helper.ImpHTTP, helper.ImpApp, helper.ImpCutil)
	g.AddImport(helper.AppImport(m.PackageWithGroup("")))
	g.AddImport(helper.ViewImport(m.PackageWithGroup("v")))

	var prefix string
	if len(m.Group) > 0 {
		prefix = defaultPrefix
	}
	cl, err := controllerList(g, m, nil, args.Models.WithController(), args.Enums, prefix)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(cl, controllerDetail(g, args.Models, m, nil, args.Audit(m), prefix))
	g.AddBlocks(
		controllerCreateForm(g, m, nil, args.Models, prefix), controllerRandom(m, prefix), controllerCreate(m, nil, prefix),
		controllerEditForm(m, nil, prefix), controllerEdit(m, nil, prefix), controllerDelete(m, nil, prefix),
	)
	g.AddBlocks(controllerModelFromPath(m), controllerModelFromForm(m))
	return g.Render(linebreak)
}

func controllerArgFor(col *model.Column, b *golang.Block, retVal string, indent int) {
	ind := util.StringJoin(lo.Times(indent, func(_ int) string {
		return "\t"
	}), "")
	add := func(s string, args ...any) {
		b.WF(ind+s, args...)
	}
	switch col.Type.Key() {
	case types.KeyBool:
		add("%sArg, err := cutil.PathBool(r, %q)", col.Camel(), col.Camel())
		add("if err != nil {")
		add("\treturn %s, errors.Wrap(err, \"must provide [%s] as a boolean argument\")", retVal, col.Camel())
		add("}")
	case types.KeyInt:
		add("%sArgStr, err := cutil.PathString(r, %q, false)", col.Camel(), col.Camel())
		add("if err != nil {")
		add("\treturn %s, errors.Wrap(err, \"must provide [%s] as an argument\")", retVal, col.Camel())
		add("}")
		bits := types.Bits(col.Type)
		if bits == 0 {
			add("%sArgX, err := strconv.ParseInt(%sArgStr, 10, 64)", col.Camel(), col.Camel())
		} else {
			add("%sArgX, err := strconv.ParseInt(%sArgStr, 10, %d)", col.Camel(), col.Camel(), bits)
		}
		add("if err != nil {")
		add("\treturn %s, errors.Wrap(err, \"field [%s] must be a valid a valid integer\")", retVal, col.Camel())
		add("}")
		if bits == 0 {
			add("%sArg := int(%sArgX)", col.Camel(), col.Camel())
		} else {
			add("%sArg := int%d(%sArgX)", col.Camel(), bits, col.Camel())
		}
	case types.KeyFloat:
		add("%sArgStr, err := cutil.PathString(r, %q, false)", col.Camel(), col.Camel())
		add("if err != nil {")
		add("\treturn %s, errors.Wrap(err, \"must provide [%s] as an argument\")", retVal, col.Camel())
		add("}")
		add("%sArg, err := strconv.ParseFloat(%sArgStr, 64)", col.Camel(), col.Camel())
		add("if err != nil {")
		add("\treturn %s, errors.Wrap(err, \"field [%s] must be a valid a valid floating-point number\")", retVal, col.Camel())
		add("}")
	case types.KeyString:
		add("%sArg, err := cutil.PathString(r, %q, false)", col.Camel(), col.Camel())
		add("if err != nil {")
		add("\treturn %s, errors.Wrap(err, \"must provide [%s] as a string argument\")", retVal, col.Camel())
		add("}")
	case types.KeyList:
		if !types.IsStringList(col.Type) {
			add("// ERROR: invalid list argument [%s]", col.Type.String())
			break
		}
		add("%sArg, err := cutil.PathArray(r, %q)", col.Camel(), col.Camel())
		add("if err != nil {")
		add("\treturn %s, errors.Wrap(err, \"must provide [%s] as an comma-separated argument\")", retVal, col.Camel())
		add("}")
	case types.KeyUUID:
		add("%sArgStr, err := cutil.PathString(r, %q, false)", col.Camel(), col.Camel())
		add("if err != nil {")
		add("\treturn %s, errors.Wrap(err, \"must provide [%s] as an argument\")", retVal, col.Camel())
		add("}")
		add("%sArgP := util.UUIDFromString(%sArgStr)", col.Camel(), col.Camel())
		add("if %sArgP == nil {", col.Camel())
		add("\treturn %s, errors.Errorf(\"argument [%s] (%%%%s) is not a valid UUID\", %sArgStr)", retVal, col.Camel(), col.Camel())
		add("}")
		add("%sArg := *%sArgP", col.Camel(), col.Camel())
	default:
		add("// ERROR: unhandled controller arg type [%s]", col.Type.String())
	}
}

func blockFor(m *model.Model, prefix string, grp *model.Column, keys ...string) *golang.Block {
	properKeys := lo.Map(keys, func(k string, _ int) string {
		return util.StringToTitle(k)
	})
	name := m.Proper() + withGroupName(util.StringJoin(properKeys, ""), grp)
	ret := golang.NewBlock(name, "func")
	ret.WF("func %s(w http.ResponseWriter, r *http.Request) {", name)
	var grpStr string
	if grp != nil {
		grpStr = grp.Name + "."
	}
	ret.WF("\t%sAct(\"%s.%s%s\", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {", prefix, m.Package, grpStr, util.StringJoin(keys, "."))
	return ret
}

func withGroupName(s string, grp *model.Column) string {
	if grp == nil {
		return s
	}
	return s + "By" + grp.Proper()
}
