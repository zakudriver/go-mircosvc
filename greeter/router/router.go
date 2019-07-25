package router

import (
	"github.com/Zhan9Yunhua/blog-svr/config"
	"github.com/Zhan9Yunhua/blog-svr/utils"
	"github.com/gin-gonic/gin"
)

func NewRoute(e *gin.Engine) {
	gp := e.Group(config.SvrCfg.APIPrefix)

	injectRouter(gp)

	e.NoRoute(func(ctx *gin.Context) {
		utils.ErrRes(ctx, "not found")
	})
}

func injectRouter(group *gin.RouterGroup) {
	InjectUserRouter(group)
}
