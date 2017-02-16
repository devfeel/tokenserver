package config

import (
	"encoding/xml"
)

//配置信息
type AppConfig struct {
	XMLName    xml.Name    `xml:"config"`
	Log        Log         `xml:"log"`
	HttpServer HttpServer  `xml:"httpserver"`
	Redises    []RedisInfo `xml:"redises>redis"`
}

//全局配置
type HttpServer struct {
	HttpPort  int `xml:"httpport,attr"`
	PProfPort int `xml:"pprofport,attr"`
}

//log配置
type Log struct {
	FilePath string `xml:"filepath,attr"`
}

//Redis信息
type RedisInfo struct {
	ID       string `xml:"id,attr"`
	ServerIP string `xml:"serverip,attr"`
	KeyPre   string `xml:"keypre,attr"`
}
