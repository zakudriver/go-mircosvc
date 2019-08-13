package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Zhan9Yunhua/blog-svr/gateway/config"
	"github.com/go-kit/kit/log"
)

func NewLogger() (log.Logger, error) {
	conf := config.GetConfig()

	path, err := filepath.Abs(conf.LogPath)
	if err != nil {
		return nil, err
	}
	fmt.Println(path)

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
