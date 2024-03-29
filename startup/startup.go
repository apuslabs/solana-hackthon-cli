package startup

import (
	"solana-hackthon-cli/ca"
	"solana-hackthon-cli/computer"
	"solana-hackthon-cli/computer/monitor"
	"solana-hackthon-cli/config"
)

// 启动预设，检查环境信息。查询计算机配置，生成节点keypaire， 注册节点信息(如果第一次启动的话)
func Startup() {
	config.Init()
	ca.Init()
	gpuNode := computer.Init()
	RegisterGpuNode(gpuNode)
}

func RegisterGpuNode(gpuNode monitor.GpuNode) {
	gpuNode.Owner = config.OwnerPubkey
	gpuNode.Id = ca.GetPubkey()
	gpuNode.Price = config.Price
	gpuNode.Endpoint = config.Endpoint
	// http request
}
