package logger

import (
	"context"
	"fmt"
	"github/hsj/golite/config"
	"log/slog"
	"path/filepath"
	"strings"
)

const (
	LevelDebug   = slog.LevelDebug
	LevelTrace   = slog.Level(-2)
	LevelInfo    = slog.LevelInfo
	LevelWarning = slog.LevelWarn
	LevelError   = slog.LevelError
	LevelFatal   = slog.Level(12)
)

type Logger interface {
	Debug(ctx context.Context, format string, args ...any)
	Trace(ctx context.Context, format string, args ...any)
	Info(ctx context.Context, format string, args ...any)
	Warning(ctx context.Context, format string, args ...any)
	Fatal(ctx context.Context, format string, args ...any)
}

var LevelNames = map[slog.Leveler]string{
	LevelTrace: "TRACE",
	LevelFatal: "FATAL",
}

var LevelMap = map[string]slog.Level{
	"TRACE": LevelTrace,
	"DEBUG": LevelDebug,
	"INFO":  LevelInfo,
	"WARN":  LevelWarning,
	"ERROR": LevelError,
	"FATAL": LevelFatal,
}

type LogConfig struct {
	Dir      string `toml:"dir"`
	FileName string `toml:"filename"`
	MinLevel string `toml:"level"`
	Format   string `toml:"format"`
}

func parse(conf string) (*LogConfig, error) {
	var logConfig LogConfig
	if err := config.Parse(conf, &logConfig); err != nil {
		return nil, err
	}
	if logConfig.Dir == "" {
		logConfig.Dir = "logs"
	}
	absDir, err := filepath.Abs(logConfig.Dir)
	if err != nil {
		return nil, err
	}
	logConfig.Dir = absDir
	return &logConfig, nil
}

func NewLogger(ctx context.Context, conf string) (Logger, error) {
	logConf, err := parse(conf)
	if err != nil {
		return nil, err
	}

	logLevel, ok := LevelMap[strings.ToUpper(logConf.MinLevel)]
	if !ok {
		return nil, fmt.Errorf("invalid log level: %s", logConf.MinLevel)
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
		ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
			if attr.Key == slog.LevelKey {
				level := attr.Value.Any().(slog.Level)
				levelLabel, exists := LevelNames[level]
				if !exists {
					levelLabel = level.String()
				}
				attr.Value = slog.StringValue(levelLabel)
			}
			return attr
		},
	}

	if logConf.Dir != "" && logConf.FileName != "" {
		return NewTextLogger(ctx, logConf, opts)
	}

	return NewConsoleLogger(ctx, logConf, opts)
}
