package computer

import (
	"fmt"
	"os/exec"
)

var cuda_version string

// 检查 nvidia-smi 查看是否有检查docker环境版本
func Init() {
	CheckNvidiaSmi()
	CheckDocker()
	cuda_version = GetCudaVersion()
	RegisterNode()
}

func CheckNvidiaSmi() {
	_, err := exec.Command("nvidia-smi").Output()
	if err != nil {
		panic("nividia-smi error; msg:" + err.Error())
	}
}

func GetCudaVersion() string {
	version, err := exec.Command("nvidia-smi", "--query-gpu", "cuda_version").Output()
	if err != nil {
		panic("nividia-smi error; msg:" + err.Error())
	}
	fmt.Printf("cuda_version: %s\n", string(version))
	return string(version)
}

func GetGpuCards() []string {
	_, err := exec.Command("nvidia-smi").Output()
	if err != nil {
		panic("nividia-smi error; msg:" + err.Error())
	}
	return []string{}
}

func CheckDocker() {
	_, err := exec.Command("docker", "version").Output()
	if err != nil {
		panic("docker version error; msg:" + err.Error())
	}
}

func RegisterNode() {

}
