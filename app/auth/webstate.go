package auth

import (
	"encoding/base64"
	"net/url"

	"github.com/kyleu/projectforge/app/util"
	"github.com/markbates/goth"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

func setState(ctx *fasthttp.RequestCtx) string {
	state := ctx.Request.URI().QueryArgs().Peek("state")
	if len(state) > 0 {
		return string(state)
	}

	nonceBytes := util.RandomBytes(64)

	return base64.URLEncoding.EncodeToString(nonceBytes)
}

func validateState(ctx *fasthttp.RequestCtx, sess goth.Session) error {
	rawAuthURL, err := sess.GetAuthURL()
	if err != nil {
		return err
	}

	authURL, err := url.Parse(rawAuthURL)
	if err != nil {
		return err
	}

	originalState := authURL.Query().Get("state")
	qs := string(ctx.Request.URI().QueryArgs().Peek("state"))
	if originalState != "" && (originalState != qs) {
		return errors.New("state token mismatch")
	}
	return nil
}
