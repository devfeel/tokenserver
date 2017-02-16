##1、Create GlobalID:
创建全局ID
<br>api-url: http://127.0.0.1/global/createid?idtype=?&appid=?&module=?
<br>HttpMethod：Get
<br>Get Param：
* idtype: 代表请求ID方式，目前支持：UUID、Number、TimeNumber
* appid: 代表请求的应用
* module: 代表请求应用某特定模块
<br>返回结构：HandlerResponse
* 0: 成功
* -100001：idtype is empty
* -100002：appid & module is empty
* -201001：no config RedisID_Global Redis
* -201002：create global number error [redis]
* Response.Message = GlobalID //string
<br>result:
* UUID - {"RetCode":0,"RetMsg":"ok","Message":"b6b81359c127554ee2a34bb36457026b"}
* Number - {"RetCode":0,"RetMsg":"","Message":"64"}
* TimeNumber - {"RetCode":0,"RetMsg":"","Message":"2017021611304200000002"}
