package upgrade

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/go-github/v39/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	"{{{ .Package }}}/app/util"
)

func createGitHubClient(ctx context.Context) *github.Client {
	client := http.DefaultClient
	if token := util.GetEnv("github_token"); token != "" {
		src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		client = oauth2.NewClient(ctx, src)
	}
	return github.NewClient(client)
}

func (s *Service) getRelease(ctx context.Context, n string) (*github.RepositoryRelease, error) {
	org, repo, err := parseSource()
	if err != nil {
		return nil, err
	}

	var rel *github.RepositoryRelease
	var res *github.Response
	if n == "" {
		rel, res, err = s.client.Repositories.GetLatestRelease(ctx, org, repo)
	} else {
		if !strings.HasPrefix(n, "v") {
			n = "v" + n
		}
		rel, res, err = s.client.Repositories.GetReleaseByTag(ctx, org, repo, n)
	}
	if err != nil {
		if res != nil && res.StatusCode == http.StatusNotFound {
			return nil, errors.Errorf("can't access repository at [%s]", util.AppSource)
		}
	}
	return rel, err
}
