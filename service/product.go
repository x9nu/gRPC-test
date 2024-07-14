package service

import (
	context "context"

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
