package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"solana-hackthon-cli/ca"
	"strings"
)

// AI模型起始端口
const AGENT_START_PORT int64 = 10100

var index_Port = AGENT_START_PORT

var agentMap = make(map[string]ca.Agent)
var agentTempMap = make(map[string]int)
var hashPortMap = make(map[string]int64)

func startAgents(containers []types.Container) {
	// 合约拉取agentinfo
	agents := ca.Agents()
	// agentinfo设置map缓存
	for _, a := range agents {
		agentMap[a.ModelHash] = a
		agentTempMap[a.ModelHash] = 0
	}
	// 记录已经跑起来的容器和端口，创建的时候略过
	for _, container := range containers {
		cname := strings.TrimLeft(container.Names[0], "/")
		if agent, ok := agentMap[cname]; ok {
			port := int64(container.Ports[0].PublicPort)
			hashPortMap[agent.ModelHash] = port
			agentTempMap[agent.ModelHash] = 1
			if index_Port <= port {
				index_Port = port + 1
			}
			continue
		}
	}
	// 从现有容器中最大的port往后跑
	for k, v := range agentTempMap {
		if v == 1 {
			continue
		}
		agent := agentMap[k]
		dockerfileds := DockerFileds{
			Image:         agent.DockerImageHref,
			Port:          agent.ApiDefaultPort,
			HostPort:      index_Port,
			ContainerName: agent.ModelHash,
		}
		err := startImage(dockerfileds)
		if err != nil {
			fmt.Printf("create ai agent [%s] canister failed, msg: %s\n", agent.ModelHash, err.Error())
			continue
		}
		index_Port = index_Port + 1
	}
}

func GetPort(hash string) int64 {
	return hashPortMap[hash]
}
