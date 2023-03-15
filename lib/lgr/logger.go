package lgr

import (
	"log"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DebugLevel = "DEBUG"
	InfoLevel  = "INFO"
	WarnLevel  = "WARN"
	ErrorLevel = "ERROR"
)

type Logger struct{ *zap.SugaredLogger }

func getLevel() zapcore.Level {
	lvl := os.Getenv("LOG_LEVEL")
	switch strings.ToUpper(lvl) {
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	default:
		return zapcore.DebugLevel
	}
}

func getEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:          "message",
		LevelKey:            "level",
		TimeKey:             "time",
		NameKey:             "name",
		CallerKey:           "caller",
		FunctionKey:         "",
		StacktraceKey:       "stacktrace",
		SkipLineEnding:      false,
		LineEnding:          "\n",
		EncodeLevel:         zapcore.CapitalLevelEncoder,
		EncodeTime:          zapcore.ISO8601TimeEncoder,
		EncodeDuration:      zapcore.NanosDurationEncoder,
		EncodeCaller:        zapcore.ShortCallerEncoder,
		NewReflectedEncoder: nil,
		ConsoleSeparator:    "\t",
	}
}

func getConsoleCore() zapcore.Core {
	cfg := getEncoderConfig()
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg),
		zapcore.Lock(os.Stderr),
		getLevel(),
	)
}

func getJSONCore() zapcore.Core {
	cfg := getEncoderConfig()

	file, err := os.Create(os.Getenv("LOG_FILE"))
	if err != nil {
		log.Fatalf("creating logger file: %v\n", err)
	}

	return zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg),
		zapcore.Lock(file),
		getLevel(),
	)
}

func New() Logger {
	cores := []zapcore.Core{
		getConsoleCore(),
		getJSONCore(),
	}

	core := zapcore.NewTee(cores...)

	return Logger{zap.New(core).Sugar()}
}
