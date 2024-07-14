package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"grpc-test/client/auth"
	"grpc-test/client/service"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	// // 添加证书 - 单向认证
	// creds, errCreds := credentials.NewClientTLSFromFile("../cert/server.crt", "*.x9nu.cn")
	// if errCreds != nil {
	// 	log.Fatal("client:creds err", errCreds)
	// }

	cert, _ := tls.LoadX509KeyPair("../cert/client.pem", "../cert/client.key")
	// 创建一个新的、空的 CertPool
	certPool := x509.NewCertPool()
	ca, _ := os.ReadFile("../cert/ca.crt")
	// 尝试解析所传入的 PEM 编码的证书。如果解析成功会将其加到 CertPool 中，便于后面的使用
	certPool.AppendCertsFromPEM(ca)
	// 构建基于 TLS 的 TransportCredentials 选项
	creds := credentials.NewTLS(&tls.Config{
		// 设置证书链，允许包含一个或多个
		Certificates: []tls.Certificate{cert},
		// 要求必须校验客户端的证书。可以根据实际情况选用以下参数
		ServerName: "*.x9nu.cn",
		RootCAs:    certPool,
	})
	// 此处可以替代成jwt或oath的方式
	token := &auth.Authentication{
		//带&传递它的指针节约资源。如果不带，赋值或作为参数传递时，都会在内存中创建结构体的一个新副本，导致额外的内存消耗
		User:     "admin",
		Password: "admin",
	}

	conn, err := grpc.NewClient(":8002", grpc.WithTransportCredentials(creds), grpc.WithPerRPCCredentials(token))
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
