//go:build test_all || !func_test
// +build test_all !func_test

package filesystem_test

import (
	"testing"

	"projectforge.dev/projectforge/app/lib/filesystem"
)

func TestFileMode_ToFS(t *testing.T) {
	t.Parallel()

	t.Run("converts correctly", func(t *testing.T) {
		t.Parallel()
		mode := filesystem.FileMode(0o755)
		fsMode := mode.ToFS()
		if fsMode != 0o755 {
			t.Errorf("expected 0755, got %o", fsMode)
		}
	})
}

func TestFileInfo_Equal(t *testing.T) {
	t.Parallel()

	t.Run("equal files", func(t *testing.T) {
		t.Parallel()
		f1 := &filesystem.FileInfo{Name: "test.txt", Size: 100, Mode: 0o644, IsDir: false}
		f2 := &filesystem.FileInfo{Name: "test.txt", Size: 100, Mode: 0o644, IsDir: false}
		if !f1.Equal(f2) {
			t.Error("expected files to be equal")
		}
	})

	t.Run("different name", func(t *testing.T) {
		t.Parallel()
		f1 := &filesystem.FileInfo{Name: "test.txt", Size: 100, Mode: 0o644, IsDir: false}
		f2 := &filesystem.FileInfo{Name: "other.txt", Size: 100, Mode: 0o644, IsDir: false}
		if f1.Equal(f2) {
			t.Error("expected files to be different")
		}
	})

	t.Run("different size", func(t *testing.T) {
		t.Parallel()
		f1 := &filesystem.FileInfo{Name: "test.txt", Size: 100, Mode: 0o644, IsDir: false}
		f2 := &filesystem.FileInfo{Name: "test.txt", Size: 200, Mode: 0o644, IsDir: false}
		if f1.Equal(f2) {
			t.Error("expected files to be different")
		}
	})

	t.Run("different mode", func(t *testing.T) {
		t.Parallel()
		f1 := &filesystem.FileInfo{Name: "test.txt", Size: 100, Mode: 0o644, IsDir: false}
		f2 := &filesystem.FileInfo{Name: "test.txt", Size: 100, Mode: 0o755, IsDir: false}
		if f1.Equal(f2) {
			t.Error("expected files to be different")
		}
	})

	t.Run("different IsDir", func(t *testing.T) {
		t.Parallel()
		f1 := &filesystem.FileInfo{Name: "test", Size: 100, Mode: 0o755, IsDir: false}
		f2 := &filesystem.FileInfo{Name: "test", Size: 100, Mode: 0o755, IsDir: true}
		if f1.Equal(f2) {
			t.Error("expected files to be different")
		}
	})
}

func TestFileInfo_QueryEscapedPath(t *testing.T) {
	t.Parallel()

	t.Run("simple path", func(t *testing.T) {
		t.Parallel()
		f := &filesystem.FileInfo{Name: "file.txt"}
		result := f.QueryEscapedPath("dir1", "dir2")
		expected := "/dir1/dir2/file.txt"
		if result != expected {
			t.Errorf("expected %s, got %s", expected, result)
		}
	})

	t.Run("path with spaces", func(t *testing.T) {
		t.Parallel()
		f := &filesystem.FileInfo{Name: "my file.txt"}
		result := f.QueryEscapedPath("my dir")
		expected := "/my+dir/my+file.txt"
		if result != expected {
			t.Errorf("expected %s, got %s", expected, result)
		}
	})

	t.Run("no parent directories", func(t *testing.T) {
		t.Parallel()
		f := &filesystem.FileInfo{Name: "file.txt"}
		result := f.QueryEscapedPath()
		expected := "/file.txt"
		if result != expected {
			t.Errorf("expected %s, got %s", expected, result)
		}
	})
}

func TestFileInfos_Sorted(t *testing.T) {
	t.Parallel()

	t.Run("sorts by name", func(t *testing.T) {
		t.Parallel()
		infos := filesystem.FileInfos{
			{Name: "c.txt"},
			{Name: "a.txt"},
			{Name: "b.txt"},
		}
		sorted := infos.Sorted()
		if sorted[0].Name != "a.txt" || sorted[1].Name != "b.txt" || sorted[2].Name != "c.txt" {
			t.Errorf("expected sorted order, got %v", sorted)
		}
	})
}

func TestFileInfos_Equal(t *testing.T) {
	t.Parallel()

	t.Run("equal slices", func(t *testing.T) {
		t.Parallel()
		f1 := filesystem.FileInfos{{Name: "a.txt", Size: 100}, {Name: "b.txt", Size: 200}}
		f2 := filesystem.FileInfos{{Name: "a.txt", Size: 100}, {Name: "b.txt", Size: 200}}
		if !f1.Equal(f2) {
			t.Error("expected slices to be equal")
		}
	})

	t.Run("different lengths", func(t *testing.T) {
		t.Parallel()
		f1 := filesystem.FileInfos{{Name: "a.txt"}}
		f2 := filesystem.FileInfos{{Name: "a.txt"}, {Name: "b.txt"}}
		if f1.Equal(f2) {
			t.Error("expected slices to be different")
		}
	})

	t.Run("different content", func(t *testing.T) {
		t.Parallel()
		f1 := filesystem.FileInfos{{Name: "a.txt", Size: 100}}
		f2 := filesystem.FileInfos{{Name: "a.txt", Size: 200}}
		if f1.Equal(f2) {
			t.Error("expected slices to be different")
		}
	})
}

func TestDefaultModes(t *testing.T) {
	t.Parallel()

	t.Run("DirectoryMode is 755", func(t *testing.T) {
		t.Parallel()
		if filesystem.DirectoryMode != 0o755 {
			t.Errorf("expected 0755, got %o", filesystem.DirectoryMode)
		}
	})

	t.Run("DefaultMode is 644", func(t *testing.T) {
		t.Parallel()
		if filesystem.DefaultMode != 0o644 {
			t.Errorf("expected 0644, got %o", filesystem.DefaultMode)
		}
	})
}
