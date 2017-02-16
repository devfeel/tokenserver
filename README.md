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
