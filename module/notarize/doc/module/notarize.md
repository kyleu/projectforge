# Notarize

The **`notarize`** module provides Apple code signing and notarization capabilities for your application. It enables seamless distribution of macOS applications through Apple's notarization service, ensuring your applications meet Apple's security requirements for Gatekeeper.

## Overview

This module automates the Apple notarization process for macOS builds, providing:

- **Automated Notarization**: Submits built applications to Apple's notarization service
- **Multi-Architecture Support**: Handles Intel (amd64), Apple Silicon (arm64), and universal builds
- **DMG Distribution**: Notarizes DMG installers for easy distribution
- **Build Integration**: Seamlessly integrates with your application's build system

## Key Features

### Apple Platform Support
- Intel x64 (darwin_amd64) builds
- Apple Silicon ARM64 (darwin_arm64) builds
- Universal (darwin_all) builds combining both architectures
- DMG installer packages

### Security & Compliance
- Apple notarization service integration
- Gatekeeper compatibility
- Code signing certificate support
- Team ID and Apple ID authentication

### Developer Experience
- Automated submission process
- Build system integration
- Environment variable configuration
- Conditional execution for testing

## Requirements

### Apple Developer Account
- Valid Apple Developer account
- App-specific password for notarization
- Team ID for organization accounts
- Code signing certificate installed

### Environment Variables
- `apple_email` - Apple ID email address
- `apple_team_id` - Developer team identifier
- `apple_password` - App-specific password for notarization
- `publish_test` - Set to "true" to skip notarization during testing

### Project Configuration
- Project's `SigningIdentity` must be configured
- `notarize` build option must be enabled in project settings
- Valid code signing certificate in keychain

## Build Integration

The module provides a shell script (`notarize.sh`) that automatically submits built DMG files to Apple's notarization service:

```bash
# Submits three DMG variants for notarization:
# - Intel x64 build
# - Apple Silicon ARM64 build
# - Universal build (both architectures)
```

### Build Process
1. Application is built and signed with configured identity
2. DMG installers are created for each architecture
3. Notarization script submits DMGs to Apple's service
4. Apple validates and notarizes the applications
5. Email notifications confirm notarization status

## Usage Notes

### Performance Considerations
- Notarization typically takes 2-15 minutes per submission
- Apple sends email notifications for each submission
- Process requires active internet connection
- Large applications may take longer to process

### Best Practices
- Test builds locally before notarizing
- Monitor email for notarization status updates
- Keep Apple credentials secure and up-to-date

### Troubleshooting
- Verify Apple ID and team ID are correct
- Ensure app-specific password is current
- Check code signing certificate validity
- Review Apple's notarization logs for errors

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/notarize
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Apple Developer Documentation](https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution) - Official notarization guide
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
