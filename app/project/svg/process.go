package svg

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

const (
	ghLineAwesome = "https://raw.githubusercontent.com/icons8/line-awesome/master/svg/"
	svgRoot       = "client/src/svg/"
)

func AddToProject(prj string, fs filesystem.FileLoader, src string, tgt string) (*SVG, error) {
	ret, err := load(src, tgt)
	if err != nil {
		return nil, err
	}
	dst := svgPath(prj, tgt)
	err = fs.WriteFile(dst, []byte(ret.Markup), filesystem.DefaultMode, true)
	if err != nil {
		return nil, errors.Wrap(err, "unable to write SVG to file ["+tgt+".svg]")
	}
	return ret, nil
}

func svgPath(prj string, key string) string {
	if strings.Contains(key, "@") {
		x, y := util.StringSplit(key, '@', true)
		return fmt.Sprintf("%s%s/wwwroot/svg/%s%s", util.StringToTitle(prj), util.StringToTitle(y), x, util.ExtSVG)
	}
	return svgRoot + key + util.ExtSVG
}

func load(src string, tgt string) (*SVG, error) {
	tgt, _ = util.StringSplit(tgt, '@', true)
	url := src
	if !strings.HasPrefix(src, "http") {
		url = ghLineAwesome + src + util.ExtSVG
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	rsp, b, err := util.NewHTTPRequest(ctx, http.MethodGet, url).RunSimple()
	if err != nil {
		if !strings.HasPrefix(src, "http") {
			origErr := err
			origURL := url
			url = ghLineAwesome + src + "-solid.svg"
			rsp, b, err = util.NewHTTPRequest(ctx, http.MethodGet, url).RunSimple()
			if err != nil {
				return nil, errors.Wrapf(origErr, "unable to call URL [%s]", origURL)
			}
		} else {
			return nil, errors.Wrapf(err, "unable to call URL [%s]: %d", url, rsp.StatusCode)
		}
	}
	return Transform(tgt, b, url)
}
