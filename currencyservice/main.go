package main

import (
	"currencyservice/handler"
	"fmt"
	"net"
	"strconv"

	pb "currencyservice/proto"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
)

const PORT = 50012
const ADDRESS = "127.0.0.1"

func main() {
	ipport := ADDRESS + ":" + strconv.Itoa(PORT)

	/*
		注册到consul上
		// 初始化consul配置
	*/

	consulConfig := api.DefaultConfig()
	// create consul object

	consulClient, err_consul := api.NewClient(consulConfig)
	if err_consul != nil {
		fmt.Println("consul create object error:", err_consul)
	}

	// 告诉consul 即将注册到服务的消息
	reg := api.AgentServiceRegistration{
		Name:    "currencyservice",
		Port:    PORT,
		Address: ADDRESS,
		Tags:    []string{"currencyservice"},
	}

	// 注册grpc服务到consul上
	err_agent := consulClient.Agent().ServiceRegister(&reg)
	if err_agent != nil {
		fmt.Println("register server error:", err_agent)
	}

	// ----------------------grpc server----------------------
	// 初始化grpc 对象
	grpcServer := grpc.NewServer()
	// 注册服务
	pb.RegisterCurrencyServiceServer(grpcServer, &handler.CurrencyService{})

	// 监听端口
	listen, err := net.Listen("tcp", ipport)
	if err != nil {
		fmt.Println("grpc listen error:", err)
		return
	}
	defer listen.Close()

	// 启动服务
	fmt.Println("grpc server start success")
	err_grpc := grpcServer.Serve(listen)
	if err_grpc != nil {
		fmt.Println("grpc server start error:", err_grpc)
		return
	}
}
