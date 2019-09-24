package logger

import (
	"github.com/go-kit/kit/log/level"
	"io"
	"os"
	"path/filepath"

	"github.com/go-kit/kit/log"
)

func NewLogger(path string) log.Logger {
	logger, err := handleLogger(path)
	if err != nil {
		panic(err)
	}

	return logger
}

func handleLogger(logPath string) (log.Logger, error) {
	path, err := filepath.Abs(logPath)
	if err != nil {
		return nil, err
	}

	logfile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	defer logfile.Close()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(io.MultiWriter(os.Stderr, logfile))
		logger = level.NewFilter(logger, level.AllowInfo())
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
		// logger = log.With(logger, "component", "http")
	}

	return logger, nil
}
