package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"strings"
)

// AI模型起始端口
const AGENT_START_PORT int64 = 9000

var index_Port = AGENT_START_PORT

var agentMap = make(map[string]Agent)
var agentTempMap = make(map[string]int)
var hashPortMap = make(map[string]int64)

type Agent struct {
	Owner           string `json:"owner"`
	Post            string `json:"post"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	ModelHash       string `json:"model_hash"` // container name: imageherf-hash
	ModelType       string `json:"model_type"`
	ApiType         string `json:"api_type"`
	ApiDoc          string `json:"api_doc"`
	Price           string `json:"price"`
	DockerImageHref string `json:"docker_image_href"` // image name
	ApiDefaultPort  int64  `json:"api_default_port"`  // image port
}

// 查询agent信息
func Agents() []Agent {
	return []Agent{}
}

func startAgents(containers []types.Container) {
	// 合约拉取agentinfo
	agents := Agents()
	// agentinfo设置map缓存
	for _, agent := range agents {
		agentMap[agent.ModelHash] = agent
		agentTempMap[agent.ModelHash] = 0
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
