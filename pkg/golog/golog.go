package golog

import "go.uber.org/zap"

func New() (*zap.Logger, error) {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	loggerConfig.Encoding = "json"
	loggerConfig.DisableCaller = true
	loggerConfig.OutputPaths = []string{"stdout"}
	if logger, err := loggerConfig.Build(); err != nil {
		return nil, err
	} else {
		return logger, nil
	}
}
