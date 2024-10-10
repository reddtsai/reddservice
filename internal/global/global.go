package global

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	loggerOnce  sync.Once
	Logger      *zap.Logger
	SugarLogger *zap.SugaredLogger
	configOnce  sync.Once
	Config      *Configuration
)

// here value is set by ldflags
var (
	CONFIG_PATH = "conf.d"
)

func init() {
	initConfiguration()
	initZapLogger()
}

func initZapLogger() {
	loggerOnce.Do(func() {
		config := zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(zapcore.Level(Config.LogOpts.Level))
		logger, err := config.Build()
		if err != nil {
			panic(err)
		}
		Logger = logger
		zap.ReplaceGlobals(Logger)
		SugarLogger = Logger.Sugar()
	})
}
