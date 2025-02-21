package service

import (
	"context"
	cart "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/cart"
	common "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/common"
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/infra/rpc"
	gateutils "github.com/PiaoAdmin/gomall/app/hertz_gateway/utils"
	rpccart "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/cart"
	"github.com/cloudwego/hertz/pkg/app"
)

type AddCartItemService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewAddCartItemService(Context context.Context, RequestContext *app.RequestContext) *AddCartItemService {
	return &AddCartItemService{RequestContext: RequestContext, Context: Context}
}

func (h *AddCartItemService) Run(req *cart.AddCartItemReq) (resp *common.Empty, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	_, err = rpc.CartClient.AddItem(h.Context, &rpccart.AddItemReq{
		UserId: uint32(gateutils.GetUserIdFromCtx(h.Context)),
		Item: &rpccart.CartItem{
			ProductId: req.ProductId,
			Quantity:  req.ProductNum,
		},
	})
	if err != nil {
		return nil, err
	}
	return
}
