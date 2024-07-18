package util

import "go.uber.org/zap"

type Logger = *zap.SugaredLogger

var RootLogger Logger
