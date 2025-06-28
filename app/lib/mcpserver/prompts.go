package mcpserver

import (
	"projectforge.dev/projectforge/app/util"
)

var ProjectPrompt = &Prompt{
	Name:        "project_usage",
	Description: "A simple prompt that helps the system build, test, and work with an application managed by Project Forge" + util.AppName,
	Content:     `This application is written using Go, and is a web application managed by Project Forge.`,
}
