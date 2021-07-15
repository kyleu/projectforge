package svg

import (
	"os"

	"github.com/kyleu/projectforge/app/filesystem"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

func Load(src string, tgt string) (*SVG, error) {
	url := "https://raw.githubusercontent.com/icons8/line-awesome/master/svg/" + src + ".svg"
	var b []byte
	cl := &fasthttp.Client{}
	status, b, err := cl.Get(b, url)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to call URL [%s]", url)
	}
	if status != 200 {
		return nil, errors.Errorf("received status [%d] while calling URL [%s]", status, url)
	}

	return Transform(tgt, b)
}

func Save(src string, tgt string) (*SVG, error) {
	ret, err := Load(src, tgt)
	if err != nil {
		return nil, err
	}
	dst := "client/src/svg/" + tgt + ".svg"
	err = os.WriteFile(dst, []byte(ret.Markup), filesystem.DefaultMode)
	if err != nil {
		return nil, errors.Wrap(err, "unable to write SVG to file ["+tgt+".svg]")
	}
	return ret, nil
}
