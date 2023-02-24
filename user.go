package main

import "net"

type User struct {
	Name   string
	Addr   string
	C      chan string
	Conn   net.Conn
	Server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	addr := conn.RemoteAddr().String()
	user := &User{
		Name:   addr,
		Addr:   addr,
		C:      make(chan string),
		Conn:   conn,
		Server: server,
	}

	// 启动监听当前user chan消息的goroutine
	go user.ListenMessage()

	return user
}

// 监听当前user的chan，一旦有消息就发送给对端客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.Conn.Write([]byte(msg + "\n"))
	}
}

// 用户上线
func (this *User) Online() {

	// 用户上线，加入到OnlineMap中
	this.Server.mapLock.Lock()
	this.Server.OnlineMap[this.Name] = this
	this.Server.mapLock.Unlock()

	// 广播用户上线消息
	this.Server.BroadCast(this, "用户上线")
}

// 用户下线
func (this *User) Offline() {

	// 将用户从OnlineMap中删除
	this.Server.mapLock.Lock()
	delete(this.Server.OnlineMap, this.Name)
	this.Server.mapLock.Unlock()

	// 广播用户下线消息
	this.Server.BroadCast(this, "下线")
}

// 用户处理消息的业务
func (this *User) DoMessage(msg string) {
	this.Server.BroadCast(this, msg)
}
