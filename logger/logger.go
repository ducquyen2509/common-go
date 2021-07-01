package logger

import (
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"sync"
)

var instance *zap.Logger
var sugarInstance *zap.SugaredLogger
var once sync.Once

func getLogger() *zap.Logger {
	once.Do(func() {
		instance = newLogger()
	})
	return instance
}

func getSugaredLogger() *zap.SugaredLogger {
	if sugarInstance == nil {
		sugarInstance = getLogger().Sugar()
	}
	return sugarInstance
}

func newLogger() (logger *zap.Logger) {

	if len(os.Getenv("debug")) > 0 {
		encoderConfig := zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger = zap.New(
			zapcore.NewCore(
				zapcore.NewConsoleEncoder(encoderConfig),
				zapcore.AddSync(colorable.NewColorableStdout()),
				zapcore.DebugLevel,
			),
			zap.Development(),
			zap.AddCaller(),
			zap.AddStacktrace(zapcore.WarnLevel),
		)
		return logger
	}

	logger, _ = zap.NewProduction()
	return logger
}

func Debug(message string, fields ...zap.Field) {
	getLogger().Debug(message, fields...)
}

func Debugf(template string, args ...interface{}) {
	getSugaredLogger().Debugf(template, args...)
}

func Error(message string, fields ...zap.Field) {
	getLogger().Error(message, fields...)
}

func Errorf(template string, args ...interface{}) {
	getSugaredLogger().Errorf(template, args...)
}

func Fatal(message string, fields ...zap.Field) {
	getLogger().Fatal(message, fields...)
}

func Fatalf(template string, args ...interface{}) {
	getSugaredLogger().Fatalf(template, args...)
}

func Info(message string, fields ...zap.Field) {
	getLogger().Info(message, fields...)
}

func Infof(template string, args ...interface{}) {
	getSugaredLogger().Infof(template, args...)
}

func Warn(message string, fields ...zap.Field) {
	getLogger().Warn(message, fields...)
}

func Warnf(template string, args ...interface{}) {
	getSugaredLogger().Warnf(template, args...)
}

func AddHook(hook func(zapcore.Entry) error) {
	instance = getLogger().WithOptions(zap.Hooks(hook))
	sugarInstance = instance.Sugar()
}

func WithField(key, val string) *zap.Logger {
	return getLogger().With(zap.String(key, val))
}

func With(fields ...zap.Field) *zap.Logger {
	return getLogger().With(fields...)
}

func WithRequest(r *http.Request) *zap.Logger {
	return getLogger().With(
		zap.Any("method", r.Method),
		zap.Any("host", r.Host),
		zap.Any("path", r.URL.Path),
	)
}

func SugaredWithRequest(r *http.Request) *zap.SugaredLogger {
	return getSugaredLogger().With(
		zap.Any("method", r.Method),
		zap.Any("host", r.Host),
		zap.Any("path", r.URL.Path),
	)
}
