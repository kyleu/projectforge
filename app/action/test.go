package action

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

func onTest(ctx context.Context, params *Params) *Result {
	methodName := params.Cfg.GetStringOpt("method")
	params.Logger.Infof("running test method [%s]...", methodName)
	switch methodName {
	case "":
		return errorResult(errors.New("must provide test method"), params.Cfg, params.Logger)
	case "bootstrap":
		return bootstrap(ctx, params)
	default:
		return errorResult(errors.Errorf("invalid test method [%s]", methodName), params.Cfg, params.Logger)
	}
}

func bootstrap(ctx context.Context, params *Params) *Result {
	params.Cfg.Add(
		"path", "./testproject",
		"name", "Test Project",
		"summary", "A Test Project!",
		"package", "projectforge.dev/projectforge/testproject",
		"homepage", "https://projectforge.dev",
	)

	err := wipeIfNeeded(params.Cfg, params.Logger)
	if err != nil {
		return errorResult(err, params.Cfg, params.Logger)
	}

	return onCreate(ctx, params)
}

func wipeIfNeeded(cfg util.ValueMap, logger *zap.SugaredLogger) error {
	shouldWipe, _ := cfg.ParseBool("wipe", true, true)
	if shouldWipe {
		fs := filesystem.NewFileSystem(".", logger)
		path := cfg.GetStringOpt("path")
		if path == "" {
			return errors.New("must provide [path] as an argument")
		}
		if fs.Exists(path) {
			logger.Infof("removing existing directory [%s]", path)
			_ = fs.RemoveRecursive(path)
		}
	}
	return nil
}
