source = ["./build/dist/darwin_darwin_amd64_v1/{{{ .Exec }}}"]
bundle_id = "{{{ .Info.Bundle }}}"

//notarize {
//  path = "./build/dist/{{{ .Exec }}}_{{{ .Version }}}_darwin_amd64_desktop.dmg"
//  bundle_id = "{{{ .Info.Bundle }}}"
//}

apple_id {
  username = "{{{ .Info.AuthorEmail }}}"
  password = "@env:APPLE_PASSWORD"
}

sign {
  application_identity = "{{{ .Info.SigningIdentity }}}"
}

dmg {
  output_path = "./build/dist/{{{ .Exec }}}_{{{ .Version }}}_darwin_amd64.dmg"
  volume_name = "{{{ .Name }}}"
}

zip {
  output_path = "./build/dist/{{{ .Exec }}}_{{{ .Version }}}_darwin_amd64_notarized.zip"
}
