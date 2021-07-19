package main

import "tunnelDemo/pkg/tunnel"

const (
	controlAddr = ":8009"
	tunnelAddr  = ":8008"
	visitAddr   = ":8007"
)

func main() {
	conf := &tunnel.ServerConfig{
		ControlPost:controlAddr,
		VisitorPost:visitAddr,
		TunnelPost:tunnelAddr,
	}
	err := tunnel.ServerRun(conf)
	if err != nil {
		panic(err)
	}
}
