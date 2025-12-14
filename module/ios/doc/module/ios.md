# iOS

The **`ios`** module enables your application to build native iOS applications using a system WebView-based approach.
It provides the necessary templates, build scripts, and configuration to package your web application as a native iOS app.

## Overview

This module provides:

- **Native iOS App Template**: Complete Xcode project template with WebView integration
- **Build Automation**: Scripts to build and package iOS applications from your web app
- **Asset Management**: Icon and launch screen configuration through your application's configuration
- **WebView Bridge**: Seamless integration between your web application and iOS runtime

## Key Features

### WebView-Based Architecture
- Leverages native iOS WebView (WKWebView) for optimal performance
- Full access to your web application's functionality
- Native iOS chrome with web content

### Build Integration
- Automated build process via `./bin/ios.sh` script
- Integrates with the normal application build system
- Generates ready-to-submit Xcode projects

### Asset Configuration
- Configure app icons through the Project Forge UI or CLI
- Customize launch screens and app metadata
- Automatic asset generation for multiple iOS device sizes

### Development Workflow
- Test locally with iOS Simulator
- Debug with Safari Web Inspector
- Hot reload during development

## Requirements

### Build Prerequisites
- **macOS**: Required for iOS development and Xcode
- **Xcode**: Latest version with iOS SDK
- **iOS Build Option**: Must be enabled in your application's configuration

### Target Compatibility
- **iOS 12.0+**: Minimum supported iOS version
- **iPhone and iPad**: Universal app support
- **All device orientations**: Responsive design support

## Usage

### 1. Enable iOS Builds

In the Project Forge UI, enable "iOS" in build options, or add the "iOS" build option in your application's configuration file

### 2. Configure App Settings

Set up your iOS app metadata through your application's configuration or the Project Forge UI:
- App name and bundle identifier
- App icons (various sizes generated automatically)
- Launch screen configuration
- App Store metadata

### 3. Build the iOS App

```bash
# Build iOS application
./bin/ios.sh

# The script will:
# - Copy iOS template to build directory
# - Configure app settings and assets
# - Generate Xcode project
# - Build iOS app bundle
```

### 4. Development and Testing

```bash
# Open generated Xcode project
open ./build/dist/mobile_ios_app_arm64/YourApp.xcodeproj

# Or build from command line
cd ./build/dist/mobile_ios_app_arm64
xcodebuild -scheme YourApp -destination 'platform=iOS Simulator,name=iPhone 14'
```

## Project Structure

### Generated Files

After building, you'll find:

```
./build/dist/mobile_ios_app_arm64/
├── YourApp.xcodeproj/          # Xcode project file
├── YourApp/                    # iOS app source
│   ├── Info.plist             # App configuration
│   ├── AppDelegate.swift      # iOS app delegate
│   ├── ViewController.swift   # WebView controller
│   └── Assets.xcassets/       # App icons and images
└── YourAppTests/              # Unit test template
```

### Template Source

The iOS template is located at:

```
./tools/ios/
├── template.xcodeproj/        # Xcode project template
├── template/                  # iOS source templates
└── scripts/                   # Build automation scripts
```

## Configuration

### App Settings

Configure your iOS settings:

- **Bundle Identifier**: Unique app identifier (e.g., `com.yourcompany.yourapp`)
- **Display Name**: App name shown on iOS home screen

### WebView Configuration

The module automatically configures:
- **Local Server**: Starts your server on an available port
- **Network Access**: Allows external network requests
- **JavaScript**: Enabled with full API access
- **Storage**: Persistent local storage support

## Development Tips

### Testing
- Use iOS Simulator for rapid testing
- Test on physical devices for performance validation
- Verify network connectivity and offline behavior

### Debugging
- Enable Safari Web Inspector for debugging web content
- Use Xcode debugger for native iOS issues
- Monitor WebView console logs

## Limitations

### Build Dependencies
- Requires macOS for building and testing
- Xcode installation and setup required
- Apple Developer account needed for App Store distribution

## Troubleshooting

### Common Issues

**Build fails with Xcode errors:**
```bash
# Ensure Xcode command line tools are installed
sudo xcode-select --install

# Verify iOS build option is enabled in your application's configuration
```

**App crashes on launch:**
- Check that your web server is accessible from iOS app
- Verify network permissions in Info.plist
- Review WebView console logs

**Icons not showing:**
- Ensure icons are configured in your application's configuration
- Check that icon files are generated in Assets.xcassets
- Verify icon sizes meet iOS requirements
- Run "Rebuild SVG" in the Project Forge UI

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/ios
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
