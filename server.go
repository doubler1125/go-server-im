package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip        string
	Port      int
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	Message chan string
}

// 创建一个server
func NewServer(ip string, port int) *Server {
	// 创建一个server对象
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// 启动服务
func (this *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.listen err:", err)
		return
	}
	defer listener.Close()

	// 启动监听message的goroutine
	go this.ListenMessager()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}

		// do hander
		go this.Handler(conn)
	}
}

func (this *Server) Handler(conn net.Conn) {

	user := NewUser(conn, this)

	user.Online()

	isActive := make(chan bool)

	// 接收客户端消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Printf("Conn read err:", err)
				return
			}

			// 提取用户的消息(去除'\n')
			msg := string(buf[:n-1])

			// 将得到的消息广播
			user.DoMessage(msg)

			isActive <- true
		}
	}()

	// 当前handler阻塞
	for {
		select {
		case <-isActive:
			// 当前用户是活跃的，应该重置定时器；不做任何处理，为了激活select更新下面这个定时器
		case <-time.After(time.Second * 300):
			user.SendMsg("你被踢了")

			// 销毁资源
			close(user.C)

			// 关闭连接
			conn.Close()

			// 退出当前handler
			return // 或者 runtime.Goexit()
		}
	}
}

// 广播消息的方法
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg
}

// 监听Message广播chan的goroutine，一旦有消息就发送给全部的在线User
func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message

		// 将msg发送给全部的在线user
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}
