package filesystem_test

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/filesystem"
	"{{{ .Package }}}/app/lib/log"
)

const testFile = "foo.txt"

func TestFileSystem(t *testing.T) {
	t.Parallel()
	fixWD()
	fs, err := filesystem.NewFileSystem("tmp", false, "")
	if err != nil {
		t.Errorf("error starting filesystem [tmp]: %+v", err)
	}
	err = testFS("test-filesystem", fs, true)
	if err != nil {
		t.Errorf("error testing filesystem [%s]: %+v", fs.String(), err)
	}
}

func testFS(testDir string, fs filesystem.FileLoader, removeWhenDone bool) error {
	fixWD()
	logger, _ := log.CreateTestLogger()
	content := "Hello, test world!"
	if fs.Exists(testDir) {
		err := fs.RemoveRecursive(testDir, logger)
		if err != nil {
			return errors.Wrapf(err, "error removing test directory [%s] before starting", testDir)
		}
	}

	if err := fs.CreateDirectory(testDir); err != nil {
		return err
	}

	if err := fs.WriteFile(path.Join(testDir, testFile), []byte(content), filesystem.DefaultMode, false); err != nil {
		return err
	}

	if b, err := fs.ReadFile(path.Join(testDir, testFile)); err != nil || string(b) != content {
		if err != nil {
			return err
		}
		return errors.Errorf("content [%s] didn't match [%s]", string(b), content)
	}

	files := fs.ListFiles(testDir, nil, logger)
	if len(files) != 1 {
		return errors.Errorf("expected [%d] files, observed [%d]", 1, len(files))
	}
	if files[0].Name != testFile {
		return errors.Errorf("expected [%s] filename, observed [%s]", testFile, files[0].Name)
	}
	if files[0].Size != 18 {
		return errors.Errorf("expected [%d] file size, observed [%d]", 18, files[0].Size)
	}

	if err := fs.Remove(path.Join(testDir, testFile), logger); err != nil {
		return err
	}

	files = fs.ListFiles(testDir, nil, logger)
	if len(files) != 0 {
		return errors.Errorf("expected [%d] files, observed [%d]", 0, len(files))
	}

	if removeWhenDone {
		if err := fs.RemoveRecursive(testDir, logger); err != nil {
			return errors.Wrapf(err, "error removing test directory [%s] after completion", testDir)
		}
	}

	return nil
}

func fixWD() {
	wd, _ := os.Getwd()
	origwd := wd
	wd = strings.TrimSuffix(wd, "/filesystem")
	wd = strings.TrimSuffix(wd, "/lib")
	wd = strings.TrimSuffix(wd, "/app")
	if wd != origwd {
		_ = os.Chdir(wd)
	}
}
