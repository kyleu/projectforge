package script

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/export/model"
)

func NotebookScript(p *project.Project, args *model.Args, addHeader bool, linebreak string) (*file.File, error) {
	var content []string
	w := func(msg string, args ...any) {
		content = append(content, fmt.Sprintf(msg, args...))
	}
	if addHeader {
		w("// " + file.HeaderContent)
	}
	w(`import * as d3 from "npm:d3";`)
	w("")
	w("export async function load(u, t) {")
	w("  const response = await fetch(u);")
	w("  if (!response.ok) throw new Error(`fetch failed: ${response.status}`);")
	w("  if (t === \"csv\") {")
	w("    return d3.csvParse(await response.text());")
	w("  }")
	w("  if (t === \"json\" || t === \"\") {")
	w("    return await response.json();")
	w("  }")
	w("  return await response.text();")
	w("}")
	w("")
	w(`const defaultOpts = {limit: 0, offset: 0, order: "", t: "json", q: "", extra: {}};`)
	w("")
	w("export function urlFor(key, path, {limit, offset, order, t, q, extra} = defaultOpts) {")
	w("  let ret = `http://localhost:" + fmt.Sprint(p.Port) + "/${path}?t=${t}`;")
	w(`  let prefix = key !== "" ? key + "." : "";`)
	w(`  if (order && order !== "") {`)
	w("    ret += `&${prefix}o=${order}`;")
	w("  }")
	w("  if (limit > 0) {")
	w("    ret += `&${prefix}l=${limit}`;")
	w("  }")
	w("  if (offset > 0) {")
	w("    ret += `&${prefix}x=${offset}`;")
	w("  }")
	w(`  if (q && q !== "") {`)
	w("    ret += `&q=${encodeURIComponent(q)}`;")
	w("  }")
	w("  if (extra) {")
	w("    for (const [key, value] of Object.entries(extra)) {")
	w("      ret += `&${encodeURIComponent(key)}=${encodeURIComponent(value)}`;\n")
	w("    }")
	w("  }")
	w("  return ret;")
	w("}")

	for _, m := range args.Models {
		w("")
		w("export async function %s(opts) {", m.CamelPlural())
		w("  return await load(urlFor(%q, %q, opts), opts?.t);", m.Package, m.Route())
		w("}")
		_ = m
	}
	w("")

	return &file.File{
		Type:    file.TypeJavaScript,
		Path:    []string{"notebook", "docs", "components"},
		Name:    p.Key + ".js",
		Mode:    filesystem.DefaultMode,
		Content: strings.Join(content, linebreak),
	}, nil
}
