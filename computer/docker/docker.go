package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"strconv"
	"strings"
	"time"
)

var dockerClient *client.Client

type DockerFileds struct {
	Image         string
	Port          int64
	HostPort      int64
	ContainerName string // agent hash
}

// 拉取和管理docker容器，缓存容器：port映射 host网络模式
func Init() {
	dc, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println(" Unable to create docker client; msg: ", err.Error())
		panic(err)
	}
	dockerClient = dc
	containers, err := dockerClient.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		fmt.Println(" can not search docker canisters; msg: ", err.Error())
		panic(err)
	}
	startTelegraf(containers)
	startAgents(containers)
	loopPullAgent()
}

func loopPullAgent() {
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				defer func() {
					if err := recover(); err != nil {
						fmt.Println(" loop pull agent error; msg:", err)
					}
				}()
				containers, err := dockerClient.ContainerList(context.Background(), container.ListOptions{})
				if err != nil {
					fmt.Println(" can not search docker canisters; msg: ", err.Error())
				}
				startAgents(containers)
			}
		}
	}()
}

// 不是正在运行的容器一律删除重新创建启动
func startImage(fileds DockerFileds) error {
	err := pullImage(fileds.Image)
	if err != nil {
		return err
	}
	err = clearContainer(fileds.ContainerName)
	if err != nil {
		return err
	}
	id, err := createContainer(fileds)
	if err != nil {
		return err
	}
	err = startContainer(id)
	if err != nil {
		return err
	}
	fmt.Println("create and run docker canister success; id: ", id)
	return nil
}

func pullImage(imageName string) error {
	ctx := context.Background()
	pull, err := dockerClient.ImagePull(ctx, imageName, image.PullOptions{})
	defer pull.Close()
	if err != nil {
		return err
	}
	return nil
}

func createContainer(dockerfileds DockerFileds) (string, error) {
	ctx := context.Background()
	port := nat.Port(fmt.Sprintf("%d/tcp", dockerfileds.Port))
	hostPort := strconv.FormatInt(dockerfileds.HostPort, 10)
	createResponse, err := dockerClient.ContainerCreate(ctx,
		&container.Config{
			Image: dockerfileds.Image,
			ExposedPorts: nat.PortSet{
				port: {},
			},
			Cmd: strslice.StrSlice{"--gpus=all"},
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				port: []nat.PortBinding{nat.PortBinding{
					HostIP:   "0.0.0.0", //docker容器映射的宿主机的ip
					HostPort: hostPort,  //docker 容器映射到宿主机的端口
				}},
			},
		},
		nil,
		nil,
		dockerfileds.ContainerName)
	if err != nil {
		fmt.Println("create canister err; canister name: ", dockerfileds.ContainerName)
		return "", err
	}
	return createResponse.ID, nil
}

func startContainer(id string) error {
	ctx := context.Background()
	err := dockerClient.ContainerStart(ctx, id, container.StartOptions{})
	if err != nil {
		fmt.Println("failed to start container: ", id)
	}
	return nil
}

func clearContainer(containerName string) error {
	err := dockerClient.ContainerRemove(context.Background(), containerName, container.RemoveOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "No such container") {
			return nil
		}
		return err
	}
	return nil
}
