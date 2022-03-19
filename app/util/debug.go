// Content managed by Project Forge, see [projectforge.md] for details.
package util

import (
	"runtime"
	"runtime/debug"
)

type DebugInfo struct {
	Tags *OrderedMap[string]
	Mods []*debug.Module
}

func GetDebugInfo() (*DebugInfo, error) {
	tags := NewOrderedMap[string](false, 10)
	tags.Append("Go Version", runtime.Version())
	tags.Append("CPU Architecture", runtime.GOARCH)
	tags.Append("Operating System", runtime.GOOS)

	mods, ok := debug.ReadBuildInfo()
	if !ok {
		return nil, nil
	}
	return &DebugInfo{Tags: tags, Mods: mods.Deps}, nil
}
