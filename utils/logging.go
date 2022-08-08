package utils

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() (logger *zap.Logger, err error) {
	prodEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	devEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	logFile, _ := os.OpenFile("logging.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	writer := zapcore.AddSync(logFile)

	core := zapcore.NewTee(
		zapcore.NewCore(prodEncoder, writer, zap.DebugLevel),
		zapcore.NewCore(devEncoder, zapcore.AddSync(os.Stdout), zap.DebugLevel),
	)

	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	defer logger.Sync()

	return
}
