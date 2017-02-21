# TokenServer
token服务，提供token一致性服务以及相关的全局ID生成服务等

##API列表
* <a href="https://github.com/devfeel/tokenserver/blob/master/README/tokenmessage.MD#1create-token">Create Token</a>
* <a href="https://github.com/devfeel/tokenserver/blob/master/README/tokenmessage.MD#2query-token">Query Token</a>
* <a href="https://github.com/devfeel/tokenserver/blob/master/README/tokenmessage.MD#3verify-token">Verify Token</a>
* <a href="https://github.com/devfeel/tokenserver/blob/master/README/globalmessage.MD#1create-globalid">Create GlobalID</a>

##配置说明 
```
<?xml version="1.0" encoding="UTF-8"?>
<config>
<httpserver httpport="8201" pprofport="8202"></httpserver>
<log filepath="/home/devfeel/tokenserver/logs"></log>
<redises>
    <redis id="tokenredis" serverip="127.0.0.1:6379" keypre="devfeel:TokenServer"></redis>
    <redis id="globalredis" serverip="127.0.0.1:6379" keypre="devfeel:TokenServer"></redis>
</redises>
</config>
```

##运行说明
假设执行文件安装在/home/devfeel/tokenserver/目录下：
<br>
* tokenserver    可执行文件
* innerlogs      程序内部日志目录
* logs           程序业务日志目录
* app.conf       程序配置文件

*默认配置下，会监听两个端口，一个为业务端口，一个为pprof端口


##外部依赖
* mgo - gopkg.in/mgo.v2/bson
* redigo - github.com/garyburd/redigo/redis
