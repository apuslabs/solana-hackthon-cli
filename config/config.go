package config

import (
	"flag"
)

var OwnerPubkey string
var ServerAddress string
var Price int64
var Endpoint string

var ProgramId string

func init() {
	flag.StringVar(&OwnerPubkey, "ownerpubkey", "", "node owner publickey")
	flag.StringVar(&ServerAddress, "serveraddress", "", "register server address: https://host:port")
	flag.StringVar(&ProgramId, "programid", "", "apus solana contract programid")
	flag.Int64Var(&Price, "price", 0, "price for ai task")
	flag.StringVar(&Endpoint, "endpoint", "", "endpoint for access this node: ip/domain")
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

	if ProgramId == "" {
		panic("programid must be set: apus solana contract programid")
	}

	if Endpoint == "" {
		panic("programid must be set: endpoint for access this node: ip/domain")
	}
}
