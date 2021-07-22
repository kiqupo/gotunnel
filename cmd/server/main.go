package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"tunnelDemo/pkg/tunnel"
)

const (
	// 测试端口
	tunnelPost  = 8008
	visitPost   = 8007
	// pprof性能监听端口
	monitorPost = ":6060"
)

func main() {
	conf := &tunnel.ServerConfig{
		ConnectCount: 5,
		MaxStream:    5,
	}
	//server := tunnel.ServerTunnel(conf)
	//go server.Run()

	tunnel.ServerTunnel(conf)

	// 模拟分配端口
	tunnel.RegisterController("213", tunnelPost, visitPost)

	pprofMonitor()
}

// pprof性能监听
func pprofMonitor()  {
	log.Println(http.ListenAndServe(monitorPost, nil))
}
