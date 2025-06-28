package mcpserver

import (
	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/util"
)

var (
	AllTools             = Tools{ListProjectsTool, GetProjectTool}
	AllResources         = Resources{}
	AllResourceTemplates = ResourceTemplates{ProjectContentResource}
	AllPrompts           = Prompts{ProjectPrompt}
)

func WireLibrary(s *Server, as *app.State, logger util.Logger) error {
	if err := s.AddTools(as, logger, AllTools...); err != nil {
		return err
	}
	if err := s.AddResources(as, logger, AllResources...); err != nil {
		return err
	}
	if err := s.AddResourceTemplates(as, logger, AllResourceTemplates...); err != nil {
		return err
	}
	if err := s.AddPrompts(as, logger, AllPrompts...); err != nil {
		return err
	}
	return nil
}
