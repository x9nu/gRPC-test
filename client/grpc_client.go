package main

import (
	"context"
	"fmt"
	"grpc-test/client/service"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(":8002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("conn err")
	}
	defer conn.Close()

	prodCli := service.NewProductServiceClient(conn)
	request := &service.ProductRequest{ProdId: 233}
	resp, err := prodCli.GetProductStock(context.Background(), request)
	if err != nil {
		log.Fatal("get stock err")
	}
	fmt.Println("调用gRPC方法成功，ProdStock = ", resp.ProdStock)
}
