package startup

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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
	if ca.Inited() {
		return
	}
	gpuNode.Owner = config.OwnerPubkey
	gpuNode.Id = "null"
	// gpuNode.Sk = ca.GetPubkey().SecretKey
	gpuNode.Price = config.Price
	gpuNode.Endpoint = config.Endpoint
	// http request
	err := register(gpuNode)
	if err != nil {
		fmt.Println("注册机器节点上链失败", err.Error())
		panic(err)
	}
	fmt.Println("注册机器节点上链成功")
}

type RegisterGPUNodeResponse struct {
	Id string `json:"id"`
	Sk string `json:"sk"`
	Tx string `json:"tx"`
}

func register(gpuNode monitor.GpuNode) error {
	jsonData, err := json.Marshal(gpuNode)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/register-gpu-node", config.ServerAddress)
	resp, err := http.Post(url, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		// read the body as text
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		return errors.New(buf.String())
	}
	var registerResponse RegisterGPUNodeResponse
	err = json.NewDecoder(resp.Body).Decode(&registerResponse)
	if err != nil {
		return err
	}
	ca.SaveLocalKey(registerResponse.Id, registerResponse.Sk)
	return nil
}
