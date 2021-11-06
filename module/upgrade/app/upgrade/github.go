package upgrade

import (
	"context"
	"net/http"
	"os"

	"github.com/google/go-github/v39/github"
	"github.com/pkg/errors"
	"github.com/tcnksm/go-gitconfig"
	"golang.org/x/oauth2"

	"{{{ .Package }}}/app/util"
)

func createGithubClient() *github.Client {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		token, _ = gitconfig.GithubToken()
	}

	client := http.DefaultClient
	if token != "" {
		src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		client = oauth2.NewClient(context.Background(), src)
	}
	return github.NewClient(client)
}

func (s *Service) latestRelease(ctx context.Context) (*github.RepositoryRelease, error) {
	org, repo, err := parseSource()
	rel, res, err := s.client.Repositories.GetLatestRelease(ctx, org, repo)
	if err != nil {
		if res != nil && res.StatusCode == 404 {
			err = nil
			return nil, errors.Errorf("can't access repository at [%s]", util.AppSource)
		}
	}
	return rel, nil
}
