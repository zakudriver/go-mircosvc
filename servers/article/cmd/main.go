package main

import (
	"github.com/kum0/go-mircosvc/servers/article/config"
	"github.com/kum0/go-mircosvc/shared/logger"
	sharedZipkin "github.com/kum0/go-mircosvc/shared/zipkin"

	"fmt"

	"github.com/kum0/go-mircosvc/servers/article/endpoints"
	"github.com/kum0/go-mircosvc/shared/db"
	sharedEtcd "github.com/kum0/go-mircosvc/shared/etcd"
	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
)

func main() {
	conf := config.GetConfig()
	log, f := logger.NewLogger(conf.LogPath)
	defer f.Close()

	zipkinTracer, reporter := sharedZipkin.NewZipkin(log, conf.ZipkinAddr, "localhost:"+conf.GrpcPort,
		conf.ServiceName)
	defer reporter.Close()

	opentracing.SetGlobalTracer(zipkinot.Wrap(zipkinTracer))
	// tracer := opentracing.GlobalTracer()
	{
		etcdClient := sharedEtcd.NewEtcd(conf.EtcdAddr)
		register := sharedEtcd.Register("/usersvc", "localhost:"+conf.GrpcPort, etcdClient, log)
		defer register.Register()
	}

	var svc endpoints.ArticleServicer
	{
		mdb := db.NewMysql(conf.MysqlUsername, conf.MysqlPassword, conf.MysqlAddr, conf.MysqlAuthsource)
		svc = endpoints.NewArticleService(mdb)
	}

	fmt.Println(svc)
}
