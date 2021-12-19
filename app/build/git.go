package build

import (
	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
)

func gitStatus(path string) (interface{}, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		if err.Error() == "repository does not exist" {
			return "no-repo", nil
		}
		return nil, errors.Wrapf(err, "unable to open git repo from path [%s]", path)
	}

	_, _ = repo.Notes()

	ret := map[string]interface{}{"repo": repo}
	return ret, nil
}
