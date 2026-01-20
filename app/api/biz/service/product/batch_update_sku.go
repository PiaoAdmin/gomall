package service

import (
	"context"

	apiProduct "github.com/PiaoAdmin/pmall/app/api/biz/model/api/product"
	"github.com/PiaoAdmin/pmall/app/api/rpc"
	"github.com/PiaoAdmin/pmall/rpc_gen/product"
	"github.com/cloudwego/hertz/pkg/app"
)

type BatchUpdateSkuService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewBatchUpdateSkuService(ctx context.Context, c *app.RequestContext) *BatchUpdateSkuService {
	return &BatchUpdateSkuService{
		RequestContext: c,
		Context:        ctx,
	}
}

func (s *BatchUpdateSkuService) Run(req *apiProduct.BatchUpdateSkuRequest) (resp *apiProduct.BatchUpdateSkuResponse, err error) {
	skus := make([]*product.ProductSKU, 0, len(req.Items))
	for _, item := range req.Items {
		skus = append(skus, &product.ProductSKU{
			Id:    item.SkuId,
			Price: item.Price,
			Stock: item.Stock,
		})
	}

	// 调用 RPC BatchUpdateSku
	rpcResp, err := rpc.ProductClient.BatchUpdateSku(s.Context, &product.BatchUpdateSkuRequest{
		Skus: skus,
	})
	if err != nil {
		return nil, err
	}

	return &apiProduct.BatchUpdateSkuResponse{
		Success:      rpcResp.Success,
		UpdatedCount: int32(len(req.Items)),
		Message:      "SKU 更新成功",
	}, nil
}
