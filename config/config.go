package config

import (
	"flag"
)

var OwnerPubkey string
var ServerAddress string
var Price string
var Endpoint string

var SkipGpu bool

func init() {
	flag.StringVar(&OwnerPubkey, "ownerpubkey", "", "node owner publickey")
	flag.StringVar(&ServerAddress, "serveraddress", "", "register server address: https://host:port")
	flag.StringVar(&Price, "price", "0", "price for ai task")
	flag.StringVar(&Endpoint, "endpoint", "", "endpoint for access this node: ip/domain")
	flag.BoolVar(&SkipGpu, "skipgpu", false, "")
}

// 定义配置，读取配置。提供全局配置调用方法
func Init() {
	flag.Parse()

	if OwnerPubkey == "" {
		panic("ownerpubkey must be set: node owner publickey")
	}
	if ServerAddress == "" {
		panic("serveraddress must be set: register server address: https://host:port")
	}

	if Endpoint == "" {
		panic("programid must be set: endpoint for access this node: ip/domain")
	}
}
