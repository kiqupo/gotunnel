package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	_ "net/http/pprof"
	"tunnelDemo/pkg/tunnel"
)

const (
	// 远端的服务控制通道，用来传递控制信息，如出现新连接和心跳
	remoteControlAddr = "127.0.0.1:8009"
	// 远端服务端口，用来建立隧道
	remoteServerAddr  = "127.0.0.1:8008"
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

func runTunnel()  {
	conf := &tunnel.ClientConfig{
		ControllerAddr:remoteControlAddr,
		TunnelAddr:remoteServerAddr,
		LocalServerAddr:localServerAddr,
	}
	err := tunnel.ClientRun(conf)
	if err != nil {
		log.Fatal(err)
	}
}

// 模拟HTTP服务
func runHttp()  {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})
	r.Run(":8080") // 监听并在 0.0.0.0:8080 上启动服务
}

// pprof性能监听
func pprofMonitor()  {
	log.Println(http.ListenAndServe(monitorPost, nil))
}