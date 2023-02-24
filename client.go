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
	flag       int // 客户端模式
}

func NewClient(ip string, port int) *Client {

	// 创建客户端对象
	client := &Client{
		ServerIp:   ip,
		ServerPort: port,
		flag:       999,
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

func (this *Client) menu() bool {
	var flag int

	fmt.Sprintln("1.公聊模式")
	fmt.Sprintln("2.私聊模式")
	fmt.Sprintln("3.更新用户名")
	fmt.Sprintln("0.退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		this.flag = flag
		return true
	} else {
		fmt.Println("请输入合法范围内的数字")
		return false
	}
}

func (client *Client) run() {
	for client.flag != 0 {
		// 直到选一个模式为止
		for client.menu() != true {
		}

		switch client.flag {
		case 1:
			fmt.Println("选择公聊模式")
			break
		case 2:
			fmt.Println("选择私聊模式")
			break
		case 3:
			fmt.Println("更新用户名")
			break
		}
	}
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
	client.run()
}
