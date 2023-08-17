package scripting

import (
	"strings"

	"github.com/dop251/goja"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

func LoadVM(key string, src string, logger util.Logger) (any, *goja.Runtime, error) {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	prepend := func(args []any) []any {
		return append([]any{"[" + key + "]: "}, args...)
	}
	err := vm.Set("console", map[string]any{
		"debug": func(args ...any) {
			logger.Debug(prepend(args)...)
		},
		"log": func(args ...any) {
			logger.Info(prepend(args)...)
		},
		"info": func(args ...any) {
			logger.Info(prepend(args)...)
		},
		"warn": func(args ...any) {
			logger.Warn(prepend(args)...)
		},
		"error": func(args ...any) {
			logger.Error(prepend(args)...)
		},
	})
	if err != nil {
		return nil, nil, err
	}
	res, err := vm.RunString(src)
	if err != nil {
		return nil, nil, err
	}
	return res.Export(), vm, nil
}

func RunFn(vm *goja.Runtime, fn string, args ...any) (any, error) {
	tFn, ok := goja.AssertFunction(vm.Get(fn))
	if !ok {
		return "", errors.Errorf("script must have a function named [%s]", fn)
	}
	jsArgs := lo.Map(args, func(x any, _ int) goja.Value {
		return vm.ToValue(x)
	})
	return tFn(goja.Undefined(), jsArgs...)
}

func LoadExamples(vm *goja.Runtime) (map[string][][]any, error) {
	exJS := vm.Get("examples")
	if exJS == nil {
		return nil, nil
	}
	ex := exJS.Export()
	ret := map[string][][]any{}
	err := util.CycleJSON(ex, &ret)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse examples, must be of form {\"function-name\": [ [args], [args] ]}")
	}
	return ret, nil
}

func RunExamples(vm *goja.Runtime) (map[string]map[string]any, error) {
	ex, err := LoadExamples(vm)
	if err != nil {
		return nil, err
	}
	ret := make(map[string]map[string]any, len(ex))
	for f, tests := range ex {
		res, e := RunExamplesFunc(vm, f, tests...)
		if e != nil {
			return nil, e
		}
		ret[f] = res
	}
	return ret, nil
}

func RunExamplesFunc(vm *goja.Runtime, f string, tests ...[]any) (map[string]any, error) {
	ret := map[string]any{}
	for _, testArgs := range tests {
		xKey := strings.TrimPrefix(strings.TrimSuffix(util.ToJSONCompact(testArgs), "]"), "[")
		x, err := RunFn(vm, f, testArgs...)
		if err != nil {
			return nil, err
		}
		xKey = f + "(" + xKey + ")"
		ret[xKey] = x
	}
	return ret, nil
}
