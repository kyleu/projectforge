//go:build android
package util

import (
	_ "golang.org/x/mobile/geom" // This file is only used to keep mobile from being evicted.
)

func _() {}
