package controller

import (
	"os"

	"github.com/go-gem/sessions"
	"github.com/gorilla/securecookie"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"{{{ .Package}}}/app/controller/cutil"
	"{{{ .Package}}}/app/util"
	"{{{ .Package}}}/app/web"
)

var sessionKey = func() string {
	x := os.Getenv("SESSION_KEY")
	if x == "" {
		x = util.AppKey + "_random_secret_key"
	}
	return x
}()

var store *sessions.CookieStore

func initStore(keyPairs ...[]byte) *sessions.CookieStore {
	ret := sessions.NewCookieStore(keyPairs...)
	for _, x := range ret.Codecs {
		c, ok := x.(*securecookie.SecureCookie)
		if ok {
			c.MaxLength(65536)
		}
	}
	return ret
}

func loadPageState(ctx *fasthttp.RequestCtx, logger *zap.SugaredLogger) *cutil.PageState {
	logger = logger.With(zap.String("path", string(ctx.Request.URI().Path())))

	if store == nil {
		store = initStore([]byte(sessionKey))
	}
	session, err := store.Get(ctx, util.AppKey)
	if err != nil {
		logger.Warnf("error retrieving session: %+v", err)
		session, err = store.New(ctx, util.AppKey)
		if err != nil {
			logger.Warnf("error creating new session: %+v", err)
		}
	}
	flashes := util.StringArrayFromInterfaces(session.Flashes())
	if len(flashes) > 0 {
		err = web.SaveSession(ctx, session, logger)
		if err != nil {
			logger.Warnf("can't save session: %+v", err)
		}
	}

	prof, err := loadProfile(session)
	if err != nil {
		logger.Warnf("can't load profile: %+v", err)
	}

	var a web.Accounts
	authX, ok := session.Values["auth"]
	if ok {
		authS, ok := authX.(string)
		if ok {
			a = web.AccountsFromString(authS)
		}
	}

	return &cutil.PageState{
		Method:   string(ctx.Method()),
		URI:      ctx.Request.URI(),
		Flashes:  flashes,
		Session:  session,
		Profile:  prof,
		Accounts: a,
		Icons:    initialIcons,
		Logger:   logger,
	}
}
