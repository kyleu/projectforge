package upgrade

import (
	"context"
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/v39/github"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"{{{ .Package }}}/app/util"
)

type Service struct {
	logger *zap.SugaredLogger
	client *github.Client
}

func NewService(logger *zap.SugaredLogger) *Service {
	return &Service{logger: logger, client: createGithubClient()}
}

func (s *Service) UpgradeIfNeeded(ctx context.Context, version string, force bool) error {
	currVersion, err := semver.NewVersion(version)
	if err != nil {
		return errors.Wrapf(err, "unable to parse current version from [%s]", version)
	}
	latestRelease, err := s.latestRelease(ctx)
	if err != nil {
		return err
	}
	latestVersion, err := semver.NewVersion(strings.TrimPrefix(*latestRelease.TagName, "v"))
	if err != nil {
		return errors.Wrapf(err, "unable to parse version of latest release from [%s]", *latestRelease.TagName)
	}
	if !force {
		if latestVersion.Equal(*currVersion) {
			msg := "no action needed, already using [%s], the latest version"
			s.logger.Infof(msg, currVersion.String())
			return nil
		}
		if latestVersion.LessThan(*currVersion) {
			msg := "no action needed, you're using [%s], which is somehow higher than [%s], the latest version"
			s.logger.Infof(msg, currVersion.String(), latestVersion.String())
			return nil
		}
	}

	s.logger.Infof("upgrading from [%s] to [%s]...", currVersion.String(), latestVersion.String())
	zipped, err := downloadAsset(*latestVersion, latestRelease)
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
	s.logger.Infof("successfully upgraded [%s] to version [%s]", util.AppKey, latestVersion.String())
	return nil
}
