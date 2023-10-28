// Package filesystem - Content managed by Project Forge, see [projectforge.md] for details.
package filesystem

import (
	"io"
	"net/http"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/util"
)

func (f *FileSystem) Download(url string, path string, overwrite bool, logger util.Logger) (int, error) {
	if f.Exists(path) && !overwrite {
		return 0, errors.Errorf("file [%s] exists", path)
	}
	w, err := f.FileWriter(path, true, false)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to open file at path [%s]", path)
	}
	rsp, err := http.Get(url)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to load url [%s]", url)
	}
	if rsp.StatusCode != 200 {
		return 0, errors.Errorf("response from url [%s] has status [%d], expected [200]", url, rsp.StatusCode)
	}
	x, err := io.Copy(w, rsp.Body)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to save response from url [%s] to file [%s]", url, path)
	}
	return int(x), nil
}
