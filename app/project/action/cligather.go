package action

import (
	"context"

	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func cliGather(_ context.Context, p *project.Project, mods module.Modules, _ util.Logger) error {
	for _, q := range CreatePrompts() {
		switch q.Kind {
		case CreatePromptString:
			curr := CreatePromptDefaultString(p, q.Key)
			val, err := promptString(q.Label, q.Query, curr)
			if err != nil {
				return err
			}
			if err := ApplyCreatePromptString(p, q.Key, val); err != nil {
				return err
			}
		case CreatePromptModules:
			val, err := promptModules(CreatePromptDefaultModules(p), mods)
			if err != nil {
				return err
			}
			ApplyCreatePromptModules(p, val)
		}
	}
	return nil
}
