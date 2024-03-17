# Content managed by Project Forge, see [projectforge.md] for details.
source = ["./build/dist/darwin_darwin_all/projectforge"]
bundle_id = "com.kyleu.projectforge"

notarize {
  path = "./build/dist/projectforge_1.1.2_darwin_all_desktop.dmg"
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
  output_path = "./build/dist/projectforge_1.1.2_darwin_all.dmg"
  volume_name = "Project Forge"
}

zip {
  output_path = "./build/dist/projectforge_1.1.2_darwin_all_notarized.zip"
}
