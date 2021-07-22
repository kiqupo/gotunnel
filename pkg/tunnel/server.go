// Package tunnel TODO 1.ControllerManager控制通道连接池管理，2.Controller对象管理
package tunnel

import (
	"errors"
	"github.com/hashicorp/yamux"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

type Server struct {
	bucket sync.Map //map[int]*Controller
	conf   *ServerConfig
}

// Controller 管理SA通道
type Controller struct {
	SAID string

	//与SA通道端口
	TunnelPort int

	//用户请求端口
	UserPort int

	TunnelListener *net.TCPListener
	Connecters     []*Connecter
	CreatTime      time.Time

	UserListener *net.TCPListener
}

type Connecter struct {
	Conn    *net.TCPConn
	Session *yamux.Session
}

type UserConn struct {
	addTime time.Time
	accept  *net.TCPConn
}

type ServerConfig struct {
	// 固定TCP连接数
	ConnectCount int `json:"connect_count"`

	// 每个tcp连接最大stream数
	MaxStream int `json:"max_stream"`
}

func ServerTunnel(conf *ServerConfig) *Server {
	once.Do(func() {
		server = &Server{
			conf: conf,
		}
	})
	return server
}

func RegisterController(saId string, tunnelPort, userPort int) (err error) {
	if server != nil {
		c := &Controller{
			SAID:       saId,
			TunnelPort: tunnelPort,
			UserPort:   userPort,
		}
		err = c.TunnelListen()
		if err != nil {
			return
		}
		err = c.UserListen()
		if err != nil {
			return
		}
		server.bucket.Store(userPort, c)
		return
	}
	return errors.New("TunnelServer单例未初始化")
}

func UnRegisterController(userPort int) error {
	if server != nil {

	}
	return errors.New("单例未初始化")
}

// TunnelListen 监听SA的TCP连接请求
func (c *Controller) TunnelListen() error {
	port := strconv.Itoa(c.TunnelPort)
	tcpListener, err := CreateTCPListener(":" + port)
	if err != nil {
		return err
	}
	if c.TunnelListener == nil {
		c.TunnelListener = tcpListener
	}
	log.Println("[监听通道tcp端口]：" + port)
	go c.ConnListen()
	return err
}

func (c *Controller) ConnListen() {
	for {
		tcpConn, err := c.TunnelListener.AcceptTCP()
		if err != nil {
			log.Println(c.TunnelPort, ":[端口连接]:", tcpConn.RemoteAddr().String(), ":创建session错误:", err)
			continue
		}
		log.Println(c.TunnelPort, ":[端口新TCP连接]:", tcpConn.RemoteAddr().String())

		var session *yamux.Session
		session, err = yamux.Client(tcpConn, nil)
		if err != nil {
			log.Println(c.TunnelPort, ":[端口连接]:", tcpConn.RemoteAddr().String(), ":创建session错误:", err)
			continue
		}
		newconn := &Connecter{
			Conn:    tcpConn,
			Session: session,
		}
		c.Connecters = append(c.Connecters, newconn)
	}
}

// UserListen 监听用户请求
func (c *Controller) UserListen() error {
	port := strconv.Itoa(c.UserPort)
	tcpListener, err := CreateTCPListener(":" + port)
	if err != nil {
		return err
	}
	if c.UserListener == nil {
		c.UserListener = tcpListener
	}
	log.Println("[监听用户tcp端口]：" + port)
	go c.RequestListen()
	return err
}

func (c *Controller) RequestListen() {
	for {
		userConn, err := c.UserListener.AcceptTCP()
		if err != nil {
			continue
		}
		log.Println("[新用户请求]：" + userConn.RemoteAddr().String())
		c.establishTunnel(userConn)
	}
}

func (c *Controller) establishTunnel(reqConn *net.TCPConn) {
	for _, connecter := range c.Connecters {
		if connecter.Session.NumStreams() < server.conf.MaxStream {
			stream, _ := connecter.Session.Open()
			go Join2Conn(reqConn, stream)
			return
		}
	}
}