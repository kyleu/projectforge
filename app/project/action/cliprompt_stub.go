//go:build js
// +build js

package action

import "github.com/pkg/errors"

func promptString(_ string, _ string, _ string, _ ...bool) (string, error) {
	return "", errors.New("interactive prompts are unavailable on this platform")
}
