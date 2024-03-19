package svg

import (
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

const (
	svgPath       = "client/src/svg"
	ghLineAwesome = "https://raw.githubusercontent.com/icons8/line-awesome/master/svg/"
)

func AddToProject(fs filesystem.FileLoader, src string, tgt string) (*SVG, error) {
	ret, err := load(src, tgt)
	if err != nil {
		return nil, err
	}
	dst := filepath.Join(svgPath, tgt+util.ExtSVG)
	err = fs.WriteFile(dst, []byte(ret.Markup), filesystem.DefaultMode, true)
	if err != nil {
		return nil, errors.Wrap(err, "unable to write SVG to file ["+tgt+".svg]")
	}
	return ret, nil
}

func load(src string, tgt string) (*SVG, error) {
	url := src
	if !strings.HasPrefix(src, "http") {
		url = ghLineAwesome + src + util.ExtSVG
	}
	cl := http.DefaultClient
	rsp, err := cl.Get(url)
	if err != nil || rsp.StatusCode == 404 {
		if !strings.HasPrefix(src, "http") {
			origErr := err
			origURL := url
			url = ghLineAwesome + src + "-solid.svg"
			rsp, err = cl.Get(url)
			if err != nil {
				return nil, errors.Wrapf(origErr, "unable to call URL [%s]", origURL)
			}
		} else {
			return nil, errors.Wrapf(err, "unable to call URL [%s]: %d", url, rsp.StatusCode)
		}
	}
	if rsp.StatusCode != 200 {
		return nil, errors.Errorf("received status [%d] while calling URL [%s]", rsp.StatusCode, url)
	}
	b, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read response")
	}

	return Transform(tgt, b, url)
}
