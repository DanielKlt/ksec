package environment

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var env *environment

type environment struct {
	vp *viper.Viper
	lg *zap.Logger
}

func GetEnv() *environment {
	if env == nil {
		env = &environment{}
	}
	return env
}

func (e *environment) GetViper() *viper.Viper {
	if e.vp == nil {
		e.vp = viper.New()
	}
	return e.vp
}

func (e *environment) GetLogger() *zap.Logger {
	if e.lg == nil {
		newLogger()
	}
	return e.lg
}

func getLoggerLevel() string {
	if v, ok := os.LookupEnv("LOG_LEVEL"); ok {
		return v
	}
	return "info"
}

func newLogger() {
	var cfg zap.Config
	cfg.Level.UnmarshalText([]byte(getLoggerLevel()))
	cfg.Encoding = "json"
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}
	cfg.EncoderConfig.MessageKey = "message"
	cfg.EncoderConfig.LevelKey = "level"
	cfg.EncoderConfig.EncodeLevel = zap.NewProductionEncoderConfig().EncodeLevel
	lg, err := cfg.Build()
	if err != nil {
		panic("failed to create logger")
	}
	defer lg.Sync()
	env.lg = lg
}
