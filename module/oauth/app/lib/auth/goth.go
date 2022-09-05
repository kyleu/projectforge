package auth

import (
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/amazon"
	"github.com/markbates/goth/providers/auth0"
	"github.com/markbates/goth/providers/azuread"
	"github.com/markbates/goth/providers/battlenet"
	"github.com/markbates/goth/providers/bitbucket"
	"github.com/markbates/goth/providers/box"
	"github.com/markbates/goth/providers/cloudfoundry"
	"github.com/markbates/goth/providers/dailymotion"
	"github.com/markbates/goth/providers/deezer"
	"github.com/markbates/goth/providers/digitalocean"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/dropbox"
	"github.com/markbates/goth/providers/eveonline"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/fitbit"
	"github.com/markbates/goth/providers/gitea"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gitlab"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/gplus"
	"github.com/markbates/goth/providers/heroku"
	"github.com/markbates/goth/providers/influxcloud"
	"github.com/markbates/goth/providers/instagram"
	"github.com/markbates/goth/providers/intercom"
	"github.com/markbates/goth/providers/kakao"
	"github.com/markbates/goth/providers/lastfm"
	"github.com/markbates/goth/providers/line"
	"github.com/markbates/goth/providers/linkedin"
	"github.com/markbates/goth/providers/mailru"
	"github.com/markbates/goth/providers/mastodon"
	"github.com/markbates/goth/providers/meetup"
	"github.com/markbates/goth/providers/microsoftonline"
	"github.com/markbates/goth/providers/naver"
	"github.com/markbates/goth/providers/nextcloud"
	"github.com/markbates/goth/providers/okta"
	"github.com/markbates/goth/providers/onedrive"
	"github.com/markbates/goth/providers/openidConnect"
	"github.com/markbates/goth/providers/oura"
	"github.com/markbates/goth/providers/paypal"
	"github.com/markbates/goth/providers/salesforce"
	"github.com/markbates/goth/providers/seatalk"
	"github.com/markbates/goth/providers/shopify"
	"github.com/markbates/goth/providers/slack"
	"github.com/markbates/goth/providers/soundcloud"
	"github.com/markbates/goth/providers/spotify"
	"github.com/markbates/goth/providers/steam"
	"github.com/markbates/goth/providers/strava"
	"github.com/markbates/goth/providers/stripe"
	"github.com/markbates/goth/providers/tumblr"
	"github.com/markbates/goth/providers/twitch"
	"github.com/markbates/goth/providers/twitter"
	"github.com/markbates/goth/providers/typetalk"
	"github.com/markbates/goth/providers/uber"
	"github.com/markbates/goth/providers/vk"
	"github.com/markbates/goth/providers/wecom"
	"github.com/markbates/goth/providers/wepay"
	"github.com/markbates/goth/providers/xero"
	"github.com/markbates/goth/providers/yahoo"
	"github.com/markbates/goth/providers/yammer"
	"github.com/markbates/goth/providers/yandex"
	"github.com/markbates/goth/providers/zoom"
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/auth/msfix"
	"{{{ .Package }}}/app/util"
)

//nolint:cyclop, funlen, gocyclo, maintidx
func toGoth(id string, k string, s string, c string, scopes ...string) (goth.Provider, error) {
	switch id {
	case "amazon":
		return amazon.New(k, s, c, scopes...), nil
	case auth0Key:
		return auth0.New(k, s, c, util.GetEnv("auth0_domain"), scopes...), nil
	case "azuread":
		return azuread.New(k, s, c, nil, scopes...), nil
	case "battlenet":
		return battlenet.New(k, s, c, scopes...), nil
	case "bitbucket":
		return bitbucket.New(k, s, c, scopes...), nil
	case "box":
		return box.New(k, s, c, scopes...), nil
	case "cloudfoundry":
		return cloudfoundry.New("", k, s, c, scopes...), nil
	case "dailymotion":
		return dailymotion.New(k, s, c, append([]string{"email"}, scopes...)...), nil
	case "deezer":
		return deezer.New(k, s, c, append([]string{"email"}, scopes...)...), nil
	case "digitalocean":
		return digitalocean.New(k, s, c, append([]string{"read"}, scopes...)...), nil
	case "discord":
		return discord.New(k, s, c, append([]string{discord.ScopeIdentify, discord.ScopeEmail}, scopes...)...), nil
	case "dropbox":
		return dropbox.New(k, s, c, scopes...), nil
	case "eveonline":
		return eveonline.New(k, s, c, scopes...), nil
	case "facebook":
		return facebook.New(k, s, c, scopes...), nil
	case "fitbit":
		return fitbit.New(k, s, c, scopes...), nil
	case "gitea":
		return gitea.New(k, s, c, scopes...), nil
	case "github":
		return github.New(k, s, c, scopes...), nil
	case "gitlab":
		return gitlab.New(k, s, c, scopes...), nil
	case "google":
		return google.New(k, s, c, scopes...), nil
	case "gplus":
		return gplus.New(k, s, c, scopes...), nil
	case "heroku":
		return heroku.New(k, s, c, scopes...), nil
	case "influxcloud":
		return influxcloud.New(k, s, c, scopes...), nil
	case "instagram":
		return instagram.New(k, s, c, scopes...), nil
	case "intercom":
		return intercom.New(k, s, c, scopes...), nil
	case "kakao":
		return kakao.New(k, s, c, scopes...), nil
	case "lastfm":
		return lastfm.New(k, s, c), nil
	case "line":
		return line.New(k, s, c, append([]string{"profile", "openid", "email"}, scopes...)...), nil
	case "linkedin":
		return linkedin.New(k, s, c, scopes...), nil
	case "mailru":
		return mailru.New(k, s, c, append([]string{"read:accounts"}, scopes...)...), nil
	case "mastodon":
		return mastodon.New(k, s, c, append([]string{"read:accounts"}, scopes...)...), nil
	case "meetup":
		return meetup.New(k, s, c, scopes...), nil
	case microsoftKey:
		return msfix.New(k, s, c, util.GetEnv("microsoft_tenant"), scopes...), nil
	case "microsoftonline":
		return microsoftonline.New(k, s, c, scopes...), nil
	case "naver":
		return naver.New(k, s, c), nil
	case nextcloudKey:
		return nextcloud.NewCustomisedDNS(k, s, c, util.GetEnv("nextcloud_url"), scopes...), nil
	case "okta":
		return okta.New(k, s, c, "openid", append([]string{"profile", "email"}, scopes...)...), nil
	case "onedrive":
		return onedrive.New(k, s, c, scopes...), nil
	case OpenIDConnectKey:
		return openidConnect.New(k, s, c, util.GetEnv("openid_connect_url"), append([]string{"profile", "email"}, scopes...)...)
	case "oura":
		return oura.New(k, s, c, scopes...), nil
	case "paypal":
		return paypal.New(k, s, c, scopes...), nil
	case "salesforce":
		return salesforce.New(k, s, c, scopes...), nil
	case "seatalk":
		return seatalk.New(k, s, c, scopes...), nil
	case "shopify":
		return shopify.New(k, s, c, append([]string{shopify.ScopeReadCustomers, shopify.ScopeReadOrders}, scopes...)...), nil
	case "slack":
		return slack.New(k, s, c, append([]string{slack.ScopeUserRead, "users:read.email"}, scopes...)...), nil
	case "soundcloud":
		return soundcloud.New(k, s, c, scopes...), nil
	case "spotify":
		return spotify.New(k, s, c, scopes...), nil
	case "steam":
		return steam.New(k, c), nil
	case "strava":
		return strava.New(k, s, c, scopes...), nil
	case "stripe":
		return stripe.New(k, s, c, scopes...), nil
	case "tumblr":
		return tumblr.New(k, s, c), nil
	case "twitch":
		return twitch.New(k, s, c, scopes...), nil
	case "twitter":
		return twitter.New(k, s, c), nil
	case "typetalk":
		return typetalk.New(k, s, c, append([]string{"my"}, scopes...)...), nil
	case "uber":
		return uber.New(k, s, c, scopes...), nil
	case "vk":
		return vk.New(k, s, c, scopes...), nil
	case "wecom":
		return wecom.New(k, s, util.GetEnv("wecom_agent_id"), c), nil
	case "wepay":
		return wepay.New(k, s, c, append([]string{"view_user"}, scopes...)...), nil
	case "xero":
		return xero.New(k, s, c), nil
	case "yahoo":
		return yahoo.New(k, s, c, scopes...), nil
	case "yammer":
		return yammer.New(k, s, c, scopes...), nil
	case "yandex":
		return yandex.New(k, s, c, scopes...), nil
	case "zoom":
		return zoom.New(k, s, c, scopes...), nil
	default:
		return nil, errors.Errorf("invalid user provider [%s]", id)
	}
}
