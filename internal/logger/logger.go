package logger

import (
	"log/slog"
	"os"
)

// ----------------------------------------------------------------
// LOGGER
// ----------------------------------------------------------------

type LevelT int

const (
	DEBUG LevelT = LevelT(slog.LevelDebug)
	INFO  LevelT = LevelT(slog.LevelInfo)
	WARN  LevelT = LevelT(slog.LevelWarn)
	ERROR LevelT = LevelT(slog.LevelError)
)

type LoggerT struct {
	Logger       *slog.Logger
	CommonFields map[string]any
}

func NewLogger(level LevelT, commonFields map[string]any) (logger LoggerT) {
	logger.CommonFields = commonFields

	opts := &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.Level(level),
	}
	jsonHandler := slog.NewJSONHandler(os.Stdout, opts)
	logger.Logger = slog.New(jsonHandler)

	return logger
}

func GetLevel(levelStr string) (l LevelT) {
	levelMap := map[string]LevelT{
		"debug": DEBUG,
		"info":  INFO,
		"warn":  WARN,
		"error": ERROR,
	}

	l, ok := levelMap[levelStr]
	if !ok {
		l = DEBUG
	}

	return l
}

func (l *LoggerT) Debug(msg string, extraFields map[string]any) {
	extraFieldsArr := []any{}

	for k, v := range l.CommonFields {
		extraFieldsArr = append(extraFieldsArr, k, v)
	}

	for k, v := range extraFields {
		extraFieldsArr = append(extraFieldsArr, k, v)
	}

	l.Logger.Debug(msg, extraFieldsArr...)
}

func (l *LoggerT) Info(msg string, extraFields map[string]any) {
	extraFieldsArr := []any{}

	for k, v := range l.CommonFields {
		extraFieldsArr = append(extraFieldsArr, k, v)
	}

	for k, v := range extraFields {
		extraFieldsArr = append(extraFieldsArr, k, v)
	}

	l.Logger.Info(msg, extraFieldsArr...)
}

func (l *LoggerT) Warn(msg string, extraFields map[string]any) {
	extraFieldsArr := []any{}

	for k, v := range l.CommonFields {
		extraFieldsArr = append(extraFieldsArr, k, v)
	}

	for k, v := range extraFields {
		extraFieldsArr = append(extraFieldsArr, k, v)
	}

	l.Logger.Warn(msg, extraFieldsArr...)
}

func (l *LoggerT) Error(msg string, extraFields map[string]any) {
	extraFieldsArr := []any{}

	for k, v := range l.CommonFields {
		extraFieldsArr = append(extraFieldsArr, k, v)
	}

	for k, v := range extraFields {
		extraFieldsArr = append(extraFieldsArr, k, v)
	}

	l.Logger.Error(msg, extraFieldsArr...)
}

func (l *LoggerT) Fatal(msg string, extraFields map[string]any) {
	extraFieldsArr := []any{}

	for k, v := range l.CommonFields {
		extraFieldsArr = append(extraFieldsArr, k, v)
	}

	for k, v := range extraFields {
		extraFieldsArr = append(extraFieldsArr, k, v)
	}

	l.Logger.Error(msg, extraFieldsArr...)
	os.Exit(1)
}
