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
	simpleIcons   = "https://raw.githubusercontent.com/simple-icons/simple-icons/develop/icons/"
	svgRoot       = "client/src/svg/"
)

func AddToProject(ctx context.Context, prj string, fs filesystem.FileLoader, src string, tgt string) (*SVG, error) {
	ret, err := load(ctx, src, tgt)
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

func AddContentToProject(prj string, fs filesystem.FileLoader, tgt string, content string, srcURL string) (*SVG, error) {
	dst := svgPath(prj, tgt)
	if tgt == "" {
		return nil, errors.New("svg key may not be empty")
	}
	ret, err := Transform(tgt, []byte(content), srcURL)
	if err != nil {
		return nil, err
	}
	if ret.Markup == "" {
		return nil, errors.New("svg markup may not be empty")
	}
	err = fs.WriteFile(dst, []byte(ret.Markup), filesystem.DefaultMode, true)
	if err != nil {
		return nil, errors.Wrap(err, "unable to write SVG content to file ["+tgt+".svg]")
	}
	return ret, nil
}

func svgPath(prj string, key string) string {
	if strings.Contains(key, "@") {
		x, y := util.StringCut(key, '@', true)
		return fmt.Sprintf("%s%s/wwwroot/svg/%s%s", util.StringToTitle(prj), util.StringToTitle(y), x, util.ExtSVG)
	}
	return svgRoot + key + util.ExtSVG
}

func load(ctx context.Context, src string, tgt string) (*SVG, error) {
	tgt, _ = util.StringCut(tgt, '@', true)
	test := func(u string) ([]byte, error) {
		ctx, cancel := context.WithTimeout(ctx, time.Minute)
		defer cancel()
		st, _, b, err := util.NewHTTPRequest(ctx, http.MethodGet, u).RunSimple()
		return b, errors.Wrapf(err, "unable to call URL [%s]: %d", u, st)
	}

	get := func(u string) (*SVG, error) {
		if !strings.HasSuffix(u, ".svg") {
			u += util.ExtSVG
		}
		b, err := test(u)
		if err != nil {
			return nil, err
		}
		return Transform(tgt, b, u)
	}

	if strings.HasPrefix(src, "http") {
		b, err := test(src)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}
		return Transform(tgt, b, src)
	}

	if strings.HasPrefix(src, "brand-") {
		src = strings.TrimPrefix(src, "brand-")
		tgt = strings.TrimPrefix(tgt, "brand-")
		return get(simpleIcons + src)
	}

	if ret, err := get(ghLineAwesome + src + "-solid.svg"); err == nil {
		return ret, nil
	}

	return get(ghLineAwesome + src)
}
