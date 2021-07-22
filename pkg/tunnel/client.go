package tunnel

import (
	"fmt"
	"github.com/hashicorp/yamux"
	"log"
	"net"
)

type Client struct {
	conf   *ClientConfig
}

type ClientConfig struct {
	// 固定TCP连接数
	ConnectCount int `json:"connect_count"`

	// 需要连接的通道地址
	TunnelAddr string

	// 本地服务地址
	LocalServerAddr string
}

func ClientTunnel(conf *ClientConfig) *Client {
	once.Do(func() {
		client = &Client{
			conf:   conf,
		}
	})
	return client
}

func (c *Client)ConnectTunnel()  {
	remote := connectRemote(c.conf.TunnelAddr)
	session, _ := yamux.Server(remote, nil)
	for {
		// 建立多个流通路
		stream, err := session.Accept()
		if err != nil {
			fmt.Println("session over:",err)
			break
		}
		log.Println("[stream Accept]：")
		go clientTunnel(c.conf.LocalServerAddr, stream)
	}
}

func ClientRun(conf *ClientConfig) {
	client := ClientTunnel(conf)
	for i := 0; i < conf.ConnectCount; i++ {
		go client.ConnectTunnel()
	}
}

func clientTunnel(localServerAddr string,conn net.Conn) {
	local := connectLocal(localServerAddr)

	if local != nil && conn != nil {
		Join2Conn(local, conn)
	} else {
		if local != nil {
			_ = local.Close()
		}
		if conn != nil {
			_ = conn.Close()
		}
	}
}

func connectLocal(localServerAddr string) *net.TCPConn {
	conn, err := CreateTCPConn(localServerAddr)
	if err != nil {
		log.Println("[连接本地服务失败]" + err.Error())
	}
	log.Println("[连接本地服务成功]：" + localServerAddr)
	return conn
}

func connectRemote(remoteTunnelAddr string) *net.TCPConn {
	conn, err := CreateTCPConn(remoteTunnelAddr)
	if err != nil {
		log.Println("[连接远端服务失败]" + err.Error())
	}
	log.Println("[连接SC服务成功]" + remoteTunnelAddr)
	return conn
}