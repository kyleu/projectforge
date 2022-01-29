// Package site $PF_IGNORE$
package site

import (
	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/lib/menu"
	"github.com/kyleu/projectforge/app/lib/user"
)

const (
	keyInstall     = "install"
	keyDownload    = "download"
	keyFeatures    = "features"
	keyContrib     = "contributing"
	keyCustomizing = "customizing"
	keyTech        = "technology"
)

func Menu(as *app.State, _ *user.Profile, _ user.Accounts) menu.Items {
	return menu.Items{
		{Key: keyInstall, Title: "Install", Icon: "code", Route: "/" + keyInstall},
		{Key: keyDownload, Title: "Download", Icon: "download", Route: "/" + keyDownload},
		{Key: keyFeatures, Title: "Features", Icon: "bolt", Route: "/" + keyFeatures, Children: featuresMenu(as)},
		{Key: keyCustomizing, Title: "Customizing", Icon: "code", Route: "/" + keyCustomizing},
		{Key: keyContrib, Title: "Contributing", Icon: "cog", Route: "/" + keyContrib},
		{Key: keyTech, Title: "Technology", Icon: "cog", Route: "/" + keyTech},
	}
}
