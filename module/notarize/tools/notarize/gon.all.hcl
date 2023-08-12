source = ["./build/dist/darwin_darwin_all/{{{ .Exec }}}"]
bundle_id = "{{{ .Info.Bundle }}}"{{{ if .BuildDesktop }}}

notarize {
  path = "./build/dist/{{{ .Exec }}}_{{{ .Version }}}_darwin_all_desktop.dmg"
  bundle_id = "{{{ .Info.Bundle }}}"
}{{{ end }}}

apple_id {
  username = "{{{ .Info.NotarizationEmail }}}"
  password = "@env:APPLE_PASSWORD"
}

sign {
  application_identity = "{{{ .Info.SigningIdentity }}}"
}

dmg {
  output_path = "./build/dist/{{{ .Exec }}}_{{{ .Version }}}_darwin_all.dmg"
  volume_name = "{{{ .Name }}}"
}

zip {
  output_path = "./build/dist/{{{ .Exec }}}_{{{ .Version }}}_darwin_all_notarized.zip"
}
