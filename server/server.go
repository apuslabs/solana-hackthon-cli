package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"solana-hackthon-cli/computer/docker"
	"solana-hackthon-cli/computer/monitor"
	"strconv"
)

// 启动web服务，提供查询接口
func Init() {
	r := gin.Default()

	r.GET("/healthCheck", HealthCheckHandler)

	host := fmt.Sprintf("0.0.0.0:%d", 8080)
	if err := r.Run(host); err != nil {
		fmt.Printf("start service failed, err:%v\n", err)
		panic(err)
	}
}

func HealthCheckHandler(c *gin.Context) {
	agent := c.Query("agent")
	if agent == "" {
		c.JSON(http.StatusOK, Response{Code: 400, Msg: "agent must be set", Data: ""})
		return
	}
	health := monitor.GetHealth()
	if !health {
		c.JSON(http.StatusOK, Response{Code: 200, Msg: "", Data: Health{Busy: true, Port: "0"}})
		return
	}
	port := docker.GetPort(agent)
	if port == 0 {
		c.JSON(http.StatusOK, Response{Code: 200, Msg: "", Data: Health{Busy: true, Port: "0"}})
		return
	}
	c.JSON(http.StatusOK, Response{Code: 200, Msg: "", Data: Health{Busy: false, Port: strconv.FormatInt(port, 10)}})

}
