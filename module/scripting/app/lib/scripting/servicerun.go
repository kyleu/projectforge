package scripting

import (
	"strings"

	"github.com/dop251/goja"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

func (s *Service) LoadVM(src string) (*goja.Runtime, error) {
	vm := goja.New()
	_, err := vm.RunString(src)
	if err != nil {
		return nil, err
	}
	return vm, nil
}

func (s *Service) RunPath(pth string, fn string, logger util.Logger) (any, error) {
	src, err := s.LoadScript(pth, logger)
	if err != nil {
		return "", err
	}
	return s.RunScript(src, fn)
}

func (s *Service) RunScript(src string, fn string, args ...any) (any, error) {
	vm, err := s.LoadVM(src)
	if err != nil {
		return "", err
	}
	return s.RunFn(vm, fn, args...)
}

func (s *Service) RunFn(vm *goja.Runtime, fn string, args ...any) (any, error) {
	tFn, ok := goja.AssertFunction(vm.Get(fn))
	if !ok {
		return "", errors.Errorf("script must have a function named [%s]", fn)
	}
	jsArgs := lo.Map(args, func(x any, _ int) goja.Value {
		return vm.ToValue(x)
	})
	return tFn(goja.Undefined(), jsArgs...)
}

func (s *Service) LoadExamples(vm *goja.Runtime) (map[string][][]any, error) {
	exJS := vm.Get("examples")
	if exJS == nil {
		return nil, nil // errors.New("script must have a variable named [examples]")
	}
	ex := exJS.Export()
	ret := map[string][][]any{}
	err := util.CycleJSON(ex, &ret)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse examples, must be of form {\"function-name\": [ [args], [args] ]}")
	}
	return ret, nil
}

func (s *Service) RunExamples(vm *goja.Runtime) (map[string]map[string]any, error) {
	ex, err := s.LoadExamples(vm)
	if err != nil {
		return nil, err
	}
	ret := make(map[string]map[string]any, len(ex))
	for f, tests := range ex {
		res, e := s.RunExamplesFunc(vm, f, tests...)
		if e != nil {
			return nil, e
		}
		ret[f] = res
	}
	return ret, nil
}

func (s *Service) RunExamplesFunc(vm *goja.Runtime, f string, tests ...[]any) (map[string]any, error) {
	ret := map[string]any{}
	for _, testArgs := range tests {
		x, err := s.RunFn(vm, f, testArgs...)
		if err != nil {
			return nil, err
		}
		xKey := strings.TrimPrefix(strings.TrimSuffix(util.ToJSONCompact(testArgs), "]"), "[")
		xKey = f + "(" + xKey + ")"
		ret[xKey] = x
	}
	return ret, nil
}
