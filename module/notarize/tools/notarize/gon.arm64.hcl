source = ["./build/dist/darwin_darwin_arm64/{{{ .Exec }}}"]
bundle_id = "{{{ .Info.Bundle }}}"

notarize {
  path = "./build/dist/{{{ .Exec }}}_{{{ .Version }}}_darwin_arm64_desktop.dmg"
  bundle_id = "{{{ .Info.Bundle }}}"
}

apple_id {
  username = "{{{ .Info.AuthorEmail }}}"
  password = "@env:APPLE_PASSWORD"
}

sign {
  application_identity = "{{{ .Info.SigningIdentity }}}"
}

dmg {
  output_path = "./build/dist/{{{ .Exec }}}_{{{ .Version }}}_darwin_arm64.dmg"
  volume_name = "{{{ .Name }}}"
}

zip {
  output_path = "./build/dist/{{{ .Exec }}}_{{{ .Version }}}_darwin_arm64_notarized.zip"
}
