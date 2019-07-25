package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/web"
	"log"
	"net/http"

	user "github.com/Zhan9Yunhua/blog-svr/servers/user/proto/user"
)

func main() {
	service := web.NewService(
		web.Name("go.micro.web.greeter"),
	)

	service.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			r.ParseForm()

			name := r.Form.Get("name")
			// if len(name) == 0 {
			// 	name = "World"
			// }

			cl := user.NewSayService("go.micro.srv.greeter", client.DefaultClient)
			rsp, err := cl.Hello(context.Background(), &user.Request{
				Name: name,
			})

			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			// w.Write([]byte(`<html><body><h1>` + rsp.Msg + `</h1></body></html>`))

			r, _ := json.Marshal(map[string]interface{}{"name": rsp.Msg})
			w.Write(r)
			return
		}

		fmt.Fprint(w, `<html><body><h1>Enter Name<h1><form method=post><input name=name type=text /></form></body></html>`)
	})

	if err := service.Init(func(o *web.Options) {
		o.Address = ":8080"
	}); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
