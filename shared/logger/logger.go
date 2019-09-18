package logger

import (
	"io"
	"os"
	"path/filepath"

	"github.com/go-kit/kit/log"

	lg "github.com/Zhan9Yunhua/logger"
)

func NewLogger(path string) log.Logger {
	logger, err := handleLogger(path)
	if err != nil {
		lg.Fatalln(err)
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

	var kitlogger log.Logger
	{
		kitlogger = log.NewLogfmtLogger(io.MultiWriter(os.Stderr, logfile))
		kitlogger = log.With(kitlogger, "ts", log.DefaultTimestampUTC)
	}

	return log.With(kitlogger, "component", "http"), nil
}
