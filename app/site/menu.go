// Package site $PF_IGNORE$
package site

import (
	"context"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/lib/menu"
	"projectforge.dev/projectforge/app/lib/user"
)

const (
	keyInstall     = "install"
	keyDownload    = "download"
	keyFeatures    = "features"
	keyContrib     = "contributing"
	keyCustomizing = "customizing"
	keyTech        = "technology"
	keyFAQ         = "faq"
)

func Menu(ctx context.Context, as *app.State, _ *user.Profile, _ user.Accounts) menu.Items {
	return menu.Items{
		{Key: keyInstall, Title: "Install", Icon: "code", Route: "/" + keyInstall},
		{Key: keyDownload, Title: "Download", Icon: "download", Route: "/" + keyDownload},
		{Key: keyFeatures, Title: "Features", Icon: "bolt", Route: "/" + keyFeatures, Children: featuresMenu(ctx, as.Services.Modules)},
		{Key: keyCustomizing, Title: "Customizing", Icon: "code", Route: "/" + keyCustomizing},
		{Key: keyContrib, Title: "Contributing", Icon: "cog", Route: "/" + keyContrib},
		{Key: keyTech, Title: "Technology", Icon: "shield", Route: "/" + keyTech},
		{Key: keyFAQ, Title: "FAQ", Icon: "question", Route: "/" + keyFAQ},
	}
}
