package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var sugar *zap.SugaredLogger

func Sync() {
	logger.Sync()
}

func InitZapLogger() *zap.Logger {
	logger = zap.NewExample()
	sugar = zap.NewExample().Sugar()
	return logger
}

func ZapCoreFileds(zapFileds ...zapcore.Field) []zapcore.Field {
	return zapFileds
}

func SugerArgs(args ...interface{}) []interface{} {
	return args
}

func Info(msg string, zapFields ...zapcore.Field) {
	logger.Info(msg, zapFields...)
}

// - Example
// logger.Infof("test implement log zap args1: %s, args2: %d",
// 	logger.SugerArgs("test2", 9999),
// 	logger.ZapCoreFileds(
// 		zap.String("uuid", "1234"),
// 		zap.String("trace_id", "uasd1efg2"),),
// )

func Infof(template string, args []interface{}, zapFileds []zapcore.Field) {
	if len(zapFileds) == 0 || zapFileds == nil {
		sugar.Infof(template, args...)
	} else {
		sugar.WithOptions(zap.Fields(zapFileds...)).Infof(template, args...)
	}
}

func Error(msg string, zapFields ...zapcore.Field) {
	logger.Error(msg, zapFields...)
}

func Errorf(template string, args []interface{}, zapFileds []zapcore.Field) {
	if len(zapFileds) == 0 || zapFileds == nil {
		sugar.Errorf(template, args...)
	} else {
		sugar.WithOptions(zap.Fields(zapFileds...)).Errorf(template, args...)
	}
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Warnf(template string, args []interface{}, zapFileds []zapcore.Field) {
	if len(zapFileds) == 0 || zapFileds == nil {
		sugar.Warnf(template, args...)
	} else {
		sugar.WithOptions(zap.Fields(zapFileds...)).Warnf(template, args...)
	}
}

func Fatal(msg string, fields ...zapcore.Field) {
	logger.Fatal(msg, fields...)
}

func Panic(msg string, fields ...zapcore.Field) {
	logger.Panic(msg, fields...)
}
