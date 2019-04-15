package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

var LocalLogger, _ = zap.NewProduction()

func InitLocalLog() bool {

	level := zap.NewAtomicLevel()
	level.SetLevel(zapcore.InfoLevel)

	logConfig := zap.Config{
		Level:       level,
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "Time",
			LevelKey:       "Level",
			NameKey:        "Name",
			CallerKey:      "Caller",
			MessageKey:     "Msg",
			StacktraceKey:  "St",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"localLog.txt", "stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	var err error
	LocalLogger, err = logConfig.Build()
	if err != nil {
		log.Fatal("LocalLogger Build Fail!: ", err.Error())
		return false
	}

	return true
}
