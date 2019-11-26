package loggerService

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

const (
	requestIdKey = "requestId"
)

var mainLoggerServiceInstance *LoggerService

type LoggerService struct {
	logger *zap.SugaredLogger
}

func GetMainLogger() *LoggerService {
	return mainLoggerServiceInstance
}

func createLogger() (*LoggerService, error) {
	config := zap.NewProductionConfig()
	config.Encoding = "console"
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig.CallerKey = ""

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &LoggerService{logger: logger.Sugar()}, nil
}

func (log *LoggerService) Info(ctx context.Context, msg string, keysAndValues ...interface{}) {
	log.logger.Infow(log.decorateLogMessage(ctx, msg), keysAndValues...)
}

func (log *LoggerService) Warn(ctx context.Context, msg string, keysAndValues ...interface{}) {
	log.logger.Warnw(log.decorateLogMessage(ctx, msg), keysAndValues...)
}

func (log *LoggerService) Error(ctx context.Context, msg string, keysAndValues ...interface{}) {
	log.logger.Errorw(log.decorateLogMessage(ctx, msg), keysAndValues...)
}

func (log *LoggerService) Fatal(ctx context.Context, msg string, keysAndValues ...interface{}) {
	log.logger.Fatalw(log.decorateLogMessage(ctx, msg), keysAndValues...)
}

func (log *LoggerService) decorateLogMessage(ctx context.Context, msg string) string {
	if ctx != nil {
		if requestId, ok := ctx.Value(requestIdKey).(string); ok {
			return fmt.Sprintf("[%s] %s", requestId, msg)
		}
	}

	return msg
}

func init() {
	var err error
	mainLoggerServiceInstance, err = createLogger()

	if err != nil {
		log.Fatalf("can't init logger")
	}
}
