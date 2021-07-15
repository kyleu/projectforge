package site

import (
	"github.com/kyleu/projectforge/app/auth"
	"github.com/kyleu/projectforge/app/menu"
	"github.com/kyleu/projectforge/app/user"
)

const (
	keyIntro      = "intro"
	keyInstall    = "install"
	keyDownload   = "download"
	keyQuickStart = "quickstart"
	keyContrib    = "contrib"
)

func Menu(p *user.Profile, s auth.Sessions) menu.Items {
	return menu.Items{
		{Key: keyIntro, Title: "Introduction", Icon: "heart", Route: "/intro"},
		{Key: keyInstall, Title: "Install", Icon: "code", Route: "/install"},
		{Key: keyDownload, Title: "Download", Icon: "download", Route: "/download"},
		{Key: keyQuickStart, Title: "Quick Start", Icon: "bolt", Route: "/quickstart"},
		{Key: keyContrib, Title: "Contributing", Icon: "cog", Route: "/contrib"},
	}
}
