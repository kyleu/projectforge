source = ["./build/dist/darwin_amd64_darwin_amd64/$PF_EXECUTABLE$"]
bundle_id = "$PF_BUNDLE$"

//notarize {
//  path = "./build/dist/$PF_EXECUTABLE$_desktop_$PF_VERSION$_macos_x86_64.dmg"
//  bundle_id = "$PF_BUNDLE$"
//}

apple_id {
  username = "$PF_AUTHOR_EMAIL$"
  password = "@env:APPLE_PASSWORD"
}

sign {
  application_identity = "$PF_SIGNING_IDENTITY$"
}

dmg {
  output_path = "./build/dist/$PF_EXECUTABLE$_$PF_VERSION$_macos_x86_64.dmg"
  volume_name = "$PF_NAME$"
}

zip {
  output_path = "./build/dist/$PF_EXECUTABLE$_$PF_VERSION$_macos_x86_64_notarized.zip"
}
