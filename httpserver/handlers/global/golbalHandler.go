package global

import (
	"devfeel/tokenserver/config"
	"devfeel/tokenserver/const"
	"devfeel/tokenserver/framework/crypto"
	"devfeel/tokenserver/framework/json"
	"devfeel/tokenserver/framework/log"
	"devfeel/tokenserver/framework/redis"
	"devfeel/tokenserver/httpserver/model"
	"fmt"
	"github.com/devfeel/dotweb"
	"strconv"
	"time"
)

const (
	IdType_UUID       = "uuid"
	IdType_Number     = "number"
	IdType_TimeNumber = "timenumber"
)

/*创建GlobalID
Author: panxinming
CreateTime: 2017-02-16 11:00
HttpMethod：Get
Route Param:
idtype: 代表请求ID方式，目前支持：UUID、Number、TimeNumber
Get Param：
appid: 代表请求的应用
module: 代表请求应用某特定模块

返回结构：HandlerResponse
0: 成功
-100001：idtype is empty
-100002：appid & module is empty
-201001：no config RedisID_Global Redis
-201002：create global number error [redis]
Response.Message = GlobalID //string

UPDATE LOG:
1、初始版本 --2017-02-16 11:00 by pxm
*/
func CreateGlobalID(ctx *dotweb.HttpContext) {
	result := &models.HandlerResponse{RetCode: 0, RetMsg: ""}
	idtype := ctx.QueryString("idtype")
	appid := ctx.QueryString("appid")   //代表请求的应用
	module := ctx.QueryString("module") //代表请求应用某特定模块
	var code string

	defer func() {
		logger.Info("TokenServer::CreateGlobalID ["+idtype+"] ["+appid+"] ["+module+"] => "+jsonutil.GetJsonString(result), constdefine.LogTarget_Global)
		ctx.WriteJson(result)
	}()

	if idtype == "" {
		result.RetCode = -100001
		result.RetMsg = "idtype is empty"
		return
	}
	if appid == "" || module == "" {
		result.RetCode = -100002
		result.RetMsg = "appid & module is empty"
		return
	}

	if idtype == IdType_UUID {
		code = cryptos.GetGuid()
		result.RetCode = 0
		result.RetMsg = "ok"
		result.Message = code
		return
	}

	if idtype == IdType_TimeNumber {
		createTimeNumberID(appid, module, result)
		return
	}

	if idtype == IdType_Number {
		createNumberID(appid, module, result)
		return
	}

	//no match idtype
	result.RetCode = -200001
	result.RetMsg = "idtype[" + idtype + "] is not support"
}

//基于Redis创建连续数字
func createNumberID(appid, module string, result *models.HandlerResponse) {
	code := ""

	//获取redis配置
	redisInfo, exists := config.GetRedisInfo(constdefine.RedisID_Global)
	if !exists {
		result.RetCode = -201001
		result.RetMsg = "no config RedisID_Global redis"
		return
	}

	key := redisInfo.KeyPre + ":Global_Number:" + appid + ":" + module
	redisClient := redisutil.GetRedisClient(redisInfo.ServerIP)
	val, err := redisClient.INCR(key)
	if err != nil {
		result.RetCode = -201002
		result.RetMsg = "create global number error [redis] => " + err.Error()
		return
	}

	code = strconv.Itoa(val)
	result.Message = code
}

//基于Redis创建TimeNumber
func createTimeNumberID(appid, module string, result *models.HandlerResponse) {
	code := ""
	timeLayout := "20060102150405"

	//获取redis配置
	redisInfo, exists := config.GetRedisInfo(constdefine.RedisID_Global)
	if !exists {
		result.RetCode = -202001
		result.RetMsg = "no config RedisID_Global redis"
		return
	}

	key := redisInfo.KeyPre + ":Global_TimeNumber:" + appid + ":" + module
	redisClient := redisutil.GetRedisClient(redisInfo.ServerIP)
	val, err := redisClient.INCR(key)
	if err != nil {
		result.RetCode = -202002
		result.RetMsg = "create global number error [redis] => " + err.Error()
		return
	}
	//创建8位补0字符串
	code = time.Now().Format(timeLayout) + fmt.Sprintf("%08d", val)
	result.Message = code
}
