package logger

import (
	"io"
	"os"
	"path/filepath"

	"github.com/go-kit/kit/log/level"

	"github.com/go-kit/kit/log"
)

func NewLogger(path string) (log.Logger, *os.File) {
	path, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(io.MultiWriter(os.Stderr, f))
		logger = level.NewFilter(logger, level.AllowInfo())
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	return logger, f
}

func handleLogger(logPath string) (log.Logger, error) {
	path, err := filepath.Abs(logPath)
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(io.MultiWriter(os.Stderr, f))
		logger = level.NewFilter(logger, level.AllowInfo())
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	return logger, nil
}
