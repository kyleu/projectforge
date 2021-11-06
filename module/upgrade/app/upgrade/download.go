package upgrade

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"

	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/v39/github"
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/util"
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

func downloadAsset(version semver.Version, release *github.RepositoryRelease) ([]byte, error) {
	candidate := assetFor(version)
	var match *github.ReleaseAsset
	for _, a := range release.Assets {
		if a.Name != nil && (*a.Name) == candidate {
			match = a
			break
		}
	}
	if match == nil {
		return nil, errors.Errorf("no asset available for version [%s] with name [%s]", version.String(), candidate)
	}
	if match.BrowserDownloadURL == nil {
		return nil, errors.Errorf("no asset url available in asset [%s]", candidate)
	}
	rsp, err := http.DefaultClient.Get(*match.BrowserDownloadURL)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to download asset from [%s]", match.BrowserDownloadURL)
	}
	defer func() { _ = rsp.Body.Close() }()
	return ioutil.ReadAll(rsp.Body)
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
		ret, err = ioutil.ReadAll(reader)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to read file [%s] from zip", f.Name)
		}
	}
	return ret, nil
}
