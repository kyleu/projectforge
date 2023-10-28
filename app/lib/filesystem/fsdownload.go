// Package filesystem - Content managed by Project Forge, see [projectforge.md] for details.
package filesystem

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/util"
)

func (f *FileSystem) Download(ctx context.Context, url string, path string, overwrite bool, _ util.Logger) (int, error) {
	if !strings.HasPrefix(url, "https://") {
		return 0, errors.New("only [https] requests are supported")
	}
	if f.Exists(path) && !overwrite {
		return 0, errors.Errorf("file [%s] exists", path)
	}
	w, err := f.FileWriter(path, true, false)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to open file at path [%s]", path)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to make request for url [%s]", url)
	}
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to load url [%s]", url)
	}
	defer func() {
		_ = rsp.Body.Close()
	}()
	if rsp.StatusCode != http.StatusOK {
		return 0, errors.Errorf("response from url [%s] has status [%d], expected [200]", url, rsp.StatusCode)
	}
	x, err := io.Copy(w, rsp.Body)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to save response from url [%s] to file [%s]", url, path)
	}
	return int(x), nil
}
