package main

import (
	"github.com/kum0/go-mircosvc/servers/article/config"
	"github.com/kum0/go-mircosvc/shared/logger"
	sharedZipkin "github.com/kum0/go-mircosvc/shared/zipkin"
)

func main() {
	conf := config.GetConfig()
	log, f := logger.NewLogger(conf.LogPath)
	defer f.Close()

	_, reporter := sharedZipkin.NewZipkin(log, conf.ZipkinAddr, "localhost:"+conf.GrpcPort,
		conf.ServiceName)
	defer reporter.Close()
}
