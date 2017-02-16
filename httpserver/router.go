package httpserver

import (
	"devfeel/tokenserver/httpserver/handlers/global"
	"devfeel/tokenserver/httpserver/handlers/token"
	"github.com/devfeel/dotweb"
)

func InitRoute(dotweb *dotweb.Dotweb) {
	//token
	dotweb.HttpServer.POST("/token/create", token.CreateToken)
	dotweb.HttpServer.POST("/token/verify", token.VerifyToken)
	dotweb.HttpServer.GET("/token/query", token.QueryToken)
	//global
	dotweb.HttpServer.GET("/global/createid", global.CreateGlobalID)
}
