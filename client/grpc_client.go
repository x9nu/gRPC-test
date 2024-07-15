package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"grpc-test/client/auth"
	"grpc-test/client/service"
	"io"
	"log"
	"os"
	"time"

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
	///* 普通RPC start */
	// request := &service.ProductRequest{ProdId: 233}
	// resp, err := prodCli.GetProductStock(context.Background(), request)
	// if err != nil {
	// 	log.Fatal("get stock err", err)
	// }
	// fmt.Println("调用gRPC方法成功，ProdStock = ", resp.ProdStock, resp.User, resp.Data)
	///* 普通RPC end */

	// /* 客户端流RPC start */
	// stream, err := prodCli.UpdateProductStockClientStream(context.Background())
	// if err != nil {
	// 	log.Fatal("获取流出错", err)
	// }
	// rsp := make(chan struct{}, 1)
	// go prodRequest(stream, rsp)
	// select {
	// case <-rsp:
	// 	recv, err := stream.CloseAndRecv()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	stock := recv.ProdStock
	// 	fmt.Println("库存:", stock)
	// }
	// /* 客户端流RPC end */

	// /* 服务端流 RPC start */
	// request := &service.ProductRequest{ProdId: 233}
	// stream, err := prodCli.GetProductStockServerStream(context.Background(), request)
	// if err != nil {
	// 	log.Fatal("获取流出错", err)
	// }
	// for {
	// 	recv, err := stream.Recv()
	// 	if err != nil {
	// 		if err == io.EOF {
	// 			fmt.Println("客户端接收数据完成")
	// 			err := stream.CloseSend()
	// 			if err != nil {
	// 				log.Fatal("关闭流出错", err)
	// 			}
	// 			break
	// 		}
	// 		log.Fatal("获取流出错", err)
	// 	}
	// 	fmt.Println("客户端收到的流", recv.ProdStock)
	// }
	// /* 服务端流 RPC end */

	/* 双向流 RPC start */
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	stream, err := prodCli.HelloBidirectionalStream(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// 原来的for {}中是没有结束的，实际使用必须得有退出策略，比如下面的三种方案：
	// - 设置消息计数：你可以像之前讨论过的那样，设定一个消息计数，当发送和接收的消息达到预定数量时，退出循环。
	// - 使用上下文（context）：你可以使用带有超时或取消能力的上下文，当超时发生或调用取消函数时，退出循环。
	// - 监听服务端的流关闭：当服务端关闭流时，stream.Recv() 将会返回 io.EOF 错误，这可以作为一个退出循环的信号。
	for {
		select {
		case <-ctx.Done():
			// Context has been canceled, exit the loop.
			log.Println("Context canceled, exiting loop.")
			return
		default:
			// Continue with normal operation.
			request := &service.ProductRequest{ProdId: 1}
			err = stream.Send(request)
			if err != nil {
				log.Fatal(err)
			}

			recv, err := stream.Recv()
			if err != nil {
				if err == io.EOF { // Service has closed the stream, exit the loop.
					log.Println("Service closed the stream, exiting loop.")
					return
				}
				log.Fatal(err)
			}
			fmt.Println("客户端收到的流", recv.ProdStock)
		}
		/* 双向流 RPC end */
	}
}

func prodRequest(stream service.ProductService_UpdateProductStockClientStreamClient, rsp chan struct{}) {
	// 模拟 10 个请求
	count := 0
	for {
		request := &service.ProductRequest{ProdId: 233}
		err := stream.Send(request)
		if err != nil {
			log.Fatal("发送流出错", err)
		} else {
			count++
		}
		fmt.Println("发送请求", count)
		if count > 9 {
			rsp <- struct{}{}
			break
		}
	}
}
