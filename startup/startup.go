package startup

import (
	"solana-hackthon-cli/ca"
	"solana-hackthon-cli/computer"
	"solana-hackthon-cli/computer/docker"
	"solana-hackthon-cli/config"
)

// 启动预设，检查环境信息。查询计算机配置，生成节点keypaire， 注册节点信息(如果第一次启动的话)
func Startup() {
	config.Init()
	ca.Init()
	computer.Init()
	docker.Init()
	RegisterGpuNode()
}

func RegisterGpuNode() {

}
