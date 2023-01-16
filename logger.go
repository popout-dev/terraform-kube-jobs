package main

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/rs/zerolog/pkgerrors"
)

type Logger struct{}

func NewLogger() Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	return Logger{}
}

func (logger *Logger) Info(message string) {
	log.Info().Msg(message)
}

func (logger *Logger) Error(err error) {
	log.Error().Stack().Err(err).Msg("")
}

func (logger *Logger) Fatal(err error) {
	log.Fatal().Stack().Err(err).Msg("")
}

func (logger *Logger) StackTrace(err error) error {
	return errors.Wrap(err, "")
}

type Key int

const (
	LoggerKey Key = iota
)

func AddLoggerToContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, LoggerKey, logger)
}

func LoggerFromContext(ctx context.Context) (Logger, error) {
	loggerValue := ctx.Value(LoggerKey)

	logger, ok := loggerValue.(Logger)
	if !ok {
		return Logger{}, errors.New("unexpected logger type")
	}

	return logger, nil
}
