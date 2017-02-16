package httpserver

import (
	"devfeel/tokenserver/config"
	"devfeel/tokenserver/framework/log"
	"fmt"
	"github.com/devfeel/dotweb"
	"strconv"
)

func StartServer() error {

	//初始化DotServer
	dotweb := dotweb.New()

	//设置dotserver日志目录
	dotweb.SetLogPath(config.CurrentConfig.Log.FilePath)

	//设置路由
	InitRoute(dotweb)

	innerLogger := logger.GetInnerLogger()

	//启动监控服务
	pprofport := config.CurrentConfig.HttpServer.PProfPort
	go dotweb.StartPProfServer(pprofport)

	// 开始服务
	port := config.CurrentConfig.HttpServer.HttpPort
	innerLogger.Debug("dotweb.StartServer => " + strconv.Itoa(port))
	err := dotweb.StartServer(port)
	return err
}

func ReSetServer() {
	//初始化应用信息
	fmt.Println("ReSetServer")
}
