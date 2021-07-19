package main

import (
	"tunnelDemo/pkg/tunnel"
)

const (
	// 远端的服务控制通道，用来传递控制信息，如出现新连接和心跳
	remoteControlAddr = "127.0.0.1:8009"
	// 远端服务端口，用来建立隧道
	remoteServerAddr  = "127.0.0.1:8008"
	// 本地需要映射的服务端口
	localServerAddr = "127.0.0.1:8080"
)

func main() {
	conf := &tunnel.ClientConfig{
		ControllerAddr:remoteControlAddr,
		TunnelAddr:remoteServerAddr,
		LocalServerAddr:localServerAddr,
	}
	err := tunnel.ClientRun(conf)
	if err != nil {
		panic(err)
	}
}
