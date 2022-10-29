package controller

import (
	"strings"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

const defaultPrefix = "controller."

func Controller(m *model.Model, args *model.Args, addHeader bool) (*file.File, error) {
	fn := m.Package
	if len(m.Group) > 0 {
		fn = m.GroupString("c", "") + "/" + fn
	}
	g := golang.NewFile(m.LastGroup("c", "controller"), []string{"app", "controller"}, fn)
	for _, imp := range helper.ImportsForTypes("parse", m.PKs().Types()...) {
		g.AddImport(imp)
	}
	if len(m.Group) > 0 {
		g.AddImport(helper.ImpAppController)
	}
	g.AddImport(helper.ImpFmt, helper.ImpErrors, helper.ImpFastHTTP, helper.ImpApp, helper.ImpCutil)
	g.AddImport(helper.AppImport("app/" + m.PackageWithGroup("")))
	g.AddImport(helper.AppImport("views/" + m.PackageWithGroup("v")))

	var prefix string
	if len(m.Group) > 0 {
		prefix = defaultPrefix
	}
	cl, err := controllerList(m, nil, args.Models, args.Enums, g, prefix)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(cl, controllerDetail(args.Models, m, nil, prefix))
	if m.IsRevision() {
		g.AddBlocks(controllerRevision(m, prefix))
	}
	g.AddBlocks(
		controllerCreateForm(m, nil, prefix), controllerCreateFormRandom(m, prefix), controllerCreate(m, g, nil, prefix),
		controllerEditForm(m, nil, prefix), controllerEdit(m, g, nil, prefix), controllerDelete(m, g, nil, prefix),
	)
	if m.IsHistory() {
		g.AddBlocks(controllerHistory(m, prefix))
	}
	g.AddBlocks(controllerModelFromPath(m), controllerModelFromForm(m))
	return g.Render(addHeader)
}

func controllerArgFor(col *model.Column, b *golang.Block, retVal string, indent int) {
	ind := ""
	for i := 0; i < indent; i++ {
		ind += "\t"
	}
	add := func(s string, args ...any) {
		b.W(ind+s, args...)
	}
	switch col.Type.Key() {
	case types.KeyBool:
		add("%sArg, err := cutil.RCRequiredBool(rc, %q)", col.Camel(), col.Camel())
		add("if err != nil {")
		add("\treturn %s, errors.Wrap(err, \"must provide [%s] as an argument\")", retVal, col.Camel())
		add("}")
		add("if err != nil {")
		add("\treturn %s, errors.Wrap(err, \"field [%s] must be a valid a valid boolean\")", retVal, col.Camel())
		add("}")
	case types.KeyInt:
		add("%sArgStr, err := cutil.RCRequiredString(rc, %q, false)", col.Camel(), col.Camel())
		add("if err != nil {")
		add("\treturn %s, errors.Wrap(err, \"must provide [%s] as an argument\")", retVal, col.Camel())
		add("}")
		add("%sArg, err := strconv.Atoi(%sArgStr)", col.Camel(), col.Camel())
		add("if err != nil {")
		add("\treturn %s, errors.Wrap(err, \"field [%s] must be a valid a valid integer\")", retVal, col.Camel())
		add("}")
	case types.KeyFloat:
		add("%sArgStr, err := cutil.RCRequiredString(rc, %q, false)", col.Camel(), col.Camel())
		add("if err != nil {")
		add("\treturn %s, errors.Wrap(err, \"must provide [%s] as an argument\")", retVal, col.Camel())
		add("}")
		add("%sArg, err := strconv.ParseFloat(%sArgStr, 64)", col.Camel(), col.Camel())
		add("if err != nil {")
		add("\treturn %s, errors.Wrap(err, \"field [%s] must be a valid a valid floating-point number\")", retVal, col.Camel())
		add("}")
	case types.KeyString:
		add("%sArg, err := cutil.RCRequiredString(rc, %q, false)", col.Camel(), col.Camel())
		add("if err != nil {")
		add("\treturn %s, errors.Wrap(err, \"must provide [%s] as an argument\")", retVal, col.Camel())
		add("}")
	case types.KeyList:
		if !types.IsStringList(col.Type) {
			add("// ERROR: invalid list argument [%s]", col.Type.String())
			break
		}
		add("%sArg, err := cutil.RCRequiredArray(rc, %q)", col.Camel(), col.Camel())
		add("if err != nil {")
		add("\treturn %s, errors.Wrap(err, \"must provide [%s] as an comma-separated argument\")", retVal, col.Camel())
		add("}")
	case types.KeyUUID:
		add("%sArgStr, err := cutil.RCRequiredString(rc, %q, false)", col.Camel(), col.Camel())
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
	properKeys := make([]string, 0, len(keys))
	for _, k := range keys {
		properKeys = append(properKeys, util.StringToTitle(k))
	}
	name := m.Proper() + withGroupName(strings.Join(properKeys, ""), grp)
	ret := golang.NewBlock(name, "func")
	ret.W("func %s(rc *fasthttp.RequestCtx) {", name)
	grpStr := ""
	if grp != nil {
		grpStr = grp.Name + "."
	}
	ret.W("\t%sAct(\"%s.%s%s\", rc, func(as *app.State, ps *cutil.PageState) (string, error) {", prefix, m.Package, grpStr, strings.Join(keys, "."))
	return ret
}

func withGroupName(s string, grp *model.Column) string {
	if grp == nil {
		return s
	}
	return s + "By" + grp.Proper()
}
