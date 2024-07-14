package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"grpc-test/service"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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

	/*
		实现token认证，需要合法的用户名和密码
		实现一个拦截器
	*/
	var authInterceptor grpc.UnaryServerInterceptor
	authInterceptor = func(
		ctx context.Context,
		req any, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		// 拦截请求，验证token
		err = Auth(ctx)
		if err != nil {
			return
		}
		// 成功处理请求
		return handler(ctx, req)
	}
	rpcServer := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(authInterceptor))
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

func Auth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return fmt.Errorf("metadata not found")
	}
	fmt.Println("md:", md)
	var user string
	var password string
	// 此处不能直接user:=md["user"][0]因为不确定有数据，若没数据，直接读可能会out of range
	if val, ok := md["user"]; ok {
		user = val[0]
	}
	if val, ok := md["password"]; ok {
		password = val[0]
	}
	if user != "admin" || password != "admin" {
		return status.Errorf(codes.Unauthenticated, "token error")
	}
	return nil
}
