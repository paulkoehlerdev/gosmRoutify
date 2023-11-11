package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"strings"
)

type LogLevel string

const (
	LevelError   LogLevel = "error"
	LevelWarning          = "warning"
	LevelInfo             = "info"
	LevelDebug            = "debug"
)

type Logger interface {
	Info() Logger
	Warn() Logger
	Error() Logger
	Debug() Logger

	Msg(string)
	Msgf(string, ...any)

	WithAttrs(...any) Logger
}

type impl struct {
	logger *slog.Logger
	level  slog.Level
}

func New(level LogLevel, writer io.Writer) Logger {
	return &impl{
		logger: slog.New(slog.NewTextHandler(writer, &slog.HandlerOptions{Level: slog.LevelDebug})),
		level:  parseLogLevel(level),
	}
}

func parseLogLevel(level LogLevel) slog.Level {
	switch level {
	case LevelError:
		return slog.LevelError
	case LevelWarning:
		return slog.LevelWarn
	case LevelInfo:
		return slog.LevelInfo
	case LevelDebug:
		return slog.LevelDebug
	}
	return slog.LevelInfo
}

func getStackTrace() string {
	buf := make([]byte, 1<<16)
	n := runtime.Stack(buf, false)
	buf = buf[:n]

	stackStr := strings.ReplaceAll(string(buf), "\n\t", ", ")
	stackStr = strings.ReplaceAll(stackStr, "\n", ", ")
	return stackStr
}

func (i *impl) Info() Logger {
	return i.withLevel(slog.LevelInfo)
}

func (i *impl) Warn() Logger {
	return i.withLevel(slog.LevelWarn)
}

func (i *impl) Error() Logger {
	return i.withAttrs("stack", getStackTrace()).withLevel(slog.LevelError)
}

func (i *impl) Debug() Logger {
	return i.withLevel(slog.LevelDebug)
}

func (i *impl) Msg(msg string) {
	i.logger.Log(context.Background(), i.level, msg)
}

func (i *impl) Msgf(msg string, attr ...any) {
	i.Msg(fmt.Sprintf(msg, attr...))
}

func (i *impl) WithAttrs(args ...any) Logger {
	return i.withAttrs(args...)
}

func (i *impl) withAttrs(args ...any) *impl {
	return &impl{
		logger: i.logger.With(args...),
		level:  i.level,
	}
}

func (i *impl) withLevel(level slog.Level) Logger {
	return &impl{
		logger: i.logger,
		level:  level,
	}
}
