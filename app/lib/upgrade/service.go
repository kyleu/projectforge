// Content managed by Project Forge, see [projectforge.md] for details.
package upgrade

import (
	"context"
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/v39/github"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/kyleu/projectforge/app/util"
)

type Service struct {
	logger *zap.SugaredLogger
	client *github.Client
}

func NewService(logger *zap.SugaredLogger) *Service {
	return &Service{logger: logger, client: createGithubClient()}
}

func (s *Service) UpgradeIfNeeded(ctx context.Context, o string, n string, force bool) error {
	currVersion, err := semver.NewVersion(o)
	if err != nil {
		return errors.Wrapf(err, "unable to parse current version from [%s]", o)
	}
	tgtRelease, err := s.getRelease(ctx, n)
	if err != nil {
		return err
	}
	tgtVersion, err := semver.NewVersion(strings.TrimPrefix(*tgtRelease.TagName, "v"))
	if err != nil {
		return errors.Wrapf(err, "unable to parse version of latest release from [%s]", *tgtRelease.TagName)
	}
	if !force {
		if tgtVersion.Equal(*currVersion) {
			const msg = "no action needed, already using [%s], the latest version"
			s.logger.Infof(msg, currVersion.String())
			return nil
		}
		if tgtVersion.LessThan(*currVersion) {
			const msg = "no action needed, you're using [%s], which is somehow higher than [%s], the latest version"
			s.logger.Infof(msg, currVersion.String(), tgtVersion.String())
			return nil
		}
	}

	s.logger.Infof("upgrading from [%s] to [%s]...", currVersion.String(), tgtVersion.String())
	zipped, err := s.downloadAsset(*tgtVersion, tgtRelease)
	if err != nil {
		return err
	}

	fc, err := unzip(zipped)
	if err != nil {
		return err
	}

	err = overwrite(fc)
	if err != nil {
		return errors.Wrap(err, "unable to overwrite binary installation")
	}
	s.logger.Infof("successfully upgraded [%s] to version [%s]", util.AppKey, tgtVersion.String())
	return nil
}
