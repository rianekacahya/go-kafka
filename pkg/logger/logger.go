package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

var (
	instance *zap.Logger
	sugarInstance *zap.SugaredLogger
	once sync.Once
)

func GetLogger() *zap.Logger {
	once.Do(func() {
		instance = newLogger()
	})
	return instance
}

func GetSugaredLogger() *zap.SugaredLogger {
	if sugarInstance == nil {
		sugarInstance = GetLogger().Sugar()
	}
	return sugarInstance
}

func newLogger() *zap.Logger {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	loggerConfig.Encoding = "json"
	loggerConfig.DisableCaller = true
	loggerConfig.OutputPaths = []string{"stdout"}
	if logger, err := loggerConfig.Build(); err != nil {
		panic(err)
	} else {
		return logger
	}
}

// Debug logs the message at debug level with additional fields, if any
func Debug(message string, fields ...zap.Field) {
	GetLogger().Debug(message, fields...)
}

// Debugf allows Sprintf style formatting and logs at debug level
func Debugf(template string, args ...interface{}) {
	GetSugaredLogger().Debugf(template, args...)
}

// Error logs the message at error level and prints stacktrace with additional fields, if any
func Error(message string, fields ...zap.Field) {
	GetLogger().Error(message, fields...)
}

// Errorf allows Sprintf style formatting, logs at error level and prints stacktrace
func Errorf(template string, args ...interface{}) {
	GetSugaredLogger().Errorf(template, args...)
}

// Fatal logs the message at fatal level with additional fields, if any and exits
func Fatal(message string, fields ...zap.Field) {
	GetLogger().Fatal(message, fields...)
}

// Fatalf allows Sprintf style formatting, logs at fatal level and exits
func Fatalf(template string, args ...interface{}) {
	GetSugaredLogger().Fatalf(template, args...)
}

// Info logs the message at info level with additional fields, if any
func Info(message string, fields ...zap.Field) {
	GetLogger().Info(message, fields...)
}

// Infof allows Sprintf style formatting and logs at info level
func Infof(template string, args ...interface{}) {
	GetSugaredLogger().Infof(template, args...)
}

// Warn logs the message at warn level with additional fields, if any
func Warn(message string, fields ...zap.Field) {
	GetLogger().Warn(message, fields...)
}

// Warnf allows Sprintf style formatting and logs at warn level
func Warnf(template string, args ...interface{}) {
	GetSugaredLogger().Warnf(template, args...)
}

// AddHook adds func(zapcore.Entry) error) to the logger lifecycle
func AddHook(hook func(zapcore.Entry) error) {
	instance = GetLogger().WithOptions(zap.Hooks(hook))
	sugarInstance = instance.Sugar()
}
