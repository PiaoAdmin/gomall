package service

import (
	"context"

	apiOrder "github.com/PiaoAdmin/pmall/app/api/biz/model/api/order"
	"github.com/PiaoAdmin/pmall/app/api/md/jwt"
	"github.com/PiaoAdmin/pmall/app/api/rpc"
	"github.com/PiaoAdmin/pmall/common/errs"
	"github.com/PiaoAdmin/pmall/rpc_gen/cart"
	orderrpc "github.com/PiaoAdmin/pmall/rpc_gen/order"
	"github.com/cloudwego/hertz/pkg/app"
)

type PlaceOrderService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewPlaceOrderService(ctx context.Context, c *app.RequestContext) *PlaceOrderService {
	return &PlaceOrderService{RequestContext: c, Context: ctx}
}

func (s *PlaceOrderService) Run(req *apiOrder.PlaceOrderReq) (resp *apiOrder.PlaceOrderResp, err error) {
	claims := jwt.ExtractClaims(s.Context, s.RequestContext)
	userID := uint64(claims[jwt.JwtMiddleware.IdentityKey].(float64))

	// 调用购物车 RPC 服务
	cartResp, err := rpc.CartClient.GetCartDetails(s.Context, &cart.GetCartDetailsRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}
	if len(cartResp.Items) == 0 {
		return nil, errs.New(40007, "cart is empty")
	}
	// map items
	items := make([]*orderrpc.CartItem, 0, len(cartResp.Items))
	for _, it := range cartResp.Items {
		items = append(items, &orderrpc.CartItem{
			SkuId:    it.SkuId,
			Quantity: it.Quantity,
			SkuName:  it.SkuName,
			Price:    it.Price,
		})
	}

	rpcReq := &orderrpc.PlaceOrderReq{
		UserId: userID,
		Email:  req.Email,
		Items:  items,
	}
	if req.ShippingAddress != nil {
		rpcReq.ShippingAddress = &orderrpc.Address{
			Name:          req.ShippingAddress.Name,
			StreetAddress: req.ShippingAddress.StreetAddress,
			City:          req.ShippingAddress.City,
			ZipCode:       req.ShippingAddress.ZipCode,
		}
	}

	rpcResp, err := rpc.OrderClient.PlaceOrder(s.Context, rpcReq)
	if err != nil {
		return nil, err
	}

	resp = &apiOrder.PlaceOrderResp{Order: &apiOrder.OrderResultDTO{OrderId: rpcResp.Order.GetOrderId()}}
	return resp, nil
}
