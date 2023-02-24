package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	Conn       net.Conn
}

func NewClient(ip string, port int) *Client {

	// 创建客户端对象
	client := &Client{
		ServerIp:   ip,
		ServerPort: port,
	}

	// 链接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return nil
	}

	client.Conn = conn
	return client
}

var serverIP string
var serverPort int

// ./client -p 127.0.0.1 -port 8888
// init会在main之前自动执行
func init() {
	flag.StringVar(&serverIP, "ip", "127.0.0.1", "设置服务器IP地址")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口")
}

func main() {

	// 命令行解析
	flag.Parse()

	client := NewClient(serverIP, serverPort)
	if client == nil {
		fmt.Println("》》》》链接服务器失败")
		return
	}

	fmt.Println(">>>>链接服务器成功")

	// 启动客户端的业务
	select {}

}
