//go:build test_all || !func_test
// +build test_all !func_test

package filesystem_test

import (
	"testing"

	"projectforge.dev/projectforge/app/lib/filesystem"
)

func TestMemFS(t *testing.T) {
	t.Parallel()

	t.Run("returns non-nil filesystem", func(t *testing.T) {
		t.Parallel()
		fs := filesystem.MemFS()
		if fs == nil {
			t.Error("expected non-nil filesystem")
		}
	})

	t.Run("returns same instance", func(t *testing.T) {
		t.Parallel()
		fs1 := filesystem.MemFS()
		fs2 := filesystem.MemFS()
		if fs1 != fs2 {
			t.Error("expected same instance")
		}
	})
}

func TestOSFS(t *testing.T) {
	t.Parallel()

	t.Run("returns non-nil filesystem", func(t *testing.T) {
		t.Parallel()
		fs := filesystem.OSFS()
		if fs == nil {
			t.Error("expected non-nil filesystem")
		}
	})

	t.Run("returns same instance", func(t *testing.T) {
		t.Parallel()
		fs1 := filesystem.OSFS()
		fs2 := filesystem.OSFS()
		if fs1 != fs2 {
			t.Error("expected same instance")
		}
	})
}

func TestNewFileSystem(t *testing.T) {
	t.Parallel()

	t.Run("creates filesystem with empty mode", func(t *testing.T) {
		t.Parallel()
		fs, err := filesystem.NewFileSystem("tmp", false, "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if fs == nil {
			t.Error("expected non-nil filesystem")
		}
	})

	t.Run("creates memory filesystem", func(t *testing.T) {
		t.Parallel()
		fs, err := filesystem.NewFileSystem("tmp", false, "memory")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if fs == nil {
			t.Error("expected non-nil filesystem")
		}
	})

	t.Run("creates file filesystem", func(t *testing.T) {
		t.Parallel()
		fs, err := filesystem.NewFileSystem("tmp", false, "file")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if fs == nil {
			t.Error("expected non-nil filesystem")
		}
	})

	t.Run("creates readonly filesystem", func(t *testing.T) {
		t.Parallel()
		fs, err := filesystem.NewFileSystem("tmp", true, "memory")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if fs == nil {
			t.Error("expected non-nil filesystem")
		}
	})

	t.Run("returns error for invalid mode", func(t *testing.T) {
		t.Parallel()
		_, err := filesystem.NewFileSystem("tmp", false, "invalid")
		if err == nil {
			t.Error("expected error for invalid mode")
		}
	})

	t.Run("strips ./ prefix from root", func(t *testing.T) {
		t.Parallel()
		fs, _ := filesystem.NewFileSystem("./tmp", false, "memory")
		if fs.Root() != "tmp" {
			t.Errorf("expected 'tmp', got '%s'", fs.Root())
		}
	})
}

func TestFileSystem_Root(t *testing.T) {
	t.Parallel()

	t.Run("returns root path", func(t *testing.T) {
		t.Parallel()
		fs, _ := filesystem.NewFileSystem("myroot", false, "memory")
		if fs.Root() != "myroot" {
			t.Errorf("expected 'myroot', got '%s'", fs.Root())
		}
	})
}

func TestFileSystem_Clone(t *testing.T) {
	t.Parallel()

	t.Run("returns new instance with same root", func(t *testing.T) {
		t.Parallel()
		fs, _ := filesystem.NewFileSystem("tmp", false, "memory")
		clone := fs.Clone()
		if clone == nil {
			t.Error("expected non-nil clone")
		}
	})
}

func TestFileSystem_String(t *testing.T) {
	t.Parallel()

	t.Run("formats correctly", func(t *testing.T) {
		t.Parallel()
		fs, _ := filesystem.NewFileSystem("myroot", false, "memory")
		result := fs.String()
		if result != "fs://myroot" {
			t.Errorf("expected 'fs://myroot', got '%s'", result)
		}
	})
}

func TestFileSystem_Exists(t *testing.T) {
	t.Parallel()

	t.Run("returns false for non-existent path", func(t *testing.T) {
		t.Parallel()
		fs, _ := filesystem.NewFileSystem("tmp", false, "memory")
		if fs.Exists("nonexistent") {
			t.Error("expected false for non-existent path")
		}
	})
}

func TestFileSystem_IsDir(t *testing.T) {
	t.Parallel()

	t.Run("returns false for non-existent path", func(t *testing.T) {
		t.Parallel()
		fs, _ := filesystem.NewFileSystem("tmp", false, "memory")
		if fs.IsDir("nonexistent") {
			t.Error("expected false for non-existent path")
		}
	})
}

func TestFileSystem_CreateDirectory(t *testing.T) {
	t.Parallel()

	t.Run("creates directory", func(t *testing.T) {
		t.Parallel()
		fs, _ := filesystem.NewFileSystem("tmp/test-create-dir", false, "memory")
		err := fs.CreateDirectory("subdir")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
