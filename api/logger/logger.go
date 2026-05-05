package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	logger zerolog.Logger
}

var logInstance *Logger

func getLoggerWithParams(format, level string) *Logger {
	var writer io.Writer

	if format == "json" {
		writer = os.Stdout
	} else {
		writer = zerolog.ConsoleWriter{Out: os.Stdout}
	}

	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}

	logger := zerolog.New(writer).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(lvl)
	zerolog.TimeFieldFormat = time.RFC3339
	return &Logger{
		logger: logger,
	}
}

func init() {
	logInstance = getLoggerWithParams("text", "info")
}

type Config struct {
	Format string
	Level  string
}

func Configure(cfg *Config) {

	if cfg == nil {
		cfg = &Config{
			Format: "text",
			Level:  "info",
		}
	}

	logInstance = getLoggerWithParams(cfg.Format, cfg.Level)
}

func Error(ctx context.Context, message ...any) {
	logInstance.logger.Error().Fields(GetMetadata(ctx)).Msg(getMessage(message...))
}

func Log(ctx context.Context) {
	logInstance.logger.Log().Fields(GetMetadata(ctx)).Send()
}

func Info(ctx context.Context, message ...any) {
	logInstance.logger.Info().Fields(GetMetadata(ctx)).Msg(getMessage(message...))
}

func Debug(ctx context.Context, message ...any) {
	logInstance.logger.Debug().Fields(GetMetadata(ctx)).Msg(getMessage(message...))
}

func Warn(ctx context.Context, message ...any) {
	logInstance.logger.Warn().Fields(GetMetadata(ctx)).Msg(getMessage(message...))
}

func Fatal(ctx context.Context, message ...any) {
	logInstance.logger.Fatal().Fields(GetMetadata(ctx)).Msg(getMessage(message...))
}

func Panic(ctx context.Context, message ...any) {
	logInstance.logger.Panic().Fields(GetMetadata(ctx)).Msg(getMessage(message...))
}

func getMessage(message ...any) string {
	// convert the any array to a single string message
	var msg string
	for i, m := range message {
		msg += fmt.Sprintf("%v", m)
		if i < len(message)-1 {
			msg += " "
		}
	}
	return msg
}

type loggerCtx struct{}

func SetMetadata(ctx context.Context, metadata map[string]interface{}) context.Context {
	if len(metadata) == 0 {
		return ctx
	}

	mergedMetadata := make(map[string]interface{})
	if existingMetadata := GetMetadata(ctx); existingMetadata != nil {
		for k, v := range existingMetadata {
			mergedMetadata[k] = v
		}
	}
	for k, v := range metadata {
		mergedMetadata[k] = v
	}

	return context.WithValue(ctx, loggerCtx{}, mergedMetadata)
}

func GetMetadata(ctx context.Context) map[string]interface{} {
	if ctx == nil {
		return nil
	}
	if metadata, ok := ctx.Value(loggerCtx{}).(map[string]interface{}); ok {
		return metadata
	}
	return nil
}
