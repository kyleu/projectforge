source = ["./build/dist/darwin_darwin_amd64/projectforge"]
bundle_id = "com.kyleu.projectforge"

//notarize {
//  path = "./build/dist/projectforge_0.1.15_macos_x86_64_desktop.dmg"
//  bundle_id = "com.kyleu.projectforge"
//}

apple_id {
  username = "kyle@kyleu.com"
  password = "@env:APPLE_PASSWORD"
}

sign {
  application_identity = "Developer ID Application: Kyle Unverferth (C6S478FYLD)"
}

dmg {
  output_path = "./build/dist/projectforge_0.1.15_macos_x86_64.dmg"
  volume_name = "Project Forge"
}

zip {
  output_path = "./build/dist/projectforge_0.1.15_macos_x86_64_notarized.zip"
}
