package svg

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"projectforge.dev/projectforge/app/lib/filesystem"
)

const svgPath = "client/src/svg"

func AddToProject(fs filesystem.FileLoader, src string, tgt string) (*SVG, error) {
	ret, err := load(src, tgt)
	if err != nil {
		return nil, err
	}
	dst := filepath.Join(svgPath, tgt+".svg")
	err = fs.WriteFile(dst, []byte(ret.Markup), filesystem.DefaultMode, true)
	if err != nil {
		return nil, errors.Wrap(err, "unable to write SVG to file ["+tgt+".svg]")
	}
	return ret, nil
}

func load(src string, tgt string) (*SVG, error) {
	url := src
	if !strings.HasPrefix(src, "http") {
		url = "https://raw.githubusercontent.com/icons8/line-awesome/master/svg/" + src + ".svg"
	}
	var b []byte
	cl := &fasthttp.Client{}
	status, b, err := cl.Get(b, url)
	if err != nil {
		if !strings.HasPrefix(src, "http") {
			origErr := err
			origURL := url
			url = "https://raw.githubusercontent.com/icons8/line-awesome/master/svg/" + src + "-solid.svg"
			_, _, err = cl.Get(b, url)
			if err != nil {
				return nil, errors.Wrapf(origErr, "unable to call URL [%s]", origURL)
			}
		}
		return nil, errors.Wrapf(err, "unable to call URL [%s]: %d", url, status)
	}
	if status != 200 {
		return nil, errors.Errorf("received status [%d] while calling URL [%s]", status, url)
	}

	return Transform(tgt, b, url)
}
