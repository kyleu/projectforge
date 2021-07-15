package auth

import (
	"os"

	"github.com/kyleu/projectforge/app/auth/msfix"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/amazon"
	"github.com/markbates/goth/providers/apple"
	"github.com/markbates/goth/providers/auth0"
	"github.com/markbates/goth/providers/azuread"
	"github.com/markbates/goth/providers/battlenet"
	"github.com/markbates/goth/providers/bitbucket"
	"github.com/markbates/goth/providers/box"
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
	"github.com/markbates/goth/providers/instagram"
	"github.com/markbates/goth/providers/intercom"
	"github.com/markbates/goth/providers/kakao"
	"github.com/markbates/goth/providers/lastfm"
	"github.com/markbates/goth/providers/line"
	"github.com/markbates/goth/providers/linkedin"
	"github.com/markbates/goth/providers/mastodon"
	"github.com/markbates/goth/providers/meetup"
	"github.com/markbates/goth/providers/microsoftonline"
	"github.com/markbates/goth/providers/naver"
	"github.com/markbates/goth/providers/nextcloud"
	"github.com/markbates/goth/providers/okta"
	"github.com/markbates/goth/providers/onedrive"
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
	"github.com/markbates/goth/providers/twitch"
	"github.com/markbates/goth/providers/twitter"
	"github.com/markbates/goth/providers/typetalk"
	"github.com/markbates/goth/providers/uber"
	"github.com/markbates/goth/providers/vk"
	"github.com/markbates/goth/providers/wepay"
	"github.com/markbates/goth/providers/xero"
	"github.com/markbates/goth/providers/yahoo"
	"github.com/markbates/goth/providers/yammer"
	"github.com/markbates/goth/providers/yandex"
	"github.com/pkg/errors"
)

// nolint
func toGoth(id string, k string, s string, c string) (goth.Provider, error) {
	switch id {
	case "amazon":
		return amazon.New(k, s, c), nil
	case "apple":
		return apple.New(k, s, c, nil, apple.ScopeName, apple.ScopeEmail), nil
	case auth0Key:
		return auth0.New(k, s, c, os.Getenv("auth0_domain")), nil
	case "azuread":
		return azuread.New(k, s, c, nil), nil
	case "battlenet":
		return battlenet.New(k, s, c), nil
	case "bitbucket":
		return bitbucket.New(k, s, c), nil
	case "box":
		return box.New(k, s, c), nil
	case "dailymotion":
		return dailymotion.New(k, s, c, "email"), nil
	case "deezer":
		return deezer.New(k, s, c, "email"), nil
	case "digitalocean":
		return digitalocean.New(k, s, c, "read"), nil
	case "discord":
		return discord.New(k, s, c, discord.ScopeIdentify, discord.ScopeEmail), nil
	case "dropbox":
		return dropbox.New(k, s, c), nil
	case "eveonline":
		return eveonline.New(k, s, c), nil
	case "facebook":
		return facebook.New(k, s, c), nil
	case "fitbit":
		return fitbit.New(k, s, c), nil
	case "gitea":
		return gitea.New(k, s, c), nil
	case "github":
		return github.New(k, s, c), nil
	case "gitlab":
		return gitlab.New(k, s, c), nil
	case "google":
		return google.New(k, s, c), nil
	case "gplus":
		return gplus.New(k, s, c), nil
	case "heroku":
		return heroku.New(k, s, c), nil
	case "instagram":
		return instagram.New(k, s, c), nil
	case "intercom":
		return intercom.New(k, s, c), nil
	case "kakao":
		return kakao.New(k, s, c), nil
	case "lastfm":
		return lastfm.New(k, s, c), nil
	case "line":
		return line.New(k, s, c, "profile", "openid", "email"), nil
	case "linkedin":
		return linkedin.New(k, s, c), nil
	case "mastodon":
		return mastodon.New(k, s, c, "read:accounts"), nil
	case "meetup":
		return meetup.New(k, s, c), nil
	case microsoftKey:
		return msfix.New(k, s, c, os.Getenv("microsoft_tenant")), nil
	case "microsoftonline":
		return microsoftonline.New(k, s, c), nil
	case "naver":
		return naver.New(k, s, c), nil
	case nextcloudKey:
		return nextcloud.NewCustomisedDNS(k, s, c, os.Getenv("nextcloud_url")), nil
	case "okta":
		return okta.New(k, s, c, "openid", "profile", "email"), nil
	case "onedrive":
		return onedrive.New(k, s, c), nil
	case "paypal":
		return paypal.New(k, s, c), nil
	case "salesforce":
		return salesforce.New(k, s, c), nil
	case "seatalk":
		return seatalk.New(k, s, c), nil
	case "shopify":
		return shopify.New(k, s, c, shopify.ScopeReadCustomers, shopify.ScopeReadOrders), nil
	case "slack":
		return slack.New(k, s, c), nil
	case "soundcloud":
		return soundcloud.New(k, s, c), nil
	case "spotify":
		return spotify.New(k, s, c), nil
	case "steam":
		return steam.New(k, c), nil
	case "strava":
		return strava.New(k, s, c), nil
	case "stripe":
		return stripe.New(k, s, c), nil
	case "twitch":
		return twitch.New(k, s, c), nil
	case "twitter":
		return twitter.New(k, s, c), nil
	case "typetalk":
		return typetalk.New(k, s, c, "my"), nil
	case "uber":
		return uber.New(k, s, c), nil
	case "vk":
		return vk.New(k, s, c), nil
	case "wepay":
		return wepay.New(k, s, c, "view_user"), nil
	case "xero":
		return xero.New(k, s, c), nil
	case "yahoo":
		return yahoo.New(k, s, c), nil
	case "yammer":
		return yammer.New(k, s, c), nil
	case "yandex":
		return yandex.New(k, s, c), nil
	default:
		return nil, errors.Errorf("invalid user provider [%s]", id)
	}
}
