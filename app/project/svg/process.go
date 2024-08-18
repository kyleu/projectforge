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
	dst := svgPath(prj, ret.Key)
	err = fs.WriteFile(dst, []byte(ret.Markup), filesystem.DefaultMode, true)
	if err != nil {
		return nil, errors.Wrap(err, "unable to write SVG to file ["+ret.Key+".svg]")
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
	var url string
	test := func(u string) ([]byte, error) {
		url = u
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		rsp, b, err := util.NewHTTPRequest(ctx, http.MethodGet, u).RunSimple()
		return b, errors.Wrapf(err, "unable to call URL [%s]: %d", u, rsp.StatusCode)
	}

	if strings.HasPrefix(src, "http") {
		b, err := test(src)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}
		return Transform(tgt, b, src)
	}

	b, err := test(ghLineAwesome + src + util.ExtSVG)
	if err != nil {
		const suffix = "-solid.svg"
		b, err = test(ghLineAwesome + src + suffix)
		if err != nil {
			b, err = test(ghLineAwesome + src + suffix)
			if err != nil {
				return nil, err
			}
		}
	}
	return Transform(tgt, b, url)
}
