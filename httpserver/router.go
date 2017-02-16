package httpserver

import (
	"github.com/devfeel/dotweb"
	"github.com/devfeel/tokenserver/httpserver/handlers/global"
	"github.com/devfeel/tokenserver/httpserver/handlers/token"
)

func InitRoute(dotweb *dotweb.Dotweb) {
	//token
	dotweb.HttpServer.POST("/token/create", token.CreateToken)
	dotweb.HttpServer.POST("/token/verify", token.VerifyToken)
	dotweb.HttpServer.GET("/token/query", token.QueryToken)
	//global
	dotweb.HttpServer.GET("/global/createid", global.CreateGlobalID)
}
