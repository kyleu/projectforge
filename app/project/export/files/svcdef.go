package files

import (
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func ServiceDefinition(p *project.Project) util.ValueMap {
	links := []util.ValueMap{
		{"name": "repo", "type": "repo", "url": p.Info.Sourcecode},
		{"name": "url", "type": "url", "url": p.Info.Homepage},
	}
	for _, x := range p.Info.Deployments {
		links = append(links, util.ValueMap{"name": "deployment", "type": "link", "url": x})
	}
	ret := util.ValueMap{
		"schema-version": "v2",
		"dd-service":     p.Key,
		"team":           p.Info.AuthorID,
		"dd-team":        p.Info.AuthorID,
		"contacts": []util.ValueMap{
			{
				"name":    p.Info.AuthorName,
				"type":    "email",
				"contact": p.Info.AuthorEmail,
			},
		},
		"links": links,
		"repos": []util.ValueMap{
			{"name": "sourcecode", "provider": "github", "url": p.Info.Sourcecode},
		},
		"docs": []util.ValueMap{
			{"name": "sourcecode", "provider": "github", "url": p.Info.Sourcecode},
		},
		"integrations": []util.ValueMap{},
		"extensions":   []util.ValueMap{},
		"tags":         p.Tags,
	}
	return ret
}
