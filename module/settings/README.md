# [Settings] Project Forge Module

This directory contains the files used by the "settings" module of [Project Forge](https://projectforge.dev).

## Purpose

Provides a framework for managing file-backed application settings.

When enabled, this module:

- Persists a typed `Settings` struct to `settings.json` in the config directory
- Exposes admin UI pages at `/admin` and `/admin/settings`
- Uses `SettingsFieldDescs` to render editable fields in the UI

Customize `app/lib/settings/settings.go` to add or rename settings, and keep the field descriptors in sync so the admin editor renders correctly.

This module depends on the filesystem module for storage.

## Usage

When the "settings" module is enabled in your project, all of the files in this directory (except this readme) will be processed and included in your application.

See [doc/module/settings.md](doc/module/settings.md) for usage information.
