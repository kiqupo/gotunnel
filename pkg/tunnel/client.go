package tunnel

import (
	"bufio"
	"errors"
	"github.com/hashicorp/yamux"
	"io"
	"log"
	"net"
)

type ClientConfig struct {
	// 控制通道地址
	ControllerAddr string

	// 需要连接的通道地址
	TunnelAddr string

	// 本地服务地址
	LocalServerAddr string
}

func ClientRun(conf *ClientConfig) error{
	ctrlConn, err := CreateTCPConn(conf.ControllerAddr)
	if err != nil {
		log.Println("[连接失败]" + conf.ControllerAddr + err.Error())
		return err
	}
	log.Println("[已连接]" + conf.ControllerAddr)

	reader := bufio.NewReader(ctrlConn)

	remote := connectRemote(conf.TunnelAddr)
	session, _ := yamux.Client(remote, nil)
	for {
		s, err := reader.ReadString('\n')
		if err != nil || err == io.EOF {
			break
		}

		// 当有新连接信号出现时，新建一个tcp连接
		if s == NewConnection+"\n" {
			stream, _ := session.Open()
			go ClientTunnel(conf.LocalServerAddr,stream)
		}
	}

	log.Println("[已断开]" + conf.ControllerAddr)
	return errors.New("控制已经断开")
}

func ClientTunnel(localServerAddr string,conn net.Conn) {
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
	log.Println("[连接本地服务成功]" + remoteTunnelAddr)
	return conn
}