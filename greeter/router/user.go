package router

import (
	"github.com/gin-gonic/gin"

	userPb "github.com/Zhan9Yunhua/blog-svr/servers/user/proto/user"
	"github.com/micro/go-micro/client"
)

const (
	userPrefix = "user"
)

var (
	userCl userPb.SayService
)


func InjectUserRouter(router *gin.RouterGroup) {
	userCl = userPb.NewSayService("go.micro.greeter", client.DefaultClient)
	gp := router.Group(userPrefix)

	user := new(UserRouter)

	gp.POST("/register", user.register)
}

type UserRouter struct {
}

func (user *UserRouter) register(ctx *gin.Context) {

}
