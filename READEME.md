## Golang即时通讯小系统

### 构建基础Server
### 用户上线功能
### 用户消息广播
### 用户业务封装
### 在线用户查询
### 修改用户名
### 超时强踢
### 私聊用户
### 客户端实现
### 客户端模式选择
### 客户端更新用户名
### 公聊模式
### 私聊模式

```
server:
go build -o server main.go server.go user.go
./server

client:
go build -o client client.go
./client
```