# Content managed by Project Forge, see [projectforge.md] for details.
source = ["./build/dist/darwin_darwin_amd64_v1/projectforge"]
bundle_id = "com.kyleu.projectforge"

notarize {
  path = "./build/dist/projectforge_0.11.10_darwin_amd64_desktop.dmg"
  bundle_id = "com.kyleu.projectforge"
}

apple_id {
  username = "kyle@kyleu.com"
  password = "@env:APPLE_PASSWORD"
}

sign {
  application_identity = "Developer ID Application: Kyle Unverferth (C6S478FYLD)"
}

dmg {
  output_path = "./build/dist/projectforge_0.11.10_darwin_amd64.dmg"
  volume_name = "Project Forge"
}

zip {
  output_path = "./build/dist/projectforge_0.11.10_darwin_amd64_notarized.zip"
}
