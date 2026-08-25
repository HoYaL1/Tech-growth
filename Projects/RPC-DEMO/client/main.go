package main

import (
	"log"
	"net/rpc"
)

// 抄server里 阐明 需要传递/接收什么样的数据
type (
	GetUserReq struct {
		Id string `json:"id"`
	}

	GetUserResp struct {
		Id    string
		Nanme string
		Phone string
	}
)

func main() {
	client, err := rpc.Dial("tcp", "localhost:1234") //Dial 建立 TCP 连接，打通客户端 ↔ 服务端的网络管道（只做 1 次）
	//// server的listen监听端口，接收本端Dial发出的TCP连接握手报文
	if err != nil {
		log.Fatal("请求服务器失败", err)
	}

	defer client.Close() //defer会在函数主体内容运行完后执行 释放资源

	//链接成功后需要发起请求 则需要定义好 请求 和 响应 变量
	var (
		rep  = GetUserReq{Id: "2"}
		resp GetUserResp
	)

	err = client.Call("UserServer.GetUser", rep, &resp) // func (client *rpc.Client) Call(serviceMethod string, args any, reply any) error
	//Call 在已经打通的管道里面，发送一次 RPC 函数调用，可以调用很多次
	if err != nil {
		log.Panicln("请求失败", err)
		return
	}

	log.Println(resp)
}
