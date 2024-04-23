package logger

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var Slogger *zap.SugaredLogger

func InitLogger() {

	// Настройка конфигурации логгера
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	// Настройка логгера с конфигом
	var err error
	Logger, err = config.Build()
	if err != nil {
		fmt.Printf("Ошибка настройки логгера: %v\n", err)
	}
	Slogger = Logger.Sugar()
}
