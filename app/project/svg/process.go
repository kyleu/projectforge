package svg

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

const ghLineAwesome = "https://raw.githubusercontent.com/icons8/line-awesome/master/svg/"

func AddToProject(prj string, fs filesystem.FileLoader, src string, tgt string) (*SVG, error) {
	ret, err := load(src, tgt)
	if err != nil {
		return nil, err
	}
	dst := filepath.Join(svgPath(prj, tgt))
	err = fs.WriteFile(dst, []byte(ret.Markup), filesystem.DefaultMode, true)
	if err != nil {
		return nil, errors.Wrap(err, "unable to write SVG to file ["+tgt+".svg]")
	}
	return ret, nil
}

func svgPath(prj string, key string) string {
	if strings.Contains(key, "@") {
		x, y := util.StringSplit(key, '@', true)
		return fmt.Sprintf("%s%s/wwwroot/svg/%s%s", util.StringToTitle(prj), util.StringToTitle(x), y, util.ExtSVG)
	}
	return "client/src/svg/" + key + util.ExtSVG
}

func load(src string, tgt string) (*SVG, error) {
	url := src
	if !strings.HasPrefix(src, "http") {
		url = ghLineAwesome + src + util.ExtSVG
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	rsp, err := http.DefaultClient.Do(r)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil || rsp.StatusCode == 404 {
		if !strings.HasPrefix(src, "http") {
			origErr := err
			origURL := url
			url = ghLineAwesome + src + "-solid.svg"
			r, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
			if reqErr != nil {
				return nil, reqErr
			}
			rsp, reqErr = http.DefaultClient.Do(r)
			defer func() { _ = rsp.Body.Close() }()
			if reqErr != nil {
				return nil, errors.Wrapf(origErr, "unable to call URL [%s]", origURL)
			}
		} else {
			return nil, errors.Wrapf(err, "unable to call URL [%s]: %d", url, rsp.StatusCode)
		}
	}
	if rsp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("received status [%d] while calling URL [%s]", rsp.StatusCode, url)
	}
	b, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read response")
	}

	return Transform(tgt, b, url)
}
