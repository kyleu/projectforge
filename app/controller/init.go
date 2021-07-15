package controller

import (
	"github.com/kyleu/projectforge/app/module"
	"github.com/kyleu/projectforge/app/project"
	"go.uber.org/zap"
)

var (
	controllerLogger     *zap.SugaredLogger
	controllerModuleSvc  *module.Service
	controllerProjectSvc *project.Service
)

// This method is where you'll initialize app-specific dependencies.
func initApp() {
	controllerLogger = _currentAppState.Logger
	controllerModuleSvc = module.NewService(_currentAppState.Logger)
	controllerProjectSvc = project.NewService(_currentAppState.Logger)
}

// This method is where you'll initialize dependencies for the marketing site.
func initSite() {
}
