package utils

import (
	"net/http"

	"github.com/Zhan9Yunhua/logger"

	"github.com/gin-gonic/gin"
)

func ErrRes(ctx *gin.Context, a ...interface{}) {
	ctx.JSON(http.StatusOK, resHandler(append(a, Error)))
	ctx.Abort()
}

func SerErrRes(ctx *gin.Context, a ...interface{}) {
	ctx.JSON(http.StatusOK, resHandler(append(a, SerError)))
	ctx.Abort()
}

func OkRes(ctx *gin.Context, a ...interface{}) {
	ctx.JSON(http.StatusOK, resHandler(append(a, OK)))
	ctx.Abort()
}

func resHandler(a []interface{}) (res gin.H) {
	var err error
	var code Status
	var msg string
	var data gin.H

	for _, v := range a {
		switch t := v.(type) {
		case error:
			err = t
		case Status:
			code = t
		case string:
			msg = t
		case gin.H:
			data = t
		}
	}

	if code > OK {
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
