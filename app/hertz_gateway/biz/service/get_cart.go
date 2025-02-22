package service

import (
	"context"
	"strconv"

	common "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/common"
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/infra/rpc"
	gateutils "github.com/PiaoAdmin/gomall/app/hertz_gateway/utils"
	"github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/cart"
	rpcproduct "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/product"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

type GetCartService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetCartService(Context context.Context, RequestContext *app.RequestContext) *GetCartService {
	return &GetCartService{RequestContext: RequestContext, Context: Context}
}

func (h *GetCartService) Run(req *common.Empty) (resp map[string]any, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	cartResp, err := rpc.CartClient.GetCart(h.Context, &cart.GetCartReq{
		UserId: uint32(gateutils.GetUserIdFromCtx(h.Context)),
	})
	if err != nil {
		return nil, err
	}
	var items []map[string]string
	var total float64
	for _, item := range cartResp.Cart.Items {
		productResp, err := rpc.ProductClient.GetProduct(h.Context, &rpcproduct.GetProductReq{Id: int64(item.ProductId)})
		if err != nil {
			continue
		}
		p := productResp.Product
		items = append(items, map[string]string{
			// "Name":        p.Name,
			// "Description": p.Description,
			"Price": strconv.FormatFloat(float64(p.Price), 'f', 2, 64),
			// "Picture":     p.Picture,
			"Qty": strconv.Itoa(int(item.Quantity)),
		})
		total += float64(p.Price) * float64(item.Quantity)
	}
	return utils.H{
		"tilte": "Cart",
		"items": items,
		"total": strconv.FormatFloat(float64(total), 'f', 2, 64),
	}, nil

}
