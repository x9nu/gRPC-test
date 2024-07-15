package service

import (
	context "context"
	"fmt"
	"io"

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
			if err == io.EOF { //正常处理完成
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
