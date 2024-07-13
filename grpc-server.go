package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"grpc-test/service"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	// // 添加证书 - 单向认证
	// creds, errCreds := credentials.NewServerTLSFromFile("./cert/server.crt", "./cert/server.key")
	// if errCreds != nil {
	// 	log.Fatal("server:证书生成失败", errCreds)
	// }

	cert, err := tls.LoadX509KeyPair("cert/server.crt", "cert/server.key")
	if err != nil {
		log.Fatal("证书读取错误", err)
	}
	// 创建一个新的、空的 CertPool
	certPool := x509.NewCertPool()
	ca, err := os.ReadFile("cert/ca.crt")
	if err != nil {
		log.Fatal("ca证书读取错误", err)
	}
	// 尝试解析所传入的 PEM 编码的证书。如果解析成功会将其加到 CertPool 中，便于后面的使用
	certPool.AppendCertsFromPEM(ca)
	// 构建基于 TLS 的 TransportCredentials 选项
	creds := credentials.NewTLS(&tls.Config{
		// 设置证书链，允许包含一个或多个
		Certificates: []tls.Certificate{cert},
		// 要求必须校验客户端的证书。可以根据实际情况选用以下参数
		ClientAuth: tls.RequireAndVerifyClientCert,
		// 设置根证书的集合，校验方式使用 ClientAuth 中设定的模式
		ClientCAs: certPool,
	})

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
