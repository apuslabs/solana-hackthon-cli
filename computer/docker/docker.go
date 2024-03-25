package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"strconv"
	"strings"
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
	dc, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation(), client.WithHost("tcp://:2375"))
	if err != nil {
		fmt.Println(" Unable to create docker client; msg: ", err.Error())
		panic(err)
	}
	dockerClient = dc
	containers, err := dockerClient.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		fmt.Println(" can not search docker canisters; msg: ", err.Error())
		panic(err)
	}
	startTelegraf(containers)
	startAgents(containers)
	keepalive()
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
	id, err := createImage(fileds)
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

func createImage(dockerfileds DockerFileds) (string, error) {
	ctx := context.Background()
	port := nat.Port(fmt.Sprintf("%d/tcp", dockerfileds.Port))
	hostPort := strconv.FormatInt(dockerfileds.HostPort, 10)
	createResponse, err := dockerClient.ContainerCreate(ctx,
		&container.Config{
			Image: dockerfileds.Image,
			ExposedPorts: nat.PortSet{
				port: {},
			},
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

func keepalive() {
	// channel通信，启动线程设置channel信号，另一个线程监听。hover错误直接return. 防止进程挂掉
	containers, err := dockerClient.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		fmt.Println(" can not search docker canisters; msg: ", err.Error())
	}
	startTelegraf(containers)
	startAgents(containers)
}
