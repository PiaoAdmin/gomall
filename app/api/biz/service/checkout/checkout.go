package checkout

import (
	"context"

	apiCheckout "github.com/PiaoAdmin/pmall/app/api/biz/model/api/checkout"
	"github.com/PiaoAdmin/pmall/app/api/md/jwt"
	"github.com/PiaoAdmin/pmall/app/api/rpc"
	"github.com/PiaoAdmin/pmall/common/errs"
	"github.com/PiaoAdmin/pmall/rpc_gen/cart"
	checkoutrpc "github.com/PiaoAdmin/pmall/rpc_gen/checkout"
	"github.com/cloudwego/hertz/pkg/app"
)

type CheckoutService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewCheckoutService(ctx context.Context, c *app.RequestContext) *CheckoutService {
	return &CheckoutService{RequestContext: c, Context: ctx}
}

func (s *CheckoutService) Run(req *apiCheckout.CheckoutReq) (*apiCheckout.CheckoutResp, error) {
	claims := jwt.ExtractClaims(s.Context, s.RequestContext)
	userID := uint64(claims[jwt.JwtMiddleware.IdentityKey].(float64))

	cartResp, err := rpc.CartClient.GetCartDetails(s.Context, &cart.GetCartDetailsRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	if len(cartResp.Items) == 0 {
		return nil, errs.New(40007, "cart is empty")
	}

	items := make([]*checkoutrpc.CheckoutItem, 0, len(cartResp.Items))
	for _, it := range cartResp.Items {
		items = append(items, &checkoutrpc.CheckoutItem{
			SkuId:    it.SkuId,
			Quantity: it.Quantity,
		})
	}

	rpcReq := &checkoutrpc.CheckoutRequest{
		UserId: userID,
		Items:  items,
	}
	if req.ShippingAddress != nil {
		rpcReq.ShippingAddress = &checkoutrpc.Address{
			Name:          req.ShippingAddress.Name,
			StreetAddress: req.ShippingAddress.StreetAddress,
			City:          req.ShippingAddress.City,
			ZipCode:       req.ShippingAddress.ZipCode,
		}
	}

	rpcResp, err := rpc.CheckoutClient.Checkout(s.Context, rpcReq)
	if err != nil {
		return nil, err
	}

	respItems := make([]*apiCheckout.CheckoutItemDTO, 0, len(rpcResp.Items))
	for _, it := range rpcResp.Items {
		respItems = append(respItems, &apiCheckout.CheckoutItemDTO{
			SkuId:       it.SkuId,
			Quantity:    it.Quantity,
			SkuName:     it.SkuName,
			SkuImage:    it.SkuImage,
			Price:       it.Price,
			MarketPrice: it.MarketPrice,
			SpuId:       it.SpuId,
			SkuSpecData: it.SkuSpecData,
		})
	}

	return &apiCheckout.CheckoutResp{
		OrderId:     rpcResp.OrderId,
		TotalAmount: rpcResp.TotalAmount,
		Items:       respItems,
	}, nil
}
