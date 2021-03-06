package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"tunnelDemo/pkg/gotunnel"
	"tunnelDemo/test"

	"github.com/gin-gonic/gin"
)

const (
	// 远端服务端口，用来建立隧道
	remoteServerAddr = "127.0.0.1:8008"
	// 本地需要映射的服务端口
	localServerAddr = "127.0.0.1:8080"

	// pprof性能监听端口
	monitorPost = ":6061"
)

func main() {
	go pprofMonitor()
	go runTunnel()
	runHttp()
}

func runTunnel() {
	// 模拟已请求到端口与固定连接数
	conf := &gotunnel.ClientConfig{
		ConnectCount:    5,
		TunnelAddr:      remoteServerAddr,
		LocalServerAddr: localServerAddr,
	}
	gotunnel.ClientRun(conf)
}

// 模拟HTTP服务
func runHttp() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})
	r.GET("/file", test.Download)
	r.GET("/ws", test.WsHandle)
	r.Run(":8080")
}

// pprof性能监听
func pprofMonitor() {
	log.Println(http.ListenAndServe(monitorPost, nil))
}
