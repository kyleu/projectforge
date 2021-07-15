source = ["./build/dist/darwin_arm64_darwin_arm64/$PF_EXECUTABLE$"]
bundle_id = "$PF_BUNDLE$"

apple_id {
  username = "$PF_AUTHOR_EMAIL$"
  password = "@env:APPLE_PASSWORD"
}

sign {
  application_identity = "$PF_SIGNING_IDENTITY$"
}

dmg {
  output_path = "./build/dist/$PF_EXECUTABLE$_0.0.0_macos_arm64.dmg"
  volume_name = "$PF_NAME$"
}

zip {
  output_path = "./build/dist/$PF_EXECUTABLE$_0.0.0_macos_arm64_notarized.zip"
}
