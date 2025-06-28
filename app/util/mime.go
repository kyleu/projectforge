package util

import (
	"mime"
	"strings"
)

func ExtensionFromMIME(mt string) string {
	if mt == "" {
		mt = "application/json"
	}
	if mt == "text/markdown" {
		return "md"
	}
	mts, _ := mime.ExtensionsByType(mt)
	if len(mts) == 0 {
		return mt
	}
	return strings.TrimPrefix(mts[0], ".")
}

func MIMEFromExtension(ext string) string {
	ext = strings.TrimPrefix(ext, ".")
	if ext == "" {
		ext = "txt"
	}
	if ext == "md" {
		return "text/markdown"
	}
	mt := mime.TypeByExtension(ext)
	if mt == "" {
		mt = "application/json"
	}
	return mt
}
