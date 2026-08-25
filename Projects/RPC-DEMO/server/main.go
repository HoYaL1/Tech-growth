package main

import (
	"errors"
	"log"
	"net"
	"net/rpc"
)

// 定义客户端传过来的请求
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

type UserServer struct {
}

// 服务定义  format: "Service.Method"  UserServer.GetUser
func (*UserServer) GetUser(req GetUserReq, resp *GetUserResp) error {
	if user, ok := users[req.Id]; ok {
		*resp = GetUserResp{
			Id:    user.Id,
			Nanme: user.Name,
			Phone: user.Phone}
		return nil
	}
	return errors.New("用户不存在")
}

// 完成服务的监听工作
func main() {
	//创建服务
	userServer := new(UserServer)
	//服务注册到RPC里
	rpc.Register(userServer)
	//监听，服务器一直监听有没有客户端链接
	listener, error := net.Listen("tcp", ":1234") //network  port
	if error != nil {
		log.Fatal("监听失败", error)
	}

	log.Println("服务启动成功")

	for {
		connection, err := listener.Accept() //connection 链接就相当于在客户端和服务端之间建立了一条双向通信的管道
		// 每个客户端和服务器之前都有一个conn
		if err != nil {
			log.Println("接受客户端链接失败", err)
			continue
		}

		//并发处理客户端请求
		go rpc.ServeConn(connection)
	}
}
