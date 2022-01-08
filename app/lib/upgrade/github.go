package upgrade

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/google/go-github/v39/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	"github.com/kyleu/projectforge/app/util"
)

func createGithubClient() *github.Client {
	client := http.DefaultClient
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		client = oauth2.NewClient(context.Background(), src)
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
		if res != nil && res.StatusCode == 404 {
			return nil, errors.Errorf("can't access repository at [%s]", util.AppSource)
		}
	}
	return rel, err
}
