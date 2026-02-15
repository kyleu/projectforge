//go:build js || aix
// +build js aix

package action

import "github.com/pkg/errors"

func promptString(_ string, _ string, _ string) (string, error) {
	return "", errors.New("interactive prompts are unavailable on this platform")
}
