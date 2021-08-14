package action

import (
	"context"

	"github.com/kyleu/projectforge/app/filesystem"
	"github.com/kyleu/projectforge/app/module"
	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func onTest(ctx context.Context, cfg util.ValueMap, mSvc *module.Service, pSvc *project.Service, logger *zap.SugaredLogger) *Result {
	methodName, _ := cfg.GetString("method", true)
	logger.Infof("running test method [%s]...", methodName)
	switch methodName {
	case "":
		return errorResult(errors.New("must provide test method"), cfg, logger)
	case "bootstrap":
		return bootstrap(ctx, cfg, mSvc, pSvc, logger)
	default:
		return errorResult(errors.Errorf("invalid test method [%s]", methodName), cfg, logger)
	}
}

func bootstrap(ctx context.Context, cfg util.ValueMap, mSvc *module.Service, pSvc *project.Service, logger *zap.SugaredLogger) *Result {
	cfg.Add(
		"path", "./testproject",
		"name", "Test Project",
		"summary", "A Test Project!",
		"package", "github.com/kyleu/projectforge/testproject",
		"homepage", "https://projectforge.dev",
	)

	err := wipeIfNeeded(cfg, logger)
	if err != nil {
		return errorResult(err, cfg, logger)
	}

	return onCreate(ctx, "testproject", cfg, mSvc, pSvc, logger)
}

func wipeIfNeeded(cfg util.ValueMap, logger *zap.SugaredLogger) error {
	shouldWipe, _ := cfg.GetBool("wipe")
	if shouldWipe {
		fs := filesystem.NewFileSystem(".", logger)
		path, _ := cfg.GetString("path", true)
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
