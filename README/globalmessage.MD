##1、Create GlobalID:
描述：创建全局ID
<br>接口: http://127.0.0.1/global/createid?idtype=?&appid=?&module=?
<br>HttpMethod：Get
Get Param：
* idtype: 代表请求ID方式，目前支持：uuid、number、timenumber、mongodb
* appid: 代表请求的应用
* module: 代表请求应用某特定模块

返回结构：HandlerResponse
* 0: 成功
* -100001：idtype is empty
* -100002：appid & module is empty
* -201001：no config RedisID_Global Redis
* -201002：create global number error [redis]
* Response.Message = GlobalID //string

result:
* UUID - {"RetCode":0,"RetMsg":"ok","Message":"b6b81359c127554ee2a34bb36457026b"}
* Number - {"RetCode":0,"RetMsg":"","Message":"64"}
* TimeNumber - {"RetCode":0,"RetMsg":"","Message":"2017021611304200000002"}
* MongoDB - {"RetCode":0,"RetMsg":"ok","Message":"58abd9dae138236e64ebfdf8"}
