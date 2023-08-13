package module

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/file/diff"
	"projectforge.dev/projectforge/app/util"
)

func (s *Service) UpdateFile(mods Modules, d *diff.Diff, logger util.Logger) ([]string, error) {
	var ret []string
	for _, mod := range mods {
		loader := s.GetFilesystem(mod.Key)
		if !loader.Exists(d.Path) {
			continue
		}
		mode, b, err := fileContent(loader, d.Path)
		if err != nil {
			return nil, err
		}

		newContent, err := reverseDiff(string(b), d, logger)
		if err != nil {
			return nil, err
		}
		if len(newContent) > 0 {
			if bytes.Equal(b, newContent) {
				ret = append(ret, fmt.Sprintf("no changes required to [%s] for module [%s]", d.Path, mod.Key))
			} else {
				_ = mode
				return nil, errors.Errorf("unable to merge change into module")
				//err = loader.WriteFile(d.Path, newContent, mode, true)
				//if err != nil {
				//	return nil, err
				//}
				//ret = append(ret, fmt.Sprintf("wrote [%d] bytes to [%s] for module [%s]", len(newContent), d.Path, mod.Key))
			}
		}
	}
	return ret, nil
}

func reverseDiff(dest string, d *diff.Diff, logger util.Logger) ([]byte, error) {
	logger.Debugf("reversing [%d] changes for file [%s]", len(d.Changes), d.Path)
	for _, ch := range d.Changes {
		var preCtx, postCtx, addedCtx, deletedCtx []string
		var isPost, hasTemplate bool
		for _, chLine := range ch.Lines {
			switch chLine.T {
			case "context":
				if strings.Contains(chLine.V, "{{{") {
					hasTemplate = true
				}
				if isPost {
					postCtx = append(postCtx, chLine.V)
				} else {
					preCtx = append(preCtx, chLine.V)
				}
			case "added":
				isPost = true
				addedCtx = append(addedCtx, chLine.V)
			case "deleted":
				isPost = true
				deletedCtx = append(deletedCtx, chLine.V)
			default:
				return nil, errors.Errorf("unable to handle change with type [%s]", chLine.T)
			}
		}

		pre, _, post := strings.Join(preCtx, ""), strings.Join(postCtx, ""), strings.Join(addedCtx, "")
		deleted := strings.Join(deletedCtx, "")

		preIdx := strings.Index(dest, pre)
		if preIdx == -1 {
			if hasTemplate {
				return nil, errors.New("module file with template does not contain pre-content context lines")
			}
			return nil, errors.New("module file does not contain pre-content context lines")
		}
		preIdx += len(pre)
		postIdx := preIdx + strings.Index(dest[preIdx:], post)
		if postIdx == -1 {
			if hasTemplate {
				return nil, errors.New("module file with template does not contain post-content context lines")
			}
			return nil, errors.New("module file does not contain post-content context lines")
		}

		dest = dest[:preIdx] + deleted + dest[postIdx:]
	}

	return []byte(dest), nil
}
