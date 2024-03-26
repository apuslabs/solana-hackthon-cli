package monitor

import (
	"context"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"log"
)

// 监控计算机节点性能，定时更新健康状态
const token = "1B2G7-MFj2ee9FfoR1YQWG_jVUbCtNP8k4Nhh8h3cCnR-U9n0zta0oFzz3bRUPGQR_edXHbDn1Tx4g0buZEBBQ=="
const url = "http://180.166.208.2:8086"
const org = "apus"
const bucket = "gpu"

var client influxdb2.Client

var health = false

func GetHealth() bool {
	return health
}

func Init() {
	client = influxdb2.NewClient(url, token)
	queryAPI := client.QueryAPI(org)
	query := `from(bucket: "gpu")
            |> range(start: -10m)
            |> filter(fn: (r) => r._measurement == "measurement1")`
	_, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		panic(err)
	}
}

func RefreshHealth() {

}

func GetCpu() {
	queryAPI := client.QueryAPI(org)
	query := `from(bucket: "gpu")
            |> range(start: -10m)
            |> filter(fn: (r) => r._measurement == "measurement1")`
	results, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	for results.Next() {
		fmt.Println(results.Record())
	}
	if err := results.Err(); err != nil {
		log.Fatal(err)
	}
}

func GetMemory() {

}

func GetDisk() {

}

func GetGpu() {

}
