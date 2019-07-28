package main

import (
	"time"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/api"
	"github.com/micro/go-micro/server"

	ha "github.com/micro/go-api/handler/http"
)

func main() {
	service := micro.NewService(
		micro.Name("go.micro.api.user"),
		micro.Version("latest"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*15),
	)

	service.Server().Init(
		server.Wait(nil),
	)

	// Register Handler
	example.RegisterExampleHandler(service.Server(), new(handler.Example),
		api.WithEndpoint(&api.Endpoint{
			Name: "User.Register",
			Path: []string{"/user/register"},
			Method: []string{"POST"},
			Handler: ha.Handler,
		}))
}