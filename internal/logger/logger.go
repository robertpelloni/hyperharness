package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func init() {
	InitLogger(false)
}

func InitLogger(debug bool) {
	level := zapcore.InfoLevel
	if debug {
		level = zapcore.DebugLevel
	}

	// Create a core that writes to both stdout and a file
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	// Open file for writing logs (append mode)
	file, _ := os.OpenFile("hyperharness.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(file), level),
	)

	Log = zap.New(core, zap.AddCaller())
}
