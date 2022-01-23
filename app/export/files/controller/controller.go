package controller

import (
	"strings"

	"github.com/kyleu/projectforge/app/export/files/helper"
	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/file"
	"github.com/kyleu/projectforge/app/util"
)

func Controller(m *model.Model, args *model.Args, addHeader bool) (*file.File, error) {
	g := golang.NewFile("controller", []string{"app", "controller"}, m.Package)
	for _, imp := range helper.ImportsForTypes("parse", m.PKs().Types()...) {
		g.AddImport(imp)
	}
	g.AddImport(helper.ImpFmt, helper.ImpErrors, helper.ImpFastHTTP, helper.ImpApp, helper.ImpCutil)
	g.AddImport(helper.AppImport("app/" + m.Package))
	g.AddImport(helper.AppImport("views/v" + m.Package))
	g.AddBlocks(controllerTitle(m), controllerList(m, nil), controllerDetail(args.Models, m, nil))
	if m.IsRevision() {
		g.AddBlocks(controllerRevision(m))
	}
	g.AddBlocks(
		controllerCreateForm(m, nil), controllerCreateFormRandom(m), controllerCreate(m, g, nil),
		controllerEditForm(m, nil), controllerEdit(m, g, nil), controllerDelete(m, g, nil),
	)
	if m.IsHistory() {
		g.AddBlocks(controllerHistory(m))
	}
	g.AddBlocks(controllerModelFromPath(m), controllerModelFromForm(m))
	return g.Render(addHeader)
}

func controllerArgFor(col *model.Column, b *golang.Block, retVal string, indent int) {
	ind := ""
	for i := 0; i < indent; i++ {
		ind += "\t"
	}
	add := func(s string, args ...interface{}) {
		b.W(ind+s, args...)
	}
	switch col.Type.Key {
	case model.TypeInt.Key:
		add("%sArgStr, err := RCRequiredString(rc, %q, false)", col.Camel(), col.Camel())
		add("if err != nil {")
		add("\treturn %s, errors.Wrap(err, \"must provide [%s] as an argument\")", retVal, col.Camel())
		add("}")
		add("%sArg, err := strconv.Atoi(%sArgStr)", col.Camel(), col.Camel())
		add("if err != nil {")
		add("\treturn %s, errors.Wrap(err, \"field [%s] must be a valid a valid integer\")", retVal, col.Camel())
		add("}")
	case model.TypeString.Key:
		add("%sArg, err := RCRequiredString(rc, %q, false)", col.Camel(), col.Camel())
		add("if err != nil {")
		add("\treturn %s, errors.Wrap(err, \"must provide [%s] as an argument\")", retVal, col.Camel())
		add("}")
	case model.TypeUUID.Key:
		add("%sArgStr, err := RCRequiredString(rc, %q, false)", col.Camel(), col.Camel())
		add("if err != nil {")
		add("\treturn %s, errors.Wrap(err, \"must provide [%s] as an argument\")", retVal, col.Camel())
		add("}")
		add("%sArgP := util.UUIDFromString(%sArgStr)", col.Camel(), col.Camel())
		add("if %sArgP == nil {", col.Camel())
		add("\treturn %s, errors.Errorf(\"argument [%s] (%%%%s) is not a valid UUID\", %sArgStr)", retVal, col.Camel(), col.Camel())
		add("}")
		add("%sArg := *%sArgP", col.Camel(), col.Camel())
	default:
		add("ERROR: unhandled controller arg type [%s]", col.Type.String())
	}
}

func blockFor(m *model.Model, grp *model.Column, keys ...string) *golang.Block {
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
	ret.W("\tact(\"%s.%s%s\", rc, func(as *app.State, ps *cutil.PageState) (string, error) {", m.Package, grpStr, strings.Join(keys, "."))
	return ret
}

func withGroupName(s string, grp *model.Column) string {
	if grp == nil {
		return s
	}
	return s + "By" + grp.Proper()
}

func controllerTitle(m *model.Model) *golang.Block {
	ret := golang.NewBlock("Title", "func")
	ret.W("const %sDefaultTitle = \"%s\"", m.Camel(), m.TitlePlural())
	return ret
}
