package main

import (
	"context"
	"fmt"
	"grpc-test/client/service"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	// 添加证书
	creds, errCreds := credentials.NewClientTLSFromFile("../cert/server.crt", "*.x9nu.cn")
	if errCreds != nil {
		log.Fatal("client:creds err", errCreds)
	}
	conn, err := grpc.NewClient(":8002", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal("conn err")
	}
	defer conn.Close()

	prodCli := service.NewProductServiceClient(conn)
	request := &service.ProductRequest{ProdId: 233}
	resp, err := prodCli.GetProductStock(context.Background(), request)
	if err != nil {
		log.Fatal("get stock err", err)
	}
	fmt.Println("调用gRPC方法成功，ProdStock = ", resp.ProdStock)
}
