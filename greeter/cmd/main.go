// package main
//
// import (
// 	"fmt"
// 	"io"
// 	"os"
//
// 	"github.com/Zhan9Yunhua/blog-svr/config"
// 	"github.com/Zhan9Yunhua/blog-svr/greeter/router"
// 	"github.com/Zhan9Yunhua/blog-svr/utils"
// 	"github.com/Zhan9Yunhua/logger"
// 	"github.com/gin-gonic/gin"
// 	"github.com/micro/go-micro/web"
//
// )
//
// func main() {
// 	// service := web.NewService(
// 	// 	web.Name("go.micro.web.greeter"),
// 	// )
// 	//
// 	// service.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 	// 	if r.Method == http.MethodGet {
// 	// 		r.ParseForm()
// 	//
// 	// 		name := r.Form.Get("name")
// 	// 		// if len(name) == 0 {
// 	// 		// 	name = "World"
// 	// 		// }
// 	//
// 	// 		cl := user.NewSayService("go.micro.srv.greeter", client.DefaultClient)
// 	// 		rsp, err := cl.Hello(context.Background(), &user.Request{
// 	// 			Name: name,
// 	// 		})
// 	//
// 	// 		if err != nil {
// 	// 			http.Error(w, err.Error(), 500)
// 	// 			return
// 	// 		}
// 	//
// 	// 		// w.Write([]byte(`<html><body><h1>` + rsp.Msg + `</h1></body></html>`))
// 	//
// 	// 		r, _ := json.Marshal(map[string]interface{}{"name": rsp.Msg})
// 	// 		w.Write(r)
// 	// 		return
// 	// 	}
// 	//
// 	// 	fmt.Fprint(w, `<html><body><h1>Enter Name<h1><form method=post><input name=name type=text /></form></body></html>`)
// 	// })
// 	//
// 	// if err := service.Init(func(o *web.Options) {
// 	// 	o.Address = ":8080"
// 	// }); err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	//
// 	// if err := service.Run(); err != nil {
// 	// 	log.Fatal(err)
// 	// }
//
// 	service := web.NewService(
// 		web.Name("go.micro.api.greeter"),
// 	)
//
// 	service.Init()
//
//
// 	fmt.Println("Gin Version:", gin.Version)
//
// 	if config.SvrCfg.Env != config.DevelopmentMode {
// 		gin.SetMode(gin.ReleaseMode)
//
// 		gin.DisableConsoleColor()
//
// 		f, err := utils.SafeOpenFile(config.SvrCfg.LogDir, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
// 		if err != nil {
// 			logger.Fatalln(err)
// 		}
// 		gin.DefaultWriter = io.MultiWriter(f)
// 	}
//
// 	app := gin.New()
//
// 	app.Use(gin.Logger())
//
// 	app.Use(gin.Recovery())
//
// 	router.NewRoute(app)
//
// 	service.Handle("/", app)
//
// 	if err := service.Run(); err != nil {
// 		logger.Fatalln(err)
// 	}
// }


package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/web"
)

// exampleCall will handle /example/call
func exampleCall(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// get name
	name := r.Form.Get("name")

	if len(name) == 0 {
		http.Error(
			w,
			errors.BadRequest("go.micro.api.example", "no content").Error(),
			400,
		)
		return
	}

	// marshal response
	b, _ := json.Marshal(map[string]interface{}{
		"message": "got your message " + name,
	})

	// write response
	w.Write(b)
}

// exampleFooBar will handle /example/foo/bar
func exampleFooBar(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(
			w,
			errors.BadRequest("go.micro.api.example", "require post").Error(),
			400,
		)
		return
	}

	if len(r.Header.Get("Content-Type")) == 0 {
		http.Error(
			w,
			errors.BadRequest("go.micro.api.example", "need content-type").Error(),
			400,
		)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(
			w,
			errors.BadRequest("go.micro.api.example", "expect application/json").Error(),
			400,
		)
		return
	}

	// do something
}

func main() {
	// we're using go-web for convenience since it registers with discovery
	service := web.NewService(
		web.Name("go.micro.api.example"),
	)

	service.HandleFunc("/example/call", exampleCall)
	service.HandleFunc("/example/foo/bar", exampleFooBar)

	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
