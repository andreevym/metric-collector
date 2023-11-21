package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func NewLogger(level string) (*zap.Logger, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return zl, err
}

func Logger(level string) (*zap.Logger, error) {
	if Log != nil {
		return Log, nil
	}

	var err error
	Log, err = NewLogger(level)
	if err != nil {
		return nil, err
	}

	return Log, nil
}
