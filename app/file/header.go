package file

import (
	"strings"

	"go.uber.org/zap"
)

const HeaderContent = "Content managed by Project Forge, see [projectforge.md] for details."

func ContainsHeader(s string) bool {
	return strings.Contains(s, HeaderContent)
}

func contentWithHeader(t Type, c string, logger *zap.SugaredLogger) string {
	if strings.Contains(c, "$PF_IGNORE$") {
		return c
	}
	switch t.Key {
	case TypeBatch.Key:
		return "rem " + HeaderContent + "\n" + c
	case TypeCodeowners.Key:
		return "# " + HeaderContent + "\n" + c
	case TypeConf.Key:
		return c
	case TypeCSS.Key:
		return "/* " + HeaderContent + " */\n" + c
	case TypeDocker.Key:
		return "# " + HeaderContent + "\n" + c
	case TypeEntitlements.Key:
		return secondLine(c, "<!-- "+HeaderContent+" -->")
	case TypeGitIgnore.Key:
		return c
	case TypeGo.Key:
		return "// " + HeaderContent + "\n" + c
	case TypeGoMod.Key:
		return "// " + HeaderContent + "\n" + c
	case TypeGradle.Key:
		return "// " + HeaderContent + "\n" + c
	case TypeHCL.Key:
		return "# " + HeaderContent + "\n" + c
	case TypeHTML.Key:
		return "<!-- " + HeaderContent + " -->\n" + c
	case TypeIcons.Key:
		return c
	case TypeIgnore.Key:
		return c
	case TypeJavaScript.Key:
		return "// " + HeaderContent + "\n" + c
	case TypeJSON.Key:
		return c
	case TypeKotlin.Key:
		return "// " + HeaderContent + "\n" + c
	case TypeMakefile.Key:
		return "# " + HeaderContent + "\n" + c
	case TypeMarkdown.Key:
		return "<!--- " + HeaderContent + " -->\n" + c
	case TypePList.Key:
		return c
	case TypeProperties.Key:
		return "# " + HeaderContent + "\n" + c
	case TypeShell.Key:
		return secondLine(c, "# "+HeaderContent)
	case TypeSQL.Key:
		return "-- " + HeaderContent + "\n" + c
	case TypeSVG.Key:
		return c
	case TypeSwift.Key:
		return "// " + HeaderContent + "\n" + c
	case TypeTypeScript.Key:
		return "// " + HeaderContent + "\n" + c
	case TypeXML.Key:
		return secondLine(c, "<!-- "+HeaderContent+" -->")
	case TypeYAML.Key:
		return "# " + HeaderContent + "\n" + c
	default:
		logger.Warnf("unhandled header for file type [%s]", t.Title)
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
