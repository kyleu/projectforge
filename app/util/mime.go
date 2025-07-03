package util

import (
	"mime"
	"strings"
)

const (
	MIMETypeJSON     = "application/json"
	MIMETypeMarkdown = "text/markdown"
	extMarkdown      = "md"
)

func ExtensionFromMIME(mt string) string {
	if mt == "" {
		mt = MIMETypeJSON
	}
	if mt == MIMETypeMarkdown {
		return extMarkdown
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
	if ext == extMarkdown {
		return MIMETypeMarkdown
	}
	mt := mime.TypeByExtension(ext)
	if mt == "" {
		mt = MIMETypeJSON
	}
	return mt
}
