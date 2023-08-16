package scripting

import (
	"path/filepath"
	"strings"

	"github.com/dop251/goja"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"{{{ .Package }}}/app/lib/filesystem"
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

func (s *Service) LoadScript(pth string, logger util.Logger) (string, error) {
	logger.Infof("loading script [%s]", pth)
	filePath := filepath.Join(s.Path, pth)
	b, err := s.FS.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	sc := string(b)
	return sc, nil
}

func (s *Service) RunPath(pth string, fn string, logger util.Logger) (any, error) {
	src, err := s.LoadScript(pth, logger)
	if err != nil {
		return "", err
	}
	return s.RunScript(src, fn)
}

func (s *Service) RunScript(src string, fn string, args ...any) (any, error) {
	vm := goja.New()
	_, err := vm.RunString(src)
	if err != nil {
		return "", err
	}

	tFn, ok := goja.AssertFunction(vm.Get(fn))
	if !ok {
		return "", errors.Errorf("script must have a function named [%s]", fn)
	}

	jsArgs := lo.Map(args, func(x any, _ int) goja.Value {
		return vm.ToValue(x)
	})

	res, err := tFn(goja.Undefined(), jsArgs...)
	return res.Export(), nil
}
