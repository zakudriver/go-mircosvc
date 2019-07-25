package router

import "github.com/gin-gonic/gin"

type UserRouter struct {
	Prefix string
}

func (user *UserRouter)Inject(route *gin.RouterGroup) {
	gp := route.Group(user.Prefix)

	gp.POST("/signin", user.signin)
}

func (user *UserRouter)signin(ctx *gin.Context)  {

}
