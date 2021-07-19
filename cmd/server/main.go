package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"tunnelDemo/pkg/tunnel"
)

const (
	controlPost = ":8009"
	tunnelPost  = ":8008"
	visitPost   = ":8007"
	// pprof性能监听端口
	monitorPost = ":6060"
)

func main() {
	go pprofMonitor()

	conf := &tunnel.ServerConfig{
		ControlPost:controlPost,
		VisitorPost:visitPost,
		TunnelPost:tunnelPost,
	}
	err := tunnel.ServerRun(conf)
	if err != nil {
		log.Fatal(err)
	}
}

// pprof性能监听
func pprofMonitor()  {
	log.Println(http.ListenAndServe(monitorPost, nil))
}
