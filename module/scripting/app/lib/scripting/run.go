package scripting

import (
	"github.com/dop251/goja"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

func LoadVM(key string, src string, logger util.Logger) (any, *goja.Runtime, error) {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	err := wireFunctions(key, vm, logger)
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

func wireFunctions(key string, vm *goja.Runtime, logger util.Logger) error {
	prepend := func(args []any) []any {
		return append([]any{"[" + key + "]: "}, args...)
	}
	console := map[string]any{
		"debug": func(args ...any) { logger.Debug(prepend(args)...) },
		"log":   func(args ...any) { logger.Info(prepend(args)...) },
		"info":  func(args ...any) { logger.Info(prepend(args)...) },
		"warn":  func(args ...any) { logger.Warn(prepend(args)...) },
		"error": func(args ...any) { logger.Error(prepend(args)...) },
	}
	if err := vm.Set("console", console); err != nil {
		return err
	}
	if err := vm.Set("randomString", util.RandomString); err != nil {
		return err
	}
	if err := vm.Set("microsToMillis", util.MicrosToMillis); err != nil {
		return err
	}
	if err := vm.Set("getEnv", util.GetEnv); err != nil {
		return err
	}
	if err := vm.Set("replaceEnvVars", util.ReplaceEnvVars); err != nil {
		return err
	}
	return nil
}
