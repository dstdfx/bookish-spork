package log

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	stdout = "stdout"
	stderr = "stderr"
)

// InitLoggerOpts contains options to the InitLogger function.
type InitLoggerOpts struct {
	File      string
	UseStdout bool
	Debug     bool
}

// InitLogger initializes the Logger from the provided options.
func InitLogger(opts InitLoggerOpts) (*zap.Logger, error) {
	// Configure loglevel.
	loglevel := zap.NewAtomicLevelAt(zap.InfoLevel)
	if opts.Debug {
		loglevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	// Configure output paths.
	outputPaths, errPaths, err := outputConfig(
		opts.File,
		opts.UseStdout,
	)
	if err != nil {
		return nil, err
	}

	// Create a zap logger instance.
	cfg := zapConfig(loglevel, outputPaths, errPaths)

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	// Remove current caller with logger.go.
	// All callers will be added from calling functions in external applications.
	logger = logger.WithOptions(zap.AddCallerSkip(1))

	return logger, nil
}

func outputConfig(file string, useStdout bool) ([]string, []string, error) {
	var outputPaths []string
	errPaths := []string{stderr}

	if useStdout {
		outputPaths = append(outputPaths, stdout)
	}

	if file != "" {
		if _, err := os.Stat(file); err != nil {
			logDir := filepath.Dir(file)
			if _, err := os.Stat(logDir); err != nil {
				return nil, nil, fmt.Errorf("directory %s doesn't exist: %s", logDir, err)
			}
		}

		outputPaths = append(outputPaths, file)
		errPaths = append(errPaths, file)
	}

	return outputPaths, errPaths, nil
}

func zapConfig(level zap.AtomicLevel, outputPaths, errPaths []string) zap.Config {
	return zap.Config{
		Level:    level,
		Encoding: "json",
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			CallerKey:      "caller",
			MessageKey:     "msg",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      outputPaths,
		ErrorOutputPaths: errPaths,
	}
}
