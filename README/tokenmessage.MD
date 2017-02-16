##1、Create Token:
描述：创建Token
<br>接口: http://127.0.0.1/token/create
<br>HttpMethod：Post
<br>Post Param：TokenInfo
<br>返回结构：HandlerResponse
* 0: 成功
* -100001：post data not legal
* -100002：AppID is empty
* -200001：no config RedisID_Token Redis
* -200002：create token error [uuid]
* -200003：create token error [redis set]
* Response.Message = TokenInfo

其中TokenInfo结构：
```go
type TokenInfo struct {
    Token       string
    AppID       string
    TokenBody   string
    LifeSeconds int //有效时间，单位为秒，如果输入不合法，默认为1800秒
}
```
测试数据：
```
{
    "Token":"",
    "AppID":"10000",
    "TokenBody":"1",
    "LifeSeconds":1800
}
```

result:
```
{"RetCode":0,"RetMsg":"ok[OK]","Message":{"Token":"a8d3400acb3217e10e946d1e0ca107db","AppID":"10000","TokenBody":"1","LifeSeconds":30}}
```

**************************************************************************************************

##2、Query Token:
描述：查询token,不影响token有效性，调用完token依然存在，如需一次性token验证，请调用verify方法
<br>接口:http://192.168.8.178:8201/token/query?appid=?&token=?
<br>HttpMethod：Get
<br>Get Param：
* appid: 应用编码
* token: token
<br>返回结构：HandlerResponse
* 0: 成功
* -100001：querystring token|appid is empty
* -200001：no config RedisID_Token Redis
* -200002：query token error [redis get]
* Response.Message = TokenInfo

result:
```
{"RetCode":0,"RetMsg":"ok","Message":"{\"Token\":\"a8d3400acb3217e10e946d1e0ca107db\",\"AppID\":\"10000\",\"TokenBody\":\"1\",\"LifeSeconds\":1800}"}
```
其中TokenInfo结构：
```go
type TokenInfo struct {
    Token       string
    AppID       string
    TokenBody   string
    LifeSeconds int //有效时间，单位为秒
}
```

**************************************************************************************************

##3、Verify Token:
描述：验证token,一次有效，调用完token移除，如需多次验证，请调用query方法
<br>接口：http://192.168.8.178:8201/token/verify
<br>HttpMethod：Post
<br>Post Param：VerifyTokenRequest
<br>返回结构：HandlerResponse
* 0: 成功
* -100001：post data not legal
* -100002：AppID is empty
* -200001：no config RedisID_Token Redis
* -201001: query token exists error
* -201002: redis token not exists
* -202001: get token-locker error
* -202002: token-locker locked by other
* -203001: query token error
* -203002: redis token data not legal
* -203003: token body is not match
* Response.Message = TokenInfo

其中VerifyTokenRequest结构：
```go
type VerifyTokenRequest struct {
    Token       string
    AppID       string
    TokenBody   string
    IsCheckBody bool //是否需要验证Body是否一致
}
```
测试数据：
```
{
    "Token":"a8d3400acb3217e10e946d1e0ca107db",
    "AppID":"10000",
    "TokenBody":"1",
    "IsCheckBody":true
}
```
result:
```
{"RetCode":0,"RetMsg":"ok","Message":{"Token":"9d8ccbdd35a61aacadf80ea91e4f90bb","AppID":"10000","TokenBody":"1","LifeSeconds":1800}}
```

其中TokenInfo结构：
```go
type TokenInfo struct {
    Token       string
    AppID       string
    TokenBody   string
    LifeSeconds int //有效时间，单位为秒
}
```
