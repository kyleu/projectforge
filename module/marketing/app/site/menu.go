package site

import (
	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/menu"
	"{{{ .Package }}}/app/user"
)

const (
	keyInstall    = "install"
	keyDownload   = "download"
	keyQuickStart = "quickstart"
	keyContrib    = "contributing"
	keyTech       = "technology"
)

func Menu(as *app.State, _ *user.Profile, _ user.Accounts) menu.Items {
	return menu.Items{
		{Key: keyInstall, Title: "Install", Icon: "code", Route: "/" + keyInstall},
		{Key: keyDownload, Title: "Download", Icon: "download", Route: "/" + keyDownload},
		{Key: keyQuickStart, Title: "Quick Start", Icon: "bolt", Route: "/" + keyQuickStart},
		{Key: keyContrib, Title: "Contributing", Icon: "cog", Route: "/" + keyContrib},
		{Key: keyTech, Title: "Technology", Icon: "cog", Route: "/" + keyTech},
	}
}
