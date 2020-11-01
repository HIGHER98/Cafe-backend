package app

import (
	"cafe/pkg/e"

	"github.com/gin-gonic/gin"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Response setting gin.JSON
func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	g.C.Header("Access-Control-Allow-Origin", "*")
	g.C.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
	if g.C.Request.Method == "OPTIONS" {
		g.C.AbortWithStatus(204)
		return
	}
	g.C.JSON(httpCode, Response{
		Code: errCode,
		Msg:  e.GetMsg(errCode),
		Data: data,
	})
	return
}
