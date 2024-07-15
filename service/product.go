package service

import (
	context "context"
	"fmt"
	"io"
	"log"

	"google.golang.org/protobuf/types/known/anypb"
)

var ProductService = &productService{}

type productService struct {
	UnimplementedProductServiceServer
}

func (p *productService) GetProductStock(ctx context.Context, prodReq *ProductRequest) (*ProductResponse, error) {
	// 实现具体业务逻辑
	stock := p.GetStockByID(prodReq.ProdId)
	user := User{Username: "test_proto_import"}
	content := Content{Msg: "any msg"}
	any, _ := anypb.New(&content) // 因为proto文件中ProductResponse指定的data是 any 类型，所以要转成 `anypb` 类型

	return &ProductResponse{ProdStock: stock, User: &user, Data: any}, nil
}

func (p *productService) GetStockByID(id int32) int32 {
	return 100
}

func (p *productService) UpdateProductStockClientStream(stream ProductService_UpdateProductStockClientStreamServer) error {
	count := 0
	// 源源不断地接受客户端发送过来的消息，直到收10个请求就相应一次
	for {
		recv, err := stream.Recv()
		if err != nil {
			if err == io.EOF { // stream.Recv()中接收到 io.EOF，表示客户端不再发送消息并且已经关闭了流
				return nil
			}
			return err
		} else {
			count++
		}
		if count > 9 {
			err := stream.SendAndClose(&ProductResponse{ProdStock: recv.ProdId})
			if err != nil {
				return err
			}
		}
		fmt.Println("服务端接收到的流：", recv.ProdId)
	}
}

func (p *productService) GetProductStockServerStream(req *ProductRequest, stream ProductService_GetProductStockServerStreamServer) error {
	count := 0
	for {
		rsp := &ProductResponse{ProdStock: req.ProdId}
		err := stream.Send(rsp)
		if err != nil {
			return err
		} else {
			count++
		}
		if count > 10 {
			return nil
		}
	}
}

func (p *productService) HelloBidirectionalStream(stream ProductService_HelloBidirectionalStreamServer) error {
	for {
		recv, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				// 客户端关闭了流，可以优雅地退出循环。
				log.Println("客户端关闭了流，服务端结束处理。")
				return nil
			}
			return err
		}
		fmt.Println("服务端接收到客户端的消息", recv.ProdId)
		// 处理业务逻辑...

		// 发送响应
		err = stream.Send(&ProductResponse{ProdStock: recv.ProdId})
		if err != nil {
			return err
		}
	}
}
