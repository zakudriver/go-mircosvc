package utils

import (
	"github.com/Zhan9Yunhua/blog-svr/common"
	"net/http"

	"github.com/Zhan9Yunhua/logger"

	"github.com/gin-gonic/gin"
)

func ErrRes(ctx *gin.Context, a ...interface{}) {
	ctx.JSON(http.StatusOK, resHandler(append(a, common.Error)))
	ctx.Abort()
}

func SerErrRes(ctx *gin.Context, a ...interface{}) {
	ctx.JSON(http.StatusOK, resHandler(append(a, common.SerError)))
	ctx.Abort()
}

func OkRes(ctx *gin.Context, a ...interface{}) {
	ctx.JSON(http.StatusOK, resHandler(append(a, common.OK)))
	ctx.Abort()
}

func resHandler(a []interface{}) (res gin.H) {
	var err error
	var code common.Status
	var msg string
	var data gin.H

	for _, v := range a {
		switch t := v.(type) {
		case error:
			err = t
		case common.Status:
			code = t
		case string:
			msg = t
		case gin.H:
			data = t
		}
	}

	if code > common.OK {
		if err == nil {
			logger.Infof(">> msg: %s.\n", msg)
		} else {
			logger.Infof(">> msg: %s. error: %s.\n", msg, err.Error())
		}
	}

	res = gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	}
	if data == nil {
		delete(res, "data")
	}
	return
}
