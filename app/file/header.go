package file

import (
	"strings"

	"projectforge.dev/projectforge/app/util"
)

const HeaderContent = "Content managed by Project Forge, see [projectforge.md] for details."

func ContainsHeader(s string) bool {
	return strings.Contains(s, HeaderContent) || strings.Contains(s, "$PF_IGNORE$")
}

func contentWithHeader(filename string, t Type, c string, linebreak string, logger util.Logger) string {
	if strings.Contains(c, IgnorePattern) {
		return c
	}
	switch t.Key {
	case TypeBatch.Key:
		return secondLine(c, "rem "+HeaderContent, linebreak)
	case TypeCodeowners.Key, TypeDocker.Key, TypeDockerIgnore.Key, TypeEnv.Key, TypeGraphQL.Key, TypeHCL.Key, TypeMakefile.Key, TypeProperties.Key, TypeYAML.Key:
		return "# " + HeaderContent + linebreak + c
	case TypeConf.Key, TypeEditorConfig.Key, TypeESLint.Key, TypeGitIgnore.Key:
		return c
	case TypeIcons.Key, TypeIgnore.Key, TypeJSON.Key, TypePList.Key, TypeProtobuf.Key, TypeSVG.Key:
		return c
	case TypeCSS.Key:
		return "/* " + HeaderContent + " */" + linebreak + c
	case TypeGo.Key, TypeGoMod.Key, TypeGradle.Key, TypeJavaScript.Key, TypeKotlin.Key, TypeSwift.Key, TypeTypeScript.Key:
		return "// " + HeaderContent + linebreak + c
	case TypeHTML.Key:
		return "<!-- " + HeaderContent + " -->" + linebreak + c
	case TypeMarkdown.Key:
		return "<!--- " + HeaderContent + " -->" + linebreak + c
	case TypeSQL.Key:
		return "-- " + HeaderContent + linebreak + c
	case TypeShell.Key:
		return secondLine(c, "# "+HeaderContent, linebreak)
	case TypeEntitlements.Key, TypeXML.Key:
		return secondLine(c, "<!-- "+HeaderContent+" -->", linebreak)
	default:
		logger.Warnf("unhandled header for file [%s], of type [%s]", filename, t.Key)
		return c
	}
}

func secondLine(content string, rplc string, linebreak string) string {
	idx := strings.Index(content, linebreak)
	if idx == -1 {
		return content
	}
	return content[0:idx] + linebreak + rplc + content[idx:]
}
