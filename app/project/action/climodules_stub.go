//go:build js || aix
// +build js aix

package action

import (
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/module"
)

func promptModules(_ []string, _ module.Modules) ([]string, error) {
	return nil, errors.New("interactive prompts are unavailable on this platform")
}
