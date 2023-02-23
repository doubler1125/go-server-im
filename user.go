package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	Conn net.Conn
}

func NewUser(conn net.Conn) *User {
	addr := conn.RemoteAddr().String()
	user := &User{
		Name: addr,
		Addr: addr,
		C:    make(chan string),
		Conn: conn,
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
