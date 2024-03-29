package computer

import (
	"fmt"
	"os/exec"
	"regexp"
	"solana-hackthon-cli/computer/docker"
	"solana-hackthon-cli/computer/monitor"
	"strings"
)

var cuda_version string

// 检查 nvidia-smi 查看是否有检查docker环境版本
func Init() monitor.GpuNode {
	cuda_version = CheckNvidiaSmi()
	fmt.Println("cuda version: ", cuda_version)
	GpuCardCheck()
	CheckDocker()
	CheckDockerComputer()
	CheckGit()

	docker.Init()
	// 最好等一段时间
	
	node := monitor.Init()
	node.CudaVersion = cuda_version

	monitor.RefreshHealth()
	return node
}

func CheckNvidiaSmi() string {
	informationByte, err := exec.Command("nvidia-smi").Output()
	if err != nil {
		panic("cmd nividia-smi error; msg:" + err.Error())
	}
	info := string(informationByte)
	result := regexp.MustCompile("CUDA Version:\\s+(\\d+(\\.\\d+)+)").FindString(info)
	return regexp.MustCompile("\\d+(\\.\\d+)+").FindString(result)
}

func GpuCardCheck() {
	infomationByte, err := exec.Command("nvidia-smi", "-L").Output()
	if err != nil {
		panic("cmd nividia-smi error; msg:" + err.Error())
	}
	info := string(infomationByte)
	if !strings.Contains(info, "UUID") {
		panic("not found GPU")
	}
}

func CheckDocker() {
	dockerVersionByte, err := exec.Command("docker", "-v").Output()
	if err != nil {
		panic("cmd docker version error; msg:" + err.Error())
	}
	version := regexp.MustCompile("\\d+(\\.\\d+)+").FindString(string(dockerVersionByte))
	fmt.Println("docker version: ", version)
}

func CheckDockerComputer() {
	dockerVersionByte, err := exec.Command("docker-compose", "-v").Output()
	if err != nil {
		panic("cmd docker-compose -v error; msg:" + err.Error())
	}
	version := regexp.MustCompile("v\\d+(\\.\\d+)+").FindString(string(dockerVersionByte))
	fmt.Println("docker-compose version: ", version)
}

func CheckGit() {
	_, err := exec.Command("git", "version").Output()
	if err != nil {
		panic("cmd git version error; msg:" + err.Error())
	}
}
