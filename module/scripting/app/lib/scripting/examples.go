package scripting

import (
	"strings"

	"github.com/dop251/goja"
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/util"
)

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
