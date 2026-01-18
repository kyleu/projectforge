//go:build test_all || !func_test
// +build test_all !func_test

package filesystem_test

import (
	"testing"

	"projectforge.dev/projectforge/app/lib/filesystem"
)

func TestImageType(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    []string
		expected string
	}{
		{name: "empty path", input: []string{}, expected: ""},
		{name: "no extension", input: []string{"filename"}, expected: ""},
		{name: "bmp extension", input: []string{"image.bmp"}, expected: "bmp"},
		{name: "gif extension", input: []string{"image.gif"}, expected: "gif"},
		{name: "jpg extension", input: []string{"image.jpg"}, expected: "jpg"},
		{name: "jpeg extension", input: []string{"image.jpeg"}, expected: "jpeg"},
		{name: "png extension", input: []string{"image.png"}, expected: "png"},
		{name: "unsupported extension", input: []string{"document.pdf"}, expected: ""},
		{name: "multiple path segments", input: []string{"dir", "subdir", "image.png"}, expected: "png"},
		{name: "uppercase extension", input: []string{"IMAGE.PNG"}, expected: ""},
		{name: "double extension", input: []string{"file.tar.gz"}, expected: ""},
		{name: "hidden file with extension", input: []string{".hidden.jpg"}, expected: "jpg"},
	}

	for _, tc := range cases {
		c := tc
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			result := filesystem.ImageType(c.input...)
			if result != c.expected {
				t.Errorf("ImageType(%v) = %s, expected %s", c.input, result, c.expected)
			}
		})
	}
}

func TestExifExtract_NoExifData(t *testing.T) {
	t.Parallel()

	t.Run("returns error for non-image data", func(t *testing.T) {
		t.Parallel()
		data := []byte("not an image")
		_, err := filesystem.ExifExtract(data)
		if err == nil {
			t.Error("expected error for non-image data")
		}
	})

	t.Run("returns error for empty data", func(t *testing.T) {
		t.Parallel()
		data := []byte{}
		_, err := filesystem.ExifExtract(data)
		if err == nil {
			t.Error("expected error for empty data")
		}
	})
}
