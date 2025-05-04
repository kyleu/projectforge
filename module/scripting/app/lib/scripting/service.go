package scripting

import (
	"context"
	"strings"

	"github.com/samber/lo"

	"{{{ .Package }}}/app/lib/filesystem"
	"{{{ .Package }}}/app/lib/filter"
	"{{{ .Package }}}/app/lib/search/result"
	"{{{ .Package }}}/app/util"
)

type Service struct {
	FS   filesystem.FileLoader `json:"-"`
	Path string                `json:"path,omitempty"`
}

func NewService(fs filesystem.FileLoader, pth string) *Service {
	return &Service{FS: fs, Path: pth}
}

func (s *Service) ListScripts(logger util.Logger) []string {
	files, _ := s.FS.ListFilesRecursive(s.Path, nil, logger)
	return lo.FilterMap(files, func(x string, _ int) (string, bool) {
		return x, strings.HasSuffix(x, ".js")
	})
}

func (s *Service) ListScriptSizes(logger util.Logger) ([]string, map[string]int) {
	files := s.ListScripts(logger)
	return files, lo.SliceToMap(files, func(scr string) (string, int) {
		return scr, s.Size(scr)
	})
}

func (s *Service) LoadScript(pth string, logger util.Logger) (string, string, error) {
	logger.Infof("loading script [%s]", pth)
	filePath := util.StringFilePath(s.Path, pth)
	b, err := s.FS.ReadFile(filePath)
	if err != nil {
		pth += ".js"
		b, err = s.FS.ReadFile(filePath + ".js")
		if err != nil {
			return "", "", err
		}
	}
	sc := string(b)
	return pth, sc, nil
}

func (s *Service) SaveScript(pth string, content string, logger util.Logger) error {
	logger.Infof("saving script [%s]", pth)
	filePath := util.StringFilePath(s.Path, pth)
	return s.FS.WriteFile(filePath, []byte(content), filesystem.DefaultMode, true)
}

func (s *Service) DeleteScript(pth string, logger util.Logger) error {
	filePath := util.StringFilePath(s.Path, pth)
	return s.FS.Remove(filePath, logger)
}

func (s *Service) Size(scr string) int {
	filePath := util.StringFilePath(s.Path, scr)
	st, err := s.FS.Stat(filePath)
	if err != nil {
		return 0
	}
	return int(st.Size)
}

func (s *Service) SearchScripts(_ context.Context, _ filter.ParamSet, q string, logger util.Logger) (result.Results, error) {
	f := func(fn string, _ int) (*result.Result, bool) {
		fn, scr, _ := s.LoadScript(fn, logger)
		res := result.NewResult("script", fn, "/admin/scripting/"+fn, fn, "cog", scr, scr, q)
		if len(res.Matches) > 0 {
			return res, true
		}
		return nil, false
	}
	return lo.FilterMap(s.ListScripts(logger), f), nil
}
