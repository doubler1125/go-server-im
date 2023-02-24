package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
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

	client.conn = conn
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

	// 单独开启一个goroutine去处理server的回执消息
	go client.DealResponse()

	// 启动客户端的业务
	client.run()
}

func (client *Client) run() {
	for client.flag != 0 {
		// 直到选一个模式为止
		for client.menu() != true {
		}

		switch client.flag {
		case 1:
			client.PublicChat()
			break
		case 2:
			fmt.Println("选择私聊模式")
			break
		case 3:
			client.updateName()
			break
		}
	}
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

func (client *Client) updateName() bool {

	fmt.Println("请输入用户名")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("client.conn.write err:", err)
		return false
	}
	return true
}

// 处理server返回的消息，直接显示到标准输出即可
func (client *Client) DealResponse() {

	// 一旦client.conn有数据，直接拷贝到stdout标准输出上，永久阻塞监听
	io.Copy(os.Stdout, client.conn)
	// 相当于下面的写法
	// for {
	// 	buf := make([]byte, 4096)
	// 	client.conn.Read(buf)
	// 	fmt.Println(buf)
	// }
}

func (client *Client) PublicChat() {
	var chatMsg string

	fmt.Println(">>>>>请输入聊天内容,输入exit退出")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		// 发送给服务端
		if len(chatMsg) != 0 {
			msg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(msg))
			if err != nil {
				fmt.Println("client.conn.write err:", err)
				break
			}
		}

		chatMsg = ""
		fmt.Println("请输入聊天内容,输入exit退出")
		fmt.Scanln(&chatMsg)
	}
}
