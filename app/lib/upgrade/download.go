// Content managed by Project Forge, see [projectforge.md] for details.
package upgrade

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"runtime"

	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/v39/github"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

func assetFor(version semver.Version) string {
	o := runtime.GOOS
	if o == "darwin" {
		o = "macos"
	}
	arch := runtime.GOARCH
	if arch == "amd64" {
		arch = "x86_64"
	}
	return fmt.Sprintf("%s_%s_%s_%s.zip", util.AppKey, version.String(), o, arch)
}

func (s *Service) downloadAsset(ctx context.Context, version semver.Version, release *github.RepositoryRelease) ([]byte, error) {
	candidate := assetFor(version)
	match := lo.FindOrElse(release.Assets, nil, func(a *github.ReleaseAsset) bool {
		return a.Name != nil && (*a.Name) == candidate
	})
	if match == nil {
		return nil, errors.Errorf("no asset available for version [%s] with name [%s]", version.String(), candidate)
	}
	if match.BrowserDownloadURL == nil || *match.BrowserDownloadURL == "" {
		return nil, errors.Errorf("no asset url available in asset [%s]", candidate)
	}

	org, repo, err := parseSource()
	if err != nil {
		return nil, err
	}

	rsp, _, err := s.client.Repositories.DownloadReleaseAsset(ctx, org, repo, *match.ID, s.client.Client())
	if err != nil {
		return nil, errors.Wrapf(err, "unable to download asset from [%s]", *match.BrowserDownloadURL)
	}
	return io.ReadAll(rsp)
}

func unzip(zipped []byte) ([]byte, error) {
	r, err := zip.NewReader(bytes.NewReader(zipped), int64(len(zipped)))
	if err != nil {
		return nil, errors.Wrapf(err, "unable to unzip response of size [%d]", len(zipped))
	}
	var ret []byte
	for _, f := range r.File {
		if ret != nil {
			return nil, errors.New("multiple files found in zip")
		}
		reader, err := f.Open()
		if err != nil {
			return nil, errors.Wrapf(err, "unable to open file [%s] from zip", f.Name)
		}
		ret, err = io.ReadAll(reader)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to read file [%s] from zip", f.Name)
		}
	}
	return ret, nil
}
