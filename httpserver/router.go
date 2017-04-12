package httpserver

import (
	"github.com/devfeel/dotweb"
	"github.com/devfeel/tokenserver/httpserver/handlers/global"
	"github.com/devfeel/tokenserver/httpserver/handlers/token"
)

func InitRoute(dotweb *dotweb.DotWeb) {
	//token
	dotweb.HttpServer.Router().POST("/token/create", token.CreateToken)
	dotweb.HttpServer.Router().POST("/token/verify", token.VerifyToken)
	dotweb.HttpServer.Router().GET("/token/query", token.QueryToken)
	//global
	dotweb.HttpServer.Router().GET("/global/createid", global.CreateGlobalID)
}
