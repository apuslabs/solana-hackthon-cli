package docker

import (
	"github.com/docker/docker/api/types"
	"strings"
)

const image_name = "telegraf"
const image_port = 8086
const host_port = 10001
const container_name = "telegraf"

// 启动监控程序
func startTelegraf(containers []types.Container) {
	// 判断是否已经有telegraf
	flag := false
	for _, container := range containers {
		if strings.Contains(container.Names[0], container_name) {
			flag = true
		}
	}
	if !flag {
		dockerfileds := DockerFileds{
			Image:         image_name,
			Port:          image_port,
			HostPort:      host_port,
			ContainerName: container_name,
		}
		err := startImage(dockerfileds)
		if err != nil {
			panic(err)
		}
	}
}
