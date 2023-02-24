package main

import (
	"net"
	"strings"
)

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

// 给当前用户的客户单发送消息
func (this *User) SendMsg(msg string) {
	this.Conn.Write([]byte(msg))
}

// 用户处理消息的业务
func (this *User) DoMessage(msg string) {

	if msg == "who" {

		this.Server.mapLock.Lock()
		for _, user := range this.Server.OnlineMap {
			sendMsg := "[" + user.Addr + "]" + user.Name + "在线\n"
			this.SendMsg(sendMsg)
		}
		this.Server.mapLock.Unlock()

	} else if len(msg) > 7 && msg[:7] == "rename|" {

		newName := strings.Split(msg, "|")[1]
		if _, ok := this.Server.OnlineMap[newName]; ok {
			this.SendMsg("用户名已被使用\n")
			return
		}

		this.Server.mapLock.Lock()
		delete(this.Server.OnlineMap, this.Name)
		this.Server.OnlineMap[newName] = this
		this.Server.mapLock.Unlock()

		this.Name = newName
		this.SendMsg("您已更新用户名:" + this.Name + "\n")

	} else if len(msg) > 3 && msg[:3] == "to|" {

		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			this.SendMsg("消息格式不正确，请使用\"to|张三|消息内容\"格式")
			return
		}

		remoteUser, ok := this.Server.OnlineMap[remoteName]
		if !ok {
			this.SendMsg("用户不存在")
			return
		}

		sendMsg := strings.Split(msg, "|")[2]
		if sendMsg == "" {
			this.SendMsg("消息内容为空")
			return
		}

		remoteUser.SendMsg(this.Name + "说:" + sendMsg)

	} else {
		this.Server.BroadCast(this, msg)
	}

}
