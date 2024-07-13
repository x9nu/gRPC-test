package main

import (
	"fmt"
	"grpc-test/service"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	// 添加证书
	creds, errCreds := credentials.NewServerTLSFromFile("./cert/server.crt", "./cert/server.key")
	if errCreds != nil {
		log.Fatal("server:证书生成失败", errCreds)
	}

	rpcServer := grpc.NewServer(grpc.Creds(creds))
	service.RegisterProductServiceServer(rpcServer, service.ProductService)
	listen, err := net.Listen("tcp", ":8002")
	if err != nil {
		log.Fatal("listen error")
	}
	err = rpcServer.Serve(listen)
	defer rpcServer.Stop()
	if err != nil {
		log.Fatal("rpcServer.Serve err:", err)
	}
	fmt.Println("rpcServer.Serve success")
}
