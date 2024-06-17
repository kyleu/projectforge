package site

import (
	"context"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/lib/menu"
	"projectforge.dev/projectforge/app/lib/user"
	"projectforge.dev/projectforge/app/util"
)

const (
	keyAbout       = "about"
	keyComponents  = "components"
	keyContrib     = "contributing"
	keyCustomizing = "customizing"
	keyDownload    = "download"
	keyFAQ         = "faq"
	keyFeatures    = "features"
	keyGallery     = "gallery"
	keyInstall     = "install"
	keyTech        = "technology"
)

func Menu(_ context.Context, as *app.State, _ *user.Profile, logger util.Logger) menu.Items {
	return menu.Items{
		{Key: keyGallery, Title: "Gallery", Icon: "flag", Route: "/" + keyGallery},
		{Key: keyInstall, Title: "Install", Icon: "code", Route: "/" + keyInstall},
		{Key: keyDownload, Title: "Download", Icon: "download", Route: "/" + keyDownload},
		{Key: keyFeatures, Title: "Features", Icon: "bolt", Route: "/" + keyFeatures, Children: featuresMenu(as.Services.Modules)},
		{Key: keyComponents, Title: "Components", Icon: "dna", Route: "/" + keyComponents, Children: componentsMenu(logger)},
		{Key: keyCustomizing, Title: "Customizing", Icon: "code", Route: "/" + keyCustomizing},
		{Key: keyContrib, Title: "Contributing", Icon: "cog", Route: "/" + keyContrib},
		{Key: keyTech, Title: "Technology", Icon: "shield", Route: "/" + keyTech},
		{Key: keyFAQ, Title: "FAQ", Icon: "book", Route: "/" + keyFAQ},
		{Key: keyAbout, Title: "About", Icon: "question", Route: "/" + keyAbout},
	}
}
