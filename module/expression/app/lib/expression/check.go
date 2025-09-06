package expression

import (
	"sync"

	"github.com/google/cel-go/cel"

	"{{{ .Package }}}/app/util"
)

var (
	expressionMU    = &sync.Mutex{}
	expressionCache = make(map[string]*Expression)
)

type Engine struct {
	env *cel.Env
}

func NewEngine(opts ...cel.EnvOption) (*Engine, error) {
	var err error
	env, err := cel.NewEnv(opts...)
	return &Engine{env: env}, err
}

func (e *Engine) Check(as string, params map[string]any, logger util.Logger) (bool, error) {
	expressionMU.Lock()
	ex, ok := expressionCache[as]
	expressionMU.Unlock()
	if !ok {
		ex = NewExpression("temp", "a temporary expression", as)

		err := ex.Compile(e)
		if err != nil {
			logger.Errorf("error compiling expression [%v]: %+v", as, err)
			return false, err
		}
		expressionMU.Lock()
		expressionCache[as] = ex
		expressionMU.Unlock()
	}

	rsp, _, err := ex.Run(params)
	if err != nil {
		logger.Debugf("error running expression [%v]: %v", as, err.Error())
		return false, err
	}
	ret := CheckResult(rsp, logger)
	return ret, nil
}

func (e *Engine) Compile(as string, logger util.Logger) (*Expression, error) {
	expressionMU.Lock()
	ex, ok := expressionCache[as]
	expressionMU.Unlock()
	if !ok {
		ex = NewExpression("temp", "a temporary expression", as)

		err := ex.Compile(e)
		if err != nil {
			logger.Errorf("error compiling expression [%v]: %+v", as, err)
			return nil, err
		}
		expressionMU.Lock()
		expressionCache[as] = ex
		expressionMU.Unlock()
	}
	return ex, nil
}
