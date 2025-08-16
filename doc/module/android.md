# Android

The **`android`** module enables [Project Forge](https://projectforge.dev) applications to be deployed as native Android applications using a webview-based architecture. It provides a complete Android build pipeline that wraps your web application in a native Android container.

## Overview

This module provides:

- **Native Android App**: Complete Android project template with webview integration
- **Build Pipeline**: Automated build scripts using gomobile and Gradle
- **Cross-Platform**: Leverages your existing web application without code changes
- **Distribution Ready**: Generates APK files ready for Google Play Store or sideloading

## Key Features

### Native Integration
- Native Android webview container
- Full-screen web application experience
- Native navigation and lifecycle management
- Platform-specific optimizations

### Build System
- **gomobile** integration for Go-to-Android compilation
- **Gradle** build system for Android project management
- Automated AAR (Android Archive) generation
- APK packaging and signing support

### Deployment
- Debug and release build configurations
- ZIP packaging for easy distribution
- Android Studio project generation
- Compatible with Google Play Store requirements

## Architecture

The Android module creates a hybrid application architecture:

```
Android App Container
├── Native Android Activity
├── WebView Component
│   └── Your Project Forge Web App
├── Go Mobile Library (AAR)
│   └── Core Application Logic
└── Android Resources
    ├── Icons and Assets
    ├── Manifest Configuration
    └── Build Configuration
```

## Package Structure

### Build Infrastructure

- **`bin/build/android.sh`** - Main build script for Android applications
  - Compiles Go code to Android Archive (AAR) format
  - Sets up Android Studio project structure
  - Builds APK files for distribution

### Templates

- **`tools/android/`** - Complete Android Studio project template
  - Android manifest configuration
  - Gradle build files
  - WebView activity implementation
  - Resource files and icons

## Build Process

The Android build process follows these steps:

1. **Go Compilation**: Uses `gomobile bind` to compile Go code to Android AAR
2. **Project Setup**: Copies Android template to build directory
3. **Library Integration**: Includes compiled AAR in Android project
4. **APK Build**: Uses Gradle to build final APK file
5. **Packaging**: Creates distribution ZIP files

## Requirements

### Development Dependencies
- **Go 1.25+** with gomobile support
- **Android SDK** with API level 21+
- **Gradle** build system
- **Java Development Kit (JDK) 8+**

### Runtime Requirements
- **Android 5.0** (API level 21) or higher
- **WebView** component (standard on all Android devices)
- **Network connectivity** (for web application features)

## Configuration

### Android Manifest
The module provides a configurable Android manifest template supporting:
- Application name and package configuration
- Permission declarations
- Target SDK version settings
- Icon and theme customization

### Build Options
Configure builds through Project Forge settings:
- Target Android API level
- Application signing configuration
- Build variants (debug/release)
- Icon and splash screen assets

## Usage

### Prerequisites
Ensure the `android` build option is enabled in your Project Forge configuration.

### Building Android App

```bash
# Build Android application
./bin/build/android.sh [version]

# Example with version
./bin/build/android.sh 1.0.0
```

### Build Outputs

After successful build, find outputs in:
- **`./build/dist/mobile_android_arm64/`** - Android Studio project
- **`./build/dist/{app_name}_android_aar.zip`** - Go mobile library
- **`./build/dist/{app_name}_android_apk.zip`** - Installable APK files

### Development Workflow

1. **Development**: Develop and test your web application normally
2. **Android Build**: Run the Android build script
3. **Testing**: Install APK on device or emulator for testing
4. **Iteration**: Rebuild after making changes to your web app

## Customization

### Icons and Branding
- Replace icon files in the Android template
- Modify application name in manifest
- Configure splash screen and themes
- Set package name and version information

### Advanced Configuration
- Modify Gradle build scripts for custom dependencies
- Adjust WebView settings for performance
- Configure native Android features (permissions, services)
- Customize build variants and signing configurations

## Troubleshooting

### Common Issues

**Build Fails - Missing Dependencies**
```bash
# Ensure gomobile is installed
go install golang.org/x/mobile/cmd/gomobile@latest
gomobile init
```

**APK Install Fails**
- Check Android device/emulator API level (minimum API 21)
- Verify developer options and USB debugging enabled
- Ensure sufficient storage space

**WebView Not Loading**
- Review Android log output for errors

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/android
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)
