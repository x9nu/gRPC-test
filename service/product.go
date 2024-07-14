package service

import context "context"

var ProductService = &productService{}

type productService struct {
	UnimplementedProductServiceServer
}

func (p *productService) GetProductStock(ctx context.Context, prodReq *ProductRequest) (*ProductResponse, error) {
	// 实现具体业务逻辑
	stock := p.GetStockByID(prodReq.ProdId)
	user := User{Username: "test_proto_import"}
	return &ProductResponse{ProdStock: stock, User: &user}, nil
}

func (p *productService) GetStockByID(id int32) int32 {
	return 100
}
