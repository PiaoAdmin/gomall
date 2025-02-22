package service

import (
	"context"
	"github.com/PiaoAdmin/gomall/app/cart/biz/dal/mysql"
	"github.com/PiaoAdmin/gomall/app/cart/biz/model"
	cart "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/cart"
	"github.com/cloudwego/kitex/pkg/kerrors"
)

type AddItemService struct {
	ctx context.Context
} // NewAddItemService new AddItemService
func NewAddItemService(ctx context.Context) *AddItemService {
	return &AddItemService{ctx: ctx}
}

// Run create note info
func (s *AddItemService) Run(req *cart.AddItemReq) (resp *cart.AddItemResp, err error) {
	// Finish your business logic.这里是校验代码，因为我不知道商品服务那边怎么写，所以这里注释掉了
	//productResp, err := rpc.ProductCLient.GetProduct(s.ctx, &product.GetProductReq{Id: req.Item.ProductId})
	//if err != nil {
	//	return nil, err
	//}
	//if productResp == nil || productResp.Product.Id == 0 {
	//	return nil, kerrors.NewBizStatusError(40004, "product not found")
	//}
	cartItem := &model.Cart{
		UserId:    req.UserId,
		ProductId: req.Item.ProductId,
		Qty:       req.Item.Quantity,
	}
	err = model.AddItem(s.ctx, mysql.DB, cartItem)
	if err != nil {
		return nil, kerrors.NewBizStatusError(50000, err.Error())
	}
	return &cart.AddItemResp{}, nil
}
