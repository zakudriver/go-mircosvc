package router

import (
	"github.com/gin-gonic/gin"
	"github.com/Zhan9Yunhua/blog-svr/config"
)

func Route(e *gin.Engine) {
	api := e.Group(config.SvrCfg.APIPrefix)

	user.Inject(api)

	e.NoRoute(func(ctx *gin.Context) {
		common.ErrRes(ctx, "not found")
	})
}
