package token

import (
	"github.com/devfeel/dotweb"
	"github.com/devfeel/tokenserver/config"
	"github.com/devfeel/tokenserver/const"
	"github.com/devfeel/tokenserver/framework/crypto"
	"github.com/devfeel/tokenserver/framework/json"
	"github.com/devfeel/tokenserver/framework/log"
	"github.com/devfeel/tokenserver/framework/redis"
	"github.com/devfeel/tokenserver/httpserver/model"
	"strconv"
)

type TokenInfo struct {
	Token       string
	AppID       string
	TokenBody   string
	LifeSeconds int //有效时间，单位为秒
}

type VerifyTokenRequest struct {
	Token       string
	AppID       string
	TokenBody   string
	IsCheckBody bool //是否需要验证Body是否一致
}

const (
	Token_DefaultLifeSeconds = 60 * 30 //token默认有效期，默认30分钟有效期
)

/*创建token
Author: panxinming
CreateTime: 2017-02-10 12:00
HttpMethod：Post
Post Param：TokenInfo
返回结构：HandlerResponse
0: 成功
-100001：post data not legal
-100002：AppID is empty
-200001：no config RedisID_Token Redis
-200002：create token error [uuid]
-200003：create token error [redis set]
Response.Message = TokenInfo

UPDATE LOG:
1、初始版本 --2017-02-10 12:00 by pxm
*/
func CreateToken(ctx *dotweb.HttpContext) {
	result := &models.HandlerResponse{RetCode: 0, RetMsg: ""}
	var tokenInfo TokenInfo
	//获取Post内容
	postContent := string(ctx.PostBody())

	defer func() {
		logger.Info("TokenServer::CreateToken["+postContent+"] => "+jsonutil.GetJsonString(result), constdefine.LogTarget_Token)
		ctx.WriteJson(result)
	}()

	//解析提交数据
	if err_jsonunmar := jsonutil.Unmarshal(postContent, &tokenInfo); err_jsonunmar != nil {
		result.RetCode = -100001
		result.RetMsg = "post data not legal => " + err_jsonunmar.Error()
		return
	}
	if tokenInfo.AppID == "" {
		result.RetCode = -100002
		result.RetMsg = "AppID is empty"
		return
	}
	//默认值处理
	if tokenInfo.LifeSeconds <= 0 {
		tokenInfo.LifeSeconds = Token_DefaultLifeSeconds
	}

	//获取redis配置
	redisInfo, exists := config.GetRedisInfo(constdefine.RedisID_Token)
	if !exists {
		result.RetCode = -200001
		result.RetMsg = "no config RedisID_Token Redis"
		return
	}

	//创建token
	token := cryptos.GetGuid()
	if token == "" {
		result.RetCode = -200002
		result.RetMsg = "create token error [uuid]"
		return
	}
	tokenInfo.Token = token

	key := redisInfo.KeyPre + ":Token:" + tokenInfo.AppID + "_" + tokenInfo.Token
	value := jsonutil.GetJsonString(tokenInfo)
	redisClient := redisutil.GetRedisClient(redisInfo.ServerIP)
	val, err := redisClient.SetWithExpire(key, value, tokenInfo.LifeSeconds)
	if err != nil {
		result.RetCode = -200003
		result.RetMsg = "create token error [redis set] => " + err.Error()
		return
	}

	result.RetCode = 0
	result.RetMsg = "ok[" + val + "]"
	result.Message = tokenInfo
	//设置原子锁，默认有效期为token的两倍，初始值为1
	key = redisInfo.KeyPre + ":TokenLock:" + tokenInfo.AppID + "_" + tokenInfo.Token
	redisClient.SetWithExpire(key, "1", tokenInfo.LifeSeconds*2)
}

/*验证token
Author: panxinming
CreateTime: 2017-02-10 12:00
HttpMethod：Post
Post Param：VerifyTokenRequest
返回结构：HandlerResponse
0: 成功
-100001：post data not legal
-100002：AppID is empty
-200001：no config RedisID_Token Redis
-201001: query token exists error
-201002: redis token not exists
-202001: get token-locker error
-202002: token-locker locked by other
-203001: query token error
-203002: redis token data not legal
-203003: token body is not match

Response.Message = TokenInfo

UPDATE LOG:
1、初始版本 --2017-02-10 12:00 by pxm
*/
func VerifyToken(ctx *dotweb.HttpContext) {
	result := &models.HandlerResponse{RetCode: 0, RetMsg: ""}
	var verifyToken VerifyTokenRequest
	//获取Post内容
	postContent := string(ctx.PostBody())

	defer func() {
		logger.Info("TokenServer::VerifyToken["+postContent+"] => "+jsonutil.GetJsonString(result), constdefine.LogTarget_Token)
		ctx.WriteJson(result)
	}()

	//解析提交数据
	err_jsonunmar := jsonutil.Unmarshal(postContent, &verifyToken)
	if err_jsonunmar != nil {
		result.RetCode = -100001
		result.RetMsg = "post data not legal => " + err_jsonunmar.Error()
		return
	}
	if verifyToken.AppID == "" {
		result.RetCode = -100002
		result.RetMsg = "AppID is empty"
		return
	}

	//获取Redis配置
	redisInfo, exists := config.GetRedisInfo(constdefine.RedisID_Token)
	if !exists {
		result.RetCode = -200001
		result.RetMsg = "no config RedisID_Token Redis"
		return
	}

	keyLocker := redisInfo.KeyPre + ":TokenLock:" + verifyToken.AppID + "_" + verifyToken.Token
	keyToken := redisInfo.KeyPre + ":Token:" + verifyToken.AppID + "_" + verifyToken.Token

	//创建Redis链接
	var redisClient *redisutil.RedisClient
	redisClient = redisutil.GetRedisClient(redisInfo.ServerIP)

	//检查token是否存在
	numExists, errExists := redisClient.Exists(keyToken)
	if errExists != nil {
		result.RetCode = -201001
		result.RetMsg = "query token exists error => " + errExists.Error()
		return
	}
	if numExists == 0 {
		result.RetCode = -201002
		result.RetMsg = "redis token not exists"
		return
	}

	//加锁
	valIncr, errIncr := redisClient.DECR(keyLocker)
	if errIncr != nil {
		result.RetCode = -202001
		result.RetMsg = "get token-locker error => " + errIncr.Error()
		return
	}
	if valIncr != 0 {
		result.RetCode = -202002
		result.RetMsg = "token-locker locked by other [" + strconv.Itoa(valIncr) + "]"
		//归还当前锁
		redisClient.INCR(keyLocker)
		return
	}

	var redisToken TokenInfo
	//检查Token
	valToken, errToken := redisClient.Get(keyToken)
	if errToken != nil {
		result.RetCode = -203001
		result.RetMsg = "query token error => " + errToken.Error()
		//归还当前锁
		redisClient.INCR(keyLocker)
		return
	}
	err_jsonunmar = jsonutil.Unmarshal(valToken, &redisToken)
	if err_jsonunmar != nil {
		result.RetCode = -203002
		result.RetMsg = "redis token data not legal [" + valToken + "] => " + err_jsonunmar.Error()
		//归还当前锁
		redisClient.INCR(keyLocker)
		return
	}
	if verifyToken.IsCheckBody {
		if verifyToken.TokenBody != redisToken.TokenBody {
			result.RetCode = -201003
			result.RetMsg = "token body is not match [" + valToken + "]"
			//归还当前锁
			redisClient.INCR(keyLocker)
			return
		}
	}

	result.RetCode = 0
	result.RetMsg = "ok"
	result.Message = redisToken
	//一切正常，删除相关token信息
	redisClient.Del(keyLocker)
	redisClient.Del(keyToken)
}

/*查询token
Author: panxinming
CreateTime: 2017-02-10 12:00
HttpMethod：Get
Get Param：
appid: 应用编码
token: token
返回结构：HandlerResponse
0: 成功
-100001：querystring token|appid is empty
-200001：no config RedisID_Token Redis
-200002：query token error [redis get]
Response.Message = TokenInfo

UPDATE LOG:
1、初始版本 --2017-02-10 12:00 by pxm
*/
func QueryToken(ctx *dotweb.HttpContext) {
	result := &models.HandlerResponse{RetCode: 0, RetMsg: ""}
	token := ctx.QueryString("token")
	appId := ctx.QueryString("appid")

	defer func() {
		logger.Info("TokenServer::QueryToken["+ctx.RawQuery()+"] => "+jsonutil.GetJsonString(result), constdefine.LogTarget_Token)
		ctx.WriteJson(result)
	}()

	if token == "" || appId == "" {
		result.RetCode = -100001
		result.RetMsg = "querystring token|appid is empty"
		return
	}

	//处理数据
	redisInfo, exists := config.GetRedisInfo(constdefine.RedisID_Token)
	if !exists {
		result.RetCode = -200001
		result.RetMsg = "no config RedisID_Token Redis"
		return
	}
	key := redisInfo.KeyPre + ":Token:" + appId + "_" + token
	redisClient := redisutil.GetRedisClient(redisInfo.ServerIP)
	val, err := redisClient.Get(key)
	if err != nil {
		result.RetCode = -200002
		result.RetMsg = "query token error [get] => " + err.Error()
		return
	}

	if val == "" {
		result.RetCode = -200003
		result.RetMsg = "query token not exists"
		return
	}

	result.RetCode = 0
	result.RetMsg = "ok"
	result.Message = val
}
