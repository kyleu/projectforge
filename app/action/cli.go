package action

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/app/util"
)

func cliProject(p *project.Project, modKeys []string) error {
	reader := bufio.NewReader(os.Stdin)

	promptString := func(query string, curr string) string {
		print(query)
		if curr != "" {
			print(" (default: " + curr + ")")
		}
		println()
		print(" > ")
		text, err := reader.ReadString('\n')
		if err != nil {
			println("error: " + err.Error())
		}
		text = strings.TrimSuffix(text, "\n")
		if text == "" {
			text = curr
		}
		return text
	}

	if p.Key == "TODO" {
		p.Key = ""
	}
	p.Key = promptString("Enter a project key; must only contain alphanumerics", p.Key)

	if p.Name == "" {
		p.Name = p.Key
	}
	p.Name = promptString("Enter a project name; use title case and spaces if needed", p.Name)

	if p.Icon == "" {
		p.Icon = "star"
	}

	if p.Exec == "" {
		p.Exec = p.Key
	}

	p.Version = promptString("Enter a version, such as 0.0.0", p.Version)

	p.Info.Org = promptString("Enter the github organization that owns this project", p.Info.Org)

	if p.Package == "" {
		p.Package = "github.com/" + p.Info.Org + "/" + p.Key
	}
	p.Package = promptString("Enter your project's package", p.Package)

	if p.Port == 0 {
		p.Port = 20000
	}
	p.Port, _ = strconv.Atoi(promptString("Enter the default port your http server will run on", fmt.Sprint(p.Port)))

	modPromptString := "Enter the modules your project will use as a comma-separated list from choices [" + strings.Join(modKeys, ", ") + "]"
	mods := promptString(modPromptString, strings.Join(p.Modules, ", "))
	p.Modules = util.SplitAndTrim(mods, ",")

	if p.Info.License == "" {
		p.Info.License = "Proprietary"
	}
	p.Info.License = promptString("Enter the license used by this project", p.Info.License)

	return nil
}
