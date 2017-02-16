package config

import (
	"encoding/xml"
	"github.com/devfeel/tokenserver/framework/json"
	"github.com/devfeel/tokenserver/framework/log"
	"io/ioutil"
	"os"
	"sync"
)

const ()

var (
	CurrentConfig  AppConfig
	CurrentBaseDir string
	innerLogger    *logger.InnerLogger
	redisMap       map[string]*RedisInfo
	redisMutex     *sync.RWMutex
)

func init() {
	//初始化读写锁
	redisMutex = new(sync.RWMutex)
	innerLogger = logger.GetInnerLogger()
}

func SetBaseDir(baseDir string) {
	CurrentBaseDir = baseDir
}

//初始化配置文件
func InitConfig(configFile string) *AppConfig {
	innerLogger.Info("AppConfig::InitConfig 配置文件[" + configFile + "]开始...")
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		innerLogger.Warn("AppConfig::InitConfig 配置文件[" + configFile + "]无法解析 - " + err.Error())
		os.Exit(1)
	}

	var result AppConfig
	err = xml.Unmarshal(content, &result)
	if err != nil {
		innerLogger.Warn("AppConfig::InitConfig 配置文件[" + configFile + "]解析失败 - " + err.Error())
		os.Exit(1)
	}

	CurrentConfig = result

	//初始化RedisMap
	tmpRedisMap := make(map[string]*RedisInfo)
	for k, v := range result.Redises {
		tmpRedisMap[v.ID] = &result.Redises[k]
		innerLogger.Info("AppConfig::InitConfig Load RedisInfo => " + jsonutil.GetJsonString(v))
	}

	redisMutex.Lock()
	redisMap = tmpRedisMap
	redisMutex.Unlock()

	innerLogger.Info("AppConfig::InitConfig 配置文件[" + configFile + "]完成")

	return &CurrentConfig
}

func GetRedisInfo(redisID string) (*RedisInfo, bool) {
	redisMutex.RLock()
	defer redisMutex.RUnlock()
	redis, exists := redisMap[redisID]
	return redis, exists
}
