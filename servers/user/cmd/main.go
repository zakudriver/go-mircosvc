package main

import (
	"context"
	"log"

	"github.com/micro/go-micro"

	user "github.com/Zhan9Yunhua/blog-svr/servers/user/proto/user"
	// "github.com/micro/go-micro/api"


	"github.com/micro/go-api"
	proto "github.com/micro/go-api/proto"
	rapi "github.com/micro/go-micro/api/handler/api"

	apip "github.com/micro/go-micro/api/proto"
)

func main() {
	service := micro.NewService(
		micro.Name("go.micro.api.user"),
	)

	service.Init()

	// Register Handler
	// user.RegisterUserHandler(service.Server(), new(Say), api.WithEndpoint(
	// 	&api.Endpoint{
	// 		Name:"Aa.Login",
	// 		Path: []string{"/aa"},
	// 		Method: []string{"GET"},
	// 		Handler: rapi.Handler,
	// 	},
	// ))

	user.RegisterUserHandler(service.Server(), new(Say), api.WithEndpoint(
		&api.Endpoint{
			Name:    "Aa.Login",
			Path:    []string{"/aa"},
			Method:  []string{"GET"},
			Handler: rapi.Handler,
		},
	))

	// service.Server().Handle(service.Server().NewHandler(new(Say)))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

type Say struct {
}

func (s *Say) Login(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	log.Print("Received Say.Hello API request")
	log.Println(req.Method)

	// b, _ := json.Marshal(map[string]string{
	// 	"message": "ok",
	// 	"method": req.Method,
	// })

	return nil
}
