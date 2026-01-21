package service

import (
	"context"

	apiOrder "github.com/PiaoAdmin/pmall/app/api/biz/model/api/order"
	"github.com/PiaoAdmin/pmall/app/api/md/jwt"
	"github.com/PiaoAdmin/pmall/app/api/rpc"
	orderrpc "github.com/PiaoAdmin/pmall/rpc_gen/order"
	"github.com/cloudwego/hertz/pkg/app"
)

type ListOrderService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewListOrderService(ctx context.Context, c *app.RequestContext) *ListOrderService {
	return &ListOrderService{RequestContext: c, Context: ctx}
}

func (s *ListOrderService) Run(req *apiOrder.ListOrderReq) (resp *apiOrder.ListOrderResp, err error) {
	claims := jwt.ExtractClaims(s.Context, s.RequestContext)
	userID := uint64(claims[jwt.JwtMiddleware.IdentityKey].(float64))

	rpcReq := &orderrpc.ListOrderReq{UserId: userID}
	rpcResp, err := rpc.OrderClient.ListOrder(s.Context, rpcReq)
	if err != nil {
		return nil, err
	}

	out := &apiOrder.ListOrderResp{Orders: make([]*apiOrder.OrderDTO, 0, len(rpcResp.Orders))}
	for _, o := range rpcResp.Orders {
		dto := &apiOrder.OrderDTO{
			OrderId:   o.OrderId,
			UserId:    o.UserId,
			Status:    o.Status,
			CreatedAt: int32(o.GetCreatedAt()),
		}
		if o.ShippingAddress != nil {
			dto.ShippingAddress = &apiOrder.AddressDTO{
				Name:          o.ShippingAddress.Name,
				StreetAddress: o.ShippingAddress.StreetAddress,
				City:          o.ShippingAddress.City,
				ZipCode:       o.ShippingAddress.ZipCode,
			}
		}
		items := make([]*apiOrder.OrderItem, 0, len(o.Items))
		for _, it := range o.Items {
			items = append(items, &apiOrder.OrderItem{
				SkuId:    it.SkuId,
				SkuName:  it.SkuName,
				Quantity: it.Quantity,
				Price:    it.Price,
			})
		}
		dto.Items = items
		out.Orders = append(out.Orders, dto)
	}

	return out, nil
}
