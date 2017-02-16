package main

import (
	"flag"
	"fmt"
	"github.com/devfeel/tokenserver/config"
	"github.com/devfeel/tokenserver/framework/file"
	"github.com/devfeel/tokenserver/framework/log"
	"github.com/devfeel/tokenserver/httpserver"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var (
	innerLogger           *logger.InnerLogger
	configFile            string
	DefaultConfigFileName string
)

func init() {
	innerLogger = logger.GetInnerLogger()
	DefaultConfigFileName = fileutil.GetCurrentDirectory() + "/app.conf"
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			strLog := "TechService.tokenserver:main recover error => " + fmt.Sprintln(err)
			os.Stdout.Write([]byte(strLog))
			innerLogger.Error(strLog)

			buf := make([]byte, 4096)
			n := runtime.Stack(buf, true)
			innerLogger.Error(string(buf[:n]))
			os.Stdout.Write(buf[:n])
		}
	}()

	currentBaseDir := fileutil.GetCurrentDirectory()
	flag.StringVar(&configFile, "config", "", "配置文件路径")
	if configFile == "" {
		configFile = DefaultConfigFileName
	}

	//启动内部日志服务
	logger.StartInnerLogHandler(currentBaseDir)

	//加载xml配置文件
	appconfig := config.InitConfig(configFile)

	//设置基本目录
	config.SetBaseDir(currentBaseDir)

	//启动日志服务
	logger.StartLogHandler(appconfig.Log.FilePath)

	//监听系统信号
	go listenSignal()

	err := httpserver.StartServer()
	if err != nil {
		innerLogger.Warn("HttpServer.StartServer失败 " + err.Error())
		fmt.Println("HttpServer.StartServer失败 " + err.Error())
	}

}

func listenSignal() {
	c := make(chan os.Signal, 1)
	//syscall.SIGSTOP
	signal.Notify(c, syscall.SIGHUP)
	for {
		s := <-c
		innerLogger.Info("signal::ListenSignal [" + s.String() + "]")
		switch s {
		case syscall.SIGHUP: //配置重载
			innerLogger.Info("signal::ListenSignal reload config begin...")
			//重新加载配置文件
			config.InitConfig(configFile)
			innerLogger.Info("signal::ListenSignal reload config end")
		default:
			return
		}
	}
}
