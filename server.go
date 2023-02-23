package main

import (
	"fmt"
	"net"
	"sync"
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

	user := NewUser(conn)

	// 用户上线，加入到OnlineMap中
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()

	// 广播用户上线消息
	this.BroadCast(user, "用户上线")
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
