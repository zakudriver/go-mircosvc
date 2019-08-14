package logger

import (
	"fmt"
	"github.com/Zhan9Yunhua/blog-svr/servers/user/config"
	"github.com/go-kit/kit/log"
	"io"
	"os"
	"path/filepath"

	lg "github.com/Zhan9Yunhua/logger"
)

func NewLogger() log.Logger {
	logger, err := handleLogger()
	if err != nil {
		lg.Fatalln(err)
	}

	return logger
}

func handleLogger() (log.Logger, error) {
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
