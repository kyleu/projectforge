package module

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

func (s *Service) Download(_ context.Context, key string, url string, logger util.Logger) error {
	if url == "" {
		return errors.New("must provide URL")
	}
	logger.Infof("downloading module [%s] from URL [%s]", key, url)
	response, err := util.NewHTTPRequest(context.Background(), http.MethodGet, url).Run()
	if err != nil {
		return errors.Wrapf(err, "unable to retrieve module [%s] from [%s]", key, url)
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		return errors.Errorf("module [%s] load request to [%s] returned status code [%d]", key, url, response.StatusCode)
	}

	b, err := io.ReadAll(response.Body)
	if err != nil {
		return errors.Errorf("unable to read body from module [%s] load request to [%s]", key, url)
	}

	r, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return errors.Errorf("unable to unzip body from module [%s] load request to [%s]", key, url)
	}

	_ = s.config.RemoveRecursive(key, logger)
	for _, f := range r.File {
		fn := util.StringFilePath(key, f.Name)
		if f.FileInfo().IsDir() {
			err = s.config.CreateDirectory(fn)
			if err != nil {
				return errors.Errorf("unable to create directory [%s] for module [%s] from [%s]", fn, key, url)
			}
			continue
		}
		o, err := f.Open()
		if err != nil {
			return errors.Errorf("unable to read file [%s] from module [%s] from [%s]", fn, key, url)
		}
		content, err := io.ReadAll(o)
		if err != nil {
			return errors.Errorf("unable to read content of [%s] for module [%s] from [%s]", fn, key, url)
		}
		err = s.config.WriteFile(fn, content, filesystem.FileMode(f.Mode()), false)
		if err != nil {
			return errors.Errorf("unable to write file [%s] for module [%s] from [%s]", fn, key, url)
		}
	}

	return nil
}
