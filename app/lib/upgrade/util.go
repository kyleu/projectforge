// Package upgrade - Content managed by Project Forge, see [projectforge.md] for details.
package upgrade

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/util"
)

func parseSource() (string, string, error) {
	remainder := util.AppSource
	idx := strings.LastIndex(remainder, "/")
	if idx == -1 {
		return "", "", errors.Errorf("invalid app source [%s]", remainder)
	}
	repo := remainder[idx+1:]
	remainder = remainder[:idx]
	idx = strings.LastIndex(remainder, "/")
	if idx == -1 {
		return "", "", errors.Errorf("no org provided in source [%s]", remainder)
	}
	org := remainder[idx+1:]
	return org, repo, nil
}

func overwrite(content []byte) error {
	cmdFN, err := os.Executable()
	if err != nil {
		return err
	}

	stat, err := os.Lstat(cmdFN)
	if err != nil {
		return errors.Wrapf(err, "unable to read file [%s]", cmdFN)
	}
	if stat.Mode()&os.ModeSymlink != 0 {
		var p string
		p, err = filepath.EvalSymlinks(cmdFN)
		if err != nil {
			return errors.Wrapf(err, "unable to resolve symlink [%s] for executable [%s]", p, cmdFN)
		}
		cmdFN = p
	}

	dir := filepath.Dir(cmdFN)
	fn := filepath.Base(cmdFN)
	newFN := filepath.Join(dir, fmt.Sprintf("%s.new", fn))
	newFile, err := os.OpenFile(newFN, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, stat.Mode())
	if err != nil {
		return errors.Wrapf(err, "unable to create file [%s]", newFN)
	}
	defer func() { _ = newFile.Close() }()
	_, err = io.Copy(newFile, bytes.NewReader(content))
	if err != nil {
		return err
	}
	_ = newFile.Close()

	oldFN := filepath.Join(dir, fmt.Sprintf("%s.old", fn))
	_ = os.Remove(oldFN)

	err = os.Rename(cmdFN, oldFN)
	if err != nil {
		return errors.Wrapf(err, "unable to rename file [%s] to [%s]", cmdFN, oldFN)
	}

	err = os.Rename(newFN, cmdFN)
	if err != nil {
		_ = os.Rename(oldFN, cmdFN)
		return err
	}

	_ = os.Remove(oldFN)

	return nil
}
