package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"tunnelDemo/pkg/gotunnel"
)

const (
	// 测试端口
	tunnelPost = 8008
	visitPost  = 8007
	// pprof性能监听端口
	monitorPost = ":6060"
)

func main() {
	conf := &gotunnel.ServerConfig{
		ConnectCount: 5,
		MaxStream:    5,
	}

	gotunnel.ServerTunnel(conf)

	// 模拟分配端口
	gotunnel.RegisterController("213", tunnelPost, visitPost)

	pprofMonitor()
}

// pprof性能监听
func pprofMonitor() {
	log.Println(http.ListenAndServe(monitorPost, nil))
}
