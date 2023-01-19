package action

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func ServiceDefinition(p *project.Project) util.ValueMap {
	team := p.Info.Team
	if team == "" {
		team = p.Info.AuthorID
	}

	links := []util.ValueMap{
		{"name": "repo", "type": "repo", "url": p.Info.Sourcecode},
		{"name": "url", "type": "url", "url": p.Info.Homepage},
	}
	for _, x := range p.Info.Deployments {
		links = append(links, util.ValueMap{"name": "deployment", "type": "link", "url": x})
	}

	tags := make([]string, 0, len(p.Tags))
	tags = append(tags, fmt.Sprintf("service:%s", p.CleanKey()))
	//for _, x := range p.Tags {
	//	tags = append(tags, fmt.Sprintf("%s:%s", x, x))
	//}
	for _, x := range p.Modules {
		switch x {
		case "export":
			tags = append(tags, "codegen:true")
		case "expression":
			tags = append(tags, x+":cel")
		case "graphql", "gqlgen":
			tags = append(tags, "graphql:"+x)
		case "grpc":
			tags = append(tags, "transport:"+x)
		case "mysql":
			tags = append(tags, "database:"+x)
		case "oauth":
			tags = append(tags, "auth:"+x)
		case "postgres":
			tags = append(tags, "database:"+x)
		case "queue":
			tags = append(tags, x+":kafka")
		case "sqlite":
			tags = append(tags, "database:"+x)
		case "temporal":
			tags = append(tags, "workflow:"+x)
		case "wasm":
			tags = append(tags, "build:"+x)
		}
	}

	gh := p.Info.AuthorID
	if !strings.Contains(gh, "/") {
		gh = "https://github.com/" + gh
	}

	contacts := []util.ValueMap{
		{"name": p.Info.AuthorName, "type": "email", "contact": p.Info.AuthorEmail},
		{"name": p.Info.AuthorName, "type": "github", "contact": gh},
	}
	for _, x := range p.Info.Channels {
		l, r := util.StringSplit(x, ':', true)
		if r == "" {
			r = l
			l = "unknown"
		}
		ch, u := util.StringSplit(r, ',', true)
		if u == "" {
			u = ch
		}
		contacts = append(contacts, util.ValueMap{"name": ch, "type": l, "contact": u})
	}

	docs := []util.ValueMap{{"name": "Source Code", "provider": "github", "url": p.Info.Sourcecode}}
	for _, x := range p.Info.Docs {
		docs = append(docs, util.ValueMap{"name": x.Name, "provider": x.Provider, "url": x.URL})
	}
	ret := util.ValueMap{
		"schema-version": "v2",
		"dd-service":     p.Key,
		"team":           team,
		"dd-team":        team,
		"contacts":       contacts,
		"links":          links,
		"repos": []util.ValueMap{
			{"name": "sourcecode", "provider": "github", "url": p.Info.Sourcecode},
		},
		"docs":         docs,
		"integrations": util.ValueMap{},
		"extensions":   util.ValueMap{},
		"tags":         tags,
	}
	return ret
}
