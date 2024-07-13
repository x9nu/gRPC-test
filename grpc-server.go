package main

import (
	"fmt"
	"grpc-test/service"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	rpcServer := grpc.NewServer()
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
