source = ["./build/dist/darwin_amd64_darwin_amd64/{{{ .Exec }}}"]
bundle_id = "{{{ .Bundle }}}"

//notarize {
//  path = "./build/dist/{{{ .Exec }}}_desktop_{{{ .Version }}}_macos_x86_64.dmg"
//  bundle_id = "{{{ .Bundle }}}"
//}

apple_id {
  username = "{{{ .AuthorEmail }}}"
  password = "@env:APPLE_PASSWORD"
}

sign {
  application_identity = "{{{ .SigningIdentity }}}"
}

dmg {
  output_path = "./build/dist/{{{ .Exec }}}_{{{ .Version }}}_macos_x86_64.dmg"
  volume_name = "{{{ .Name }}}"
}

zip {
  output_path = "./build/dist/{{{ .Exec }}}_{{{ .Version }}}_macos_x86_64_notarized.zip"
}
