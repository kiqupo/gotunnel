// Package tunnel TODO 1.ControllerManager控制通道连接池管理，2.Controller对象管理
package tunnel

import (
	"fmt"
	"github.com/hashicorp/yamux"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

var (
	clientConn         *net.TCPConn
	connectionPool     map[string]*ConnMatch
	connectionPoolLock sync.Mutex
)

type ConnMatch struct {
	addTime time.Time
	accept  *net.TCPConn
}

type ServerConfig struct {
	// 控制通道监听端口
	ControlPost string

	// 用户请求监听端口
	VisitorPost string

	// 数据通道监听端口
	TunnelPost string
}

func ServerRun(conf *ServerConfig) (err error) {
	defer func() {
		if err = recover().(error); err != nil {
			fmt.Println(err)
		}
	}()

	connectionPool = make(map[string]*ConnMatch, 32)
	go serveControlChannel(conf.ControlPost)
	go userRequestListener(conf.VisitorPost)
	go saTunnelListener(conf.TunnelPost)
	cleanConnectionPool()
	return err
}

// 创建一个控制通道，用于传递控制消息，如：心跳，创建新连接
func serveControlChannel(controlAddr string) {
	tcpListener, err := CreateTCPListener(controlAddr)
	if err != nil {
		panic(err)
	}

	log.Println("[已监听控制通道]" + controlAddr)
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}

		log.Println("[控制通道新连接]" + tcpConn.RemoteAddr().String())
		// 如果当前已经有一个客户端存在，则丢弃这个链接
		if clientConn != nil {
			_ = tcpConn.Close()
		} else {
			clientConn = tcpConn
			go keepAlive(tcpConn)
		}
	}
}

// 和客户端保持一个心跳链接
func keepAlive(conn *net.TCPConn) {
	for {
		if conn == nil {
			return
		}
		_, err := conn.Write(([]byte)(KeepAlive + "\n"))
		if err != nil {
			log.Println("[已断开客户端连接]", conn.RemoteAddr())
			clientConn = nil
			return
		}
		time.Sleep(time.Second * 3)
	}
}

// 监听来自用户的请求
func userRequestListener(visitAddr string) {
	tcpListener, err := CreateTCPListener(visitAddr)
	if err != nil {
		panic(err)
	}
	defer tcpListener.Close()
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			continue
		}
		log.Println("[新用户请求]：" + tcpConn.RemoteAddr().String())
		addConn2Pool(tcpConn)
		// 通知客户端建立通道
		sendMessage(NewConnection + "\n")
	}
}

// 将用户来的连接放入连接池中
func addConn2Pool(accept *net.TCPConn) {
	connectionPoolLock.Lock()
	defer connectionPoolLock.Unlock()

	now := time.Now()
	connectionPool[strconv.FormatInt(now.UnixNano(), 10)] = &ConnMatch{now, accept}
}

// 发送给客户端新消息
func sendMessage(message string) {
	if clientConn == nil {
		log.Println("[无已连接的客户端]")
		return
	}
	_, err := clientConn.Write([]byte(message))
	if err != nil {
		log.Println("[发送消息异常]: message: ", message)
	}
}

// 接收sa来的请求并建立隧道
func saTunnelListener(tunnelAddr string) {
	tcpListener, err := CreateTCPListener(tunnelAddr)
	if err != nil {
		panic(err)
	}
	defer tcpListener.Close()

	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			continue
		}
		go sessionServer(tcpConn)
	}
}

func sessionServer(tcpConn *net.TCPConn)  {
	session, _ := yamux.Server(tcpConn, nil)
	for {
		// 建立多个流通路
		stream, err := session.Accept()
		if err != nil {
			fmt.Println("session over.")
			return
		}
		log.Println("[新客户端请求隧道]：" + tcpConn.RemoteAddr().String())
		go establishTunnel(stream)
	}
}

func establishTunnel(tunnel net.Conn) {
	connectionPoolLock.Lock()
	defer connectionPoolLock.Unlock()

	for key, connMatch := range connectionPool {
		if connMatch.accept != nil {
			go Join2Conn(connMatch.accept, tunnel)
			delete(connectionPool, key)
			return
		}
	}

	_ = tunnel.Close()
}

func cleanConnectionPool() {
	for {
		connectionPoolLock.Lock()
		for key, connMatch := range connectionPool {
			if time.Now().Sub(connMatch.addTime) > time.Second*10 {
				_ = connMatch.accept.Close()
				delete(connectionPool, key)
			}
		}
		connectionPoolLock.Unlock()
		time.Sleep(5 * time.Second)
	}
}
