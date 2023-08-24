package controller

import (
	"bufio"
	"bytes"

	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/views"
)

func Testbed(rc *fasthttp.RequestCtx) {
	Act("testbed", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret := "Testbed!"
		ps.Data = ret
		raw := `GET /logo.svg HTTP/1.1
Host: localhost:40009
accept: image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8
sec-ch-ua: "Chromium";v="116", "Not)A;Brand";v="24", "Google Chrome";v="116"
sec-ch-ua-mobile: ?0
sec-ch-ua-platform: "macOS"
user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36

`
		req := fasthttp.Request{}
		err := req.Read(bufio.NewReader(bytes.NewReader([]byte(raw))))
		if err != nil {
			return "", err
		}

		return Render(rc, as, &views.Testbed{Param: ret}, ps)
	})
}
