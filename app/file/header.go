package file

import (
	"strings"

	"projectforge.dev/projectforge/app/util"
)

const HeaderContent = "Content managed by Project Forge, see [projectforge.md] for details."

func ContainsHeader(s string) bool {
	return strings.Contains(s, HeaderContent) || strings.Contains(s, "$PF_IGNORE$")
}

func contentWithHeader(filename string, t Type, c string, logger util.Logger) string {
	if strings.Contains(c, IgnorePattern) {
		return c
	}
	switch t.Key {
	case TypeBatch.Key:
		return "rem " + HeaderContent + "\n" + c
	case TypeCodeowners.Key, TypeDocker.Key, TypeYAML.Key, TypeProperties.Key, TypeMakefile.Key, TypeHCL.Key:
		return "# " + HeaderContent + "\n" + c
	case TypeConf.Key, TypeEditorConfig.Key, TypeESLint.Key, TypeGitIgnore.Key:
		return c
	case TypeIcons.Key, TypeIgnore.Key, TypeJSON.Key, TypePList.Key, TypeProtobuf.Key, TypeSVG.Key:
		return c
	case TypeCSS.Key:
		return "/* " + HeaderContent + " */\n" + c
	case TypeGo.Key, TypeGoMod.Key, TypeGradle.Key, TypeJavaScript.Key, TypeKotlin.Key, TypeSwift.Key, TypeTypeScript.Key:
		return "// " + HeaderContent + "\n" + c
	case TypeHTML.Key:
		return "<!-- " + HeaderContent + " -->\n" + c
	case TypeMarkdown.Key:
		return "<!--- " + HeaderContent + " -->\n" + c
	case TypeSQL.Key:
		return "-- " + HeaderContent + "\n" + c
	case TypeShell.Key:
		return secondLine(c, "# "+HeaderContent)
	case TypeEntitlements.Key, TypeXML.Key:
		return secondLine(c, "<!-- "+HeaderContent+" -->")
	default:
		logger.Warnf("unhandled header for file [%s], of type [%s]", filename, t.Key)
		return c
	}
}

func secondLine(content string, rplc string) string {
	idx := strings.Index(content, "\n")
	if idx == -1 {
		return content
	}
	return content[0:idx] + "\n" + rplc + content[idx:]
}
